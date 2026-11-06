import { getWsUrl } from '@/api/app/index'
import { ref, reactive } from 'vue'
import reconnectableSocket from '../../utils/http/websocket/reconnectableSocket'

const wsUrl = ref(null)
const ws = ref(null)
const events = reactive({})

const defaultState = { wsUrl, ws, events, version: 0 }
export const useIM = (isolated = false) => {
  const state = isolated ? { wsUrl: ref(null), ws: ref(null), events: reactive({}), version: 0 } : defaultState
  const { wsUrl, ws, events } = state
  async function connect(openid, cb) {
    const version = ++state.version
    let res = await getWsUrl({ openid: openid })

    if (version !== state.version) return
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
      if (version === state.version) emit('message', data)
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
    state.version++
    Object.keys(events).forEach((event) => {
      events[event] = []

      delete events[event]
    })

    ws.value && ws.value.close()
    ws.value = null
  }

  return {
    connect,
    on: on,
    close
  }
}
