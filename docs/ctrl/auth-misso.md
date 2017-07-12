## auth api howto(misso)


#### 1 获得服务端sessionid

先随便访问一个页面或者api地址

```
curl -i -X GET -H 'Accept: text/plain' 'http://c3-op-mon-dev01.bj/v1.0/auth/info'
HTTP/1.1 200 OK
Server: nginx
Date: Wed, 17 May 2017 01:12:59 GMT
Content-Type: text/plain; charset=utf-8
Connection: keep-alive
Set-Cookie: falconSessionId=4ae67ec530eaad11414485c265f2043a; Path=/; Expires=Thu, 18 May 2017 01:12:59 GMT; Max-Age=86400; HttpOnly
Content-Length: 76

{
  "user": null,
  "reader": false,
  "operator": false,
  "admin": false
}
```

在 response header 里取出 Set-Cookie 后的字段 ( falcon 目前的 cookie 变量名为 falconSessionId， 请以 response header 里获取的变量名为准 ), 之后的所有请求将该值设置到 request header里


#### 2 使用misso的token认证

将misso的token加入request header, 访问 v1.0/auth/callback/misso, 成功后会返回302, 之后使用登录时的 cookie， 访问api

```
curl -i -X GET -H 'Cookie: falconSessionId=4ae67ec530eaad11414485c265f2043a' -H 'Authorization: service-ocean;********************************;********************************' 'http://c3-op-mon-dev01.bj/v1.0/auth/callback/misso'
HTTP/1.1 302 Found
Server: nginx
Date: Wed, 17 May 2017 01:35:27 GMT
Content-Type: text/html; charset=utf-8
Connection: keep-alive
Location: /#
Content-Length: 25

<a href="/#">Found</a>.
```

#### 3 访问api

获取用户信息

```
curl -X GET -H 'Cookie: falconSessionId=4ae67ec530eaad11414485c265f2043a' 'http://c3-op-mon-dev01.bj/v1.0/auth/info'
{
  "user": {
    "id": 2982,
    "uuid": "service-ocean@misso",
    "name": "service-ocean@misso",
    "cname": "",
    "email": "",
    "phone": "",
    "im": "",
    "qq": "",
    "disabled": 0,
    "ctime": "2017-05-04T11:33:35+08:00"
  },
  "reader": true,
  "operator": true,
  "admin": true
}
```

logout

```
curl -X GET -H 'Cookie: falconSessionId=4ae67ec530eaad11414485c265f2043a' 'http://c3-op-mon-dev01.bj/v1.0/auth/logout'
"logout success!"
```

#### 4 其他

- cookie失效，重复步骤１
- 如果访问api出现401 Unauthorized, 重复步骤2
- 所有的访问，请保持cookie
- 收到的 response header 如果出现 Set-Cookie, 更新 cookie
- misso 的 token 需要申请，使用时会限制来源ip地址，使用白名单策略

