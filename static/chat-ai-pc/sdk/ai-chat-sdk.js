var d = Object.defineProperty;
var m = (i, t, e) => t in i ? d(i, t, { enumerable: !0, configurable: !0, writable: !0, value: e }) : i[t] = e;
var s = (i, t, e) => (m(i, typeof t != "symbol" ? t + "" : t, e), e);
function p(i) {
  if (!i)
    return "";
  const t = [];
  for (let e in i)
    if (i.hasOwnProperty(e)) {
      let a = i[e];
      (Array.isArray(a) || typeof a == "object" && a !== null) && (a = JSON.stringify(a)), a !== void 0 && t.push(encodeURIComponent(e) + "=" + encodeURIComponent(a));
    }
  return t.join("&");
}
function g() {
  const i = document.getElementById("ai_chat_js"), t = document.createElement("link"), e = new URL(i.src).origin;
  t.type = "text/css", t.rel = "stylesheet", t.href = e + "/sdk/style.css", document.getElementsByTagName("head")[0].appendChild(t);
}
function u(i) {
  return new Promise((t, e) => {
    let a = new Image();
    a.src = i, a.onload = () => t(a), a.onerror = e;
  });
}
async function f(i) {
  let t = document.createElement("img");
  return t.src = i.buttonIcon, t.style.width = "50px", t.style.height = "50px", t;
}
async function v(i) {
  let t = document.createElement("div");
  t.className = "chat-wiki-avatar_type2";
  let e = document.createElement("img");
  e.src = i.buttonIcon, e.className = "chat-wiki-avatar_type2_icon", t.appendChild(e);
  let a = document.createElement("div");
  return a.className = "chat-wiki-avatar_type2_text", a.innerText = i.buttonText, t.appendChild(a), t;
}
async function w(i) {
  try {
    const t = await u(i.buttonIcon);
    return t.style.width = t.width + "px", t.style.height = t.height + "px", t;
  } catch (t) {
    console.log(t);
  }
}
async function E(i) {
  return i.displayType === 3 ? w(i) : i.displayType === 2 ? v(i) : f(i);
}
class y {
  constructor(t) {
    s(this, "avatarElWrapper", null);
    s(this, "avatarContentEl", null);
    s(this, "avatarEl", null);
    s(this, "click", null);
    s(this, "left", 0);
    s(this, "top", 0);
    s(this, "right", 0);
    s(this, "bottom", 0);
    s(this, "width", 50);
    s(this, "height", 50);
    s(this, "initialX", 0);
    s(this, "initialY", 0);
    s(this, "initialMouseX", 0);
    s(this, "initialMouseY", 0);
    // 拖拽状态标志位
    s(this, "dragging", !1);
    s(this, "config", {
      displayType: 1,
      buttonText: "快来聊聊吧~",
      buttonIcon: "",
      bottomMargin: 32,
      rightMargin: 32,
      showUnreadCount: 1,
      showNewMessageTip: 1
    });
    // 拖拽移动
    s(this, "handleDrag", (t) => {
      this.dragging = !0;
      const e = t.clientX - this.initialMouseX, a = t.clientY - this.initialMouseY, n = window.innerWidth, r = window.innerHeight;
      this.left = Math.max(0, Math.min(e, n - this.width)), this.top = Math.max(0, Math.min(a, r - this.height)), this.right = n - this.left - this.width, this.bottom = r - this.top - this.height, this.updataPosition();
    });
    // 拖拽结束
    s(this, "handleDragEnd", () => {
      const t = Math.abs(this.initialX - this.left), e = Math.abs(this.initialY - this.top);
      t <= 3 && e <= 3 && this.handleClick(), this.dragging = !1, document.removeEventListener("mousemove", this.handleDrag), document.removeEventListener("mouseup", this.handleDragEnd);
    });
  }
  init(t) {
    const { config: e } = t;
    this.config = e.floatBtn, this.insertAvatar();
  }
  // 设置初始位置
  setInitialPosition() {
    const t = window.innerWidth, e = window.innerHeight;
    this.top = e - this.height - this.config.bottomMargin * 1, this.left = t - this.width - this.config.rightMargin * 1, this.bottom = this.config.bottomMargin * 1, this.right = this.config.rightMargin * 1, this.updataPosition();
  }
  updataPosition() {
    this.avatarElWrapper.style.left = this.left + "px", this.avatarElWrapper.style.top = this.top + "px";
  }
  onWindowResize() {
    const t = window.innerWidth;
    let a = window.innerHeight - this.height - this.bottom, n = t - this.width - this.right;
    this.top = Math.max(0, a), this.left = Math.max(0, n), this.updataPosition();
  }
  // 启用拖拽
  enableDrag() {
    window.addEventListener("resize", this.onWindowResize.bind(this)), this.avatarEl.addEventListener("mousedown", (t) => {
      this.initialX = this.left, this.initialY = this.top, this.initialMouseX = t.clientX - this.avatarElWrapper.getBoundingClientRect().left, this.initialMouseY = t.clientY - this.avatarElWrapper.getBoundingClientRect().top, document.addEventListener("mousemove", this.handleDrag), document.addEventListener("mouseup", this.handleDragEnd), t.preventDefault();
    });
  }
  insertAvatar() {
    document.getElementById("zm_chat-wiki-avatar") || (this.avatarElWrapper = document.createElement("div"), this.avatarElWrapper.className = "zm_chat-wiki-avatar-wrapper", this.avatarElWrapper.style.display = "block", this.avatarElWrapper.style.left = "-99999px", this.avatarElWrapper.style.top = "-99999px", this.avatarContentEl = document.createElement("div"), this.avatarContentEl.className = "zm_chat-wiki-avatar-content", this.avatarEl = document.createElement("div"), this.avatarEl.id = "zm_chat-wiki-avatar", this.avatarContentEl.appendChild(this.avatarEl), this.avatarElWrapper.appendChild(this.avatarContentEl), document.body.appendChild(this.avatarElWrapper), E(this.config).then((t) => {
      this.avatarEl.appendChild(t), this.width = this.avatarEl.offsetWidth, this.height = this.avatarEl.offsetHeight, this.setInitialPosition(), this.enableDrag();
    }).catch((t) => {
      console.error("图片加载失败");
    }));
  }
  handleClick(t) {
    this.hide(), c.open();
  }
  removeAvatar() {
    this.avatarElWrapper && document.body.removeChild(this.avatarElWrapper);
  }
  show() {
    this.avatarElWrapper.style.display = "block";
  }
  hide() {
    this.avatarElWrapper.style.display = "none";
  }
}
const o = new y();
class W {
  constructor() {
    s(this, "dotEl", null);
  }
  create(t) {
    t.value == 0 && !this.dotEl || (t.value > 0 && !this.dotEl && (this.dotEl = document.createElement("div"), o.avatarContentEl.appendChild(this.dotEl)), this.dotEl && t.value <= 0 ? (this.dotEl.remove(), this.dotEl = null) : (this.dotEl.className = `ai-dot${Number(t.value) > 9 ? " ai-dot-plus" : ""}`, this.dotEl.textContent = t.value));
  }
}
const M = new W(), b = `<svg fill="none" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" class="design-iconfont">
  <path d="M5 5L11 11" stroke="#595959" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/>
  <path d="M5 11L11 5" stroke="#595959" stroke-width="1.3" stroke-linecap="round" stroke-linejoin="round"/>
</svg>`, C = {
  1: {
    value: "text",
    label: "text"
  },
  2: {
    value: "menu",
    label: "[菜单]"
  },
  3: {
    value: "image",
    label: "[图片]"
  }
};
class k {
  constructor() {
    s(this, "listDomWrapper", null);
    s(this, "timer", null);
  }
  create(t) {
    if (t.length === 0 && this.listDomWrapper) {
      this.remove();
      return;
    }
    this.listDomWrapper && this.remove(), this.listDomWrapper = document.createElement("div"), this.listDomWrapper.className = "new-message-list-wrapper";
    let e = document.createElement("div");
    e.className = "new-message-list";
    for (let n = 0; n < t.length; n++) {
      let r = t[n], l = this.getMessageDom(r);
      e.innerHTML += l;
    }
    this.listDomWrapper.appendChild(e), o.avatarContentEl.appendChild(this.listDomWrapper), this.listDomWrapper.querySelectorAll(".close-btn").forEach((n) => {
      n.addEventListener("click", function(r) {
        r.stopPropagation();
        const l = this.closest(".message-item");
        l && l.remove();
      });
    }), this.timer && (clearTimeout(this.timer), this.timer = null), this.timer = setTimeout(() => {
      this.remove();
    }, 5 * 1e3);
  }
  getMessageDom(t) {
    let e = t.content;
    return (t.msg_type == 2 || t.msg_type == 3) && (e = C[t.msg_type].label), `<div class="message-item">
      <div class="ai-assistant">
        <span class="close-btn">${b}</span>
        <div class="message-header">
          <img class="ai-icon" src="${t.avatar}" />
          <span class="ai-name">${t.robot_name}</span>
        </div>
        <div class="message-content">${e}</div>
      </div>
    </div>`;
  }
  remove() {
    this.listDomWrapper && (this.listDomWrapper.remove(), this.listDomWrapper = null);
  }
}
const h = new k();
class D {
  constructor() {
    s(this, "iframe", null);
    s(this, "iframeSrc", null);
    s(this, "config", {});
  }
  init(t) {
    return this.config = t, this.iframeSrc = this.config.iframeSrc, this.insertAiChat(), this.getExpectedOrigin(), window.addEventListener("message", this.handleMessage.bind(this), !1), {
      open: this.open.bind(this),
      close: this.close.bind(this)
    };
  }
  removeAiChat() {
    const t = document.getElementById("zm_chat-wiki-iframe");
    t && (document.body.removeChild(t), this.iframe = null);
  }
  onInit(t) {
    o.init(t);
  }
  open() {
    this.iframe.style.display = "block", this.postMessage("openWindow", {});
  }
  close() {
    o.show(), this.iframe.style.display = "none", this.postMessage("closeWindow", {});
  }
  onClose() {
    o.show(), this.iframe.style.display = "none";
  }
  createDot(t) {
    M.create({ value: t });
  }
  createNewMessage(t) {
    let e = t || [];
    if (e.length === 0) {
      h.remove();
      return;
    }
    e = e.slice(-1), h.create(e);
  }
  getExpectedOrigin() {
    return new URL(this.iframeSrc).origin;
  }
  handleMessage(t) {
    if (!t.origin || !(new URL(t.origin).origin === this.getExpectedOrigin()))
      return;
    let n = t.data;
    n && (n.action === "closeChat" && this.onClose(), n.action === "init" && this.onInit(n.data), n.action === "dot" && this.createDot(n.data), n.action === "newMessage" && this.createNewMessage(n.data));
  }
  postMessage(t, e) {
    if (e)
      try {
        e = JSON.parse(JSON.stringify(e));
      } catch (a) {
        console.error("Failed to stringify data:", a);
        return;
      }
    if (this.iframe.contentWindow && typeof this.iframe.contentWindow.postMessage == "function")
      try {
        this.iframe.contentWindow.postMessage({ action: t, data: e }, "*");
      } catch (a) {
        console.error("Failed to post message:", a);
      }
    else
      console.warn(
        "frame.contentWindow is not available or postMessage is not supported."
      );
  }
  insertAiChat() {
    if (document.getElementById("zm_chat-wiki-iframe"))
      return;
    const t = p(this.config.params);
    this.iframe = document.createElement("iframe"), this.iframe.id = "zm_chat-wiki-iframe", this.iframe.src = this.iframeSrc + "?" + t, this.iframe.style.display = "none", document.body.appendChild(this.iframe);
  }
}
const c = new D();
function x() {
  let i = {
    iframeSrc: "/#/chat",
    remote: "",
    params: {}
  };
  const t = document.getElementById("ai_chat_js");
  if (t) {
    let e = t.getAttribute("data-json"), a = new URL(t.src).origin;
    i.iframeSrc = a + "/web/#/chat";
    try {
      i.params = JSON.parse(e);
    } catch (n) {
      console.error("Failed to stringify data:", n);
      return;
    }
  }
  return c.init(i);
}
g();
window.AiChatSDK = x();
