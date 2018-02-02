webpackJsonp([1],{402:function(e,t,n){"use strict";function a(e){return e&&e.__esModule?e:{default:e}}Object.defineProperty(t,"__esModule",{value:!0});var r=n(39),o=a(r),s=n(33),l=a(s),u=n(88),d=a(u),i=n(18),c=a(i),f=n(48),p=a(f),E=n(23),h=a(E),m=n(41),C=a(m),S=n(77),g=a(S),y=n(6),R=a(y),v=n(116),w=a(v),_=n(21),T=a(_),A=n(2),I=a(A),N=n(7),D=a(N),b=n(5),x=a(b),O=n(4),L=a(O);n(40),n(34),n(115),n(19),n(32),n(89),n(42),n(78),n(174),n(22);var H=n(0),k=a(H),V=n(12),K=n(3),U=a(K),P=n(946);n(949);var B=o.default.Item,F=C.default.Option,q=function(e){function t(){var e,n,a,r;(0,I.default)(this,t);for(var o=arguments.length,s=Array(o),l=0;l<o;l++)s[l]=arguments[l];return n=a=(0,x.default)(this,(e=t.__proto__||Object.getPrototypeOf(t)).call.apply(e,[this].concat(s))),a.state={endpointInputBorderColor:"#d9d9d9",endpoint:"",tag:"",counter:"",searchEndpointsText:"",searchCountersText:"",selectedEndpointRowKeys:[],selectedCounterRowKeys:[],selectedEndpointRows:[],selectedCounterRows:[],endpointsLimit:"50",countersLimit:"50"},a.changeEndpoint=function(e){var t=e.target.value;a.setState({endpoint:t}),0===t.replace(/\s+/g,"").length?a.setState({endpointInputBorderColor:"red"}):a.setState({endpointInputBorderColor:"#d9d9d9"})},a.changeTag=function(e){a.setState({tag:e.target.value})},a.changeCounter=function(e){a.setState({counter:e.target.value})},a.onInputEndpointChange=function(e){a.setState({searchEndpointsText:e.target.value,selectedEndpointRowKeys:[],selectedEndpointRows:[]})},a.onInputCounterChange=function(e){a.setState({searchCountersText:e.target.value,selectedCounterRowKeys:[],selectedCounterRows:[]})},a.queryEndpoints=function(e){var t=a.state,n=t.endpoint,r=t.tag;if(""===n.replace(/[\s]/g,"")&&""===r.replace(/[\s]/g,""))return T.default.warning({title:"提示",content:"请填写有效Endpoint或者标签！"}),!1;a.setState({selectedEndpointRowKeys:[],selectedEndpointRows:[]}),"submit"===e.type?(e.preventDefault(),a.props.getEndpoints("getEndpoints",n,r,a.state.endpointsLimit)):(a.setState({endpointsLimit:e}),a.props.getEndpoints("getEndpoints",n,r,e))},a.queryCounters=function(e){var t=a.state.selectedEndpointRows;if(0===t.length)return T.default.warning({title:"提示",content:"请勾选要查看的endpoints"}),!1;a.setState({selectedCounterRowKeys:[],selectedCounterRows:[]});var n=t.map(function(e){return e.id}).join(),r=a.state.counter;"submit"===e.type?(e.preventDefault(),a.props.getCounters(n,r,a.state.countersLimit)):(a.setState({countersLimit:e}),a.props.getCounters(n,r,e))},a.showCharts=function(e,t){var n=a.state.selectedEndpointRows,r=a.state.selectedCounterRows;if(0===n.length||void 0===t&&0===r.length)return T.default.warning({title:"提示",content:"请勾选要查看的endpoints和counters！"}),!1;var o=n.map(function(e){return e.endpoint}),s=void 0===t?r.map(function(e){return e.counter}):[t],l=(0,U.default)().unix(),u=(0,U.default)().subtract(1,"hour").unix();a.props.saveSearchText({endpoint:a.state.endpoint,tag:a.state.tag,counter:a.state.counter,selectedEndpointRows:a.state.selectedEndpointRows,selectedCounterRows:a.state.selectedCounterRows}),a.props.getId(o,s,e,u,l,t)},r=n,(0,x.default)(a,r)}return(0,L.default)(t,e),(0,D.default)(t,[{key:"componentDidMount",value:function(){var e=this.props.dashboard.toJS(),t=e.searchText,n=t.endpoint,a=t.tag,r=t.counter,o=t.selectedEndpointRows,s=t.selectedCounterRows,l=0===o.length?[]:o.map(function(e){return e.id}),u=0===s.length?[]:s.map(function(e){return e.key});this.setState({endpoint:n,tag:a,counter:r,selectedEndpointRows:o,selectedCounterRows:s,selectedEndpointRowKeys:l,selectedCounterRowKeys:u})}},{key:"render",value:function(){var e=this,t=this.props.dashboard.toJS(),n=t.endpoints,a=t.counters,r=t.isLoading,s=[{title:"Endpoints",dataIndex:"endpoint",key:"id"}],u={selections:!0,selectedRowKeys:this.state.selectedEndpointRowKeys,onChange:function(t,n){e.setState({selectedEndpointRowKeys:t,selectedEndpointRows:n})},onSelectInvert:function(t){var a=t;0!==t.length&&(a=t.map(function(e){return n.filter(function(t){return t.id===e})[0]})),e.setState({selectedEndpointRows:a,selectedEndpointRowKeys:t})}},i={selections:!0,selectedRowKeys:this.state.selectedCounterRowKeys,onChange:function(t,n){e.setState({selectedCounterRowKeys:t,selectedCounterRows:n})},onSelectInvert:function(t){var n=t;0!==t.length&&(n=t.map(function(e){return a.filter(function(t){return t.key===e})[0]})),e.setState({selectedCounterRows:n,selectedCounterRowKeys:t})}},f=0,E=a.map(function(e){var t=e;return f+=1,t.key=f,t}),m=[{title:"Counters",dataIndex:"counter",key:"key",render:function(t,n){return k.default.createElement("a",{onClick:function(){return e.showCharts("Endpoint",n.counter)}},t)}},{title:"类型",dataIndex:"type",width:60},{title:"频率",dataIndex:"step",width:50}],S=k.default.createElement(w.default,null,k.default.createElement(w.default.Item,null,k.default.createElement(c.default,{type:"dashed",onClick:function(){return e.showCharts("Endpoint")}},"Endpoint视角")),k.default.createElement(w.default.Item,null,k.default.createElement(c.default,{type:"dashed",onClick:function(){return e.showCharts("Counter")}},"Counter视角")),k.default.createElement(w.default.Item,null,k.default.createElement(c.default,{type:"dashed",onClick:function(){return e.showCharts("组合")}},"组合视角"))),y={labelCol:{span:4},wrapperCol:{span:19}},v={labelCol:{span:2},wrapperCol:{span:21}},_={wrapperCol:{sm:{span:19,offset:4}}},T={wrapperCol:{sm:{span:21,offset:2}}};return k.default.createElement("div",{id:"dashboard-container"},k.default.createElement("div",{className:"endpoints"},k.default.createElement(o.default,{onSubmit:this.queryEndpoints},k.default.createElement("div",{className:"content"},k.default.createElement("div",{className:"header"},k.default.createElement("span",{className:"num"},"1"),k.default.createElement("span",{className:"title"},"搜索Endpoints")),k.default.createElement("div",{className:"up"},k.default.createElement("div",null,k.default.createElement(B,(0,R.default)({},y,{label:"Endpoint"}),k.default.createElement(p.default,{id:"txtDashBoardEndPoint",placeholder:"可以用空格分割多个搜索关键字",onChange:this.changeEndpoint,value:this.state.endpoint,style:{borderColor:this.state.endpointInputBorderColor}})),k.default.createElement(B,(0,R.default)({},y,{label:"标签"}),k.default.createElement(p.default,{id:"txtDashBoardTag",placeholder:"eg:job=appstore-web",value:this.state.tag,onChange:this.changeTag})),k.default.createElement(B,_,k.default.createElement(c.default,{type:"primary",htmlType:"submit"},"全局搜索")))),0===n.length?k.default.createElement("div",{className:"down"},"无数据"):k.default.createElement("div",null,k.default.createElement("div",{id:"table-header"},k.default.createElement(B,null,k.default.createElement(C.default,{value:this.state.endpointsLimit,style:{width:80,position:"relative",top:-1},onChange:this.queryEndpoints},k.default.createElement(F,{value:"50"},"Limit50"),k.default.createElement(F,{value:"100"},"Limit100"),k.default.createElement(F,{value:"500"},"Limit500"))),k.default.createElement(B,{className:"check"},k.default.createElement(p.default,{style:{width:200},placeholder:"请输入过滤信息",suffix:k.default.createElement(h.default,{type:"filter"}),value:this.state.searchEndpointsText,onChange:this.onInputEndpointChange}))),k.default.createElement("div",{id:"endpoints-results"},k.default.createElement(l.default,{className:"narrow-rows",rowKey:function(e){return e.id},rowSelection:u,columns:s,dataSource:n.filter(function(t){var n=e.state.searchEndpointsText;try{var a=new RegExp(n,"gi");return t.endpoint.match(a)}catch(e){return!0}}),pagination:!1})))))),0===r?k.default.createElement("div",{className:"dashboard-searchdata-loading"},k.default.createElement(g.default,{size:"large"})):"",k.default.createElement("div",{className:"counters"},k.default.createElement("p",{className:"header"},k.default.createElement("span",{className:"num"},"2"),k.default.createElement("span",{className:"title"},"搜索Counters")),k.default.createElement("div",{className:"up"},k.default.createElement(o.default,{onSubmit:this.queryCounters},k.default.createElement("div",{className:"content"},k.default.createElement("div",null,k.default.createElement(B,(0,R.default)({},v,{label:"Counter"}),k.default.createElement(p.default,{id:"txtDashBoardCounter",placeholder:"可以用空格分割多个搜索关键字",value:this.state.counter,onChange:this.changeCounter})),k.default.createElement(B,T,k.default.createElement(c.default,{type:"primary",htmlType:"submit"},"搜索")))))),0===a.length?k.default.createElement("div",{className:"down"},"无数据"):k.default.createElement("div",{id:"counters-results"},k.default.createElement(o.default,null,k.default.createElement("div",{className:"title"},k.default.createElement(B,null,k.default.createElement(C.default,{value:this.state.countersLimit,style:{width:80},onChange:this.queryCounters},k.default.createElement(F,{value:"50"},"Limit50"),k.default.createElement(F,{value:"100"},"Limit100"),k.default.createElement(F,{value:"500"},"Limit500"))),k.default.createElement(B,null,k.default.createElement(p.default,{style:{width:140},placeholder:"请输入过滤信息",suffix:k.default.createElement(h.default,{type:"filter"}),value:this.state.searchCountersText,onChange:this.onInputCounterChange})),k.default.createElement(B,{className:"check"},k.default.createElement(d.default,{overlay:S},k.default.createElement(c.default,{type:"default"},"看图")))),k.default.createElement(l.default,{className:"narrow-rows",rowSelection:i,columns:m,dataSource:E.filter(function(t){var n=e.state.searchCountersText.replace(/([.?*+^$[\]\\(){}|-])/g,"\\$1"),a=new RegExp(n,"gi");return t.counter.match(a)}),pagination:!1})))))}}]),t}(k.default.PureComponent),G={getEndpoints:P.getEndpoints,getCounters:P.getCounters,getCharts:P.getCharts,saveSearchText:P.saveSearchText,getId:P.getId},J=function(e){return{dashboard:e.dashboard}};t.default=(0,V.connect)(J,G)(q)},946:function(e,t,n){"use strict";function a(e){return e&&e.__esModule?e:{default:e}}function r(){return{type:"RESET_DASHBOARD_STATE"}}function o(e){return{type:"SAVE_SEARCH_DATA",payload:e}}function s(e){return{type:"SAVE_SEARCH_TEXT",searchText:e}}function l(e){return{type:R,payload:{endpoints:e}}}function u(e){return{type:v,counters:e}}function d(e){return{type:w,chartsData:e}}function i(e,t,n,a){return function(r){return(0,y.fetch)(e,{params:{tag:n,query:t,limit:a}}).then(function(e){if(r(l(e)),0===e.length)return C.default.warning({title:"提示",content:"响应数据是空！"}),!1}).catch(function(){C.default.warning({title:"提示",content:"endpoints获取失败！"})})}}function c(e,t,n){return function(a){return a({type:"START_REQUEST_COUNTERS"}),(0,y.fetch)("getCounters",{params:{query:t,limit:n,ids:e}}).then(function(e){if(a(u(e)),0===e.length)return C.default.warning({title:"提示",content:"响应数据是空！"}),!1}).catch(function(){a({type:"RECEIVE_FAIL_COUNTERS"}),C.default.warning({title:"提示",content:"counters获取失败！"})})}}function f(e,t,n,a,r,o){return _=0,T=0,function(s){s({type:"CLEAR_DASHBOARD_VIEWS"}),s({type:"START_REQUEST_CHARTSDATA"});var l=[];if("Endpoint"===a){_=n.length;for(var u=_;u;)l.push({rsp:[]}),u-=1;n.sort().forEach(function(n,a){var u={consol_fun:e,counters:[n],end_time:o,hostnames:t,start_time:r};l[a].params=u,(0,y.fetch)("getCharts",{body:u}).then(function(e){l[a].rsp=e,s(d(l)),T+=1,_===T&&s({type:"FINISHED_RESPONSE"})}).catch(function(){T+=1,_===T&&s({type:"FINISHED_RESPONSE"})})})}else if("Counter"===a){_=t.length;for(var i=_;i;)l.push({rsp:[]}),i-=1;var c=t.sort();c.forEach(function(t,a){var u={consol_fun:e,counters:n,end_time:o,hostnames:[t],start_time:r};l[a].params=u,(0,y.fetch)("getCharts",{body:u}).then(function(e){l[a].rsp=e,s(d(l)),T+=1,_===T&&s({type:"FINISHED_RESPONSE"})}).catch(function(){T+=1,_===T&&s({type:"FINISHED_RESPONSE"})})})}else{_=1;var f={consol_fun:e,counters:n,end_time:o,hostnames:t,start_time:r},p=JSON.stringify(f);l.push({rsp:[],params:f}),(0,y.fetch)("getCharts",p).then(function(e){l[0].rsp=e,s(d(l)),T+=1,_===T&&s({type:"FINISHED_RESPONSE"})}).catch(function(){T+=1,_===T&&s({type:"FINISHED_RESPONSE"})})}}}function p(e,t,n,a,r,s){return function(l){var u=(JSON.stringify({endpoints:e,counters:t}),window.open("about:blank","_blank"));(0,y.fetch)("tmpGraphAdd",{body:{endpoints:e,counters:t}}).then(function(d){var i={id:d.id,endpoints:e,counters:t,title:n,startTime:a,endTime:r,cf:"AVERAGE",method:"nosum"};l(o(i));var c="h";if("Counter"===n&&(c="k"),"组合"===n&&(c="a"),void 0===s){var f="id="+d.id+"&graph_type="+c+"&cf=AVERAGE&start="+a+"&end="+r;return u.location.href="/dashboard/charts?"+f,!1}var p="id="+d.id+"&graph_type=h&cf=AVERAGE&start=-3600";u.location.href="/chart?"+p}).catch(function(){C.default.warning({title:"提示",content:"响应失败！"})})}}function E(e,t,n,a,r,s){return function(l){(0,y.fetch)("tmpGraphGet",{path:"/"+e}).then(function(u){var d=u.endpoints,i=u.counters,c="Endpoint";"k"===t&&(c="Counter"),"a"===t&&(c="组合"),l(o({id:e,endpoints:d,counters:i,title:c,startTime:n,endTime:a,cf:r,method:s})),l(f(r,d,i,c,n,a))}).catch(function(){C.default.warning({title:"提示",content:"响应失败！"})})}}function h(e,t){return"RESET_DASHBOARD_STATE"===t.type&&(e=void 0),I(e,t)}Object.defineProperty(t,"__esModule",{value:!0});var m=n(21),C=a(m);t.resetDashboardState=r,t.saveSearchData=o,t.saveSearchText=s,t.getEndpoints=i,t.getCounters=c,t.getCharts=f,t.getId=p,t.getChartParams=E,t.default=h,n(22);var S=n(173),g=a(S),y=n(16),R="RECEIVE_SUCCESS_ENDPOINTS",v="RECEIVE_SUCCESS_COUNTERS",w="RECEIVE_SUCCESS_CHARTSDATA",_=0,T=0,A=g.default.fromJS({endpoints:[],counters:[],searchText:{endpoint:"",tag:"",counter:"",selectedEndpointRows:[],selectedCounterRows:[]},selectedEndpoints:[],selectedCounters:[],title:"",start:"",end:"",cf:"AVERAGE",method:"nosum",chartsData:[],isLoading:-1}),I=function(){var e=arguments.length>0&&void 0!==arguments[0]?arguments[0]:A,t=arguments[1];switch(t.type){case"RECEIVE_SUCCESS_ENDPOINTS":return e.set("endpoints",t.payload.endpoints);case"START_REQUEST_COUNTERS":case"START_REQUEST_CHARTSDATA":return e.set("isLoading",0);case"RECEIVE_SUCCESS_COUNTERS":return e.merge({isLoading:1,counters:t.counters});case"RECEIVE_FAIL_COUNTERS":return e.set("isLoading",1);case"SAVE_SEARCH_TEXT":return e.set("searchText",t.searchText);case"SAVE_SEARCH_DATA":return e.merge({id:t.payload.id,selectedEndpoints:t.payload.endpoints,selectedCounters:t.payload.counters,title:t.payload.title,start:t.payload.startTime,end:t.payload.endTime,cf:t.payload.cf,method:t.payload.method});case"RECEIVE_SUCCESS_CHART_PARAMS":return e.set("searchData",t.payload);case"CLEAR_DASHBOARD_VIEWS":return e.set("chartsData",[]);case"RECEIVE_SUCCESS_CHARTSDATA":return e.set("chartsData",t.chartsData);case"FINISHED_RESPONSE":return e.set("isLoading",1);default:return e}}},949:function(e,t){}});
//# sourceMappingURL=1.27a9529f1e8ebf479d9c.js.map