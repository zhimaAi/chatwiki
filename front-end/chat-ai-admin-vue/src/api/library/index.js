import SSE from '@/utils/http/sse'
import request from '@/utils/http/axios'
import { useUserStore } from '@/stores/modules/user'

const baseURL = import.meta.env.VITE_BASE_API_URL

export const saveLibrarySearch = (data) => {
  return request.post({
    url: '/manage/saveLibrarySearch',
    data: data
  })
}

export const getLibrarySearch = (params = {}) => {
  return request.get({
    url: '/manage/getLibrarySearch',
    params: params
  })
}

export const searchLirary = (data) => {
  const userStore = useUserStore()
  
  return new SSE({
    token: userStore.getToken ?? '',
    url: baseURL + '/manage/libraryAiSummary',
    data: data
  })
}

export const getLibraryList = (params = {}) => {
  return request.get({
    url: '/manage/getLibraryList',
    params: params
  })
}

export const generateSimilarQuestions = (data) => {
  return request.post({
    url: '/manage/generateSimilarQuestions',
    data
  })
}

// AI生成提示词_ai分段
export const generateAiPrompt = (data) => {
  return request.post({
    url: '/manage/generateAiPrompt',
    data
  })
}

export const deleteLibrary = ({ id }) => {
  return request.post({
    url: '/manage/deleteLibrary',
    data: {
      id
    }
  })
}

export const createLibrary = (data) => {
  return request.post({
    headers: {
      //'Content-Type': 'multipart/form-data'
    },
    url: '/manage/createLibrary',
    data: data
  })
}

export const getLibraryFileList = ({ status, library_id, file_name = undefined, page = 1, size = 20, group_id,sort_field,sort_type }) => {
  return request.get({
    url: '/manage/getLibFileList',
    params: {
      status,
      library_id,
      file_name,
      page,
      size,
      group_id,
      sort_field,
      sort_type
    }
  })
}

export const getLibraryInfo = ({ id }) => {
  return request.get({
    url: '/manage/getLibraryInfo',
    params: {
      id
    }
  })
}

export const delLibraryFile = ({ id }) => {
  return request.post({
    url: '/manage/delLibraryFile',
    data: {
      id
    }
  })
}

export const addLibraryFile = (data) => {
  return request.post({
    url: '/manage/addLibraryFile',
    data: data
  })
}

export const readLibFileExcelTitle = (data) => {
  return request.post({
    url: '/manage/readLibFileExcelTitle',
    data: data
  })
}

export const editLibrary = (data) => {
  return request.post({
    url: '/manage/editLibrary',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    data: data
  })
}

// 获取分隔符列表
export const getSeparatorsList = () => {
  return request.get({
    url: '/manage/getSeparatorsList',
    params: {}
  })
}

// 获取文档拆分
export const getLibFileSplit = (params) => {
  return request.get({
    url: '/manage/getLibFileSplit',
    params: params
  })
}

// 获取文档拆分 //新接口
export const getLibFileSplitPreview = (params) => {
  return request.get({
    url: '/manage/getLibFileSplitPreview',
    params: params
  })
}

// 段落重新分段
export const getSplitParagraph = (params) => {
  return request.get({
    url: '/manage/getSplitParagraph',
    params: params
  })
}

export const getLibFileInfo = (params = {}) => {
  return request.get({
    url: '/manage/getLibFileInfo',
    params: params
  })
}

export const saveLibFileSplit = (data) => {
  return request.post({
    url: '/manage/saveLibFileSplit',
    data: data
  })
}

// 保存段落重新分段
export const saveSplitParagraph = (data) => {
  return request.post({
    url: '/manage/saveSplitParagraph',
    data: data
  })
}

export const getParagraphList = (params = {}) => {
  return request.get({
    url: '/manage/getParagraphList',
    params: params
  })
}

export const editParagraph = (data) => {
  return request.post({
    url: '/manage/editParagraph',
    data: data
  })
}

export const deleteParagraph = (data) => {
  return request.post({
    url: '/manage/deleteParagraph',
    data: data
  })
}

export const libraryRecallTest = (data) => {
  return request.post({
    url: '/manage/libraryRecallTest',
    data: data
  })
}

export const getLibFileExcelTitle = (params) => {
  return request.get({
    url: '/manage/getLibFileExcelTitle',
    params: params
  })
}

export const editLibFile = (data) => {
  return request.post({
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    url: '/manage/editLibFile',
    data: data
  })
}


export const getLibraryRobotInfo = (params = {}) => {
  return request.get({
    url: '/manage/getLibraryRobotInfo',
    params: params
  })
}

export const createGraph = (data) => {
  return request.post({
    url: '/manage/constructGraph',
    data: data
  })
}

export const reconstructVector = (data) => {
  return request.post({
    url: '/manage/reconstructVector',
    data: data
  })
}

export const reconstructCategoryVector = (data) => {
  return request.post({
    url: '/manage/reconstructCategoryVector',
    data: data
  })
}

export const reconstructGraph = (data) => {
  return request.post({
    url: '/manage/reconstructGraph',
    data: data
  })
}

// 获取ai分段
export const getLibFileSplitAiChunks = (params) => {
  return request.get({
    url: '/manage/getLibFileSplitAiChunks',
    params: params
  })
}

export const getParagraphCount = (params = {}) => {
  return request.get({
    url: '/manage/getParagraphCount',
    params: params
  })
}

export const getCategoryList = (params = {}) => {
  return request.get({
    url: '/manage/getCategoryList',
    params: params
  })
}

export const saveCategory = (data) => {
  return request.post({
    url: '/manage/saveCategory',
    data: data
  })
}

export const updateParagraphCategory = (data) => {
  return request.post({
    url: '/manage/updateParagraphCategory',
    data: data
  })
}

export const restudyLibraryFile = (data) => {
  return request.post({
    url: '/manage/restudyLibraryFile',
    data: data
  })
}

export const getCategoryParagraphList = (params = {}) => {
  return request.get({
    url: '/manage/getCategoryParagraphList',
    params: params
  })
}

export const saveCategoryParagraph = (data) => {
  return request.post({
    url: '/manage/saveCategoryParagraph',
    data: data
  })
}

export const manualCrawl = (data) => {
  return request.post({
    url: '/manage/manualCrawl',
    data: data
  })
}


// 知识库关联机器人
export const relationRobot = (data) => {
  return request.post({
    url: '/manage/relationRobot',
    data: data
  })
}


// pdf取消解析
export const cancelOcrPdf = (data) => {
  return request.post({
    url: '/manage/cancelOcrPdf',
    data: data
  })
}

export const getLibRawFileOnePage = (params = {}) => {
  return request.get({
    url: '/manage/getLibRawFileOnePage',
    params: params
  })
}

// 获取知识图谱数据
export const getFileGraphInfo = (params) => {
  return request.get({
    url: '/manage/getFileGraphInfo',
    params: params
  })
}

export const getLibraryGroup = (params) => {
  return request.get({
    url: '/manage/getLibraryGroup',
    params: params
  })
}

export const sortLibararyListGroup = (data) => {
  return request.post({
    url: '/manage/sortLibararyListGroup',
    data: data
  })
}

export const saveLibraryGroup = (data) => {
  return request.post({
    url: '/manage/saveLibraryGroup',
    data: data
  })
}

export const relationLibraryGroup = (data) => {
  return request.post({
    url: '/manage/relationLibraryGroup',
    data: data
  })
}

export const deleteLibraryGroup = (data) => {
  return request.post({
    url: '/manage/deleteLibraryGroup',
    data: data
  })
}

export const setParagraphGroup = (data) => {
  return request.post({
    url: '/manage/setParagraphGroup',
    data: data
  })
}

export const getFAQFileList = (params) => {
  return request.get({
    url: '/manage/getFAQFileList',
    params: params
  })
}

export const addFAQFile = (data) => {
  return request.post({
    url: '/manage/addFAQFile',
    data: data
  })
}

export const importParagraph = (data) => {
  return request.post({
    url: '/manage/importParagraph',
    data: data
  })
}

export const deleteFAQFile = (data) => {
  return request.post({
    url: '/manage/deleteFAQFile',
    data: data
  })
}

export const renewFAQFileData = (data) => {
  return request.post({
    url: '/manage/renewFAQFileData',
    data: data
  })
}

export const getFAQFileQAList = (params) => {
  return request.get({
    url: '/manage/getFAQFileQAList',
    params: params
  })
}

export const deleteFAQFileQA = (data) => {
  return request.post({
    url: '/manage/deleteFAQFileQA',
    data: data
  })
}

export const saveFAQFileQA = (data) => {
  return request.post({
    url: '/manage/saveFAQFileQA',
    data: data
  })
}

export const getFAQConfig = (params) => {
  return request.get({
    url: '/manage/getFAQConfig',
    params: params
  })
}

export const getFAQFileInfo = (params) => {
  return request.get({
    url: '/manage/getFAQFileInfo',
    params: params
  })
}

export const getFAQFileChunks = (params) => {
  return request.get({
    url: '/manage/getFAQFileChunks',
    params: params
  })
}

export const createExportLibFileTask = (params) => {
  return request.get({
    url: '/manage/createExportLibFileTask',
    params: params
  })
}

export const getLibraryListGroup = (params) => {
  return request.get({
    url: '/manage/getLibraryListGroup',
    params: params
  })
}


export const saveLibraryListGroup = (data) => {
  return request.post({
    url: '/manage/saveLibraryListGroup',
    data: data
  })
}

export const deleteLibraryListGroup = (data) => {
  return request.post({
    url: '/manage/deleteLibraryListGroup',
    data: data
  })
}


export const getLibFileRecycleList = (params) => {
  return request.get({
    url: '/manage/getLibFileRecycleList',
    params: params
  })
}

export const delRecycleLibraryFile = (data) => {
  return request.post({
    url: '/manage/delRecycleLibraryFile',
    data: data
  })
}


export const restoreRecycleLibraryFile = (data) => {
  return request.post({
    url: '/manage/restoreRecycleLibraryFile',
    data: data
  })
}

export const statLibraryTotal = (data) => {
  return request.post({
    url: '/manage/statLibraryTotal',
    data: data
  })
}

export const statLibraryDataSort = (data) => {
  return request.post({
    url: '/manage/statLibraryDataSort',
    data: data
  })
}

export const statLibrarySort = (data) => {
  return request.post({
    url: '/manage/statLibrarySort',
    data: data
  })
}

export const statLibraryDataRobotDetail = (data) => {
  return request.post({
    url: '/manage/statLibraryDataRobotDetail',
    data: data
  })
}

export const statLibraryRobotDetail = (data) => {
  return request.post({
    url: '/manage/statLibraryRobotDetail',
    data: data
  })
}