    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="/static/js/jquery.min.js"></script>
    <script src="/static/js/bootstrap.min.js"></script>
    <!-- Just to make our placeholder images work. Don't actually copy the next line! -->
    <script src="/static/js/holder.min.js"></script>
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <script src="/static/js/ie10-viewport-bug-workaround.js"></script>
    <script src="/static/layer/layer.min.js"></script>
    <script src="/static/js/custom.js"></script>
    <script type="text/javascript">

function err_message_quietly(msg, f) {
	$.layer({
		title : false,
		closeBtn : false,
		time : 2,
		dialog : {
			msg : msg
		},
		end : f
	});
}

function ok_message_quietly(msg, f) {
	$.layer({
		title : false,
		closeBtn : false,
		time : 1,
		dialog : {
			msg : msg,
			type : 1
		},
		end : f
	});
}

function my_confirm(msg, btns, yes_func, no_func) {
	$.layer({
		shade : [ 0 ],
		area : [ 'auto', 'auto' ],
		dialog : {
			msg : msg,
			btns : 2,
			type : 4,
			btn : btns,
			yes : yes_func,
			no : no_func
		}
	});
}

// - business function -

function set_role(uid, obj) {
	var role = obj.checked ? '1' : '0';
	$.post('/target-user/role?id='+uid, {
		'role' : role
	}, function(json) {
		if (json.msg.length > 0) {
			err_message_quietly(json.msg);
			location.reload();
		} else {
			if (role == '1') {
				ok_message_quietly('成功设置为管理员：）');
			} else if (role == '0') {
				ok_message_quietly('成功取消管理员权限：）');
			}
		}
	});
}

function user_detail(uid) {
	$("#user_detail_div").load("/user/detail?id=" + uid);
	$.layer({
		type : 1,
		shade : [ 0.5, '#000' ],
		shadeClose : true,
		closeBtn : [ 0, true ],
		area : [ '450px', '240px' ],
		title : false,
		border : [ 0 ],
		page : {
			dom : '#user_detail_div'
		}
	});
}

function edit_user(id, method) {
	jQuery[method]('/v1.0/user/'+id, JSON.stringify({
		"name" : $("#name").val(),
		"cname" : $("#cname").val(),
		"email" : $("#email").val(),
		"phone" : $("#phone").val(),
		"im" : $("#im").val(),
		"qq" : $("#qq").val()
	}), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly(method +" user success");
		}
	});
}

function delete_user(id) {
	my_confirm("真的要删除么？", [ '确定', '取消' ], function() {
		$.delete('/v1.0/user/'+id, {}, function(json) {
		        if (json.code != 200) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete user successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
	});
}

function edit_host(id, method) {
	jQuery[method]('/v1.0/host/'+id+'?tag='+$("#tag").val(), JSON.stringify({
		"Uuid" : $("#uuid").val(),
		"name" : $("#name").val(),
		"type" : $("#type").val(),
		"status" : $("#status").val(),
		"loc" : $("#loc").val(),
		"idc" : $("#idc").val()
	}), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly(method +" host success");
		}
	});
}

function delete_host(id) {
	my_confirm("真的要删除么？", [ '确定', '取消' ], function() {
		$.delete('/v1.0/host/'+id, {}, function(json) {
			if (json.code != 200) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete host successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
	});
}

function edit_tag(id, method) {
	jQuery[method]('/v1.0/tag/'+id, JSON.stringify({
		"name" : $("#name").val()
	}), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly(method +" tag success");
		}
	});
}

function delete_tag(id) {
	my_confirm("真的要删除么？", [ '确定', '取消' ], function() {
		$.delete('/v1.0/tag/'+id, {}, function(json) {
			if (json.code != 200) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete tag successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
	});
}

function edit_role(id, method) {
	jQuery[method]('/v1.0/role/'+id, JSON.stringify({
		"name" : $("#name").val(),
		"cname" : $("#cname").val(),
		"note" : $("#note").val()
	}), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly(method +" role success");
		}
	});
}

function delete_role(id) {
	my_confirm("真的要删除么？", [ '确定', '取消' ], function() {
		$.delete('/v1.0/role/'+id, {}, function(json) {
			if (json.code != 200) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete role successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
	});
}

function edit_system(id, method) {
	jQuery[method]('/v1.0/system/'+id, JSON.stringify({
		"name" : $("#name").val(),
		"cname" : $("#cname").val(),
		"developers" : $("#developers").val(),
		"email" : $("#email").val()
	}), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly(method +" system success");
		}
	});
}

function delete_system(id) {
	my_confirm("真的要删除么？", [ '确定', '取消' ], function() {
		$.delete('/v1.0/system/'+id, {}, function(json) {
			if (json.code != 200) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete system successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
	});
}

function edit_token(id, method) {
	jQuery[method]('/v1.0/token/'+id, JSON.stringify({
		"name" : $("#name").val(),
		"system_id" : parseInt($("#system_id").val()),
		"cname" : $("#cname").val(),
		"note" : $("#note").val()
	}), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly(method +" token success");
		}
	});
}

function delete_token(id) {
	my_confirm("真的要删除么？", [ '确定', '取消' ], function() {
		$.delete('/v1.0/token/'+id, {}, function(json) {
			if (json.code != 200) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete token successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
	});
}

function edit_config(module){
	$.post('/settings/config/'+module, JSON.stringify(
		$('#form-config').serializeObject()
	), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly("update sysconfig success");
		}
	});
}

jQuery.each( [ "put", "delete" ], function( i, method ) {
  jQuery[ method ] = function( url, data, callback, type ) {
    if ( jQuery.isFunction( data ) ) {
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
    </script>
  </body>
</html>
