{{template "layout/head.tpl" .}}
<div class="container">
  <div class="page-header">
    <h1 class="page-header">{{.H1}} {{.Role.Id}}</h1>
    <form>
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

      <button type="button" class="btn btn-default" id="edit-role" onclick="edit_role('{{.Role.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    </from>
  </div>
</div>
{{template "layout/foot.tpl" .}}
