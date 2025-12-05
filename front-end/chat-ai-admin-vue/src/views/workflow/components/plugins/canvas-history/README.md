# CanvasHistory 插件

`CanvasHistory` 是一个为 LogicFlow 设计的增强型历史记录管理插件，提供了强大的撤销（Undo）和重做（Redo）功能。它能够捕获画布上的关键操作，生成状态快照，并允许用户在不同的历史状态之间轻松切换。

## 功能特性

-   **自动状态捕获**: 自动监听画布上的关键变更事件（如节点添加、拖拽、连线变更等），并保存状态快照。
-   **撤销与重做**:
    -   `Ctrl` + `Z` (或 `Cmd` + `Z`): 撤销上一步操作。
    -   `Ctrl` + `Y` (或 `Cmd` + `Shift` + `Z`): 重做已撤销的操作。
-   **智能事件处理**: 忽略在输入框或可编辑元素中的快捷键，避免与文本输入冲突。
-   **状态恢复**: 不仅恢复节点和边的位置与属性，还能恢复元素的原有选中状态。
-   **丰富的 API**: 提供了一系列 API 用于外部控制，如手动保存状态、清空历史、检查撤销/重做状态等。
-   **事务管理**: 支持将多个连续操作合并为一个历史记录单元，优化历史记录的粒度。
-   **回调机制**: 允许注册回调函数，在历史记录发生变化时获得通知，并获取详细的状态信息。

## 如何使用

1.  **引入插件**:

    ```javascript
    import { CanvasHistory } from './canvas-history';
    ```

2.  **注册插件**:

    在初始化 LogicFlow 实例时，将 `CanvasHistory` 插件注册到 `plugins` 数组中。

    ```javascript
    import LogicFlow from '@logicflow/core';
    import { CanvasHistory } from './canvas-history';

    const lf = new LogicFlow({
      // ... 其他 LogicFlow 配置
      keyboard: {
        enabled: false, // 开启自定义快捷键
        shortcuts: [
          // 屏蔽以下自带的快捷键
          { keys: ["cmd + z", "ctrl + z"], callback: () => {} },
          { keys: ["cmd + y", "ctrl + y"], callback: () => {} },
          { keys: ["cmd + shift + z", "ctrl + shift + z"], callback: () => {} }
        ]
      },
      plugins: [CanvasHistory]
    });
    ```

> **重要**: 为了让此插件的撤销/重做功能正常工作，必须在 LogicFlow 初始化时通过 `keyboard` 选项屏蔽掉内置的 `Ctrl+Z` 和 `Ctrl+Y` 相关快捷键。

## 快捷键

| 快捷键                       | 功能                               |
| ---------------------------- | ---------------------------------- |
| `Ctrl` + `Z` (Windows/Linux) | 撤销上一步操作                     |
| `Cmd` + `Z` (macOS)          | 撤销上一步操作                     |
| `Ctrl` + `Y` (Windows/Linux) | 重做上一步操作                     |
| `Ctrl`+`Shift`+`Z` (Win/Linux) | 重做上一步操作                     |
| `Cmd` + `Shift` + `Z` (macOS)  | 重做上一步操作                     |

## API 和方法

### 核心方法

-   **`undo()`**:
    执行一次撤销操作。

-   **`redo()`**:
    执行一次重做操作。

-   **`saveCurrentState(operation: string)`**:
    手动保存当前画布的状态快照。`operation` 参数用于描述本次操作的类型。

-   **`setInitialState(graphData: object)`**:
    设置初始历史状态，通常在加载已有画布数据后调用，以建立历史记录的基线。

-   **`clearHistory()`**:
    清空所有历史记录和重做记录，并将当前画布状态作为新的初始状态。

### 状态检查

-   **`canUndo(): boolean`**:
    检查当前是否可以执行撤销操作。

-   **`canRedo(): boolean`**:
    检查当前是否可以执行重做操作。

-   **`getHistorySize(): number`**:
    获取当前撤销栈中的历史记录数量。

-   **`getRedoSize(): number`**:
    获取当前重做栈中的记录数量。

### 事务管理

-   **`beginTransaction()`**:
    开始一个事务，暂停历史记录的自动保存。用于将一系列连续的、细粒度的操作合并。

-   **`commitTransaction(operation: string)`**:
    提交一个事务，恢复历史记录的自动保存，并手动保存一次包含所有变更的最终快照。

### 事件回调

-   **`onHistoryChange(callback: function)`**:
    注册一个回调函数，当历史记录（撤销/重做栈）发生变化时被调用。回调函数会接收到一个包含详细状态的对象。

    **回调参数 (`historyDetails`):**
    ```typescript
    {
      canUndo: boolean,
      canRedo: boolean,
      historySize: number, // 可撤销的步骤数
      redoSize: number,    // 可重做的步骤数
      maxSize: number,     // 历史记录最大容量
      currentCanvas: object, // 当前画布状态的快照
      historyDetails: object[], // 完整的历史记录栈信息
      redoDetails: object[],    // 完整的重做栈信息
      currentHistoryIndex: number // 当前状态在历史栈中的索引
    }
    ```

    **使用示例**:
    ```javascript
    const history = lf.extension.canvasHistory;
    history.onHistoryChange(({ canUndo, canRedo }) => {
      // 例如，根据 canUndo 和 canRedo 的状态更新 UI 按钮的可用性
      document.getElementById('undo-btn').disabled = !canUndo;
      document.getElementById('redo-btn').disabled = !canRedo;
    });
    ```

## 插件选项

可以在注册插件时传入配置选项。

-   **`maxHistorySize`**:
    -   **类型**: `number`
    -   **默认值**: `50`
    -   **描述**: 设置历史记录栈的最大容量，超过此容量的旧记录将被丢弃。

    **配置示例**:
    ```javascript
    const lf = new LogicFlow({
      // ...
      plugins: [
        [CanvasHistory, { maxHistorySize: 100 }]
      ]
    });