{{template "layout/head.tpl" .}}
<div class="container-fluid">
  <div class="row">
    <div class="col-sm-3 col-md-2 sidebar">
      <div class="zTreeDemoBackground left">
      	<ul id="tree" class="ztree"></ul>
      </div>
    </div>
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <ul class="nav nav-pills">
      {{range .Links}} {{if eq .Text $.CurLink}}
        <li class="active"><a href="#">{{.Text}}<span class="sr-only">(current)</span></a></li>
      {{else}}
        <li{{if .Disabled}} class="disabled"{{end}}><a href="{{.Url}}">{{.Text}}</a> </li>
      {{end}} {{end}}
      </ul>
