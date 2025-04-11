
import './style.css'
import AiChatWidget from './ai-chat'

export function init() {
  let config = {
    iframeSrc: import.meta.env.VITE_AI_CHAT_BASE_URL + '/#/chat',
    remote: '',
    params: {}
  };

  const sdkEl = document.getElementById("ai_chat_js")

  if(sdkEl){
    let params = sdkEl.getAttribute("data-json")
    let origin = new URL(sdkEl.src).origin

    if(import.meta.env.DEV){
      // 开发者模式iframe地址使用本地地址
      config.iframeSrc = import.meta.env.VITE_AI_CHAT_BASE_URL + '/#/chat'
    }else{
      config.iframeSrc = origin + '/web/#/chat'
    }
    
    try{
      config.params = JSON.parse(params)
    } catch (error) {
      console.error('Failed to stringify data:', error);
      return;
    }
  }

  return AiChatWidget.init(config)
}

