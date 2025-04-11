import AiAvatar from "./ai-avatar";
import { closeIcon } from "./icon";

const msgTypeMap = {
  1: {
    value: "text",
    label: "text",
  },
  2: {
    value: "menu",
    label: "[菜单]",
  },
  3: {
    value: "image",
    label: "[图片]",
  },
};

class NewMssage {
  listDomWrapper = null;
  timer = null;

  create(list) {
    if (list.length === 0 && this.listDomWrapper) {
      this.remove();
      return;
    }
    
    if (this.listDomWrapper) {
      this.remove();
    }

    this.listDomWrapper = document.createElement("div");
    this.listDomWrapper.className = "new-message-list-wrapper";

    let listDom = document.createElement("div");
    listDom.className = "new-message-list";

    // 遍历list，生成dom
    for (let i = 0; i < list.length; i++) {
      let msg = list[i];
      let msgDom = this.getMessageDom(msg);
      listDom.innerHTML += msgDom;
    }

    this.listDomWrapper.appendChild(listDom);

    AiAvatar.avatarContentEl.appendChild(this.listDomWrapper);

    // 添加关闭按钮点击事件
    const closeBtns = this.listDomWrapper.querySelectorAll('.close-btn');
    closeBtns.forEach(btn => {
      btn.addEventListener('click', function(e) {
        e.stopPropagation(); // 阻止事件冒泡
        const messageItem = this.closest('.message-item');
        if (messageItem) {
          messageItem.remove();
        }
      });
    });

    if(this.timer){
      clearTimeout(this.timer)
      this.timer = null
    }

    this.timer = setTimeout(() => {
      this.remove();
    }, 5 * 1000);
  }

  getMessageDom(msg) {
    let content = msg.content;
    if (msg.msg_type == 2 || msg.msg_type == 3) {
      content = msgTypeMap[msg.msg_type].label;
    }

    let msgDom = `<div class="message-item">
      <div class="ai-assistant">
        <span class="close-btn">${closeIcon}</span>
        <div class="message-header">
          <img class="ai-icon" src="${msg.avatar}" />
          <span class="ai-name">${msg.robot_name}</span>
        </div>
        <div class="message-content">${content}</div>
      </div>
    </div>`;
    return msgDom;
  }

  remove() {
    if (this.listDomWrapper) {
      this.listDomWrapper.remove();
      this.listDomWrapper = null;
    }
  }
}

export default new NewMssage();
