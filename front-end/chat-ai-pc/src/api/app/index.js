import request from '@/utils/http/axios'

// 公共的上传文件
export const uploadFile = ({ file, category, extraData }, onUploadProgress) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/upload',
    data: {
      file,
      category,
      ...extraData
    },
    onUploadProgress
  })
}
