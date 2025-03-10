export default class Socket {
  constructor(url, options) {
    this.debug = options.debug || false
    this.socketOpen = false
    this.socketMsgQueue = []
    this.heartbeatInterval = 5 * 1000
    this.heartbeatTimer = null
    this.disableHeartbeat = options.disableHeartbeat
    // 禁止重连，手动关闭socket（退出登录）等场景不需要在重连了
    this.disableReconnect = false
    this.reconnectAttempts = 0
    this.reconnectInterval = 5 * 1000
    this.reconnectTimer = null
    this.events = {}

    this.url = url
    this.ws = null

    this.options = options

    this.connect()

    this.log = function () {}

    if (this.debug) {
      this.log = console.log
    }
  }

  connect() {
    if (this.ws && this.socketOpen) {
      // console.log('Socket is already connected');
      return
    }

    this.ws = new WebSocket(this.url)

    this.ws.onopen = () => {
      this.log('WebSocket open')
      this.onOpen()
    }

    this.ws.onmessage = (event) => {
      this.log('WebSocket message')
      this.log(event)
      this.onMessage(event)
    }

    this.ws.onerror = (error) => {
      this.log('WebSocket error')
      this.log(error)
      this.onError(error)
    }

    this.ws.onclose = (event) => {
      this.log('WebSocket closed')
      this.log(event)
      this.onClose(event)
    }
  }

  onOpen() {
    this.socketOpen = true
    this.reconnectAttempts = 0 // 重连成功重置重连次数
    this.disableReconnect = false // 恢复断线重连功能

    // 清除重连定时器
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    this.sendHeartBeat() // 发送心跳包

    // 离线消息重发
    for (var i = 0; i < this.socketMsgQueue.length; i++) {
      this.send(this.socketMsgQueue[i])
    }

    this.socketMsgQueue = []

    this.emit('open')
  }

  onMessage(event) {
    // 可以在这里添加接收到消息后的逻辑
    if (event && event.data === 'ping') {
      this.pong(event)
    } else {
      this.emit('message', JSON.parse(event.data))
    }
  }

  onError(error) {
    // 可以在这里添加错误处理逻辑
    this.emit('error', error)
  }

  onClose(event) {
    // 可以在这里添加断开连接后的逻辑，例如重新连接等
    this.emit('close', event)

    this.ws = null
    this.socketOpen = false

    // 清除重连定时器
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    // 停止心跳包
    this.stopHeartBeat()

    // 判断是否需要重连
    if (this.shouldReconnect() && !this.disableReconnect) {
      this.reconnect()
    }
  }

  send(message) {
    if (this.ws && this.socketOpen) {
      this.ws.send(JSON.stringify(message))
    } else {
      this.socketMsgQueue.push(message)
    }
  }
  pong() {
    this.ws.send('pong')
  }
  sendHeartBeat() {
    if (this.disableHeartbeat) {
      return
    }
    clearTimeout(this.heartbeatTimer) // 清除之前的心跳超时定时器
    // 每5秒发送一次心跳包
    this.heartbeatTimer = setTimeout(() => {
      this.send({
        type: 'ping'
      })

      // 继续发送心跳包
      this.sendHeartBeat()
      this.emit('ping')
    }, this.heartbeatInterval)
  }

  stopHeartBeat() {
    if (this.disableHeartbeat) {
      return
    }
    clearTimeout(this.heartbeatTimer) // 清除心跳超时定时器
  }

  // 是否需要重连
  shouldReconnect() {
    const maxReconnectAttempts = this.options.maxReconnectAttempts || Infinity
    return this.reconnectAttempts < maxReconnectAttempts
  }

  reconnect() {
    if (this.reconnectTimer) {
      return
    }

    this.disableReconnect = false
    this.reconnectAttempts++

    // 每5秒尝试连接一次
    this.reconnectTimer = setTimeout(() => {
      if (this.ws) {
        this.ws.close()
        this.socketOpen = false
        this.ws = null
      } else {
        this.connect()
      }
    }, this.reconnectInterval)
  }

  close(disableReconnect = true) {
    this.disableReconnect = disableReconnect

    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
    }

    if (this.ws) {
      this.ws.close()
    }
  }

  // 事件发射器
  emit(event, ...args) {
    if (this.events && Array.isArray(this.events[event])) {
      this.events[event].forEach((listener) => listener(...args))
    }
  }

  // 添加事件监听器
  on(event, listener) {
    if (!this.events) {
      this.events = {}
    }
    // 如果 events[event] 不存在或者不是数组，则初始化为数组
    if (!this.events[event]) {
      this.events[event] = []
    }
    // 添加监听器到数组
    this.events[event].push(listener)
  }

  // 移除事件监听器
  off(event, listener) {
    if (this.events && Array.isArray(this.events[event])) {
      // 如果没有指定 listener，则移除所有监听器
      if (!listener) {
        this.events[event] = []
      } else {
        // 否则，移除指定的监听器
        const index = this.events[event].indexOf(listener)
        if (index !== -1) {
          this.events[event].splice(index, 1)
        }
      }
    }
  }
}
