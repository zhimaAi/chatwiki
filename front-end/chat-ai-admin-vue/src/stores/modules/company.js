import { defineStore } from 'pinia'
import { getCompany } from '@/api/user'
import { useI18n } from '@/hooks/web/useI18n'

const defaultAvatar = 'https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/favicon.ico'

function getDevFixedDomains() {
  const protocol = window.location.protocol || 'http:'
  const port = window.location.port ? `:${window.location.port}` : ''

  return {
    admin_domain: `${protocol}//localhost${port}`,
    agent_domain: `${protocol}//127.0.0.1${port}`
  }
}

const topNavigateDefaultData = [
  {
    id: 'explore',
    name: '探索',
    open: true,
    isDisabled: true,
  },
  {
    id: 'robot',
    name: '机器人',
    open: true,
    isDisabled: true,
  },
  {
    id: 'clawbot',
    name: 'Clawbot',
    open: true
  },
  {
    id: 'library',
    name: '知识库',
    open: true,
    isDisabled: true,
  },
  {
    id: 'PublicLibrary-new',
    name: '文档',
    open: false
  },
  {
    id: 'library-search-new',
    name: '搜索',
    open: false
  },
  {
    id: 'chat-monitor',
    name: '会话',
    open: true
  },
  {
    id: 'workbench',
    name: '工作台',
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
      admin_domain: '',
      agent_domain: '',
      is_public_network: '',
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
      const { t } = useI18n()
      // 开发环境固定用 localhost/127.0.0.1 模拟 admin/agent 分域，避免依赖后台真实域名配置。
      const domainData = import.meta.env.DEV ? getDevFixedDomains() : {}
      const companyData = data ? { ...data, ...domainData } : data

      this.companyInfo = companyData;
      this.name = companyData ? companyData.name : '';
      this.id = companyData ? companyData.id : '';
      this.ali_ocr_switch = companyData ? companyData.ali_ocr_switch : '';
      this.ali_ocr_key = companyData ? companyData.ali_ocr_key : '';
      this.ali_ocr_secret = companyData ? companyData.ali_ocr_secret : '2';
      this.admin_domain = companyData ? companyData.admin_domain : '';
      this.agent_domain = companyData ? companyData.agent_domain : '';
      this.is_public_network = companyData ? companyData.is_public_network : '';

       // t(`navigation.${item.id}`, item.name)
       let navs = []
      if(companyData && companyData.top_navigate){
        let top_navigate = JSON.parse(companyData.top_navigate)
        let topNavigateIds = top_navigate.map(item => item.id)

        navs = top_navigate
          .map(item => {
            let defaultItem = topNavigateDefaultData.find(it => it.id == item.id)
            return defaultItem ? { ...defaultItem, ...item } : null
          })
          .filter(Boolean)
          .concat(topNavigateDefaultData.filter(item => !topNavigateIds.includes(item.id)))
      }else{
        navs = topNavigateDefaultData
      }

      navs.forEach(item => {
        item.name = t(`navigation.${item.id}`, item.name)
      })

      this.top_navigate = navs

      setFavicon(companyData?.avatar || defaultAvatar)
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
