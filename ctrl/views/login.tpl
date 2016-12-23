{{template "layout/head.tpl" .}}

    <div class="container">
{{range .Modules}}
        <h2 class="form-signin-heading">{{str2html .Name}}</h2>
	{{str2html .Html}}
{{end}}
    </div> <!-- /container -->


{{template "layout/foot.tpl" .}}
