<template>
<div id="app">
  <navbar></navbar>
  <notification></notification>
  <router-view></router-view>
</div>
</template>

<script>
import navbar from './navbar'
import notification from './notification'

export default {
  components: {
    navbar,
    notification
  },
  data () {
    return { }
  }

}
</script>

<style>
</style>
