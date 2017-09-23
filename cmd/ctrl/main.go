/*
 * Copyright 2016,2017 falcon Author. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync/atomic"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/parse"
	"github.com/yubo/gotool/flags"

	_ "github.com/yubo/falcon/ctrl/ctrl"
)

var opts falcon.CmdOpts

const (
	MODULE_NAME = "\x1B[32m[MAIN]\x1B[0m "
)

func init() {
	flag.StringVar(&opts.ConfigFile, "config",
		"./etc/falcon.conf", "falcon config file")

	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s [OPTIONS] COMMAND "+
		"start|stop|reload\n", os.Args[0])
	flags.NewCommand("help", "show help information", helpHandle, flag.ExitOnError)
	flags.NewCommand("start", "start falcon", start, flag.ExitOnError)
	flags.NewCommand("stop", "stop falcon", stop, flag.ExitOnError)
	flags.NewCommand("parse", "just parse falcon ConfigFile", parseHandle,
		flag.ExitOnError)
	flags.NewCommand("reload", "reload falcon", reload, flag.ExitOnError)

	flags.NewCommand("version", "show falcon version information",
		version, flag.ExitOnError)

	flags.NewCommand("git", "show falcon git version information",
		git, flag.ExitOnError)

	flags.NewCommand("changelog", "show falcon changelog information",
		changelog, flag.ExitOnError)
}

func helpHandle(arg interface{}) {
	flags.Usage()
}

func signalNotify(p *falcon.Process) {
	sigs := make(chan os.Signal, 1)

	glog.Infof(MODULE_NAME+"[%d] register signal notify", p.Pid)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
		syscall.SIGUSR1)
	atomic.StoreUint32(&p.Status, falcon.APP_STATUS_RUNNING)

	glog.Infof(MODULE_NAME+"[%d] register signal notify", p.Pid)

	for {
		s := <-sigs
		glog.Infof(MODULE_NAME+"recv %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			pidfile := fmt.Sprintf("%s.%d", p.Config.PidFile, p.Pid)
			glog.Info(MODULE_NAME + "exiting")
			atomic.StoreUint32(&p.Status, falcon.APP_STATUS_EXIT)
			os.Rename(p.Config.PidFile, pidfile)

			for i, n := 0, len(p.Module); i < n; i++ {
				p.Module[n-i-1].Stop()
			}

			glog.Infof(MODULE_NAME+"pid:%d exit", p.Pid)
			os.Remove(pidfile)
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Info(MODULE_NAME + "reload")

			// reparse config, get new config
			// newConfig := parse.Parse(p.Config.ConfigFile, false)
			newConfig := parse.Parse(p.Config.ConfigFile, false)

			// check config diff
			if len(newConfig.Conf) != len(p.Config.Conf) {
				glog.Error("not support add/del module\n")
				break
			}

			for i, config := range newConfig.Conf {
				m, ok := falcon.ModuleTpls[falcon.GetType(config)]
				if !ok {
					glog.Exitf("%s's module not support, you should"+
						" import module ", falcon.GetType(config))
					break
				}
				newM := m.New(config)
				if newM.Name() != p.Module[i].Name() {
					glog.Exitf("%s's module not support,"+
						" not support add/del/disable module",
						falcon.GetType(config))
					break
				}
			}

			// do it
			atomic.StoreUint32(&p.Status, falcon.APP_STATUS_RELOAD)
			falcon.SetGlog(newConfig)
			for i, m := range p.Module {
				m.Reload(newConfig.Conf[i])
			}
			atomic.StoreUint32(&p.Status, falcon.APP_STATUS_RUNNING)
		default:
			for _, m := range p.Module {
				m.Signal(s)
			}
		}
	}

}

func start(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := parse.Parse(opts.ConfigFile, false)
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
	signalNotify(app)
}

func stop(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := parse.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGTERM); err != nil {
		glog.Fatal(err)
	}
}

func reload(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := parse.Parse(opts.ConfigFile, false)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGUSR1); err != nil {
		glog.Fatal(err)
	}
}

func parseHandle(arg interface{}) {
	opts := arg.(*falcon.CmdOpts)
	c := parse.Parse(opts.ConfigFile, true)
	dir, _ := os.Getwd()
	glog.Infof("work dir :%s", dir)
	glog.Infof("\n%s", c)
}

func version(arg interface{}) {
	fmt.Printf("%s\n", falcon.VERSION)
}

func git(arg interface{}) {
	fmt.Println(falcon.COMMIT)
}

func changelog(arg interface{}) {
	fmt.Println(falcon.CHANGELOG)
}

func main() {
	flags.Parse()
	cmd := flags.CommandLine.Cmd

	if cmd != nil && cmd.Action != nil {
		opts.Args = cmd.Flag.Args()
		cmd.Action(&opts)
	} else {
		opts.Args = flag.Args()
		start(&opts)
	}
}
