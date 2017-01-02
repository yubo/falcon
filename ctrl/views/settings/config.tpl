{{template "layout/head.tpl" .}}
    <div class="container-fluid">
      <div class="row">
        {{template "layout/sidebar.tpl" .}}
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
          <h1 class="page-header">{{.Moudle}} config</h1>
          <form id="form-config">
{{range $key, $value := .Config}}
            <div class="form-group">
              <label for="name">{{$key}}</label>
              <input type="text" name="{{$key}}" id="{{$key}}" class="form-control" value="{{$value}}" />
            </div>
{{end}}
            <button type="button" class="btn btn-default" id="edit-config" onclick="edit_config('{{.Moudle}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
          </from>

        </div>
      </div>
    </div>
{{template "layout/foot.tpl" .}}
