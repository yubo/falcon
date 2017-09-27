//
// Copyright 2017 falcon Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beego

import (
	"context"
	"fmt"
	"net"
	"net/http/fcgi"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/utils"
)

type BeegoModule struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func (p *BeegoModule) PreStart() error {
	AddAPPStartHook(
		registerMime,
		registerDefaultErrorHandler,
		registerSession,
		registerTemplate,
		registerAdmin,
		registerGzip,
	)

	for _, hk := range hooks {
		if err := hk(); err != nil {
			panic(err)
		}
	}
	return nil
}

func (p *BeegoModule) Start() error {
	p.ctx, p.cancel = context.WithCancel(context.Background())
	BeeApp.ModuleRun(p)
	return nil
}

func (p *BeegoModule) Stop() error {
	p.cancel()
	return nil
}

// Run beego module.
func (app *App) ModuleRun(p *BeegoModule) {
	addr := BConfig.Listen.HTTPAddr

	if BConfig.Listen.HTTPPort != 0 {
		addr = fmt.Sprintf("%s:%d", BConfig.Listen.HTTPAddr, BConfig.Listen.HTTPPort)
	}

	var (
		err error
		l   net.Listener
	)

	// run cgi server
	if BConfig.Listen.EnableFcgi {
		if BConfig.Listen.EnableStdIo {
			if err = fcgi.Serve(nil, app.Handlers); err == nil { // standard I/O
				logs.Info("Use FCGI via standard I/O")
			} else {
				logs.Critical("Cannot use FCGI via standard I/O", err)
			}
			return
		}
		if BConfig.Listen.HTTPPort == 0 {
			// remove the Socket file before start
			if utils.FileExists(addr) {
				os.Remove(addr)
			}
			l, err = net.Listen("unix", addr)
		} else {
			l, err = net.Listen("tcp", addr)
		}
		if err != nil {
			logs.Critical("Listen: ", err)
		}
		go func() {
			if err = fcgi.Serve(l, app.Handlers); err != nil {
				logs.Critical("fcgi.Serve: ", err)
			}
		}()
		go func() {
			<-p.ctx.Done()
			l.Close()
		}()
		return
	}

	app.Server.Handler = app.Handlers
	app.Server.ReadTimeout = time.Duration(BConfig.Listen.ServerTimeOut) * time.Second
	app.Server.WriteTimeout = time.Duration(BConfig.Listen.ServerTimeOut) * time.Second
	app.Server.ErrorLog = logs.GetLogger("HTTP")

	// run normal mode
	if BConfig.Listen.EnableHTTPS {
		go func() {
			time.Sleep(20 * time.Microsecond)
			if BConfig.Listen.HTTPSPort != 0 {
				app.Server.Addr = fmt.Sprintf("%s:%d", BConfig.Listen.HTTPSAddr, BConfig.Listen.HTTPSPort)
			} else if BConfig.Listen.EnableHTTP {
				BeeLogger.Info("Start https server error, confict with http.Please reset https port")
				return
			}
			logs.Info("https server Running on https://%s", app.Server.Addr)
			if err := app.Server.ListenAndServeTLS(BConfig.Listen.HTTPSCertFile, BConfig.Listen.HTTPSKeyFile); err != nil {
				logs.Critical("ListenAndServeTLS: ", err)
				time.Sleep(100 * time.Microsecond)
			}
			go func() {
				<-p.ctx.Done()
				app.Server.Shutdown(context.Background())
			}()
		}()
	}
	if BConfig.Listen.EnableHTTP {
		go func() {
			app.Server.Addr = addr
			logs.Info("http server Running on http://%s", app.Server.Addr)
			if BConfig.Listen.ListenTCP4 {
				ln, err := net.Listen("tcp4", app.Server.Addr)
				if err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					return
				}
				if err = app.Server.Serve(ln); err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					return
				}
			} else {
				if err := app.Server.ListenAndServe(); err != nil {
					logs.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
				}
			}
			go func() {
				<-p.ctx.Done()
				app.Server.Shutdown(context.Background())
			}()
		}()
	}
}
