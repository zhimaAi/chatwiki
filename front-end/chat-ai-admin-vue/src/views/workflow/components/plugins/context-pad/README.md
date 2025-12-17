# ContextPad 插件

`ContextPad` 是一个为 `LogicFlow` 工作流引擎设计的插件，它提供了一个上下文菜单（或称快捷面板），允许用户在画布的特定位置快速添加新的节点。

## 文件结构

- `index.js`: 插件的主要逻辑，负责处理菜单的创建、显示、隐藏和交互。
- `index.vue`: 菜单的 `Vue` UI 组件，展示了可供选择添加的节点列表。

## 主要功能

`ContextPad` 插件通过监听 `LogicFlow` 的事件，动态地在画布上创建和管理一个上下文菜单。

### 1. 插件初始化 (`constructor`)

-   **DOM 创建**: 创建一个 `div` 元素作为菜单的容器 (`__menuDOM`)。
-   **Vue 实例**: 使用 `Vue 3` 的 `createApp` 和 `h` API 创建一个 `Vue` 应用实例。这个实例负责渲染 `index.vue` 组件，该组件构成了菜单的界面。
-   **数据提供**: 通过 `provide` 向 `Vue` 组件注入了获取当前节点 (`_activeData`) 和 `GraphModel` 的方法。
-   **交互回调**: 定义了 `onClickItem` 回调函数。当用户在菜单中点击一个选项时，该函数会被触发，进而通过 `lf.graphModel.eventCenter.emit` 触发一个 `custom:addNode` 事件，以在工作流中添加新节点。

### 2. 插件渲染与事件监听 (`render`)

-   **挂载**: 在插件首次渲染时，将菜单的 `DOM` 元素 (`__menuDOM`) 添加到 `LogicFlow` 的容器中，并挂载 `Vue` 应用。
-   **事件监听**:
    -   `custom:showPopupMenu`: 当此自定义事件被触发时，插件会记录下锚点数据 (`anchorData`) 和节点模型 (`model`)，并调用 `createMenu()` 来显示上下文菜单。
    -   `node:click`, `edge:click`, `blank:click`: 当用户点击节点、边或画布空白处时，插件会调用 `hideMenu()` 来隐藏菜单。
    -   `node:delete`, `edge:delete`, `node:drag`, `graph:transform`: 在菜单显示期间，会监听这些事件，以便在节点被删除、拖拽或画布变换时自动隐藏菜单。

### 3. 菜单的定位与显示

-   **`getContextMenuPosition()`**: 此方法根据 `_activeAnchorData`（通常是节点的某个锚点）的坐标，计算出菜单应该显示的初始 `HTML` 坐标。
-   **`showMenu()`**:
    1.  调用 `getContextMenuPosition()` 获取菜单的理想位置。
    2.  进行复杂的边界检测：
        -   确保菜单不会超出画布的底部。如果超出，会自动向上调整位置。
        -   确保菜单不会超出画布的顶部。如果超出，会将其固定在顶部附近。
    3.  设置菜单的 `left`, `top`, `height` 等 `CSS` 属性，使其在计算出的位置正确显示。
    4.  将菜单的 `display` 样式设置为 `flex`，使其可见。

### 4. 菜单的隐藏 (`hideMenu`)

-   将菜单的 `display` 样式设置为 `none`。
-   移除在 `showMenu` 中添加的事件监听器，以避免不必要的计算和内存泄漏。

## 使用流程

1.  外部模块（如某个节点）触发 `custom:showPopupMenu` 事件，并传递锚点数据和节点模型。
2.  `ContextPad` 插件监听到该事件，显示上下文菜单 (`index.vue`)。
3.  菜单被定位在锚点的右侧，并根据画布边界进行智能调整。
4.  用户在菜单中选择一个要添加的节点。
5.  `onClickItem` 回调被触发，进而触发 `custom:addNode` 事件，并附带了新节点类型、原始节点模型和锚点数据。
6.  其他监听 `custom:addNode` 事件的模块（如 `workflow-canvas.vue`）会处理节点的实际添加逻辑。
7.  菜单被自动隐藏。
