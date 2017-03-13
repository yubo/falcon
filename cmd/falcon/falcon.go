/*
 * Copyright 2017 yubo. All rights reserved.
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

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	_ "github.com/yubo/falcon/agent"
	_ "github.com/yubo/falcon/agent/plugin"
	_ "github.com/yubo/falcon/agent/plugin/emulator"
	_ "github.com/yubo/falcon/backend"
	_ "github.com/yubo/falcon/ctrl"
	_ "github.com/yubo/falcon/ctrl/api/models/auth"
	_ "github.com/yubo/falcon/ctrl/api/routers"
	_ "github.com/yubo/falcon/loadbalance"
	"github.com/yubo/gotool/flags"
)

var opts falcon.CmdOpts

func init() {
	flag.StringVar(&opts.ConfigFile, "config",
		"/etc/falcon/falcon.conf", "falcon config file")

	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s [OPTIONS] COMMAND ",
		"start|stop|reload\n", os.Args[0])
	flags.NewCommand("help", "show help information", help, flag.ExitOnError)
	flags.NewCommand("start", "start falcon", start, flag.ExitOnError)
	flags.NewCommand("stop", "stop falcon", stop, flag.ExitOnError)
	flags.NewCommand("parse", "just parse falcon ConfigFile", parse, flag.ExitOnError)
	flags.NewCommand("reload", "reload falcon", reload, flag.ExitOnError)
}

func help(arg interface{}) {
	flags.Usage()
}

func start(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := falcon.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Check(); err != nil {
		glog.Fatal(err)
	}
	if err := app.Save(); err != nil {
		glog.Fatal(err)
	}

	dir, _ := os.Getwd()
	glog.V(4).Infof("work dir :%s", dir)
	glog.V(4).Infof("\n%s", c)

	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Start()
}

func stop(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := falcon.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGTERM); err != nil {
		glog.Fatal(err)
	}
}

func reload(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := falcon.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGUSR1); err != nil {
		glog.Fatal(err)
	}
}

func parse(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := falcon.Parse(opts.ConfigFile, true)
	dir, _ := os.Getwd()
	glog.Infof("work dir :%s", dir)
	glog.Infof("\n%s", c)
}

func main() {
	flags.Parse()
	cmd := flags.CommandLine.Cmd

	if cmd != nil && cmd.Action != nil {
		opts.Args = cmd.Flag.Args()
		cmd.Action(&opts)
	} else {
		//flags.Usage()
		opts.Args = flag.Args()
		start(&opts)
	}
}
