{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">{{.H1}} {{.Rule.Id}}</h1>
<div class="container">
  <form id="form-meta">
    <div class="form-group">
      <label for="name">Template Name</label>
      <input type="text" name="name" id="name" class="form-control" value="{{.Rule.Name}}" />
    </div>
  
    <div class="form-group">
      <label for="pid">Parent Name</label>
      <input type="text" name="pid" id="pid" class="form-control" value="{{.Rule.Pid}}" />
    </div>
  
    <div class="form-group">
      <label for="sendto">Send to</label>
      <input type="text" name="sendto" id="sendto" class="form-control" value="{{.Rule.SendTo}}" />
    </div>

    <div class="form-group">
      <label for="callback">callback</label>
      <input type="text" name="url" id="url" class="form-control" value="{{.Rule.Url}}" />
      <div class="mt10">
          <label class="checkbox-inline"> <input type="checkbox" id="before_callback_sms" > 回调之前发提醒短信 </label>
          <label class="checkbox-inline"> <input type="checkbox" id="before_callback_mail" > 回调之前发提醒邮件 </label>
          <label class="checkbox-inline"> <input type="checkbox" id="after_callback_sms" > 回调之后发结果短信 </label>
          <label class="checkbox-inline"> <input type="checkbox" id="after_callback_mail" > 回调之后发结果邮件 </label>
      </div>


    </div>

    <button type="button" class="btn btn-default" id="edit-meta" onclick="edit_meta('rule','{{.Rule.Id}}','{{.Method}}');"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
    <a href="#" onClick="refresh_portion('/rule');" class="btn btn-default"> <span class="glyphicon glyphicon-arrow-left"></span>返回</a>
  </from>
</div>
{{template "layout/sidebar_foot.tpl" .}}
