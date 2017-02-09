<template>
<div class="container">
    <label for="inputUsername" class="sr-only">Username</label>
    <input name="username" type="text" class="form-control" placeholder="Username" ref="username"
      required autofocus>
    <label for="inputPassword" class="sr-only">Password</label>
    <input name="password" type="password" class="form-control" placeholder="Password" ref="password"
      required>
    <button class="btn btn-lg btn-primary btn-block" @click="login">Sign in</button>
</div>
</template>

<script>
export default {
  data () {
    return { }
  },
  methods: {
    login () {
      this.$store.dispatch('login/login', {
        username: this.$refs.username.value,
        password: this.$refs.password.value,
        method: 'ldap',
        router: this.$router
      })
    }
  }
}
</script>

<style>
</style>
