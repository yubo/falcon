package storage

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/open-falcon/graph/api"
	"github.com/open-falcon/graph/http"
	"github.com/open-falcon/graph/index"
	"github.com/open-falcon/graph/rrdtool"
	"github.com/yubo/falcon/specs"
)

const (
	GAUGE           = "GAUGE"
	DERIVE          = "DERIVE"
	COUNTER         = "COUNTER"
	CACHE_TIME      = 1800000 //ms
	FLUSH_DISK_STEP = 1000    //ms
	DEFAULT_STEP    = 60      //s
	MIN_STEP        = 30      //s
)

const (
	GRAPH_F_MISS uint32 = 1 << iota
	GRAPH_F_ERR
	GRAPH_F_SENDING
	GRAPH_F_FETCHING
)

var (
	configFile string
	configPtr  unsafe.Pointer
	configOpts Options = defaultOptions
	doneChan   []chan chan error
	mutex      sync.Mutex
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func start_signal(pid int, cfg *Options) {
	sigs := make(chan os.Signal, 1)
	log.Println(pid, "register signal notify")
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		s := <-sigs
		log.Println("recv", s)

		switch s {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			log.Println("gracefull shut down")
			if cfg.Http {
				http.Close_chan <- 1
				<-http.Close_done_chan
			}
			log.Println("http stop ok")

			if cfg.Rpc {
				api.Close_chan <- 1
				<-api.Close_done_chan
			}
			log.Println("rpc stop ok")

			rrdtool.Out_done_chan <- 1
			rrdtool.FlushAll(true)
			log.Println("rrdtool stop ok")

			log.Println(pid, "exit")
			os.Exit(0)
		}
	}
}

func Handle(arg interface{}) {
	config := arg.(*specs.CmdOptions).ConfigFile

	parse(config)

	// init db
	initDB()

	// rrdtool before api for disable loopback connection
	rrdtool.Start()

	// start api
	go api.Start()
	// start indexing
	index.Start()
	// start http server
	go http.Start()

	start_signal(os.Getpid(), Config())
}
