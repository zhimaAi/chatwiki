import request from '@/utils/http/axios'

export function getTriggerConfigList(params = {}) {
  return request.get({url: '/manage/getTriggerConfigList', params})
}