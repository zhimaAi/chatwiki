import request from "@/utils/http/axios/index.js";

export const getTemplateCates = (params = {}) => {
  return request.get({url: '/manage/getRobotTemplateCategoryList', params})
}

export const getTemplates = (params = {}) => {
  return request.get({url: '/manage/getRobotTemplateList', params})
}

export const getTplDetailMain = (params = {}) => {
  return request.get({url: '/manage/commonGetRobotTemplateDetail', params})
}

export const useRobotTemplate = (data = {}) => {
  return request.post({url: '/manage/useRobotTemplate', data})
}
