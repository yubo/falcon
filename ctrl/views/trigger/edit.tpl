{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">{{.H1}} {{.Role.Id}}</h1>
<div class="container">
  <form id="form-meta">
    <div class="form-group">
      <label for="name">Role Name</label>
      <input type="text" name="name" id="name" class="form-control" value="{{.Role.Name}}" />
    </div>
  
    <div class="form-group">
      <label for="cname">Role Cname</label>
      <input type="text" name="cname" id="cname" class="form-control" value="{{.Role.Cname}}" />
    </div>
  
    <div class="form-group">
      <label for="note">Role Note</label>
      <input type="text" name="note" id="note" class="form-control" value="{{.Role.Note}}" />
    </div>
  
    <button type="button" class="btn btn-default" id="edit-meta" onclick="edit_meta('role','{{.Role.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    <a href="#" onClick="refresh_portion('/role');" class="btn btn-default"> <span class="glyphicon glyphicon-arrow-left"></span>返回</a>
  </from>
</div>
{{template "layout/sidebar_foot.tpl" .}}
