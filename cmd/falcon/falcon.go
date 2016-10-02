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

	_ "github.com/yubo/falcon"
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/flags"
)

var opts specs.CmdOpts

func init() {
	flag.StringVar(&opts.ConfigFile, "config",
		"/etc/falcon/falcon.conf", "falcon config file")

	flags.CommandLine.Usage = fmt.Sprintf("Usage: %s [OPTIONS] COMMAND ",
		"start|stop|reload\n", os.Args[0])

	flags.NewCommand("help", "show help information",
		help, flag.ExitOnError)
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
