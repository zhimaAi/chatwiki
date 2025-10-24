import { defineStore } from 'pinia'
import { store } from '../index'
import { LIBRARY_OPEN_AVATAR } from '@/constants/index'
import { getLibDocPartner, getLibraryInfo, editLibrary } from '@/api/public-library'
import { getDomainList } from '@/api/user/index'
import { getLibraryList } from '@/api/library/index'
import { getIconTemplateById } from '@/config/open-doc/icon-template'

const defaultIconTemplateConfig = getIconTemplateById(1);

export const usePublicLibraryStore = defineStore('publicLibrary', {
  state: () => {
    return {
      library_key: '',
      library_id: '',
      operate_rights: '', // 4：管理权限  2：编辑权限 0:无
      libraryInfo: {
        type: 1,
        avatar: '',
        avatar_file: '',
        library_intro: '',
        library_name: '',
        library_key: '',
        library_id: '',
        model_config_id: '',
        model_define: '',
        use_model: '',
        is_offline: false,
        share_url: '',
        access_rights: '', // 0：未公开 1：公开
        icon_template_config_id: 1,
      },
      homePreviewStyle: 'pc', // pc  mobile
      domainList: [],
      libraryList: [],
      docTreeState: [],
      docTreeStateMap: {},
      iconTemplateConfig: {
        ...defaultIconTemplateConfig
      }
    }
  },
  actions: {
    setDocTreeState(val){
      let stateMap = {}

      const loop = (arr) => {
        arr.forEach(item => {
          if(!item.doc_icon){
            item.doc_icon = this.getDocIconByLevel(item.level, item.is_dir);
          }

          stateMap[item.id] = item;

          if(item.children && item.children.length) {
            loop(item.children)
          }
         })
      }

      loop(val)
      
      this.docTreeState = val;
      this.docTreeStateMap = stateMap;
    },
    getDocIconByLevel(level, isDir){
      if(!level){
        level = 0;
      }

      let iconConfig = this.iconTemplateConfig.levels[level];

      if(isDir == 1){
        return iconConfig.folder_icon;
      }

      return iconConfig.doc_icon;
     },
    async getMyDomainList(){
      getDomainList().then((res) => {
        this.domainList = res.data || []
      })
    },
    async getLibraryList(){
      getLibraryList({type: 1}).then((res) => { 
        this.libraryList = res.data || []
      })
    },
    async getInfo(query) {
      this.library_key = query.library_key
      this.library_id = query.library_id

      await this.getLibDocPartner()
      await this.getLibraryInfo()
    },
    async getLibDocPartner() {
      try {
        let res = await getLibDocPartner({ library_key: this.library_key, library_id: this.library_id })

        this.operate_rights = res.data.operate_rights
      } catch (err){
        console.log(err)
      }
    },
    async getLibraryInfo(){
      try {
        let res = await getLibraryInfo({ id: this.library_id })
        this.libraryInfo = {
          ...this.libraryInfo,
          ...res.data
        }
        this.libraryInfo.avatar = res.data.avatar || LIBRARY_OPEN_AVATAR;
        this.libraryInfo.avatar_file = res.data.avatar_file || LIBRARY_OPEN_AVATAR;

        this.iconTemplateConfig = getIconTemplateById(res.data.icon_template_config_id);
      } catch (err){
        console.log(err)
      }
    },
    checkPermission(permissions) {
      if (!permissions || permissions.includes('*') || permissions.length === 0) {
        return true
      }
      return permissions.includes(this.operate_rights)
    },
    async saveEditLibrary(data){
      try{
        let res = await editLibrary(data)

        await this.getLibraryInfo()

        return res
      }catch (err){
        console.log(err)
      }
    },
    changeHomePreviewStyle(){
      this.homePreviewStyle = this.homePreviewStyle === 'pc' ? 'mobile' : 'pc'
    }
  }
})

export const usePublicLibraryStoreWithOut = () => {
  return usePublicLibraryStore(store)
}
