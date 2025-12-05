# LogicFlow 自定义分组节点 (CustomGroupNode)

## 概述

`CustomGroupNode` 是一个为 LogicFlow 设计的高度定制化的分组节点。它不仅仅是一个简单的容器，而是通过 `foreignObject` 嵌入了一个完整的 Vue 组件 (`CustomGroupComponent`) 来渲染其内部 UI，从而实现了复杂的交互和视图。

该节点的核心功能是作为一个“循环”或“作用域”容器，在工作流中管理一组子节点，并处理其内部的变量、连接逻辑和交互行为。

## 主要功能

-   **Vue 组件集成**: 通过 `foreignObject` 将 Vue 组件 (`./index.vue`) 渲染为节点的视图，实现了丰富的 UI 和交互逻辑。
-   **动态内容渲染**: Vue 组件可以根据节点的 `properties`、`isSelected`、`isHovered` 等状态动态更新其视图。
-   **自定义锚点**: 实现了高度定制化的锚点，包括点击锚点弹出菜单、根据连接状态改变外观（加号/箭头）等。
-   **严格的连接规则**: 通过 `sourceRules` 和 `targetRules` 定义了严格的边连接逻辑，例如：
    -   只允许从右侧锚点连出，连接到左侧锚点。
    -   分组内的节点只能与同组内的其他节点连接。
-   **父级变量收集**: 具备 `getAllParentVariable` 方法，能够向上遍历图，收集所有前置节点的输出变量，为分组内部提供上下文。
-   **事件驱动**: 通过 `provide/inject` 和 `eventCenter` 与 Vue 子组件深度交互，例如子组件可以触发父画布添加新节点、更新自身尺寸和数据。
-   **自定义样式**: 通过覆盖 `getNodeStyle`、`getOutlineStyle` 等方法，实现了独特的节点外观，包括透明背景、无边框等（因为主要样式由 Vue 组件控制）。

## 文件结构

-   `index.js`: 定义了 `CustomGroupNode` 的 `view` 和 `model`，是该插件的核心逻辑所在。
-   `index.vue`: 定义了节点的 UI 和交互，作为 Vue 组件被 `index.js` 渲染。

## 实现详解

### `CustomGroupModel` (数据模型)

继承自 `dynamicGroup.model`，主要负责：

1.  **初始化**: 在 `initNodeData` 中设置节点的默认属性，如尺寸、`radius`，并关闭了 LogicFlow 的默认折叠、缩放等行为。
2.  **样式定义**: 在 `getNodeStyle` 和 `getOutlineStyle` 中将节点的默认背景和边框设置为透明或无，因为 UI 完全由 Vue 组件接管。
3.  **连接规则**: 在 `setAttributes` 中定义了 `sourceRules` 和 `targetRules`，用于在创建边时进行校验。
4.  **锚点定义**: `getDefaultAnchor` 方法定义了左右两个自定义锚点。
5.  **变量收集**: `getAllParentVariable` 是一个核心业务方法，它会回溯查找所有父节点，并聚合它们的输出变量，供分组节点内部使用。

### `CustomGroup` (视图)

继承自 `dynamicGroup.view`，主要负责：

1.  **渲染机制**:
    -   重写 `getShape` 方法，在 LogicFlow 的原生 SVG `rect` 基础上，追加一个 `<foreignObject>`。
    -   在 `<foreignObject>` 内部创建一个 `<div>` 作为 Vue 组件的挂载点。
2.  **Vue 实例管理**:
    -   `mountVueComponentOnNextTick` 方法负责创建和挂载 Vue 应用实例。
    -   首次渲染时，使用 `createApp` 创建 Vue 实例并 `mount` 到 `<div>` 上。
    -   后续更新时，只更新 Vue 组件的 `props` (`isSelected`, `isHovered`, `properties`)，避免了重复创建实例的开销。
3.  **数据通信**:
    -   通过 `provide` 向 Vue 子组件注入了多个方法，如 `getNode`、`getGraph`、`addNode`、`setData` 等，使子组件能够与 LogicFlow 画布进行通信。
4.  **自定义锚点渲染**:
    -   重写 `getAnchorShape` 方法，根据锚点是否已被连接 (`sourceAnchorIds`, `targetAnchorIds`)，动态渲染出不同的 SVG 图形（带箭头的圆或带加号的圆）。
5.  **生命周期**:
    -   在 `destroy` 方法中，调用 `app.unmount()` 来确保 Vue 实例被正确销毁，防止内存泄漏。

## 如何使用

该自定义节点通过 `register` 方法在 LogicFlow 实例中注册。

```javascript
import LogicFlow from '@logicflow/core';
import { register } from '@logicflow/vue-node-registry';
import customGroupNode from './nodes/custom-group-node/index.js';

// ... 初始化 lf 实例
const lf = new LogicFlow({
  // ...
});

// 注册自定义分组节点
register(customGroupNode, lf);

// 添加节点
lf.addNode({
  type: 'custom-group',
  x: 300,
  y: 300,
  properties: {
    // ... 自定义属性
  }
});
```
