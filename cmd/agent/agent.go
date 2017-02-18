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
	"strings"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/agent"
	"github.com/yubo/gotool/flags"
)

var (
	opts      falcon.CmdOpts
	app       *falcon.Process
	ag        agent.Agent
	pidfile   string
	ifpre     string
	upstreams string
)

func init() {
	host, _ := os.Hostname()
	def := falcon.ConfAgentDef

	flag.StringVar(&pidfile, "p", "/tmp/agnet.pid", "pid file path")
	flag.IntVar(&ag.Conf.Params.Debug, "d", def.Params.Debug, "debug level")
	flag.StringVar(&ag.Conf.Params.Host, "host", host, "hostname")
	flag.BoolVar(&ag.Conf.Params.Rpc, "rpc", def.Params.Rpc, "enable rpc")
	flag.BoolVar(&ag.Conf.Params.Http, "http", def.Params.Http, "enable http")
	flag.StringVar(&ag.Conf.Params.RpcAddr, "ra", def.Params.RpcAddr, "rpc addr")
	flag.StringVar(&ag.Conf.Params.HttpAddr, "ha", def.Params.HttpAddr, "http addr")
	flag.StringVar(&ag.Conf.Params.CtrlAddr, "ca", def.Params.CtrlAddr, "ctrl addr")
	flag.StringVar(&ifpre, "if", strings.Join(def.IfPre, ","), "interface prefix")
	flag.StringVar(&upstreams, "upstreams", strings.Join(def.Upstreams, ","), "interface prefix")
	flag.IntVar(&ag.Conf.Interval, "interval", def.Interval, "interval for collecting data(s)")
	flag.IntVar(&ag.Conf.PayloadSize, "payloadSize", def.PayloadSize, "meta number per rpc call")
	flag.IntVar(&ag.Conf.Params.ConnTimeout, "conntimeout", def.Params.ConnTimeout, "conntimeout(ms)")
	flag.IntVar(&ag.Conf.Params.CallTimeout, "calltimeout", def.Params.CallTimeout, "calltimeout(ms)")

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
	if ifpre == "" || upstreams == "" {
		glog.Fatal(falcon.ErrParam)
	}

	ag.Conf.IfPre = strings.Split(ifpre, ",")
	ag.Conf.Upstreams = strings.Split(upstreams, ",")

	app := falcon.NewProcess(pidfile, []falcon.Module{falcon.Module(&ag)})

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
	app := falcon.NewProcess(pidfile, []falcon.Module{falcon.Module(&ag)})
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
		flags.Usage()
	}
}
