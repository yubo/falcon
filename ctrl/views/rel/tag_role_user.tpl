{{if not .Portion}}
{{template "layout/tree_head.tpl" .}}
<br />
<p class="page-header">在节点(tag)上设置用户(user)和角色(role)之间的关系, 角色是权限(token)和用户在节点上的容器</p>
<form id="rel" class="navbar-form navbar-left">
  <div class="input-group">
    <span class="input-group-addon">role</span> 
    <input type="text" id="role" class="js-states form-control" style="width:197px;"/>
  </div>
  <div class="input-group">
    <span class="input-group-addon">users</span> 
    <input type="text" id="users" class="js-states form-control" style="width:300px;"/>
  </div>
  <button type="button" class="btn btn-primary" onclick="bind_role_users();">Bind</button>
</form>
<div id="content">
{{end}}

{{if .TagRoleUser}}
  <table class="table table-striped">
    <thead> <tr> <th>user</th> <th>role</th> <th>tag</th> <th class="pull-right">command</th> </tr> </thead>
    <tbody> {{range .TagRoleUser}}
      <tr> <td>{{.UserName}}</td> <td>{{.RoleName}}</td> <td>{{.TagName}}</td> <td>
          <div class="pull-right">     
            <a href="#" onClick="unbind_tag_role_user('{{.TagId}}','{{.RoleId}}','{{.UserId}}','{{$.URL}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-remove"></span> </a>
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
		$("#users").select2("val", "");
	});

	generate_select2("#role", "要绑定角色的名字", false, "/v1.0/role/search");
	generate_select2("#users", "要绑定用户的名字/email", true, "/v1.0/user/search");
});

</script>

{{template "layout/tree_foot.tpl" .}}
{{end}}
