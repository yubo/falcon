export const isBrowser = typeof window !== 'undefined' && window.document && document.createElement
export const isNode = !isBrowser && typeof global !== 'undefined'

const { _, axios } = window

export function fetch (opts = {}) {
  return new Promise((resolve, reject) => {
    axios(_.merge({
      baseURL: '/v1.0'
    }, opts))
    .then((res) => {
      resolve(res)
    })
    .catch((err) => {
      reject(err)
    })
  })
}
