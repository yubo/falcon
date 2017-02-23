<template>
<div id="content" class="main">
  <div v-loading.lock="loading">
    <el-form label-position="right" label-width="80px" :model="userform">
      <el-form-item label="name">  <el-input :disabled="hasName" v-model="userform.name"> </el-input> </el-form-item>
      <el-form-item label="uuid">  <el-input :disabled="true" v-model="userform.uuid"> </el-input> </el-form-item>
      <el-form-item label="cname"> <el-input v-model="userform.cname"></el-input> </el-form-item>
      <el-form-item label="email"> <el-input v-model="userform.email"></el-input> </el-form-item>
      <el-form-item label="phone"> <el-input v-model="userform.phone"></el-input> </el-form-item>
      <el-form-item label="im">    <el-input v-model="userform.im">   </el-input> </el-form-item>
      <el-form-item label="qq">    <el-input v-model="userform.qq">   </el-input> </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="putData">Update</el-button>
        <el-button @click="fetchData">Reset</el-button>
      </el-form-item>
    </el-form>
  </div>
</div>
</template>

<script>
import { fetch, Msg } from 'src/utils'
export default {
  data () {
    return {
      loading: false,
      hasName: true,
      userform: {
        id: 0,
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
      for (var k in this.userform) {
        this.userform[k] = this.$store.state.auth.user[k]
      }
      this.hasName = (this.userform.name !== '')
    },
    putData () {
      this.loading = true
      console.log(this.userform)
      // update
      fetch({
        method: 'put',
        url: 'settings/profile',
        data: this.userform
      }).then((res) => {
        Msg.success('update success')
        this.$store.commit('auth/m_set_user', res.data)
        this.loading = false
      }).catch((err) => {
        Msg.error('update failed', err)
        this.loading = false
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
