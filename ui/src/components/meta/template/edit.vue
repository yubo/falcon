<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 v-if="isEdit" class="page-header">edit template</h1>
  <h1 v-else class="page-header">add template </h1>
  <div v-loading.lock="loading">
    <el-form label-position="right" label-width="80px" :model="form">
      <el-form-item label="name">   <el-input v-model="form.template.name"></el-input> </el-form-item>
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
       <el-input v-model="form.action.url"></el-input>
       <el-checkbox v-model.number="bcs">before callback sms</el-checkbox>
       <el-checkbox v-model.number="bcm">before callback mail</el-checkbox>
       <el-checkbox v-model.number="acs">after callback sms</el-checkbox>
       <el-checkbox v-model.number="acm">after callback mail</el-checkbox>
      </el-form-item>
      <el-form-item v-if="isEdit">
        <el-button type="primary" @click="updateTpl">Update</el-button>
        <el-button type="primary" @click="strategyAdd">Add strategy</el-button>
      </el-form-item>
      <el-form-item v-else>
        <el-button type="primary" @click="createtpl">Create</el-button>
      </el-form-item>
    </el-form>
  </div>
  <div v-loading.lock='loading' v-if="isEdit">
    <el-table :data="tableData" border style="width: 100%" class="mt20">
      <el-table-column prop="metricTags" label="metricTags"> </el-table-column>
      <el-table-column prop="fun" label="fun"> </el-table-column>
      <el-table-column prop="maxStep" label="maxStep"> </el-table-column>
      <el-table-column prop="priority" label="priority"> </el-table-column>
      <el-table-column prop="runTime" label="runTime"> </el-table-column>
      <el-table-column label="command">
        <template scope="scope">
          <el-button @click="strategyClone(scope.row.id)" type="text" size="small">CLONE</el-button>
          <el-button @click="strategyEdit(scope.row.id)" type="text" size="small">EDIT</el-button>
          <el-button @click="strategyDelete(scope.row.id)" type="text" size="small">DELETE</el-button>
        </template>
      </el-table-column>
    </el-table>

  </div>
  <div v-loading.lock='loading' v-if="showStrategy" class="mt20">
    <div class="form" role="form">
      <div class="form-group">
        <label for="metric" class="col-sm-2 control-label"> metric: </label>
        <el-select
          style="width: 100%"
          placeholder="metric name"
          v-model="sform.metric"
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
      </div><div class="form-group">
        <label for="tags"> tags:</label><el-input style="width: 100%" v-model="sform.tags"></el-input> 
      </div><div class="form-group">
        <label for="note">note:</label><el-input style="width: 100%" v-model="sform.note"></el-input>
      </div><div class="form-inline mt10" role="form"><div class="form-group">
          触发条件:
          <el-input v-model="sform.fun" style="width:100px;"></el-input>
          <el-select
            style="width:90px;"
            v-model="sform.op"
            filterable>
            <el-option
              v-for="item in optionOps"
              :label="item"
              :value="item">
            </el-option>
          </el-select>
          <el-input v-model="sform.condition" style="width:100px;"> </el-input>
          触发函数举例：all(#2)、sum(#3)、avg(#2)、min(#2)、max(#4)、diff(#5)、pdiff(#5) #后面数字表示最近几个点，不能大于10
      </div></div><div class="form-inline mt10" role="form"><div class="form-group">
          最大报警次数:
          <el-input v-model.number="sform.maxStep"></el-input>
           alarm level:
          <el-select
            style="width:60px;"
            v-model="sform.priority"
            filterable>
            <el-option
              v-for="item in optionPrioritys"
              :label="item"
              :value="item">
            </el-option>
          </el-select>
      </div></div><div class="form-inline mt10" role="form"><div class="form-group">
          生效时间（不指定全天生效）：
          <el-input v-model="sform.runBegin" placeholder="00:00"></el-input> 到
          <el-input v-model="sform.runEnd" placeholder="23:59"></el-input>
      </div></div><div class="form-group" v-if="strategyId">
        <el-button type="primary" @click="updateStrategy">Update</el-button>
      </div><div class="form-group" v-else>
        <el-button type="primary" @click="createStrategy">Create</el-button>
      </div>
    </div>
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
      sloading: false,
      showStrategy: false,
      optionTemplates: [],
      optionUics: [],
      optionMetrics: [],
      optionOps: ['==', '!=', '<', '<=', '>', '>='],
      optionPrioritys: [0, 1, 2, 3, 4, 5],
      tableData: [],
      uics: [],
      bcs: false,
      bcm: false,
      acs: false,
      acm: false,
      templateId: 0,
      strategyId: 0,
      pid: '',
      form: {
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
      sform: {
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
  created () {
    this.templateId = this.$route.query.id
    this.fetchData()
  },
  methods: {
    strategyAdd () {
      this.showStrategy = true
    },
    strategyClone (id) {
      this.showStrategy = true
      this.fetchStrategy(id)
    },
    strategyEdit (id) {
      this.showStrategy = true
      this.strategyId = id
      this.fetchStrategy(id)
    },
    strategyDelete (id) {
    },
    getMetrics (query) {
      if (query !== '') {
        this.sloading = true
        fetch({
          router: this.$router,
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
          Message.error(err.response.data)
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
          router: this.$router,
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
          Message.error(err.response.data)
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
          router: this.$router,
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
          Message.error(err.response.data)
          this.sloading = false
        })
      } else {
        this.optionUics = []
      }
    },
    fetchData () {
      if (!this.templateId) {
        return
      }
      this.loading = true
      fetch({
        router: this.$router,
        method: 'get',
        url: 'template/' + this.templateId
      }).then((res) => {
        for (var k in this.form) {
          if (typeof this.form[k] === 'object') {
            for (var k1 in this.form[k]) {
              if (res.data[k] && res.data[k][k1]) {
                this.form[k][k1] = res.data[k][k1]
              }
            }
          } else {
            this.form[k] = res.data[k]
          }
        }
        if (this.form.action.uic) {
          this.uics = this.form.action.uic.split(',')
          this.optionUics = this.uics.map((v) => {
            return {name: v}
          })
        } else {
          this.uics = []
        }

        if (this.form.template.pid) {
          this.pid = this.form.template.pid.toString()
          this.optionTemplates = [{name: this.form.pname, id: this.pid}]
        } else {
          this.pid = ''
        }
        this.bcs = !!this.form.action.beforeCallbackSms
        this.bcm = !!this.form.action.beforeCallbackMail
        this.acs = !!this.form.action.afterCallbackSms
        this.acm = !!this.form.action.afterCallbackMail
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    updateTpl () {
      this.loading = true
      this.form.action.uic = this.uics.join(',')
      this.form.template.pid = +this.pid
      this.form.action.beforeCallbackSms = +this.bcs
      this.form.action.beforeCallbackMail = +this.bcm
      this.form.action.afterCallbackSms = +this.acs
      this.form.action.afterCallbackMail = +this.acm
      fetch({
        router: this.$router,
        method: 'put',
        url: 'template/' + this.templateId,
        data: JSON.stringify(this.form)
      }).then((res) => {
        Message.success('update success')
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    createtpl () {
      // create
      this.loading = true
      this.form.action.uic = this.uics.join(',')
      this.form.template.pid = +this.pid
      this.form.action.beforeCallbackSms = +this.bcs
      this.form.action.beforeCallbackMail = +this.bcm
      this.form.action.afterCallbackSms = +this.acs
      this.form.action.afterCallbackMail = +this.acm
      fetch({
        router: this.$router,
        method: 'post',
        url: 'template',
        data: JSON.stringify(this.form)
      }).then((res) => {
        Message.success('create success')
        this.templateId = res.data.id
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },

    fetchStrategy (id) {
      this.loading = true
      fetch({
        router: this.$router,
        method: 'get',
        url: 'stragegy/' + id
      }).then((res) => {
        for (var k in this.sform) {
          this.sform[k] = res.data[k]
        }
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    updateStrategy () {
      this.loading = true
      fetch({
        router: this.$router,
        method: 'put',
        url: 'strategy/' + this.stragegyId,
        data: JSON.stringify(this.sform)
      }).then((res) => {
        Message.success('update success')
        this.showStrategy = false
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    createStrategy () {
      // create
      this.loading = true
      this.sform.tplId = this.templateId
      fetch({
        router: this.$router,
        method: 'post',
        url: 'strategy',
        data: JSON.stringify(this.sform)
      }).then((res) => {
        Message.success('create success')
        this.showStrategy = false
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    goBack () {
      this.$router.push(-1)
    }
  },
  computed: {
    isEdit () {
      return this.templateId
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.el-input {
  width: 180px;
}
.el-select {
  width: 280px;
}
</style>
