export class CanvasHistory {
  static pluginName = 'canvasHistory'

  constructor({ lf }) {
    this.lf = lf
    this.historyStack = [] // 撤销栈
    this.redoStack = [] // 重做栈
    this.isProcessing = false // 防止递归调用
    this.historyChangeCallbacks = [] // 历史记录变化回调
    this._isSelectionDrop = false; // 
    this._isSelectionCompleted = false;

  }

  /**
   * 创建状态快照
   */
  createSnapshot() {
    const graphData = this.lf.getGraphData()
    const selectElements = this.lf.getSelectElements(true)
    let selectedElementIds = [...selectElements.nodes, ...selectElements.edges].map(item => item.id)
    return {
      nodes: JSON.parse(JSON.stringify(graphData.nodes || [])),
      edges: JSON.parse(JSON.stringify(graphData.edges || [])),
      selectedElementIds: selectedElementIds || [],
      timestamp: Date.now(),
      operation: 'snapshot'
    }
  }

  /**
   * 保存当前状态到历史栈
   */
  saveCurrentState(operation = 'change') {
    if (this.isProcessing) return

    const snapshot = this.createSnapshot()
    snapshot.operation = operation

    // 添加到历史栈
    this.historyStack.push(snapshot)

    // 限制历史记录数量
    if (this.historyStack.length > this.maxHistorySize) {
      this.historyStack.shift()
    }

    // 新操作时清空重做栈
    this.redoStack = []

    // 触发回调
    this.notifyHistoryChange()
  }

  /**
   * 撤销操作
   */
  undo() {
    if (this.historyStack.length <= 1 || this.isProcessing) return

    // 当前状态移到重做栈
    const currentState = this.historyStack.pop()
    this.redoStack.push(currentState)

    // 恢复到上一个状态
    const previousState = this.historyStack[this.historyStack.length - 1]
    this.restoreState(previousState)

    this.notifyHistoryChange()
  }

  /**
   * 重做操作
   */
  redo() {
    if (this.redoStack.length === 0 || this.isProcessing) return

    // 从重做栈取出状态
    const stateToRestore = this.redoStack.pop()
    this.historyStack.push(stateToRestore)

    // 恢复到该状态
    this.restoreState(stateToRestore)

    this.notifyHistoryChange()
  }

  /**
   * 恢复到指定状态
   */
  restoreState(state) {
    this.isProcessing = true

    try {
      // 清空当前画布
      this.lf.clearSelectElements()
      this.lf.clearData()

      state.nodes.forEach(item => {
        item.width = item.properties.width
        item.height = item.properties.height
      })

      this.lf.graphModel.graphDataToModel({
        nodes: state.nodes || [],
        edges: state.edges || [],
      })
      // 恢复节点
      // if (state.nodes && state.nodes.length > 0) {
      //   state.nodes.forEach(node => {
      //     this.lf.addNode(node)
      //   })
      // }

      // 恢复连线
      // if (state.edges && state.edges.length > 0) {
      //   state.edges.forEach(edge => {
      //     this.lf.addEdge(edge)
      //   })
      // }

      // 恢复选中状态
      if (state.selectedElementIds && state.selectedElementIds.length > 0) {
        state.selectedElementIds.forEach(id => {
          this.lf.selectElementById(id, true)
        })
      } 
    } finally {
      this.isProcessing = false
    }
  }

  /**
   * 历史记录键盘事件处理 - 仅处理撤销/重做快捷键
   */
  handleHistoryKeyDown = (event) => {
    // 检查是否在输入框中，避免冲突
    if (this.isInputElement(event.target)) return

    const { ctrlKey, key, metaKey } = event
    const isCtrlOrCmd = ctrlKey || metaKey

    // Ctrl+Z 撤销
    if (isCtrlOrCmd && key.toLowerCase() === 'z' && !event.shiftKey) {
      event.preventDefault()
      event.stopPropagation()
      this.undo()
      return
    }

    // Ctrl+Y 或 Ctrl+Shift+Z 重做
    if (isCtrlOrCmd && (key.toLowerCase() === 'y' || (key.toLowerCase() === 'z' && event.shiftKey))) {
      event.preventDefault()
      event.stopPropagation()
      this.redo()
      return
    }
  }

  /**
   * 判断是否为输入框元素
   */
  isInputElement(element) {
    const inputTypes = ['INPUT', 'TEXTAREA', 'SELECT']
    
    if (inputTypes.includes(element.tagName)) return true
    
    // 检查contenteditable属性
    let current = element
    while (current && current !== document.body) {
      if (current.hasAttribute && current.hasAttribute('contenteditable')) {
        return true
      }
      current = current.parentElement
    }
    
    return false
  }

  /**
   * 注册LogicFlow事件监听
   */
  registerLogicFlowEvents() {
    // 只记录影响画布结构的事件
    const eventsToRecord = [
      'node:add',
      // 'node:delete', // 改为自定义控制
      'node:drop',
      'edge:add',
      // 'edge:delete', // 改为自定义控制
      'edge:adjust',
      // 'selection:drop'
    ]

    eventsToRecord.forEach(eventName => {
      this.lf.on(eventName, () => {
        this.saveCurrentState(eventName)
      })
    })

    // selection:selected 也会触发 selection:drop，所以需要判断是否有选中元素
    this.lf.on('selection:drop', () => {
      this._isSelectionDrop = true;
      this._isSelectionCompleted = false; 
      // 延时一下等待selection:selected后在执行onselectionDrop
      setTimeout(() => {
        this.onselectionDrop()
      }, 20)
    })

    // 配合selection:drop用来区分是操作了选区还是移动了选区
    this.lf.on('selection:selected', () => {
      if(this._isSelectionDrop){
        this._isSelectionDrop = false;
        this._isSelectionCompleted = true;
      }
    })

    // 画布单击
    this.lf.on('blank:click', () => {
      this.clearSelectedStatus()
    })

    // 元素点击
    this.lf.on('element:click', () => {
      this.clearSelectedStatus()
    })
  }

  clearSelectedStatus(){
    this._isSelectionDrop = false;
    this._isSelectionCompleted = false;
  }

  onselectionDrop(){
    if(!this._isSelectionCompleted){
      this.saveCurrentState('selection:drop')
    }

    this._isSelectionDrop = false;
  }

  /**
   * 初始化历史记录功能
   */
  // eslint-disable-next-line no-unused-vars
  render(lf, toolOverlay, options) {
    this.maxHistorySize = (options && options.maxHistorySize) || 50
    // 注册历史记录相关键盘事件（仅撤销/重做）
    document.addEventListener('keydown', this.handleHistoryKeyDown)

    // 注册LogicFlow事件
    this.registerLogicFlowEvents()
  }

  /**
   * 从外部设置初始状态，通常在加载已有画布数据后调用
   * @param {Object} graphData - LogicFlow的图数据
   */
  setInitialState(graphData) {
    this.historyStack = []
    this.redoStack = []

    const snapshot = {
      nodes: JSON.parse(JSON.stringify(graphData.nodes || [])),
      edges: JSON.parse(JSON.stringify(graphData.edges || [])),
      timestamp: Date.now(),
      operation: 'init'
    }

    this.historyStack.push(snapshot)
    this.notifyHistoryChange()
  }

  /**
   * 检查是否可以撤销
   */
  canUndo() {
    return this.historyStack.length > 1
  }

  /**
   * 检查是否可以重做
   */
  canRedo() {
    return this.redoStack.length > 0
  }

  /**
   * 获取历史记录数量
   */
  getHistorySize() {
    return Math.max(0, this.historyStack.length - 1)
  }

  /**
   * 获取重做记录数量
   */
  getRedoSize() {
    return this.redoStack.length
  }

  /**
   * 注册历史记录变化回调
   */
  onHistoryChange(callback) {
    if (typeof callback === 'function') {
      this.historyChangeCallbacks.push(callback)
    }
  }

  /**
   * 通知历史记录变化
   */
  notifyHistoryChange() {
    this.historyChangeCallbacks.forEach(callback => {
      try {
        // 构建历史记录详细信息
        const historyDetails = this.historyStack.map((snapshot, index) => ({
          index: index,
          timestamp: snapshot.timestamp,
          operation: snapshot.operation,
          nodeCount: snapshot.nodes ? snapshot.nodes.length : 0,
          edgeCount: snapshot.edges ? snapshot.edges.length : 0,
          nodes: snapshot.nodes || [],
          edges: snapshot.edges || [],
          selectedElementIds: snapshot.selectedElementIds || []
        }))

        // 构建重做记录详细信息
        const redoDetails = this.redoStack.map((snapshot, index) => ({
          index: index,
          timestamp: snapshot.timestamp,
          operation: snapshot.operation,
          nodeCount: snapshot.nodes ? snapshot.nodes.length : 0,
          edgeCount: snapshot.edges ? snapshot.edges.length : 0,
          nodes: snapshot.nodes || [],
          edges: snapshot.edges || [],
          selectedElementIds: snapshot.selectedElementIds || []
        }))
        // 获取当前画布状态
        const currentHistoryIndex = this.historyStack.length - 1

        callback({
          canUndo: this.canUndo(),
          canRedo: this.canRedo(),
          historySize: this.getHistorySize(),
          redoSize: this.getRedoSize(),
          maxSize: this.maxHistorySize,
          // 当前画布状态
          currentCanvas: historyDetails[currentHistoryIndex],
          // 历史记录详情
          historyDetails: historyDetails,
          // 重做记录详情
          redoDetails: redoDetails,
          // 当前历史索引（0为初始状态）
          currentHistoryIndex: currentHistoryIndex
        })
      } catch (error) {
        console.error('History change callback error:', error)
      }
    })
  }

  /**
   * 清空历史记录
   */
  clearHistory() {
    this.historyStack = [this.createSnapshot()]
    this.redoStack = []
    this.notifyHistoryChange()
  }

  /**
   * 替换上一个历史状态
   */
  replaceLastState(operation = 'replace') {
    if (this.isProcessing) return

    if (this.canUndo()) {
      this.historyStack.pop()
    }

    // 不直接调用 saveCurrentState 是为了避免清空 redoStack
    const snapshot = this.createSnapshot()
    snapshot.operation = operation
    this.historyStack.push(snapshot)

    // 限制历史记录数量
    if (this.historyStack.length > this.maxHistorySize) {
      this.historyStack.shift()
    }
    
    this.notifyHistoryChange()
  }

  /**
   * 开始一个事务，暂停历史记录的自动保存
   */
  beginTransaction() {
    this.isProcessing = true;
  }

  /**
   * 提交一个事务，恢复历史记录的自动保存，并手动保存一次快照
   * @param {string} operation - 本次操作的类型描述
   */
  commitTransaction(operation = 'commit') {
    this.isProcessing = false;
    this.saveCurrentState(operation);
  }

  /**
   * 销毁插件
   */
  destroy() {
    // 移除历史记录键盘事件监听
    document.removeEventListener('keydown', this.handleHistoryKeyDown)
    
    // 清空回调
    this.historyChangeCallbacks = []
    
    // 清空状态
    this.historyStack = []
    this.redoStack = []
  }
}