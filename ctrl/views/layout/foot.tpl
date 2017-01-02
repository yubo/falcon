    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="/static/js/jquery.min.js"></script>
    <script src="/static/js/bootstrap.min.js"></script>
    <!-- Just to make our placeholder images work. Don't actually copy the next line! -->
    <script src="/static/js/holder.min.js"></script>
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <script src="/static/js/ie10-viewport-bug-workaround.js"></script>
    <script src="/static/layer/layer.js"></script>
    <script src="/static/js/custom.js"></script>
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

function delete_meta(meta, id) {
	layer.confirm("真的要删除么？", {
		icon: 3,
		btn: [ '确定', '取消' ]
	}, function() {
		$.delete('/v1.0/'+meta+'/'+id, {}, function(json) {
			if (json.code != 200) {
				of_emsg(json.msg, 2000);
			} else {
				of_imsg('delete '+meta+' successfully', 1000,
					 function(){location.reload();});
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
    </script>
  </body>
</html>
