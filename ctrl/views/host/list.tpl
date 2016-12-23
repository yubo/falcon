{{template "layout/head.tpl" .}}
  <div class="container">
    <div class="page-header">
      <h1>Host</h1>
      <div class="table-responsive">
        <table class="table table-striped">
          <thead>
            <tr>
              <th>name</th>
              <th>uuid</th>
              <th>type</th>
              <th>status</th>
              <th>loc</th>
              <th>idc</th>
              <th>created</th>
              <th class="pull-right">command</th>
            </tr>
          </thead>
          <tbody> {{range .Hosts}}
            <tr>
              <td>{{.Name}}</td>
              <td>{{.Uuid}}</td>
              <td>{{.Type}}</td>
              <td>{{.Status}}</td>
              <td>{{.Loc}}</td>
              <td>{{.Idc}}</td>
              <td>{{dateformat .Create_time "2006-01-02 15:04:05"}}</td>
              <td>
                <div class="pull-right">     
                  <a href="/host/edit/{{.Id}}" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-edit"></span> </a>                 
                  <span class="cut-line">¦</span>
                  <a href="javascript:delete_host('{{.Id}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-trash"></span> </a>
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

