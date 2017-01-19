{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">Tag</h1>
<div class="table-responsive">
  <a href="#" onClick="refresh_portion('/tag/add');" class="btn btn-default pull-right"><span class="glyphicon glyphicon-plus"></span>Add</a>
  <table class="table table-striped">
    <thead>
      <tr>
        <th>name</th>
        <th>created</th>
        <th class="pull-right">command</th>
      </tr>
    </thead>
    <tbody> {{range .Tags}}
      <tr>
        <td>{{.Name}}</td>
        <td>{{dateformat .Create_time "2006-01-02 15:04:05"}}</td>
        <td>
          <div class="pull-right">     
            <a href="#" onClick="refresh_portion('/tag/edit/{{.Id}}');" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-edit"></span> </a>                 
            <span class="cut-line">¦</span>
            <a href="#" onClick="delete_meta('tag','{{.Id}}','{{$.URL}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-trash"></span> </a>
          </div> 
        </td>
      </tr> {{end}}
    </tbody>
  </table>
</div>
{{template "layout/paginator.tpl" .}}  
{{template "layout/sidebar_foot.tpl" .}}
