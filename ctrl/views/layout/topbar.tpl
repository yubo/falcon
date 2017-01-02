<ul class="nav nav-pills">
{{range .Links}}
{{if eq .Text $.CurLink}}
    <li class="active"><a href="#">{{.Text}}<span class="sr-only">(current)</span></a></li>
{{else}}
    <li{{if .Disabled}} class="disabled"{{end}}><a href="{{.Url}}">{{.Text}}</a> </li>
{{end}}
{{end}}
</ul>

