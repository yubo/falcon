<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <!-- The above 3 meta tags *must* come first in the head; any other head content must come *after* these tags -->
    <meta name="description" content="">
    <meta name="author" content="">
    <link rel="icon" href="/static/favicon.ico">
    <title>Falcon</title>
    <!-- Bootstrap core CSS -->
    <link href="/static/css/bootstrap.min.css" rel="stylesheet">
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <link href="/static/css/ie10-viewport-bug-workaround.css" rel="stylesheet">
    <!-- Custom styles for this template -->
    <link href="/static/css/zTreeStyle/zTreeStyle.css" rel="stylesheet">
    <link href="/static/css/select2.css" rel="stylesheet">
    <link href="/static/css/select2-bootstrap.css" rel="stylesheet">
    <link href="/static/css/custom.css" rel="stylesheet">
    <!-- Just for debugging purposes. Don't actually copy these 2 lines! -->
    <!--[if lt IE 9]><script src="../../assets/js/ie8-responsive-file-warning.js"></script><![endif]-->
    <script src="/static/js/ie-emulation-modes-warning.js"></script>
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
    <script src="/static/js/jquery.min.js"></script>
    <script src="/static/js/bootstrap.min.js"></script>
    <script src="/static/js/jquery.ztree.core.min.js"></script>
    <script src="/static/js/jquery.ztree.exedit.min.js"></script>

    <script src="/static/js/select2.min.js"></script>
    <!-- IE10 viewport hack for Surface/desktop Windows 8 bug -->
    <script src="/static/js/ie10-viewport-bug-workaround.js"></script>
    <script src="/static/layer/layer.js"></script>
    <script src="/static/js/custom.js"></script>

  </head>
  <body>
    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container-fluid">
        <div class="navbar-header">
          <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
            <span class="sr-only">Toggle navigation</span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
            <span class="icon-bar"></span>
          </button>
          <a class="navbar-brand" href="#">Falcon</a>
        </div>
        <div id="navbar" class="navbar-collapse collapse">
          <ul class="nav navbar-nav navbar-right">

{{range .HeadLinks}}
  {{if .SubLinks}}
            <li class="dropdown">
              <a href="{{.Url}}" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">{{.Text}}<span class="caret"></span></a>
              <ul class="dropdown-menu">
    {{range .SubLinks}}
      {{if .Text}}
                <li{{if .Disabled}} class="disabled"{{end}}><a href="{{.Url}}">{{.Text}}</a></li>
      {{else}}
                <li role="separator" class="divider"></li>
      {{end}}
    {{end}}
              </ul>
            </li>
  {{else}} 
	    <li><a href="{{.Url}}">{{.Text}}</a></li>
  {{end}}
{{end}}
          </ul>
{{if .Search}}
          <form id="search" class="navbar-form navbar-right">
            <input id="query" name="{{.Search.Name}}" value="{{.Query}}" type="text" class="form-control" placeholder="{{.Search.Placeholder}}">
          </form>
{{end}}
        </div>
      </div>
    </nav>

