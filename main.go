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
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon/lib/core"
	"github.com/yubo/falcon/modules/agent"
	"github.com/yubo/falcon/modules/alarm"
	"github.com/yubo/falcon/modules/ctrl"
	"github.com/yubo/falcon/modules/service"
	"github.com/yubo/falcon/modules/sys"
	"github.com/yubo/falcon/modules/transfer"
	"github.com/yubo/gotool/flags"

	_ "github.com/yubo/falcon/modules/agent/plugin"
)

const (
	PROCESS_NAME = "falcon"
)

type arrayString []string

func (i *arrayString) String() string {
	return strings.Join([]string(*i), ",")
}

func (i *arrayString) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var (
	opts struct {
		ValueFiles arrayString
		ConfigFile string
		Module     string
	}
)

func init() {

	core.RegisterModule(&agent.Agent{})
	core.RegisterModule(&alarm.Alarm{})
	core.RegisterModule(&ctrl.Ctrl{})
	core.RegisterModule(&service.Service{})
	core.RegisterModule(&sys.Sys{})
	core.RegisterModule(&transfer.Transfer{})

	flag.StringVar(&opts.ConfigFile, "config", fmt.Sprintf("./%s.yaml", PROCESS_NAME), "app config file")
	flag.Var(&opts.ValueFiles, "f", "app values file")
	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s COMMAND start|stop|reload|stats\n", os.Args[0])

	flags.NewCommand("start", "start falcon", startHandle, flag.ExitOnError)
	flags.NewCommand("stop", "stop falcon", stopHandle, flag.ExitOnError)
	flags.NewCommand("reload", "reload falcon", reloadHandle, flag.ExitOnError)

	cmd := flags.NewCommand("stats", "show falcon modules stats", statsHandle, flag.ExitOnError)
	cmd.StringVar(&opts.Module, "m", "all", "module name")

	flags.NewCommand("parse", "just parse falcon ConfigFile", parseHandle, flag.ExitOnError)
	flags.NewCommand("help", "show help information", helpHandle, flag.ExitOnError)
	flags.NewCommand("version", "show falcon version information", versionHandle, flag.ExitOnError)
	flags.NewCommand("git", "show falcon git version information", gitHandle, flag.ExitOnError)
	flags.NewCommand("changelog", "show falcon changelog information", changelogHandle, flag.ExitOnError)

}

func main() {
	flags.Parse()
	if len(os.Args) == 1 {
		startHandle(nil)
	} else {
		flags.Exec()
	}
}

func signalNotify(p *core.Process) {
	sigs := make(chan os.Signal, 1)

	glog.Infof("[%d] register signal notify", p.Pid)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	atomic.StoreUint32(&p.Status, core.APP_STATUS_RUNNING)

	for {
		s := <-sigs
		glog.Infof("recv %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			atomic.StoreUint32(&p.Status, core.APP_STATUS_EXIT)
			p.Stop()
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Infof("reload")

			atomic.StoreUint32(&p.Status, core.APP_STATUS_RELOAD)
			if err := p.Reload(); err != nil {
				glog.Errorf("not support add/del module\n")
			}
			atomic.StoreUint32(&p.Status, core.APP_STATUS_RUNNING)
		default:
			p.Signal(s)
		}
	}

}

/* handle */
func startHandle(arg interface{}) {
	app, err := core.NewProcess(opts.ConfigFile, BaseConf, []string(opts.ValueFiles))
	if err != nil {
		glog.Fatal(err)
	}

	if err := app.Start(); err != nil {
		glog.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)

	glog.Infof("[%d] register signal notify", app.Pid)
	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	atomic.StoreUint32(&app.Status, core.APP_STATUS_RUNNING)

	for {
		s := <-sigs
		glog.Infof("recv %v", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			atomic.StoreUint32(&app.Status, core.APP_STATUS_EXIT)
			app.Stop()
			os.Exit(0)
		case syscall.SIGUSR1:
			glog.Infof("reload")

			atomic.StoreUint32(&app.Status, core.APP_STATUS_RELOAD)
			if err := app.Reload(); err != nil {
				glog.Errorf("not support add/del module\n")
			}
			atomic.StoreUint32(&app.Status, core.APP_STATUS_RUNNING)
		default:
			app.Signal(s)
		}
	}
}

func stopHandle(arg interface{}) {
	app, err := core.NewProcess(opts.ConfigFile, BaseConf, []string(opts.ValueFiles))
	if err != nil {
		glog.Fatal(err)
	}
	if err := app.Kill(syscall.SIGTERM); err != nil {
		glog.Fatal(err)
	}
}

func reloadHandle(arg interface{}) {
	app, err := core.NewProcess(opts.ConfigFile, BaseConf, []string(opts.ValueFiles))
	if err != nil {
		glog.Fatal(err)
	}

	if err := app.Kill(syscall.SIGUSR1); err != nil {
		glog.Fatal(err)
	}
}

func helpHandle(arg interface{}) {
	flags.Usage()
}

func parseHandle(arg interface{}) {
	app, err := core.NewProcess(opts.ConfigFile, BaseConf, []string(opts.ValueFiles))
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Printf("%s\n", app.Configer)
}

func versionHandle(arg interface{}) {
	fmt.Printf("%s\n", VERSION)
}

func gitHandle(arg interface{}) {
	fmt.Println(COMMIT)
}

func changelogHandle(arg interface{}) {
	fmt.Println(CHANGELOG)
}

func statsHandle(arg interface{}) {
	app, err := core.NewProcess(opts.ConfigFile, BaseConf, []string(opts.ValueFiles))
	if err != nil {
		glog.Fatal(err)
	}

	err = app.Stats(opts.Module)
	if err != nil {
		glog.Fatal(err)
	}
}
