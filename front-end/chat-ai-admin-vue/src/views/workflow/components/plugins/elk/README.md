# Elk.js LogicFlow 插件文档

## 概述

该文件是一个为 `LogicFlow` 设计的插件，集成了 `elkjs` 库，以实现对工作流图的自动分层布局。此插件能够处理包含普通节点、分组节点和嵌套节点的复杂图结构，并自动计算节点位置和边的路径，以生成清晰、美观的布局。

## 主要功能

-   **自动布局**: 调用 `layout` 方法，对传入的 `graphData` 进行全自动布局。
-   **层级结构支持**: 通过 `arrayToTree` 工具函数，将扁平的节点列表转换为 `elkjs` 所需的树形结构，完美支持嵌套节点和分组。
-   **节点类型适配**: 能够识别不同类型的节点（如 `custom-group`, `group-start-node`, `group-end-node` 等），并为它们应用不同的布局参数和尺寸约束。
-   **边路径计算**: 自动计算边的路径，并生成平滑的贝塞尔曲线。
-   **配置灵活**: 提供了丰富的 `layoutOptions` 配置项，可以精细调整布局算法、方向、间距等参数。

## 文件结构

```javascript
// elk.js

// 1. 辅助函数：arrayToTree
function arrayToTree(flatData, idKey = 'id', parentIdKey = 'parentId', childrenKey = 'children') {
  // ... 将扁平化数组转换为树形结构
}

// 2. 主类：Elk
class Elk {
  static pluginName = 'elk';

  constructor({ lf }) {
    // ... 初始化
  }

  async layout(graphData) {
    // ... 核心布局逻辑
  }
}

// 3. 导出
export { Elk };
```

## `Elk` 类详解

### `constructor({ lf })`

-   **`lf`**: LogicFlow 的实例。插件通过该实例与 LogicFlow 交互，例如在布局完成后更新画布。

### `async layout(graphData)`

这是插件的核心方法，负责执行整个布局流程。

-   **参数**: `graphData` (`object`) - LogicFlow 导出的图数据，包含 `nodes` 和 `edges` 数组。
-   **返回值**: `Promise<void>` - 该方法是异步的，在布局计算和画布更新完成后结束。

#### 布局流程：

1.  **数据转换**:
    -   将输入的 `graphData` 转换为 `elkjs` 需要的格式。
    -   遍历 `nodes`，根据节点类型（`custom-group`, `group-start-node` 等）设置不同的尺寸（`width`, `height`）和布局选项（`layoutOptions`）。
    -   利用 `arrayToTree` 函数将带有 `parentId` 的节点列表转换为层级结构。
    -   转换 `edges` 数据，指定边的源节点（`sources`）和目标节点（`targets`）。

2.  **创建 ELK 图**:
    -   构建一个 `elkGraph` 对象，包含 `id`, `layoutOptions`, `children` (转换后的节点树) 和 `edges`。
    -   `layoutOptions` 在这里被集中配置，用于控制全局布局行为，例如：
        -   `'elk.algorithm': 'layered'` (使用分层算法)
        -   `'elk.direction': 'RIGHT'` (布局从左到右)
        -   `'elk.spacing.nodeNode': 200` (节点间距)

3.  **执行布局**:
    -   调用 `this.elk.layout(elkGraph)`，`elkjs` 会异步计算出所有节点和边的最佳位置和路径。

4.  **更新节点位置**:
    -   递归遍历 `elkjs` 返回的 `newElkGraph`。
    -   根据计算出的 `x`, `y` 坐标更新 `originalGraphData` 中对应节点的坐标。
    -   对于嵌套节点，其坐标是相对于父节点计算的，该方法会正确处理这种相对关系，计算出最终的绝对坐标。

5.  **更新边路径**:
    -   遍历 `newElkGraph.edges`，获取计算出的边路径信息 (`sections`)。
    -   为 `originalGraphData` 中的每条边生成 `pointsList`，用于绘制平滑的贝塞尔曲线。
    -   保留原始的锚点信息（`sourceAnchorId`, `targetAnchorId`）。

6.  **更新画布**:
    -   调用 `this.lf.graphModel.graphDataToModel(originalGraphData)`，将更新后的图数据应用到 LogicFlow 画布上，完成视图的重新渲染。

## `arrayToTree` 函数

这是一个独立的工具函数，用于将包含 `id` 和 `parentId` 的扁平节点数组转换为树形结构。

-   **参数**:
    -   `flatData` (`Array`): 扁平的节点数据。
    -   `idKey` (`string`): 节点唯一标识的键名，默认为 `'id'`。
    -   `parentIdKey` (`string`): 父节点标识的键名，默认为 `'parentId'`。
    -   `childrenKey` (`string`): 子节点数组的键名，默认为 `'children'`。
-   **返回值**: `Array` - 转换后的树形结构数组。

## 如何使用

1.  **注册插件**:
    ```javascript
    import LogicFlow from '@logicflow/core';
    import { Elk } from '@/views/workflow/components/plugins/elk/elk.js';

    const lf = new LogicFlow({
      // ... 其他配置
      plugins: [Elk]
    });
    ```

2.  **调用布局**:
    在需要执行自动布局的地方（例如，点击一个“自动布局”按钮），获取当前图数据并调用 `layout` 方法。
    ```javascript
    // 获取 LogicFlow 实例上的插件
    const elk = lf.extension.elk;

    // 获取图数据并执行布局
    const graphData = lf.getGraphData();
    await elk.layout(graphData);
