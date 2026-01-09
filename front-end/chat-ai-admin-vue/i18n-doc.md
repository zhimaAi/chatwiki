# 国际化（i18n）模块详细说明

## 1. 整体架构

### 1.1 技术栈

- **核心库**: vue-i18n 9.12.1
- **状态管理**: Pinia (locale store)
- **UI 组件库集成**: Ant Design Vue locale
- **持久化**: localStorage (via pinia-plugin-persistedstate)

### 1.2 核心文件结构

```
src/
├── locales/
│   ├── index.js                    # i18n 初始化入口
│   ├── helper.js                   # 工具函数（消息生成、HTML lang 设置）
│   └── lang/                       # 语言资源目录
│       ├── zh-CN.js                # 中文入口文件
│       ├── en-US.js                # 英文入口文件
│       ├── zh-CN/                  # 中文资源
│       │   ├── common.json
│       │   ├── layout.json
│       │   ├── routes/
│       │   └── views/
│       └── en-US/                  # 英文资源
│           ├── common.json
│           ├── layout.json
│           ├── routes/
│           └── views/
├── hooks/web/
│   ├── useI18n.js                  # 翻译 hook
│   └── useLocale.js                # 语言切换 hook
└── stores/modules/
    └── locale.js                   # 语言状态管理
```

### 1.3 架构图

```
┌─────────────────────────────────────────────────────────┐
│                     Vue Components                       │
│  (使用 useI18n('namespace') 获取翻译函数)                 │
└──────────────────┬──────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│                    useI18n Hook                          │
│  - 处理 namespace 前缀                                   │
│  - 调用 vue-i18n 的 t() 方法                             │
└──────────────────┬──────────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────────┐
│                   vue-i18n 实例                          │
│  - i18n.global.t(key)                                    │
│  - 管理语言包                                             │
└──────────────────┬──────────────────────────────────────┘
                   │
         ┌─────────┴─────────┐
         ▼                   ▼
┌─────────────────┐  ┌─────────────────┐
│  zh-CN.js       │  │  en-US.js       │
│  (动态加载)     │  │  (动态加载)     │
└─────────────────┘  └─────────────────┘
         │                   │
         └─────────┬─────────┘
                   │
         ┌─────────┴─────────┐
         ▼                   ▼
┌─────────────────┐  ┌─────────────────┐
│  zh-CN/*.json   │  │  en-US/*.json   │
│  (通过 genMessage│  │  (通过 genMessage│
│   合并)         │  │   合并)         │
└─────────────────┘  └─────────────────┘

语言切换流程：
LocaleDropdown.vue
    │
    ▼
useLocale.changeLocale()
    │
    ├─→ 动态 import 目标语言包
    ├─→ i18n.setLocaleMessage()
    ├─→ i18n.global.locale = locale
    ├─→ localeStore.setCurrentLocale()
    ├─→ setHtmlPageLang()
    └─→ window.location.reload()
```

## 2. 核心文件详细说明

### 2.1 src/locales/index.js - i18n 初始化

**职责**: 创建和配置 vue-i18n 实例

**核心配置**:

- `legacy: false` - 使用 Composition API 模式
- `locale` - 当前语言（从 localeStore 读取）
- `fallbackLocale` - 回退语言
- `messages` - 语言包对象
- `availableLocales` - 可用语言列表
- `silentTranslationWarn` - 生产环境静默警告
- `missingWarn` - 开发环境显示缺失警告

**初始化流程**:

1. 从 localeStore 获取当前语言
2. 动态导入对应的语言入口文件 (zh-CN.js 或 en-US.js)
3. 设置 HTML lang 属性
4. 创建 i18n 实例并挂载到 app

### 2.2 src/locales/helper.js - 工具函数

#### setHtmlPageLang(locale)

- 设置 HTML 文档的 lang 属性
- 用于 SEO 和浏览器语言检测

#### genMessage(langs, prefix)

- **功能**: 将所有 JSON 语言文件合并成统一的消息对象
- **输入**:
  - `langs` - import.meta.glob 返回的模块对象
  - `prefix` - 语言前缀 ('zh-CN' 或 'en-US')
- **输出**: 嵌套对象，按照文件路径组织
- **转换规则**:
  - `./zh-CN/common.json` → `{ common: {...} }`
  - `./zh-CN/views/user/login.json` → `{ views: { user: { login: {...} } } }`
- **使用 lodash.set**: 自动创建嵌套结构

### 2.3 src/hooks/web/useI18n.js - 翻译 Hook

**职责**: 提供带有 namespace 的翻译函数

**参数**:

- `namespace` - 命名空间前缀（如 'views.user.login'）

**返回**:

- `t` - 翻译函数，自动添加 namespace 前缀
- 其他 vue-i18n 方法（te, rt 等）

**逻辑**:

```javascript
// 如果没有 namespace
t('key') → 'key'

// 如果有 namespace
useI18n('views.user.login')
t('accountLogin') → t('views.user.login.accountLogin')

// 如果 key 本身包含 namespace
t('common.add') → t('common.add')  // 不重复添加
```

**使用示例**:

```vue
<script setup>
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('views.user.login')
// 使用: {{ t('accountLogin') }}
</script>
```

### 2.4 src/hooks/web/useLocale.js - 语言切换 Hook

**职责**: 提供语言切换功能

**核心方法**:

- `changeLocale(locale)` - 切换语言

**切换流程**:

1. 动态 import 目标语言包: `import(`../../locales/lang/${locale}.js`)`
2. 调用 `i18n.global.setLocaleMessage(locale, langModule.default)`
3. 更新 i18n 的当前语言（兼容 legacy 和 composition 模式）
4. 更新 localeStore 状态
5. 设置 HTML lang 属性

### 2.5 src/stores/modules/locale.js - 语言状态管理

**职责**: 管理语言状态和 Ant Design Vue locale

**State**:

```javascript
{
  currentLocale: {
    lang: 'zh-CN',              // 当前语言代码
    antvLocale: zhCn            // Ant Design Vue 语言包
  },
  localeMap: [                  // 支持的语言列表
    { lang: 'zh-CN', name: '简体中文' },
    { lang: 'en-US', name: 'English' }
  ]
}
```

**Getters**:

- `getCurrentLocale()` - 获取当前语言配置
- `getLocaleMap()` - 获取所有支持的语言
- `getSelectedLocale()` - 获取当前语言对象（从 localeMap 中查找）

**Actions**:

- `setCurrentLocale(localeMap)` - 设置当前语言
  - 更新 lang
  - 更新 antvLocale（根据 antvLocaleMap 映射）
  - 持久化到 localStorage

**Ant Design Vue 集成**:

```javascript
const antvLocaleMap = {
  'zh-CN': zhCn, // ant-design-vue/es/locale/zh_CN
  'en-US': en, // ant-design-vue/es/locale/en_US
  en: en // 兼容旧版本
}
```

## 3. 语言资源组织

### 3.1 支持的语言

- **zh-CN** (简体中文)
- **en-US** (美式英语)

### 3.2 语言文件结构

#### 入口文件 (zh-CN.js / en-US.js)

```javascript
import { genMessage } from '../helper'
const modulesFiles = import.meta.glob('./zh-CN/**/*.json', { eager: true })
export default {
  ...genMessage(modulesFiles, 'zh-CN')
}
```

#### 资源目录结构

```
lang/
├── zh-CN/
│   ├── common.json              # 通用翻译（错误消息、操作按钮等）
│   ├── layout.json              # 布局组件翻译
│   ├── routes/
│   │   └── basic.json           # 路由相关翻译
│   └── views/
│       └── user/
│           ├── login.json       # 登录页翻译
│           ├── model.json       # 模型管理翻译
│           ├── left-menus.json  # 左侧菜单翻译
│           ├── enterprise.json  # 企业设置翻译
│           └── ...
└── en-US/
    └── (相同结构)
```

### 3.3 翻译 Key 命名规范

**规则**:

- 使用点号分隔的层级结构
- 对应文件路径：`views/user/login.json` 中的 key 访问路径为 `views.user.login.key`
- 使用小写字母和下划线
- 语义化命名

**示例**:

```json
// common.json
{
  "add": "添加",
  "delete_successful": "删除成功",
  "api_request_failed": "请求出错，请稍候重试"
}

// layout.json
{
  "header": {
    "tooltip_error_log": "错误日志",
    "dropdown_item_login_out": "退出系统"
  }
}

// views/user/login.json
{
  "account_login": "账号登录",
  "please_number": "请输入账号",
  "profession_one_title": "金融行业"
}
```

### 3.4 翻译内容格式

**类型**:

- 简单字符串: `"add": "添加"`
- 嵌套对象: `"header": { "tooltip": "..." }`
- 占位符支持（vue-i18n 标准语法）

## 4. 语言切换完整流程

### 4.1 初始化流程（应用启动）

```
main.js
  │
  ▼
setupI18n(app)
  │
  ├─→ createI18nOptions()
  │   │
  │   ├─→ localeStore.getCurrentLocale (从 localStorage 读取)
  │   ├─→ import(`./lang/${locale.lang}.js`)
  │   ├─→ setHtmlPageLang(locale.lang)
  │   └─→ return i18n options
  │
  ├─→ i18n = createI18n(options)
  └─→ app.use(i18n)
```

### 4.2 切换语言流程

```
用户点击语言切换
  │
  ▼
LocaleDropdown.vue.setLang()
  │
  ▼
useLocale().changeLocale(key)
  │
  ├─→ import(`../../locales/lang/${locale}.js`)  // 动态加载
  ├─→ i18n.global.setLocaleMessage(locale, messages)
  ├─→ setI18nLanguage(locale)
  │   │
  │   ├─→ i18n.global.locale = locale  (兼容 legacy 模式)
  │   ├─→ localeStore.setCurrentLocale({ lang })
  │   └─→ setHtmlPageLang(locale)
  │
  └─→ window.location.reload()  // 刷新页面
```

### 4.3 页面刷新后

1. localeStore 从 localStorage 恢复语言设置
2. setupI18n 重新初始化，加载对应语言包
3. 所有组件使用新的语言资源

## 5. 与其他模块的集成

### 5.1 Ant Design Vue 集成

- **位置**: `src/components/global-config-provider/index.vue`
- **方式**: 通过 ConfigProvider 的 locale 属性
- **数据源**: localeStore.currentLocale.antvLocale

```vue
<ConfigProvider :locale="currentLocale.antvLocale">
  <slot />
</ConfigProvider>
```

### 5.2 Pinia 状态管理集成

- **持久化**: 使用 `pinia-plugin-persistedstate`
- **存储 key**: `lang`
- **存储位置**: localStorage

### 5.3 Vue Router 集成

- **路由标题**: 可以使用 i18n 翻译
- **面包屑**: layout.json 中定义了相关翻译

### 5.4 组件使用方式

**方式 1: 带命名空间（推荐）**

```vue
<script setup>
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('views.user.login')
</script>

<template>
  <div>{{ t('accountLogin') }}</div>
  <!-- 自动添加前缀 -->
</template>
```

**方式 2: 不带命名空间**

```vue
<script setup>
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n()
</script>

<template>
  <div>{{ t('views.user.login.accountLogin') }}</div>
  <!-- 完整路径 -->
</template>
```

## 6. 配置项说明

### 6.1 i18n 配置

- **legacy: false** - 启用 Composition API
- **fallbackLocale** - 回退语言（默认当前语言）
- **sync: true** - 同步语言到根组件
- **silentTranslationWarn** - 生产环境静默翻译警告
- **missingWarn** - 开发环境显示缺失翻译警告
- **silentFallbackWarn** - 生产环境静默回退警告

### 6.2 Ant Design Vue 配置

- 通过 localeStore 自动映射
- 支持的语言：zh-CN, en-US
- 兼容旧版本 'en'

## 7. 注意事项和最佳实践

### 7.1 添加新翻译

1. 在对应语言的 JSON 文件中添加 key-value
2. 确保中英文文件保持同步
3. 使用语义化的 key 命名
4. 避免硬编码文本

### 7.2 添加新语言

1. 在 `src/locales/lang/` 下创建新语言目录
2. 创建对应的 JSON 文件
3. 创建语言入口文件 (如 `fr-FR.js`)
4. 在 `locale.js` 中更新 `localeMap` 和 `antvLocaleMap`

### 7.3 性能优化

- 使用动态导入语言包，减少初始加载体积
- 按需加载，只加载当前需要的语言
- JSON 文件格式，解析速度快

### 7.4 开发建议

- 在开发环境启用警告，方便发现缺失翻译
- 使用统一的命名规范
- 保持文件结构清晰，按功能模块组织
- 定期检查翻译完整性

## 8. 常见问题

### Q1: 为什么语言切换后需要刷新页面？

A: 为了确保所有组件和第三方库（如 Ant Design Vue）都能正确应用新语言。

### Q2: 如何在组件外使用翻译？

A: 可以直接导入 i18n 实例：

```javascript
import { i18n } from '@/locales'
const t = i18n.global.t
```

### Q3: 如何支持占位符？

A: 使用 vue-i18n 的标准语法：

```json
{
  "welcome": "欢迎, {name}!"
}
```

```javascript
t('welcome', { name: 'John' })
```

### Q4: 如何处理复数形式？

A: 使用 vue-i18n 的复数语法：

```json
{
  "item": "0 项 | 1 项 | {n} 项"
}
```

## 9. 相关文件清单

### 核心文件

- `src/locales/index.js` - i18n 初始化
- `src/locales/helper.js` - 工具函数
- `src/hooks/web/useI18n.js` - 翻译 hook
- `src/hooks/web/useLocale.js` - 语言切换 hook
- `src/stores/modules/locale.js` - 状态管理

### 语言资源

- `src/locales/lang/zh-CN.js` - 中文入口
- `src/locales/lang/en-US.js` - 英文入口
- `src/locales/lang/zh-CN/common.json` - 通用中文翻译
- `src/locales/lang/en-US/common.json` - 通用英文翻译
- `src/locales/lang/zh-CN/layout.json` - 布局中文翻译
- `src/locales/lang/en-US/layout.json` - 布局英文翻译
- `src/locales/lang/zh-CN/views/**` - 各页面中文翻译
- `src/locales/lang/en-US/views/**` - 各页面英文翻译

### 集成点

- `src/layouts/AdminLayout/compoents/locale-dropdown.vue` - 语言切换组件
- `src/components/global-config-provider/index.vue` - Ant Design Vue 集成
- `src/main.js` - 应用入口，初始化 i18n

## i18n Agent

```markdown
# Role
你是一个精通 Vue/React 国际化架构的前端专家。

# Context
我的项目 i18n 目录结构如下（基于 `src/locales/lang`）：
- `src/locales/lang/zh-CN/`：中文源文件
- `src/locales/lang/en-US/`：英文目标文件
- 结构特点：模块化嵌套，JSON 文件路径与源码目录结构存在映射关系。

# Goal
请分析**当前打开的代码文件**（Active File），提取其中的硬编码中文，并将其抽取到对应的 i18n 模块文件中，并自动重构代码。

# Task Workflow

1. **Analyze Path & Namespace (分析路径与命名空间)**:
   - 获取当前文件的相对路径（例如 `src/views/user/model/components/AddModelAlert.vue`）。
   - **推导 Full Key Path (完整键路径)**: 
     - 结合目录结构，推导出完整的 i18n 路径。
     - 例如：`views.user.model.components.add-model-alert`。
   - **推导 JSON 存储位置**:
     - 基于路径的前几层找到对应的 JSON 文件（例如 `locales/lang/{lang}/views/user/model/components/add-model-alert.json`）。

   - 获取当前文件的相对路径（例如 `src/views/library-search/index.vue`）。
     - **推导 Full Key Path (完整键路径)**: 
       - 结合目录结构，推导出完整的 i18n 路径。
       - 例如：`views.library-search.index`。
     - **推导 JSON 存储位置**:
       - 基于路径的前几层找到对应的 JSON 文件（例如 `locales/lang/{lang}/views/library-search/index.json`）。

   - 获取当前文件的相对路径（例如 `src/views/login/login.vue`）。
     - **推导 Full Key Path (完整键路径)**: 
       - 结合目录结构，推导出完整的 i18n 路径。
       - 例如：`views.login.login`。
     - **推导 JSON 存储位置**:
       - 基于路径的前几层找到对应的 JSON 文件（例如 `locales/lang/{lang}/views/login/login.json`）。

   - 获取当前文件的相对路径（例如 `src/views/public-library/home/index.vue`）。
     - **推导 Full Key Path (完整键路径)**: 
       - 结合目录结构，推导出完整的 i18n 路径。
       - 例如：`views.public-library.home.index`。
     - **推导 JSON 存储位置**:
       - 基于路径的前几层找到对应的 JSON 文件（例如 `locales/lang/{lang}/views/public-library/home/index.json`）。

2. **Extract & Translate (提取与翻译)**:
   - 提取硬编码中文。
   - **JSON 处理**:
     - 如果目标 JSON 文件不存在，请自动创建。
     - 同时向 `zh-CN` 和 `en-US` 追加内容。
     - Key 命名使用 `snake_case` (例如 `submit_btn`)。

3. **Code Refactor (代码重构 - 核心规则)**:
   - **统一采用命名空间模式 (Namespace Mode)** (不区分层级深度):
     - **Import**: 必须使用自定义 Hook: `import { useI18n } from '@/hooks/web/useI18n'`。
     - **Script**: 提取 Key 的前缀路径（通常是文件路径映射）作为 Namespace 传入 Hook。
       - *示例*: 如果完整 Key 路径是 `views.user.model.components.add-model-alert.title`。
       - *Namespace*: `views.user.model.components.add-model-alert`。
       - *写法*: `const { t } = useI18n('views.user.model.components.add-model-alert')`。
       - **Template**: 直接使用短 Key，例如 `t('title')` (不要使用完整路径)。

       - *示例*: 如果完整 Key 路径是 `views.library-search.index.all`。
       - *Namespace*: `views.library-search`。
       - *写法*: `const { t } = useI18n('views.library-search.index')`。
       - **Template**: 直接使用短 Key，例如 `t('all')` (不要使用完整路径)。
     

   - **变量处理**:
     - 如果包含变量 `${name}`，JSON 中保留 `{name}`，代码中使用 `t('key', { name: val })`。

# Constraints
- **强制约束**: 所有 i18n 引入必须来自 `@/hooks/web/useI18n`，**禁止**从 `vue-i18n` 引入。
- 保持 `zh-CN` 和 `en-US` 结构完全一致。
- 命名空间字符串应基于文件路径生成（kebab-case 连接）。
- 英语翻译要简练。
- 直接生成代码并 Apply，无需确认。
- 翻译完成后需要再次确认翻译文件路径是否正确。
- 始终用中文沟通。 
```
