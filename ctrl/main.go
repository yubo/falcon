/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syscall"

	"github.com/astaxie/beego"
	"github.com/golang/glog"
	"github.com/yubo/falcon/ctrl/core"
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/flags"
)

var (
	opts      specs.CmdOpts
	app       *specs.Process
	ct        ctrl.Ctrl
	pidfile   string
	ifpre     string
	upstreams string
)

func init() {
	host, _ := os.Hostname()

	flag.StringVar(&pidfile, "p", "/tmp/ctrl.pid", "pid file path")
	flag.IntVar(&ct.Params.Debug, "d", 0, "debug level")
	flag.StringVar(&ct.Params.Host, "host", host, "hostname")
	flag.BoolVar(&ct.Params.Rpc, "rpc", true, "enable rpc")
	flag.BoolVar(&ct.Params.Http, "http", false, "enable http")
	flag.StringVar(&ct.Params.RpcAddr, "ra", "127.0.0.1:1988", "rpc addr")
	flag.StringVar(&ct.Params.HttpAddr, "ha", "127.0.0.1:1989", "http addr")
	flag.StringVar(&opts.ConfigFile, "config",
		"./conf/app.conf", "ctrl config file")

	beego.BConfig.AppName = opts.ConfigFile

	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s [OPTIONS] COMMAND ",
		"start|stop\n", os.Args[0])

	flags.NewCommand("start", "start agent",
		start, flag.ExitOnError)

	flags.NewCommand("stop", "stop agent",
		stop, flag.ExitOnError)

	flags.NewCommand("help", "show help information",
		help, flag.ExitOnError)
}

func start(arg interface{}) {
	app := specs.NewProcess(pidfile, []specs.Module{specs.Module(&ct)})

	if err := app.Check(); err != nil {
		glog.Fatal(err)
	}
	if err := app.Save(); err != nil {
		glog.Fatal(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Start()
}

func stop(arg interface{}) {
	app := specs.NewProcess(pidfile, []specs.Module{specs.Module(&ct)})
	if err := app.Kill(syscall.SIGTERM); err != nil {
		glog.Fatal(err)
	}
}

func help(arg interface{}) {
	flags.Usage()
}

func main() {
	flags.Parse()
	cmd := flags.CommandLine.Cmd

	if cmd != nil && cmd.Action != nil {
		opts.Args = cmd.Flag.Args()
		cmd.Action(&opts)
	} else {
		start(&opts)
	}
}
