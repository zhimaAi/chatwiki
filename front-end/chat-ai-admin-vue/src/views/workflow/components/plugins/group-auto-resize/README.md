# LogicFlow 分组自动调整尺寸插件 (GroupAutoResize)

## 概述

`GroupAutoResize` 是一个为 `LogicFlow` 设计的交互增强插件。它的核心功能是监听分组（Group）内部子节点的拖拽事件，当子节点被拖动到超出其父分组的边界时，该插件会自动扩大父分组的尺寸，以确保子节点始终被完全包含在内。

这个插件极大地改善了用户在操作复杂的分组流程图时的体验，避免了手动调整分组大小的繁琐操作。

## 主要功能

-   **边界检测**：实时监测拖拽中的子节点是否越过父分组的内边距（`padding`）边界。
-   **自动调整尺寸**：当检测到边界穿越时，自动计算并扩大父分组的 `width` 和 `height`。
-   **中心点校正**：在调整尺寸的同时，同步调整父分组的 `x` 和 `y` 坐标，以确保分组是从正确的方向（例如，向左或向上）扩展，而不是仅仅向右下角变大。
-   **平滑体验**：通过 `resize_padding` 配置，可以在子节点接触边界时提供一个缓冲区域，使得尺寸调整更加平滑自然。
-   **性能保护**：在处理拖拽事件时，会临时解绑自身，避免因 API 调用而产生的无限循环，并在操作完成后重新绑定，确保插件的稳定运行。

## 实现详解

### `constructor({ lf })`

-   **`lf`**: LogicFlow 的实例。
-   **`this.padding`**: 定义了父分组的内边距。插件会基于这个内边距来计算子节点可移动的安全区域。
-   **`this.resize_padding`**: 一个额外的缓冲空间。当子节点触碰到边界时，父分组会额外扩大这个值的宽度或高度，避免子节点紧贴边缘。
-   **事件监听**: 在构造函数中，通过 `lf.on('node:drag', this.handleNodeDrag)` 来监听所有节点的拖拽事件。

### `handleNodeDrag({ data })`

这是插件的核心处理函数，在每次节点拖拽时被触发。

1.  **事件解绑**：在函数开始时，立刻通过 `this.lf.off('node:drag', this.handleNodeDrag)` 解绑事件。这是为了防止在后续调用 LogicFlow 的 API（如修改 `groupNode` 的属性）时，再次触发 `node:drag` 事件，从而导致死循环。

2.  **节点检查**：
    -   获取当前拖拽的节点 (`model`)。
    -   检查该节点是否属于一个分组（通过检查 `props.loop_parent_key`）。如果不是子节点，则不做任何处理。

3.  **边界计算**：
    -   获取子节点的位置和尺寸 (`x`, `y`, `width`, `height`)。
    -   获取父分组的位置和尺寸 (`groupX`, `groupY`, `groupWidth`, `groupHeight`)。
    -   计算出子节点的四条边（`nodeLeft`, `nodeRight`, `nodeTop`, `nodeBottom`）。
    -   根据父分组的尺寸和 `this.padding`，计算出父分组内部的安全区域四条边（`groupLeft`, `groupRight`, `groupTop`, `groupBottom`）。

4.  **尺寸调整**：
    -   通过 `if` 语句逐一判断子节点的边是否超出了父分组的安全区域。
    -   如果超出，计算需要增加的差值 (`diff`)，该差值等于超出部分的距离加上 `this.resize_padding`。
    -   根据超出的方向，更新 `newGroupWidth`、`newGroupHeight`，并同时调整 `newGroupX` 和 `newGroupY` 来校正中心点。例如，如果子节点向左超出了边界，父分组的宽度会增加，同时其中心点 `x` 会向左移动 `diff / 2`。
    -   设置 `needsResize = true` 标志。

5.  **应用更新**：
    -   如果 `needsResize` 为 `true`，则将计算出的新尺寸和位置（`newGroupX`, `newGroupY`, `newGroupWidth`, `newGroupHeight`）直接赋值给 `groupNode` 模型。

6.  **事件重绑**：在 `finally` 块中，无论前面的逻辑是否出错，都通过 `this.lf.on('node:drag', this.handleNodeDrag)` 重新绑定事件，确保插件在下一次拖拽时依然可用。

## 如何使用

该插件作为 LogicFlow 的一个标准插件被注册和使用。

1.  **导入插件**:
    ```javascript
    import { GroupAutoResize } from './plugins/group-auto-resize';
    ```

2.  **注册插件**:
    在初始化 LogicFlow 实例时，将其添加到 `plugins` 数组中。

    ```javascript
    import LogicFlow from '@logicflow/core';
    
    const lf = new LogicFlow({
      container: document.querySelector('#app'),
      // ... 其他配置
      plugins: [
        // ... 其他插件
        GroupAutoResize 
      ]
    });
    ```

之后，该插件会自动生效，无需任何额外调用。
