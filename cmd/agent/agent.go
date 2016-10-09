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
	"github.com/yubo/falcon/agent"
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/flags"
)

var (
	opts      specs.CmdOpts
	app       *specs.Process
	ag        agent.Agent
	pidfile   string
	ifpre     string
	upstreams string
)

func init() {
	host, _ := os.Hostname()

	flag.StringVar(&pidfile, "p", "/tmp/agnet.pid", "pid file path")
	flag.IntVar(&ag.Params.Debug, "d", 0, "debug level")
	flag.StringVar(&ag.Params.Host, "host", host, "hostname")
	flag.BoolVar(&ag.Params.Rpc, "rpc", true, "enable rpc")
	flag.BoolVar(&ag.Params.Http, "http", true, "enable http")
	flag.StringVar(&ag.Params.RpcAddr, "ra", "127.0.0.1:1988", "rpc addr")
	flag.StringVar(&ag.Params.HttpAddr, "ha", "127.0.0.1:1989", "http addr")
	flag.StringVar(&ag.Params.CtrlAddr, "ca", "127.0.0.1:8001", "ctrl addr")
	flag.StringVar(&ifpre, "if", "eth,em", "interface prefix")
	flag.IntVar(&ag.Interval, "interval", 60, "interval for collecting data(s)")
	flag.IntVar(&ag.Batch, "batch", 60, "batch number per send")
	flag.IntVar(&ag.Params.ConnTimeout, "conntimeout", 1000, "conntimeout(ms)")
	flag.IntVar(&ag.Params.CallTimeout, "calltimeout", 5000, "calltimeout(ms)")

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
		glog.Fatal(specs.ErrParam)
	}

	ag.IfPre = strings.Split(ifpre, ",")

	app := specs.NewProcess(pidfile, []specs.Module{specs.Module(&ag)})

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
	app := specs.NewProcess(pidfile, []specs.Module{specs.Module(&ag)})
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
