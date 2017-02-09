<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 v-if="isEdit" class="page-header">edit role</h1>
  <h1 v-else class="page-header">add role </h1>
  <div v-loading.lock="loading">
    <el-form label-position="right" label-width="80px" :model="roleform">
      <el-form-item label="name">   <el-input v-model="roleform.name"></el-input> </el-form-item>
      <el-form-item label="cname">  <el-input v-model="roleform.cname"></el-input> </el-form-item>
      <el-form-item label="note">   <el-input v-model="roleform.note"></el-input> </el-form-item>
      <el-form-item v-if="isEdit">
        <el-button type="primary" @click="putData">Update</el-button>
        <el-button @click="fetchData">Reset</el-button>
      </el-form-item>
      <el-form-item v-else>
        <el-button type="primary" @click="postData">Create</el-button>
        <el-button @click="goBack">Cancel</el-button>
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
      roleform: {
        name: '',
        cname: '',
        note: ''
      }
    }
  },
  created () {
    this.fetchData()
  },
  methods: {
    fetchData () {
      if (!this.$route.query.id) {
        return
      }
      this.loading = true
      fetch({
        router: this.$router,
        method: 'get',
        url: 'role/' + this.$route.query.id
      }).then((res) => {
        for (var k in this.roleform) {
          this.roleform[k] = res.data[k]
        }
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    putData () {
      console.log('submit!')
      this.loading = true
      // update
      fetch({
        router: this.$router,
        method: 'put',
        url: 'role/' + this.$route.query.id,
        data: JSON.stringify(this.roleform)
      }).then((res) => {
        Message.success('update success')
        this.loading = false
      }).catch((err) => {
        Message.error(err.response.data)
        this.loading = false
      })
    },
    postData () {
      // create
      fetch({
        router: this.$router,
        method: 'post',
        url: 'role',
        data: JSON.stringify(this.roleform)
      }).then((res) => {
        Message.success('update success')
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
      return this.$route.query.id
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
