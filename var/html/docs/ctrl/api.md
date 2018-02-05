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
<a id="user-content-接口设计" class="anchor" href="#%E6%8E%A5%E5%8F%A3%E8%AE%BE%E8%AE%A1" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>接口设计</h2>
<h4>
<a id="user-content-ctrl-接口使用http-restful-架构" class="anchor" href="#ctrl-%E6%8E%A5%E5%8F%A3%E4%BD%BF%E7%94%A8http-restful-%E6%9E%B6%E6%9E%84" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>ctrl 接口使用http RESTful 架构</h4>
<table>
<thead>
<tr>
<th>api</th>
<th>method</th>
</tr>
</thead>
<tbody>
<tr>
<td>create</td>
<td>POST</td>
</tr>
<tr>
<td>delet</td>
<td>DELETE</td>
</tr>
<tr>
<td>edit</td>
<td>PUT</td>
</tr>
<tr>
<td>list/search/get</td>
<td>GET</td>
</tr>
</tbody>
</table>
<h4>
<a id="user-content-接口的数据类型和命名规则以user为例" class="anchor" href="#%E6%8E%A5%E5%8F%A3%E7%9A%84%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B%E5%92%8C%E5%91%BD%E5%90%8D%E8%A7%84%E5%88%99%E4%BB%A5user%E4%B8%BA%E4%BE%8B" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>接口的数据类型和命名规则,以user为例</h4>
<table>
<thead>
<tr>
<th>api</th>
<th>struct name</th>
</tr>
</thead>
<tbody>
<tr>
<td>create</td>
<td>UserCreate</td>
</tr>
<tr>
<td>edit</td>
<td>UserUpdate</td>
</tr>
<tr>
<td>get</td>
<td>User</td>
</tr>
</tbody>
</table>
</article></body>
