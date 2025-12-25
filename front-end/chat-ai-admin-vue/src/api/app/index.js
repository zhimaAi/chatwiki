import request from '@/utils/http/axios'

// 公共的上传文件
export const uploadFile = ({ file, category }, onUploadProgress) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/upload',
    data: {
      file,
      category
    },
    onUploadProgress
  })
}

// 获取ws连接地址
export const getWsUrl = ({ openid }) => {
  return request.get({
    url: '/chat/getWsUrl',
    params: {
      openid,
      debug: import.meta.env.DEV ? 1 : 0
    }
  })
}
