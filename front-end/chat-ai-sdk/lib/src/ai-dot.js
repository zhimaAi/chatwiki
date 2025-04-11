import AiAvatar from "./ai-avatar";

class AiDot {
  dotEl = null;

  create(config) {
    if(config.value == 0 && !this.dotEl){
      return;
    }

    if(config.value > 0 && !this.dotEl){
      this.dotEl = document.createElement("div");
      AiAvatar.avatarContentEl.appendChild(this.dotEl);
    }

    if(this.dotEl && config.value <= 0){
      // 移除dotEl
      this.dotEl.remove();
      this.dotEl = null;
    }else{
      // 如果dotEl存在则更新它的文本内容
      this.dotEl.className = `ai-dot${Number(config.value) > 9 ? ' ai-dot-plus' : ''}`;
      this.dotEl.textContent = config.value;
    }
  }
}

export default new AiDot();
