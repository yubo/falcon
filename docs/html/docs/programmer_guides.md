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
<a id="user-content-profile" class="anchor" href="#profile" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>profile</h2>
<pre><code>falcon start -cpu/heap
^C

$go tool pprof  /tmp/cpu.prof 
Entering interactive mode (type "help" for commands)
(pprof) top
46220ms of 46350ms total (99.72%)
Dropped 94 nodes (cum &lt;= 231.75ms)
      flat  flat%   sum%        cum   cum%
   14880ms 32.10% 32.10%    14880ms 32.10%  runtime.chanrecv
   12690ms 27.38% 59.48%    27570ms 59.48%  runtime.selectnbrecv
   10080ms 21.75% 81.23%    10080ms 21.75%  github.com/yubo/falcon/lib/tsdb.(*KeyListWriter).writeOneKey
    7860ms 16.96% 98.19%    46220ms 99.72%  github.com/yubo/falcon/lib/tsdb.(*KeyListWriter).startWriterThread.func1
     710ms  1.53% 99.72%      710ms  1.53%  github.com/yubo/falcon/lib/tsdb.(*KeyListWriter).startWriterThread
         0     0% 99.72%    46330ms   100%  runtime.goexit
</code></pre>
<p>see also</p>
<pre><code>https://blog.golang.org/profiling-go-programs
</code></pre>
</article></body>
