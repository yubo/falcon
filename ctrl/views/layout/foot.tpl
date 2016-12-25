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

function update_profile() {
	$.post('/me/profile', {
		'cnname' : $("#cnname").val(),
		'email' : $("#email").val(),
		'phone' : $("#phone").val(),
		'im' : $("#im").val(),
		'qq' : $("#qq").val()
	}, function(json) {
		if (json.msg.length > 0) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly("更新成功：）");
		}
	});
}

function change_password() {
	$.post('/me/chpwd', {
		'old_password' : $("#old_password").val(),
		'new_password' : $("#new_password").val(),
		'repeat_password' : $("#repeat_password").val()
	}, function(json) {
		if (json.msg.length > 0) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly("密码修改成功：）");
		}
	});
}

function register() {
	$.post('/auth/register', {
		'name' : $('#name').val(),
		'password' : $("#password").val(),
		'repeat_password' : $("#repeat_password").val()
	}, function(json) {
		if (json.msg.length > 0) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly('sign up successfully', function() {
				location.href = '/auth/login';
			});
		}
	});
}

function login() {
	var raw = $('#ldap').prop('checked');
	if (raw) {
		useLdap = '1'
	} else {
		useLdap = '0'
	}
	$.post('/auth/login', {
		'name' : $('#name').val(),
		'password' : $("#password").val(),
		'ldap' : useLdap,
		'sig': $("#sig").val(),
		'callback': $("#callback").val()
	}, function(json) {
		if (json.msg.length > 0) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly('sign in successfully', function() {
				var redirect_url = '/me/info';
				if (json.data.length > 0) {
					redirect_url = json.data;
				}
				location.href = redirect_url;
			});
		}
	});
}



function create_team() {
	$.post('/me/team/c', {
		'name' : $("#name").val(),
		'resume' : $("#resume").val(),
		'users' : $("#users").val()
	}, function(json) {
		if (json.msg.length > 0) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly('create team successfully');
		}
	});
}

function edit_team(tid) {
	$.post('/target-team/edit', {
		'resume' : $("#resume").val(),
		'users' : $("#users").val(),
		'id': tid
	}, function(json) {
		if (json.msg.length > 0) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly('edit team successfully');
		}
	});
}



function delete_team(id) {
	my_confirm("真的真的要删除么？", [ '确定', '取消' ], function() {
		$.get('/target-team/delete?id='+id, {}, function(json) {
			if (json.msg.length > 0) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete team successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
	});
}

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
		"cnname" : $("#cnname").val(),
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
	jQuery[method]('/v1.0/host/'+id, JSON.stringify({
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

function edit_scope(id, method) {
	jQuery[method]('/v1.0/scope/'+id, JSON.stringify({
		"name" : $("#name").val(),
		"system_id" : parseInt($("#system_id").val()),
		"cname" : $("#cname").val(),
		"note" : $("#note").val()
	}), function(json) {
		if (json.code != 200) {
			err_message_quietly(json.msg);
		} else {
			ok_message_quietly(method +" scope success");
		}
	});
}

function delete_scope(id) {
	my_confirm("真的要删除么？", [ '确定', '取消' ], function() {
		$.delete('/v1.0/scope/'+id, {}, function(json) {
			if (json.code != 200) {
				err_message_quietly(json.msg);
			} else {
				ok_message_quietly('delete scope successfully', function() {
					location.reload();
				});
			}
		});
	}, function() {
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
    </script>
  </body>
</html>
