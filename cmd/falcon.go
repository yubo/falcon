package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/open-falcon/graph/api"
	"github.com/open-falcon/graph/http"
	"github.com/open-falcon/graph/index"
	"github.com/open-falcon/graph/rrdtool"
	"github.com/yubo/falcon"
	"github.com/yubo/falcon/storage"
	"github.com/yubo/gotool/flags"
)

func init() {
	flags.CommandLine.Usage = fmt.Sprintf(
		"Usage: %s COMMAND [OPTIONS] host[:port]\n\n",
		os.Args[0])

	flags.NewCommand("version", "show falcon version information",
		version_handle, flag.ExitOnError)

	flags.NewCommand("git", "show falcon git version information",
		git_handle, flag.ExitOnError)

	flags.NewCommand("changelog", "show falcon changelog information",
		changelog_handle, flag.ExitOnError)
}

func main() {
	flags.Parse()

	cmd := flags.CommandLine.Cmd
	if cmd != nil && cmd.Action != nil {
		if err != nil {
			fmt.Println("cannot connection to dpvs server")
			return
		}

		cmd.Action(&dpvs.CallOptions{Opt: CmdOpt,
			Args: cmd.Flag.Args()})
	} else {
		flags.Usage()
	}

}
