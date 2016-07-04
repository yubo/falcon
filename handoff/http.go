package handoff

import (
	"net"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/yubo/falcon/specs"
)

var (
	httpEvent  *specs.RoutineEvent
	httpConfig HandoffOpts
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func httpRoutes() {
}

func httpStart(config HandoffOpts) {
	if !config.Http {
		glog.Info("http.Start warning, not enabled")
		return
	}

	httpConfig = config

	addr := httpConfig.HttpAddr
	if addr == "" {
		return
	}
	s := &http.Server{
		Addr:           addr,
		MaxHeaderBytes: 1 << 30,
	}
	glog.Infof("http listening %s", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatal(err)
	}

	l := ln.(*net.TCPListener)
	registerEventChan(httpEvent)

	go s.Serve(tcpKeepAliveListener{l})

	go func() {
		select {
		case event := <-httpEvent.E:
			if event.Method == specs.ROUTINE_EVENT_M_EXIT {
				l.Close()
				event.Done <- nil
				return
			}
		}
	}()

}
