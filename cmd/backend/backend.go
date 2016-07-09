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

	"github.com/yubo/falcon"
	"github.com/yubo/falcon/backend"
	"github.com/yubo/falcon/specs"
	"github.com/yubo/gotool/flags"
)

var opts specs.CmdOpts

func init() {
	flags.CommandLine.Usage = fmt.Sprintf(
		"Usage: %s COMMAND [OPTIONS] start|stop|status|reload\n",
		os.Args[0])

	flag.StringVar(&opts.ConfigFile, "config",
		"/etc/falcon/backend.conf", "backend config file")

	flags.NewCommand("version", "show falcon version information",
		falcon.Version_handle, flag.ExitOnError)

	flags.NewCommand("git", "show falcon git version information",
		falcon.Git_handle, flag.ExitOnError)

	flags.NewCommand("changelog", "show falcon changelog information",
		falcon.Changelog_handle, flag.ExitOnError)

	flags.NewCommand("help", "show help information",
		falcon.Help_handle, flag.ExitOnError)
}

func main() {
	flags.Parse()

	cmd := flags.CommandLine.Cmd
	if cmd != nil && cmd.Action != nil {
		opts.Args = cmd.Flag.Args()
		cmd.Action(&opts)
	} else {
		opts.Args = flag.Args()
		backend.Handle(&opts)
	}

}
