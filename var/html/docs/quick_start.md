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
<a id="user-content-falcon-qucik-start" class="anchor" href="#falcon-qucik-start" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>Falcon Qucik Start</h2>
<h2>
<a id="user-content-overview" class="anchor" href="#overview" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>overview</h2>
<p><a href="img/falcon-overview.svg?raw=true" target="_blank"><img src="img/falcon-overview.svg?raw=true" alt="" style="max-width:100%;"></a></p>
<h2>
<a id="user-content-install--start" class="anchor" href="#install--start" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>install &amp;&amp; start</h2>
<pre><code>#install gcc make automake libtool golang ...

# install protoc
wget https://github.com/google/protobuf/archive/v3.4.1.tar.gz
tar xzvf v3.4.1.tar.gz
cd protobuf-3.4.1
./autogen.sh
./configure
make
sudo make install

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/coreos/protodoc
go get -u github.com/beego/bee

# build falcon
git clone https://github.com/yubo/falcon
cd falcon
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
<a id="user-content-api-flow" class="anchor" href="#api-flow" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>API flow</h2>
<p><a href="img/falcon-api.svg?raw=true" target="_blank"><img src="img/falcon-api.svg?raw=true" alt="" style="max-width:100%;"></a></p>
<h2>
<a id="user-content-trigger" class="anchor" href="#trigger" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>trigger</h2>
<h4>
<a id="user-content-event" class="anchor" href="#event" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>event</h4>
<p><a href="img/falcon-event.svg?raw=true" target="_blank"><img src="img/falcon-event.svg?raw=true" alt="" style="max-width:100%;"></a></p>
<h4>
<a id="user-content-action" class="anchor" href="#action" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>action</h4>
<p><a href="img/falcon-action.svg?raw=true" target="_blank"><img src="img/falcon-action.svg?raw=true" alt="" style="max-width:100%;"></a></p>
<h2>
<a id="user-content-benchmark" class="anchor" href="#benchmark" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>benchmark</h2>
<pre><code>cd backend
go test -bench=Add -benchtime=20s
go test -bench=.
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
