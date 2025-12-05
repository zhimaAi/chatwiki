export class GroupAutoResize {
  static pluginName = 'groupAutoResize';

  constructor({ lf }) {
    this.lf = lf;
    this.padding = { top: 94, right: 16, bottom: 16, left: 16 };
    // resize_padding 用于在调整大小时，在子节点和父节点之间保留一些额外的空间
    this.resize_padding = 30;

    this.handleNodeDrag = this.handleNodeDrag.bind(this);
    lf.on('node:drag', this.handleNodeDrag);
  }

  handleNodeDrag({ data }) {
    // 1. 解绑事件，避免在调用API时触发循环
    this.lf.off('node:drag', this.handleNodeDrag);

    try {
      const model = this.lf.graphModel.getNodeModelById(data.id);
      if (!model) return;

      const { x, y, width, height } = model;
      const props = model.getProperties();

      if (props.loop_parent_key) {
        const groupNode = this.lf.graphModel.getNodeModelById(props.loop_parent_key);
        if (groupNode) {
          const { x: groupX, y: groupY, width: groupWidth, height: groupHeight } = groupNode;

          // 子节点的四条边
          const nodeLeft = x - width / 2;
          const nodeRight = x + width / 2;
          const nodeTop = y - height / 2;
          const nodeBottom = y + height / 2;

          // 父容器的四条边（考虑内边距）
          const groupLeft = groupX - groupWidth / 2 + this.padding.left;
          const groupRight = groupX + groupWidth / 2 - this.padding.right;
          const groupTop = groupY - groupHeight / 2 + this.padding.top;
          const groupBottom = groupY + groupHeight / 2 - this.padding.bottom;

          let newGroupWidth = groupWidth;
          let newGroupHeight = groupHeight;
          let newGroupX = groupX;
          let newGroupY = groupY;
          let needsResize = false;

          // 2. 判断并计算新的尺寸和位置
          if (nodeLeft < groupLeft) {
            const diff = groupLeft - nodeLeft + this.resize_padding;
            newGroupWidth += diff;
            newGroupX -= diff / 2; // 中心点向左移动
            needsResize = true;
          }
          if (nodeRight > groupRight) {
            const diff = nodeRight - groupRight + this.resize_padding;
            newGroupWidth += diff;
            newGroupX += diff / 2; // 中心点向右移动
            needsResize = true;
          }
          if (nodeTop < groupTop) {
            const diff = groupTop - nodeTop + this.resize_padding;
            newGroupHeight += diff;
            newGroupY -= diff / 2; // 中心点向上移动
            needsResize = true;
          }
          if (nodeBottom > groupBottom) {
            const diff = nodeBottom - groupBottom + this.resize_padding;
            newGroupHeight += diff;
            newGroupY += diff / 2; // 中心点向下移动
            needsResize = true;
          }

          // 3. 如果需要，则调用 setAttributes 更新节点属性
          if (needsResize) {
            groupNode.x = newGroupX
            groupNode.y = newGroupY
            groupNode._width = groupNode.width = newGroupWidth
            groupNode._height = groupNode.height = newGroupHeight
          }
        }
      }
    } finally {
      // 3. 在finally块中重新绑定事件，确保无论是否出错都能恢复监听
      this.lf.on('node:drag', this.handleNodeDrag);
    }
  }
}