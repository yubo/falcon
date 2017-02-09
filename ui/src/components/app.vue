<template>
<div id="app">
  <navbar></navbar>
  <router-view></router-view>
</div>
</template>

<script>
import navbar from './navbar'
import { fetch } from 'src/utils'
import { Message } from 'element-ui'

export default {
  components: {
    navbar
  },
  data () {
    return { }
  },
  create () {
    if (!this.$store.config) {
      fetch({
        router: this.$router,
        method: 'get',
        url: 'settings/config/ui'
      }).then((res) => {
        this.$store.commit('m_set_config', res.data)
      }).catch((err) => {
        Message.error(err.response.data)
      })
    }
  }
}
</script>

<style>
</style>
