{{template "layout/head.tpl" .}}
<div class="container-fluid">
  <div class="row">
    {{template "layout/sidebar.tpl" .}}
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <div class="page-header">
        <h1 class="page-header">{{.H1}} {{.Host.Id}}</h1>
        <form id="form-meta">
          <div class="form-group">
            <label for="name">机器名</label>
            <input type="text" name="name" id="name" class="form-control" value="{{.Host.Name}}" />
          </div>

          <div class="form-group">
            <label for="uuid">uuid</label>
            <input type="text" name="uuid" id="uuid" class="form-control" value="{{.Host.Uuid}}" />
          </div>

          <div class="form-group">
            <label for="type">type</label>
            <input type="text" name="type" id="type" class="form-control" value="{{.Host.Type}}" />
          </div>

          <div class="form-group">
            <label for="status">status</label>
            <input type="text" name="status" id="status" class="form-control" value="{{.User.Status}}" />
          </div>

          <div class="form-group">
            <label for="loc">loc</label>
            <input type="text" name="loc" id="loc" class="form-control" value="{{.User.Loc}}" />
          </div>

          <div class="form-group">
            <label for="idc">idc</label>
            <input type="text" name="idc" id="idc" class="form-control" value="{{.User.Idc}}" />
          </div>

          <button type="button" class="btn btn-default" id="edit-meta" onclick="edit_meta('host','{{.Host.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
          <a href="/host" class="btn btn-default"> <span class="glyphicon glyphicon-arrow-left"></span>返回</a>
        </from>
      </div>
    </div>
  </div>
</div>
{{template "layout/foot.tpl" .}}
