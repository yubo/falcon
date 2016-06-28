package storage

type MigrateOpt struct {
	Enable      bool              `hcl:"enable"`
	Concurrency int               `hcl:"concurrency"`
	Replicas    int               `hcl:"replicas"`
	Cluster     map[string]string `hcl:"cluster"`
}

type Options struct {
	Debug       bool       `hcl:"debug"`
	Http        bool       `hcl:"http"`
	HttpAddr    string     `hcl:"http_addr"`
	Rpc         bool       `hcl:"rpc"`
	RpcAddr     string     `hcl:"rpc_addr"`
	RrdStorage  string     `hcl:"rrd_storage"`
	Dsn         string     `hcl:"dsn"`
	DbMaxIdle   int        `hcl:"db_max_idle"`
	CallTimeout int        `hcl:"call_timeout"`
	Migrate     MigrateOpt `hcl:"migrate"`
}

var (
	defaultOptions = Options{
		Debug:      false,
		Http:       true,
		HttpAddr:   "0.0.0.0:6071",
		Rpc:        true,
		RpcAddr:    "0.0.0.0:6070",
		RrdStorage: "/home/work/data/6070",
		Dsn:        "root:@tcp(127.0.0.1:3306)/graph?loc=Local&parseTime=true",
		DbMaxIdle:  4,
		Migrate: MigrateOpt{
			Enable:      false,
			Concurrency: 2,
			Replicas:    500,
			Cluster: map[string]string{
				"graph-00": "127.0.0.1:6070",
			},
		},
	}
)
