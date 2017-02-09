<template>
  <div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
    <div class="form-inline" role="form">
      <div class="form-group">
        <input type="text" v-model="query" @keyup.enter="handleQuery" class="form-control" placeholder="template name">
      </div>
      <button type="button" @click="handleQuery" class="btn btn-primary">search</button>
      <div class="pull-right">
        <router-link to="/meta/template/edit" class="btn btn-default"><span class="glyphicon glyphicon-plus"></span>Add</router-link>
      </div>
    </div>

    <el-table v-loading="loading" :data="tableData" border style="width: 100%" class="mt20">
      <el-table-column prop="name"    label="name"> </el-table-column>
      <el-table-column prop="pname"   label="pname"> </el-table-column>
      <el-table-column prop="creator" label="creator"> </el-table-column>
      <el-table-column label="command">
        <template scope="scope">
          <el-button @click="handleClone(scope.row)" type="text" size="small">CLONE</el-button>
          <el-button @click="handleBinded(scope.row)" type="text" size="small">BINDED</el-button>
          <el-button @click="handleEdit(scope.row)" type="text" size="small">EDIT</el-button>
          <el-button @click="handleDelete(scope.row)" type="text" size="small">DELETE</el-button>
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
  </div>
</template>

<script>
import { fetch } from 'src/utils'
import { Message, MessageBox } from 'element-ui'
export default {
  data () {
    return {
      loading: false,
      reQuery: true,
      query: '',
      per: 10,
      cur: 1,
      total: 0,
      pageSizes: [5, 10, 20, 50],
      tableData: []
    }
  },
  methods: {
    sizeChange (per) {
      this.per = per
      this.query = this.$route.query.query
      this.fetchData()
    },
    curChange (cur) {
      this.cur = cur
      this.query = this.$route.query.query
      this.fetchData()
    },

    handleQuery () {
      this.reFetchData()
    },

    reFetchData () {
      fetch({
        router: this.$router,
        method: 'get',
        url: 'template/cnt',
        params: {query: this.query}
      }).then((res) => {
        this.total = res.data.total
        this.fetchData()
      }).catch((err) => {
        Message.error(err.response.data)
      })
    },

    fetchData (opts = {query: this.query, per: this.per, offset: this.offset}) {
      this.loading = true
      fetch({
        router: this.$router,
        method: 'get',
        url: 'template/search',
        params: opts
      }).then((res) => {
        this.tableData = res.data
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    handleClone (template) {
    },
    handleBinded (template) {
    },
    handleEdit (template) {
      this.$router.push({
        path: '/meta/template/edit',
        query: {id: template.id}
      })
    },
    handleDelete (template) {
      MessageBox.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        fetch({
          router: this.$router,
          method: 'delete',
          url: 'template/' + template.id
        }).then((res) => {
          Message.success('success!')
          this.total = this.total - 1
          this.fetchData()
        }).catch((err) => {
          Message.error(err.response.data)
        })
      }).catch(() => {
        Message.info('cancel')
      })
    }
  },
  computed: {
    offset () {
      return (this.per * (this.cur - 1))
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
    this.reFetchData()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
