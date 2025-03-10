import { objectToQueryString } from './util'

class AiChatWidget {
  iframe = null;
  iframeSrc = null
  onClose = null;
  onInit = null;
  
  config = {}
  
  constructor(config) {
    if (config) {
      this.config = config;
      this.iframeSrc = this.config.iframeSrc

      this.init();
    }
  }

  init() {
    if (this.iframeSrc) {
      this.insertAiChat();
    }
    
    this.getExpectedOrigin()
    window.addEventListener("message", this.handleMessage.bind(this), false);
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
      this.close();
    }

    if (res.action === "init") {
      this.onInit(res.data);
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

  removeAiChat() {
    const iframe = document.getElementById("zm_chat-wiki-iframe");
    if (iframe) {
      document.body.removeChild(iframe);
      this.iframe = null;
    }
  }

  open() {
    this.iframe.style.display = "block";
    this.postMessage('openWindow', {});
  }

  close() {
    if (this.onClose) {
      this.iframe.style.display = "none";
      this.onClose();
    }
  }
}

export default AiChatWidget;
