{{template "layout/head.tpl" .}}
<div class="container">
  <div class="page-header">
    <h1 class="page-header">{{.H1}} {{.System.Name}}</h1>
    <form>
      <div class="form-group">
        <label for="name">Scope Name</label>
        <input type="text" name="name" id="name" class="form-control" value="{{.Scope.Name}}" />
      </div>
      <div class="form-group">
        <label for="cname">Scope Cname</label>
        <input type="text" name="cname" id="cname" class="form-control" value="{{.Scope.Cname}}" />
      </div>
      <div class="form-group">
        <label for="developers">Scope Note</label>
        <input type="text" name="developers" id="developers" class="form-control" value="{{.Scope.Note}}" />
      </div>

      <input type="hidden" name="system_id" id="system_id" class="form-control" value="{{.Scope.System_id}}" />
      <button type="button" class="btn btn-default" id="update_system" onclick="edit_scope('{{if gt .Scope.Id 0}}{{.Scope.Id}}{{end}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    </from>
  </div>
</div>
{{template "layout/foot.tpl" .}}
