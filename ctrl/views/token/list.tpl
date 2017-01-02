{{template "layout/head.tpl" .}}
<div class="container-fluid">
  <div class="row">
    {{template "layout/sidebar.tpl" .}}
    <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
      <div class="page-header">
        <h1>Token</h1>
        <div class="table-responsive">
          <a href="/token/add" class="btn btn-default pull-right"><span class="glyphicon glyphicon-plus"></span>Add</a>
          <table class="table table-striped">
            <thead>
              <tr>
                <th>name</th>
                <th>cname</th>
                <th>note</th>
                <th>created</th>
                <th class="pull-right">command</th>
              </tr>
            </thead>
            <tbody> {{range .Tokens}}
              <tr>
                <td>{{.Name}}</td>
                <td>{{.Cname}}</td>
                <td>{{.Note}}</td>
                <td>{{dateformat .Create_time "2006-01-02 15:04:05"}}</td>
                <td>
                  <div class="pull-right">     
                    <a href="/token/edit/{{.Id}}" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-edit"></span> </a>                 
                    <span class="cut-line">¦</span>
                    <a href="javascript:delete_meta('token','{{.Id}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-trash"></span> </a>
                  </div> 
                </td>
              </tr> {{end}}
            </tbody>
          </table>
        </div>
        {{template "layout/paginator.tpl" .}}  
      </div>
    </div>
  </div>
</div>
{{template "layout/foot.tpl" .}}

