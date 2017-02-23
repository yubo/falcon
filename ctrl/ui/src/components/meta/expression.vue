<template>
  <div id="content" class="main">
    <div class="form-inline" role="form">
      <div class="form-group">
        <input type="text" v-model="query" @keyup.enter="handleQuery" class="form-control" placeholder="expression name">
      </div>
      <button type="button" @click="handleQuery" class="btn btn-primary">search</button>
      <input type="checkbox" v-model="mine" class="form-control">
      <span>mine</span>
      <div class="pull-right">
        <button :disabled="!isOperator" type="button" @click="editObj" class="btn btn-default"><span class="glyphicon glyphicon-plus"></span>Add</button>
      </div>
    </div>

    <el-table v-loading="loading" :data="tableData" border style="width: 100%" class="mt20">
      <el-table-column prop="expression"   label="expression"> </el-table-column>
      <el-table-column prop="cfun"   label="fun"> </el-table-column>
      <el-table-column prop="uic"   label="uic"> </el-table-column>
      <el-table-column prop="note"   label="note"> </el-table-column>
      <el-table-column prop="creator"  label="creator"> </el-table-column>
      <el-table-column label="command">
        <template scope="scope">
          <el-button :disabled="!isOperator" @click="pauseObj(scope)" type="text" size="small" >{{scope.row.pause ? 'ACTIVE' : 'PAUSE'}}</el-button>
          <el-button :disabled="!isOperator" @click="editObj(scope.row)" type="text" size="small">EDIT</el-button>
          <el-button :disabled="!isOperator" @click="deleteObj(scope.row)" type="text" size="small">DELETE</el-button>
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

    <el-dialog size="large" :title="isEdit ? 'edit expression' : 'add expression'" v-model="editVisible" :close-on-click-modal="false">
      <div v-loading.lock="dloading">
        <el-form label-position="right" label-width="80px" :model="objForm">
          <el-form-item label="expression"><el-input v-model="objForm.expression.expression"></el-input> </el-form-item>
          <el-form-item label="disabled"><el-switch on-text="" off-text="" v-model="pause"></el-input> </el-form-item>
          <el-form-item label="触发条件">
              <el-input v-model="objForm.expression.fun" style="width:100px;"></el-input>
              <el-select
                style="width:90px;"
                v-model="objForm.expression.op"
                filterable>
                <el-option
                  v-for="item in optionOps"
                  :label="item"
                  :value="item">
                </el-option>
              </el-select>
              <el-input v-model="objForm.expression.condition" style="width:100px;"> </el-input>
              :alarm(); callback();
          </el-form-item>
          <el-form-item label="uic">
            <el-select
              style="width: 100%"
              placeholder="uic"
              v-model="uics"
              multiple
              filterable
              remote
              :remote-method="getUics"
              :loading="sloading">
              <el-option
                v-for="item in optionUics"
                :label="item.name"
                :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="最大报警次数">
              <el-input v-model.number="objForm.action.maxStep" style="width:100px;"></el-input> alarm level:
              <el-select
                style="width:60px;"
                v-model="objForm.expression.priority"
                filterable>
                <el-option
                  v-for="item in optionPrioritys"
                  :label="item"
                  :value="item">
                </el-option>
              </el-select>
          </el-form-item>
          <el-form-item label="callback">
           <el-input v-model="objForm.action.url"></el-input>
           <el-checkbox v-model.number="bcs">before callback sms</el-checkbox>
           <el-checkbox v-model.number="bcm">before callback mail</el-checkbox>
           <el-checkbox v-model.number="acs">after callback sms</el-checkbox>
           <el-checkbox v-model.number="acm">after callback mail</el-checkbox>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitObj">submit</el-button>
            <el-button @click="editVisible = false">Cancel</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { fetch, Msg } from 'src/utils'
export default {
  data () {
    return {
      loading: false,
      dloading: false,
      sloading: false,
      query: '',
      mine: true,
      per: 10,
      cur: 1,
      total: 0,
      pageSizes: [5, 10, 20, 50],
      tableData: [],
      curId: 0,
      editVisible: false,
      optionUics: [],
      optionOps: ['==', '!=', '<', '<=', '>', '>='],
      optionPrioritys: [0, 1, 2, 3, 4, 5],
      uics: [],
      bcs: false,
      bcm: false,
      acs: false,
      acm: false,
      pause: false,
      objForm: {
        expression: {
          expression: '',
          op: '',
          condition: '',
          maxStep: '',
          priority: '',
          msg: '',
          fun: '',
          pause: 0
        },
        action: {
          uic: '',
          url: '',
          sendsms: 0,
          sendmail: 0,
          beforeCallbackSms: 0,
          beforeCallbackMail: 0,
          afterCallbackSms: 0,
          afterCallbackMail: 0
        }
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
      this.reFetchObjs()
    },

    getUics (query) {
      if (query !== '') {
        this.sloading = true
        fetch({
          method: 'get',
          url: 'team/search',
          params: {
            query: query,
            per: 10
          }
        }).then((res) => {
          this.optionUics = res.data
          this.sloading = false
        }).catch((err) => {
          Msg.error('get failed', err)
          this.sloading = false
        })
      } else {
        this.optionUics = []
      }
    },

    reFetchObjs () {
      fetch({
        method: 'get',
        url: 'expression/cnt',
        params: {query: this.query, mint: this.mine}
      }).then((res) => {
        this.total = res.data.total
        this.fetchObjs()
      }).catch((err) => {
        Msg.error('get failed', err)
      })
    },

    fetchObjs (opts = {query: this.query, mine: this.mine, per: this.per, offset: this.offset}) {
      this.loading = true
      fetch({
        method: 'get',
        url: 'expression/search',
        params: opts
      }).then((res) => {
        this.tableData = res.data ? res.data.map((v) => {
          return {
            id: v.id,
            expression: v.expression,
            cfun: v.fun + v.op + v.condition,
            uic: v.uic,
            note: 'Max(' + v.max_step + '), P(' + v.priority + ')',
            pause: v.pause,
            creator: v.creator
          }
        }) : []
        this.loading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    },
    // set expression/action from obj
    setObj (obj = {expression: {}, action: {}}) {
      for (var k in this.objForm) {
        if (typeof this.objForm[k] === 'object') {
          for (var k1 in this.objForm[k]) {
            this.objForm[k][k1] = obj[k][k1]
          }
        } else {
          this.objForm[k] = obj[k]
        }
      }
    },
    // fetch expression/action from v1.0/expression/:id
    // fetch strategys from v1.0/strategy/search?tid=?&query=?&per=...
    fetchObj () {
      if (!this.curId) {
        return
      }
      this.dloading = true
      fetch({
        method: 'get',
        url: 'expression/' + this.curId
      }).then((res) => {
        this.setObj(res.data)
        if (this.objForm.action.uic) {
          this.uics = this.objForm.action.uic.split(',')
          this.optionUics = this.uics.map((v) => {
            return {name: v}
          })
        } else {
          this.uics = []
        }
        this.bcs = !!this.objForm.action.beforeCallbackSms
        this.bcm = !!this.objForm.action.beforeCallbackMail
        this.acs = !!this.objForm.action.afterCallbackSms
        this.acm = !!this.objForm.action.afterCallbackMail
        this.pause = !!this.objForm.expression.pause
        this.dloading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.dloading = false
        this.editVisible = false
      })
    },

    pauseObj (scope) {
      console.log(scope)
      let status = +!scope.row.pause
      fetch({
        method: 'put',
        url: 'expression/pause',
        params: {id: scope.row.id, pause: status}
      }).then((res) => {
        Msg.success('update success')
        this.tableData[scope.$index].pause = status
      }).catch((err) => {
        Msg.error('update failed', err)
      })
    },
    editObj (obj = {}) {
      this.editVisible = true
      this.curId = obj.id
      if (!obj.id) {
        this.setObj()
        return
      }
      this.fetchObj()
    },
    submitObj () {
      this.dloading = true
      this.objForm.action.uic = this.uics.join(',')
      this.objForm.action.beforeCallbackSms = +this.bcs
      this.objForm.action.beforeCallbackMail = +this.bcm
      this.objForm.action.afterCallbackSms = +this.acs
      this.objForm.action.afterCallbackMail = +this.acm
      this.objForm.expression.pause = +this.pause
      fetch({
        method: this.isEdit ? 'put' : 'post',
        url: this.isEdit ? 'expression/' + this.curId : 'expression',
        data: this.objForm
      }).then((res) => {
        Msg.success('update success')
        if (!this.isEdit) {
          this.total++
        }
        this.fetchObjs()
        this.dloading = false
        this.editVisible = false
      }).catch((err) => {
        Msg.error('update failed', err)
        this.dloading = false
      })
    },
    deleteObj (obj) {
      Msg.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        fetch({
          method: 'delete',
          url: 'expression/' + obj.id
        }).then((res) => {
          Msg.success('success!')
          this.total--
          this.fetchObjs()
        }).catch((err) => {
          Msg.error('delete failed', err)
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
    isEdit () {
      return this.curId
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
    this.reFetchObjs()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
