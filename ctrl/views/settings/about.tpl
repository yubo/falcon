{{template "layout/head.tpl" .}}
    <div class="container-fluid">
      <div class="row">
        {{template "layout/sidebar.tpl" .}}
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
          <h1 class="page-header">{{.H1}}</h1>
          <div class="container">
{{if .Me}}
            <table class="table table-striped">
              <tbody>
                <tr><td>id         </td><td>{{.Me.Id}}    </td></tr>
                <tr><td>uuid       </td><td>{{.Me.Uuid}}  </td></tr>
                <tr><td>用户名     </td><td>{{.Me.Name}}  </td></tr>
                <tr><td>中文名     </td><td>{{.Me.Cname}} </td></tr>
                <tr><td>email      </td><td>{{.Me.Email}} </td></tr>
                <tr><td>phone      </td><td>{{.Me.Phone}} </td></tr>
                <tr><td>im         </td><td>{{.Me.IM}}    </td></tr>
                <tr><td>qq         </td><td>{{.Me.QQ}}    </td></tr>
              </tbody>
            </table>
{{else}}
not login
{{end}}
          </div>
        </div>
      </div>
    </div>
{{template "layout/foot.tpl" .}}
