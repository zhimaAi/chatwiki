import { h } from '@logicflow/core'

const EXTRA_PAN_PADDING = 300 // 屏幕 CSS 像素，不随画布缩放变化
const EDGE_THRESHOLD = 40
const PAN_SPEED = 240 // 屏幕 CSS 像素 / 秒
const MAX_FRAME_DELTA = 50
const MIN_PAN_DISTANCE = 0.01

function getAnchorKey(nodeId, anchorId) {
  return JSON.stringify([nodeId, anchorId])
}

function getNodeBounds(nodes) {
  let bounds = null

  for (const node of nodes) {
    if (node.visible === false || node.virtual) continue

    const { x, y, width, height } = node
    if (![x, y, width, height].every(Number.isFinite)) continue

    const left = x - width / 2
    const right = x + width / 2
    const top = y - height / 2
    const bottom = y + height / 2

    if (!bounds) {
      bounds = { left, right, top, bottom }
    } else {
      bounds.left = Math.min(bounds.left, left)
      bounds.right = Math.max(bounds.right, right)
      bounds.top = Math.min(bounds.top, top)
      bounds.bottom = Math.max(bounds.bottom, bottom)
    }
  }

  return bounds
}

// 只扩展当前画布的节点 view，通过稳定 ref 取得原生锚点实例，保留节点外观与连接规则。
export function withEdgeAutoPan(NodeView, controller) {
  if (!controller || !NodeView) return NodeView

  return class extends NodeView {
    getAnchors() {
      const anchors = super.getAnchors()
      if (!this.autoPanAnchorRefs) this.autoPanAnchorRefs = new Map()

      return anchors.map((anchor) => {
        const { nodeModel, anchorData } = anchor.props
        const key = getAnchorKey(nodeModel.id, anchorData.id)

        if (!this.autoPanAnchorRefs.has(key)) {
          this.autoPanAnchorRefs.set(key, (instance) => controller.registerAnchor(key, instance))
        }

        return h(anchor.type, {
          ...anchor.props,
          key: anchor.key ?? key,
          ref: this.autoPanAnchorRefs.get(key)
        })
      })
    }
  }
}

export class EdgeAutoPan {
  static pluginName = 'edgeAutoPan'

  constructor({ lf }) {
    this.lf = lf
    this.anchors = new Map()
    this.activeAnchorKey = null
    this.sourceNodeId = null
    this.pointerEvent = null
    this.pointerReleased = false
    this.frameId = null
    this.lastFrameTime = null
    this.destroyed = false

    this.handleDragStart = this.handleDragStart.bind(this)
    this.handleDrag = this.handleDrag.bind(this)
    this.endDrag = this.endDrag.bind(this)
    this.stopPan = this.stopPan.bind(this)
    this.handlePointerUp = this.handlePointerUp.bind(this)
    this.handleVisibilityChange = this.handleVisibilityChange.bind(this)
    this.handleNodeDelete = this.handleNodeDelete.bind(this)
    this.tick = this.tick.bind(this)

    this.flowListeners = [
      ['anchor:dragstart', this.handleDragStart],
      ['anchor:drag', this.handleDrag],
      ['anchor:dragend', this.endDrag],
      ['node:delete', this.handleNodeDelete]
    ]
    this.flowListeners.forEach(([event, listener]) => lf.on(event, listener))
    document.addEventListener('pointerup', this.handlePointerUp, true)
    document.addEventListener('pointercancel', this.endDrag, true)
    document.addEventListener('visibilitychange', this.handleVisibilityChange)
    window.addEventListener('blur', this.stopPan)
  }

  registerAnchor(key, instance) {
    if (this.destroyed) return

    if (instance) {
      this.anchors.set(key, instance)
    } else {
      this.anchors.delete(key)
      if (key === this.activeAnchorKey) this.endDrag()
    }
  }

  getActiveAnchor() {
    const anchor = this.anchors.get(this.activeAnchorKey)
    // 临时连线状态属于 LogicFlow 内部视图能力，升级依赖时只需复核这一处适配。
    if (!anchor?.state || typeof anchor.setState !== 'function' ||
        typeof anchor.moveAnchorEnd !== 'function') return null
    if (!this.lf.getNodeModelById(this.sourceNodeId)) return null
    return anchor
  }

  handleDragStart({ data, nodeModel, e }) {
    this.endDrag()
    if (this.destroyed) return

    this.activeAnchorKey = getAnchorKey(nodeModel.id, data.id)
    this.sourceNodeId = nodeModel.id
    this.pointerEvent = e
    this.pointerReleased = false
    if (!this.getActiveAnchor()) this.endDrag()
  }

  handleDrag({ data, nodeModel, e }) {
    if (this.destroyed || !this.activeAnchorKey || this.pointerReleased) return
    if (getAnchorKey(nodeModel.id, data.id) !== this.activeAnchorKey) return
    if (!e || !(e.buttons & 1)) {
      this.endDrag()
      return
    }

    this.pointerEvent = e
    if (document.hidden || !document.hasFocus() || !this.getActiveAnchor() || !this.getPanMovement(1)) {
      this.stopPan()
      return
    }

    if (this.frameId === null) {
      this.lastFrameTime = performance.now()
      this.frameId = requestAnimationFrame(this.tick)
    }
  }

  getPanMovement(distance) {
    const { graphModel } = this.lf
    const { width, height, transformModel, editConfigModel } = graphModel
    if (!this.pointerEvent || width <= 0 || height <= 0 || editConfigModel.stopMoveGraph) return null

    const { clientX, clientY } = this.pointerEvent
    const { domOverlayPosition } = graphModel.getPointByClient({ x: clientX, y: clientY })
    // 贴边或越出容器仍属于向该侧拖线，不能停止平移；仅触发判断使用限定后的坐标。
    const x = Math.max(0, Math.min(domOverlayPosition.x, width))
    const y = Math.max(0, Math.min(domOverlayPosition.y, height))

    // Vue 节点测量会直接更新模型尺寸，拖线期间读取当前外轮廓，避免缓存过期导致越界。
    const bounds = getNodeBounds(graphModel.nodes)
    if (!bounds) return null

    const [left, top] = transformModel.CanvasPointToHtmlPoint([bounds.left, bounds.top])
    const [right, bottom] = transformModel.CanvasPointToHtmlPoint([bounds.right, bounds.bottom])
    const thresholdX = Math.min(EDGE_THRESHOLD, width / 2)
    const thresholdY = Math.min(EDGE_THRESHOLD, height / 2)
    let directionX = 0
    let directionY = 0
    let availableX = 0
    let availableY = 0

    // 指针向右/下拖动时，视野向右/下推进，画布内容需要向左/上平移。
    if (x < thresholdX) {
      directionX = 1
      availableX = Math.max(0, EXTRA_PAN_PADDING - left)
    } else if (x > width - thresholdX) {
      directionX = -1
      availableX = Math.max(0, right + EXTRA_PAN_PADDING - width)
    }
    if (y < thresholdY) {
      directionY = 1
      availableY = Math.max(0, EXTRA_PAN_PADDING - top)
    } else if (y > height - thresholdY) {
      directionY = -1
      availableY = Math.max(0, bottom + EXTRA_PAN_PADDING - height)
    }

    if (availableX < MIN_PAN_DISTANCE) directionX = 0
    if (availableY < MIN_PAN_DISTANCE) directionY = 0
    const directionLength = Math.hypot(directionX, directionY)
    if (!directionLength) return null

    const axisDistance = distance / directionLength
    return {
      x: directionX * Math.min(axisDistance, availableX),
      y: directionY * Math.min(axisDistance, availableY)
    }
  }

  tick(timestamp) {
    this.frameId = null
    if (this.destroyed || this.pointerReleased || document.hidden || !document.hasFocus() || !this.getActiveAnchor()) {
      this.stopPan()
      return
    }

    const deltaMs = Math.min(Math.max(timestamp - this.lastFrameTime, 0), MAX_FRAME_DELTA)
    this.lastFrameTime = timestamp
    if (deltaMs === 0) {
      this.frameId = requestAnimationFrame(this.tick)
      return
    }

    const movement = this.getPanMovement(PAN_SPEED * deltaMs / 1000)
    if (!movement) {
      this.stopPan()
      return
    }

    const { transformModel } = this.lf.graphModel
    const previousX = transformModel.TRANSLATE_X
    const previousY = transformModel.TRANSLATE_Y
    transformModel.translate(movement.x, movement.y)
    if (previousX === transformModel.TRANSLATE_X && previousY === transformModel.TRANSLATE_Y) {
      this.stopPan()
      return
    }

    this.syncAnchorEnd()
    if (!this.destroyed && this.activeAnchorKey && !this.pointerReleased) {
      this.frameId = requestAnimationFrame(this.tick)
    }
  }

  syncAnchorEnd() {
    const anchor = this.getActiveAnchor()
    if (!anchor || !this.pointerEvent) return

    const { clientX, clientY } = this.pointerEvent
    // 平移量是屏幕像素，末端是画布坐标；重新映射鼠标位置才能兼容任意缩放。
    const { canvasOverlayPosition: { x, y } } = this.lf.graphModel.getPointByClient({ x: clientX, y: clientY })
    anchor.setState({ endX: x, endY: y })
    anchor.moveAnchorEnd(x, y, this.pointerEvent)
  }

  handlePointerUp(event) {
    if (!this.pointerEvent || event.pointerId !== this.pointerEvent.pointerId) return
    // 捕获阶段先停动画，最终创建连线仍交给原生结束流程，避免松手后多移动一帧。
    this.pointerReleased = true
    this.stopPan()
  }

  handleVisibilityChange() {
    if (document.hidden) this.stopPan()
  }

  handleNodeDelete({ data }) {
    if (data.id === this.sourceNodeId) this.endDrag()
  }

  stopPan() {
    if (this.frameId !== null) cancelAnimationFrame(this.frameId)
    this.frameId = null
    this.lastFrameTime = null
  }

  endDrag() {
    this.stopPan()
    this.activeAnchorKey = null
    this.sourceNodeId = null
    this.pointerEvent = null
    this.pointerReleased = false
  }

  destroy() {
    if (this.destroyed) return
    this.destroyed = true
    this.endDrag()
    this.anchors.clear()
    this.flowListeners.forEach(([event, listener]) => this.lf.off(event, listener))
    document.removeEventListener('pointerup', this.handlePointerUp, true)
    document.removeEventListener('pointercancel', this.endDrag, true)
    document.removeEventListener('visibilitychange', this.handleVisibilityChange)
    window.removeEventListener('blur', this.stopPan)
  }
}
