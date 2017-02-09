<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 v-if="isEdit" class="page-header">edit team</h1>
  <h1 v-else class="page-header">add team </h1>
  <div v-loading.lock="loading">
    <el-form label-position="right" label-width="80px" :model="teamform">
      <el-form-item label="name">   <el-input v-model="teamform.name"></el-input> </el-form-item>
      <el-form-item label="note">   <el-input v-model="teamform.note"></el-input> </el-form-item>
      <el-form-item label="member">
        <el-select
          style="width: 100%"
          placeholder="user name"
          v-model="users"
          multiple
          filterable
          remote
          :remote-method="getUsers"
          :loading="sloading">
          <el-option
            v-for="user in optionUsers"
            :key="user.id"
            :label="user.name"
            :value="user.id">
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item v-if="isEdit">
        <el-button type="primary" @click="updateData">Update</el-button>
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
      isEdit: false,
      editId: 0,
      loading: false,
      sloading: false,
      users: [],
      optionUsers: [],
      teamform: {
        name: '',
        note: ''
      }
    }
  },
  methods: {
    getUsers (query) {
      if (query !== '') {
        this.sloading = true
        fetch({
          router: this.$router,
          method: 'get',
          url: 'user/search',
          params: {
            query: query,
            per: 10
          }
        }).then((res) => {
          this.optionUsers = res.data
          this.sloading = false
        }).catch((err) => {
          Message.error(err.response.data)
          this.sloading = false
        })
      } else {
        this.optionUsers = []
      }
    },
    fetchData () {
      this.loading = true
      fetch({
        router: this.$router,
        method: 'get',
        url: 'team/' + this.editId
      }).then((res) => {
        for (var k in this.teamform) {
          this.teamform[k] = res.data[k]
        }
        this.loading = false
        this.fetchMember()
      }).catch((err) => {
        Message.error(err)
        this.loading = false
      })
    },
    fetchMember () {
      this.loading = true
      fetch({
        router: this.$router,
        method: 'get',
        url: 'team/' + this.editId + '/member'
      }).then((res) => {
        this.optionUsers = res.data.users
        this.users = this.optionUsers.map((i) => {
          return i.id
        })
        this.loading = false
      }).catch((err) => {
        Message.error(err)
        this.loading = false
      })
    },
    updateData () {
      console.log('submit!')
      this.loading = true
      // update
      fetch({
        router: this.$router,
        method: 'put',
        url: 'team/' + this.editId,
        data: JSON.stringify(this.teamform)
      }).then((res) => {
        Message.success('update success')
        this.loading = false
        this.updateMember()
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
        url: 'team',
        data: JSON.stringify(this.teamform)
      }).then((res) => {
        this.editId = res.data.id
        this.isEdit = true
        this.loading = false
        this.updateMember()
      }).catch((err) => {
        Message.error(err)
        this.loading = false
      })
    },
    updateMember () {
      this.loading = true
      fetch({
        router: this.$router,
        method: 'put',
        url: 'team/' + this.editId + '/member',
        data: JSON.stringify({uids: this.users})
      }).then((res) => {
        Message.success('update success')
        this.loading = false
      }).catch((err) => {
        console.log(err)
        Message.error(err)
        this.loading = false
      })
    },
    goBack () {
      console.log(this.users)
      // this.$router.push(-1)
    }
  },
  created () {
    if (this.$route.query.id) {
      this.editId = this.$route.query.id
      this.isEdit = true
      this.fetchData()
    }
  },
  computed: {
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
