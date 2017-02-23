<template>
<div id="content" class="main">
  <div class="form-group">
    <button type="button" @click="debug('populate')" class="btn btn-primary">populate</button>
  </div>
  <div class="form-group">
    <button type="button" @click="debug('reset_db')" class="btn btn-primary">reset db</button>
  </div>
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
        method: 'get',
        url: '/admin/debug/' + action
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
