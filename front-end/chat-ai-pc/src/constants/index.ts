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
export const REQUEST_TIMEOUT = 5 * 60 * 1000

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
// SDK AVAVATAR
export const DEFAULT_SDK_FLOAT_AVATAR = new URL('@/assets/img/sdk_float_avatar.svg', import.meta.url).href
// 默认的头部图像 
export const DEFAULT_HEAD_IMAGE = new URL('@/assets/img/web_app_default_avatar.png', import.meta.url).href
