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
<a id="user-content-install--start" class="anchor" href="#install--start" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>install &amp;&amp; start</h2>
<pre><code>## falcon
git clone git@git.n.xiaomi.com:falcon/falcon-lite.git
cd falcon-lite
make
make install

## nginx
sudo apt-get install nginx
mv /etc/nginx/conf.d/falcon.conf.example /etc/nginx/conf.d/falcon.conf
#edit /etc/nginx/conf.d/falcon.conf
#edit /etc/nginx/mime.types
  -    text/html                             html htm shtml;
  +    text/html                             html htm shtml md;

## mysql
sudp apt-get isntall mysql
mysql -u root -p &lt; ./scripts/db_schema/01_database.sql
mysql -u root -p &lt; ./scripts/db_schema/02_database_user.sql
mysql -u falcon -p1234 falcon &lt; ./scripts/db_schema/03_falcon.sql
mysql -u falcon -p1234 alarm &lt; ./scripts/db_schema/04_alarm.sql
mysql -u falcon -p1234 idx &lt; ./scripts/db_schema/05_index.sql

## start falcon
cp /etc/falcon/falcon.example.conf /etc/falcon/falcon.conf
sudo service falcon start
</code></pre>
<h2>
<a id="user-content-filedir-list" class="anchor" href="#filedir-list" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>file/dir list</h2>
<table>
<thead>
<tr>
<th>dir</th>
<th>desc</th>
</tr>
</thead>
<tbody>
<tr>
<td>/etc/falcon/falcon.conf</td>
<td>config</td>
</tr>
<tr>
<td>/etc/init.d/falcon</td>
<td>init.d script</td>
</tr>
<tr>
<td>/opt/falcon/log</td>
<td>log</td>
</tr>
<tr>
<td>/opt/falcon/tsdb</td>
<td>tsdb storage directry</td>
</tr>
<tr>
<td>/opt/falcon/emu_tpl</td>
<td>emulator template file directry</td>
</tr>
<tr>
<td>/opt/falcon/html</td>
<td>nginx document root</td>
</tr>
<tr>
<td>/sbin/falcon</td>
<td>falcon binary executable file(all module)</td>
</tr>
<tr>
<td>/sbin/agent</td>
<td>falcon-agent binary executable file(just single module)</td>
</tr>
</tbody>
</table>
</article></body>
