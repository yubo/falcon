<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.10.0/github-markdown.min.css">
<style>
	.markdown-body {
		box-sizing: border-box;
		min-width: 200px;
		max-width: 980px;
		margin: 0 auto;
		padding: 45px;
	}

	@media (max-width: 767px) {
		.markdown-body {
			padding: 15px;
		}
	}
</style> </head>
<body> <article class="markdown-body">

<h2>
<a id="user-content-auth-api-howtomisso" class="anchor" href="#auth-api-howtomisso" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>auth api howto(misso)</h2>
<h4>
<a id="user-content-1-获得服务端sessionid" class="anchor" href="#1-%E8%8E%B7%E5%BE%97%E6%9C%8D%E5%8A%A1%E7%AB%AFsessionid" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>1 获得服务端sessionid</h4>
<p>先随便访问一个页面或者api地址</p>
<pre><code>curl -i -X GET -H 'Accept: text/plain' 'http://c3-op-mon-dev01.bj/v1.0/auth/info'
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
</code></pre>
<p>在 response header 里取出 Set-Cookie 后的字段 ( falcon 目前的 cookie 变量名为 falconSessionId， 请以 response header 里获取的变量名为准 ), 之后的所有请求将该值设置到 request header里</p>
<h4>
<a id="user-content-2-使用misso的token认证" class="anchor" href="#2-%E4%BD%BF%E7%94%A8misso%E7%9A%84token%E8%AE%A4%E8%AF%81" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>2 使用misso的token认证</h4>
<p>将misso的token加入request header, 访问 v1.0/auth/callback/misso, 成功后会返回302, 之后使用登录时的 cookie， 访问api</p>
<pre><code>curl -i -X GET -H 'Cookie: falconSessionId=4ae67ec530eaad11414485c265f2043a' -H 'Authorization: service-ocean;********************************;********************************' 'http://c3-op-mon-dev01.bj/v1.0/auth/callback/misso'
HTTP/1.1 302 Found
Server: nginx
Date: Wed, 17 May 2017 01:35:27 GMT
Content-Type: text/html; charset=utf-8
Connection: keep-alive
Location: /#
Content-Length: 25

&lt;a href="/#"&gt;Found&lt;/a&gt;.
</code></pre>
<h4>
<a id="user-content-3-访问api" class="anchor" href="#3-%E8%AE%BF%E9%97%AEapi" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>3 访问api</h4>
<p>获取用户信息</p>
<pre><code>curl -X GET -H 'Cookie: falconSessionId=4ae67ec530eaad11414485c265f2043a' 'http://c3-op-mon-dev01.bj/v1.0/auth/info'
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
</code></pre>
<p>logout</p>
<pre><code>curl -X GET -H 'Cookie: falconSessionId=4ae67ec530eaad11414485c265f2043a' 'http://c3-op-mon-dev01.bj/v1.0/auth/logout'
"logout success!"
</code></pre>
<h4>
<a id="user-content-4-其他" class="anchor" href="#4-%E5%85%B6%E4%BB%96" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>4 其他</h4>
<ul>
<li>cookie失效，重复步骤１</li>
<li>如果访问api出现401 Unauthorized, 重复步骤2</li>
<li>所有的访问，请保持cookie</li>
<li>收到的 response header 如果出现 Set-Cookie, 更新 cookie</li>
<li>misso 的 token 需要申请，使用时会限制来源ip地址，使用白名单策略</li>
</ul>
</article></body>
