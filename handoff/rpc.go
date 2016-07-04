package handoff

import (
	"container/list"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	rpcEvent    *specs.RoutineEvent
	rpcConnects connList
	rpcConfig   HandoffOpts
)

type connList struct {
	sync.RWMutex
	list *list.List
}

func (l *connList) insert(c net.Conn) *list.Element {
	l.Lock()
	defer l.Unlock()
	return l.list.PushBack(c)
}
func (l *connList) remove(e *list.Element) net.Conn {
	l.Lock()
	defer l.Unlock()
	return l.list.Remove(e).(net.Conn)
}

type falcon int

func (this *falcon) Ping(req specs.Null, resp *specs.RpcResp) error {
	return nil
}

func (t *falcon) Update(args []*specs.MetaData,
	reply *specs.HandoffResp) error {
	reply.Invalid = 0
	now := time.Now().Unix()

	items := []*specs.MetaData{}
	for _, v := range args {
		if v == nil {
			reply.Invalid += 1
			continue
		}

		if v.K == "" || v.Host == "" {
			reply.Invalid += 1
			continue
		}

		if v.Type != specs.COUNTER &&
			v.Type != specs.GAUGE &&
			v.Type != specs.DERIVE {
			reply.Invalid += 1
			continue
		}

		if v.Step <= 0 {
			reply.Invalid += 1
			continue
		}

		if len(v.K)+len(v.Tags) > 510 {
			reply.Invalid += 1
			continue
		}

		if v.Ts <= 0 || v.Ts > now*2 {
			v.Ts = now
		}

		items = append(items, &specs.MetaData{
			K:    v.K,
			Host: v.Host,
			Ts:   v.Ts,
			Step: v.Step,
			Type: v.Type,
			Tags: sortTags(v.Tags),
		})
	}

	appUpdateChan <- &items

	reply.Message = "ok"
	reply.Total = len(args)

	statInc(ST_RPC_UPDATE, 1)
	statInc(ST_RPC_UPDATE_CNT, len(items))
	statInc(ST_RPC_UPDATE_ERR, reply.Invalid)

	return nil
}

func _rpcStart(config *HandoffOpts,
	listener **net.TCPListener) (err error) {
	var addr *net.TCPAddr

	if !config.Rpc {
		return nil
	}

	addr, err = net.ResolveTCPAddr("tcp", config.RpcAddr)
	if err != nil {
		glog.Fatalf("rpc.Start error, net.ResolveTCPAddr failed, %s", err)
	}

	*listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		glog.Fatalf("rpc.Start error, listen %s failed, %s",
			config.RpcAddr, err)
	} else {
		glog.Infof("rpc.Start ok, listening on %s", config.RpcAddr)
	}

	go func() {
		var tempDelay time.Duration // how long to sleep on accept failure
		for {
			conn, err := (*listener).Accept()
			if err != nil {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			tempDelay = 0
			go func() {
				e := rpcConnects.insert(conn)
				defer rpcConnects.remove(e)
				jsonrpc.ServeConn(conn)
			}()
		}
	}()
	return err
}

func _rpcStop(config *HandoffOpts,
	listener *net.TCPListener) (err error) {
	if listener == nil {
		return specs.ErrNoent
	}

	listener.Close()
	rpcConnects.Lock()
	for e := rpcConnects.list.Front(); e != nil; e = e.Next() {
		e.Value.(net.Conn).Close()
	}
	rpcConnects.Unlock()

	return nil
}

func rpcStart(config HandoffOpts) {
	var rpcListener *net.TCPListener

	rpc.Register(new(falcon))
	registerEventChan(rpcEvent)
	rpcConfig = config

	_rpcStart(&rpcConfig, &rpcListener)

	go func() {
		select {
		case event := <-rpcEvent.E:
			if event.Method == specs.ROUTINE_EVENT_M_EXIT {
				_rpcStop(&rpcConfig, rpcListener)

				return
			} else if event.Method == specs.ROUTINE_EVENT_M_RELOAD {
				_rpcStop(&rpcConfig, rpcListener)

				glog.V(3).Infof("old:\n%s\n new:\n%s",
					rpcConfig, appConfig)
				rpcConfig = appConfig
				_rpcStart(&rpcConfig, &rpcListener)
			}
		}
	}()
}
