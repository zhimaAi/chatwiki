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
      headers: {
        'App-Type': 'yun_h5',
      },
      signal: this.controller.signal,
      // 允许在页面隐藏时继续接收消息(开启后不再触发自动重连的问题)
      openWhenHidden: true,
      body: formdata,
      onopen: async (response) => {
        if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
          if (typeof this.onOpen === 'function') {
            this.onOpen()
          }

          return
        }

        console.error(`连接失败: ${response.status} ${response.statusText}`)
        throw new Error(`连接失败: ${response.status} ${response.statusText}`);
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
        
        this.abort()
      },
      onerror: (error: any) => {
        console.log(error)
        // 明确指定错误类型为any
        if (typeof this.onError === 'function') {
          this.onError(error)
        }

        this.abort();
        throw error;
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
