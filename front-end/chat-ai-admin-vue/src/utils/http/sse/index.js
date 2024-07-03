import { fetchEventSource, EventStreamContentType } from '@/libs/fetch-event-source'

export default class SSE {
  onOpen = undefined
  onClose = undefined
  onError = undefined
  onMessage = undefined

  opt = {
    url: '',
    data: {}
  }

  controller = new AbortController()

  constructor(opt = { url: '', data: {} }) {
    this.socketMsgQueue = []

    this.opt = { ...this.opt, ...opt }

    this.open()
  }

  open() {
    const that = this

    let formdata = new FormData()

    for (let key in that.opt.data) {
      formdata.append(key, that.opt.data[key])
    }

    fetchEventSource(this.opt.url, {
      method: 'POST',
      headers: {
        // 'Content-Type': 'application/x-www-form-urlencoded;charset=UTF-8'
      },
      signal: that.controller.signal,
      // 允许在页面隐藏时继续接收消息(开启后不再触发自动重连的问题)
      openWhenHidden: true,
      body: formdata,
      async onopen(response) {
        if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
          if (typeof that.onOpen === 'function') {
            that.onOpen()
          }
          return // everything's good
        } else if (response.status >= 400 && response.status < 500 && response.status !== 429) {
          // client-side errors are usually non-retriable:
          throw new Error('连接出错')
        } else {
          throw new Error('连接出错')
        }
      },
      onmessage(res) {
        if (typeof that.onMessage === 'function') {
          that.onMessage(res)
        }
      },
      onclose() {
        if (typeof that.onClose === 'function') {
          that.onClose()
        }
        that.controller.abort()
      },
      onerror(err) {
        console.log(err)
        // 可以在这里添加错误处理逻辑
        if (typeof that.onError === 'function') {
          that.onError(err)
        }
        throw err
      }
    })
  }

  abort = () => {
    this.controller.abort()
  }
}
