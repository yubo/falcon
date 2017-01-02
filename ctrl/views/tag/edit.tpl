{{template "layout/head.tpl" .}}
<div class="container-fluid">
  <div class="row">
    {{template "layout/sidebar.tpl" .}}
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <div class="page-header">
        <h1 class="page-header">{{.H1}} {{.Tag.Id}}</h1>
        <form id="form-meta">
          <div class="form-group">
            <label for="name">Tag Name</label>
            <input type="text" name="name" id="name" class="form-control" value="{{.Tag.Name}}" />
          </div>

          <button type="button" class="btn btn-default" id="edit-meta" onclick="edit_meta('tag','{{.Tag.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
          <a href="/tag" class="btn btn-default"> <span class="glyphicon glyphicon-arrow-left"></span>返回</a>
        </from>
      </div>
    </div>
  </div>
</div>
{{template "layout/foot.tpl" .}}
