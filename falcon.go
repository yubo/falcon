/*
 * Copyright 2016 yubo. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package falcon

import (
	"flag"
	"os"
	"runtime"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon/conf"
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/flags"
)

const (
	MODULE_NAME = "\x1B[32m[FALCON]\x1B[0m "
)

func init() {
	flags.NewCommand("start", "start falcon",
		start, flag.ExitOnError)

	flags.NewCommand("stop", "stop falcon",
		stop, flag.ExitOnError)

	flags.NewCommand("parse", "just parse falcon ConfigFile",
		parse, flag.ExitOnError)

	flags.NewCommand("reload", "reload falcon",
		reload, flag.ExitOnError)

}

func start(arg interface{}) {
	opts := arg.(*specs.CmdOpts)
	conf := conf.Parse(opts.ConfigFile, false)
	app := specs.NewProcess(conf.PidFile, conf.Modules)

	if err := app.Check(); err != nil {
		glog.Fatal(err)
	}
	if err := app.Save(); err != nil {
		glog.Fatal(err)
	}

	dir, _ := os.Getwd()
	glog.V(4).Infof("work dir :%s", dir)
	glog.V(4).Infof("\n%s", conf)

	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Start()
}

func stop(arg interface{}) {
	opts := arg.(*specs.CmdOpts)
	conf := conf.Parse(opts.ConfigFile, false)
	app := specs.NewProcess(conf.PidFile, conf.Modules)

	if err := app.Kill(syscall.SIGTERM); err != nil {
		glog.Fatal(err)
	}
}

func parse(arg interface{}) {
	opts := arg.(*specs.CmdOpts)
	conf := conf.Parse(opts.ConfigFile, true)
	dir, _ := os.Getwd()
	glog.Infof("work dir :%s", dir)
	glog.Infof("\n%s", conf)
}

func reload(arg interface{}) {
	opts := arg.(*specs.CmdOpts)
	conf := conf.Parse(opts.ConfigFile, false)
	app := specs.NewProcess(conf.PidFile, conf.Modules)

	if err := app.Kill(syscall.SIGUSR1); err != nil {
		glog.Fatal(err)
	}
}
