{{template "layout/head.tpl" .}}
  <div class="container">
    <div class="page-header">
      <h1>System</h1>
      <div class="table-responsive">
        <table class="table table-striped">
          <thead>
            <tr>
              <th>name</th>
              <th>cname</th>
              <th>developers</th>
              <th>email</th>
              <th>created</th>
              <th>scope cmd</th>
              <th class="pull-right">system cmd</th>
            </tr>
          </thead>
          <tbody> {{range .Systems}}
            <tr>
              <td>{{.Name}}</td>
              <td>{{.Cname}}</td>
              <td>{{.Developers}}</td>
              <td>{{.Email}}</td>
              <td>{{dateformat .Create_time "2006-01-02 15:04:05"}}</td>
              <td>
                  <a href="/scope/add/{{.Id}}" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-plus"></span> </a>
                  <span class="cut-line">¦</span>
                  <a href="/scope/{{.Id}}" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-th-list"></span> </a>
              </td>
               <td>
                <div class="pull-right">     
                  <a href="/system/edit/{{.Id}}" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-edit"></span> </a>
                  <span class="cut-line">¦</span>
                  <a href="javascript:delete_system('{{.Id}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-trash"></span> </a>
                </div> 
              </td>
            </tr> {{end}}
          </tbody>
        </table>
      </div>
      {{template "layout/paginator.tpl" .}}  
    </div>
  </div>
{{template "layout/foot.tpl" .}}

