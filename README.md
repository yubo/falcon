# lite-falcon

简化版本的falcon

出于性能测试的需要，有了这个精简版本，官方的版本请参考[open-falcon](https://github.com/open-falcon/)

![][lite_falcon_img]

## 使用
```
make
```

## benchmark
```
cd backend
go test -bench=Add -benchtime=20s
go test -bench=.
```


## 特点
- 整合／精简了代码,依赖包
- 简化配置
- 多硬盘支持(单机多写)

## 模块对应关系

| lite-falcon |   falcon                            |
|-------------|-------------------------------------|
|   agent     |   agent, _aggregator, nodata, task_ |
|   handoff   |   transfer, _query_                 |
|   backend   |   graph, _judge, sender_            |
|   _no plan_ |   _gateway,  hbs_                   |


[lite_falcon_img]:https://cdn.rawgit.com/yubo/falcon/master/doc/img/lite-falcon.svg
