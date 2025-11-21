/**
 * 自定义键盘快捷键插件
 * @class CustomKeyboard
 * @property {string} pluginName - 插件名称
 * @property {object} lf - LogicFlow 实例
 * @property {boolean} isCtrlPressed - Ctrl 键是否被按下
 * @property {number} zoom - 当前缩放值
 * @property {object|null} clipboard - 剪贴板，用于存储复制的节点数据
 * @property {number} basePointX - 粘贴节点的 X 轴偏移量
 * @property {number} basePointY - 粘贴节点的 Y 轴偏移量
 */
class CustomKeyboard {
  static pluginName = 'custom-keyboard'
  constructor({ lf }) {
    this.lf = lf
    this.isCtrlPressed = false
    this.zoom = 1
    this.clipboard = null // 用于存储复制的节点数据
    this.basePointX = 20 // 粘贴时节点的 X 轴基础偏移量
    this.basePointY = 20 // 粘贴时节点的 Y 轴基础偏移量
  }

  /**
   * 插件渲染方法，用于绑定事件监听器
   */
  render() {
    this.lf.container.tabIndex = 0
    this.lf.container.addEventListener('wheel', this.wheelEventListener, { passive: false })
    document.addEventListener('keydown', this.keydownListener, { passive: false })
    document.addEventListener('keyup', this.keyupListener, { passive: false })
  }

  /**
   * 插件销毁方法，用于移除事件监听器
   */
  destroy() {
    this.lf.container.removeEventListener('wheel', this.wheelEventListener)
    document.removeEventListener('keydown', this.keydownListener)
    document.removeEventListener('keyup', this.keyupListener)
  }

  /**
   * 鼠标滚轮事件监听器
   * - Shift + 滚轮: 画布横向滚动
   * - Ctrl + 滚轮: 画布缩放
   * @param {WheelEvent} e - 滚轮事件对象
   */
  wheelEventListener = (e) => {
    // Shift + 滚轮: 画布横向滚动
    if (e.shiftKey) {
      e.preventDefault()
      const transform = this.lf.getTransform()
      let SCALE_X = transform.SCALE_X
      let SCALE_Y = transform.SCALE_Y
      if (e.deltaY > 0) {
        this.lf.translate(-100 * SCALE_X, 100 * SCALE_Y)
      } else {
        this.lf.translate(100 * SCALE_X, -100 * SCALE_Y)
      }
      // Ctrl + 滚轮: 画布缩放
    } else if (e.ctrlKey) {
      e.preventDefault()
      const { transformModel } = this.lf.graphModel
      const [canvasX, canvasY] = transformModel.HtmlPointToCanvasPoint([e.clientX, e.clientY])

      if (e.deltaY < 0) {
        // 放大
        this.zoom += 0.04
      } else {
        // 缩小
        this.zoom -= 0.04
      }

      this.lf.zoom(this.zoom, [canvasX, canvasY])
    }
  }

  /**
   * 键盘按下事件监听器
   * - Ctrl + C: 复制节点
   * - Ctrl + V: 粘贴节点
   * - Delete: 删除节点
   * - Ctrl (长按): 激活套索选择
   * @param {KeyboardEvent} e - 键盘事件对象
   */
  keydownListener = (e) => {
    const key = e.key.toLowerCase()
    const {
      graphModel: { eventCenter }
    } = this.lf

    // 判断事件是否发生在画布上
    const target = e.target
    const isCanvas = target === document.body || this.lf.container.contains(target)

    if (e.ctrlKey) {
      // Ctrl + C: 复制节点
      if (key === 'c') {
        if (!isCanvas) return
        e.preventDefault()
        let selectedElements = this.lf.getSelectElements(true)
        selectedElements = JSON.parse(JSON.stringify(selectedElements))
        // 去掉边的信息
        selectedElements.edges = []

        // 过滤掉开始节点和分组节点
        const validNodes = selectedElements.nodes.filter(
          (node) => node.type !== 'start-node' && node.type !== 'custom-group'
        )

        if (validNodes.length > 0) {
          // 重置粘贴偏移量
          this.basePointX = 20
          this.basePointY = 20
          // 计算选中节点的左上角基点
          let basePoint = { x: Infinity, y: Infinity }
          validNodes.forEach((node) => {
            // 清空下一个节点的键值
            node.id = '';
            node.nodeSortKey = '';
            node.properties.id = '';
            node.properties.next_node_key = '';
            node.properties.nodeSortKey = '';
            node.properties.node_key = '';
            
            if (!node.width && node.properties.width) {
              node.width = node.properties.width
            }
            if (!node.height && node.properties.height) {
              node.height = node.properties.height
            }
            if (node.x < basePoint.x) basePoint.x = node.x
            if (node.y < basePoint.y) basePoint.y = node.y
          })
          // 将节点数据和基点存入剪贴板
          this.clipboard = {
            nodes: validNodes,
            basePoint: basePoint
          }
        } else {
          this.clipboard = null
        }
        return
      }

      // Ctrl + V: 粘贴节点
      if (key === 'v') {
        if (!isCanvas) return
        e.preventDefault()
        if (this.clipboard && this.clipboard.nodes.length > 0) {
          const pasteBasePoint = {
            x: this.clipboard.basePoint.x + this.basePointX,
            y: this.clipboard.basePoint.y + this.basePointY
          }

          // 触发自定义粘贴事件
          eventCenter.emit('custom:paste', {
            originalNodes: this.clipboard.nodes,
            basePoint: this.clipboard.basePoint,
            pasteBasePoint: pasteBasePoint
          })

          // 增加下一次粘贴的偏移量
          this.basePointX += 20
          this.basePointY += 20
        }
        return
      }
    }

    // Delete: 删除节点
    if (key === 'delete' && isCanvas) {
      eventCenter.emit('custom:keyoard:delete')
      return
    }

    // Ctrl (长按): 激活套索选择
    if (e.ctrlKey && !this.isCtrlPressed) {
      e.preventDefault()
      this.isCtrlPressed = true
      this.lf.extension.selectionSelect.openSelectionSelect()
      this.lf.once('selection:selected', () => {
        this.isCtrlPressed = false
        this.lf.extension.selectionSelect.closeSelectionSelect()
      })
    }
  }

  /**
   * 键盘松开事件监听器
   * - Ctrl: 关闭套索选择
   * @param {KeyboardEvent} e - 键盘事件对象
   */
  keyupListener = (e) => {
    // 松开 Ctrl 键，关闭套索选择
    if (!e.ctrlKey && this.isCtrlPressed) {
      e.preventDefault()
      this.isCtrlPressed = false
      this.lf.extension.selectionSelect.closeSelectionSelect()
    }
  }
}

export { CustomKeyboard }