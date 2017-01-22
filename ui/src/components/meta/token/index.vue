<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <h1 class="page-header">token</h1>
  <p>tag配置</p>
  <div class="container">
      <div class="form-group">
        <label for="name">string(配置string)</label>
        <input type="text" name="string" id="string" class="form-control" value="cop=xiaomi" />
      </div>
    
      <div class="form-group">
        <label for="name">slices(配置slices)</label>
        <input type="text" name="slices" id="slices" class="form-control" value="s1,s2,s3" />
      </div>
    
    <button type="button" class="btn btn-default" id="edit-config"> <span class="glyphicon glyphicon-floppy-disk"></span> 更新 </button>
  </div>
</div>
</template>

<script>
export default {
  data () {
    return {}
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
