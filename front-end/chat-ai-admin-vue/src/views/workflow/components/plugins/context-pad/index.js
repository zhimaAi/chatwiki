import { createApp, h } from 'vue'
import ContextPadView from "./index.vue";

const WIDTH = 444
const HEIGHT = 620;

export class ContextPad {
  constructor({ lf }) {
    this.menuTypeMap = new Map();
    this.lf = lf;
    this.__menuDOM = document.createElement("div");
    this.__menuDOM.className = "lf-inner-context";
    this.__menuDOM.style.width = `${WIDTH}px`;
    this.__menuDOM.style.height = `${HEIGHT}px`;


    this.isMounted = false

    this.r = h(ContextPadView, {
      properties: {},
      isSelected: false,
      isHovered:false,
      onClickItem: ({node, event}) => {
        this.hideMenu()
        
        lf.graphModel.eventCenter.emit('custom:addNode', {
          node: node,
          event: event,
          model: this._activeData,
          anchorData: this._activeAnchorData,
        })
      },
    })

    this.app = createApp({
      render: () => this.r,
      provide: () => ({
        getNode: () => this._activeData,
        getGraph: () => lf.graphModel,
      }),
      mounted: () => {},
    })
  }

  render(lf, container) {
    this.container = container;

    if (!this.isMounted) {
      this.isMounted = true
      this.container.appendChild(this.__menuDOM);
      this.app.mount(this.__menuDOM)
    } else {
      // this.r.component.props.properties = this.props.model.getProperties()
    }

    lf.on("custom:showPopupMenu", ({anchorData, model}) => {
      this._activeData = model;
      this._activeAnchorData = anchorData;
      this.createMenu();
    })

    lf.on("node:click", ({ data }) => {
      // this._activeData = data;
      // this.createMenu();
      this.hideMenu();
    });
    lf.on("edge:click", ({ data }) => {
      // 获取右上角坐标
      // this._activeData = data;
      // this.createMenu();
      this.hideMenu();
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

    return this.lf.graphModel.transformModel.CanvasPointToHtmlPoint([x, y]);
  }
  showMenu() {
    // 获取菜单的理想位置
    const [x, y] = this.getContextMenuPosition();
    // 获取画布的高度和宽度，减去适当偏移量
    const canvasHeight = this.lf.graphModel.height - 64;
    const canvasWidth = this.lf.graphModel.width - 64;

    // 显示菜单
    this.__menuDOM.style.display = "flex";
    
    // 计算菜单的实际高度，取525和（画布高度-20）中的较小值
    const menuHeight = Math.min(525, canvasHeight - 20);
    const menuWidth = WIDTH;

    // 处理水平位置
    let left = x + 10;
    
    // 判断菜单是否会超出画布右侧
    if (left + menuWidth > canvasWidth) {
      // 如果超出，将菜单位置调整到锚点左侧
      left = x - menuWidth - 20;
    }
    
    // 如果菜单左侧超出画布左侧，则设置为4像素
    if (left < 0) {
      left = 4;
    }
    
    // 设置菜单的水平位置
    this.__menuDOM.style.left = `${left}px`;

    let top = y;
    
    // 判断菜单是否会超出画布底部
    if(y + menuHeight > canvasHeight){
      // 如果超出，向上调整菜单位置，使其紧贴画布底部
      top = y - (y + menuHeight - canvasHeight) - 10
    }else{
      // 如果未超出，向上微调10像素
      top = y - 10
    }

    // 如果菜单顶部超出画布顶部，则设置为4像素
    if(top < 0){
      top = 4
    }
    
    // 设置菜单的垂直位置
    this.__menuDOM.style.top = `${top}px`;
    // 设置菜单的高度，如果计算出的高度有效则使用，否则使用默认高度
    this.__menuDOM.style.height = `${menuHeight > 0 ? menuHeight : HEIGHT}px`;

    // 如果菜单当前未显示，则添加事件监听器，在节点删除、拖拽或画布变换时隐藏菜单
    !this.isShow && this.lf.on("node:delete,edge:delete,node:drag,graph:transform", this.listenDelete);

    // 标记菜单为已显示
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
      // this.container.removeChild(this.__menuDOM);
    }

    this.lf.off("node:delete,edge:delete,node:drag,graph:transform",  this.listenDelete);

    this.isShow = false;
  }

  listenDelete = () => {
    this.hideMenu();
  }
}
