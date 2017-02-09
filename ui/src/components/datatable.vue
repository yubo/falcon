<template>
<div id="content" class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
  <div class="form-inline" role="form">
    <div class="form-group">
      <input type="text" v-model="searchBy" class="form-control" placeholder="name or cnname">
    </div>
    <button type="button" class="btn btn-primary">search</button>
    <input type="checkbox" id="iamcreator">
    <span>只显示我创建的</span>

    <div class="pull-right">
      <a href="#" class="btn btn-default"><span class="glyphicon glyphicon-plus"></span>Add</a>
    </div>
  </div>

  <div class="mt20">
    <table class="table table-striped table-bordered">
      <thead> <tr> <th v-for="(column, idx) in dataTable.columns">{{column.text}}</th> </tr> </thead>
      <tbody>
        <tr v-for="(row, row_idx) in filteredRows">
          <td v-for="(item, item_idx) in row"> {{item.value}} </td>
        </tr>
      </tbody>
    </table>
      <!--
      <div class="v-table-footer-info">
        Showing {{firstRow + 1}} to {{lastRow}} of {{filteredRows.length}} items
      </div>
      -->

    <ul class="pagination pagination-sm mt0" v-if="lastPage !== 1">
      <li :class="{disabled: currentPage == 1}">
        <a @click="togglePage('prev')" >Prev</a>
      </li>
      <li :class="{current: currentPage == 1}">
        <a class="v-table-footer-page-btn" href="#" @click="togglePage(1)">1</a>
      </li>
      <li class="disabled" v-if="currentPage >= 5 && lastPage > 10"><a>...</a></li>
      <li>
        <a href="#" :class="{current: currentPage == page + 1}"
        @click="togglePage(page + 1)" v-for="(page, idx) in centerPartPage">{{page + 1}}</a>
      </li>
      <li class="disabled" v-if="lastPage > 10 && lastPage - currentPage > 5"><a>...</a></li>
      <li>
        <a href="#" :class="{current: currentPage == page + 1}"
          @click="togglePage(page + 1)" v-for="(page, idx) in lastPartPage">{{page + 1}}</a>
      </li>
      <li :class="{disabled: currentPage == lastPage}">
        <a href="#" @click="togglePage('next')" >Next</a>
      </li>
    </ul>

    <div class="pull-right">
      <span>Show
      <select v-model="dataTable.options.pageCount" @change="onChangePageCount()">
        <option>5</option>
        <option>10</option>
        <option>20</option>
        <option>50</option>
      </select>
      items each page
      </span>
    </div>

  </div>
</div>
</template>

<script>
export default {
  props: ['dataTable'],
  data () {
    return {
      currentPage: 1,
      searchBy: '',
      rows: [],
      sort: {
        sortBy: '',
        desc: true
      }
    }
  },
  computed: {
    filteredRows () {
      return this.dataTable.rows
    },

    lastPage () {
      return Math.ceil(this.filteredRows.length / this.dataTable.options.pageCount)
    },

    centerPartPage () {
      if (this.lastPage > 10 && this.currentPage >= 5) {
        if (this.lastPage - this.currentPage > 5) {
          return this.currentPage === this.lastPage ? [this.currentPage - 3, this.currentPage - 2, this.currentPage - 1] : [this.currentPage - 2, this.currentPage - 1, this.currentPage]
        } else {
          const r = []

          for (let i = this.lastPage - 6; i < this.lastPage; i++) {
            r.push(i)
          }
          return r
        }
      } else if (this.lastPage > 10) {
        const r = []

        for (let i = 1; i < 5; i++) {
          r.push(i)
        }
        return r
      } else {
        const r = []

        for (let i = 1; i < this.lastPage; i++) {
          r.push(i)
        }
        return r
      }
    },

    lastPartPage () {
      if (this.lastPage > 10 && this.lastPage - this.currentPage > 5) {
        return [this.lastPage - 1]
      } else {
        return []
      }
    },
    firstRow () {
      return this.currentPage === 1 ? 0 : this.dataTable.options.pageCount * (this.currentPage - 1)
    },
    lastRow () {
      return this.dataTable.options.pageCount * this.currentPage > this.filteredRows.length ? this.filteredRows.length : this.dataTable.options.pageCount * this.currentPage
    }
  },
  watch: {
    'dataTable.rows' (rows) {
      rows.forEach((row, index) => {
        for (let key in row) {
          const column = this.dataTable.columns.filter((column) => {
            return column.value === key
          })[0]

          row[key] = Object.assign({
            editable: column.editable,
            editing: false,
            tmpValue: ''
          }, row[key])
        }

        this.dataTable.rows[index] = row
      })
    },

    'dataTable.columns' (columns) {
      columns.forEach((column, index) => {
        column = Object.assign({
          editable: false,
          sortable: false
        }, column)

        this.dataTable.columns[index] = column
      })
    },

    'searchBy' (val) {
      if (val) {
        this.currentPage = 1
      }
    }
  },

  filters: {
    pagination (rows, currentPage, pageCount) {
      return this.getPageRows(rows, currentPage, pageCount)
    }
  },

  methods: {
    onChangePageCount () {
      this.currentPage = 1
    },

    getPageRows (rows) {
      return rows.slice(this.firstRow, this.lastRow)
    },

    togglePage (page) {
      switch (page) {
        case 'prev':
          if (this.currentPage <= 1) return
          this.currentPage--
          break
        case 'next':
          if (this.currentPage >= this.lastPage) return
          this.currentPage++
          break
        default:
          if (this.currentPage === page) return
          this.currentPage = page
      }
      if (this.dataTable.onPageChanged) {
        this.dataTable.onPageChanged(this.currentPage)
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
