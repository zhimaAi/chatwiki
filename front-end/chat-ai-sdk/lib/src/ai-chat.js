import { objectToQueryString } from './util'
import AiAvatar from './ai-avatar'
import AiDot from './ai-dot'
import NewMessage from './new-message'

class AiChatWidget {
  iframe = null;
  iframeSrc = null
  config = {}
  
  constructor() {
    
  }

  init(config) {
    this.config = config;
    this.iframeSrc = this.config.iframeSrc;

    this.insertAiChat();
    
    this.getExpectedOrigin()
    window.addEventListener("message", this.handleMessage.bind(this), false);

    return {
      open: this.open.bind(this),
      close: this.close.bind(this),
    };
  }

  removeAiChat() {
    const iframe = document.getElementById("zm_chat-wiki-iframe");
    if (iframe) {
      document.body.removeChild(iframe);
      this.iframe = null;
    }
  }

  onInit(data) {
    AiAvatar.init(data)
  }

  open() {
    this.iframe.style.display = "block";
    this.postMessage('openWindow', {});
  }

  close() {
    AiAvatar.show();
    this.iframe.style.display = "none";
    this.postMessage('closeWindow', {});
  }

  onClose() {
    AiAvatar.show()
    this.iframe.style.display = "none";
  }

  createDot(data) {
    AiDot.create({value: data})
  }

  createNewMessage(data) {
    let list = data || [];
    
    if(list.length === 0){
      NewMessage.remove()
      return
    }

    // 截取数组的最后一个元素
    list = list.slice(-1);

    NewMessage.create(list)
  }

  getExpectedOrigin(){
    let origin = new URL(this.iframeSrc).origin

    return origin
  }

  handleMessage(event) {
    if (!event.origin) {
      return;
    }

    const url = new URL(event.origin);
    const isExpectedOrigin = url.origin === this.getExpectedOrigin();

    if (!isExpectedOrigin) {
      return;
    }

    let res = event.data;

    if(import.meta.env.DEV){
      console.log('Received message from parent:', res);
    }
    
    if(!res){
      return
    }

    if (res.action === "closeChat") {
      this.onClose();
    }

    if (res.action === "init") {
      this.onInit(res.data)
    }
    // 更新未读消息数
    if (res.action === "dot") {
      this.createDot(res.data);
    }
    // 更新未读消息
    if (res.action === "newMessage") {
      this.createNewMessage(res.data);
    }
  }

  postMessage(action, data) {
    if (data) {
      try {
        data = JSON.parse(JSON.stringify(data));
      } catch (error) {
        console.error("Failed to stringify data:", error);
        return;
      }
    }

    if (this.iframe.contentWindow && typeof this.iframe.contentWindow.postMessage === "function") {
      try {
        this.iframe.contentWindow.postMessage({ action: action, data }, "*");
      } catch (error) {
        console.error("Failed to post message:", error);
      }
    } else {
      console.warn(
        "frame.contentWindow is not available or postMessage is not supported."
      );
    }
  }

  insertAiChat() {
    if (document.getElementById("zm_chat-wiki-iframe")) {
      return;
    }
    const queryStr = objectToQueryString(this.config.params)
 
    this.iframe = document.createElement("iframe");
    this.iframe.id = "zm_chat-wiki-iframe";
    this.iframe.src = this.iframeSrc + '?' + queryStr;
    this.iframe.style.display = "none";

    document.body.appendChild(this.iframe);
  }

  
}

export default new AiChatWidget();
