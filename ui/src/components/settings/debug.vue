<template>
<div v-loading.lock="loading" id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 class="page-header">debug</h1>
    <div class="form-group">
    <button type="button" @click="debug('populate')" class="btn btn-primary">populate</button>
    </div>
    <div class="form-group">
    <button type="button" @click="debug('reset_db')" class="btn btn-primary">reset db</button>
</div>
</template>

<script>
import { fetch } from 'src/utils'
import { Notification } from 'element-ui'
export default {
  data () {
    return {
      loading: false
    }
  },
  methods: {
    debug (action) {
      fetch({
        router: this.$router,
        method: 'get',
        url: '/settings/debug/' + action
      }).then((res) => {
        Notification.success({title: 'Success', message: res.data})
        this.loading = false
      }).catch((err) => {
        Notification.error({title: 'Error', message: err.response.data})
        this.loading = false
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
