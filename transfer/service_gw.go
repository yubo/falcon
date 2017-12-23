/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package transfer

import (
	"net/http"
	"path"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
)

type ServiceGwModule struct {
	enable   bool
	ctx      context.Context
	cancel   context.CancelFunc
	upstream string
	address  string
}

func (p *ServiceGwModule) prestart(transfer *Transfer) error {
	p.upstream = transfer.Conf.Configer.Str(C_GRPC_ADDR)
	p.address = transfer.Conf.Configer.Str(C_HTTP_ADDR)
	p.enable, _ = transfer.Conf.Configer.Bool(C_HTTP_ENABLE)
	return nil
}

func (p *ServiceGwModule) start(transfer *Transfer) error {
	if !p.enable {
		glog.Info(MODULE_NAME + "http.Start not enabled")
		return nil
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	mux := http.NewServeMux()

	err := Gateway(p.ctx, mux, p.upstream)
	if err != nil {
		return nil
	}

	server := &http.Server{Addr: p.address, Handler: mux}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			p.cancel()
		}
		return
	}()

	go func() {
		<-p.ctx.Done()
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		server.Shutdown(ctx)
	}()
	return nil
}

func (p *ServiceGwModule) stop(transfer *Transfer) error {
	if !p.enable {
		return nil
	}
	p.cancel()
	return nil
}

func (p *ServiceGwModule) reload(transfer *Transfer) error {
	if p.enable {
		p.stop(transfer)
		time.Sleep(time.Second)
	}
	p.prestart(transfer)
	return p.start(transfer)
}

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(ctx context.Context, address string,
	opts ...runtime.ServeMuxOption) (http.Handler, error) {

	mux := runtime.NewServeMux(opts...)

	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDialer(falcon.Dialer),
		grpc.WithBlock(),
	}

	err := RegisterTransferHandlerFromEndpoint(ctx, mux, address, dialOpts)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func Gateway(ctx context.Context, mux *http.ServeMux, upstream string,
	opts ...runtime.ServeMuxOption) error {

	mux.HandleFunc("/swagger/", serveSwagger)

	gw, err := newGateway(ctx, upstream, opts...)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	return nil
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		glog.Errorf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	glog.Infof("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join("transfer", p)
	http.ServeFile(w, r, p)
}
