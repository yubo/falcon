<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 class="page-header">About Falcon</h1>
  <div v-html="compiledMarkdown"></div>
</div>
</template>

<script>
import marked from 'marked'
export default {
  data () {
    return {
      md: `# about
- 1
- 2
- 3
  * 1
  * 2
  * 3
`
    }
  },
  computed: {
    compiledMarkdown: function () {
      return marked(this.md, { sanitize: true })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
