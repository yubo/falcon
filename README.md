# falcon-alpha
[![wercker status](https://app.wercker.com/status/264bf495c340505f479d787192a213f4/s/master "wercker status")](https://app.wercker.com/project/byKey/264bf495c340505f479d787192a213f4)
[![codecov](https://codecov.io/gh/yubo/falcon/branch/master/graph/badge.svg)](https://codecov.io/gh/yubo/falcon)
[![License (3-Clause BSD)](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg?style=flat-square)](http://opensource.org/licenses/BSD-3-Clause)

falcon

This is an experimental version, the official version, please refer to[open-falcon](https://github.com/open-falcon/)

## overview
![](docs/img/falcon-overview.svg?raw=true)

## build & install
```sh
go get github.com/yubo/falcon
```

[qucik start](docs/quick_start.md)

## FEATURE
- Independent management module, include host, user, role, permission, organizational structure
- Support google, github, WeChat(weixin) oauth
- Transplanted facebook [beringei](https://github.com/facebookincubator/beringei)
- Single binary & config file & process
- Use gRPC instead of json rpc
- Use gRPC gateway for RESTful API


## TODO
- Centralized configuration(etcd)

## resource

* [quick start](docs/quick_start.md)
* [user guides](docs/user_guides.md)
* [programmer guides](docs/programmer_guides.md)
* [http api](docs/swagger/index.html)
* [grpc api reference](docs/api_reference.md)
* [https://github.com/facebookincubator/beringei](https://github.com/facebookincubator/beringei)
* [https://github.com/huangaz/tsdb](https://github.com/huangaz/tsdb)
