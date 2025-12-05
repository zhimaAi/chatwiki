# CustomKeyboard 插件

`CustomKeyboard` 是一个为 LogicFlow 工作流设计的插件，旨在增强用户体验，提供一套丰富的键盘快捷键，使得画布操作更加高效和直观。

## 功能特性

- **画布缩放**: 使用 `Ctrl` + `鼠标滚轮` 对画布进行平滑缩放。
- **画布平移**: 使用 `Shift` + `鼠标滚row` 对画布进行横向滚动。
- **节点复制与粘贴**:
  - `Ctrl` + `C`: 复制选中的一个或多个节点（开始节点和分组节点除外）。
  - `Ctrl` + `V`: 在画布上粘贴已复制的节点，每次粘贴会自动产生偏移，避免节点重叠。
- **节点删除**: 使用 `Delete` 键删除选中的节点和边。
- **套索选择**: 按住 `Ctrl` 键激活套索工具，方便用户框选多个节点。

## 如何使用

1.  **引入插件**:

    ```javascript
    import { CustomKeyboard } from './custom-keyboard';
    ```

2.  **注册插件**:

    在初始化 LogicFlow 实例时，将 `CustomKeyboard` 插件注册到 `plugins` 数组中。

    ```javascript
    import LogicFlow from '@logicflow/core';
    import { CustomKeyboard } from './custom-keyboard';

    const lf = new LogicFlow({
      // ... 其他 LogicFlow 配置
      keyboard: {
        enabled: false, // 开启自定义快捷键
        shortcuts: [
          // 屏蔽以下自带的快捷键
          { keys: ["cmd + c", "ctrl + c"], callback: () => {} },
          { keys: ["cmd + v", "ctrl + v"], callback: () => {} },
          { keys: ["backspace"], callback: () => {} }
        ]
      },
      plugins: [CustomKeyboard]
    });
    ```

> **重要**: 为了防止与 LogicFlow 的内置快捷键冲突，建议在初始化时通过 `keyboard` 选项屏蔽掉 `Ctrl+C`, `Ctrl+V`, 和 `Delete/Backspace` 的默认行为。

## 快捷键列表

| 快捷键                 | 功能                                               |
| ---------------------- | -------------------------------------------------- |
| `Ctrl` + `鼠标滚轮`    | 以鼠标指针位置为中心进行画布缩放                   |
| `Shift` + `鼠标滚轮`   | 横向平移画布                                       |
| `Ctrl` + `C`           | 复制选中的节点（开始节点和分组节点除外）           |
| `Ctrl` + `V`           | 粘贴节点，每次粘贴会自动偏移                     |
| `Delete`               | 删除选中的节点或边                                 |
| `Ctrl` (长按)          | 激活套索选择工具，用于框选多个节点和边             |

## 自定义事件

该插件会触发以下自定义事件，您可以在业务代码中监听这些事件以实现更复杂的功能。

-   **`custom:paste`**:
    -   **触发时机**: 当用户执行粘贴操作 (`Ctrl` + `V`) 时触发。
    -   **事件数据**:
        ```typescript
        {
          originalNodes: object[], // 被复制的原始节点数据
          basePoint: { x: number, y: number }, // 原始节点集合的左上角坐标
          pasteBasePoint: { x: number, y: number } // 本次粘贴操作的目标左上角坐标
        }
        ```
    -   **使用示例**:
        ```javascript
        lf.on('custom:paste', ({ originalNodes, pasteBasePoint }) => {
          // 在此处理粘贴逻辑，例如创建新节点
          const newNodes = originalNodes.map(node => {
            // ... 根据业务需求调整节点属性
            return {
              ...node,
              x: node.x - originalNodes[0].x + pasteBasePoint.x,
              y: node.y - originalNodes[0].y + pasteBasePoint.y,
            };
          });
          lf.addElements({ nodes: newNodes, edges: [] });
        });
        ```

-   **`custom:keyoard:delete`**:
    -   **触发时机**: 当用户按下 `Delete` 键时触发。
    -   **事件数据**: 无
    -   **使用示例**:
        ```javascript
        lf.on('custom:keyoard:delete', () => {
          const selectedElements = lf.getSelectElements();
          // 在此处理删除逻辑
          lf.deleteElements(selectedElements.nodes.map(node => node.id));
        });
        ```

## 主要属性和方法

### 属性

| 属性名          | 类型          | 描述                                         |
| --------------- | ------------- | -------------------------------------------- |
| `pluginName`    | `string`      | 插件名称，固定为 `'custom-keyboard'`。       |
| `lf`            | `object`      | LogicFlow 实例。                             |
| `isCtrlPressed` | `boolean`     | 标记 `Ctrl` 键是否被按下。                   |
| `zoom`          | `number`      | 当前画布的缩放值。                           |
| `clipboard`     | `object`      | 剪贴板，用于存储复制的节点数据和基点信息。   |
| `basePointX`    | `number`      | 粘贴节点的 X 轴偏移量，默认为 `20`。         |
| `basePointY`    | `number`      | 粘贴节点的 Y 轴偏移量，默认为 `20`。         |

### 方法

-   **`render()`**:
    插件的渲染方法，在插件注册时被 LogicFlow 自动调用。主要用于绑定键盘和鼠标事件监听器。

-   **`destroy()`**:
    插件的销毁方法，在 LogicFlow 实例销毁时自动调用。主要用于移除事件监听器，防止内存泄漏。