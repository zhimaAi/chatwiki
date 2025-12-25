import request from "@/utils/http/axios/index.js";

export function getMcpSquareTypeList(params = {}) {
  return request.get({url: '/manage/getMcpSquareTypeList', params})
}

export function getMcpSquareList(params = {}) {
  return request.get({url: '/manage/getMcpSquareList', params})
}
