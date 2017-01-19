{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">{{.H1}} {{.Tag.Id}}</h1>
<div class="container">
  <form id="form-meta">
    <div class="form-group">
      <label for="name">Tag Name</label>
      <input type="text" name="name" id="name" class="form-control" value="{{.Tag.Name}}" />
    </div>
  
    <button type="button" class="btn btn-default" id="edit-meta" onclick="edit_meta('tag','{{.Tag.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    <a href="#" onClick="refresh_portion('/tag');" class="btn btn-default"> <span class="glyphicon glyphicon-arrow-left"></span>返回</a>
  </from>
</div>
{{template "layout/sidebar_foot.tpl" .}}
