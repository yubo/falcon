{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">{{.H1}}</h1>
<div class="container">
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
</div>
{{template "layout/sidebar_foot.tpl" .}}
