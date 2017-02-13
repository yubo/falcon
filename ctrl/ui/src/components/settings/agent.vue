<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 class="page-header">Agent Configurations</h1>
  <div v-loading.lock="loading">
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

      <el-form-item label="interval"> <el-input v-model.number="form.interval"> </el-input> </el-form-item>
      <el-form-item label="payloadSize"> <el-input v-model.number="form.payloadSize"> </el-input> </el-form-item>

      <el-form-item label="ifPre"> 
        <el-select
            v-model="form.ifPre"
            multiple
            filterable
            allow-create
            placeholder="interface prefix name">
            <el-option
              v-for="item in ifPreOpts"
              :label="item"
              :value="item">
            </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="upstreams"> 
        <el-select
            v-model="form.upstreams"
            multiple
            filterable
            allow-create
            placeholder="upstreams name">
            <el-option
              v-for="item in upstreamsOpts"
              :label="item"
              :value="item">
            </el-option>
        </el-select>
      </el-form-item>


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
      ifPreOpts: [],
      upstreamsOpts: [],
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
        interval: 0,
        payloadSize: 0,
        ifPre: [],
        upstreams: []
      }
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    fetchData () {
      this.loading = true
      fetch({
        method: 'get',
        url: 'settings/config/agent'
      }).then((res) => {
        for (var k in this.form) {
          if (k === 'params') {
            for (var k1 in this.form[k]) {
              if (res.data[k] && res.data[k][k1]) {
                this.form[k][k1] = res.data[k][k1]
              }
            }
          } else {
            this.form[k] = res.data[k]
          }
        }
        this.ifPreOpts = this.form.ifPre
        this.upstreamsOpts = this.form.upstreams
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
        url: 'settings/config/agent',
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
