<template>
  <div v-show="alertShow" v-bind:class="alertClass">
    <button type="button" class="close" @click="alertClose"><span aria-hidden="true">&times;</span></button>
    <strong>{{ alertLevel }}</strong> {{ alertMsg }}
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  methods: {
    alertClose () {
      this.$store.dispatch('notification_close')
    }
  },
  computed: {
    alertClass () {
      return {
        alert: true,
        'alert-success': this.$store.state.notification.level === 'SUCCESS',
        'alert-info': this.$store.state.notification.level === 'INFO',
        'alert-warning': this.$store.state.notification.level === 'WARNING',
        'alert-danger': this.$store.state.notification.level === 'DANGER' ||
          this.$store.state.notification.level === 'ERROR'
      }
    },
    ...mapState({
      alertLevel: state => state.notification.level,
      alertMsg: state => state.notification.msg,
      alertShow: state => state.notification.show
    })
  },
  data () {
    return {
    }
  }
}
</script>

<style scoped>

</style>
