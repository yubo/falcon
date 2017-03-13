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
          <el-form-item label="lease key"> <el-input v-model="form.leasekey"></el-input></el-form-item>
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
          <el-form-item label="cluster"> <el-input v-model="form.judge_cluster"></el-input> (judge-00=127.0.0.1:6080,judge-01=127.0.0.1:6081) </el-form-item>
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
          <el-form-item label="cluster">
            <el-select
              style="width: 100%"
              placeholder="graph cluster"
              v-model="graph_cluster"
              multiple
              filterable
              remote
              :remote-method="getGraphs"
              :loading="sloading">
              <el-option
                v-for="obj in optionGraphs"
                :key="obj.name"
                :label="obj.name"
                :value="obj.name">
              </el-option>
            </el-select>

          </el-form-item>
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
          <el-form-item label="address"><el-input v-model="form.tsdb_address"></el-input>(127.0.0.1:8088)</el-form-item>
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
      sloading: false,
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
      optionGraphs: [],
      graph_cluster: [],
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
        leasekey: '',
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

    getGraphs (query) {
      if (query !== '') {
        this.sloading = true
        fetch({
          method: 'get',
          url: 'admin/online/graph',
          params: {
            query: query,
            per: 10
          }
        }).then((res) => {
          this.optionGraphs = res.data.map((v) => {
            return {name: v.value + '=' + v.key}
          })
          this.sloading = false
        }).catch((err) => {
          Msg.error('get failed', err)
          this.sloading = false
        })
      } else {
        this.optionGraphs = []
      }
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
        this.socket_enabled = this.form.socket_enable === 'true'
        this.judge_enabled = this.form.judge_enabled === 'true'
        this.graph_enabled = this.form.graph_enabled === 'true'
        this.tsdb_enabled = this.form.tsdb_enabled === 'true'
        this.graph_cluster = this.form.graph_cluster.split(',')
        this.optionGraphs = this.graph_cluster.map((v) => {
          return {name: v}
        })
        this.loading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    },
    putData () {
      this.loading = true
      let conf = {}
      this.form.enable = this.enable ? 'true' : 'false'
      this.form.httpenable = this.httpenable ? 'true' : 'false'
      this.form.rpcenable = this.rpcenable ? 'true' : 'false'
      this.form.socket_enable = this.socket_enable ? 'true' : 'false'
      this.form.judge_enabled = this.judge_enabled ? 'true' : 'false'
      this.form.graph_enabled = this.graph_enabled ? 'true' : 'false'
      this.form.tsdb_enabled = this.tsdb_enabled ? 'true' : 'false'
      this.form.graph_cluster = this.graph_cluster.join(',')

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
