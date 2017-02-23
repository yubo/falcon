<template>
  <div id="content" class="main">
    <div class="form-inline" role="form">
      <div class="form-group">
	<el-date-picker
	  v-model="timestamp"
	  type="datetimerange"
	  :picker-options="pickerOptions"
	  placeholder="select time range"
	  align="right">
	</el-date-picker>
      </div>
      <button type="button" @click="handleQuery" class="btn btn-primary">search</button>
    </div>

    <el-table v-loading="loading" :data="tableData" border style="width: 100%" class="mt20">
      <el-table-column prop="user"      label="user"> </el-table-column>
      <el-table-column prop="action_id" label="action" width="100px;"> <template scope="scope"> {{actionOps[scope.row.action_id]}} </template> </el-table-column>
      <el-table-column prop="module"    label="module" width="100px;"> <template scope="scope"> {{moduleOps[scope.row.module]}} </template> </el-table-column>
      <el-table-column prop="id"        label="obj id" width="100px;"> </el-table-column>
      <el-table-column prop="time"      label="time"> </el-table-column>
      <el-table-column label="command" width="100px;">
        <template scope="scope">
          <el-button @click="detailObj(scope.row)" type="text" size="small">detail</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pull-right mt20">
      <el-pagination
        @size-change="sizeChange"
        @current-change="curChange"
        :current-page="cur"
        :page-sizes="pageSizes"
        :page-size="per"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total">
      </el-pagination>
    </div>

    <el-dialog title="detail" v-model="detailVisible">
      {{content}}
    </el-dialog>

  </div>
</template>

<script>
import { fetch, Msg } from 'src/utils'
export default {
  data () {
    return {
      loading: false,
      timestamp: [],
      per: 10,
      cur: 1,
      total: 0,
      pageSizes: [5, 10, 20, 50],
      tableData: [],
      detailVisible: false,
      content: '',
      actionOps: ['add', 'del', 'edit'],
      moduleOps: ['host', 'role', 'system', 'tag', 'user', 'token', 'tpl', 'rule', 'template', 'trigger', 'expression', 'team'],
      pickerOptions: {
        shortcuts: [{
          text: '最近一周',
          onClick (picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近一个月',
          onClick (picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近三个月',
          onClick (picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
            picker.$emit('pick', [start, end])
          }
        }]
      }
    }
  },
  methods: {
    sizeChange (per) {
      this.per = per
      this.fetchObjs()
    },
    curChange (cur) {
      this.cur = cur
      this.fetchObjs()
    },
    handleQuery () {
      this.begin = this.timestamp[0]
      this.end = this.timestamp[1]
      this.reFetchObjs()
    },
    detailObj (obj = {}) {
      this.content = obj.data
      this.detailVisible = true
    },

    reFetchObjs () {
      fetch({
        method: 'get',
        url: 'settings/log/cnt',
        params: {begin: this.begin, end: this.end}
      }).then((res) => {
        this.total = res.data.total
        this.fetchObjs()
      }).catch((err) => {
        Msg.error('get failed', err)
      })
    },

    fetchObjs (opts = {begin: this.begin, end: this.end, per: this.per, offset: this.offset}) {
      this.loading = true
      fetch({
        method: 'get',
        url: 'settings/log/search',
        params: opts
      }).then((res) => {
        this.tableData = res.data
        this.loading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    }
  },
  computed: {
    offset () {
      return (this.per * (this.cur - 1))
    },
    begin () {
      return this.timestamp[0] ? this.timestamp[0] : ''
    },
    end () {
      return this.timestamp[1] ? this.timestamp[1] : ''
    }
  },
  created () {
    if (this.$route.query.begin && this.$route.query.end) {
      this.timestamp = [
        this.$route.query.begin, this.$route.query.end
      ]
    }
    if (this.$route.query.per) {
      this.per = this.$route.query.per
    }
    if (this.$route.query.cur) {
      this.cur = this.$route.query.cur
    }
    this.reFetchObjs()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
