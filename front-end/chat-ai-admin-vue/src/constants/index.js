import { useI18n } from '@/hooks/web/useI18n'
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
export const NO_REDIRECT_WHITE_LIST = ['/login', '/about', '/privacy_policy']

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
export const LIBRARY_NORMAL_AVATAR = new URL('@/assets/svg/ordinary-icon.svg', import.meta.url).href
// 知识库-问答知识库头像
export const LIBRARY_QA_AVATAR = new URL('@/assets/svg/faq-icon.svg', import.meta.url).href
// 知识库-公众号知识库头像
export const LIBRARY_OFFICIAL_ACCOUNT_AVATAR = new URL(
  '@/assets/svg/accounts-icon.svg',
  import.meta.url
).href
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

// 默认模板主图
export const DEFAULT_TEMPLATE_MAIN_PIC =
  'https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/modal-main-pic-new.png'
// 默认的导入csl头像
export const DEFAULT_IMPORT_CSL_AVATAR = new URL(
  '@/assets/img/import_csl_avatar.svg',
  import.meta.url
).href

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
export const REPLY_TYPE_LABEL_MAP = () => {
  const { t } = useI18n('constants.index')
  return {
    text: t('text'),
    imageText: t('imageText'),
    image: t('image'),
    card: t('card'),
    url: t('url'),
    smartMenu: t('smartMenu')
  }
}

export const REPLY_TYPE_OPTIONS = () => {
  const { t } = useI18n('constants.index')
  return [
    { label: t('imageText'), value: 'imageText' },
    { label: t('text'), value: 'text' },
    { label: t('url'), value: 'url' },
    { label: t('image'), value: 'image' },
    { label: t('card'), value: 'card' },
    { label: t('smartMenu'), value: 'smartMenu' }
  ]
}

/**
 * 关注后自动回复：类型标签映射与筛选项
 * - SUBSCRIBE_REPLY_TYPE_LABEL_MAP: 组件内展示的中文标签
 * - SUBSCRIBE_REPLY_TYPE_OPTIONS: 下拉筛选项（value 使用前端类型标识）
 */
export const SUBSCRIBE_REPLY_TYPE_LABEL_MAP = () => {
  const { t } = useI18n('constants.index')
  return {
    text: t('text'),
    image: t('image'),
    voice: t('voice'),
    video: t('video')
  }
}

export const SUBSCRIBE_REPLY_TYPE_OPTIONS = () => {
  const { t } = useI18n('constants.index')
  return [
    { label: t('text'), value: 'text' },
    { label: t('image'), value: 'image' },
    { label: t('voice'), value: 'voice' },
    { label: t('video'), value: 'video' }
  ]
}


// 关注来源选项（用于订阅回复-按来源设置）
export const SUBSCRIBE_SOURCE_OPTIONS = () => {
  const { t } = useI18n('constants.index')
  return [
    { label: t('subscribe_search'), value: 'ADD_SCENE_SEARCH' },
    { label: t('subscribe_migration'), value: 'ADD_SCENE_ACCOUNT_MIGRATION' },
    { label: t('profile_card'), value: 'ADD_SCENE_PROFILE_CARD' },
    { label: t('qr_code'), value: 'ADD_SCENE_QR_CODE' },
    { label: t('profile_link'), value: 'ADD_SCENE_PROFILE_LINK' },
    { label: t('profile_item'), value: 'ADD_SCENE_PROFILE_ITEM' },
    { label: t('paid'), value: 'ADD_SCENE_PAID' },
    { label: t('wechat_ad'), value: 'ADD_SCENE_WECHAT_ADVERTISEMENT' },
    { label: t('reprint'), value: 'ADD_SCENE_REPRINT' },
    { label: t('livestream'), value: 'ADD_SCENE_LIVESTREAM' },
    { label: t('channels'), value: 'ADD_SCENE_CHANNELS' },
    { label: t('wxa'), value: 'ADD_SCENE_WXA' },
    { label: t('other'), value: 'ADD_SCENE_OTHERS' }
  ]
}
