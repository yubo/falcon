<template>
  <el-dialog title="login" v-model="loginVisible"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="false" >
    <div v-loading="loading">
      <el-tabs v-model="activeName" @tab-click="handleClick" >
        <el-tab-pane label="ldap" name="ldap">
          <el-form label-position="right" label-width="80px" :model="ldapForm">
            <el-form-item label="username"><el-input v-model="ldapForm.username"></el-input></el-form-item>
            <el-form-item label="passworld"><el-input v-model="ldapForm.password"></el-input></el-form-item>
            <el-form-item> <el-button type="primary" @click="ldapLogin">Sign in</el-button> </el-form-item>
        </el-tab-pane>

        <el-tab-pane label="misso" name="misso">
          misso
        </el-tab-pane>
      </el-tabs>
    </div>
  </el-dialog>
</template>

<script>
export default {
  data () {
    return {
      activeName: 'ldap',
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
      console.log(tab, event)
    },
    ldapLogin () {
      this.$store.dispatch('login/login', this.ldapForm)
    }
  },
  computed: {
    login () {
      return this.$store.state.login.login
    },
    loading () {
      return this.$store.state.login.loading
    }
  },
  created () {
    this.loginVisible = !this.login
    console.log('login : ', this.login)
    if (window.Cookies.get('username') !== undefined) {
      this.$store.dispatch('login/login')
    }
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
