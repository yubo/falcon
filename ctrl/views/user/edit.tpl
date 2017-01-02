{{template "layout/head.tpl" .}}
<div class="container-fluid">
  <div class="row">
    {{template "layout/sidebar.tpl" .}}
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <div class="page-header">
        <h1 class="page-header">{{.H1}} {{.User.Id}}</h1>
        <form id="form-meta">
          <div class="form-group">
            <label for="name">用户名</label>
            <input type="text" name="name" id="name" class="form-control" value="{{.User.Name}}" />
          </div>

          <div class="form-group">
            <label for="cname">中文名</label>
            <input type="text" name="cname" id="cname" class="form-control" value="{{.User.Cname}}" />
          </div>

          <div class="form-group">
            <label for="email">Email address</label>
            <input type="email" name="email" id="email" class="form-control" value="{{.User.Email}}" />
          </div>

          <div class="form-group">
            <label for="phone">Tel</label>
            <input type="text" name="phone" id="phone" class="form-control" value="{{.User.Phone}}" />
          </div>
          <div class="form-group">
            <label for="name">IM（内部通讯工具账号，比如百度hi、米聊）</label>
            <input type="text" name="im" id="im" class="form-control" value="{{.User.IM}}" />
          </div>
          <div class="form-group">
            <label for="name">QQ</label>
            <input type="text" name="qq" id="qq" class="form-control" value="{{.User.QQ}}" />
          </div>
          <button type="button" class="btn btn-default" id="edit-meta" onclick="edit_meta('user','{{.User.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
{{if .EditUser}}
          <a href="/user" class="btn btn-default"> <span class="glyphicon glyphicon-arrow-left"></span>返回</a>
{{end}}
        </from>
      </div>
    </div>
  </div>
</div>
{{template "layout/foot.tpl" .}}
