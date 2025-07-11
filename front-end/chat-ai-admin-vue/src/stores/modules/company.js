import { defineStore } from 'pinia'
import { getCompany } from '@/api/user'
const defaultAvatar = 'https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/favicon.ico'
export const useCompanyStore = defineStore('company', {
  state: () => {
    return {
      name: '',
      id: '',
      companyInfo: undefined,
      ali_ocr_switch: '2',
      ali_ocr_key: '',
      ali_ocr_secret: ''
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