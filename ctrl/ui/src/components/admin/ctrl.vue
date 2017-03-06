<template>
<div id="content" class="main">
  <div v-loading.lock="loading">
    <el-tabs v-model="activeName" @tab-click="handleClick" >
      <el-tab-pane label="general" name="general">
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="run mode">
            <el-select v-model="form.runmode" placeholder="select run mode">
              <el-option v-for="item in optionRunModes" :label="item.name" :value="item.value"> </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="enable docs"> <el-switch :disabled="disabled.enabledocs" on-text="" off-text="" v-model="enabledocs"></el-switch> </el-form-item>
          <el-form-item label="etcd endpoints"><el-input :disabled="disabled.etcdendpoints" v-model="form.etcdendpoints"></el-input>(addr1,addr2...)</el-form-item>
          <el-form-item label="cache module"><el-input :disabled="disabled.cachemodule" v-model="form.cachemodule"></el-input></el-form-item>
          <el-form-item label="sess lifetime"> <el-input :disabled="disabled.sessiongcmaxlifetime" v-model="form.sessiongcmaxlifetime"></el-input></el-form-item>
          <el-form-item label="cookie lifetime"> <el-input :disabled="disabled.sessioncookielifetime" v-model="form.sessioncookielifetime"></el-input></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="auth" name="auth">
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="auth module"><el-input :disabled="disabled.authmodule" v-model="form.authmodule"></el-input></el-form-item>
          <el-form-item label="ldap addr"><el-input :disabled="disabled.ldapaddr" v-model="form.ldapaddr"></el-input></el-form-item>
          <el-form-item label="ldap basedn"><el-input :disabled="disabled.ldapbasedn" v-model="form.ldapbasedn"></el-input></el-form-item>
          <el-form-item label="ldap binddn"><el-input :disabled="disabled.ldapbinddn" v-model="form.ldapbinddn"></el-input></el-form-item>
          <el-form-item label="ldap bindpwd"><el-input :disabled="disabled.ldapbindpwd" v-model="form.ldapbindpwd"></el-input></el-form-item>
          <el-form-item label="ldap filter"><el-input :disabled="disabled.ldapfilter" v-model="form.ldapfilter"></el-input></el-form-item>

          <el-form-item label="misso redirect url"><el-input :disabled="disabled.missoredirecturl" v-model="form.missoredirecturl"></el-input></el-form-item>
          <el-form-item label="github client id"><el-input :disabled="disabled.githubclientid" v-model="form.githubclientid"></el-input></el-form-item>
          <el-form-item label="github client secret"><el-input :disabled="disabled.githubclientsecret" v-model="form.githubclientsecret"></el-input></el-form-item>
          <el-form-item label="github redirect url"><el-input :disabled="disabled.githubredirecturl" v-model="form.githubredirecturl"></el-input></el-form-item>
          <el-form-item label="google client id"><el-input :disabled="disabled.googleclientid" v-model="form.googleclientid"></el-input></el-form-item>
          <el-form-item label="google client secret"><el-input :disabled="disabled.googleclientsecret" v-model="form.googleclientsecret"></el-input></el-form-item>
          <el-form-item label="google redirect url"><el-input :disabled="disabled.googleredirecturl" v-model="form.googleredirecturl"></el-input></el-form-item>
        </el-form>
      </el-tab-pane>

    </el-tabs>

    <el-form label-position="right" label-width="200px" :model="form">
      <el-form-item>
        <el-button type="primary" @click="putData">Update</el-button>
        <el-button @click="fetchData">Reset</el-button>
      </el-form-item>
    </el-form>
  </div>
</div>
</template>

<script>
import { fetch, Msg } from 'src/utils'
export default {
  data () {
    return {
      loading: false,
      activeName: 'general',
      enabledocs: false,
      optionRunModes: [{
        name: 'production', value: 'prod'
      }, { name: 'develop', value: 'dev'
      }],
      form: {
        runmode: '',
        enabledocs: '',
        etcdendpoints: '',
        sessiongcmaxlifetime: '',
        sessioncookielifetime: '',
        authmodule: '',
        cachemodule: '',
        ldapaddr: '',
        ldapbasedn: '',
        ldapbinddn: '',
        ldapbindpwd: '',
        ldapfilter: '',
        missoredirecturl: '',
        githubclientid: '',
        githubclientsecret: '',
        githubredirecturl: '',
        googleclientid: '',
        googleclientsecret: '',
        googleredirecturl: ''
      },
      disabled: {
        runmode: false,
        enabledocs: false,
        sessiongcmaxlifetime: false,
        sessioncookielifetime: false,
        authmodule: false,
        cachemodule: false,
        ldapaddr: false,
        ldapbasedn: false,
        ldapbinddn: false,
        ldapbindpwd: false,
        ldapfilter: false,
        missoredirecturl: false,
        githubclientid: false,
        githubclientsecret: false,
        githubredirecturl: false,
        googleclientid: false,
        googleclientsecret: false,
        googleredirecturl: false
      }
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    handleClick (tab, event) {
    },
    fetchData () {
      this.loading = true
      fetch({
        method: 'get',
        url: 'admin/config/ctrl'
      }).then((res) => {
        for (let k in this.disabled) {
          this.disabled[k] = false
        }
        for (let k in this.form) {
          if (res.data[2] && res.data[2][k]) {
            this.form[k] = res.data[2][k]
            this.disabled[k] = true
          } else if (res.data[1] && res.data[1][k]) {
            this.form[k] = res.data[1][k]
          } else if (res.data[0] && res.data[0][k]) {
            this.form[k] = res.data[0][k]
          }
        }
        this.enabledocs = this.form.enabledocs === 'true'
        this.loading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    },
    putData () {
      this.loading = true
      let conf = {}
      this.form.enabledocs = this.enabledocs ? 'true' : 'false'

      for (let k in this.form) {
        if (!this.disabled[k] && this.form[k] !== '') {
          conf[k] = this.form[k]
        }
      }

      // update
      fetch({
        method: 'put',
        url: 'admin/config/ctrl',
        data: JSON.stringify(conf)
      }).then((res) => {
        Msg.success('update success')
        this.loading = false
      }).catch((err) => {
        Msg.error('update failed', err)
        this.loading = false
      })
    }
  }
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.el-input {
  width:380px;
}
</style>
