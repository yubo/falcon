{{if not .Portion}}
{{template "layout/tree_head.tpl" .}}
<br />
<p class="page-header">在节点(tag)上设置触发器(trigger)和规则模板(template)之间的关系, 规则模板是触发器在节点上的容器</p>
<form id="rel" class="navbar-form navbar-left">
  <div class="input-group">
    <span class="input-group-addon">template</span> 
    <input type="text" id="rule" class="js-states form-control" style="width:197px;"/>
  </div>
  <div class="input-group">
    <span class="input-group-addon">triggers</span> 
    <input type="text" id="triggers" class="js-states form-control" style="width:300px;"/>
  </div>
  <button type="button" class="btn btn-primary" onclick="bind_rule_triggers();">Bind</button>
</form>

<div id="content">
{{end}}

{{template "rel/demo_table.tpl" .}}

{{if not .Portion}}
</div>

<script type="text/javascript">
$(document).ready(function(){
    $("#rule").on("change", function(){
        $("#triggers").select2("val", "");
    });

    $("#rule").select2({
        placeholder: "要绑定规则模板的名字",
        ajax: {
            url: "/v1.0/rule/search",
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
        allowClear: true,
        minimumInputLength: 2,
        id: function(obj){return obj.id;},
        formatResult: function(obj) {return obj.name;},
        formatSelection: function(obj) {return obj.name;}
    });
    $("#triggers").select2({
        placeholder: "要绑定触发器的名字",
        multiple: true,
        ajax: {
            url: "/v1.0/trigger/search",
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
        allowClear: true,
        minimumInputLength: 2,
        id: function(obj){return obj.id;},
        formatResult: function(obj) {return obj.name;},
        formatSelection: function(obj) {return obj.name;}
    });
});

</script>

{{template "layout/tree_foot.tpl" .}}
{{end}}

