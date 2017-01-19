import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from './components/Login.vue'

Vue.use(VueRouter)

const routes = [
  { path: '/login/', name: 'login', component: Login }
]

export default new VueRouter({
  routes
})
