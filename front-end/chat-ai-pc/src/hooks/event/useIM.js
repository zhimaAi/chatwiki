import { getWsUrl } from '@/api/chat/index'
import { reactive, ref } from 'vue'
import reconnectableSocket from '../../utils/http/websocket/reconnectableSocket'

const wsUrl = ref(null)
const ws = ref(null)
const events = reactive({})

export const useIM = () => {
  async function connect(openid, cb) {
    let res = await getWsUrl({ openid: openid })

    if (res) {
      wsUrl.value = res.data.ws_url
    }

    if (ws.value) {
      ws.value.close()
    }

    ws.value = new reconnectableSocket(wsUrl.value, {
      disableHeartbeat: true
    })

    cb && cb()

    ws.value.on('message', (data) => {
      emit('message', data)
    })

    ws.value.on('error', onError)
  }

  function on(event, listener) {
    if (!events[event]) {
      events[event] = []
    }

    events[event].push(listener)
  }

  function emit(event, ...args) {
    if (events && Array.isArray(events[event])) {
      events[event].forEach((listener) => listener(...args))
    }
  }

  function onError(err) {
    console.log(err)
  }

  function close() {
    Object.keys(events).forEach((event) => {
      events[event] = []

      delete events[event]
    })

    ws.value && ws.value.close()
  }

  return {
    connect,
    on: on,
    close
  }
}
