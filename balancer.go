package falcon

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/naming"
)

type Watcher struct {
	// the channel to receives name resolution updates
	Update chan *naming.Update
	// the side channel to get to know how many updates in a batch
	Side chan int
	// the channel to notifiy update injector that the update reading is done
	readDone chan int
}

func (w *Watcher) Next() (updates []*naming.Update, err error) {
	n := <-w.Side
	if n == 0 {
		return nil, fmt.Errorf("w.side is closed")
	}
	for i := 0; i < n; i++ {
		u := <-w.Update
		if u != nil {
			updates = append(updates, u)
		}
	}
	w.readDone <- 0
	return
}

func (w *Watcher) Close() {
}

// Inject naming resolution updates to the testWatcher.
func (w *Watcher) Inject(updates []*naming.Update) {
	w.Side <- len(updates)
	for _, u := range updates {
		w.Update <- u
	}
	<-w.readDone
}

type NameResolver struct {
	W     *Watcher
	Addrs []string
	next  int
}

func (r *NameResolver) Resolve(target string) (naming.Watcher, error) {
	r.W = &Watcher{
		Update:   make(chan *naming.Update, 1),
		Side:     make(chan int, 1),
		readDone: make(chan int),
	}

	r.W.Side <- 1
	r.W.Update <- &naming.Update{
		Op:   naming.Add,
		Addr: r.Addrs[r.next],
	}
	r.next++
	if r.next == len(r.Addrs) {
		r.next = 0
	}

	go func() {
		<-r.W.readDone
	}()
	return r.W, nil
}

// RoundRobin returns a Balancer that selects addresses round-robin. It uses r to watch
// the name resolution updates and updates the addresses available correspondingly.
func RoundRobin(r naming.Resolver) grpc.Balancer {
	return &roundRobin{r: r}
}

type addrInfo struct {
	addr      grpc.Address
	connected bool
}

type roundRobin struct {
	r      naming.Resolver
	w      naming.Watcher
	addrs  []*addrInfo // all the addresses the client should potentially connect
	mu     sync.Mutex
	addrCh chan []grpc.Address // the channel to notify gRPC internals the list of addresses the client should connect to.
	next   int                 // index of the next address to return for Get()
	waitCh chan struct{}       // the channel to block when there is no connected address available
	done   bool                // The Balancer is closed.
}

func (rr *roundRobin) watchAddrUpdates() error {
	updates, err := rr.w.Next()
	if err != nil {
		grpclog.Warningf("grpc: the naming watcher stops working due to %v.", err)
		return err
	}
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, update := range updates {
		addr := grpc.Address{
			Addr:     update.Addr,
			Metadata: update.Metadata,
		}
		switch update.Op {
		case naming.Add:
			var exist bool
			for _, v := range rr.addrs {
				if addr == v.addr {
					exist = true
					grpclog.Infoln("grpc: The name resolver wanted to add an existing address: ", addr)
					break
				}
			}
			if exist {
				continue
			}
			rr.addrs = append(rr.addrs, &addrInfo{addr: addr})
		case naming.Delete:
			for i, v := range rr.addrs {
				if addr == v.addr {
					copy(rr.addrs[i:], rr.addrs[i+1:])
					rr.addrs = rr.addrs[:len(rr.addrs)-1]
					break
				}
			}
		default:
			grpclog.Errorln("Unknown update.Op ", update.Op)
		}
	}
	// Make a copy of rr.addrs and write it onto rr.addrCh so that gRPC internals gets notified.
	open := make([]grpc.Address, len(rr.addrs))
	for i, v := range rr.addrs {
		open[i] = v.addr
	}
	if rr.done {
		return grpc.ErrClientConnClosing
	}
	select {
	case <-rr.addrCh:
	default:
	}
	rr.addrCh <- open
	return nil
}

func (rr *roundRobin) Start(target string, config grpc.BalancerConfig) error {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if rr.done {
		return grpc.ErrClientConnClosing
	}
	if rr.r == nil {
		// If there is no name resolver installed, it is not needed to
		// do name resolution. In this case, target is added into rr.addrs
		// as the only address available and rr.addrCh stays nil.
		rr.addrs = append(rr.addrs, &addrInfo{addr: grpc.Address{Addr: target}})
		return nil
	}
	w, err := rr.r.Resolve(target)
	if err != nil {
		return err
	}
	rr.w = w
	rr.addrCh = make(chan []grpc.Address, 1)
	go func() {
		for {
			if err := rr.watchAddrUpdates(); err != nil {
				return
			}
		}
	}()
	return nil
}

// Up sets the connected state of addr and sends notification if there are pending
// Get() calls.
func (rr *roundRobin) Up(addr grpc.Address) func(error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	var cnt int
	for _, a := range rr.addrs {
		if a.addr == addr {
			if a.connected {
				return nil
			}
			a.connected = true
		}
		if a.connected {
			cnt++
		}
	}
	// addr is only one which is connected. Notify the Get() callers who are blocking.
	if cnt == 1 && rr.waitCh != nil {
		close(rr.waitCh)
		rr.waitCh = nil
	}
	return func(err error) {
		rr.down(addr, err)
	}
}

// down unsets the connected state of addr.
func (rr *roundRobin) down(addr grpc.Address, err error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	for _, a := range rr.addrs {
		if addr == a.addr {
			a.connected = false
			break
		}
	}
}

// Get returns the next addr in the rotation.
func (rr *roundRobin) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	var ch chan struct{}
	rr.mu.Lock()
	if rr.done {
		rr.mu.Unlock()
		err = grpc.ErrClientConnClosing
		return
	}

	if len(rr.addrs) > 0 {
		if rr.next >= len(rr.addrs) {
			rr.next = 0
		}
		next := rr.next
		for {
			a := rr.addrs[next]
			if a.connected {
				addr = a.addr
				rr.next = next
				rr.mu.Unlock()
				return
			}
			next = (next + 1) % len(rr.addrs)
			if next == rr.next {
				// Has iterated all the possible address but none is connected.
				break
			}
		}
	}
	if !opts.BlockingWait {
		if len(rr.addrs) == 0 {
			rr.mu.Unlock()
			err = grpc.Errorf(codes.Unavailable, "there is no address available")
			return
		}
		// Returns the next addr on rr.addrs for failfast RPCs.
		addr = rr.addrs[rr.next].addr
		rr.next++
		rr.mu.Unlock()
		return
	}
	// Wait on rr.waitCh for non-failfast RPCs.
	if rr.waitCh == nil {
		ch = make(chan struct{})
		rr.waitCh = ch
	} else {
		ch = rr.waitCh
	}
	rr.mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-ch:
			rr.mu.Lock()
			if rr.done {
				rr.mu.Unlock()
				err = grpc.ErrClientConnClosing
				return
			}

			if len(rr.addrs) > 0 {
				if rr.next >= len(rr.addrs) {
					rr.next = 0
				}
				next := rr.next
				for {
					a := rr.addrs[next]
					next = (next + 1) % len(rr.addrs)
					if a.connected {
						addr = a.addr
						rr.next = next
						rr.mu.Unlock()
						return
					}
					if next == rr.next {
						// Has iterated all the possible address but none is connected.
						break
					}
				}
			}
			// The newly added addr got removed by Down() again.
			if rr.waitCh == nil {
				ch = make(chan struct{})
				rr.waitCh = ch
			} else {
				ch = rr.waitCh
			}
			rr.mu.Unlock()
		}
	}
}

func (rr *roundRobin) Notify() <-chan []grpc.Address {
	return rr.addrCh
}

func (rr *roundRobin) Close() error {
	rr.mu.Lock()
	defer rr.mu.Unlock()
	if rr.done {
		return ErrBalancerClosed
	}
	rr.done = true
	if rr.w != nil {
		rr.w.Close()
	}
	if rr.waitCh != nil {
		close(rr.waitCh)
		rr.waitCh = nil
	}
	if rr.addrCh != nil {
		close(rr.addrCh)
	}
	return nil
}

func randSlice(in []string, randinit bool) {
	size := len(in)
	if size < 1 {
		return
	}

	if randinit {
		rand.Seed(time.Now().Unix())
	}

	for i := 0; i < size-1; i++ {
		//addr[size-i] <-> [0, size-i)
		src := size - i - 1
		dst := rand.Intn(src + 1)

		t := in[src]
		in[src] = in[dst]
		in[dst] = t
	}
}

func DialRr(ctx context.Context, target string, rand bool) (conn *grpc.ClientConn, r *NameResolver, err error) {
	r = &NameResolver{Addrs: strings.Split(target, ",")}

	if rand {
		randSlice(r.Addrs, false)
	}

	for i := 0; ; i++ {
		timeout, _ := context.WithTimeout(ctx, time.Second)
		conn, err = grpc.DialContext(timeout, "",
			grpc.WithInsecure(),
			grpc.WithDialer(Dialer),
			grpc.WithBlock(),
			grpc.WithBalancer(RoundRobin(r)))
		if err == nil {
			break
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}

	var updates []*naming.Update
	for _, addr := range r.Addrs {
		updates = append(updates, &naming.Update{
			Op:   naming.Add,
			Addr: addr,
		})
	}
	r.W.Inject(updates)

	return
}
