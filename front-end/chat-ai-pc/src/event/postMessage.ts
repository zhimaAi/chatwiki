function postMessage(action: string, data?: any){
    if(data){
        try {
            data = JSON.parse(JSON.stringify(data));
        } catch (error) {
            console.error('Failed to stringify data:', error);
            return;
        }
    }

    if (window.parent && typeof window.parent.postMessage === 'function') {
        try {
            window.parent.postMessage({ action: action, data }, '*');
        } catch (error) {
            console.error('Failed to post message:', error);
        }
    } else {
        console.warn('window.parent is not available or postMessage is not supported.');
    }
}


export function postInit(data: any){
    postMessage('init', data)
}

export function postCloseChat(){
    postMessage('closeChat')
} 