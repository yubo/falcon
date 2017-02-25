<template>
<div id="content" class="main">
  <el-tabs v-model="activeName" @tab-click="handleClick" >
    <el-tab-pane label="me" name="me">
    </el-tab-pane>
    <el-tab-pane label="falcon" name="falcon">
    </el-tab-pane>
  </el-tabs>
  <div class="pull-right">
    <button type="button" @click="getinfo" class="btn btn-default"><span class="glyphicon glyphicon-refresh"></span></button>
  </div>
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
----------------------| ---`
    }
  },
  methods: {
    handleClick (tab, event) {
      this[tab.name]()
    },
    getinfo () {
      this.$store.dispatch('auth/info').then(() => {
        this.me()
      })
    },
    me () {
      this.md = this.head
      this.md += '\n reader | ' + this.$store.state.auth.reader
      this.md += '\n operator | ' + this.$store.state.auth.operator
      this.md += '\n admin | ' + this.$store.state.auth.admin
      for (let k in this.$store.state.auth.user) {
        this.md += '\n' + k + ' | ' + this.$store.state.auth.user[k]
      }
      this.activeName = 'me'
      return
    },
    falcon () {
      this.md = this.head
    }
  },
  computed: {
    compiledMarkdown: function () {
      return marked(this.md, { sanitize: true })
    },
    auth () {
      return this.$store.state.auth
    }
  },
  created () {
    this[this.activeName]()
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
