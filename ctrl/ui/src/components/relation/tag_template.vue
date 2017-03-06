
<template>
  <div id="content">
    <div class="form-inline" role="form">
      <div class="form-group">
        <input type="text" v-model="query" @keyup.enter="handleQuery" class="form-control" placeholder="template name">
      </div>
      <button type="button" @click="handleQuery" class="btn btn-primary">Search</button>
      <input type="checkbox" v-model="deep" class="form-control">
      <span>搜索子节点</span>
      <input type="checkbox" v-model="mine" class="form-control">
      <span>mine</span>
      <div class="pull-right">
        <div class="input-group">
          <span class="input-group-addon">template</span> 
          <el-select
            style="width: 100%"
            placeholder="template name"
            v-model="tpls"
            multiple
            filterable
            remote
            :remote-method="getTpls"
            :loading="sloading">
            <el-option
              v-for="tpl in optionTpls"
              :key="tpl.id"
              :label="tpl.name"
              :value="tpl.id">
            </el-option>
          </el-select>
        </div>
        <button :disabled="!isOperator" type="button" class="btn btn-primary" @click="handleBind">Bind</button>
      </div>
    </div>


    <el-table v-loading.lock="loading" :data="tableData" border style="width: 100%" class="mt20" @selection-change="handleSelectionChange">
      <el-table-column :prop="curTag.name" :label="curTag.name" width="100%">
        <el-table-column type="selection"> </el-table-column>
        <el-table-column prop="id" label="id"> </el-table-column>
        <el-table-column prop="tag_name"  label="tag"> </el-table-column>
        <el-table-column prop="tpl_name"  label="template"> </el-table-column>
        <el-table-column label="command">
          <template scope="scope">
            <el-button :disabled="!isOperator" @click="unbind(scope.row)" type="danger" size="small">Unbind</el-button>
          </template>
        </el-table-column>
      </el-table-column>
    </el-table>

    <div class="mt20">
      <button :disabled="!isOperator" @click="mUnbind" type="button" class="btn btn-danger">Unbind</button>

      <div class="pull-right">
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
    </div>
  </div>
</template>

<script>
// import store from 'src/store'
import { fetch, Msg } from 'src/utils'

export default {
  data () {
    return {
      loading: false,
      sloading: false,
      deep: true,
      mine: true,
      tpls: [],
      optionTpls: [],
      query: '',
      per: 10,
      cur: 1,
      total: 0,
      pageSizes: [5, 10, 20, 50],
      multipleSelection: [],
      tableData: []
    }
  },
  watch: {
    'curTagId': function (val) {
      this.handleQuery()
    }
  },
  methods: {
    handleSelectionChange (val) {
      this.multipleSelection = val
    },
    getTpls (query) {
      if (query !== '') {
        this.sloading = true
        fetch({
          method: 'get',
          url: 'template/search',
          params: {
            query: query,
            per: 10
          }
        }).then((res) => {
          this.optionTpls = res.data
          this.sloading = false
        }).catch((err) => {
          Msg.error('get failed', err)
          this.sloading = false
        })
      } else {
        this.optionTpls = []
      }
    },
    sizeChange (per) {
      this.per = per
      this.fetchData()
    },
    curChange (cur) {
      this.cur = cur
      this.fetchData()
    },

    handleQuery () {
      this.reFetchData()
    },

    reFetchData () {
      fetch({
        method: 'get',
        url: 'rel/tag/template/cnt',
        params: { tag_id: this.curTagId, query: this.query, deep: this.deep, mine: this.mine }
      }).then((res) => {
        this.total = res.data.total
        this.fetchData()
      }).catch((err) => {
        Msg.error('get failed', err)
      })
    },

    fetchData (opts = {
      tag_id: this.curTagId,
      query: this.query,
      deep: this.deep,
      mine: this.mine,
      own: this.own,
      per: this.per,
      offset: this.offset}) {
      this.loading = true
      fetch({
        method: 'get',
        url: 'rel/tag/template/search',
        params: opts
      }).then((res) => {
        this.tableData = res.data
        this.loading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    },
    handleBind (orgs) {
      this.loading = true
      fetch({
        method: 'post',
        url: 'rel/tag/templates',
        data: {tag_id: this.curTagId, tpl_ids: this.tpls}
      }).then((res) => {
        Msg.success('success!')
        this.total++
        // loading will unset at fetchdata done
        this.fetchData()
      }).catch((err) => {
        Msg.error('update failed', err)
        this.loading = false
      })
    },
    unbind (tpl) {
      Msg.confirm('此操作将解绑定该记录, 是否继续?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        this.loading = true
        fetch({
          method: 'delete',
          url: 'rel/tag/template',
          data: {
            tag_id: this.tag_id,
            tpl_id: tpl.tpl_id
          }
        }).then((res) => {
          Msg.success('success!')
          this.total--
          this.fetchData()
        }).catch((err) => {
          Msg.error('delete failed', err)
          this.loading = false
        })
      }).catch(() => {
        Msg.info('cancel')
      })
    },
    mUnbind () {
      Msg.confirm('此操作将解绑定该记录, 是否继续?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        this.loading = true
        fetch({
          method: 'delete',
          url: 'rel/tag/templates',
          data: {
            tag_id: this.curTagId,
            tpl_ids: this.multipleSelection.map((val) => { return val.id })
          }
        }).then((res) => {
          Msg.success('success!')
          this.total = this.total - res.data.total
          this.fetchData()
        }).catch((err) => {
          Msg.error('delete failed', err)
          this.loading = false
        })
      }).catch(() => {
        Msg.info('cancel')
      })
    }
  },
  computed: {
    isOperator () {
      return this.$store.state.auth.operator
    },
    offset () {
      return (this.per * (this.cur - 1))
    },
    curTagId () {
      return this.$store.state.rel.curTag.id
    },
    curTag () {
      return this.$store.state.rel.curTag
    }
  },
  created () {
    if (this.$route.query.query) {
      this.query = this.$route.query.query
    }
    if (this.$route.query.per) {
      this.per = this.$route.query.per
    }
    if (this.$route.query.cur) {
      this.cur = this.$route.query.cur
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
