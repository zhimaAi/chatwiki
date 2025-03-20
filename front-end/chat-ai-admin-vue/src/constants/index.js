/**
 * 请求成功状态码
 */
export const SUCCESS_CODE = 0

/**
 * 请求contentType
 */
export const CONTENT_TYPE = 'application/x-www-form-urlencoded;charset=UTF-8'

/**
 * 请求超时时间
 */
export const REQUEST_TIMEOUT = 10 * 60 * 1000

/**
 * 刷新Token时间
 */
export const REFRESHTOKEN_TIMEOUT = 1 * 60 * 60 * 1000

/**
 * 不重定向白名单
 */
export const NO_REDIRECT_WHITE_LIST = ['/login', '/about']

/**
 * 不重置路由白名单
 */
export const NO_RESET_WHITE_LIST = ['Redirect', 'Login', 'NoFind', 'Root']

/**
 * 表格默认过滤列设置字段
 */
export const DEFAULT_FILTER_COLUMN = ['expand', 'selection']

/**
 * 是否根据headers->content-type自动转换数据格式
 */
export const TRANSFORM_REQUEST_DATA = true

// token key
export const TOKEN_KEY = 'TOKEN'

// 默认的用户头像
export const DEFAULT_USER_AVATAR = new URL('@/assets/img/user_avatar_2x.png', import.meta.url).href

// 英文logo
export const DEFAULT_EN_LOGO = new URL('@/assets/en_logo.svg', import.meta.url).href

// 中文logo
export const DEFAULT_ZH_LOGO = new URL('@/assets/zh_cn_logo.svg', import.meta.url).href

// 默认知识库头像
export const DEFAULT_LIBRARY_AVATAR = new URL(
  '@/assets/img/library/preview/default_avatar.svg',
  import.meta.url
).href

export const DEFAULT_LIBRARY_AVATAR2 = new URL('@/assets/img/library/t.svg', import.meta.url).href
export const DEFAULT_LIBRARY_AVATAR3 = new URL('@/assets/img/library/q.svg', import.meta.url).href
// 默认的机器人头像
export const DEFAULT_ROBOT_AVATAR = '/upload/default/robot_avatar.png'

// 默认的webApp浮动图标
export const DEFAULT_WEBAPP_ICON = new URL('@/assets/img/sdk_float_avatar.svg', import.meta.url)
  .href

// 默认prompt
export const DERAULT_ROBOT_PROMPT = `回答要求：
1、你现在是一位客服，请使用简洁、礼貌且专业的语言来回答问题
2、你只能根据知识库回答用户提问，如果你不知道答案，请回答“对不起，没有在知识库中查找到相关信息。”
3、请使用中文回答`
