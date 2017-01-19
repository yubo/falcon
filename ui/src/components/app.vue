<template>
<div id="app">
  <navbar></navbar>
  <router-view></router-view>
</div>
</template>

<script>
import navbar from './navbar'
import sideBar from './sidebar'
import container from './container'

export default {
  components: {
    navbar,
    sideBar,
    container
  },
  data () {
    return { }
  }

}
</script>

<style>
</style>
