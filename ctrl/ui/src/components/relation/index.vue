<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-sm-3">
        <div class="pull-right">
          <button type="button" @click="reload" class="btn btn-default"><span class="glyphicon glyphicon-refresh"></span></button>
        </div>
        <el-tree v-loading="loading"
          :data="tagTree"
          :props="props"
          :indent="8"
          :highlight-current="true"
          :expand-on-click-node="false"
          @current-change="handleCurrentChange">
        </el-tree>
      </div>
      <div class="col-sm-9">
        <router-view> </router-view>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      props: {
        label: 'label',
        children: 'child'
      }
    }
  },
  methods: {
    handleCurrentChange (val) {
      this.$store.commit('rel/m_cur_tag', val)
    },
    reload () {
      this.$store.dispatch('rel/load_tree')
    }
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
      this.$store.dispatch('rel/load_tree')
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
