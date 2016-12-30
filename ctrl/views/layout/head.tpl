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
    <link href="/static/css/custom.css" rel="stylesheet">

    <!-- Just for debugging purposes. Don't actually copy these 2 lines! -->
    <!--[if lt IE 9]><script src="../../assets/js/ie8-responsive-file-warning.js"></script><![endif]-->
    <script src="/static/js/ie-emulation-modes-warning.js"></script>

    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
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
{{if .Me}}
            <li class="dropdown">
              <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">management<span class="caret"></span></a>
              <ul class="dropdown-menu">
                <li><a href="/tag">tag</a></li>
                <li><a href="/role">role</a></li>
                <li><a href="/user">user</a></li>
                <li><a href="/system">system</a></li>
                <li><a href="/host">host</a></li>
                <li class="disabled"><a href="/trigger">trigger</a></li>

                <li role="separator" class="divider"></li>
                <li><a href="/tag/add">add tag</a></li>
                <li><a href="/role/add">add role</a></li>
                <li><a href="/user/add">add user</a></li>
                <li><a href="/system/add">add system</a></li>
                <li><a href="/host/add">add host</a></li>
                <li class="disabled"><a href="/trigger/add">add trigger</a></li>

              </ul>
            </li>

            <li class="dropdown">
              <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">settings<span class="caret"></span></a>
              <ul class="dropdown-menu">
                <li><a href="/settings/config/global">Global</a></li>
                <li><a href="/settings/profile">Profile</a></li>
                <li><a href="/settings/aboutme">About Me</a></li>
              </ul>
            </li>

            <li class="dropdown">
              <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">help<span class="caret"></span></a>
              <ul class="dropdown-menu">
                <li><a href="/doc" target="_blank">Doc</a></li>
                <li role="separator" class="divider"></li>
                <li><a href="/about">About</a></li>
              </ul>
            </li>
            <li><a href="/logout">{{.Me.Name}}[logout]</a></li>
{{else}}
            <li class="dropdown">
              <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">help<span class="caret"></span></a>
              <ul class="dropdown-menu">
                <li><a href="/doc" target="_blank">Doc</a></li>
                <li role="separator" class="divider"></li>
                <li><a href="/about">About Falcon</a></li>
              </ul>
            </li>
	    <li><a href="/login">[login]</a></li>
{{end}}
          </ul>
{{if .Search}}
          <form class="navbar-form navbar-right" method="get" action="{{.Search.Url}}">
            <input name="{{.Search.Name}}" value="{{.Query}}" type="text" class="form-control" placeholder="Search...">
          </form>
{{end}}
        </div>
      </div>
    </nav>

