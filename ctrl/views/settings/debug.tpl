{{template "layout/sidebar_head.tpl" .}}
<h1 class="page-header">{{.H1}}</h1>
<div class="container">
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
</div>
{{template "layout/sidebar_foot.tpl" .}}
