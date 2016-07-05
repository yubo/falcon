# lite-falcon

简化版本的falcon

出于性能测试的需要，有了这个精简版本，官方的版本请参考[open-falcon](https://github.com/open-falcon/)

## 使用
```
make
```

## 特点
- 整合／精简了代码,依赖包
- 简化配置

## 模块对应关系

| lite-falcon |   falcon   |
|-------------|------------|
|   storage   |   graph    |
|   handoff   |   transfer |
|   agent     |   agent    |

## TODO
- test/benchmark
- 多硬盘支持(单机多写)


