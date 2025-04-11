import createAvatar from "./create-avater";
import AiChatWidget from './ai-chat'
class AiAvatar {
  avatarElWrapper = null;
  avatarContentEl = null;
  avatarEl = null;
  click = null;
  
  left = 0;
  top = 0;
  right = 0;
  bottom = 0;
  width = 50;
  height = 50;

  initialX = 0;
  initialY = 0;
  initialMouseX = 0;
  initialMouseY = 0;

  // 拖拽状态标志位
  dragging = false;

  config = {
    displayType: 1,
    buttonText: "快来聊聊吧~",
    buttonIcon: "",
    bottomMargin: 32,
    rightMargin: 32,
    showUnreadCount: 1,
    showNewMessageTip: 1,
  };

  constructor(config) {
   
  }

  init(data) {
    const { config } = data;
  
    this.config = config.floatBtn;

    this.insertAvatar();
  }

  // 设置初始位置
  setInitialPosition() {
    const winWidth = window.innerWidth;
    const winHeight = window.innerHeight;

    this.top = winHeight - this.height - this.config.bottomMargin * 1;
    this.left = winWidth - this.width - this.config.rightMargin * 1;
    this.bottom = this.config.bottomMargin * 1;
    this.right = this.config.rightMargin * 1;

    this.updataPosition();
  }

  updataPosition() {
    this.avatarElWrapper.style.left = this.left + "px";
    this.avatarElWrapper.style.top = this.top + "px";
  }

  onWindowResize(){
    const winWidth = window.innerWidth;
    const winHeight = window.innerHeight;

    let top = winHeight - this.height - this.bottom;
    let left = winWidth - this.width - this.right;

    this.top = Math.max(0, top);
    this.left = Math.max(0, left);
   

    this.updataPosition();
  }
  // 启用拖拽
  enableDrag() {
    // 监听窗口大小事件并更新位置
    window.addEventListener("resize", this.onWindowResize.bind(this));
    // 监听mousedown事件
    this.avatarEl.addEventListener("mousedown", (e) => {
      this.initialX = this.left;
      this.initialY = this.top;

      // 获取鼠标相对于元素的初始位置
      this.initialMouseX =
        e.clientX - this.avatarElWrapper.getBoundingClientRect().left;
      this.initialMouseY =
        e.clientY - this.avatarElWrapper.getBoundingClientRect().top;

      // 添加对全局的mousemove和mouseup监听
      document.addEventListener("mousemove", this.handleDrag);
      document.addEventListener("mouseup", this.handleDragEnd);

      // 阻止默认行为
      e.preventDefault();
    });
  }

  // 拖拽移动
  handleDrag = (e) => {
    // 标记拖拽开始
    this.dragging = true;
    // 计算新的位置
    const newX = e.clientX - this.initialMouseX;
    const newY = e.clientY - this.initialMouseY;

    // 限制新的位置在视口内
    const winWidth = window.innerWidth;
    const winHeight = window.innerHeight;

    // 更新left和top
    this.left = Math.max(0, Math.min(newX, winWidth - this.width));
    this.top = Math.max(0, Math.min(newY, winHeight - this.height));

    this.right = winWidth - this.left - this.width;
    this.bottom = winHeight - this.top - this.height;
 
    // 更新图片位置
    this.updataPosition();
  };

  // 拖拽结束
  handleDragEnd = () => {
    const deltaX = Math.abs(this.initialX - this.left);
    const deltaY = Math.abs(this.initialY - this.top);

    if (deltaX <= 3 && deltaY <= 3) {
      this.handleClick();
    }

    // 标记拖拽结束
    this.dragging = false;

    document.removeEventListener("mousemove", this.handleDrag);
    document.removeEventListener("mouseup", this.handleDragEnd);
  };

  insertAvatar() {
    if (document.getElementById("zm_chat-wiki-avatar")) {
      return;
    }
    this.avatarElWrapper = document.createElement("div");
    this.avatarElWrapper.className = "zm_chat-wiki-avatar-wrapper";
    this.avatarElWrapper.style.display = "block";
    this.avatarElWrapper.style.left = "-99999px";
    this.avatarElWrapper.style.top = "-99999px";

    this.avatarContentEl = document.createElement("div");
    this.avatarContentEl.className = "zm_chat-wiki-avatar-content";

    this.avatarEl = document.createElement("div");
    this.avatarEl.id = "zm_chat-wiki-avatar";

    this.avatarContentEl.appendChild(this.avatarEl);
    
    this.avatarElWrapper.appendChild(this.avatarContentEl);
    
    

    document.body.appendChild(this.avatarElWrapper);
    
    // this.avatarEl.addEventListener("click", this.handleClick.bind(this));

    createAvatar(this.config).then((dom) => {
      this.avatarEl.appendChild(dom);
      
      this.width = this.avatarEl.offsetWidth;
      this.height = this.avatarEl.offsetHeight;

      this.setInitialPosition();

      this.enableDrag();
      
    }).catch(e => {
      console.error('图片加载失败');
    })
  }

  handleClick(e) {
    this.hide();
    AiChatWidget.open();
  }

  removeAvatar() {
    if (this.avatarElWrapper) {
      // this.avatarEl.removeEventListener("click", this.handleClick.bind(this));

      document.body.removeChild(this.avatarElWrapper);
    }
  }

  show() {
    this.avatarElWrapper.style.display = "block";
  }

  hide() {
    this.avatarElWrapper.style.display = "none";
  }
}

export default new AiAvatar();
