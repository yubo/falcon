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
import { fetch } from 'src/utils'
import { Message } from 'element-ui'

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
        Message.error(err.response.data)
      })
    }
  }
}
</script>

<style>
</style>
