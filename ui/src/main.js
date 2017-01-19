// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'

import app from './components/app'
import store from './store'
import router from './router'

router.push('login')

/* eslint-disable no-new */
new Vue({
  el: '#app',
  template: '<app/>',
  store,
  router,
  components: { app }
})
