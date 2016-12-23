{{template "layout/head.tpl" .}}
<div class="container">
  <div class="page-header">
    <h1 class="page-header">{{.H1}} {{.System.Id}}</h1>
    <form>
      <div class="form-group">
        <label for="name">System Name</label>
        <input type="text" name="name" id="name" class="form-control" value="{{.System.Name}}" />
      </div>
      <div class="form-group">
        <label for="cname">System Cname</label>
        <input type="text" name="cname" id="cname" class="form-control" value="{{.System.Cname}}" />
      </div>
      <div class="form-group">
        <label for="developers">System Developers</label>
        <input type="text" name="developers" id="developers" class="form-control" value="{{.System.Developers}}" />
      </div>
      <div class="form-group">
        <label for="email">System Email</label>
        <input type="text" name="email" id="email" class="form-control" value="{{.System.Email}}" />
      </div>

      <button type="button" class="btn btn-default" id="update_system" onclick="edit_system('{{.System.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    </from>
  </div>
</div>
{{template "layout/foot.tpl" .}}
