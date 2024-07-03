import { useEventBus } from '@/hooks/event/useEventBus'

const emitter = useEventBus()

window.addEventListener('message', function(event) {
    const res = event.data
    
    if(res.action === 'openWindow'){
        emitter.emit('openWindow', event.data)
    }
    
});