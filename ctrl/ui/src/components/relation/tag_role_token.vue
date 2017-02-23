<template>
  <div id="content">
    <div class="form-inline" role="form">
      <div class="form-group">
        <input type="text" v-model="query" @keyup.enter="handleQuery" class="form-control" placeholder="token name">
      </div>
      <button type="button" @click="handleQuery" class="btn btn-primary">Search</button>
      <input type="checkbox" v-model="global" class="form-control">
      <span>全局搜索</span>
      <div class="pull-right">
        <div class="input-group">
          <el-select
            style="width: 100%"
            placeholder="role name"
            v-model="roleId"
            filterable
            remote
            :remote-method="getRoles"
            :loading="sloading1">
            <el-option
              v-for="role in optionRoles"
              :key="role.id"
              :label="role.name"
              :value="role.id">
            </el-option>
          </el-select>
        </div>

        <div class="input-group">
          <el-select
            style="width: 100%"
            placeholder="token name"
            v-model="tokenId"
            filterable
            remote
            :remote-method="getTokens"
            :loading="sloading2">
            <el-option
              v-for="token in optionTokens"
              :key="token.id"
              :label="token.name"
              :value="token.id">
            </el-option>
          </el-select>
        </div>
        <button :disabled="!isOperator" type="button" class="btn btn-primary" @click="handleBind">Bind</button>
      </div>
    </div>

    <el-table v-loading.lock="loading" :data="tableData" border style="width: 100%" class="mt20">
      <el-table-column :prop="curTag.name" :label="curTag.name" width="100%">
        <el-table-column type="selection" width="55"> </el-table-column>
        <el-table-column prop="token_name" label="token"> </el-table-column>
        <el-table-column prop="role_name"  label="role"> </el-table-column>
        <el-table-column prop="tag_name"  label="tag"> </el-table-column>
        <el-table-column label="command">
          <template scope="scope">
            <el-button :disabled="!isOperator" @click="unbind(scope.row)" type="danger" size="small">Unbind</el-button>
          </template>
        </el-table-column>
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
// import store from 'src/store'
import { fetch, Msg } from 'src/utils'

export default {
  data () {
    return {
      loading: false,
      sloading1: false,
      sloading2: false,
      global: false,
      roleId: null,
      optionRoles: [],
      tokenId: null,
      optionTokens: [],
      query: '',
      per: 10,
      cur: 1,
      total: 0,
      pageSizes: [5, 10, 20, 50],
      tableData: []
    }
  },
  watch: {
    'curTagId': function (val) {
      this.handleQuery()
    }
  },
  methods: {
    getRoles (query) {
      if (query !== '') {
        this.sloading1 = true
        fetch({
          method: 'get',
          url: 'role/search',
          params: {
            query: query,
            per: 10
          }
        }).then((res) => {
          this.optionRoles = res.data
          this.sloading1 = false
        }).catch((err) => {
          Msg.error('get failed', err)
          this.sloading1 = false
        })
      } else {
        this.optionRoles = []
      }
    },
    getTokens (query) {
      if (query !== '') {
        this.sloading2 = true
        fetch({
          method: 'get',
          url: 'token/search',
          params: {
            query: query,
            per: 10
          }
        }).then((res) => {
          this.optionTokens = res.data
          this.sloading2 = false
        }).catch((err) => {
          Msg.error('get failed', err)
          this.sloading2 = false
        })
      } else {
        this.optionTokens = []
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
      this.loading = true
      fetch({
        method: 'get',
        url: 'rel/tag/role/token/cnt',
        params: { global: this.global, tag_id: this.curTagId, query: this.query }
      }).then((res) => {
        this.total = res.data.total
        this.fetchData()
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    },

    fetchData (opts = {
      global: this.global,
      tag_id: this.curTagId,
      query: this.query,
      mine: this.mine,
      per: this.per,
      offset: this.offset}) {
      fetch({
        method: 'get',
        url: 'rel/tag/role/token/search',
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
        url: 'rel/tag/role/token',
        data: {
          global: this.global,
          tag_id: this.curTagId,
          role_id: this.roleId,
          token_id: this.tokenId
        }
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
    unbind (obj) {
      Msg.confirm('此操作将解绑定该记录, 是否继续?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        console.log(obj)
        this.loading = true
        fetch({
          method: 'delete',
          url: 'rel/tag/role/token',
          data: {
            global: obj.global,
            tag_id: obj.tag_id,
            role_id: obj.role_id,
            token_id: obj.token_id
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
