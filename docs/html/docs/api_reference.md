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

<h3>
<a id="user-content-falcon-api-reference" class="anchor" href="#falcon-api-reference" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>falcon API Reference</h3>
<p>This is a generated documentation. Please read the proto files for more.</p>
<h5>
<a id="user-content-message-datapoint-tsdbtsdbproto" class="anchor" href="#message-datapoint-tsdbtsdbproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>DataPoint</code> (tsdb/tsdb.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>key</td>
<td></td>
<td>Key</td>
</tr>
<tr>
<td>value</td>
<td></td>
<td>TimeValuePair</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-datapoints-tsdbtsdbproto" class="anchor" href="#message-datapoints-tsdbtsdbproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>DataPoints</code> (tsdb/tsdb.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>key</td>
<td></td>
<td>Key</td>
</tr>
<tr>
<td>values</td>
<td></td>
<td>(slice of) TimeValuePair</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-key-tsdbtsdbproto" class="anchor" href="#message-key-tsdbtsdbproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Key</code> (tsdb/tsdb.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>key</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>shardId</td>
<td></td>
<td>int32</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-timevaluepair-tsdbtsdbproto" class="anchor" href="#message-timevaluepair-tsdbtsdbproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>TimeValuePair</code> (tsdb/tsdb.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>timestamp</td>
<td></td>
<td>int64</td>
</tr>
<tr>
<td>value</td>
<td></td>
<td>double</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-service-agent-agentagentproto" class="anchor" href="#service-agent-agentagentproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>service <code>Agent</code> (agent/agent.proto)</h5>
<table>
<thead>
<tr>
<th>Method</th>
<th>Request Type</th>
<th>Response Type</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>Put</td>
<td>PutRequest</td>
<td>PutResponse</td>
<td></td>
</tr>
<tr>
<td>GetStats</td>
<td>Empty</td>
<td>Stats</td>
<td></td>
</tr>
<tr>
<td>GetStatsName</td>
<td>Empty</td>
<td>StatsName</td>
<td></td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-empty-agentagentproto" class="anchor" href="#message-empty-agentagentproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Empty</code> (agent/agent.proto)</h5>
<p>Empty field.</p>
<h5>
<a id="user-content-message-item-agentagentproto" class="anchor" href="#message-item-agentagentproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Item</code> (agent/agent.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>metric</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>tags</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>type</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>value</td>
<td></td>
<td>double</td>
</tr>
<tr>
<td>timestamp</td>
<td></td>
<td>int64</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putrequest-agentagentproto" class="anchor" href="#message-putrequest-agentagentproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutRequest</code> (agent/agent.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>items</td>
<td></td>
<td>(slice of) Item</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putresponse-agentagentproto" class="anchor" href="#message-putresponse-agentagentproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutResponse</code> (agent/agent.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>n</td>
<td></td>
<td>int32</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-stats-agentagentproto" class="anchor" href="#message-stats-agentagentproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Stats</code> (agent/agent.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counter</td>
<td></td>
<td>(slice of) uint64</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-statsname-agentagentproto" class="anchor" href="#message-statsname-agentagentproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>StatsName</code> (agent/agent.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counterName</td>
<td></td>
<td>(slice of) bytes</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-service-alarm-alarmalarmproto" class="anchor" href="#service-alarm-alarmalarmproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>service <code>Alarm</code> (alarm/alarm.proto)</h5>
<table>
<thead>
<tr>
<th>Method</th>
<th>Request Type</th>
<th>Response Type</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>Put</td>
<td>PutRequest</td>
<td>PutResponse</td>
<td></td>
</tr>
<tr>
<td>GetStats</td>
<td>Empty</td>
<td>Stats</td>
<td></td>
</tr>
<tr>
<td>GetStatsName</td>
<td>Empty</td>
<td>StatsName</td>
<td></td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-empty-alarmalarmproto" class="anchor" href="#message-empty-alarmalarmproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Empty</code> (alarm/alarm.proto)</h5>
<p>Empty field.</p>
<h5>
<a id="user-content-message-event-alarmalarmproto" class="anchor" href="#message-event-alarmalarmproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Event</code> (alarm/alarm.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>tagId</td>
<td></td>
<td>int64</td>
</tr>
<tr>
<td>key</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>expr</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>msg</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>timestamp</td>
<td></td>
<td>int64</td>
</tr>
<tr>
<td>value</td>
<td></td>
<td>double</td>
</tr>
<tr>
<td>Priority</td>
<td></td>
<td>int32</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putrequest-alarmalarmproto" class="anchor" href="#message-putrequest-alarmalarmproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutRequest</code> (alarm/alarm.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>events</td>
<td></td>
<td>(slice of) Event</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putresponse-alarmalarmproto" class="anchor" href="#message-putresponse-alarmalarmproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutResponse</code> (alarm/alarm.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>n</td>
<td></td>
<td>int32</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-stats-alarmalarmproto" class="anchor" href="#message-stats-alarmalarmproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Stats</code> (alarm/alarm.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counter</td>
<td></td>
<td>(slice of) uint64</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-statsname-alarmalarmproto" class="anchor" href="#message-statsname-alarmalarmproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>StatsName</code> (alarm/alarm.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counterName</td>
<td></td>
<td>(slice of) bytes</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-service-service-serviceserviceproto" class="anchor" href="#service-service-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>service <code>Service</code> (service/service.proto)</h5>
<table>
<thead>
<tr>
<th>Method</th>
<th>Request Type</th>
<th>Response Type</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>Put</td>
<td>PutRequest</td>
<td>PutResponse</td>
<td></td>
</tr>
<tr>
<td>Get</td>
<td>GetRequest</td>
<td>GetResponse</td>
<td></td>
</tr>
<tr>
<td>GetStats</td>
<td>Empty</td>
<td>Stats</td>
<td></td>
</tr>
<tr>
<td>GetStatsName</td>
<td>Empty</td>
<td>StatsName</td>
<td></td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-empty-serviceserviceproto" class="anchor" href="#message-empty-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Empty</code> (service/service.proto)</h5>
<p>Empty field.</p>
<h5>
<a id="user-content-message-getrequest-serviceserviceproto" class="anchor" href="#message-getrequest-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>GetRequest</code> (service/service.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>keys</td>
<td></td>
<td>(slice of) tsdb.Key</td>
</tr>
<tr>
<td>start</td>
<td></td>
<td>int64</td>
</tr>
<tr>
<td>end</td>
<td></td>
<td>int64</td>
</tr>
<tr>
<td>consolFun</td>
<td></td>
<td>Cf</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-getresponse-serviceserviceproto" class="anchor" href="#message-getresponse-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>GetResponse</code> (service/service.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>data</td>
<td></td>
<td>(slice of) tsdb.DataPoints</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putrequest-serviceserviceproto" class="anchor" href="#message-putrequest-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutRequest</code> (service/service.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>data</td>
<td></td>
<td>(slice of) tsdb.DataPoint</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putresponse-serviceserviceproto" class="anchor" href="#message-putresponse-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutResponse</code> (service/service.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>n</td>
<td></td>
<td>int32</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-stats-serviceserviceproto" class="anchor" href="#message-stats-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Stats</code> (service/service.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counter</td>
<td></td>
<td>(slice of) uint64</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-statsname-serviceserviceproto" class="anchor" href="#message-statsname-serviceserviceproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>StatsName</code> (service/service.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counterName</td>
<td></td>
<td>(slice of) bytes</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-service-transfer-transfertransferproto" class="anchor" href="#service-transfer-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>service <code>Transfer</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Method</th>
<th>Request Type</th>
<th>Response Type</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>Put</td>
<td>PutRequest</td>
<td>PutResponse</td>
<td></td>
</tr>
<tr>
<td>Get</td>
<td>GetRequest</td>
<td>GetResponse</td>
<td></td>
</tr>
<tr>
<td>GetStats</td>
<td>Empty</td>
<td>Stats</td>
<td></td>
</tr>
<tr>
<td>GetStatsName</td>
<td>Empty</td>
<td>StatsName</td>
<td></td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-datapoint-transfertransferproto" class="anchor" href="#message-datapoint-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>DataPoint</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>key</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>value</td>
<td></td>
<td>tsdb.TimeValuePair</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-datapoints-transfertransferproto" class="anchor" href="#message-datapoints-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>DataPoints</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>key</td>
<td></td>
<td>bytes</td>
</tr>
<tr>
<td>values</td>
<td></td>
<td>(slice of) tsdb.TimeValuePair</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-empty-transfertransferproto" class="anchor" href="#message-empty-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Empty</code> (transfer/transfer.proto)</h5>
<p>Empty field.</p>
<h5>
<a id="user-content-message-getrequest-transfertransferproto" class="anchor" href="#message-getrequest-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>GetRequest</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>keys</td>
<td></td>
<td>(slice of) bytes</td>
</tr>
<tr>
<td>start</td>
<td></td>
<td>int64</td>
</tr>
<tr>
<td>end</td>
<td></td>
<td>int64</td>
</tr>
<tr>
<td>consolFun</td>
<td></td>
<td>service.Cf</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-getresponse-transfertransferproto" class="anchor" href="#message-getresponse-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>GetResponse</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>data</td>
<td></td>
<td>(slice of) DataPoints</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putrequest-transfertransferproto" class="anchor" href="#message-putrequest-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutRequest</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>data</td>
<td></td>
<td>(slice of) DataPoint</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-putresponse-transfertransferproto" class="anchor" href="#message-putresponse-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>PutResponse</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>n</td>
<td></td>
<td>int32</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-stats-transfertransferproto" class="anchor" href="#message-stats-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>Stats</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counter</td>
<td></td>
<td>(slice of) uint64</td>
</tr>
</tbody>
</table>
<h5>
<a id="user-content-message-statsname-transfertransferproto" class="anchor" href="#message-statsname-transfertransferproto" aria-hidden="true"><span aria-hidden="true" class="octicon octicon-link"></span></a>message <code>StatsName</code> (transfer/transfer.proto)</h5>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
<th>Type</th>
</tr>
</thead>
<tbody>
<tr>
<td>counterName</td>
<td></td>
<td>(slice of) bytes</td>
</tr>
</tbody>
</table>
</article></body>
