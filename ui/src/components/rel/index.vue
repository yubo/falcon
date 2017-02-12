<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-sm-3 col-md-3 sidebar">
        <el-tree v-loading="loading"
          :data="tagTree"
          :props="props"
          :highlight-current="true"
          @current-change="handleCurrentChange">
        </el-tree>
      </div>
      <div class="col-sm-9 col-sm-offset-3 col-md-9 col-md-offset-3 main">
        <ul class="nav nav-pills mt0">
          <li is="li-tpl" v-for="(obj, li_idx) in links" :obj="obj"></li>
        </ul>
        <router-view> </router-view>
      </div>
  </div>
</template>

<script>
import { liTpl } from '../tpl'

export default {
  data () {
    return {
      links: [
      { url: '/rel/tag-host', text: 'host' },
      { url: '/rel/tag-template', text: 'template' },
      { url: '/rel/tag-role-user', text: 'role user' },
      { url: '/rel/tag-role-token', text: 'role token' }
      ],
      props: {
        label: 'label',
        children: 'child'
      }
    }
  },
  methods: {
    handleCurrentChange (val) {
      this.$store.commit('rel/m_cur_tag', val)
    }
  },
  components: {
    liTpl
  },
  computed: {
    loading () {
      return this.$store.state.rel.loading
    },
    tagTree () {
      return this.$store.state.rel.tree
    }
  },
  created () {
    if (!this.$store.state.rel.loaded) {
      this.$store.commit('rel/m_load_tag')
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.sidebar {
  padding: 0px;
}
</style>
