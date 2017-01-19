{{if not .Portion}}
{{template "layout/tree_head.tpl" .}}
<br />
<p class="page-header">在节点上设置机器(host)</p>
<form id="rel" class="navbar-form navbar-left form-inline">
  <div class="input-group">
    <span class="input-group-addon">hosts</span> 
    <input type="text" id="hosts" class="js-states form-control" style="width:300px;"/>
  </div>
  <button type="button" class="btn btn-primary" onclick="bind_tag_hosts();">Bind</button>
</form>
<div id="content">
{{end}}

{{if .Hosts}}
  <table class="table table-striped">
    <thead>
      <tr>
        <th>name</th>
        <th>uuid</th>
        <th>type</th>
        <th>status</th>
        <th>loc</th>
        <th>idc</th>
        <th>created</th>
        <th class="pull-right">command</th>
      </tr>
    </thead>
    <tbody> {{range .Hosts}}
      <tr>
        <td>{{.Name}}</td>
        <td>{{.Uuid}}</td>
        <td>{{.Type}}</td>
        <td>{{.Status}}</td>
        <td>{{.Loc}}</td>
        <td>{{.Idc}}</td>
        <td>{{dateformat .Create_time "2006-01-02 15:04:05"}}</td>
        <td>
          <div class="pull-right">     
            <a href="#" onClick="unbind_tag_host('{{$.TagId}}', '{{.Id}}','{{$.URL}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-remove"></span> </a>
          </div> 
        </td>
      </tr> {{end}}
    </tbody>
  </table>
{{end}}


{{if not .Portion}}
</div>

<script type="text/javascript">
$(document).ready(function(){
	generate_select2("#hosts", "输入要绑定机器的名字", true, "/v1.0/host/search");
});
</script>
{{template "layout/tree_foot.tpl" .}}
{{end}}
