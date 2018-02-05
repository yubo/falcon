## write a agent plugin

1. first, define a struct for agent.Collector interface, and rigister it

github.com/example/plugin/demo.cos.go

```golang
func init() {
	agent.RegisterCollector(&cosCollector{})
}

type cosCollector struct{}

func (p *cosCollector) Name() (name, gname string) {
	return "cos", "demo"
}

func (p *cosCollector) Start(ctx context.Context, a *agent.Agent) error {
	interval, _ := a.Conf.Configer.Int(agent.C_INTERVAL)
	ticker := time.NewTicker(time.Second * time.Duration(interval)).C
	ch := a.PutChan

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				ch <- &agent.PutRequest{
					Dps: []*transfer.DataPoint{
						agent.GaugeValue("cos.ticker.demo",
							math.Cos(x*(float64(time.Now().Unix())+1800))),
					},
				}
			}
		}
	}()
	return nil
}

func (p *cosCollector) Collect() (ret []*transfer.DataPoint, err error) {
	return []*transfer.DataPoint{
		agent.GaugeValue("cos.collect.demo",
			math.Cos(x*float64(time.Now().Unix()))),
	}, nil

}
```

2.  import in your code


```golang
import _ "github.com/example/plugin/demo.cos.go"

```


3. add it's group name in the configuration file agent/plugins

falcon/etc/falcon.conf

```
agent {
	...
	plugins		"...,demo";
}

```

## write 



## 

## DEBUG

#### profile

```
falcon start -cpu/heap
^C

$go tool pprof  /tmp/cpu.prof 
Entering interactive mode (type "help" for commands)
(pprof) top
46220ms of 46350ms total (99.72%)
Dropped 94 nodes (cum <= 231.75ms)
      flat  flat%   sum%        cum   cum%
   14880ms 32.10% 32.10%    14880ms 32.10%  runtime.chanrecv
   12690ms 27.38% 59.48%    27570ms 59.48%  runtime.selectnbrecv
   10080ms 21.75% 81.23%    10080ms 21.75%  github.com/yubo/falcon/lib/tsdb.(*KeyListWriter).writeOneKey
    7860ms 16.96% 98.19%    46220ms 99.72%  github.com/yubo/falcon/lib/tsdb.(*KeyListWriter).startWriterThread.func1
     710ms  1.53% 99.72%      710ms  1.53%  github.com/yubo/falcon/lib/tsdb.(*KeyListWriter).startWriterThread
         0     0% 99.72%    46330ms   100%  runtime.goexit
```

see also [profiling-go-programs](https://blog.golang.org/profiling-go-programs)
