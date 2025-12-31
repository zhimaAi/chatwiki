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
    this.opt = { ...this.opt, ...opt }

    this.open()
  }

  open() {
    let formdata = new FormData()

    for (let key in this.opt.data) {
      formdata.append(key, this.opt.data[key])
    }

    let headersData = {
      'App-Type': ''
    }

    if (this.opt.token) {
      headersData = {
        'App-Type': '',
        'token': this.opt.token
      }
    }

    fetchEventSource(this.opt.url, {
      method: 'POST',
      headers: headersData,
      signal: this.controller.signal,
      // 允许在页面隐藏时继续接收消息(开启后不再触发自动重连的问题)
      openWhenHidden: true,
      body: formdata,
      onopen: async(response) => {
        if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
          if (typeof this.onOpen === 'function') {
            this.onOpen()
          }
          return // everything's good
        }
        
        throw new Error(`连接失败: ${response.status} ${response.statusText}`);
      },
      onmessage: (res) => {
        if (typeof this.onMessage === 'function') {
          this.onMessage(res)
        }
      },
      onclose: () => {
        if (typeof this.onClose === 'function') {
          this.onClose()
        }

        this.abort();
      },
      onerror: (err) => {
        console.log(err)
        // 可以在这里添加错误处理逻辑
        if (typeof this.onError === 'function') {
          this.onError(err)
        }

        this.abort();
        throw err
      }
    })
  }

  abort = () => {
    // 增加判断，避免重复 abort 报错
    if (!this.controller.signal.aborted) {
      this.controller.abort()
    }
  }
}
