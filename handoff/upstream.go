package handoff

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"

	"stathat.com/c/consistent"
)

var (
	upstreamConfig HandoffOpts
)

type backend struct {
	name    string
	streams upstream
	sched   scheduler
}

func (b *backend) init(o BackendOpt) error {

	for node, addr := range o.Upstream {
		ch := make(chan *specs.MetaData)
		b.sched.addChan(node, ch)
		b.streams.addClientChan(node, addr, ch)
	}
	return nil
}

/*   scheduler */
type scheduler interface {
	sched(string) chan *specs.MetaData
	addChan(string, chan *specs.MetaData) error
}

type schedConsistent struct {
	name    string
	consist *consistent.Consistent
	chans   map[string]chan *specs.MetaData
}

func newSchedConsistent(replicas int) *schedConsistent {
	sched := &schedConsistent{
		name:    "consistent",
		consist: consistent.New(),
		chans:   make(map[string]chan *specs.MetaData),
	}
	sched.consist.NumberOfReplicas = replicas
	return sched
}

func (s *schedConsistent) addChan(key string,
	ch chan *specs.MetaData) error {
	s.consist.Add(key)
	s.chans[key] = ch
	return nil
}

func (s *schedConsistent) sched(key string) chan *specs.MetaData {
	node, _ := s.consist.Get(key)
	return s.chans[node]
}

/* upstream */
type rpcClients struct {
	addr string
	cli  []*rpc.Client
}

type netClients struct {
	addr string
	cli  []net.Conn
}

type upstream interface {
	run(string) error
	addClientChan(string, string, chan *specs.MetaData) error
}

type upstreamFalcon struct {
	name        string
	concurrency int
	connTimeout int
	callTimeout int
	clients     map[string]rpcClients
	chans       map[string]chan *specs.MetaData
}

func newUpstreamFalcon(concurrency int,
	b BackendOpt) *upstreamFalcon {
	return &upstreamFalcon{
		name:        "falcon",
		concurrency: concurrency,
		connTimeout: b.ConnTimeout,
		callTimeout: b.CallTimeout,
		clients:     make(map[string]rpcClients),
		chans:       make(map[string]chan *specs.MetaData),
	}
}

func (tpl *upstreamFalcon) dial(address string,
	timeout int) (*rpc.Client, error) {
	return rpcDial(address, time.Duration(timeout)*time.Millisecond)
}

func (p *upstreamFalcon) addClientChan(key, addr string,
	ch chan *specs.MetaData) (err error) {
	p.chans[key] = ch
	p.clients[key] = rpcClients{
		addr: addr,
		cli:  make([]*rpc.Client, p.concurrency),
	}
	for i := 0; i < p.concurrency; i++ {
		p.clients[key].cli[i], err = p.dial(addr, p.connTimeout)
		if err != nil {
			glog.Fatalf("node:%s addr:%s err:%s\n",
				key, addr, err)
		}
	}
	return nil
}

func (p *upstreamFalcon) run(name string) error {
	for node, ch := range p.chans {
		for i := 0; i < len(p.clients[node].cli); i++ {
			go falconUpstreamWorker(name, i,
				ch, &p.clients[node].cli[i],
				p.clients[node].addr,
				p.connTimeout, p.callTimeout)
		}
	}
	return nil
}

func rpcDial(address string, timeout time.Duration) (*rpc.Client, error) {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return rpc.NewClient(conn), err
}

func reconnection(client **rpc.Client, addr string, connTimeout int) {
	var err error

	statInc(ST_CONN_ERR, 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = rpcDial(addr, time.Duration(connTimeout)*time.Millisecond)
	statInc(ST_CONN_DIAL, 1)

	for err != nil {
		//danger!! block routine
		time.Sleep(time.Millisecond * 500)
		*client, err = rpcDial(addr,
			time.Duration(connTimeout)*time.Millisecond)
		statInc(ST_CONN_DIAL, 1)
	}
}

func netRpcCall(client *rpc.Client, method string, args interface{},
	reply interface{}, timeout time.Duration) error {
	done := make(chan *rpc.Call, 1)
	client.Go(method, args, reply, done)
	select {
	case <-time.After(timeout):
		return errors.New("i/o timeout[rpc]")
	case call := <-done:
		if call.Error == nil {
			return nil
		} else {
			return call.Error
		}
	}
}

func putRpcStorageData(client **rpc.Client, item *specs.MetaData,
	addr string, connTimeout, callTimeout int) error {
	var (
		err  error
		resp *specs.RpcResp
		i    int
	)

	resp = &specs.RpcResp{}

	for i = 0; i < CONN_RETRY; i++ {
		err = netRpcCall(*client, "Storage.Put", item, resp,
			time.Duration(callTimeout)*time.Millisecond)

		if err == nil {
			goto out
		}
		if err == rpc.ErrShutdown {
			reconnection(client, addr, connTimeout)
		}
	}
out:
	return err
}

func falconUpstreamWorker(name string, idx int, ch chan *specs.MetaData,
	client **rpc.Client, addr string, connTimeout, callTimeout int) {
	var err error
	for {
		select {
		case item := <-ch:
			if err = putRpcStorageData(client, item,
				addr, connTimeout, callTimeout); err != nil {
				statInc(ST_PUT_ERR, 1)
			} else {
				statInc(ST_PUT_SUCCESS, 1)
			}
		}
	}
}

// TSDB
type upstreamTsdb struct {
	upstreamFalcon
	name        string
	concurrency int
	connTimeout int
	callTimeout int
	clients     map[string]netClients
	chans       map[string]chan *specs.MetaData
}

func newUpstreamTsdb(concurrency int, b BackendOpt) *upstreamTsdb {
	return &upstreamTsdb{
		name:        "tsdb",
		concurrency: concurrency,
		connTimeout: b.ConnTimeout,
		callTimeout: b.CallTimeout,
		clients:     make(map[string]netClients),
		chans:       make(map[string]chan *specs.MetaData),
	}
	return nil
}

func netDial(address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeout}
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		if err := tc.SetKeepAlive(true); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return conn, err
}

func (tpl *upstreamTsdb) dial(address string, timeout int) (net.Conn, error) {
	return netDial(address, time.Duration(timeout)*time.Millisecond)
}

func (p *upstreamTsdb) addClientChan(key, addr string,
	ch chan *specs.MetaData) (err error) {
	p.chans[key] = ch
	p.clients[key] = netClients{
		addr: addr,
		cli:  make([]net.Conn, p.concurrency),
	}
	for i := 0; i < p.concurrency; i++ {
		p.clients[key].cli[i], err = p.dial(addr, p.connTimeout)
		if err != nil {
			glog.Fatalf("node:%s addr:%s err:%s\n", key, addr, err)
		}
	}
	return nil
}

func (p *upstreamTsdb) run(name string) error {
	for node, ch := range p.chans {
		for i := 0; i < len(p.clients[node].cli); i++ {
			go tsdbUpstreamWorker(name, i, ch,
				&p.clients[node].cli[i],
				p.clients[node].addr,
				p.connTimeout,
				p.callTimeout)
		}
	}
	return nil
}

func tsdbUpstreamWorker(name string, idx int,
	ch chan *specs.MetaData, client *net.Conn,
	addr string, connTimeout, callTimeout int) {
	var err error
	for {
		select {
		case item := <-ch:
			if err = putTsdbData(client, item,
				addr, connTimeout, callTimeout); err != nil {
				statInc(ST_PUT_ERR, 1)
			} else {
				statInc(ST_PUT_SUCCESS, 1)
			}
		}
	}
}

func putTsdbData(client *net.Conn, item *specs.MetaData,
	addr string, connTimeout, callTimeout int) (err error) {
	var (
		i          int
		tsdbBuffer bytes.Buffer
	)

	tsdbBuffer.WriteString(item.Tsdb().TsdbString())
	tsdbBuffer.WriteString("\n")

	for i = 0; i < CONN_RETRY; i++ {
		err = tsdbSend(*client, tsdbBuffer.Bytes(),
			time.Duration(callTimeout)*time.Millisecond)

		if err == nil {
			goto out
		} else {
			tsdbReconnection(client, addr, connTimeout)
		}
	}
out:
	return err

}

func tsdbSend(client net.Conn, data []byte, timeout time.Duration) (err error) {
	done := make(chan error)
	go func() {
		_, err = client.Write(data)
		done <- err
	}()

	select {
	case <-time.After(timeout):
		return errors.New("i/o timeout[tsdb]")
	case err = <-done:
		if err != nil {
			return fmt.Errorf("call failed, err %v", err)
		}
		return nil
	}
}

func tsdbReconnection(client *net.Conn, addr string, connTimeout int) {
	var err error

	statInc(ST_CONN_ERR, 1)
	if *client != nil {
		(*client).Close()
	}

	*client, err = netDial(addr,
		time.Duration(connTimeout)*time.Millisecond)
	statInc(ST_CONN_DIAL, 1)

	for err != nil {
		//danger!! block routine
		time.Sleep(time.Millisecond * 500)
		*client, err = netDial(addr,
			time.Duration(connTimeout)*time.Millisecond)
		statInc(ST_CONN_DIAL, 1)
	}
}

func upstreamWorker(bs *[]*backend) error {
	for _, b := range *bs {
		b.streams.run(b.name)
	}
	return nil
}

func loadBalancerWorker(bs *[]*backend) {
	go func() {
		for {
			items := <-appUpdateChan
			for _, b := range *bs {
				for _, item := range *items {
					ch := b.sched.sched(item.Key())
					ch <- item
				}
			}
		}
	}()
}

func upstreamStart(config HandoffOpts) {
	upstreamConfig = config
	bs := make([]*backend, 0)
	for k, v := range upstreamConfig.Backends {
		if !v.Enable {
			continue
		}
		b := &backend{}
		b.name = k
		if v.Type == "falcon" {
			b.streams = newUpstreamFalcon(upstreamConfig.Concurrency, v)
		} else if v.Type == "tsdb" {
			b.streams = newUpstreamTsdb(upstreamConfig.Concurrency, v)
		} else {
			glog.Fatal(specs.ErrUnsupported)
		}

		if v.Sched == "consistent" {
			b.sched = newSchedConsistent(upstreamConfig.Replicas)
		} else {
			glog.Fatal(specs.ErrUnsupported)
		}

		b.init(v)
		bs = append(bs, b)
	}

	upstreamWorker(&bs)

	loadBalancerWorker(&bs)

}
