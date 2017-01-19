{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">{{.Module}}</h1>
<p>{{.Note}}</p>
<div class="container">
  <form id="form-config">
  {{range $k, $v := .Config}}
    <div class="form-group">
      <label for="name">{{$v.Key}}({{$v.Note}})</label>
      <input type="text" name="{{$v.Key}}" id="{{$v.Key}}" class="form-control" value="{{$v.Value | configFilter}}" />
    </div>
  {{end}}
  <button type="button" class="btn btn-default" id="edit-config" onclick="edit_config('{{.Module}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
  </from>
</div>
{{template "layout/sidebar_foot.tpl" .}}
