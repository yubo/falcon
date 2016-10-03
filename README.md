# lite-falcon
[![wercker status](https://app.wercker.com/status/264bf495c340505f479d787192a213f4/s/master "wercker status")](https://app.wercker.com/project/byKey/264bf495c340505f479d787192a213f4)
[![License (3-Clause BSD)](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg?style=flat-square)](http://opensource.org/licenses/BSD-3-Clause)

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
- 支持快速重启(在pc机上,100万counter,重启时间1.4s左右)
- 支持单进程多模块，多实例

## 模块对应关系

| lite-falcon |   falcon                            |
|-------------|-------------------------------------|
|   agent     |   agent, _aggregator, nodata, task_ |
|   lb        |   transfer, _query, hbs_            |
|   backend   |   graph, _judge, sender_            |
|   web       |   _fe, dashboard_                   |
|   _no plan_ |   _gateway, ...  _                  |


[lite_falcon_img]:https://cdn.rawgit.com/yubo/falcon/master/doc/img/lite-falcon.svg
