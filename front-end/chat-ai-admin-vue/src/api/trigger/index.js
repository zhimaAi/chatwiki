import request from '@/utils/http/axios'

export function getTriggerConfigList(params = {}) {
  return request.get({url: '/manage/getTriggerConfigList', params})
}

export function getTriggerOfficialMessage(params = {}) {
  return request.get({url: '/manage/getTriggerOfficialMessage', params})
}