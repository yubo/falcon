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
	"runtime/pprof"
	"sync/atomic"
	"syscall"

	"github.com/golang/glog"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/parse"
	"github.com/yubo/gotool/flags"

	"net/http"
	_ "net/http/pprof"

	_ "github.com/yubo/falcon/agent/modules"
	_ "github.com/yubo/falcon/alarm/modules"
	_ "github.com/yubo/falcon/ctrl/modules"
	_ "github.com/yubo/falcon/service/modules"
	_ "github.com/yubo/falcon/transfer/modules"
)

var (
	opts      falcon.CmdOpts
	cpu_prof  *os.File
	heap_prof *os.File
)

const (
	MODULE_NAME    = "\x1B[34m[MAIN]\x1B[0m"
	CPU_PROF_FILE  = "/tmp/cpu.prof"
	HEAP_PROF_FILE = "/tmp/heap.prof"
)

func init() {
	flag.StringVar(&opts.ConfigFile, "config", "/etc/falcon/falcon.conf", "falcon config file")

	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s COMMAND start|stop|reload|stats\n", os.Args[0])

	cmd := flags.NewCommand("start", "start falcon", start, flag.ExitOnError)
	cmd.BoolVar(&opts.CpuProfile, "cpu", false, "cpu profile(/tmp/cpu.prof)")
	cmd.BoolVar(&opts.HeapProfile, "heap", false, "heap profile(/tmp/heap.prof)")

	cmd = flags.NewCommand("stats", "show falcon modules stats", stats, flag.ExitOnError)
	cmd.StringVar(&opts.Module, "m", "all", "module name")

	flags.NewCommand("parse", "just parse falcon ConfigFile", parseHandle, flag.ExitOnError)
	flags.NewCommand("reload", "reload falcon", reload, flag.ExitOnError)
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

	go http.ListenAndServe(":18008", nil)

	for {
		s := <-sigs
		glog.Infof("%s recv %v", MODULE_NAME, s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			pidfile := fmt.Sprintf("%s.%d", p.Config.PidFile, p.Pid)
			glog.Infof("%s exiting", MODULE_NAME)
			atomic.StoreUint32(&p.Status, falcon.APP_STATUS_EXIT)
			os.Rename(p.Config.PidFile, pidfile)
			if opts.CpuProfile {
				pprof.StopCPUProfile()
				cpu_prof.Close()
			}
			if opts.HeapProfile {
				heap_prof.Close()
			}

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
	var err error

	c := parse.Parse(opts.ConfigFile)
	app := falcon.NewProcess(c)

	if err = app.Check(); err != nil {
		glog.Fatal(err)
	}
	if err = app.Save(); err != nil {
		glog.Fatal(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	if opts.CpuProfile {
		if cpu_prof, err = os.OpenFile(CPU_PROF_FILE, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			glog.Fatal(err)
		}
		pprof.StartCPUProfile(cpu_prof)
	}

	if opts.HeapProfile {
		if heap_prof, err = os.OpenFile(HEAP_PROF_FILE, os.O_RDWR|os.O_CREATE, 0644); err != nil {
			glog.Fatal(err)
		}
		pprof.WriteHeapProfile(heap_prof)
	}

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
	err := falcon.NewProcess(c).Stats(opts.Module)
	if err != nil {
		glog.Fatal(err)
	}
}

func main() {
	flags.Parse()
	flags.Exec()
}
