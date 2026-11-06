# ChatWiki 主题提取

## 当前范围

本目录保存默认主题、天青蓝和扶光金的配置。应用已接入主题 Store、动态 Provider、根节点 CSS 变量和客服左侧外观菜单。
`src/config/theme.js` 保留为默认主题的原始 Ant 配置。业务页面固定颜色、图表和全量静态反馈调用仍待后续适配，当前不代表全应用换肤已完成。
默认主题标识为 `default`，天青蓝为 `webflow`，扶光金为 `claude`。

## 配置结构

| 字段 | 含义 |
| --- | --- |
| `id` / `name` | 稳定标识与中文展示名；展示名由产品需求提供 |
| `sources` | 参考规范 URL 或项目内来源文件 |
| `palette` | 来源中明确给出的颜色，保留原始用途名称 |
| `semantic` | 新主题构建 `antTheme.token` 时使用的语义映射；运行时不再叠加到最终 Token 上 |
| `legacyBusiness` | 仅默认主题使用的历史业务背景、文字与普通链接色，独立于组件库 Token |
| `typography` | 中文界面字体回退方案及原文字体参考；数值字号和字距单位为 px |
| `radius` | 原文或现有组件圆角；数值单位为 px，`Range` 表示范围 |
| `shadow` | 来源中明确给出的 CSS 阴影 |
| `reference` | 参考间距、断点、动效等资料，不会直接改变应用布局 |
| `antTheme` | 供 ConfigProvider 使用的 Token 与组件配置；未显式配置的状态由组件库推导 |

侧栏和用户菜单直接使用应用统一的 `--app-*` 变量，不再维护侧栏专用配色配置。悬停、选中等状态使用当前 Provider 推导的 Token，不作为规范原文值。默认主题另有 `chart` 记录图表系列颜色，新主题的图表配置仍待补齐。

## 应用级接入

- `src/stores/modules/theme.js` 保存当前主题 ID，并通过现有 `useStorage` 封装持久化到 `chatwiki_theme_v1`。
- 启动时在首次挂载前恢复选择；`GlobalConfigProvider` 响应式读取对应 `antTheme`。
- `App.vue` 中的 `a-app` 已放入 Provider 内部；外观菜单通过 `App.useApp()` 取得带主题上下文的消息实例。
- `theme-variables.vue` 使用公开的 `theme.useToken()` 读取当前 Provider 的最终 Token，将统一的 `--app-*` 变量写入 html，使侧栏、用户菜单和 body 浮层也能继承。新主题的颜色与 UI 字体不再被原始 `semantic` 二次覆盖；默认主题的历史业务值通过独立的 `legacyBusiness` 显式兼容。
- 侧栏直接使用 `--app-bg-container`、`--app-text-primary`、`--app-text-secondary`、`--app-border`、`--app-primary` 和 `--app-selected`。普通悬停使用 `--app-bg-hover`，底部入口的较强悬停填充使用 `--app-fill-secondary`，对应最终 Token 的 `colorFillSecondary`。
- 会话记录、消息列表和渠道页面的已有主题变量引用也统一使用 `--app-primary`、`--app-border`，不再依赖客服模块专用变量或其固定颜色回退。
- 两套新主题显式设置 `colorLink`、`colorLinkHover`、`colorLinkActive`，普通链接与 Ant 组件共享链接 Token；不通过修改 `colorInfo` 改变信息提示色。缺少来源状态时，通过公开的 `theme.defaultAlgorithm` 和 `theme.defaultSeed` 一次性取得主色状态，填入链接配置；运行时仍由响应式 `ConfigProvider.theme` 负责切换。
- `src/assets/main.css` 的普通链接规则读取 `--app-link*`，保持低选择器优先级，让组件自身的禁用、危险等状态样式正常参与层叠。`src/assets/theme.css` 只为两套新主题映射已有语义变量和 UI 字体，不再通过主题选择器提高全局链接规则优先级。
- 同源标签页通过 storage 事件同步；事件监听随根组件卸载移除。切换不刷新页面，不给路由组件增加主题 key。
- 主题偏好仅使用 `chatwiki_theme_v1`，不读取、迁移或清理旧键。没有新键时使用默认主题，仅保存旧版外观偏好的用户需要重新选择。
- 无效存储回退默认主题；写入失败保留当前页面选择并提示无法保存。清除主题偏好时，同源其他标签页恢复默认主题。

## 默认主题的来源与边界

| 内容 | 来源 |
| --- | --- |
| 全局主色 `#2475fc`、中文字体、现有 Ant 配置 | `src/config/theme.js` |
| body 字号 15px、行高 1.6、工作流文字变量 | `src/assets/base.css` |
| 普通文字 `#595959`、删除文字 `#fb363f` | `src/assets/main.css` |
| 白色布局背景、顶栏边框 `#f2f4f7` | `src/layouts/AdminLayout/index.vue` |
| 客服工作区背景 `#f5f6f8` | `src/views/customer-service/index.vue` |
| 侧栏字体、菜单圆角 | `src/views/customer-service/components/customer-sidebar.vue` |
| 菜单浮层圆角和阴影 | `src/views/customer-service/components/customer-user-menu.vue` |
| 统计卡片背景、6px 圆角、图表系列颜色 | `src/views/robot/robot-config/statistical_analysis/hit-statics.vue` |

默认侧栏的背景、文字、边框及状态颜色统一使用应用主题变量，强调色跟随应用主色 `#2475fc`。这里只提取已经定位的公共基线和代表性组件，不代表已完成全应用固定颜色的归并。
默认 `antTheme` 直接引用现有显式配置，继续由已安装组件库提供其他默认值。`legacyBusiness` 只兼容原业务背景、文字及普通链接色，不会合并到 `antTheme.token`；这些应用级变量也供侧栏直接使用。
`base.css` 中的系统深色模式仍按原代码工作。本目录没有新增深色主题，也没有改变系统偏好行为。

## 天青蓝

来源：<https://fchangjun.github.io/awesome-design-md-cn/design-md/webflow/DESIGN.md>

- 主色、悬停蓝、辅助色、中性色、白色背景与五层阴影按原文提取。
- 链接基础色映射为品牌主色，悬停色沿用原文按钮悬停蓝；这是 ChatWiki 的链接映射。按下色使用官方主色算法推导。
- 页面、卡片及浮层统一映射为白色；原文未为后台的这些层级定义独立底色。
- 通用控件圆角取 4px，是在原文功能元素范围内做出的应用选择；卡片仍保留 4–8px 参考范围。
- `green`、`yellow`、`red` 是原文辅助色，不直接当作后台成功、警告和错误色。
- 原文徽标 10% 主色背景及按钮 6px 悬停位移只记录为参考，不直接应用到全部按钮和菜单。

## 扶光金

来源：<https://fchangjun.github.io/awesome-design-md-cn/design-md/claude/DESIGN.md>

- 保留原文陶土色 `#c96442`，产品名称“扶光金”不意味着把原色替换成金色。
- 链接基础色使用陶土主色，悬停和按下色使用官方主色算法推导，不标记为规范原文值。
- 页面取 Parchment，卡片取 Ivory；浮层取 Ivory 是 ChatWiki 的映射。
- 通用控件圆角取标准按钮的 8px；Input 的 12px 单独保留。Select、日期选择器等仍需后续逐项适配。
- 原文 `Ring Subtle` 给出 `#dedc01`，与其暖灰描述及周边色板存在疑点，本轮不纳入有效配色。
- 原文正文行高既列出 1.25–1.6，又要求不低于 1.4。本轮保留该范围为参考，未覆盖应用全局行高；接入时需统一选择。
- 深色表面仅作为原文配色资料保留，不构成新增深色主题。
- 原文蓝色聚焦边框单独记录为 `reference.focusBorder`，后续需要按输入类组件适配，不能用它覆盖主色。

## 应用适配与待补齐项

以下内容是应用映射或后续工作，不应标记为原文直接提取值：

1. `semantic` 的字段分配、通用控件圆角，以及 Claude 的 Ivory 浮层背景。
2. UI 字体链追加现有中文字体；品牌字体未下载、未安装，也未增加 `@font-face`。中文衬线标题和代码字体的完整回退方案尚未确定。
3. `antTheme` 提供基础 Token、显式链接状态与必要的组件配置。原文没有完整定义的状态通过组件库算法补齐；链接推导值会显式写入链接 Token，其余缺省状态由 Provider 推导。这不代表它们已经过视觉确认。
4. 新主题的聊天气泡、业务页面固定颜色、编辑器、工作流画布和图表系列仍需结合业务补齐。
5. 不将原文官网大标题、间距、断点和动效自动套用到后台；`reference` 不作为全局尺寸配置。
6. 静态 `message`、`notification`、`Modal.confirm` 调用尚未全量迁移；本步只接入外观菜单的消息实例。全应用页面与第三方组件仍需后续逐项适配。

当前采用官方动态主题入口与最终 Token 消费方式；业务自定义区域仍使用应用 CSS 变量衔接，不代表全应用所有页面已完成主题适配。

未执行验证。
