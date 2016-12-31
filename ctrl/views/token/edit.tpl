{{template "layout/head.tpl" .}}
<div class="container">
  <div class="page-header">
    <h1 class="page-header">{{.H1}} {{.System.Name}}</h1>
    <form>
      <div class="form-group">
        <label for="name">Token Name</label>
        <input type="text" name="name" id="name" class="form-control" value="{{.Token.Name}}" />
      </div>
      <div class="form-group">
        <label for="cname">Token Cname</label>
        <input type="text" name="cname" id="cname" class="form-control" value="{{.Token.Cname}}" />
      </div>
      <div class="form-group">
        <label for="developers">Token Note</label>
        <input type="text" name="developers" id="developers" class="form-control" value="{{.Token.Note}}" />
      </div>

      <input type="hidden" name="system_id" id="system_id" class="form-control" value="{{.Token.System_id}}" />
      <button type="button" class="btn btn-default" id="edit-system" onclick="edit_token('{{if gt .Token.Id 0}}{{.Token.Id}}{{end}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    </from>
  </div>
</div>
{{template "layout/foot.tpl" .}}
