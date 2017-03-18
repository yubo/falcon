<template>
<div id="content" class="main">
  <div v-loading.lock="loading">
    <el-tabs v-model="activeName" @tab-click="handleClick" >
      <el-tab-pane label="General" name="general">
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="debug"> <el-switch on-text="" off-text="" v-model="debug"></el-switch> </el-form-item>
          <el-form-item label="minstep"> <el-input v-model="form.minstep"></el-input> </el-form-item>
          <el-form-item label="http">
            <el-switch on-text="" off-text="" v-model="httpenable"></el-switch> 
            <el-input  v-model="form.httpaddr"></el-input>
          </el-form-item>
          <el-form-item label="rpc">
            <el-switch on-text="" off-text="" v-model="rpcenable"></el-switch> 
            <el-input  v-model="form.rpcaddr"></el-input>
          </el-form-item>
          <el-form-item label="socket">
            <el-switch on-text="" off-text="" v-model="socket_enable"></el-switch> 
          </el-form-item>
          <el-form-item label="socket listen">
            <el-input  v-model="form.socket_listen"></el-input>
          </el-form-item>
          <el-form-item label="socket timeout">
            <el-input  v-model="form.socket_timeout"></el-input>
          </el-form-item>
          <el-form-item label="leasettl"> <el-input v-model="form.leasettl"></el-input></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="judge" name="judge">
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="enable"> <el-switch on-text="" off-text="" v-model="judge_enabled"></el-switch> </el-form-item>
          <el-form-item label="batch"><el-input v-model="form.judge_batch"></el-input></el-form-item>
          <el-form-item label="conntimeout"><el-input v-model="form.judge_conntimeout"></el-input></el-form-item>
          <el-form-item label="calltimeout"><el-input v-model="form.judge_calltimeout"></el-input></el-form-item>
          <el-form-item label="maxconns"><el-input v-model="form.judge_maxconns"></el-input></el-form-item>
          <el-form-item label="maxidle"><el-input v-model="form.judge_maxidle"></el-input></el-form-item>
          <el-form-item label="replicas"><el-input v-model="form.judge_replicas"></el-input></el-form-item>
          <el-form-item label="cluster"> <el-input style="width:380px;" type="textarea" :autosize="{minRows:2, maxRows:10}" v-model="judge_cluster"></el-input></el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="graph" name="graph">
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="enable"> <el-switch on-text="" off-text="" v-model="graph_enabled"></el-switch> </el-form-item>
          <el-form-item label="batch"><el-input v-model="form.graph_batch"></el-input></el-form-item>
          <el-form-item label="conntimeout"><el-input v-model="form.graph_conntimeout"></el-input></el-form-item>
          <el-form-item label="calltimeout"><el-input v-model="form.graph_calltimeout"></el-input></el-form-item>
          <el-form-item label="maxconns"><el-input v-model="form.graph_maxconns"></el-input></el-form-item>
          <el-form-item label="maxidle"><el-input v-model="form.graph_maxidle"></el-input></el-form-item>
          <el-form-item label="replicas"><el-input v-model="form.graph_replicas"></el-input></el-form-item>
          <el-form-item label="cluster"> <el-input style="width:380px;" type="textarea" :autosize="{minRows:2, maxRows:10}" v-model="graph_cluster"></el-input> </el-form-item>
          <!--
          <el-form-item label="cluster">
            <el-checkbox-group v-model="graph_cluster" style="width: 80%">
              <el-checkbox v-for="graph in graphs" :label="graph">{{graph}}</el-checkbox>
            </el-checkbox-group>
          </el-form-item>
          -->

        </el-form>
      </el-tab-pane>

      <el-tab-pane label="tsdb" name="tsdb">
        <el-form label-position="right" label-width="200px" :model="form">
          <el-form-item label="enable"> <el-switch on-text="" off-text="" v-model="tsdb_enabled"></el-switch> </el-form-item>
          <el-form-item label="batch"><el-input v-model="form.tsdb_batch"></el-input></el-form-item>
          <el-form-item label="conntimeout"><el-input v-model="form.tsdb_conntimeout"></el-input></el-form-item>
          <el-form-item label="calltimeout"><el-input v-model="form.tsdb_calltimeout"></el-input></el-form-item>
          <el-form-item label="maxconns"><el-input v-model="form.tsdb_maxconns"></el-input></el-form-item>
          <el-form-item label="maxidle"><el-input v-model="form.tsdb_maxidle"></el-input></el-form-item>
          <el-form-item label="retry"><el-input v-model="form.tsdb_retry"></el-input></el-form-item>
          <el-form-item label="address"> <el-input style="width:380px;" type="textarea" :autosize="{minRows:2, maxRows:10}" v-model="tsdb_address"></el-input></el-form-item>
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
const { _ } = window
export default {
  data () {
    return {
      loading: false,
      activeName: 'general',
      debug: false,
      httpenable: false,
      rpcenable: false,
      socket_enable: false,
      judge_enabled: false,
      graph_enabled: false,
      tsdb_enabled: false,
      optionRunModes: [{
        name: 'production', value: 'prod'
      }, { name: 'develop', value: 'dev'
      }],
      judge_cluster: '',
      graph_cluster: '',
      tsdb_address: '',
      // graphs: [],
      form: {
        debug: '',
        minstep: '',
        httpenable: '',
        httpaddr: '',
        rpcenable: '',
        rpcaddr: '',
        socket_enable: '',
        socket_listen: '',
        socket_timeout: '',
        leasettl: '',
        judge_enabled: '',
        judge_batch: '',
        judge_conntimeout: '',
        judge_calltimeout: '',
        judge_maxconns: '',
        judge_maxidle: '',
        judge_replicas: '',
        judge_cluster: '',
        graph_enabled: '',
        graph_batch: '',
        graph_conntimeout: '',
        graph_calltimeout: '',
        graph_maxconns: '',
        graph_maxidle: '',
        graph_replicas: '',
        graph_cluster: '',
        tsdb_enabled: '',
        tsdb_batch: '',
        tsdb_conntimeout: '',
        tsdb_calltimeout: '',
        tsdb_maxconns: '',
        tsdb_maxidle: '',
        tsdb_retry: '',
        tsdb_address: ''
      }
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    handleClick (tab, event) {
    },
    getGraphs () {
      fetch({
        method: 'get',
        url: 'admin/online/graph'
      }).then((res) => {
        this.graphs = _.union(res.data.map((v) => {
          return v.value + '=' + v.key
        }), this.graphs).sort()
      }).catch((err) => {
        Msg.error('get failed', err)
      })
    },
    fetchData () {
      this.loading = true
      fetch({
        method: 'get',
        url: 'admin/config/transfer'
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
        this.socket_enable = this.form.socket_enable === 'true'
        this.judge_enabled = this.form.judge_enabled === 'true'
        this.graph_enabled = this.form.graph_enabled === 'true'
        this.tsdb_enabled = this.form.tsdb_enabled === 'true'

        this.judge_cluster = this.form.judge_cluster.split(';').join('\n')
        this.graph_cluster = this.form.graph_cluster.split(';').join('\n')
        this.tsdb_address = this.form.tsdb_address.split(';').join('\n')
        // this.graphs = this.graph_cluster.sort()
        this.loading = false
        // this.getGraphs()
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
      this.form.socket_enable = this.socket_enable ? 'true' : 'false'
      this.form.judge_enabled = this.judge_enabled ? 'true' : 'false'
      this.form.graph_enabled = this.graph_enabled ? 'true' : 'false'
      this.form.tsdb_enabled = this.tsdb_enabled ? 'true' : 'false'
      this.form.judge_cluster = this.judge_cluster.split('\n').join(';')
      this.form.graph_cluster = this.graph_cluster.split('\n').join(';')
      this.form.tsdb_address = this.tsdb_address.split('\n').join(';')
      console.log(this.form.graph_cluster)

      for (let k in this.form) {
        if (this.form[k] !== '') {
          conf[k] = this.form[k]
        }
      }

      // update
      fetch({
        method: 'put',
        url: 'admin/config/transfer',
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
