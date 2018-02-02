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
<a id="user-content-falcon-acl" class="anchor" href="#falcon-acl" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>falcon acl</h2>
<ul>
<li>global-r <code>在任何节点拥有 falcon-read token</code>
</li>
<li>global-w <code>在任何节点拥有 falcon-opterate token</code>
</li>
<li>global-a <code>在任何节点拥有 falcon-admin token</code>
</li>
<li>tag-r <code>在当前节点或上级节点拥有 falcon-read token</code>
</li>
<li>tag-w <code>在当前节点或上级节点拥有 falcon-opterate token</code>
</li>
<li>tag-a <code>在当前节点或上级节点拥有 falcon-admin token</code>
</li>
<li>owner <code>自己创建的对象</code>
</li>
<li>获取用户global权限接口 <code>http://dev02:8001/v1.0/auth/info</code>
<ul>
<li>"admin": global-a</li>
<li>"operator": global-w</li>
<li>"reader":  global-r</li>
</ul>
</li>
</ul>
<hr>
<table>
<thead>
<tr>
<th>module</th>
<th>add/edit/del</th>
<th>list/get/search</th>
</tr>
</thead>
<tbody>
<tr>
<td>tag</td>
<td>tag-w(parent)</td>
<td>tag-r</td>
</tr>
<tr>
<td>tag-tpl</td>
<td>tag-w</td>
<td>global-r</td>
</tr>
<tr>
<td>tag-host</td>
<td>tag-w</td>
<td>global-r</td>
</tr>
<tr>
<td>aggreator</td>
<td>tag-w</td>
<td>global-r</td>
</tr>
<tr>
<td>plugin</td>
<td>tag-w</td>
<td>global-r</td>
</tr>
<tr>
<td>alarm</td>
<td>owner(edit/del) / global-w(add)</td>
<td>global-r</td>
</tr>
<tr>
<td>template</td>
<td>owner(edit/del) / global-w(add)</td>
<td>global-r</td>
</tr>
<tr>
<td>expression</td>
<td>owner(edit/del) / global-w(add)</td>
<td>global-r</td>
</tr>
<tr>
<td>nodata</td>
<td>owner(edit/del) / global-w(add)</td>
<td>global-r</td>
</tr>
<tr>
<td>team</td>
<td>owner(edit/del) / global-w(add)</td>
<td>global-r</td>
</tr>
<tr>
<td>dashboard</td>
<td>global-r</td>
<td>global-r</td>
</tr>
<tr>
<td>admin, rel, token, role, user</td>
<td>tag-a</td>
<td>global-a</td>
</tr>
</tbody>
</table>
</article></body>
