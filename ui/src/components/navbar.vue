<template>
<nav class="navbar navbar-inverse navbar-fixed-top">
  <div class="container-fluid">
    <div class="navbar-header">
      <ol class="breadcrumb">
              <li><a href="#"> Falcon </a></li>
              <li><a href="#"> Settings </a></li>
              <li><a href="#"> About Me </a></li>
      </ol>
    </div>
    <div id="navbar" class="navbar-collapse collapse">
      <ul class="nav navbar-nav navbar-right">
        <li class="dropdown">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Dashboard<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li class="disabled"><a href="#">Falcon</a></li>
            <li class="disabled"><a href="#">Alarm</a></li>
            <li class="disabled"><a href="#">Graph</a></li>
          </ul>
        </li>
        <li class="dropdown">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Relation<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><a href="/rel/tag/host">Tag Host</a></li>
            <li><a href="/rel/tag/role/user">Tag Role User</a></li>
            <li><a href="/rel/tag/role/token">Tag Role Token</a></li>
            <li><a href="/rel/tag/rule/trigger">Tag Template Trigger</a></li>
          </ul>
        </li>
        <li class="dropdown">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Meta<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><a href="/tag">Tag</a></li>
            <li><a href="/host">Host</a></li>
            <li role="separator" class="divider"></li>
            <li><a href="/role">Role</a></li>
            <li><a href="/user">User</a></li>
            <li><a href="/token">Token</a></li>
            <li role="separator" class="divider"></li>
            <li><a href="/team">Team</a></li>
            <li><a href="/rule">Template</a></li>
            <li class="disabled"><a href="#">Trigger</a></li>
            <li><a href="/expression">Expression</a></li>
          </ul>
        </li>
        <li class="dropdown">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">Settings<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><a href="/settings/config/ctrl">ctrl</a></li>
            <li><a href="/settings/config/agent">agent</a></li>
            <li><a href="/settings/config/graph">graph</a></li>
            <li role="separator" class="divider"></li>
            <li><a href="/settings/profile">Profile</a></li>
            <li><a href="/settings/aboutme">About Me</a></li>
            <li role="separator" class="divider"></li>
            <li><a href="/settings/debug">Debug</a></li>
          </ul>
        </li>
        <li class="dropdown">
          <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true"
            aria-expanded="false">help<span class="caret"></span></a>
          <ul class="dropdown-menu">
            <li><a href="/doc">doc</a></li>
            <li role="separator" class="divider"></li>
            <li><a href="/about">About Falcon</a></li>
        </li>
        <li><a href="/logout">[logout]</a></li>
      </ul>
    </div>
  </div>
</nav>
</template>

<script>
export default {
  data () {
    return {
      msg: 'Welcome to falcon@xiaomi.com App 5'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}

.navbar-header > .breadcrumb {
  background-color: transparent;
  margin:0px 0px 0px 0px;
  padding: 15px 15px 15px 15px;
}
.navbar-header > .breadcrumb li a{
   font-size:14px;
   color: #9d9d9d;
   height:20px;
}

</style>
