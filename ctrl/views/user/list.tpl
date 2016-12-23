{{template "layout/head.tpl" .}}
  <div class="container">
    <div class="page-header">
      <h1>User</h1>
      <div class="table-responsive">
        <table class="table table-striped">
          <thead>
            <tr>
              <th>name</th>
              <th>uuid</th>
              <th>cnname</th>
              <th>email</th>
              <th>phone</th>
              <th>im</th>
              <th>qq</th>
              <th>created</th>
              <th class="pull-right">command</th>
            </tr>
          </thead>
          <tbody> {{range .Users}}
            <tr>
              <td>{{.Name}}</td>
              <td>{{.Uuid}}</td>
              <td>{{.Cname}}</td>
              <td>{{.Email}}</td>
              <td>{{.Phone}}</td>
              <td>{{.IM}}</td>
              <td>{{.QQ}}</td>
              <td>{{dateformat .Create_time "2006-01-02 15:04:05"}}</td>
              <td>
                <div class="pull-right">     
                  <a href="/user/edit/{{.Id}}" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-edit"></span> </a>                 
                  <span class="cut-line">¦</span>
                  <a href="javascript:delete_user('{{.Id}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-trash"></span> </a>
                </div> 
              </td>
            </tr>{{end}}
          </tbody>
        </table>
      </div>
      {{template "layout/paginator.tpl" .}}  
    </div>
  </div>
{{template "layout/foot.tpl" .}}
