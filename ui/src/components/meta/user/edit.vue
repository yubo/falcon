<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 v-if="isEdit" class="page-header">edit user</h1>
  <h1 v-else class="page-header">add user </h1>
  <div v-loading.lock="loading">
    <el-form label-position="right" label-width="80px" :model="userform">
      <el-form-item label="name">  <el-input v-model="userform.name"> </el-input> </el-form-item>
      <el-form-item label="uuid">  <el-input v-model="userform.uuid"> </el-input> </el-form-item>
      <el-form-item label="cname"> <el-input v-model="userform.cname"></el-input> </el-form-item>
      <el-form-item label="email"> <el-input v-model="userform.email"></el-input> </el-form-item>
      <el-form-item label="phone"> <el-input v-model="userform.phone"></el-input> </el-form-item>
      <el-form-item label="im">    <el-input v-model="userform.im">   </el-input> </el-form-item>
      <el-form-item label="qq">    <el-input v-model="userform.qq">   </el-input> </el-form-item>
      <el-form-item v-if="isEdit">
        <el-button type="primary" @click="putData">Update</el-button>
        <el-button @click="fetchData">Reset</el-button>
      </el-form-item> <el-form-item v-else>
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
      userform: {
        name: '',
        uuid: '',
        cname: '',
        email: '',
        phone: '',
        im: '',
        qq: ''
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
        url: 'user/' + this.$route.query.id
      }).then((res) => {
        for (var k in this.userform) {
          this.userform[k] = res.data[k]
        }
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
        router: this.$router,
        method: 'put',
        url: 'user/' + this.$route.query.id,
        data: JSON.stringify(this.userform)
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
        url: 'user',
        data: JSON.stringify(this.userform)
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
