/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package backend

import (
	"net/http"
	"path"
	"strings"

	"google.golang.org/grpc"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/yubo/falcon"
	"golang.org/x/net/context"
)

// newGateway returns a new gateway server which translates HTTP into gRPC.
func newGateway(ctx context.Context, address string,
	opts ...runtime.ServeMuxOption) (http.Handler, error) {

	mux := runtime.NewServeMux(opts...)

	dialOpts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDialer(falcon.Dialer),
		grpc.WithBlock(),
	}

	err := RegisterBackendHandlerFromEndpoint(ctx, mux, address, dialOpts)
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
	p = path.Join("backend", p)
	http.ServeFile(w, r, p)
}
