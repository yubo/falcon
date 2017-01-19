<template>
<div class="subbar">
  <ul class="nav nav-pills">
    <li><a href="/settings/config/ctrl">ctrl</a> </li>
    <li><a href="/settings/config/agent">agent</a> </li>
    <li><a href="/settings/config/graph">graph</a> </li>
    <li><a href="/settings/profile">Profile</a> </li>
    <li class="active"><a href="#">About Me<span class="sr-only">(current)</span></a></li>
    <li><a href="/settings/debug">Debug</a> </li>
  </ul>
</div>
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
.subbar{
	margin-bottom: 20px;
}
</style>
