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
export const REQUEST_TIMEOUT = 60 * 60 * 1000

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

// 知识库-普通知识库头像
export const LIBRARY_NORMAL_AVATAR = '/upload/default/library_normal_avatar.svg'
// 知识库-问答知识库头像
export const LIBRARY_QA_AVATAR = '/upload/default/library_qa_avatar.svg'
// 对外文档头像
export const LIBRARY_OPEN_AVATAR = '/upload/default/library_open_avatar.svg'

// 默认的机器人头像
export const DEFAULT_ROBOT_AVATAR = '/upload/default/robot_avatar.svg'
// 默认的工作流头像
export const DEFAULT_WORKFLOW_AVATAR = '/upload/default/workflow_avatar.svg'
// 默认的MCP头像
export const DEFAULT_MCP_AVATAR = '/upload/default/mcp_avatar.svg'
// 默认模板头像
export const DEFAULT_TEMPLATE_AVATAR = '/upload/default/template_avatar.svg'
// 默认的导入csl头像
export const DEFAULT_IMPORT_CSL_AVATAR = new URL('@/assets/img/import_csl_avatar.svg', import.meta.url).href

// 默认的webApp浮动图标
export const DEFAULT_WEBAPP_ICON = new URL('@/assets/img/sdk_float_avatar.svg', import.meta.url)
  .href


// 默认prompt
export const DERAULT_ROBOT_PROMPT = `回答要求：
1、你现在是一位客服，请使用简洁、礼貌且专业的语言来回答问题
2、你只能根据知识库回答用户提问，如果你不知道答案，请回答“对不起，没有在知识库中查找到相关信息。”
3、请使用中文回答`

export const OPEN_BOC_BASE_URL = '/open-doc'

/**
 * 关键词回复：类型标签映射与筛选项
 * - REPLY_TYPE_LABEL_MAP: 组件内展示的中文标签
 * - REPLY_TYPE_OPTIONS: 下拉筛选项（value 使用前端类型标识）
 */
export const REPLY_TYPE_LABEL_MAP = {
  text: '文本',
  imageText: '图文链接',
  image: '图片',
  card: '小程序',
  url: '链接',
  smartMenu: '智能菜单'
}

export const REPLY_TYPE_OPTIONS = [
  { label: REPLY_TYPE_LABEL_MAP.imageText, value: 'imageText' },
  { label: REPLY_TYPE_LABEL_MAP.text, value: 'text' },
  { label: REPLY_TYPE_LABEL_MAP.url, value: 'url' },
  { label: REPLY_TYPE_LABEL_MAP.image, value: 'image' },
  { label: REPLY_TYPE_LABEL_MAP.card, value: 'card' },
  { label: REPLY_TYPE_LABEL_MAP.smartMenu, value: 'smartMenu' },
]

/**
 * 关注后自动回复：类型标签映射与筛选项
 * - SUBSCRIBE_REPLY_TYPE_LABEL_MAP: 组件内展示的中文标签
 * - SUBSCRIBE_REPLY_TYPE_OPTIONS: 下拉筛选项（value 使用前端类型标识）
 */
export const SUBSCRIBE_REPLY_TYPE_LABEL_MAP = {
  text: '文本',
  image: '图片',
  voice: '音频',
  video: '视频'
}

export const SUBSCRIBE_REPLY_TYPE_OPTIONS = [
  { label: SUBSCRIBE_REPLY_TYPE_LABEL_MAP.text, value: 'text' },
  { label: SUBSCRIBE_REPLY_TYPE_LABEL_MAP.image, value: 'image' },
  { label: SUBSCRIBE_REPLY_TYPE_LABEL_MAP.voice, value: 'voice' },
  { label: SUBSCRIBE_REPLY_TYPE_LABEL_MAP.video, value: 'video' }
]

// 关注来源选项（用于订阅回复-按来源设置）
export const SUBSCRIBE_SOURCE_OPTIONS = [
  { label: '公众号搜索', value: 'ADD_SCENE_SEARCH' },
  { label: '公众号迁移', value: 'ADD_SCENE_ACCOUNT_MIGRATION' },
  { label: '名片分享', value: 'ADD_SCENE_PROFILE_CARD' },
  { label: '扫描二维码', value: 'ADD_SCENE_QR_CODE' },
  { label: '图文页内名称点击', value: 'ADD_SCENE_PROFILE_LINK' },
  { label: '图文页右上角菜单', value: 'ADD_SCENE_PROFILE_ITEM' },
  { label: '支付后关注', value: 'ADD_SCENE_PAID' },
  { label: '微信广告', value: 'ADD_SCENE_WECHAT_ADVERTISEMENT' },
  { label: '他人转载', value: 'ADD_SCENE_REPRINT' },
  { label: '视频号直播', value: 'ADD_SCENE_LIVESTREAM' },
  { label: '视频号', value: 'ADD_SCENE_CHANNELS' },
  { label: '小程序关注', value: 'ADD_SCENE_WXA' },
  { label: '其他', value: 'ADD_SCENE_OTHERS' }
]
