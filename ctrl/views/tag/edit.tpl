{{template "layout/head.tpl" .}}
<div class="container">
  <div class="page-header">
    <h1 class="page-header">{{.H1}} {{.Tag.Id}}</h1>
    <form>
      <div class="form-group">
        <label for="name">Tag Name</label>
        <input type="text" name="name" id="name" class="form-control" value="{{.Tag.Name}}" />
      </div>

      <button type="button" class="btn btn-default" id="update_tag" onclick="edit_tag('{{.Tag.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    </from>
  </div>
</div>
{{template "layout/foot.tpl" .}}
