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
<a id="user-content-write-a-agent-plugin" class="anchor" href="#write-a-agent-plugin" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>write a agent plugin</h2>
<ol>
<li>first, define a struct for agent.Collector interface, and rigister it</li>
</ol>
<p>github.com/example/plugin/demo.cos.go</p>
<div class="highlight highlight-source-go"><pre><span class="pl-k">func</span> <span class="pl-en">init</span>() {
	agent.<span class="pl-c1">RegisterCollector</span>(&amp;cosCollector{})
}

<span class="pl-k">type</span> cosCollector <span class="pl-k">struct</span>{}

<span class="pl-k">func</span> <span class="pl-en">(<span class="pl-v">p</span> *<span class="pl-v">cosCollector</span>) <span class="pl-en">Name</span></span>() (<span class="pl-v">name</span>, <span class="pl-v">gname</span> <span class="pl-v">string</span>) {
	<span class="pl-k">return</span> <span class="pl-s"><span class="pl-pds">"</span>cos<span class="pl-pds">"</span></span>, <span class="pl-s"><span class="pl-pds">"</span>demo<span class="pl-pds">"</span></span>
}

<span class="pl-k">func</span> <span class="pl-en">(<span class="pl-v">p</span> *<span class="pl-v">cosCollector</span>) <span class="pl-en">Start</span></span>(<span class="pl-v">ctx</span> <span class="pl-v">context</span>.<span class="pl-v">Context</span>, <span class="pl-v">a</span> *<span class="pl-v">agent</span>.<span class="pl-v">Agent</span>) <span class="pl-v">error</span> {
	<span class="pl-smi">interval</span>, <span class="pl-smi">_</span> <span class="pl-k">:=</span> a.<span class="pl-smi">Conf</span>.<span class="pl-smi">Configer</span>.<span class="pl-c1">Int</span>(agent.<span class="pl-smi">C_INTERVAL</span>)
	<span class="pl-smi">ticker</span> <span class="pl-k">:=</span> time.<span class="pl-c1">NewTicker</span>(time.<span class="pl-smi">Second</span> * time.<span class="pl-c1">Duration</span>(interval)).<span class="pl-smi">C</span>
	<span class="pl-smi">ch</span> <span class="pl-k">:=</span> a.<span class="pl-smi">PutChan</span>

	<span class="pl-k">go</span> <span class="pl-k">func</span>() {
		<span class="pl-k">for</span> {
			<span class="pl-k">select</span> {
			<span class="pl-k">case</span> <span class="pl-k">&lt;-</span>ctx.<span class="pl-c1">Done</span>():
				<span class="pl-k">return</span>
			<span class="pl-k">case</span> <span class="pl-k">&lt;-</span>ticker:
				ch <span class="pl-k">&lt;-</span> &amp;agent.<span class="pl-smi">PutRequest</span>{
					Dps: []*transfer.<span class="pl-smi">DataPoint</span>{
						agent.<span class="pl-c1">GaugeValue</span>(<span class="pl-s"><span class="pl-pds">"</span>cos.ticker.demo<span class="pl-pds">"</span></span>,
							math.<span class="pl-c1">Cos</span>(x*(<span class="pl-k">float64</span>(time.<span class="pl-c1">Now</span>().<span class="pl-c1">Unix</span>())+<span class="pl-c1">1800</span>))),
					},
				}
			}
		}
	}()
	<span class="pl-k">return</span> <span class="pl-c1">nil</span>
}

<span class="pl-k">func</span> <span class="pl-en">(<span class="pl-v">p</span> *<span class="pl-v">cosCollector</span>) <span class="pl-en">Collect</span></span>() (<span class="pl-v">ret</span> []*<span class="pl-v">transfer</span>.<span class="pl-v">DataPoint</span>, <span class="pl-v">err</span> <span class="pl-v">error</span>) {
	<span class="pl-k">return</span> []*transfer.<span class="pl-smi">DataPoint</span>{
		agent.<span class="pl-c1">GaugeValue</span>(<span class="pl-s"><span class="pl-pds">"</span>cos.collect.demo<span class="pl-pds">"</span></span>,
			math.<span class="pl-c1">Cos</span>(x*<span class="pl-k">float64</span>(time.<span class="pl-c1">Now</span>().<span class="pl-c1">Unix</span>()))),
	}, <span class="pl-c1">nil</span>

}</pre></div>
<ol start="2">
<li>import in your code</li>
</ol>
<div class="highlight highlight-source-go"><pre><span class="pl-k">import</span> _ <span class="pl-s"><span class="pl-pds">"</span>github.com/example/plugin/demo.cos.go<span class="pl-pds">"</span></span>
</pre></div>
<ol start="3">
<li>add it's group name in the configuration file agent/plugins</li>
</ol>
<p>falcon/etc/falcon.conf</p>
<pre><code>agent {
	...
	plugins		"...,demo";
}

</code></pre>
<h2>
<a id="user-content-write" class="anchor" href="#write" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>write</h2>
<h2></h2>
<h2>
<a id="user-content-debug" class="anchor" href="#debug" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>DEBUG</h2>
<h4>
<a id="user-content-profile" class="anchor" href="#profile" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>profile</h4>
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
<p>see also <a href="https://blog.golang.org/profiling-go-programs" rel="nofollow">profiling-go-programs</a></p>
</article></body>
