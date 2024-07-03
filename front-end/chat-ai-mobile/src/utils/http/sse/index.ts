import { fetchEventSource, EventStreamContentType } from '@/libs/fetch-event-source'
import type { EventSourceMessage } from '@/libs/fetch-event-source'

// 定义SSEOptions类型
type SSEOptions = {
  url: string
  data?: Record<string, string | Blob | File>
}

// 定义事件处理函数的类型
type EventHandler = () => void
type MessageHandler = (event: EventSourceMessage) => void
type ErrorHandler = (error: any) => void

export default class SSE {
  onOpen?: EventHandler
  onClose?: EventHandler
  onError?: ErrorHandler
  onMessage?: MessageHandler

  opt: SSEOptions = {
    url: '',
    data: {}
  }

  controller = new AbortController()

  constructor(opt: Partial<SSEOptions> = {}) {
    this.opt = { ...this.opt, ...opt }
    this.open()
  }

  open() {
    const formdata = new FormData()

    for (const key in this.opt.data) {
      if (Object.prototype.hasOwnProperty.call(this.opt.data, key)) {
        formdata.append(key, this.opt.data[key] as string | Blob | File)
      }
    }

    fetchEventSource(this.opt.url, {
      method: 'POST',
      signal: this.controller.signal,
      // 允许在页面隐藏时继续接收消息(开启后不再触发自动重连的问题)
      openWhenHidden: true,
      body: formdata,
      onopen: async (response) => {
        if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
          if (typeof this.onOpen === 'function') {
            this.onOpen()
          }
        } else if (response.status >= 400 && response.status < 500 && response.status !== 429) {
          throw new Error('连接出错')
        } else {
          throw new Error('连接出错')
        }
      },
      onmessage: (event: EventSourceMessage) => {
        // 明确指定事件类型为MessageEvent
        if (typeof this.onMessage === 'function') {
          this.onMessage(event)
        }
      },
      onclose: () => {
        if (typeof this.onClose === 'function') {
          this.onClose()
        }
        this.controller.abort()
      },
      onerror: (error: any) => {
        // 明确指定错误类型为any
        if (typeof this.onError === 'function') {
          this.onError(error)
        }
        // throw error;
      }
    })
  }

  abort() {
    this.controller.abort()
  }
}
