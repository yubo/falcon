export const isBrowser = typeof window !== 'undefined' && window.document && document.createElement
export const isNode = !isBrowser && typeof global !== 'undefined'

const { _, axios } = window

function fetch (opts = {}) {
  return new Promise((resolve, reject) => {
    axios(_.merge({
      baseURL: '/v1.0'
    }, opts))
    .then((res) => {
      resolve(res)
    })
    .catch((err) => {
      if (err.response && err.response.status === 401) {
        if (opts.router) {
          opts.router.push('/login')
        }
      }
      console.log(err.response)
      reject(err)
    })
  })
}

function vfetch (opts = {}) {
  const commit = (status, arg) => {
    const hasMutation = opts.commit && opts.mutation
    if (hasMutation) {
      return opts.commit(`${opts.mutation}.${status}`, arg)
    }
  }

  return new Promise((resolve, reject) => {
    commit('start')
    fetch(opts).then((res) => {
      commit('success', { res, ...opts.args })
      commit('end')
      resolve(res)
    }).catch((err) => {
      commit('fail', { err, ...opts.args })
      commit('end')
      reject(err)
    })
  })
}

export {
  fetch,
  vfetch
}
