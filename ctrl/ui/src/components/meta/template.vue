<template>
  <div id="content" class="main">
    <div class="form-inline" role="objForm">
      <div class="form-group">
        <input type="text" v-model="query" @keyup.enter="handleQuery" class="form-control" placeholder="template name">
      </div>
      <button type="button" @click="handleQuery" class="btn btn-primary">search</button>
      <input type="checkbox" v-model="mine" class="form-control">
      <span>mine</span>
      <div class="pull-right">
        <button :disabled="!isOperator" type="button" @click="editObj" class="btn btn-default"><span class="glyphicon glyphicon-plus"></span>Add</button>
      </div>
    </div>

    <el-table v-loading="loading" :data="tableData" border style="width: 100%" class="mt20">
      <el-table-column prop="name"    label="name"> </el-table-column>
      <el-table-column prop="pname"   label="pname"> </el-table-column>
      <el-table-column prop="creator" label="creator"> </el-table-column>
      <el-table-column label="command">
        <template scope="scope">
          <el-button :disabled="!isOperator" @click="editObj(scope.row, true)" type="text" size="small">CLONE</el-button>
          <el-button :disabled="!isOperator" @click="editObj(scope.row)" type="text" size="small">EDIT</el-button>
          <el-button :disabled="!isOperator" @click="bindObj(scope.row)" type="text" size="small">BIND</el-button>
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


    <!-- edit modal -->
    <el-dialog size="large" :title="isEdit ? 'edit template' : 'add template'" v-model="editVisible" :close-on-click-modal="false">
      <div v-loading.lock="dloading">
        <el-form label-position="right" label-width="80px" :model="objForm">

          <el-form-item label="name"><el-input v-model="objForm.template.name"></el-input> </el-form-item>
          <el-form-item label="parent template"> 
            <el-select
              style="width: 100%"
              placeholder="parent template"
              v-model="pid"
              filterable
              remote
              :remote-method="getTemplates"
              :loading="sloading">
              <el-option
                v-for="item in optionTemplates"
                :key="item.id"
                :label="item.name"
                :value="item.id">
              </el-option>
            </el-select>
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
          <el-form-item label="callback">
           <el-input v-model="objForm.action.url"></el-input>
           <el-checkbox v-model.number="bcs">before callback sms</el-checkbox>
           <el-checkbox v-model.number="bcm">before callback mail</el-checkbox>
           <el-checkbox v-model.number="acs">after callback sms</el-checkbox>
           <el-checkbox v-model.number="acm">after callback mail</el-checkbox>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="submitObj">{{isEdit ? 'update' : 'create'}}</el-button>
            <el-button v-if="isEdit" type="primary" @click="editObj2">add strategy</el-button>
            <el-button @click="editVisible = false">Cancel</el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- strategy table -->
      <div v-loading.lock='tloading' v-if="isEdit">
        <el-table :data="tableData2" border style="width: 100%" class="mt20">
          <el-table-column prop="metricTags" label="metricTags"> </el-table-column>
          <el-table-column prop="cFun" label="fun"> </el-table-column>
          <el-table-column prop="maxStep" label="maxStep"> </el-table-column>
          <el-table-column prop="priority" label="priority"> </el-table-column>
          <el-table-column prop="cRunTime" label="runTime"> </el-table-column>
          <el-table-column label="command">
            <template scope="scope">
              <el-button @click="editObj2(scope.row, false)" type="text" size="small">CLONE</el-button>
              <el-button @click="editObj2(scope.row)" type="text" size="small">EDIT</el-button>
              <el-button @click="deleteObj2(scope.row)" type="text" size="small">DELETE</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pull-right mt20">
          <el-pagination
            @size-change="sizeChange2"
            @current-change="curChange2"
            :current-page="cur2"
            :page-sizes="pageSizes2"
            :page-size="per2"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total2">
          </el-pagination>
        </div>

      </div>
    </el-dialog>

    <!-- sub modal -->
    <el-dialog :title="isEdit2 ? 'edit strategy' : 'add strategy'" v-model="editVisible2" :close-on-click-modal="false">
      <div v-loading.lock="dloading">
        <el-form label-position="right" label-width="100px" :model="objForm2">
          <el-form-item label="metric">
            <el-select
              style="width: 100%"
              placeholder="metric name"
              v-model="objForm2.metric"
              filterable
              remote
              :remote-method="getMetrics"
              :loading="sloading">
              <el-option
                v-for="item in optionMetrics"
                :label="item.name"
                :value="item.name">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="tags"><el-input v-model="objForm2.tags"></el-input></el-form-item>
          <el-form-item label="note"><el-input v-model="objForm2.note"></el-input></el-form-item>
          <el-form-item label="触发条件">
              <el-input v-model="objForm2.fun" style="width:100px;"></el-input>
              <el-select
                style="width:90px;"
                v-model="objForm2.op"
                filterable>
                <el-option
                  v-for="item in optionOps"
                  :label="item"
                  :value="item">
                </el-option>
              </el-select>
              <el-input v-model="objForm2.condition" style="width:100px;"> </el-input>
              触发函数举例：all(#2)、sum(#3)、avg(#2)、min(#2)、max(#4)、diff(#5)、pdiff(#5) #后面数字表示最近几个点，不能大于10
          </el-form-item>

          <el-form-item label="最大报警次数">
              <el-input v-model.number="objForm2.maxStep" style="width:100px;"></el-input> alarm level:
              <el-select
                style="width:60px;"
                v-model="objForm2.priority"
                filterable>
                <el-option
                  v-for="item in optionPrioritys"
                  :label="item"
                  :value="item">
                </el-option>
              </el-select>
          </el-form-item>
          <el-form-item label="生效时间">
              <el-input v-model="objForm2.runBegin" placeholder="00:00" style="width:100px;"></el-input> 到
              <el-input v-model="objForm2.runEnd" placeholder="23:59" style="width:100px;"></el-input>（不指定全天生效）
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitObj2">{{isEdit2 ? 'update' : 'create'}}</el-button>
            <el-button @click="editVisible2 = false">Cancel</el-button>
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
      tloading: false,
      query: '',
      query2: '',
      mine: true,
      per: 10,
      per2: 10,
      cur: 1,
      cur2: 1,
      total: 0,
      total2: 0,
      pageSizes: [5, 10, 20, 50],
      pageSizes2: [5, 10, 20, 50],
      tableData: [],
      tableData2: [],
      curId: 0,
      curId2: 0,
      editVisible: false,
      editVisible2: false,
      optionTemplates: [],
      optionUics: [],
      optionMetrics: [],
      optionOps: ['==', '!=', '<', '<=', '>', '>='],
      optionPrioritys: [0, 1, 2, 3, 4, 5],
      uics: [],
      bcs: false,
      bcm: false,
      acs: false,
      acm: false,
      pid: '',
      objForm: {
        template: {
          name: '',
          pid: 0
        },
        action: {
          uic: '',
          url: '',
          sendsms: 0,
          sendmail: 0,
          callback: 0,
          beforeCallbackSms: 0,
          beforeCallbackMail: 0,
          afterCallbackSms: 0,
          afterCallbackMail: 0
        },
        pname: ''
      },
      objForm2: {
        id: 0,
        tags: '',
        maxStep: 0,
        priority: 0,
        fun: '',
        op: '',
        condition: '',
        note: '',
        metric: '',
        runBegin: '',
        runEnd: '',
        tplId: 0
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
    sizeChange2 (per) {
      this.per2 = per
      this.fetchObjs2()
    },
    curChange2 (cur) {
      this.cur2 = cur
      this.fetchObjs2()
    },

    handleQuery () {
      this.reFetchObjs()
    },

    getMetrics (query) {
      if (query !== '') {
        this.sloading = true
        fetch({
          method: 'get',
          url: 'metric/search',
          params: {
            query: query,
            per: 10
          }
        }).then((res) => {
          this.optionMetrics = res.data
          this.sloading = false
        }).catch((err) => {
          Msg.error('get failed', err)
          this.sloading = false
        })
      } else {
        this.optionMetrics = []
      }
    },
    getTemplates (query) {
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
          this.optionTemplates = res.data
          this.sloading = false
        }).catch((err) => {
          Msg.error('get failed', err)
          this.sloading = false
        })
      } else {
        this.optionTemplates = []
      }
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
        url: 'template/cnt',
        params: {query: this.query}
      }).then((res) => {
        this.total = res.data.total
        this.fetchObjs()
      }).catch((err) => {
        Msg.error('get failed', err)
      })
    },
    fetchObjs (opts = {query: this.query, per: this.per, offset: this.offset}) {
      this.loading = true
      fetch({
        method: 'get',
        url: 'template/search',
        params: opts
      }).then((res) => {
        this.tableData = res.data
        this.loading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.loading = false
      })
    },
    // set template/action from obj
    setObj (obj = {template: {}, action: {}}) {
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
    // fetch template/action from v1.0/template/:id
    // fetch strategys from v1.0/strategy/search?tid=?&query=?&per=...
    fetchObj (clone = false) {
      if (!this.curId) {
        return
      }
      this.dloading = true
      fetch({
        method: 'get',
        url: 'template/' + this.curId,
        params: { clone: clone }
      }).then((res) => {
        this.setObj(res.data)
        if (clone) {
          this.total++
          this.curId = res.data.template.id
          this.fetchObjs()
        }
        if (this.objForm.action.uic) {
          this.uics = this.objForm.action.uic.split(',')
          this.optionUics = this.uics.map((v) => {
            return {name: v}
          })
        } else {
          this.uics = []
        }

        if (this.objForm.template.pid) {
          this.pid = this.objForm.template.pid.toString()
          this.optionTemplates = [{name: this.objForm.pname, id: this.pid}]
        } else {
          this.pid = ''
        }
        this.bcs = !!this.objForm.action.beforeCallbackSms
        this.bcm = !!this.objForm.action.beforeCallbackMail
        this.acs = !!this.objForm.action.afterCallbackSms
        this.acm = !!this.objForm.action.afterCallbackMail
        this.dloading = false

        this.reFetchObjs2()
      }).catch((err) => {
        Msg.error('get failed', err)
        this.dloading = false
        this.editVisible = false
      })
    },
    editObj (obj = {}, clone = false) {
      this.editVisible = true
      if (!obj.id) {
        this.setObj()
        return
      }
      this.curId = obj.id
      this.fetchObj(clone)
    },
    bindObj (obj = {}) {
    },
    submitObj () {
      this.dloading = true
      this.objForm.action.uic = this.uics.join(',')
      this.objForm.template.pid = +this.pid
      this.objForm.action.beforeCallbackSms = +this.bcs
      this.objForm.action.beforeCallbackMail = +this.bcm
      this.objForm.action.afterCallbackSms = +this.acs
      this.objForm.action.afterCallbackMail = +this.acm
      fetch({
        method: this.isEdit ? 'put' : 'post',
        url: this.isEdit ? 'template/' + this.curId : 'template',
        data: this.objForm
      }).then((res) => {
        Msg.success('submit success')
        if (!this.isEdit) {
          this.curId = res.data.id
          this.total++
        }
        this.fetchObjs()
        this.dloading = false
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
          url: 'template/' + obj.id
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
    },

    // strategy
    editObj2 (obj = {}, edit = true) {
      if (edit) {
        this.curId2 = obj.id
      }
      for (var k in this.objForm2) {
        this.objForm2[k] = obj[k]
      }
      this.editVisible2 = true
      console.log('obj', obj, 'objForm2', this.objForm2)
    },
    submitObj2 () {
      // create
      this.dloading = true
      this.objForm2.tplId = this.curId
      fetch({
        method: this.isEdit2 ? 'put' : 'post',
        url: this.isEdit2 ? 'strategy/' + this.curId2 : 'strategy',
        data: this.objForm2
      }).then((res) => {
        Msg.success('submit success')
        if (!this.isEdit2) {
          this.total2++
        }
        this.fetchObjs2()
        this.dloading = false
        this.editVisible2 = false
      }).catch((err) => {
        Msg.error('update failed', err)
        this.dloading = false
      })
    },
    deleteObj2 (obj) {
      Msg.confirm('此操作将永久删除该记录, 是否继续?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        fetch({
          method: 'delete',
          url: 'strategy/' + obj.id
        }).then((res) => {
          Msg.success('success!')
          this.total2--
          this.fetchObjs2()
        }).catch((err) => {
          Msg.error('delete failed', err)
        })
      }).catch(() => {
        Msg.info('cancel')
      })
    },
    // fetch strategys from v1.0/strategy/search?tid=?&query=?&per=...
    reFetchObjs2 () {
      this.dloading = true
      fetch({
        method: 'get',
        url: 'strategy/cnt',
        params: {tid: this.curId, query: this.query}
      }).then((res) => {
        this.total2 = res.data.total
        this.dloading = false
        this.fetchObjs2()
      }).catch((err) => {
        Msg.error('get failed', err)
        this.dloading = false
      })
    },
    fetchObjs2 (opts = {tid: this.curId, query: this.query2, per: this.per2, offset: this.offset2}) {
      if (!this.curId) {
        return
      }
      this.tloading = true
      fetch({
        method: 'get',
        url: 'strategy/search',
        params: opts
      }).then((res) => {
        this.tableData2 = res.data.map((v) => {
          var ret = v
          ret.metricTags = v.metric + '/' + v.tags
          ret.cFun = v.fun + ' ' + v.op + ' ' + v.condition
          ret.cRunTime = v.runBegin + ' ~ ' + v.runEnd
          return ret
        })
        this.tloading = false
      }).catch((err) => {
        Msg.error('get failed', err)
        this.tloading = false
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
    offset2 () {
      return (this.per2 * (this.cur2 - 1))
    },
    isEdit () {
      return this.curId
    },
    isEdit2 () {
      return this.curId2
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
