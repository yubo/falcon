# lite-falcon
[![wercker status](https://app.wercker.com/status/264bf495c340505f479d787192a213f4/s/master "wercker status")](https://app.wercker.com/project/byKey/264bf495c340505f479d787192a213f4)
[![codecov](https://codecov.io/gh/yubo/falcon/branch/master/graph/badge.svg)](https://codecov.io/gh/yubo/falcon)
[![License (3-Clause BSD)](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg?style=flat-square)](http://opensource.org/licenses/BSD-3-Clause)

简化版本的falcon

出于性能测试和试验的目的，有了这个精简版本，官方的版本请参考[open-falcon](https://github.com/open-falcon/)

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
- 后端使用共享内存，在pc机上,100万counter,重启时间1.4s左右
- 支持单进程多模块，多实例
- 集中配置

## 模块对应关系

| lite-falcon |   falcon                            |  description                                        |
|-------------|-------------------------------------|-----------------------------------------------------|
|   agent     |   agent, _aggregator, nodata, task_ | 安装在需要监控的宿主机上，采集数据，发送给lb        |
|   lb        |   transfer, _query_                 | 将接收到的请求按一定策略处理或转给后端服务(backend) |
|   backend   |   graph, _judge, sender_            | 存储、处理数据的后端服务                            |
|   ctrl      |   _fe, dashboard, hbs, portal_      | 配置、统计各组件,另提供web服务                      |
|   _no plan_ |   _gateway, ..._                    | 其他组件                                            |


[lite_falcon_img]:https://cdn.rawgit.com/yubo/falcon/master/doc/img/lite-falcon.svg
