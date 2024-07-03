import request from '@/utils/http/axios'

export const getRobotList = () => {
  return request.get({
    url: '/manage/clientSide/getRobotList'
  })
}
