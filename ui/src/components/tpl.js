const { _ } = window

var liTpl = {
  template: `<li :class="classObject"> <router-link :to="obj.url">{{ obj.text }}</router-link></li>`,
  props: ['obj'],
  computed: {
    classObject () {
      return {
        active: _.startsWith(this.$route.path, this.obj.url)
      }
    }
  }
}

export {
  liTpl
}
