export const isBrowser = typeof window !== 'undefined' && window.document && document.createElement
export const isNode = !isBrowser && typeof global !== 'undefined'

const { _, axios, Cookies } = window
const token = {
  name: Cookies.get('name'),
  sig: Cookies.get('sig')
}

export function fetch (opts = {}) {
  const hasSpin = isBrowser && opts.spin
  if (hasSpin) {
    if (typeof opts.spin === 'string') {
      opts.spin = Array.from(document.querySelectorAll(opts.spin))
    } else if (!Array.isArray(opts.spin)) {
      opts.spin = [opts.spin]
    }

    opts.spin.forEach((el) => {
      el.style.visibility = 'visible'
    })
  }

  return new Promise((resolve, reject) => {
    axios(_.merge({
      headers: {
        ApiToken: JSON.stringify(token)
      },
      baseURL: 'http://localhost:8001/v1.0'
    }, opts))
    .then((res) => {
      if (hasSpin) {
        opts.spin.forEach((el) => {
          el.style.visibility = 'hidden'
        })
      }

      resolve(res)
    })
    .catch((err) => {
      if (hasSpin) {
        opts.spin.forEach((el) => {
          el.style.visibility = 'hidden'
        })
      }
      reject(err)
    })
  })
}
