webpackJsonp([2],{946:function(t,e,n){"use strict";function r(t){return t&&t.__esModule?t:{default:t}}function a(){return{type:"RESET_DASHBOARD_STATE"}}function o(t){return{type:"SAVE_SEARCH_DATA",payload:t}}function s(t){return{type:"SAVE_SEARCH_TEXT",searchText:t}}function i(t){return{type:T,payload:{endpoints:t}}}function c(t){return{type:y,counters:t}}function u(t){return{type:m,chartsData:t}}function E(t,e,n,r){return function(a){return(0,R.fetch)(t,{params:{tag:n,query:e,limit:r}}).then(function(t){if(a(i(t)),0===t.length)return _.default.warning({title:"提示",content:"响应数据是空！"}),!1}).catch(function(){_.default.warning({title:"提示",content:"endpoints获取失败！"})})}}function d(t,e,n){return function(r){return r({type:"START_REQUEST_COUNTERS"}),(0,R.fetch)("getCounters",{params:{query:e,limit:n,ids:t}}).then(function(t){if(r(c(t)),0===t.length)return _.default.warning({title:"提示",content:"响应数据是空！"}),!1}).catch(function(){r({type:"RECEIVE_FAIL_COUNTERS"}),_.default.warning({title:"提示",content:"counters获取失败！"})})}}function f(t,e,n,r,a,o){return g=0,D=0,function(s){s({type:"CLEAR_DASHBOARD_VIEWS"}),s({type:"START_REQUEST_CHARTSDATA"});var i=[];if("Endpoint"===r){g=n.length;for(var c=g;c;)i.push({rsp:[]}),c-=1;n.sort().forEach(function(n,r){var c={consol_fun:t,counters:[n],end_time:o,hostnames:e,start_time:a};i[r].params=c,(0,R.fetch)("getCharts",{body:c}).then(function(t){i[r].rsp=t,s(u(i)),D+=1,g===D&&s({type:"FINISHED_RESPONSE"})}).catch(function(){D+=1,g===D&&s({type:"FINISHED_RESPONSE"})})})}else if("Counter"===r){g=e.length;for(var E=g;E;)i.push({rsp:[]}),E-=1;var d=e.sort();d.forEach(function(e,r){var c={consol_fun:t,counters:n,end_time:o,hostnames:[e],start_time:a};i[r].params=c,(0,R.fetch)("getCharts",{body:c}).then(function(t){i[r].rsp=t,s(u(i)),D+=1,g===D&&s({type:"FINISHED_RESPONSE"})}).catch(function(){D+=1,g===D&&s({type:"FINISHED_RESPONSE"})})})}else{g=1;var f={consol_fun:t,counters:n,end_time:o,hostnames:e,start_time:a},S=JSON.stringify(f);i.push({rsp:[],params:f}),(0,R.fetch)("getCharts",S).then(function(t){i[0].rsp=t,s(u(i)),D+=1,g===D&&s({type:"FINISHED_RESPONSE"})}).catch(function(){D+=1,g===D&&s({type:"FINISHED_RESPONSE"})})}}}function S(t,e,n,r,a,s){return function(i){var c=(JSON.stringify({endpoints:t,counters:e}),window.open("about:blank","_blank"));(0,R.fetch)("tmpGraphAdd",{body:{endpoints:t,counters:e}}).then(function(u){var E={id:u.id,endpoints:t,counters:e,title:n,startTime:r,endTime:a,cf:"AVERAGE",method:"nosum"};i(o(E));var d="h";if("Counter"===n&&(d="k"),"组合"===n&&(d="a"),void 0===s){var f="id="+u.id+"&graph_type="+d+"&cf=AVERAGE&start="+r+"&end="+a;return c.location.href="/dashboard/charts?"+f,!1}var S="id="+u.id+"&graph_type=h&cf=AVERAGE&start=-3600";c.location.href="/chart?"+S}).catch(function(){_.default.warning({title:"提示",content:"响应失败！"})})}}function p(t,e,n,r,a,s){return function(i){(0,R.fetch)("tmpGraphGet",{path:"/"+t}).then(function(c){var u=c.endpoints,E=c.counters,d="Endpoint";"k"===e&&(d="Counter"),"a"===e&&(d="组合"),i(o({id:t,endpoints:u,counters:E,title:d,startTime:n,endTime:r,cf:a,method:s})),i(f(a,u,E,d,n,r))}).catch(function(){_.default.warning({title:"提示",content:"响应失败！"})})}}function h(t,e){return"RESET_DASHBOARD_STATE"===e.type&&(t=void 0),v(t,e)}Object.defineProperty(e,"__esModule",{value:!0});var l=n(21),_=r(l);e.resetDashboardState=a,e.saveSearchData=o,e.saveSearchText=s,e.getEndpoints=E,e.getCounters=d,e.getCharts=f,e.getId=S,e.getChartParams=p,e.default=h,n(22);var A=n(173),C=r(A),R=n(16),T="RECEIVE_SUCCESS_ENDPOINTS",y="RECEIVE_SUCCESS_COUNTERS",m="RECEIVE_SUCCESS_CHARTSDATA",g=0,D=0,I=C.default.fromJS({endpoints:[],counters:[],searchText:{endpoint:"",tag:"",counter:"",selectedEndpointRows:[],selectedCounterRows:[]},selectedEndpoints:[],selectedCounters:[],title:"",start:"",end:"",cf:"AVERAGE",method:"nosum",chartsData:[],isLoading:-1}),v=function(){var t=arguments.length>0&&void 0!==arguments[0]?arguments[0]:I,e=arguments[1];switch(e.type){case"RECEIVE_SUCCESS_ENDPOINTS":return t.set("endpoints",e.payload.endpoints);case"START_REQUEST_COUNTERS":case"START_REQUEST_CHARTSDATA":return t.set("isLoading",0);case"RECEIVE_SUCCESS_COUNTERS":return t.merge({isLoading:1,counters:e.counters});case"RECEIVE_FAIL_COUNTERS":return t.set("isLoading",1);case"SAVE_SEARCH_TEXT":return t.set("searchText",e.searchText);case"SAVE_SEARCH_DATA":return t.merge({id:e.payload.id,selectedEndpoints:e.payload.endpoints,selectedCounters:e.payload.counters,title:e.payload.title,start:e.payload.startTime,end:e.payload.endTime,cf:e.payload.cf,method:e.payload.method});case"RECEIVE_SUCCESS_CHART_PARAMS":return t.set("searchData",e.payload);case"CLEAR_DASHBOARD_VIEWS":return t.set("chartsData",[]);case"RECEIVE_SUCCESS_CHARTSDATA":return t.set("chartsData",e.chartsData);case"FINISHED_RESPONSE":return t.set("isLoading",1);default:return t}}},948:function(t,e,n){"use strict";function r(t){return t&&t.__esModule?t:{default:t}}Object.defineProperty(e,"__esModule",{value:!0});var a=n(2),o=r(a),s=n(7),i=r(s),c=n(5),u=r(c),E=n(4),d=r(E),f=n(0),S=r(f),p=function(t){function e(){return(0,o.default)(this,e),(0,u.default)(this,(e.__proto__||Object.getPrototypeOf(e)).apply(this,arguments))}return(0,d.default)(e,t),(0,i.default)(e,[{key:"render",value:function(){return S.default.createElement("div",null,this.props.children)}}]),e}(S.default.PureComponent);e.default=p}});
//# sourceMappingURL=2.3ff86763bdaac636d934.js.map