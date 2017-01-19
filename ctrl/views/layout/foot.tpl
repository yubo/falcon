      <script type="text/javascript">

function of_emsg(msg, time, fn) {
	layer.open({
		icon: 2
		,time: time
		,content: msg
		,end: fn
	});
}

function of_imsg(msg, time, fn) {
	layer.open({
		icon: 1
		,time: time 
		,content: msg
		,end: fn
	});
}

function edit_meta(meta, id, method){
	jQuery[method]('/v1.0/'+meta+'/'+id, JSON.stringify(
		$('#form-meta').serializeObject()
	), function(json) {
		if (json.code != 200) {
			of_emsg(json.msg, 2000);
		} else {
			of_imsg(method + " " + meta + " success", 1000);
		}
	});
}

function delete_meta(meta, id, url) {
	layer.confirm("真的要删除么？", {
		icon: 3,
		btn: [ '确定', '取消' ]
	}, function() {
		$.delete('/v1.0/'+meta+'/'+id, {}, function(json) {
			if (json.code != 200) {
				of_emsg(json.msg, 2000);
			} else {
				of_imsg('delete '+meta+' successfully', 1000,
					 function(){refresh_portion(url);});
			}
		});
	}, function(){});
}

function edit_team_users(id, method) {
	jQuery[method]('/v1.0/team/'+id, JSON.stringify({
		"team": {
			"name": $("#name").val(),
			"note": $("#note").val()
		},
		"user_ids": $("#users").val().split(',').map(function(item){
			return parseInt(item, 10);
			})
	}), function(json) {
		if (json.code != 200) {
			of_emsg(json.msg, 2000);
		} else {
			of_imsg(method + " team  success", 1000);
		}
	});
}

function bind_tag_hosts() {
	if (!curTagId){
		of_emsg("请选取节点");
		return
	}
	$.post('/v1.0/rel/tag/hosts', JSON.stringify({
		'tag_id': parseInt(curTagId),
		"host_ids": $("#hosts").val().split(',').map(function(item){
			return parseInt(item, 10);
			})
	}), function(json) {
		if (json.code != 200) {
			of_emsg(json.msg);
		} else {
			$("#hosts").select2("val", "");
			of_imsg('bind success ('+json.data+')', 0,
				function(){refresh_portion(curUrl);});
		}
	});
}

function bind_role_tokens() {
	if (!curTagId){
		of_emsg("请选取节点");
		return
	}
	$.post('/v1.0/rel/tag/role/tokens', JSON.stringify({
		'tag_id': parseInt(curTagId),
		'role_id': parseInt($("#role").val()),
		"token_ids": $("#tokens").val().split(',').map(function(item){
			return parseInt(item, 10);
			})
	}), function(json) {
		if (json.code != 200) {
			of_emsg(json.msg);
		} else {
			$("#tokens").select2("val", "");
			of_imsg('bind success ('+json.data+')', 0,
				function(){refresh_portion(curUrl);});
		}
	});
}

function bind_role_users() {
	if (!curTagId){
		of_emsg("请选取节点");
		return
	}
	$.post('/v1.0/rel/tag/role/users', JSON.stringify({
		'tag_id': parseInt(curTagId),
		'role_id': parseInt($("#role").val()),
		"user_ids": $("#users").val().split(',').map(function(item){
			return parseInt(item, 10);
			})
	}), function(json) {
		if (json.code != 200) {
			of_emsg(json.msg);
		} else {
			$("#users").select2("val", "");
			of_imsg('bind success ('+json.data+')', 0,
				function(){refresh_portion(curUrl);});
		}
	});
}

function bind_rule_triggers() {
	if (!curTagId){
		of_emsg("请选取节点");
		return
	}
	$.post('/v1.0/rel/tag/rule/users', JSON.stringify({
		'tag_id': parseInt(curTagId),
		'rule_id': parseInt($("#rule").val()),
		"trigger_ids": $("#triggers").val().split(',').map(function(item){
			return parseInt(item, 10);
			})
	}), function(json) {
		if (json.code != 200) {
			of_emsg(json.msg);
		} else {
			$("#triggers").select2("val", "");
			of_imsg('bind success ('+json.data+')', 0,
				function(){refresh_portion(curUrl);});
		}
	});
}


function unbind_tag_host(tag_id, host_id, url) {
	layer.confirm("真的要删除么？", {
		icon: 3,
		btn: [ '确定', '取消' ]
	}, function() {
		$.delete('/v1.0/rel/tag/host', JSON.stringify({
				"tag_id": parseInt(tag_id),
				"host_id": parseInt(host_id)
		}), function(json) {
			if (json.code != 200) {
				of_emsg(json.msg, 2000);
			} else {
				of_imsg('delete tag('+tag_id+')-host('
					+host_id+') successfully', 1000,
					function(){ refresh_portion(url); });
			}
		});
	}, function(){});
}

function unbind_tag_role_user(tag_id, role_id, user_id, url) {
	layer.confirm("真的要删除么？", {
		icon: 3,
		btn: [ '确定', '取消' ]
	}, function() {
		$.delete('/v1.0/rel/tag/role/user', JSON.stringify({
				"tag_id": parseInt(tag_id),
				"role_id": parseInt(role_id),
				"user_id": parseInt(user_id)
		}), function(json) {
			if (json.code != 200) {
				of_emsg(json.msg, 2000);
			} else {
				of_imsg('delete tag('+tag_id+')-role('
					+role_id+')-user('+user_id+
					') successfully', 1000,
					function(){ refresh_portion(url); });
			}
		});
	}, function(){});
}

function unbind_tag_role_token(tag_id, role_id, token_id, url) {
	layer.confirm("真的要删除么？", {
		icon: 3,
		btn: [ '确定', '取消' ]
	}, function() {
		$.delete('/v1.0/rel/tag/role/token', JSON.stringify({
				"tag_id": parseInt(tag_id),
				"role_id": parseInt(role_id),
				"token_id": parseInt(token_id)
		}), function(json) {
			if (json.code != 200) {
				of_emsg(json.msg, 2000);
			} else {
				of_imsg('delete tag('+tag_id+')-role('
					+role_id+')-token('+token_id+
					') successfully', 1000,
					function(){ refresh_portion(url); });
			}
		});
	}, function(){});
}

function edit_config(module){
	$.post('/settings/config/'+module, JSON.stringify(
		$('#form-config').serializeObject()
	), function(json) {
		if (json.code != 200) {
			of_emsg(json.msg, 2000);
		} else {
			of_imsg("update sysconfig success", 1000);
		}
	});
}

function debug_action(action) {
	$.get('/settings/debug/'+action, {}, function(json) {
		if (json.code != 200) {
			of_emsg(json.msg);
		} else {
			of_imsg(json.data);
		}
	});
}

function refresh_portion(url) {
	curUrl = updateURLParameter(url,'portion', '1')
	$.get(curUrl, {}, function(content){
		$("#content").html(content);
	});
}

$("#search").submit(function(event){
	url = updateURLParameter(curUrl, $("#query").attr("name"),
		$("#query").val());
	url = updateURLParameter(url, "p", "0");
    	if (curTagId) {
		url = updateURLParameter(url, "tag_id", curTagId);
	}
	refresh_portion(url);
	event.preventDefault();
});

function add_param_get(url,name,value){
	url += (url.indexOf("?") == -1)? "?":"&";
	url += name + "=" + value;
	return url;
}


function updateURLParameter(url, param, paramVal){
    var newAdditionalURL = "";
    var tempArray = url.split("?");
    var baseURL = tempArray[0];
    var additionalURL = tempArray[1];
    var temp = "";
    if (additionalURL) {
        tempArray = additionalURL.split("&");
        for (i=0; i<tempArray.length; i++){
            if(tempArray[i].split('=')[0] != param){
                newAdditionalURL += temp + tempArray[i];
                temp = "&";
            }
        }
    }

    var rows_txt = temp + "" + param + "=" + paramVal;
    return baseURL + "?" + newAdditionalURL + rows_txt;
}


function generate_select2(selector, placeholder, multiple, url, init){
    $(selector).select2({
        placeholder: placeholder,
        multiple: multiple,
        ajax: {
            url: url,
            dataType: 'json',
            quietMillis: 250,
            data: function(term, page) {
                return { query: term, per: 20 };
            },
            results: function(json) {
                if (json.code == 200) {
                    return {results: json.data};
                }
            },
            cache: true
        },
	initSelection: init,
        allowClear: true,
        minimumInputLength: 2,
        id: function(obj){return obj.id;},
        formatResult: function(obj) {return obj.name;},
        formatSelection: function(obj) {return obj.name;}
    });

}

jQuery.each(["put", "delete"], function(i, method) {
	jQuery[method] = function(url, data, callback, type) {
		if (jQuery.isFunction(data)) {
			type = type || callback;
			callback = data;
			data = undefined;
		}

		return jQuery.ajax({
			url: url,
			type: method,
			dataType: type,
			data: data,
			success: callback
		});
	};
});

$.fn.serializeObject = function() {
	var o = {};
	var a = this.serializeArray();
	$.each(a, function() {
		if (o[this.name] !== undefined) {
			if (!o[this.name].push) {
				o[this.name] = [o[this.name]];
			}
			o[this.name].push(this.value || '');
		} else {
			o[this.name] = this.value || '';
		}
	});
	return o;
};

function OnRightClick(event, treeId, treeNode) {
	if (treeNode) {
		zTree.selectNode(treeNode);
		showRMenu(event.clientX, event.clientY);
	}
}
function tag_filter(treeId, parentNode, childNodes) {
	if (!childNodes) return null;
	for (var i=0, l=childNodes.length; i<l; i++) {
		childNodes[i].name = last_tag(childNodes[i].name);
		childNodes[i].target = "_self";
		childNodes[i].url = "#";
		childNodes[i].click = "click_tree("+childNodes[i].id+");";
		childNodes[i].isParent = true;
	}
	return childNodes;
}

function click_tree(tag_id) {
	curTagId = tag_id;
	refresh_portion(updateURLParameter(
		curUrl,'tag_id', tag_id));
}

function showRMenu(x, y) {
	$("#rMenu ul").show();
	rMenu.css({"top":y+"px", "left":x+"px", "visibility":"visible"});
	$("body").bind("mousedown", onBodyMouseDown);
}

function hideRMenu() {
	if (rMenu) rMenu.css({"visibility": "hidden"});
	$("body").unbind("mousedown", onBodyMouseDown);
}

function onBodyMouseDown(event){
	if (!(event.target.id == "rMenu" || $(event.target).parents("#rMenu").length>0)) {
		rMenu.css({"visibility" : "hidden"});
	}
}

function last_tag(tag){
	var pos=tag.lastIndexOf(",");
	return tag.substring(pos+1); 
}

var zTree, rMenu, curTagId, curUrl;
{{if .zTree}}
var setting = {
	async: {
		enable: true,
		url:"/v1.0/rel/zTreeNodes",
		autoParam:["id", "name=n", "level=lv"],
		dataFilter: tag_filter
	},
	callback: {
		onRightClick: OnRightClick
	}
};
$(document).ready(function(){
	$.fn.zTree.init($("#tree"), setting);
	zTree = $.fn.zTree.getZTreeObj("tree");
	rMenu = $("#rMenu");
	curUrl = {{.URL}};
});
{{else}}
$(document).ready(function(){
	curUrl = {{.URL}};
});
{{end}}

    </script>
  </body>
</html>
