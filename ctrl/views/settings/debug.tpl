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
                <tr>
		  <td>populate demo data </td>
		  <td><button type="button" class="btn btn-default" id="debug-action" onclick="debug_action('populate');"> <span class="glyphicon glyphicon-play"></span> 开始 </button> </td>

		</tr>
                <tr>
		  <td>reset database </td>
		  <td><button type="button" class="btn btn-default" id="debug-action" onclick="debug_action('reset_db');"> <span class="glyphicon glyphicon-play"></span> 开始 </button> </td>
		</tr>
                <tr>
		  <td>test msg </td>
		  <td><button type="button" class="btn btn-default" id="debug-action" onclick="debug_action('msg');"> <span class="glyphicon glyphicon-play"></span> 开始 </button> </td>
		</tr>
              </tbody>
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
