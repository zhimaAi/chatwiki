import { createApp, h } from 'vue'
import ContextPadView from "./index.vue";

const WIDTH = 130
const HEIGHT = 440

export class ContextPad {
  constructor({ lf }) {
    this.lf = lf;
    this.__menuDOM = document.createElement("div");
    this.__menuDOM.className = "lf-inner-context";
    this.__menuDOM.style.width = `${WIDTH}px`;
    this.__menuDOM.style.height = `${HEIGHT}px`;
  }

  render(lf, container) {
    this.container = container;

    lf.on("custom:showPopupMenu", ({model}) => {
      this.hideMenu();
    })

    lf.on("node:click", ({ data }) => {
      this.hideMenu();
    });

    lf.on("edge:click", ({ data }) => {
      // 获取右上角坐标
      this._activeData = data;
      this.createMenu();
    });

    lf.on("blank:click", () => {
      this.hideMenu();
    });
  }
  createMenu() {
    // this.__menuDOM = document.createElement("div");
    // this.__menuDOM.className = "lf-inner-context context-pad-wrapper";


    this.showMenu();
  }
  // 计算出菜单应该显示的位置（节点的右上角）
  getContextMenuPosition() {
    const data = this._activeAnchorData;
    let x = data.x + 5;
    let y =  data.y - 10;
    // if (Model.BaseType === "node") {
    // }
    return this.lf.graphModel.transformModel.CanvasPointToHtmlPoint([x, y]);
  }
  showMenu() {
    const [x, y] = this.getContextMenuPosition();
    this.__menuDOM.style.display = "flex";
    // 将菜单显示到对应的位置
    this.__menuDOM.style.top = `${y}px`;
    this.__menuDOM.style.left = `${x + 10}px`;

    // 菜单显示的时候，监听删除，同时隐藏
    !this.isShow && this.lf.on("node:delete,edge:delete,node:drag,graph:transform", this.listenDelete);

    this.isShow = true;
    // this.container.appendChild(this.__menuDOM);
  }
  /**
   * 隐藏菜单
   */
  hideMenu () {
    // this.__menuDOM.innerHTML = "";
    this.__menuDOM.style.display = "none";

    if (this.isShow) {
      this.container.removeChild(this.__menuDOM);
    }

    this.lf.off("node:delete,edge:delete,node:drag,graph:transform",  this.listenDelete);

    this.isShow = false;
  }

  listenDelete = () => {
    this.hideMenu();
  }
}
