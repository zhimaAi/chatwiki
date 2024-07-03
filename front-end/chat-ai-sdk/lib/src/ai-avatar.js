class AiAvatar {
  clickHandler = null;
  avatarEl = null;
  avatarSrc = "";
  left = 0;
  top = 0;
  width = 50;
  height = 50;

  initialX = 0;
  initialY = 0;
  initialMouseX = 0;
  initialMouseY = 0;

  // 拖拽状态标志位
  dragging = false;

  constructor(config) {
    if (config) {
      this.avatarSrc = config.avatarSrc;
    }
  }

  init(config) {
    this.avatarSrc = config.sdkFloatAvatar;

    this.getInitialPosition();
    this.insertAvatar();
  }

  // 获取初始位置
  getInitialPosition() {
    const winWidth = window.innerWidth;
    const winHeight = window.innerHeight;

    this.left = winWidth - this.width - 50;
    this.top = winHeight - this.height - 50;
  }

  // 启用拖拽
  enableDrag() {
    this.avatarEl.addEventListener("mousedown", (e) => {
      this.initialX = this.left;
      this.initialY = this.top;

      // 获取鼠标相对于元素的初始位置
      this.initialMouseX =
        e.clientX - this.avatarEl.getBoundingClientRect().left;
      this.initialMouseY =
        e.clientY - this.avatarEl.getBoundingClientRect().top;

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

    // 更新图片位置
    this.avatarEl.style.left = this.left + "px";
    this.avatarEl.style.top = this.top + "px";
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

    this.avatarEl = document.createElement("img");

    this.avatarEl.style.display = "block";
    this.avatarEl.style.top = this.top + "px";
    this.avatarEl.style.left = this.left + "px";
    this.avatarEl.style.width = this.width + "px";
    this.avatarEl.style.height = this.height + "px";
    this.avatarEl.src = this.avatarSrc;
    this.avatarEl.id = "zm_chat-wiki-avatar";

    // this.avatarEl.addEventListener("click", this.handleClick.bind(this));

    this.enableDrag();

    document.body.appendChild(this.avatarEl);
  }

  handleClick(e) {
    if (this.clickHandler) {
      this.clickHandler();
    }
  }

  removeAvatar() {
    if (this.avatarEl) {
      // this.avatarEl.removeEventListener("click", this.handleClick.bind(this));

      document.body.removeChild(this.avatarEl);
    }
  }

  show() {
    this.avatarEl.style.display = "block";
  }

  hide() {
    this.avatarEl.style.display = "none";
  }
}

export default AiAvatar;
