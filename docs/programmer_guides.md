## profile

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

see also
```
https://blog.golang.org/profiling-go-programs
```
