import { useEventBus } from '@/hooks/event/useEventBus'
import { useChatStore } from '@/stores/modules/chat'
const emitter = useEventBus()

window.addEventListener('message', function(event) {
    const res = event.data
    // console.log(res,'data===')
    if(res.action === 'openWindow'){
        emitter.emit('openWindow', event.data)
    }
    const chatStore = useChatStore()
    const { upDataUiStyle, updataQuickComand } = chatStore
    const type = res.type;
    switch (type) {
      case "onPreview": {
        upDataUiStyle(res.data)
        break;
      }
      case "updataQuickComand": {
        updataQuickComand(res.data)
        break;
      }
    }
});