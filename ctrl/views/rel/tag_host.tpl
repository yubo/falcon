{{template "layout/head.tpl" .}}
    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-3 col-md-2 sidebar">

        </div>
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
          {{template "layout/topbar.tpl" .}}
          <br />
          <p class="leader">在节点上设置机器(host)</p>
          <div class="container">
            {{.CurLink}}
          </div>
        </div>
      </div>
    </div>
{{template "layout/foot.tpl" .}}
