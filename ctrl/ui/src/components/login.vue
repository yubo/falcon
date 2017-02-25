<template>
  <div class="container-fluid">
    <div class="row login">
      <div class="col-sm-6 col-sm-offset-3">
        <div class="panel panel-default">
          <div class="panel-heading"> LDAP </div>
          <div class="panel-body">
            <el-form label-position="right" label-width="80px" :model="ldapForm">
              <el-form-item label="username"><el-input v-model="ldapForm.username"></el-input></el-form-item>
              <el-form-item label="passworld"><el-input v-model="ldapForm.password" type="password"></el-input></el-form-item>
              <el-form-item>
                <el-button type="primary" @click="ldapLogin">Sign in</el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </div>
      <div class="col-sm-6 col-sm-offset-3">
        <div class="panel panel-default">
          <div class="panel-body">
            <el-button v-if="auth.misso" type="primary" @click="authLogin('misso')">misso</el-button>
            <el-button v-if="auth.github" type="primary" @click="authLogin('github')">github</el-button>
            <el-button v-if="auth.google" type="primary" @click="authLogin('google')">google</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { fetch, Msg } from 'src/utils'
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
      this.$store.dispatch('auth/login', this.ldapForm).then(() => {
        this.$store.dispatch('load_config')
        this.$router.push(this.$route.query.cb ? this.$route.query.cb : '/')
      })
    },
    authLogin (module) {
      window.location.href = '/v1.0/auth/login/' + module + '?cb=' + this.$route.query.cb
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
        Msg.error('get auth modules failed', err)
      })
    }
  },
  computed: {
    isLogin () {
      return this.$store.state.auth.login
    },
    loading () {
      return this.$store.state.auth.loading
    }
  },
  created () {
    console.log('hello')
    this.fetchObjs()
  }
}
</script>

<style>
.login {
  margin-top:30px;
}
</style>
