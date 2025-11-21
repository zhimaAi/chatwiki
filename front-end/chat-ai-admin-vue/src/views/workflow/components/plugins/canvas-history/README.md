# CanvasHistory 插件

一个轻量级的画布历史记录管理插件，专注于撤销/重做功能的实现，不包含UI界面组件。

## 功能特性

### 核心功能
- **历史记录管理**: 智能记录画布结构的变更
- **撤销/重做**: 支持最多50级历史记录的撤销和重做操作
- **快捷键支持**: 
  - `Ctrl+Z`: 撤销画布变更
  - `Ctrl+Y`: 重做撤销操作
  - `Ctrl+Shift+Z`: 重做撤销操作 (Mac兼容)
- **自动记录**: 自动监听LogicFlow事件并记录结构变更
- **双栈设计**: 使用撤销栈和重做栈管理操作历史
- **状态快照**: 保存完整的画布数据结构
- **事件驱动**: 基于LogicFlow事件监听架构
- **纯功能**: 不包含UI组件，专注于核心功能实现
- **事务支持**: 支持批量操作的事务处理

## 🚀 使用方法

### 快速开始

```javascript
import CanvasHistory from './plugins/canvas-history/index.js'

// 创建LogicFlow实例
const lf = new LogicFlow({
  // LogicFlow配置
})

// 加载历史记录插件
lf.install(CanvasHistory, {
  maxHistorySize: 50 // 最大历史记录数量，默认50
})

// 监听历史记录变化
lf.extension.canvasHistory.onHistoryChange((historyInfo) => {
  console.log('历史记录状态:', historyInfo)
  // historyInfo = {
  //   canUndo: boolean,           // 是否可以撤销
  //   canRedo: boolean,           // 是否可以重做
  //   historySize: number,        // 历史记录数量
  //   redoSize: number,           // 重做记录数量
  //   maxSize: number,            // 最大记录数
  //   currentCanvas: {            // 当前画布状态
  //     nodes: array,             // 当前节点列表
  //     edges: array,             // 当前连线列表
  //     nodeCount: number,        // 当前节点数量
  //     edgeCount: number         // 当前连线数量
  //   },
  //   historyDetails: array,      // 历史记录详细信息
  //   // 每个历史记录项 = {
  //   //   index: number,           // 索引位置
  //   //   timestamp: number,       // 时间戳
  //   //   operation: string,       // 操作类型
  //   //   nodeCount: number,       // 节点数量
  //   //   edgeCount: number,       // 连线数量
  //   //   nodes: array,            // 节点列表
  //   //   edges: array             // 连线列表
  //   // }
  //   redoDetails: array,         // 重做记录详细信息
  //   currentHistoryIndex: number // 当前历史索引（0为初始状态）
  // }
})
```

### 手动控制

```javascript
// 手动撤销
lf.extension.canvasHistory.undo()

// 手动重做
lf.extension.canvasHistory.redo()

// 获取历史记录状态
const canUndo = lf.extension.canvasHistory.canUndo()
const canRedo = lf.extension.canvasHistory.canRedo()
const historySize = lf.extension.canvasHistory.getHistorySize()

// 设置初始状态（加载画布数据后调用）
lf.extension.canvasHistory.setInitialState(graphData)

// 清空历史记录
lf.extension.canvasHistory.clearHistory()
```

### 事务处理

对于需要批量执行且只希望在最后保存一次历史记录的操作，可以使用事务功能：

```javascript
// 开始事务
lf.extension.canvasHistory.beginTransaction()

// 执行多个操作
lf.addNode({id: 'node1', type: 'rect', x: 100, y: 100})
lf.addNode({id: 'node2', type: 'circle', x: 200, y: 200})
lf.addEdge({sourceNodeId: 'node1', targetNodeId: 'node2'})

// 提交事务并保存历史记录
lf.extension.canvasHistory.commitTransaction('batch-add')
```

## 📋 记录规则

### ✅ 会被记录的操作
- 添加节点 (`node:add`)
- 节点拖拽 (`node:drop`)
- 添加连线 (`edge:add`)
- 连线调整 (`edge:adjust`)

### ❌ 不会被记录的操作
- 节点属性编辑（标签、样式等）
- 画布缩放和平移
- 节点删除和连线删除（需要自定义实现）

> 注意：为了提供更灵活的删除控制，`node:delete` 和 `edge:delete` 事件默认不会自动记录到历史记录中。如果您需要记录删除操作，请在删除后手动调用 `saveCurrentState('node:delete')` 或 `saveCurrentState('edge:delete')`。

## 🔧 配置选项

### 修改历史记录数量

```javascript
// 在插件初始化后修改
const canvasHistory = lf.extension.canvasHistory
canvasHistory.maxHistorySize = 100 // 改为100条记录
```

### 回调事件

```javascript
// 监听历史记录状态变化
lf.extension.canvasHistory.onHistoryChange((state) => {
  console.log('可撤销:', state.canUndo)
  console.log('可重做:', state.canRedo)
  console.log('历史记录数:', state.historySize)
  console.log('重做记录数:', state.redoSize)
  console.log('当前画布状态:', state.currentCanvas)
  console.log('历史记录详情:', state.historyDetails)
  console.log('重做记录详情:', state.redoDetails)
  console.log('当前历史索引:', state.currentHistoryIndex)
})
```

## 🐛 注意事项

### 1. 性能优化
- 默认最多保存50条历史记录，避免内存溢出
- 使用深拷贝避免状态引用问题
- 防递归调用保护机制
- 轻量化设计，无多余UI开销

### 2. 兼容性
- 支持现代浏览器（Chrome、Firefox、Safari、Edge）
- 支持 Mac 系统的 Cmd 键
- 自动检测输入框，避免快捷键冲突
- 仅监听撤销/重做快捷键

### 3. 删除操作处理
由于删除操作需要更精确的控制，插件默认不监听删除事件。如需记录删除操作，请在执行删除后手动保存状态：

```javascript
// 删除节点或连线后手动保存状态
lf.deleteElement(elementId)
lf.extension.canvasHistory.saveCurrentState('element:delete')
```

### 4. 最佳实践
- 在复杂的流程编辑场景中使用
- 结合其他 LogicFlow 插件一起使用
- 监听历史记录变化来实现自定义UI反馈（撤销/重做按钮状态）
- 对于批量操作使用事务功能

### 5. 架构特点
- **职责分离**: 专注于历史记录管理，UI集成由外部处理
- **灵活调用**: 所有功能通过方法调用，不强制UI绑定
- **事件驱动**: 基于LogicFlow事件自动记录结构变更

## 🔍 调试模式

可以通过浏览器开发者工具查看插件状态：

```javascript
// 在浏览器控制台中
const canvasHistory = lf.extension.canvasHistory
console.log('历史栈长度:', canvasHistory.historyStack.length)
console.log('重做栈长度:', canvasHistory.redoStack.length)
console.log('当前状态:', canvasHistory.createSnapshot())
console.log('可撤销:', canvasHistory.canUndo())
console.log('可重做:', canvasHistory.canRedo())
```

## 📝 API参考

| 方法名 | 参数 | 返回值 | 描述 |
|-------|------|--------|------|
| undo() | 无 | 无 | 撤销上一步操作 |
| redo() | 无 | 无 | 重做上一步操作 |
| canUndo() | 无 | Boolean | 是否可以撤销 |
| canRedo() | 无 | Boolean | 是否可以重做 |
| getHistorySize() | 无 | Number | 获取历史记录数量 |
| getRedoSize() | 无 | Number | 获取重做记录数量 |
| clearHistory() | 无 | 无 | 清空历史记录 |
| setInitialState(graphData) | Object | 无 | 设置初始状态 |
| saveCurrentState(operation) | String | 无 | 手动保存当前状态 |
| beginTransaction() | 无 | 无 | 开始事务 |
| commitTransaction(operation) | String | 无 | 提交事务 |
| onHistoryChange(callback) | Function | 无 | 注册历史变更回调 |

## 📝 更新日志

### v2.2.0 - 回调增强
- ✅ **增强功能**: 扩展 `notifyHistoryChange` 回调信息
- ✅ **新增信息**: 当前画布状态、历史记录详情、重做记录详情、当前历史索引
- ✅ **提供**: 更丰富的历史记录管理数据，便于自定义UI实现

### v2.1.0 - 功能增强
- ✅ **新增功能**: 添加事务处理支持
- ✅ **改进**: 默认历史记录数量增加到50条
- ✅ **优化**: 删除事件监听改为可选，提供更灵活的控制

### v2.0.0 - 轻量化重构
- ✅ **重大更新**: 移除所有UI组件，专注功能实现
- ✅ 纯功能插件架构
- ✅ 简化的键盘事件处理（仅撤销/重做）
- ✅ 更灵活的删除功能集成
- ✅ 更小的包体积，更好的性能

### v1.0.0
- ✅ 初始版本发布
- ✅ 支持撤销/重做/删除功能
- ✅ 5级历史记录管理
- ✅ UI工具栏
- ✅ 键盘快捷键支持
- ✅ 回调事件系统

---

🎉 **现在您拥有了一个轻量、功能强大的画布历史记录管理工具！**