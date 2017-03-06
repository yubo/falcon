<template>
  <div class="container-fluid">
    <div class="row">
      <div class="col-sm-3">
        <div class="pull-right">
          <button type="button" @click="handleAdd" class="btn btn-default"><span class="glyphicon glyphicon-plus"></span></button>
          <button type="button" @click="handleDel" class="btn btn-default"><span class="glyphicon glyphicon-remove"></span></button>
          <button type="button" @click="handleReload" class="btn btn-default"><span class="glyphicon glyphicon-refresh"></span></button>
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

      <el-dialog :title="'add node('+tag+')'" v-model="editVisible" :close-on-click-modal="false">
        <el-form label-position="right" label-width="80px" :model="objForm">
          <el-form-item label="node">
            <el-select
              v-if="optionPrefixs.length > 0" 
              style="width:200px;"
              v-model="tagKey"
              filterable>
              <el-option
                v-for="k in optionPrefixs"
                :key="k"
                :label="k+'='"
                :value="k+'='">
              </el-option>
            </el-select>

            <el-input v-model="objForm.tagValue" style="width:200px"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="submitObj">submit</el-button>
            <el-button @click="editVisible = false">Cancel</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>

    </div>
  </div>
</template>

<script>
import { fetch, Msg } from 'src/utils'
const { _ } = window
export default {
  data () {
    return {
      optionPrefixs: [],
      editVisible: false,
      tagKey: '',
      curTag: {},
      objForm: {
        tagValue: ''
      },
      props: {
        label: 'label',
        children: 'children'
      }
    }
  },
  methods: {
    handleCurrentChange (val) {
      if (!val.ro) {
        this.$store.commit('rel/m_cur_tag', val)
      }
      this.curTag = val
    },
    handleReload () {
      this.$store.dispatch('rel/load_tree')
    },
    submitObj () {
      fetch({
        method: 'post',
        url: 'tag',
        data: {name: this.tag}
      }).then((res) => {
        Msg.success('update success')
        this.$store.commit('rel/m_add_node', {
          id: this.curTag.id,
          cid: res.data.id,
          label: this.tagKey + this.objForm.tagValue,
          name: this.tag
        })
        this.editVisible = false
      }).catch((err) => {
        Msg.error('update failed', err)
      })
    },
    handleAdd () {
      if (!this.curTag.name) {
        Msg.error('please select node')
        return
      }
      if (this.curTag.ro) {
        Msg.error('Permission denied')
        return
      }

      this.tagKey = ''
      const s = this.$store.getters.schema
      if (s === '') {
        this.editVisible = true
        return
      }

      const k1 = _.split(this.curTag.name, ',').map((v) => {
        return _.split(v, '=')[0]
      })

      let k2 = []
      for (let i = 0, key = ''; i < s.length; i++) {
        if (s[i] === ',') {
          k2.push({key: key, must: true})
          key = ''
        } else if (s[i] === ';') {
          k2.push({key: key, must: false})
          key = ''
        } else {
          key += s[i]
        }
      }

      let j = 0
      for (let i = 0; i < k1.length && j < k2.length; i++, j++) {
        for (; j < k2.length; j++) {
          if (k1[i] === k2[j].key) {
            break
          }
        }
      }

      this.optionPrefixs = []
      for (; j < k2.length; j++) {
        this.optionPrefixs.push(k2[j].key)
        if (k2[j].must) {
          break
        }
      }
      if (this.optionPrefixs.length === 0) {
        Msg.error('Permission denied, can not add node at leaf')
        return
      }
      this.editVisible = true
    },
    handleDel (obj) {
      if (!this.curTag.name) {
        Msg.error('please select node')
        return
      }
      if (this.curTag.ro) {
        Msg.error('Permission denied')
        return
      }
      Msg.confirm('此操作将永久删除该记录,以及相关的绑定关系，包括用户，机器，策略模板等 是否继续?', '提示', {
        confirmButtonText: 'Confirm',
        cancelButtonText: 'Cancel',
        type: 'warning'
      }).then(() => {
        fetch({
          method: 'delete',
          url: 'tag/' + this.curTag.id
        }).then((res) => {
          this.$store.commit('rel/m_del_node', this.curTag.id)
          this.$store.dispatch('rel/load_tree')
          Msg.success('success!')
        }).catch((err) => {
          Msg.error('delete failed', err)
        })
      }).catch(() => {
        Msg.info('cancel')
      })
    }
  },
  computed: {
    loading () {
      return this.$store.state.rel.loading
    },
    tagTree () {
      return this.$store.state.rel.tree
    },
    tag () {
      return this.curTag.name + ',' + this.tagKey + this.objForm.tagValue
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
.pull-right .btn {
  margin-top: 2px;
}

</style>
