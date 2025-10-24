import { reactive, ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { getBindLibList, getOpenDoc, getOpenHome, getOpenSearch } from '@/api/open-doc'
import { getIconTemplateById } from '@/config/open-doc/icon-template'

export const useOpenDocStore = defineStore('open-doc', () => {
  const libraryInfo = reactive({
    library_avatar: '',
    library_key: '',
    library_title: '',
    preview_key: '',
    icon_template_config_id: 1,
  })

  const seoInfo = reactive({
    seo_desc: '',
    seo_keywords: '',
    seo_title: '',
    statistics_set: '',
  })

  const previewKey = ref('')
  const token = ref('')
  const libraryList = ref([])

  const sidebarOpen = ref(false)
  const catalog = ref([])
  const catalogMap = ref({})
  const allExpanded = ref(true) // 全部展开状态，默认展开
  const forceToggleTimestamp = ref(0) // 用于触发强制切换
  const isSetStatistics = ref(false)

  const isEditPage = computed(() => {
    return token.value.length
  })

  const setPreview = (data) => {
    previewKey.value = data
  }

  const setToken = (data) => {
    token.value = data
  }

  const setLibraryInfo = (data) => {
    libraryInfo.library_avatar = data.library_avatar
    libraryInfo.library_key = data.library_key
    libraryInfo.library_title = data.library_title
    libraryInfo.preview_key = data.preview_key
    libraryInfo.icon_template_config_id = data.icon_template_config_id;
  }

  const setSeoInfo = (data) => {
    seoInfo.seo_desc = data.seo_desc
    seoInfo.seo_keywords = data.seo_keywords
    seoInfo.seo_title = data.seo_title
    seoInfo.statistics_set = data.statistics_set

    document.title = seoInfo.seo_title
    document.querySelector('meta[name="description"]').setAttribute('content', seoInfo.seo_desc)
    document.querySelector('meta[name="keywords"]').setAttribute('content', seoInfo.seo_keywords)
    // statistics_set百度统计的代码插入都页面
    if (seoInfo.statistics_set.length && !isSetStatistics.value) {
      // 检查统计代码是否已存在于页面中
      isSetStatistics.value = true
      document.head.innerHTML += seoInfo.statistics_set
    }
  }

  const setCatalog = (tree = [], icon_template_config_id) => {
    let listMap = {};
    let iconTemplate = getIconTemplateById(icon_template_config_id)

    if(!tree){
      tree = []
    }
    
    const processTree = (nodes, level) => {
      return nodes.map(node => {
        node.level = level;
        // 处理doc_
        if (!node.doc_icon) {
          let iconConfig = iconTemplate.levels[level];
          node.doc_icon =  node.is_dir == 1 ? iconConfig.folder_icon : iconConfig.doc_icon;
        }

        // 如果有子节点，递归处理它们
        if (node.children && node.children.length > 0) {
          node.children = processTree(node.children, level + 1);
        }

        listMap[node.id] = node;
        return node;
      });

      
    };

    tree = processTree(tree, 0);

    catalog.value = tree;
    catalogMap.value = listMap;
  }

  async function getLibraryList(params) {
    try {
      const res = await getBindLibList({
        library_key: params.library_key || libraryInfo.library_key || '',
      })

      if (res.code == 0 && res.data) {
        libraryList.value = res.data
      }

      return libraryList.value
    } catch (error) {
      console.log(error)
      throw '获取默认知识库失败'
    }
  }

  const getDocIcon = (data) => {
    if(data.doc_icon){
      return data.doc_icon;
    }
    let doc_icon = '';
    for(let key in catalogMap.value){
      if(catalogMap.value[key].doc_key == data.doc_key){
        doc_icon = catalogMap.value[key].doc_icon;
        break;
      }
    }

    return doc_icon;
  }


  async function getDoc(id) {
    const res = await getOpenDoc({ id })

    let data = res.data;
    let seo = {
      seo_desc: data.seo_desc || '',
      seo_keywords: data.seo_keywords || data.title,
      seo_title: data.seo_title || data.title,
      statistics_set: data.statistics_set || '',
    }

    setLibraryInfo(data)
    setSeoInfo(seo)
    setCatalog(data.catalog, data.icon_template_config_id)

    data.doc_icon = getDocIcon(data);
    return data
  }

  async function getHome(id) {
    const res = await getOpenHome({ id })

    let data = res.data
    let seo = {
      seo_desc: data.seo_desc || data.library_title,
      seo_keywords: data.seo_keywords || data.library_title,
      seo_title: data.seo_title || data.library_title,
      statistics_set: data.statistics_set || '',
    }

    setLibraryInfo(data)
    setSeoInfo(seo)
    setCatalog(data.catalog, data.icon_template_config_id)

    return data
  }

  async function getSearch(id, text) {
    const res = await getOpenSearch({ id: id, v: text })

    let data = res.data

    let seo = {
      seo_desc: data.seo_desc || libraryInfo.library_title,
      seo_keywords: data.seo_keywords || text,
      seo_title: data.seo_title || text + '-' + libraryInfo.library_title,
      statistics_set: data.statistics_set || '',
    }

    setLibraryInfo(data)
    setSeoInfo(seo)
    setCatalog(data.catalog, data.icon_template_config_id)
    return data
  }

  function toggleSidebar() {
    sidebarOpen.value = !sidebarOpen.value
  }

  function toggleAllCatalog() {
    allExpanded.value = !allExpanded.value
    forceToggleTimestamp.value = Date.now() // 更新时间戳触发组件更新
  }

  return {
    getLibraryList,
    getDoc,
    getHome,
    getSearch,
    toggleSidebar,
    toggleAllCatalog,
    setPreview,
    setToken,
    getDocIcon,
    token,
    isEditPage,
    libraryList,
    previewKey,
    libraryInfo,
    seoInfo,
    sidebarOpen,
    catalog,
    catalogMap,
    allExpanded,
    forceToggleTimestamp,
  }
})
