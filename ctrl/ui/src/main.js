// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'

import app from './components/app'
import store from './store'
import router from './router'
import {
  Tree,
  Loading,
  Form,
  FormItem,
  Input,
  Button,
  Table,
  TableColumn,
  Pagination,
  Select,
  Option,
  Switch,
  RadioGroup,
  Radio,
  Checkbox,
  Tabs,
  TabPane,
  Dialog,
  DatePicker
} from 'element-ui'

/* eslint-disable no-new */
Vue.use(Tree)
Vue.use(Loading)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Input)
Vue.use(Button)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Pagination)
Vue.use(Select)
Vue.use(Option)
Vue.use(Switch)
Vue.use(RadioGroup)
Vue.use(Radio)
Vue.use(Checkbox)
Vue.use(Tabs)
Vue.use(TabPane)
Vue.use(Dialog)
Vue.use(DatePicker)

new Vue({
  el: '#app',
  template: '<app/>',
  store,
  router,
  components: { app }
})
