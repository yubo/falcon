<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <subbar></subbar>
  <div class="container">
    <table class="table table-striped">
      <tbody>
        <tr> <td>id </td> <td>2 </td> </tr>
        <tr> <td>uuid </td> <td>cn=yubo,ou=users,dc=yubo,dc=org@ldap </td> </tr>
        <tr> <td>用户名 </td> <td>cn=yubo,ou=users,dc=yubo,dc=org@ldap </td> </tr>
        <tr> <td>中文名 </td> <td> </td> </tr>
        <tr> <td>email </td> <td> </td> </tr>
        <tr> <td>phone </td> <td> </td> </tr>
        <tr> <td>im </td> <td> </td> </tr>
        <tr> <td>qq </td> <td> </td> </tr>
      </tbody>
    </table>
  </div>
</div>
</template>

<script>
import subbar from './subbar'
export default {
  components: {
    subbar
  },
  data () {
    return {
      msg: 'Welcome to falcon@xiaomi.com App 5'
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
