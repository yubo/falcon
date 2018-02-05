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
<a id="user-content-falcon-ctrl-auth" class="anchor" href="#falcon-ctrl-auth" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>falcon ctrl auth</h2>
<h4>
<a id="user-content-google" class="anchor" href="#google" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>google</h4>
<ul>
<li>
<p><a href="https://developers.google.com/identity/protocols/OpenIDConnect#getcredentials" rel="nofollow">https://developers.google.com/identity/protocols/OpenIDConnect#getcredentials</a></p>
</li>
<li>
<p><a href="https://console.developers.google.com/apis/credentials/oauthclient" rel="nofollow">https://console.developers.google.com/apis/credentials/oauthclient</a></p>
</li>
<li>
<p>网页应用的客户端id</p>
</li>
<li>
<p>已获授权的 javascript 来源</p>
<ul>
<li><a href="http://falcon.pt.xiaomi.com" rel="nofollow">http://falcon.pt.xiaomi.com</a></li>
<li><a href="https://falcon.pt.xiaomi.com" rel="nofollow">https://falcon.pt.xiaomi.com</a></li>
</ul>
</li>
<li>
<p>已获授权的重定向 URI</p>
<ul>
<li><a href="http://falcon.pt.xiaomi.com/v1.0/auth/callback/google" rel="nofollow">http://falcon.pt.xiaomi.com/v1.0/auth/callback/google</a></li>
<li><a href="https://falcon.pt.xiaomi.com/v1.0/auth/callback/google" rel="nofollow">https://falcon.pt.xiaomi.com/v1.0/auth/callback/google</a></li>
</ul>
</li>
</ul>
<p>将得到的客户端id 和 客户端密钥 分别填写到
falcon.conf -&gt; ctrl -&gt; GoogleClientID , GoogleClientSecret</p>
<h4>
<a id="user-content-github" class="anchor" href="#github" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>github</h4>
<ul>
<li><a href="https://help.github.com/enterprise/2.11/admin/guides/user-management/using-github-oauth/">https://help.github.com/enterprise/2.11/admin/guides/user-management/using-github-oauth/</a></li>
<li><a href="https://github.com/settings/applications/">https://github.com/settings/applications/</a></li>
</ul>
<p>将得到的客户端id 和 客户端密钥 分别填写到
falcon.conf -&gt; ctrl -&gt; GithubClientID , GithubClientSecret</p>
<h4>
<a id="user-content-misso" class="anchor" href="#misso" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>misso</h4>
<ul>
<li>使用内网域名作为 misso 回调地址</li>
</ul>
</article></body>
