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

	_ "github.com/yubo/falcon/agent/modules"
	_ "github.com/yubo/falcon/alarm/modules"
	_ "github.com/yubo/falcon/ctrl/modules"
	_ "github.com/yubo/falcon/service/modules"
	_ "github.com/yubo/falcon/transfer/modules"

	"net/http"
	_ "net/http/pprof"
)

var opts falcon.CmdOpts

const (
	MODULE_NAME = "\x1B[34m[MAIN]\x1B[0m"
)

func init() {
	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s COMMAND start|stop|reload|stats\n", os.Args[0])

	cmd := flags.NewCommand("start", "start falcon", start, flag.ExitOnError)
	cmd.StringVar(&opts.ConfigFile, "config", "/etc/falcon/falcon.conf", "falcon config file")

	cmd = flags.NewCommand("parse", "just parse falcon ConfigFile", parseHandle, flag.ExitOnError)
	cmd.StringVar(&opts.ConfigFile, "config", "/etc/falcon/falcon.conf", "falcon config file")

	cmd = flags.NewCommand("reload", "reload falcon", reload, flag.ExitOnError)
	cmd.StringVar(&opts.ConfigFile, "config", "/etc/falcon/falcon.conf", "falcon config file")

	cmd = flags.NewCommand("stats", "show falcon modules stats", stats, flag.ExitOnError)
	cmd.StringVar(&opts.ConfigFile, "config", "/etc/falcon/falcon.conf", "falcon config file")
	cmd.StringVar(&opts.Module, "m", "all", "module name")

	flags.NewCommand("help", "show help information", helpHandle, flag.ExitOnError)
	flags.NewCommand("stop", "stop falcon", stop, flag.ExitOnError)
	flags.NewCommand("version", "show falcon version information", version, flag.ExitOnError)
	flags.NewCommand("git", "show falcon git version information", git, flag.ExitOnError)
	flags.NewCommand("changelog", "show falcon changelog information", changelog, flag.ExitOnError)
	flags.NewCommand("modules", "show falcon modules information", modules, flag.ExitOnError)

}

func helpHandle(arg interface{}) {
	flags.Usage()
}

func signalNotify(p *falcon.Process) {
	sigs := make(chan os.Signal, 1)

	glog.Infof("%s [%d] register signal notify", MODULE_NAME, p.Pid)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	atomic.StoreUint32(&p.Status, falcon.APP_STATUS_RUNNING)

	glog.Infof("%s [%d] register signal notify", MODULE_NAME, p.Pid)
	go http.ListenAndServe(":8008", nil)

	for {
		s := <-sigs
		glog.Infof("%s recv %v", MODULE_NAME, s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			pidfile := fmt.Sprintf("%s.%d", p.Config.PidFile, p.Pid)
			glog.Infof("%s exiting", MODULE_NAME)
			atomic.StoreUint32(&p.Status, falcon.APP_STATUS_EXIT)
			os.Rename(p.Config.PidFile, pidfile)

			for i, n := 0, len(p.Module); i < n; i++ {
				p.Module[n-i-1].Stop()
			}

			glog.Infof("%s pid:%d exit", MODULE_NAME, p.Pid)
			os.Remove(pidfile)
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Infof("%s reload", MODULE_NAME)

			// reparse config, get new config
			// newConfig := parse.Parse(p.Config.ConfigFile, false)
			newConfig := parse.Parse(p.Config.ConfigFile)

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
	c := parse.Parse(opts.ConfigFile)
	app := falcon.NewProcess(c)

	if err := app.Check(); err != nil {
		glog.Fatal(err)
	}
	if err := app.Save(); err != nil {
		glog.Fatal(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	app.Start()
	signalNotify(app)
}

func stop(arg interface{}) {
	c := parse.Parse(opts.ConfigFile)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGTERM); err != nil {
		glog.Fatal(err)
	}
}

func reload(arg interface{}) {
	c := parse.Parse(opts.ConfigFile)
	app := falcon.NewProcess(c)

	if err := app.Kill(syscall.SIGUSR1); err != nil {
		glog.Fatal(err)
	}
}

func parseHandle(arg interface{}) {
	c := parse.Parse(opts.ConfigFile)
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

func modules(arg interface{}) {
	for m, _ := range falcon.Modules {
		fmt.Printf("%s\n", m)
	}
}

func stats(arg interface{}) {
	c := parse.Parse(opts.ConfigFile)
	falcon.NewProcess(c).Stats(opts.Module)
}

func main() {
	flags.Parse()
	flags.Exec()
}
