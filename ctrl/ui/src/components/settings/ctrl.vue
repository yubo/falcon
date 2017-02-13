<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 class="page-header">Ctrl Configurations</h1>
  <div v-loading.lock="loading">
    <el-tabs v-model="activeName" @tab-click="handleClick" >
      <el-tab-pane label="core" name="core">
        <el-form label-position="right" label-width="100px" :model="form">
          <el-form-item label="name"> <el-input v-model="form.params.name"> </el-input> </el-form-item>
          <el-form-item label="disabled"> <el-switch on-text="" off-text="" v-model="form.params.disabled"></el-switch> </el-form-item>
          <el-form-item label="debug level"> <el-input v-model.number="form.params.debug" > </el-input> </el-form-item>
          <el-form-item label="http"> <el-switch on-text="" off-text="" v-model="form.params.http"></el-switch> </el-form-item>
          <el-form-item label="http addr"> <el-input v-model="form.params.httpAddr"> </el-input> </el-form-item>
          <el-form-item label="rpc"> <el-switch on-text="" off-text="" v-model="form.params.rpc"></el-switch> </el-form-item>
          <el-form-item label="rpc addr"> <el-input v-model="form.params.rpcAddr"> </el-input> </el-form-item>
          <el-form-item label="conn timeout"> <el-input v-model.number="form.params.connTimeout"> </el-input> </el-form-item>
          <el-form-item label="call timeout"> <el-input v-model.number="form.params.callTimeout"> </el-input> </el-form-item>
          <el-form-item label="concurrency"> <el-input v-model.number="form.params.concurrency"> </el-input> </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="webconfig" name="webconfig">
        <el-form label-position="right" label-width="100px" :model="form">
          <el-form-item label="app name"><el-input v-model="form.app_name"></el-input></el-form-item>
          <el-form-item label="run mode"><el-input v-model="form.run_mode"></el-input></el-form-item>
          <el-form-item label="http port"><el-input v-model.number="form.http_port"></el-input></el-form-item>
          <el-form-item label="enable docs"><el-switch on-text="" off-text="" v-model="form.enable_docs"></el-switch></el-form-item>
          <el-form-item label="sess name"><el-input v-model="form.session_name"></el-input></el-form-item>
          <el-form-item label="sess gc lifetime"><el-input v-model.number="form.session_gc_max_lifetime"></el-input></el-form-item>
          <el-form-item label="cookie lifetime"><el-input v-model.number="form.session_cookie_lifetime"></el-input></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="server" name="server">
        <el-form label-position="right" label-width="100px" :model="form">
          <el-form-item label="metric file"><el-input v-model="form.metric_file"></el-input></el-form-item>
          <el-form-item label="auth module"><el-input v-model="form.auth_module"></el-input></el-form-item>
          <el-form-item label="cache module"><el-input v-model="form.cache_module"></el-input></el-form-item>
          <el-form-item label="ldap addr"><el-input v-model="form.ldap_addr"></el-input></el-form-item>
          <el-form-item label="ldap basedn"><el-input v-model="form.ldap_basedn"></el-input></el-form-item>
          <el-form-item label="ldap binddn"><el-input v-model="form.ldap_binddn"></el-input></el-form-item>
          <el-form-item label="ldap bindpwd"><el-input v-model="form.ldap_bindpwd"></el-input></el-form-item>
          <el-form-item label="ldap filter"><el-input v-model="form.ldap_filter"></el-input></el-form-item>
          <el-form-item label="ldap tls"><el-switch on-text="" off-text="" v-model="form.ldap_tls"></el-switch></el-form-item>
          <el-form-item label="ldap debug"><el-switch on-text="" off-text="" v-model="form.ldap_debug"></el-switch></el-form-item>
        </el-form>
      </el-tab-pane>

    </el-tabs>

    <el-form label-position="right" label-width="100px" :model="form">
      <el-form-item>
        <el-button type="primary" @click="putData">Update</el-button>
        <el-button @click="fetchData">Reset</el-button>
      </el-form-item>
    </el-form>
  </div>
</div>
</template>

<script>
import { Message } from 'element-ui'
import { fetch } from 'src/utils'
export default {
  data () {
    return {
      loading: false,
      activeName: 'core',
      form: {
        params: {
          debug: 0,
          connTimeout: 0,
          callTimeout: 0,
          concurrency: 0,
          disabled: false,
          http: false,
          httpAddr: '',
          rpc: false,
          rpcAddr: '',
          name: ''
        },
        app_name: '',
        run_mode: '',
        http_port: 0,
        enable_docs: false,
        session_name: '',
        session_gc_max_lifetime: 0,
        session_cookie_lifetime: 0,
        metric_file: '',
        auth_module: '',
        cache_module: '',
        ldap_addr: '',
        ldap_basedn: '',
        ldap_binddn: '',
        ldap_bindpwd: '',
        ldap_filter: '',
        ldap_tls: '',
        ldap_debug: ''
      }
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    handleClick (tab, event) {
      console.log(tab, event)
    },
    fetchData () {
      this.loading = true
      fetch({
        method: 'get',
        url: 'settings/config/ctrl'
      }).then((res) => {
        for (var k in this.form) {
          if (typeof this.form[k] === 'object') {
            for (var k1 in this.form[k]) {
              if (res.data[k] && res.data[k][k1]) {
                this.form[k][k1] = res.data[k][k1]
              }
            }
          } else {
            this.form[k] = res.data[k]
          }
        }
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    putData () {
      this.loading = true
      // update
      fetch({
        method: 'put',
        url: 'settings/config/ctrl',
        data: JSON.stringify(this.form)
      }).then((res) => {
        Message.success('update success')
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    }
  }
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
