{{template "layout/head.tpl" .}}
    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-3 col-md-2 sidebar">
          <ul class="nav nav-sidebar">
{{range .Links}}
{{if eq .Text $.CurLink}}
            <li class="active"><a href="#">{{.Text}}<span class="sr-only">(current)</span></a></li>
{{else}}
            <li><a href="{{.Url}}">{{.Text}}</a> </li>
{{end}}
{{end}}
          </ul>
        </div>
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
{{template "user/_edit.tpl" .}}
        </div>
      </div>
    </div>
{{template "layout/foot.tpl" .}}
