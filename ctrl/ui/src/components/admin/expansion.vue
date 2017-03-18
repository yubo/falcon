<template>
  <div id="content" class="main">
    <div v-loading.lock="loading">
      <el-tabs v-model="activeName">
        <el-tab-pane label="graph" name="graph">
          <el-form label-position="right" label-width="200px">
            <el-form-item label="migrate status"> <div>{{migrating ? 'running' : 'stoped'}} </div></el-form-item>

            <el-form-item label="graph cluster"> <el-input style="width:380px;" type="textarea" :autosize="{minRows:2, maxRows:10}" v-model="graph_cluster" disabled></el-input></el-form-item>

            <el-form-item label="add endpoint"> <el-input style="width:380px;" type="textarea" :autosize="{minRows:2, maxRows:10}" v-model="new_endpoint"></el-input></el-form-item>

          </el-form>
        </el-tab-pane>
      </el-tabs>

      <el-form label-position="right" label-width="200px">
        <el-form-item>
          <el-button type="danger" @click="beginMigrate">Begin Migrate</el-button>
          <el-button type="danger" @click="finishMigrate">Finish Migrate</el-button>
          <el-button @click="fetchData">Reset</el-button>
          <el-button @click="viewOnline">view online</el-button>
        </el-form-item>
      </el-form>
    </div>
    <el-dialog title="online" v-model="onlineVisible">
      <el-table :data="graph_online" border style="width: 100%" class="mt20">
        <el-table-column prop="endpoint" label="endpoint"> </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script>
import { fetch, Msg } from 'src/utils'
export default {
  data () {
    return {
      loading: false,
      onlineVisible: false,
      activeName: 'graph',
      graph_online: [],
      migrating: false,
      graph_cluster: '',
      new_endpoint: ''
    }
  },
  created () {
    this.fetchData()
    this.getGraphOnline()
  },
  methods: {
    viewOnline () {
      this.getGraphOnline()
      this.onlineVisible = true
    },
    getGraphOnline () {
      fetch({
        method: 'get',
        url: 'admin/online/graph'
      }).then((res) => {
        this.graph_online = res.data.map((v) => {
          return {endpoint: v.key}
        }).sort()
        console.log(this.graph_online)
      }).catch((err) => {
        Msg.error('get failed', err)
      })
    },
    fetchData () {
      fetch({
        method: 'get',
        url: 'admin/expansion/graph'
      }).then((res) => {
        console.log(res.data)
        this.migrating = !!res.data.migrating
        this.graph_cluster = res.data.graph_cluster.split(';').join('\n')
        this.new_endpoint = res.data.new_endpoint.split(';').join('\n')
      }).catch((err) => {
        Msg.error('fetch failed', err)
      })
    },
    beginMigrate () {
      // update
      fetch({
        method: 'put',
        url: 'admin/expansion/graph',
        data: JSON.stringify({
          migrating: true,
          new_endpoint: this.new_endpoint.split('\n').join(';')
        })
      }).then((res) => {
        Msg.success('update success')
        this.fetchData()
      }).catch((err) => {
        Msg.error('update failed', err)
      })
    },
    finishMigrate () {
      fetch({
        method: 'put',
        url: 'admin/expansion/graph',
        data: JSON.stringify({migrating: false})
      }).then((res) => {
        Msg.success('update success')
        this.fetchData()
      }).catch((err) => {
        Msg.error('update failed', err)
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
