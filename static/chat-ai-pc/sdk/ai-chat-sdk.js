var l = Object.defineProperty;
var h = (a, t, i) => t in a ? l(a, t, { enumerable: !0, configurable: !0, writable: !0, value: i }) : a[t] = i;
var e = (a, t, i) => (h(a, typeof t != "symbol" ? t + "" : t, i), i);
function c(a) {
  if (!a)
    return "";
  const t = [];
  for (let i in a)
    if (a.hasOwnProperty(i)) {
      let n = a[i];
      (Array.isArray(n) || typeof n == "object" && n !== null) && (n = JSON.stringify(n)), n !== void 0 && t.push(encodeURIComponent(i) + "=" + encodeURIComponent(n));
    }
  return t.join("&");
}
function d() {
  const a = document.getElementById("ai_chat_js"), t = document.createElement("link"), i = new URL(a.src).origin;
  t.type = "text/css", t.rel = "stylesheet", t.href = i + "/sdk/style.css", document.getElementsByTagName("head")[0].appendChild(t);
}
class m {
  constructor(t) {
    e(this, "iframe", null);
    e(this, "iframeSrc", null);
    e(this, "onClose", null);
    e(this, "onInit", null);
    e(this, "config", {});
    t && (this.config = t, this.iframeSrc = this.config.iframeSrc, this.init());
  }
  init() {
    this.iframeSrc && this.insertAiChat(), this.getExpectedOrigin(), window.addEventListener("message", this.handleMessage.bind(this), !1);
  }
  getExpectedOrigin() {
    return new URL(this.iframeSrc).origin;
  }
  handleMessage(t) {
    if (!t.origin || !(new URL(t.origin).origin === this.getExpectedOrigin()))
      return;
    let s = t.data;
    s && (s.action === "closeChat" && this.close(), s.action === "init" && this.onInit(s.data));
  }
  postMessage(t, i) {
    if (i)
      try {
        i = JSON.parse(JSON.stringify(i));
      } catch (n) {
        console.error("Failed to stringify data:", n);
        return;
      }
    if (this.iframe.contentWindow && typeof this.iframe.contentWindow.postMessage == "function")
      try {
        this.iframe.contentWindow.postMessage({ action: t, data: i }, "*");
      } catch (n) {
        console.error("Failed to post message:", n);
      }
    else
      console.warn(
        "frame.contentWindow is not available or postMessage is not supported."
      );
  }
  insertAiChat() {
    if (document.getElementById("zm_chat-wiki-iframe"))
      return;
    const t = c(this.config.params);
    this.iframe = document.createElement("iframe"), this.iframe.id = "zm_chat-wiki-iframe", this.iframe.src = this.iframeSrc + "?" + t, this.iframe.style.display = "none", document.body.appendChild(this.iframe);
  }
  removeAiChat() {
    const t = document.getElementById("zm_chat-wiki-iframe");
    t && (document.body.removeChild(t), this.iframe = null);
  }
  open() {
    this.iframe.style.display = "block", this.postMessage("openWindow", {});
  }
  close() {
    this.onClose && (this.iframe.style.display = "none", this.onClose());
  }
}
class f {
  constructor(t) {
    e(this, "clickHandler", null);
    e(this, "avatarEl", null);
    e(this, "avatarSrc", "");
    e(this, "left", 0);
    e(this, "top", 0);
    e(this, "width", 50);
    e(this, "height", 50);
    e(this, "initialX", 0);
    e(this, "initialY", 0);
    e(this, "initialMouseX", 0);
    e(this, "initialMouseY", 0);
    // 拖拽状态标志位
    e(this, "dragging", !1);
    // 拖拽移动
    e(this, "handleDrag", (t) => {
      this.dragging = !0;
      const i = t.clientX - this.initialMouseX, n = t.clientY - this.initialMouseY, s = window.innerWidth, r = window.innerHeight;
      this.left = Math.max(0, Math.min(i, s - this.width)), this.top = Math.max(0, Math.min(n, r - this.height)), this.avatarEl.style.left = this.left + "px", this.avatarEl.style.top = this.top + "px";
    });
    // 拖拽结束
    e(this, "handleDragEnd", () => {
      const t = Math.abs(this.initialX - this.left), i = Math.abs(this.initialY - this.top);
      t <= 3 && i <= 3 && this.handleClick(), this.dragging = !1, document.removeEventListener("mousemove", this.handleDrag), document.removeEventListener("mouseup", this.handleDragEnd);
    });
    t && (this.avatarSrc = t.avatarSrc);
  }
  init(t) {
    this.avatarSrc = t.sdkFloatAvatar, this.getInitialPosition(), this.insertAvatar();
  }
  // 获取初始位置
  getInitialPosition() {
    const t = window.innerWidth, i = window.innerHeight;
    this.left = t - this.width - 50, this.top = i - this.height - 50;
  }
  // 启用拖拽
  enableDrag() {
    this.avatarEl.addEventListener("mousedown", (t) => {
      this.initialX = this.left, this.initialY = this.top, this.initialMouseX = t.clientX - this.avatarEl.getBoundingClientRect().left, this.initialMouseY = t.clientY - this.avatarEl.getBoundingClientRect().top, document.addEventListener("mousemove", this.handleDrag), document.addEventListener("mouseup", this.handleDragEnd), t.preventDefault();
    });
  }
  insertAvatar() {
    document.getElementById("zm_chat-wiki-avatar") || (this.avatarEl = document.createElement("img"), this.avatarEl.style.display = "block", this.avatarEl.style.top = this.top + "px", this.avatarEl.style.left = this.left + "px", this.avatarEl.style.width = this.width + "px", this.avatarEl.style.height = this.height + "px", this.avatarEl.src = this.avatarSrc, this.avatarEl.id = "zm_chat-wiki-avatar", this.enableDrag(), document.body.appendChild(this.avatarEl));
  }
  handleClick(t) {
    this.clickHandler && this.clickHandler();
  }
  removeAvatar() {
    this.avatarEl && document.body.removeChild(this.avatarEl);
  }
  show() {
    this.avatarEl.style.display = "block";
  }
  hide() {
    this.avatarEl.style.display = "none";
  }
}
function g() {
  let a = {
    iframeSrc: "https://zhima_chat_ai.applnk.cn/chat-ai-pc/#/chat",
    remote: "",
    params: {}
  };
  const t = document.getElementById("ai_chat_js");
  if (t) {
    let s = t.getAttribute("data-json"), r = new URL(t.src).origin;
    a.iframeSrc = r + "/web/#/chat";
    try {
      a.params = JSON.parse(s);
    } catch (o) {
      console.error("Failed to stringify data:", o);
      return;
    }
  }
  const i = new m(a), n = new f(a);
  i.onInit = (s) => {
    n.init(s);
  }, i.onClose = () => {
    n.show();
  }, n.clickHandler = () => {
    i.open(), n.hide();
  };
}
d();
window.AiChatSDK = g();
