<template>
<div id="app">
  <navbar></navbar>
  <router-view v-if="login"></router-view>
  <login></login>
</div>
</template>

<script>
import navbar from './navbar'
import login from './login'
import { fetch, Msg } from 'src/utils'

export default {
  components: {
    navbar,
    login
  },
  data () {
    return { }
  },
  computed: {
    login () {
      return this.$store.state.auth.login
    }
  },
  create () {
    if (!this.$store.config) {
      fetch({
        method: 'get',
        url: 'settings/config/ui'
      }).then((res) => {
        this.$store.commit('m_set_config', res.data)
      }).catch((err) => {
        Msg.error('get failed', err)
      })
    }
  }
}
</script>

<style>
</style>
