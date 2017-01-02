<div class="col-sm-3 col-md-2 sidebar">
  <ul class="nav nav-sidebar">
{{range .Links}}
{{if eq .Text $.CurLink}}
    <li class="active"><a href="#">{{.Text}}<span class="sr-only">(current)</span></a></li>
{{else if .Text}}
    <li{{if .Disabled}} class="disabled"{{end}}><a href="{{.Url}}">{{.Text}}</a> </li>
{{else}}
  </ul>
  <ul class="nav nav-sidebar">
{{end}}
{{end}}
  </ul>
</div>

