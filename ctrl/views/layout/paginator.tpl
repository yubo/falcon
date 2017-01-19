{{if gt .paginator.PageNums 1}}
<ul class="pagination mt0 mb20">
    {{if .paginator.HasPrev}}
    <li><a href="#" onClick="refresh_portion('{{.paginator.PageLinkFirst}}');">首页</a></li>
    <li><a href="#" onClick="refresh_portion('{{.paginator.PageLinkPrev}}');">&lt;</a></li>
    {{else}}
    <li class="disabled"><a>首页</a></li>
    <li class="disabled"><a>&lt;</a></li>
    {{end}}
    {{range $index, $page := .paginator.Pages}}
    <li{{if $.paginator.IsActive .}} class="active"{{end}}>
      <a href="#" onClick="refresh_portion('{{$.paginator.PageLink $page}}');">{{$page}}</a>
    </li>
    {{end}}
    {{if .paginator.HasNext}}
    <li><a href="#" onClick="refresh_portion('{{.paginator.PageLinkNext}}');">&gt;</a></li>
    <li><a href="#" onClick="refresh_portion('{{.paginator.PageLinkLast}}');">尾页</a></li>
    {{else}}
    <li class="disabled"><a>&gt;</a></li>
    <li class="disabled"><a>尾页</a></li>
    {{end}}
</ul>
{{end}}
