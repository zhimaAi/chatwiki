# 主题设计文档

| 主题 | 设计文档 | 文档来源 | 应用配置 |
| --- | --- | --- | --- |
| 默认主题 | [default/DESIGN.md](default/DESIGN.md) | 根据项目现有配置整理 | `src/config/themes/default.js` |
| 天青蓝（Webflow） | [webflow/DESIGN.md](webflow/DESIGN.md) | 用户提供的外部规范原文 | `src/config/themes/webflow.js` |
| 扶光金（Claude） | [claude/DESIGN.md](claude/DESIGN.md) | 用户提供的外部规范原文 | `src/config/themes/claude.js` |

## 外部原文来源

- [Webflow DESIGN.md](https://fchangjun.github.io/awesome-design-md-cn/design-md/webflow/DESIGN.md)
- [Claude DESIGN.md](https://fchangjun.github.io/awesome-design-md-cn/design-md/claude/DESIGN.md)

两份外部文档按原文保存，不在原文中追加 ChatWiki 的映射规则。默认主题文档是项目整理稿。

## 在项目中使用

外部规范包含官网排版、布局和动效参考，不会自动全部应用到后台。ChatWiki 的映射、取舍和已知边界见 [主题配置说明](../src/config/themes/README.md)，运行时以主题配置及 Provider 最终 Token 为准。

适配组件时，Ant Design Vue 组件沿用现有动态主题；自定义样式直接使用统一的 `--app-*` 语义变量。不要根据外部文档向组件写入固定主题色，也不要新增模块专属主题变量。

## Agent 主题适配提示

将下面的【组件路径】替换为目标组件路径后，发送给 Agent：

```text
请对【组件路径】进行主题适配，使其跟随应用的默认、天青蓝（Webflow）、扶光金（Claude）主题切换。

实现要求：
1. 先阅读项目 AGENTS.md 及其引用的规范、design-system/README.md 和相关主题设计文档，了解现有主题机制和组件的复用范围，再修改代码。设计文档与当前应用映射存在差异时，以现有主题配置及 Provider 最终 Token 为准。
2. Ant Design Vue 组件优先通过现有 ConfigProvider 动态主题生效；自定义样式直接使用现有 --app-* 语义变量，不新增 --customer-* 等模块专属主题变量。
3. 将固定的背景、文字、边框、圆角和阴影替换为对应主题变量，覆盖普通、悬停、选中、禁用、聚焦等组件实际存在的状态。
4. 检查 SVG 图标是否写死颜色；单色图标优先使用 currentColor。保留业务状态色、多色图标和品牌色的原有含义。
5. 若缺少必要的语义变量，在现有统一变量出口中补充，并从 Ant Design Vue 最终 Token 派生，不在各主题文件中重复定义组件样式。
6. 保持现有布局、交互、权限和业务逻辑，仅修改主题适配需要的代码。共享组件需说明其他使用方受到的影响。
7. 未经当前请求明确授权，不启动服务，不运行浏览器、测试、lint、类型检查或构建，不自动提交代码。

完成后简要说明改动范围、变量映射、共享影响及未适配项，并明确注明“未执行验证”。
```
