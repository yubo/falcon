<template>
<div id="content" class="main">
  <div v-loading.lock="loading">
    <el-tabs v-model="activeName" @tab-click="handleClick" >
      <el-tab-pane label="General" name="general">
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="debug"> <el-switch on-text="" off-text="" v-model="debug"></el-switch> </el-form-item>
          <el-form-item label="dsn"><el-input v-model="form.dsn"></el-input></el-form-item>
          <el-form-item label="db max idle"><el-input v-model="form.dbmaxidle"></el-input></el-form-item>
          <el-form-item label="http">
            <el-switch on-text="" off-text="" v-model="httpenable" desabled></el-switch> 
            <el-input  v-model="form.httpaddr"></el-input>
          </el-form-item>
          <el-form-item label="rpc">
            <el-switch on-text="" off-text="" v-model="rpcenable"></el-switch> 
            <el-input  v-model="form.rpcaddr"></el-input>
          </el-form-item>
          <el-form-item label="grpc">
            <el-switch on-text="" off-text="" v-model="grpcenable"></el-switch> 
            <el-input  v-model="form.grpcaddr"></el-input>
          </el-form-item>
          <el-form-item label="rrd_storage"> <el-input v-model="form.rrd_storage"></el-input></el-form-item>
          <el-form-item label="call timeout"> <el-input v-model="form.calltimeout"></el-input></el-form-item>
          <el-form-item label="leasettl"> <el-input v-model="form.leasettl"></el-input></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="Migrate" name="migrate">
        扩容相关操作请使用admin->expansion->graph
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="enable"> <el-switch on-text="" off-text="" v-model="migrate_enabled" disabled></el-switch> </el-form-item>
          <el-form-item label="concurrency"><el-input v-model="form.migrate_concurrency"></el-input></el-form-item>
          <el-form-item label="replicas"><el-input v-model="form.migrate_replicas" disabled></el-input></el-form-item>
          <el-form-item label="cluster"> <el-input style="width:380px;" type="textarea" :autosize="{minRows:2, maxRows:10}" v-model="migrate_cluster" disabled></el-input></el-form-item>
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
      debug: false,
      httpenable: false,
      rpcenable: false,
      grpcenable: false,
      migrate_enabled: false,
      migrate_cluster: '',
      optionRunModes: [{
        name: 'production', value: 'prod'
      }, { name: 'develop', value: 'dev'
      }],
      form: {
        debug: '',
        dsn: '',
        dbmaxidle: '',
        httpenable: '',
        httpaddr: '',
        rpcenable: '',
        rpcaddr: '',
        grpcenable: '',
        grpcaddr: '',
        rrd_storage: '',
        calltimeout: '',
        leasettl: '',
        migrate_enabled: '',
        migrate_concurrency: '',
        migrate_replicas: '',
        migrate_cluster: ''
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
        url: 'admin/config/graph'
      }).then((res) => {
        for (let k in this.form) {
          if (res.data[2] && res.data[2][k]) {
            this.form[k] = res.data[2][k]
          } else if (res.data[1] && res.data[1][k]) {
            this.form[k] = res.data[1][k]
          } else if (res.data[0] && res.data[0][k]) {
            this.form[k] = res.data[0][k]
          }
        }
        this.debug = this.form.debug === 'true'
        this.httpenable = this.form.httpenable === 'true'
        this.rpcenable = this.form.rpcenable === 'true'
        this.grpcenable = this.form.grpcenable === 'true'
        this.migrate_enabled = this.form.migrate_enabled === 'true'
        this.migrate_cluster = this.form.migrate_cluster.split(';').join('\n')
        this.loading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    },
    putData () {
      this.loading = true
      let conf = {}
      this.form.debug = this.debug ? 'true' : 'false'
      this.form.httpenable = this.httpenable ? 'true' : 'false'
      this.form.rpcenable = this.rpcenable ? 'true' : 'false'
      this.form.grpcenable = this.grpcenable ? 'true' : 'false'
      this.form.migrate_enabled = this.migrate_enabled ? 'true' : 'false'

      for (let k in this.form) {
        if (this.form[k] !== '') {
          conf[k] = this.form[k]
        }
      }

      // update
      fetch({
        method: 'put',
        url: 'admin/config/graph',
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
