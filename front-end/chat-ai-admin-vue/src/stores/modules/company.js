import { defineStore } from 'pinia'
import { getCompany } from '@/api/user'
const defaultAvatar = 'https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/favicon.ico'
const topNavigateDefaultData = [
  {
    id: 'robot',
    name: '机器人',
    open: true,
    isDisabled: true,
  },
  {
    id: 'library',
    name: '知识库',
    open: true,
    isDisabled: true,
  },
  {
    id: 'PublicLibrary',
    name: '文档',
    open: true
  },
  {
    id: 'library-search',
    name: '搜索',
    open: true
  },
  {
    id: 'chat-monitor',
    name: '会话',
    open: true
  },
  // {
  //   id: 'user',
  //   name: '系统设置',
  //   open: true
  // }
]

export const useCompanyStore = defineStore('company', {
  state: () => {
    return {
      name: '',
      id: '',
      companyInfo: undefined,
      ali_ocr_switch: '2',
      ali_ocr_key: '',
      ali_ocr_secret: '',
      topNavigateDefaultData: topNavigateDefaultData,
      top_navigate: topNavigateDefaultData,
    }
  },
  getters: {

  },
  actions: {
    async getCompanyInfo() {
      const res = await getCompany()

      if (!res) {
        return Promise.reject(res)
      }
      
      this.setCompanyInfo(res.data)

      return res
    },

    setCompanyInfo(data) {
      this.companyInfo = data;
      this.name = data ? data.name : '';
      this.id = data ? data.id : '';
      this.ali_ocr_switch = data ? data.ali_ocr_switch : '';
      this.ali_ocr_key = data ? data.ali_ocr_key : '';
      this.ali_ocr_secret = data ? data.ali_ocr_secret : '2';
      if(data && data.top_navigate){
        this.top_navigate = JSON.parse(data.top_navigate) || topNavigateDefaultData
      }else{
        this.top_navigate = topNavigateDefaultData
      }
      setFavicon(data?.avatar || defaultAvatar)
    },
  },
})


function setFavicon(url) {
  // 1. 移除已存在的 favicon 链接
  const existingLinks = document.querySelectorAll('link[rel="icon"], link[rel="shortcut icon"]');
  existingLinks.forEach(link => link.remove());

  // 2. 创建新的 favicon 链接
  const link = document.createElement('link');
  link.rel = 'icon';
  link.type = 'image/x-icon';
  link.href = url;

  // 3. 添加到文档头部
  document.head.appendChild(link);
}