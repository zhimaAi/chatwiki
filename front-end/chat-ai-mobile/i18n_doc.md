# Chat AI Mobile 多语言翻译实现文档

## 一、技术栈

- **Vue 3** - 前端框架
- **vue-i18n** (v9.12.1) - 国际化核心库
- **Pinia** - 状态管理（用于存储当前语言设置）
- **Vant** - UI 组件库（支持多语言）
- **@intlify/unplugin-vue-i18n** (v3.0.1) - Vite 插件优化

## 二、目录结构

```
src/
├── locales/                          # 多语言主目录
│   ├── index.ts                      # i18n 初始化和导出
│   ├── config.ts                     # 语言配置和类型定义
│   ├── helper.ts                     # 辅助函数
│   └── lang/                         # 语言文件目录
│       ├── zh-CN.ts                  # 中文语言入口
│       ├── en-US.ts                  # 英文语言入口
│       ├── zh-CN/                    # 中文翻译文件
│       │   ├── common.json           # 通用翻译
│       │   ├── layout.json           # 布局相关翻译
│       │   ├── routes/
│       │   │   └── basic.json        # 路由翻译
│       │   └── views/
│       │       └── user/
│       │           ├── account.json  # 用户账户相关
│       │           └── model.json    # 模型相关
│       └── en-US/                    # 英文翻译文件
│           ├── common.json
│           ├── layout.json
│           ├── routes/basic.json
│           └── views/user/
│               ├── account.json
│               └── model.json
├── hooks/web/
│   ├── useI18n.ts                    # i18n Hook
│   └── useLocale.ts                  # 语言切换 Hook
└── stores/modules/
    └── locale.ts                     # 语言状态管理
```

## 三、核心实现

### 3.1 初始化流程

```typescript
// main.ts
import { setupI18n } from '@/locales'

const setupAll = async () => {
  const app = createApp(App)
  await setupI18n(app)
  // ... 其他初始化
}
```

**初始化步骤**：
1. 创建 Vue 应用实例
2. 调用 `setupI18n(app)` 初始化 i18n
3. 从 Pinia store 中读取当前语言
4. 动态导入对应语言文件
5. 创建 i18n 实例并注入到应用

### 3.2 i18n 配置 (`src/locales/index.ts`)

```typescript
import { createI18n } from 'vue-i18n'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { setHtmlPageLang } from './helper'

export let i18n: ReturnType<typeof createI18n>

const createI18nOptions = async (): Promise<I18nOptions> => {
  const localeStore = useLocaleStoreWithOut()
  const locale = localeStore.getCurrentLocale
  const localeMap = localeStore.getLocaleMap

  // 动态导入语言文件
  const defaultLocal = await import(`./lang/${locale.lang}.ts`)
  const message = defaultLocal.default ?? {}

  // 设置 HTML lang 属性
  setHtmlPageLang(locale.lang)

  return {
    legacy: false,                    // 使用 Composition API 模式
    locale: locale.lang,               // 当前语言
    fallbackLocale: locale.lang,       // 回退语言
    messages: {
      [locale.lang]: message
    },
    availableLocales: localeMap.map((v: any) => v.lang),
    sync: true,
    silentTranslationWarn: true,       // 静默翻译警告
    missingWarn: false,
    silentFallbackWarn: true
  }
}

export const setupI18n = async (app: App<Element>) => {
  const options = await createI18nOptions()
  i18n = createI18n(options) as I18n
  app.use(i18n)
}
```

### 3.3 语言配置 (`src/locales/config.ts`)

```typescript
export const localeMap = {
  'zh-CN': 'zh-CN',
  'en-US': 'en-US'
} as const

export type LocaleType = keyof typeof localeMap

export const localeList = [
  {
    lang: localeMap['en-US'],
    label: 'English',
    icon: '🇺🇸',
    title: 'Language'
  },
  {
    lang: localeMap['zh-CN'],
    label: '简体中文',
    icon: '🇨🇳',
    title: '语言'
  }
] as const
```

### 3.4 语言文件加载机制

**语言入口文件** (`src/locales/lang/zh-CN.ts`):
```typescript
import { genMessage } from '../helper'

// 使用 Vite 的 import.meta.glob 动态导入所有 JSON 文件
const modulesFiles = import.meta.glob<Recordable>('./zh-CN/**/*.json', { eager: true })

export default {
  ...genMessage(modulesFiles, 'zh-CN')
}
```

**辅助函数** (`src/locales/helper.ts`):
```typescript
import { set } from 'lodash-es'

export function genMessage(langs: Record<string, Record<string, any>>, prefix = 'lang') {
  const obj = {}

  Object.keys(langs).forEach((key) => {
    const langFileModule = langs[key].default
    let fileName = key.replace(`./${prefix}/`, '').replace(/^\.\//, '')
    const lastIndex = fileName.lastIndexOf('.')
    fileName = fileName.substring(0, lastIndex)

    const keyList = fileName.split('/')
    const moduleName = keyList.shift()
    const objKey = keyList.join('.')

    // 将文件路径转换为嵌套对象结构
    // 例如: zh-CN/common.json -> { common: {...} }
    if (moduleName) {
      if (objKey) {
        set(obj, moduleName, obj[moduleName] || {})
        set(obj[moduleName], objKey, langFileModule)
      } else {
        set(obj, moduleName, langFileModule || {})
      }
    }
  })
  return obj
}
```

### 3.5 状态管理 (`src/stores/modules/locale.ts`)

```typescript
import { defineStore } from 'pinia'
import { Locale } from 'vant'
import { Storage } from '@/utils/Storage'
import zhCn from 'vant/es/locale/lang/zh-CN'
import enUS from 'vant/es/locale/lang/en-US'

const vantLocaleMap = {
  'zh-CN': zhCn,
  'en-US': enUS
}

export const useLocaleStore = defineStore('locales', {
  state: (): LocaleState => {
    return {
      currentLocale: {
        lang: Storage.get('lang') || 'zh-CN',  // 从本地存储读取
        vantLocale: vantLocaleMap[Storage.get('lang') || 'zh-CN']
      },
      localeMap: [
        { lang: 'zh-CN', name: '简体中文' },
        { lang: 'en-US', name: 'English' }
      ]
    }
  },
  getters: {
    getCurrentLocale(): LocaleDropdownType {
      return this.currentLocale
    },
    getLocaleMap(): LocaleDropdownType[] {
      return this.localeMap
    },
  },
  actions: {
    setCurrentLocale(localeMap) {
      this.currentLocale.lang = localeMap?.lang
      this.currentLocale.vantLocale = vantLocaleMap[localeMap?.lang]

      // 同步更新 Vant 组件库语言
      Locale.use(localeMap?.lang, this.currentLocale.vantLocale)

      // 持久化存储
      Storage.set('lang', localeMap?.lang)
    }
  }
})
```

## 四、使用方法

### 4.1 在组件中使用

```vue
<script setup lang="ts">
import { useI18n } from '@/hooks/web/useI18n'

// 方式1: 基础使用
const { t } = useI18n()

// 方式2: 带命名空间使用
const { t } = useI18n('layout')
</script>

<template>
  <!-- 方式1: 完整路径 -->
  <div>{{ t('common.add') }}</div>

  <!-- 方式2: 命名空间 -->
  <div>{{ t('header.home') }}</div>
</template>
```

### 4.2 在 JS 文件中使用

```typescript
import { useI18n } from '@/hooks/web/useI18n'

export function getErrorMsg(error) {
  const { t } = useI18n()

  switch (error.response.status) {
    case 401:
      return t('common.errMsg401')
    case 404:
      return t('common.errMsg404')
    default:
      return t('common.apiRequestFailed')
  }
}
```

### 4.3 切换语言

```typescript
import { useLocale } from '@/hooks/web/useLocale'

const { changeLocale } = useLocale()

// 切换到英文
changeLocale('en-US')

// 切换到中文
changeLocale('zh-CN')
```

**`useLocale` Hook 实现**:
```typescript
import { i18n } from '@/locales'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { setHtmlPageLang } from '@/locales/helper'

const setI18nLanguage = (locale: LocaleType) => {
  const localeStore = useLocaleStoreWithOut()

  if (i18n.mode === 'legacy') {
    i18n.global.locale = locale
  } else {
    ;(i18n.global.locale as any).value = locale
  }

  localeStore.setCurrentLocale({ lang: locale })
  setHtmlPageLang(locale)
}

export const useLocale = () => {
  const changeLocale = async (locale: LocaleType) => {
    const globalI18n = i18n.global

    // 动态加载语言包
    const langModule = await import(`../../locales/lang/${locale}.ts`)

    // 设置新的语言包
    globalI18n.setLocaleMessage(locale, langModule.default)

    // 更新语言
    setI18nLanguage(locale)
  }

  return { changeLocale }
}
```

### 4.4 翻译文件格式

**示例** (`src/locales/lang/zh-CN/common.json`):
```json
{
  "add": "添加",
  "edit": "编辑",
  "delete": "删除",
  "saveSuccess": "保存成功",
  "apiRequestFailed": "请求出错，请稍候重试",
  "errMsg401": "用户没有权限（令牌、用户名、密码错误）!",
  "operationSuccess": "操作成功"
}
```

**示例** (`src/locales/lang/en-US/common.json`):
```json
{
  "add": "Add",
  "edit": "Edit",
  "delete": "Delete",
  "saveSuccess": "Save successful",
  "apiRequestFailed": "The interface request failed, please try again later!",
  "errMsg401": "The user does not have permission (token, user name, password error)!",
  "operationSuccess": "Operation Success"
}
```

## 五、高级特性

### 5.1 命名空间支持

```typescript
// 使用命名空间
const { t } = useI18n('layout')

// 相当于 t('layout.header.home')
t('header.home')
```

### 5.2 动态参数支持

```json
// common.json
{
  "welcome": "欢迎, {name}!",
  "items": "共 {count} 个项目"
}
```

```typescript
t('common.welcome', { name: 'John' })
// 输出: 欢迎, John!

t('common.items', { count: 10 })
// 输出: 共 10 个项目
```

### 5.3 Vant UI 组件库语言同步

项目自动同步 Vant UI 组件库的语言设置，确保所有 UI 组件的提示信息也是对应语言的。

```typescript
// 在 locale store 的 setCurrentLocale 中自动处理
Locale.use(localeMap?.lang, this.currentLocale.vantLocale)
```

### 5.4 本地存储持久化

用户选择的语言会保存在本地存储中，下次访问时自动恢复：

```typescript
// 读取
Storage.get('lang') || 'zh-CN'

// 保存
Storage.set('lang', locale)
```

### 5.5 HTML lang 属性自动更新

```typescript
// 自动设置 <html> 标签的 lang 属性
export const setHtmlPageLang = (locale: LocaleType) => {
  document.querySelector('html')?.setAttribute('lang', locale)
}
```

## 六、支持的翻译文件

### 6.1 通用翻译
- `common.json` - 通用提示、按钮文本、错误消息等

### 6.2 布局相关
- `layout.json` - 页面布局、菜单、导航等

### 6.3 路由相关
- `routes/basic.json` - 路由名称、标题等

### 6.4 视图相关
- `views/user/account.json` - 用户账户页面
- `views/user/model.json` - 模型设置页面

## 七、添加新语言

### 步骤 1: 创建语言配置

在 `src/locales/lang/` 下创建新的语言文件夹，例如 `fr-FR/`：

```
lang/
├── fr-FR/
│   ├── common.json
│   ├── layout.json
│   └── views/
│       └── ...
```

### 步骤 2: 创建语言入口文件

创建 `src/locales/lang/fr-FR.ts`:

```typescript
import { genMessage } from '../helper'

const modulesFiles = import.meta.glob<Recordable>('./fr-FR/**/*.json', { eager: true })

export default {
  ...genMessage(modulesFiles, 'fr-FR')
}
```

### 步骤 3: 更新配置

在 `src/locales/config.ts` 中添加新语言：

```typescript
export const localeMap = {
  'zh-CN': 'zh-CN',
  'en-US': 'en-US',
  'fr-FR': 'fr-FR'  // 新增
} as const

export const localeList = [
  // ... 其他语言
  {
    lang: localeMap['fr-FR'],
    label: 'Français',
    icon: '🇫🇷',
    title: 'Langue'
  }
] as const
```

### 步骤 4: 更新 Vant 语言映射

在 `src/stores/modules/locale.ts` 中添加：

```typescript
import frFR from 'vant/es/locale/lang/fr-FR'

const vantLocaleMap = {
  'zh-CN': zhCn,
  'en-US': enUS,
  'fr-FR': frFR  // 新增
}
```

## 八、注意事项

1. **语言文件必须同步**：所有语言文件必须保持相同的键值结构，确保每个翻译都有对应的其他语言版本

2. **键名规范**：使用小驼峰命名法，避免特殊字符

3. **模块化组织**：按功能模块划分 JSON 文件，便于维护

4. **动态导入**：使用 Vite 的 `import.meta.glob` 实现语言文件的按需加载

5. **Vant 兼容性**：确保 Vant UI 组件库支持目标语言

6. **存储持久化**：语言设置存储在本地存储中，注意清除策略

## 九、常见问题

### Q1: 如何查找缺失的翻译？

在开发时，`silentTranslationWarn` 设为 `false` 可以看到未找到的翻译警告：

```typescript
// src/locales/index.ts
return {
  silentTranslationWarn: false,  // 显示翻译警告
  // ...
}
```

### Q2: 如何处理动态加载的语言包？

使用 `changeLocale` 方法会自动异步加载新语言包，无需手动处理

### Q3: 如何在第三方库中使用翻译？

直接引入 `useI18n` hook 使用：

```typescript
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()
export const someUtil = () => t('common.message')
```

### Q4: 如何实现语言自动切换？

根据浏览器语言自动检测：

```typescript
const browserLang = navigator.language
const supportedLangs = ['zh-CN', 'en-US']
const autoLang = supportedLangs.find(lang => browserLang.startsWith(lang))
changeLocale(autoLang || 'zh-CN')
```

## 十、最佳实践

1. **统一键名**：使用有意义的英文键名，便于维护
2. **模块化管理**：按页面或功能模块拆分翻译文件
3. **参数化**：使用参数化翻译减少重复内容
4. **命名空间**：大型项目使用命名空间避免键名冲突
5. **版本控制**：翻译文件纳入版本管理，确保团队同步
6. **定期审核**：定期检查翻译质量和缺失项
