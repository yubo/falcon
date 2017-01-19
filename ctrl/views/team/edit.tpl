{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">{{.H1}} {{.TeamUsers.Team.Id}}</h1>
<div class="container">
  <form id="form-meta">
    <div class="form-group">
      <label for="name">Team Name</label>
      <input type="text" name="name" id="name" class="form-control" value="{{.TeamUsers.Team.Name}}" />
    </div>
  
    <div class="form-group">
      <label for="note">Team Note</label>
      <input type="text" name="note" id="note" class="form-control" value="{{.TeamUsers.Team.Note}}" />
    </div>
  
    <div class="form-group">
      <label for="users">members</label>
      <input type="text" id="users" class="form-control" value="{{.User_ids}}"/>
    </div>

    <button type="button" class="btn btn-default" id="edit-meta" onclick="edit_team_users('{{.TeamUsers.Team.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    <a href="#" onClick="refresh_portion('/team');" class="btn btn-default"> <span class="glyphicon glyphicon-arrow-left"></span>返回</a>
  </from>
</div>


<script type="text/javascript">
$(document).ready(function(){
{{if and .TeamUsers .TeamUsers.Users}}
        var init = function(element, callback) {
            var users = [];
{{range .TeamUsers.Users}}
            users.push({id:{{.Id}}, name:{{.Name}}});
{{end}}
            callback(users);
        };
	generate_select2("#users", "输入要绑定用户的名字", true, "/v1.0/user/search", init);
{{else}}
	generate_select2("#users", "输入要绑定用户的名字", true, "/v1.0/user/search");
{{end}}
});

</script>


{{template "layout/sidebar_foot.tpl" .}}
