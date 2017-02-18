<template>
  <el-dialog title="login" v-model="loginVisible"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="false" >
    <div v-loading="loading">
      <el-tabs v-model="activeName" @tab-click="handleClick" >
        <el-tab-pane v-if="auth.misso || auth.github || auth.google" label="auth" name="auth">
          <el-button v-if="auth.misso" type="primary" @click="authLogin('misso')">misso</el-button>
          <el-button v-if="auth.github" type="primary" @click="authLogin('github')">github</el-button>
          <el-button v-if="auth.google" type="primary" @click="authLogin('google')">google</el-button>
        </el-tab-pane>

        <el-tab-pane v-if="auth.ldap" label="ldap" name="ldap">
          <el-form label-position="right" label-width="80px" :model="ldapForm">
            <el-form-item label="username"><el-input v-model="ldapForm.username"></el-input></el-form-item>
            <el-form-item label="passworld"><el-input v-model="ldapForm.password"></el-input></el-form-item>
            <el-form-item> <el-button type="primary" @click="ldapLogin">Sign in</el-button> </el-form-item>
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-dialog>
</template>

<script>
import { fetch } from 'src/utils'
import { Message } from 'element-ui'
export default {
  data () {
    return {
      auth: {
        misso: false,
        github: false,
        google: false,
        ldap: false
      },
      activeName: '',
      loginVisible: false,
      ldapForm: {
        username: '',
        password: '',
        method: 'ldap'
      }
    }
  },
  methods: {
    handleClick (tab, event) {
    },
    ldapLogin () {
      this.$store.dispatch('auth/login', this.ldapForm)
    },
    authLogin (module) {
      window.location.href = '/v1.0/auth/login/' + module + '?cb=' + this.$router.history.current.fullPath
    },
    fetchObjs () {
      fetch({
        method: 'get',
        url: 'auth/modules'
      }).then((res) => {
        for (let k in res.data) {
          if (this.auth[res.data[k]] !== undefined) {
            this.auth[res.data[k]] = true
          }
        }
      }).catch((err) => {
        Message.error(err.response.data)
      })
    }
  },
  computed: {
    login () {
      return this.$store.state.auth.login
    },
    loading () {
      return this.$store.state.auth.loading
    }
  },
  created () {
    this.loginVisible = !this.login
    if (!this.login) {
      this.$store.dispatch('auth/login')
    }
    this.fetchObjs()
  },
  watch: {
    'login': function (val) {
      this.loginVisible = !val
    }
  }
}
</script>

<style>
</style>
