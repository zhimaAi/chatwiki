import request from '@/utils/http/axios'

export function getVoiceList(params={}) {
  return request.get({url: '/manage/getMiniMaxVoiceList', params})
}
