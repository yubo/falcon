{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">Team</h1>
<div class="table-responsive">
  <a href="/teamusers/add" class="btn btn-default pull-right"><span class="glyphicon glyphicon-plus"></span>Add</a>
  <table class="table table-striped">
    <thead>
      <tr>
        <th>name</th>
        <th>note</th>
        <th class="pull-right">command</th>
      </tr>
    </thead>
    <tbody> {{range .Teams}}
      <tr>
        <td>{{.Name}}</td>
        <td>{{.Note}}</td>
        <td>
          <div class="pull-right">     
            <a href="/teamusers/edit/{{.Id}}" class="orange" style="text-decoration:none;"> <span class="glyphicon glyphicon-edit"></span> </a>                 
            <span class="cut-line">¦</span>
            <a href="#" onClick="delete_meta('team','{{.Id}}','{{$.URL}}');" class="orange" style="text-decoration:none;"><span class="glyphicon glyphicon-trash"></span> </a>
          </div> 
        </td>
      </tr> {{end}}
    </tbody>
  </table>
</div>
{{template "layout/paginator.tpl" .}}  
{{template "layout/sidebar_foot.tpl" .}}
