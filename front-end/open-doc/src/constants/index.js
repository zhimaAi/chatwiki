// 静态导入资源
import userAvatarUrl from '@/assets/user_avatar_2x.png'

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
export const DEFAULT_USER_AVATAR = userAvatarUrl

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

// 默认prompt
export const DERAULT_ROBOT_PROMPT = `回答要求：
1、你现在是一位客服，请使用简洁、礼貌且专业的语言来回答问题
2、你只能根据知识库回答用户提问，如果你不知道答案，请回答“对不起，没有在知识库中查找到相关信息。”
3、请使用中文回答`
