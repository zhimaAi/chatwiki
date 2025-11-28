// 贝塞尔曲线
import { BezierEdge, BezierEdgeModel, h } from "@logicflow/core";
import { createApp, h as vh } from 'vue'
import VueNode from './index.vue'
import { generateUniqueId } from '@/utils/index'

class CustomEdge extends BezierEdge {
  constructor(props) {
    super(props)
    this.isMounted = false
    this.r = vh(VueNode, {
      properties: props.model.getProperties(),
      isSelected: props.model.isSelected,
      isHovered: props.model.isHovered,
    })

    this.app = createApp({
      render: () => this.r,
      provide: () => ({
        getModel: () => props.model,
        getGraphModel: () => props.graphModel,
      }),
      mounted: () => {
        // console.log('mounted', this)
      },
    })
  }

  // 关键修复：重写 view 层的 getEndArrow 方法，手动渲染箭头
  getEndArrow() {
    const { model } = this.props;
    const { stroke, strokeWidth } = model.getArrowStyle();

    return (
      // d格式：M 尾部X 尾部Y L 尖端X 尖端Y L 尾部另一侧X 尾部另一侧Y Z
      h('path', { 
        stroke,
        strokeWidth,
        d: 'M -9 -7 0 0 -9 7',
        fill:'none'
      })
    )
  }

  getText() {
    const { model } = this.props;
    const { startPoint, endPoint, pointsList } = model; // pointsList长度为4，即两个控制点

    const id = model.id;
    // 固定元素尺寸为24×24
    const elementWidth = 24;
    const elementHeight = 24;
    const halfWidth = elementWidth / 2; // 12
    const halfHeight = elementHeight / 2; // 12

    // 1. 三次贝塞尔曲线点计算
    const getBezierPoint = (t) => {
      const { x: x0, y: y0 } = startPoint;
      const { x: x1, y: y1 } = pointsList[1]; // 第一个控制点
      const { x: x2, y: y2 } = pointsList[2]; // 第二个控制点
      const { x: x3, y: y3 } = endPoint;
      const t2 = t * t;
      const t3 = t2 * t;
      const mt = 1 - t;
      const mt2 = mt * mt;
      const mt3 = mt2 * mt;
      // 三次贝塞尔公式：B(t) = P0*(1-t)³ + 3P1*(1-t)²t + 3P2*(1-t)t² + P3*t³
      return {
        x: x0 * mt3 + 3 * x1 * mt2 * t + 3 * x2 * mt * t2 + x3 * t3,
        y: y0 * mt3 + 3 * y1 * mt2 * t + 3 * y2 * mt * t2 + y3 * t3,
      };
    };

    // 2. 高精度曲线长度计算（增加分段数）
    const getBezierLength = (segments = 5) => { // segments 增加精度
      let length = 0;
      let prevPoint = getBezierPoint(0);
      for (let i = 1; i <= segments; i++) {
        const t = i / segments;
        const currentPoint = getBezierPoint(t);
        const dx = currentPoint.x - prevPoint.x;
        const dy = currentPoint.y - prevPoint.y;
        length += Math.sqrt(dx * dx + dy * dy);
        prevPoint = currentPoint;
      }
      return length;
    };

    // 3. 二分法精准找中点（替代分段遍历）
    const findMidPoint = (totalLength) => {
      const targetLength = totalLength / 2;
      let low = 0, high = 1;
      let midT = 0.5;
      let epsilon = 0.0001; // 精度阈值

      while (high - low > epsilon) {
        midT = (low + high) / 2;
        const currentLength = getPartialLength(midT);
        if (currentLength < targetLength) {
          low = midT;
        } else {
          high = midT;
        }
      }

      // 辅助函数：计算t从0到t的曲线长度
      function getPartialLength(t) {
        let length = 0;
        let prevPoint = getBezierPoint(0);
        const segments = 5; // 内部分段数
        for (let i = 1; i <= segments; i++) {
          const currentT = i / segments * t;
          const currentPoint = getBezierPoint(currentT);
          const dx = currentPoint.x - prevPoint.x;
          const dy = currentPoint.y - prevPoint.y;
          length += Math.sqrt(dx * dx + dy * dy);
          prevPoint = currentPoint;
        }
        return length;
      }

      return getBezierPoint(midT);
    };

    // 4. 计算中点并定位元素
    const totalLength = getBezierLength();
    const midPoint = findMidPoint(totalLength);

    const positionData = {
      x: midPoint.x - halfWidth,
      y: midPoint.y - halfHeight,
      width: elementWidth,
      height: elementHeight,
    };

    const wrapperStyle = {
      width: elementWidth,
      height: elementHeight,
    };

    // 5. 挂载/更新组件
    setTimeout(() => {
      const dom = document.querySelector(`#line_${id}`)?.querySelector('.vue-app-wrapper');
      if (!this.isMounted && dom) {
        this.isMounted = true;
        this.app.mount(dom);
      } else if (this.r?.component) {
        this.r.component.props.isSelected = model.isSelected;
        this.r.component.props.isHovered = model.isHovered;
        this.r.component.props.properties = model.getProperties();
      }
    }, 0);

    return h('foreignObject', {
      ...positionData,
      id: `line_${id}`,
    }, [
      h('div', {
        id,
        style: wrapperStyle,
        className: 'vue-app-wrapper',
      }),
    ]);
  }
}

class CustomEdgeModel extends BezierEdgeModel {
  initEdgeData(data) {
    super.initEdgeData(data)
  }
  setAttributes(data) {
    super.setAttributes(data)
    // 节点拖拽时动态调整边的offset
    let offset = this.getProperties().offset;
    if(offset){
      this.offset = offset
      this.pointsList = []
      this.initPoints()
    }
  }
  createId() {
    return generateUniqueId('customEdge');
  }

  getEdgeStyle() {
    const style = super.getEdgeStyle()
    // svg属性
    if (this.isSelected || this.isHovered) {
      style.stroke = '#37e1ffff'
    } else {
      style.stroke = '#2475FC'
    }

    style.strokeWidth = 2

    return style
  }

  getArrowStyle() {
    const style = super.getArrowStyle()
    return {
      ...style,
      ...this.arrowStyle,
      refX: 0,
      refY: 0,
    };
  }

  getOutlineStyle() {
    const style = super.getOutlineStyle()
    style.stroke = 'none'
    style.hover.stroke = 'none'
    return style
  }

  /**
   * 重写此方法，使保存数据是能带上锚点数据。
   */
  getData() {
    const data = super.getData()
    data.sourceAnchorId = this.sourceAnchorId
    data.targetAnchorId = this.targetAnchorId
    return data
  }

  /**
   * 给边自定义方案，使其支持基于锚点的位置更新边的路径
   */
  updatePathByAnchor() {
    // TODO
    const sourceNodeModel = this.graphModel.getNodeModelById(this.sourceNodeId)
    const sourceAnchor = sourceNodeModel
      ?.getDefaultAnchor()
      .find((anchor) => anchor.id === this.sourceAnchorId)

    const targetNodeModel = this.graphModel.getNodeModelById(this.targetNodeId)
    const targetAnchor = targetNodeModel
      ?.getDefaultAnchor()
      .find((anchor) => anchor.id === this.targetAnchorId)

    if (sourceAnchor) {
      const startPoint = {
        x: sourceAnchor?.x,
        y: sourceAnchor?.y,
      }
      this.updateStartPoint(startPoint)
    }

    if (targetAnchor) {
      const endPoint = {
        x: targetAnchor?.x,
        y: targetAnchor?.y,
      }

      this.updateEndPoint(endPoint)
    }
    // 这里需要将原有的pointsList设置为空，才能触发bezier的自动计算control点。
    this.pointsList = []
    this.initPoints()
  }
}

export default {
  type: 'custom-bezier-edge',
  view: CustomEdge,
  model: CustomEdgeModel,
}