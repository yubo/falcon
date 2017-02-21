export const isBrowser = typeof window !== 'undefined' && window.document && document.createElement
export const isNode = !isBrowser && typeof global !== 'undefined'

const { _, axios } = window
import store from 'src/store'
import { Message, Notification, MessageBox } from 'element-ui'

function fetch (opts = {}) {
  return new Promise((resolve, reject) => {
    axios(_.merge({ baseURL: '/v1.0' }, opts)).then((res) => {
      resolve(res)
    }).catch((err) => {
      if (err.response && err.response.status === 401) {
        store.commit('auth/m_logout')
      }
      reject(err)
    })
  })
}

var Msg = {
  error: (msg, err = {}) => {
    let m = err && err.response && err.response.data ? err.response.data : msg
    if (msg.length < 64) {
      Message.error(m)
    } else {
      Notification.error({title: 'Error', message: m})
    }
  },
  info: Message.info,
  success: (msg, res = {}) => {
    let m = res && res.data ? res.data : msg
    if (msg.length < 64) {
      Message.success(m)
    } else {
      Notification.success({title: 'Success', message: m})
    }
  },
  confirm: MessageBox.confirm
}

export {
  fetch,
  Msg
}
