<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 class="page-header">About</h1>
  <el-tabs v-model="activeName" @tab-click="handleClick" >
    <el-tab-pane label="me" name="me">
    </el-tab-pane>
    <el-tab-pane label="falcon" name="falcon">
    </el-tab-pane>
  </el-tabs>
  <div v-html="compiledMarkdown"></div>
</div>
</template>

<script>
import marked from 'marked'
export default {
  data () {
    return {
      activeName: 'me',
      md: '',
      head: `
 key                  | value 
----------------------| ---
----------------------| ---`,
      falcon: '',
      me: ''
    }
  },
  methods: {
    handleClick (tab, event) {
      console.log(tab.name)
      this.md = this[tab.name]
    }
  },
  computed: {
    compiledMarkdown: function () {
      return marked(this.md, { sanitize: true })
    }
  },
  created () {
    this.me = this.head
    this.me += '\n reader | ' + this.$store.state.auth.reader
    this.me += '\n operator | ' + this.$store.state.auth.operator
    this.me += '\n admin | ' + this.$store.state.auth.admin
    for (let k in this.$store.state.auth.user) {
      this.me += '\n' + k + ' | ' + this.$store.state.auth.user[k]
    }

    this.falcon = this.head

    this.md = this[this.activeName]
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
