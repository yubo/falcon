{{if not .Portion}}
{{template "layout/tree_head.tpl" .}}
<br />
<p class="page-header">在节点(tag)上设置权限(token)和角色(role)之间的关系, 角色是权限和用户(user)在节点上的容器</p>
<form id="rel" class="navbar-form navbar-left">
  <div class="input-group">
    <span class="input-group-addon">role</span> 
    <input type="text" id="role" class="js-states form-control" style="width:197px;"/>
  </div>
  <div class="input-group">
    <span class="input-group-addon">tokens</span> 
    <input type="text" id="tokens" class="js-states form-control" style="width:300px;"/>
  </div>
  <button type="button" class="btn btn-primary" onclick="bind_role_tokens();">绑定</button>
</form>
<div id="content">
{{end}}

{{if .TagRoleToken}}
  <table class="table table-striped">
    <thead> <tr> <th>token</th> <th>role</th> <th>tag</th> <th class="pull-right">command</th> </tr> </thead>
    <tbody> {{range .TagRoleToken}}
      <tr> <td>{{.TokenName}}</td> <td>{{.RoleName}}</td> <td>{{.TagName}}</td> <td>
          <div class="pull-right">     
            <a href="#" onClick="unbind_tag_role_token('{{.TagId}}','{{.RoleId}}','{{.TokenId}}','{{$.URL}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-remove"></span> </a>
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
	$("#role").on("change", function(){
		$("#tokens").select2("val", "");
	});

	generate_select2("#role", "要绑定角色的名字", false, "/v1.0/role/search");
	generate_select2("#tokens", "要绑定权限的名字", true, "/v1.0/token/search");
});

</script>

{{template "layout/tree_foot.tpl" .}}
{{end}}
