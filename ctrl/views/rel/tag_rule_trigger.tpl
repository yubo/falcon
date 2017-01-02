{{template "layout/head.tpl" .}}
    <div class="container-fluid">
      <div class="row">
        <div class="col-sm-3 col-md-2 sidebar">

        </div>
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
          {{template "layout/topbar.tpl" .}}
          <br />
          <h1 class="page-header">{{.CurLink}}</h1>
          <p class="leader">在节点(tag)上设置触发器(trigger)和规则模板(rule)之间的关系, 规则模板是触发器在节点上的容器</p>
          <div class="container">
            {{.CurLink}}
          </div>
        </div>
      </div>
    </div>
{{template "layout/foot.tpl" .}}
