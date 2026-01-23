<style lang="less" scoped>
.home-preview {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;

  .toolbar-wrapper {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
    padding: 0 16px 0 32px;
    background-color: #fff;
    border-bottom: 1px solid #D9D9D9;

    .toolbar-right {
      display: flex;
      align-items: center;
    }

    .action-item {
      margin-right: 8px;

      &:last-child {
        margin-right: 0;
      }

      .action-name{
        margin-left: 8px;
      }
    }

    .change-home-preview-style-btn {
      margin-left: 8px;
      cursor: pointer;
    }

    .edit-user {
      display: flex;
      align-items: center;
      margin-right: 16px;

      .line {
        width: 1px;
        height: 16px;
        margin-right: 16px;
        border-radius: 1px;
        background: #d9d9d9;
      }

      .user-avatar {
        width: 24px;
        height: 24px;
        border-radius: 6px;
      }

      .user-name {
        margin-left: 8px;
        font-size: 14px;
        color: #595959;
      }

      .edit-status-text {
        margin-left: 8px;
        font-size: 14px;
        color: #8c8c8c;
      }
    }
  }

  .library-info {
    display: flex;
    flex-direction: row;
    line-height: 22px;
    padding: 24px 24px 24px 0;

    .library-logo {
      width: 32px;
      height: 32px;
      margin-right: 8px;
      border-radius: 12px;
    }

    .library-info-content {
      flex: 1;
      display: flex;
      align-items: center;
      flex-direction: row;
    }

    .library-name {
      flex: 1;
      line-height: 24px;
      height: 24px;
      font-size: 16px;
      font-weight: 600;
      color: #000000;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }

    .libray-arrow-down {
      margin-left: 12px;
      font-size: 16px;
      cursor: pointer;
    }

    .publish-status {
      display: flex;
      align-items: center;
      height: 22px;
      line-height: 22px;
      padding: 0 6px;
      border-radius: 6px;
      margin-left: 12px;
      color: #3a4559;

      &.is-publish {
        color: #21A665;
        background-color: #CAFCE4;
      }

      .status-icon {
        margin-right: 2px;
        font-size: 14px;
      }

      .status-name {
        font-size: 14px;
        font-weight: 500;
      }
    }
  }

  .preview-box {
    position: relative;
    flex: 1;
    padding: 0;
    overflow: hidden;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;

    &.is-dragging::before {
      content: '';
      display: block;
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0);
      z-index: 10;
    }

    &.is-mobile {
      align-items: flex-start;
      overflow-y: auto;
    }

    .mobile-box {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 800px;
      width: 100%;
      padding: 24px 0;
      padding-bottom: 41px;
      overflow: hidden;
    }

    .pc-box {
      height: 100%;
      width: 100%;
    }

    .iframe {
      height: 100%;
      width: 100%;
    }
  }

  .seo-setting-box {
    padding: 0 24px 32px 24px;
  }
}
</style>

<template>
  <div class="home-preview">
    <div class="toolbar-wrapper">
      <div class="toolbar-left">
        <div class="library-info">
          <div class="library-info-content">
            <wiki-dropdown :libraryInfo="libraryInfo" :libraryList="libraryList"
              @change="handleChangeLibrary"></wiki-dropdown>

            <div class="publish-status" v-if="libraryInfo.access_rights == 0">
              <ExclamationCircleFilled class="status-icon" />
              <span class="status-name">{{ t('not_published') }}</span>
            </div>
            <div class="publish-status  is-publish" v-if="libraryInfo.access_rights == 1">
              <CheckCircleFilled class="status-icon" />
              <span class="status-name">{{ t('published') }}</span>
            </div>

            <span class="change-home-preview-style-btn">
              <a-tooltip :title="homePreviewTip">
                <svg-icon name="switch-mobile" style="font-size: 14px;color:#595959;"
                  @click.stop="changeHomePreviewStyle(item)"></svg-icon>
              </a-tooltip>
            </span>
          </div>
        </div>
      </div>
      <div class="toolbar-right">
        <div class="edit-user" v-if="editUser.show && editUser.user_name">
          <img class="user-avatar" :src="editUser.avatar" alt="" />
          <span class="user-name">{{ editUser.nick_name || editUser.user_name }}</span>
          <span class="edit-status-text">{{ t('editing') }}</span>
        </div>

        <div class="action-item">
          <a-button type="primary" :loading="publishLoading" @click="handlePublish" v-if="isEdit">
            <CheckOutlined /> <span class="action-name">{{ t('publish') }}</span>
          </a-button>
          <a-button type="primary" @click="handleHomeEdit" v-else>
            <svg-icon class="action-icon" name="edit" style="font-size: 16px; color: #fff"></svg-icon>
            <span class="action-name">{{ t('edit') }}</span>
          </a-button>
        </div>

        <div class="action-item">
          <a-button @click="handleCopyShareUrl">
            <svg-icon class="action-icon" name="link-left" style="font-size: 16px"></svg-icon>
            <span class="action-name">{{ t('copy_link') }}</span>
          </a-button>
        </div>

        <div class="action-item">
          <a-button @click="handleOpenSeoSetting">
            <CheckCircleFilled style="color: #21a665; font-size: 14px" class="status-icon" v-if="isSetSeo" />
            <ExclamationCircleFilled  style="color: #ccc; font-size: 14px" class="status-icon" v-else />
            <span class="action-name">{{ t('seo') }}</span>
          </a-button>
        </div>

        <div class="action-item">
          <a-button @click="toSettingPage">
            <svg-icon name="jibenpeizhi" style="font-size: 14px;color: #595959;"></svg-icon>
            <span class="action-name">{{ t('settings') }}</span>
          </a-button>
        </div>
      </div>
    </div>

    <div class="preview-box" :class="{ 'is-mobile': styleType == 'mobile', 'is-dragging': isDragging }"
      v-if="state.doc_id">
      <div class="mobile-box" v-if="styleType == 'mobile'">
        <PhoneBox>
          <iframe ref="iframeRef" class="iframe" :src="previewUrl" frameborder="0"
            v-if="libraryInfo.library_key"></iframe>
        </PhoneBox>
      </div>

      <div class="pc-box" v-else>
        <iframe ref="iframeRef" class="iframe" :src="previewUrl" frameborder="0"
          v-if="libraryInfo.library_key"></iframe>
      </div>
    </div>
    <EditTitle ref="editTitleRef" @ok="handleSaveTitle" :confirm-loading="saveTitleLoading" />
    <EditDesc ref="editDescRef" @ok="handleSaveDesc" :confirm-loading="saveDescLoading" />
    <EditQuestionGuide ref="editQuestionGuideRef" @ok="handleSaveQuestionGuide"
      :confirm-loading="saveQuestionGuideLoading" />
    <SeoSetting ref="seoSettingRef" @ok="saveSeoOk" />
    <ShortcutSelector ref="shortcutSelectorRef" :iconTemplateConfig="iconTemplateConfig" @ok="onSelectAddShortcut" />
  </div>
</template>

<script setup>
import { generateUniqueId } from '@/utils/index'
import { useCopyShareUrl } from '@/hooks/web/useCopyShareUrl'
import { useI18n } from '@/hooks/web/useI18n'
import { OPEN_BOC_BASE_URL } from '@/constants/index'
import {
  saveDraftLibDoc,
  saveQuestionGuide,
  deleteQuestionGuide,
  saveLibDoc,
  getLibDocInfo,
  getLibDocHomeConfig,
  saveLibDocIndexQuickDoc,
  saveLibDocBannerImg
} from '@/api/public-library'
import { getUser } from '@/api/manage'
import { uploadFile } from '@/api/app'
import { useUserStore } from '@/stores/modules/user'
import { ref, computed, reactive, nextTick, onMounted, onBeforeUnmount, inject } from 'vue'
import { CheckOutlined, ExclamationCircleFilled, CheckCircleFilled } from '@ant-design/icons-vue'
import { useRouter } from 'vue-router'
import { usePublicLibraryStore } from '@/stores/modules/public-library'
import { message, Modal } from 'ant-design-vue'
import EditTitle from './components/edit-title.vue'
import EditDesc from './components/edit-desc.vue'
import EditQuestionGuide from './components/edit-question-guide.vue'
import PhoneBox from './components/phone-box.vue'
import SeoSetting from '../components/seo-setting.vue'
import WikiDropdown from '../components/wiki-dropdown.vue'
import ShortcutSelector from './components/shortcut-selector.vue'

const { t } = useI18n('views.public-library.home.index')

const router = useRouter()
// 接受isDragging
const isDragging = inject('isDragging')

const libraryStore = usePublicLibraryStore()
const { getToken, user_id } = useUserStore()

const libraryInfo = computed(() => libraryStore.libraryInfo)
const libraryId = computed(() => libraryStore.library_id)
const libraryKey = computed(() => libraryStore.library_key)
const iconTemplateConfig = computed(() => libraryStore.iconTemplateConfig)

const homePreviewStyle = computed(() => {
  return libraryStore.homePreviewStyle
})

const homePreviewTip = computed(() => {
  return homePreviewStyle.value == 'mobile' ? t('switch_to_pc') : t('switch_to_mobile')
})

const changeHomePreviewStyle = () => {
  libraryStore.changeHomePreviewStyle()
}

const iframeRef = ref(null)
const time = ref(new Date().getTime())
const isEdit = ref(false)
const state = reactive({
  doc_id: '',
  title: '',
  content: '',
  library_intro: '',
  library_name: '',
  library_key: '',
  is_pub: '',
  doc_key: '',
  seo_title: '',
  seo_desc: '',
  seo_keywords: ''
})

const isSetSeo = computed(() => { 
  return state.seo_title || state.seo_desc || state.seo_keywords
})

// 打开seo设置
const seoSettingRef = ref(null);
const handleOpenSeoSetting = () => {
  seoSettingRef.value.open({ ...state })
}

const saveSeoOk = (data) => {
  state.seo_title = data.seo_title
  state.seo_desc = data.seo_desc
  state.seo_keywords = data.seo_keywords
  // time.value = new Date().getTime()
}

const editUser = reactive({
  show: false,
  user_name: '',
  avatar: '',
  nick_name: ''
})

const getEditUser = async (userId) => {
  let res = await getUser({ id: userId })
  editUser.avatar = res.data.avatar
  editUser.user_name = res.data.user_name
  editUser.nick_name = res.data.nick_name
  return res
}
// 编辑标题
const editTitleRef = ref(null)
const saveTitleLoading = ref(false)

const handleSaveTitle = (data) => {
  state.title = data.title
  saveTitleLoading.value = true
  handleSaveDraft()
    .then(() => {
      saveTitleLoading.value = false
      editTitleRef.value.close()
    })
    .catch(() => {
      saveTitleLoading.value = false
    })
}

// 编辑描述
const editDescRef = ref(null)
const saveDescLoading = ref(false)
const handleSaveDesc = (data) => {
  state.content = data.content
  saveDescLoading.value = true
  handleSaveDraft()
    .then(() => {
      saveDescLoading.value = false
      editDescRef.value.close()
    })
    .catch(() => {
      saveDescLoading.value = false
    })
}

// 添加快捷方式
const shortcutSelectorRef = ref(null);
const shortcutList = ref([]);
const shortcutListValue = ref([]);

const handleOpenShortcutSelector = (data) => {
  console.log(data)
  shortcutSelectorRef.value.open({ library_key: libraryInfo.value.library_key, ...data })
}
const onSelectAddShortcut = (data) => {
  let key = generateUniqueId('shortcutList');
  console.log(data)
  if (!data.old_doc_id) {
    shortcutList.value.push({
      doc_id: data.doc_id,
      key
    })
  } else {
    let index = shortcutList.value.findIndex(item => item.key == data.old_doc_key)
    // 如果key不存在则用id在找一次
    if(index == -1){
      index = shortcutList.value.findIndex(item => item.doc_id == data.old_doc_id)
    }
    
    if (index > -1) {
      shortcutList.value.splice(index, 1, { doc_id: data.doc_id, key})
    } else {
      shortcutList.value.push({ doc_id: data.doc_id, key })
    }
  }

  saveLibDocIndexQuickDoc({
    library_key: libraryInfo.value.library_key,
    doc_id: state.doc_id,
    quick_doc_content: JSON.stringify(shortcutList.value.filter(item => item.doc_id))
  }).then(() => {
    time.value = new Date().getTime()
    message.success(t('save_success'))
  })
}

const handleDeleteShorcut = (data) => {
  let index = shortcutList.value.findIndex(item => item.key == data.key)

  // 如果key不存在则用id在找一次
  if(index == -1){
    index = shortcutList.value.findIndex(item => item.doc_id == data.doc_id)
  }

  if (index > -1) {
    shortcutList.value.splice(index, 1)
  }

  saveLibDocIndexQuickDoc({
    library_key: libraryInfo.value.library_key,
    doc_id: state.doc_id,
    quick_doc_content: JSON.stringify(shortcutList.value.filter(item => item.doc_id))
  }).then(() => {
    time.value = new Date().getTime()
    message.success(t('delete_success'))
  })
}

const handleShortcutSort = (data) => { 
  shortcutList.value = data.list;

  saveLibDocIndexQuickDoc({
    library_key: libraryInfo.value.library_key,
    doc_id: state.doc_id,
    quick_doc_content: JSON.stringify(shortcutList.value.filter(item => item.doc_id))
  }).then(() => {
    // time.value = new Date().getTime()
    message.success(t('save_success'))
  })
 }

// 保存
const saveDraftDoc = async () => {
  let res = await saveDraftLibDoc({
    library_key: libraryInfo.value.library_key,
    doc_id: state.doc_id,
    title: state.title,
    content: state.content,
    is_index: 1
  }).then((res) => {
    state.doc_id = res.data.doc_id

    time.value = new Date().getTime()
    return res
  })

  return res
}
const handleSaveDraft = async () => {
  let res = await saveDraftDoc()

  if (res) {
    message.success(t('save_success'))
  }
}
// 发布
const publishLoading = ref(false)
const handlePublish = () => {
  let data = {
    library_key: libraryInfo.value.library_key,
    doc_id: state.doc_id,
    title: state.title,
    content: state.content,
    is_index: 1
  }

  publishLoading.value = true

  saveLibDoc(data)
    .then((res) => {
      state.doc_id = res.data.doc_id
      publishLoading.value = false
      message.success(t('publish_success'))
      setHomeEditStatus(false)

      time.value = new Date().getTime()
    })
    .catch(() => {
      publishLoading.value = false
    })
}

// 保存问题引导
const editQuestionGuideRef = ref(null)
const saveQuestionGuideLoading = ref(false)
const questionGuide = reactive({
  id: '',
  question: ''
})

const handleSaveQuestionGuide = (data) => {
  questionGuide.question = data.question
  saveQuestionGuideLoading.value = true
  saveQuestionGuide({
    library_key: libraryInfo.value.library_key,
    id: questionGuide.id || '',
    question: questionGuide.question
  })
    .then(() => {
      time.value = new Date().getTime()
      saveQuestionGuideLoading.value = false
      editQuestionGuideRef.value.close()
    })
    .catch(() => {
      saveQuestionGuideLoading.value = false
    })
}

const styleType = computed(() => {
  return libraryStore.homePreviewStyle
})

let prviewHost = import.meta.env.VITE_OPEN_DOC_HOST || ''

const previewUrl = computed(() => {
  let url = ''
  if (import.meta.env.DEV) {
    url = `${prviewHost}/home/${libraryInfo.value.library_key}?&v=${time.value}&token=${getToken}`
  } else {
    url = `${libraryInfo.value.share_url || ''}${OPEN_BOC_BASE_URL}/home/${libraryInfo.value.library_key}?&v=${time.value}&token=${getToken}`
  }

  return url
})

const getHomeInfo = async () => {
  await getLibDocHomeConfig({ library_key: libraryInfo.value.library_key }).then(async (res) => {
    // 如果没有doc_id 则创建一个新的doc_id
    if (!res.data.id) {
      state.title = libraryInfo.value.library_name
      state.content = libraryInfo.value.library_intro

      await saveDraftDoc()

      getHomeInfo()
      return
    }

    state.library_key = libraryInfo.value.library_key;
    state.library_intro = libraryInfo.value.library_intro;
    state.library_name = libraryInfo.value.library_name;
    state.doc_id = res.data.id;
    state.title = res.data.title;
    state.content = res.data.content;
    state.seo_title = res.data.seo_title;
    state.seo_desc = res.data.seo_desc;
    state.seo_keywords = res.data.seo_keywords;
    state.is_pub = res.data.is_pub;
    state.doc_key = res.data.doc_key;

    if (res.data.quick_doc_content && res.data.quick_doc_content.length) {
      let shortcuts= JSON.parse(res.data.quick_doc_content) || [];

      shortcutList.value = shortcuts.filter(item => item.doc_id)
    }else{
      shortcutList.value = [];
    }

    if (res.data.quick_doc_content_value && res.data.quick_doc_content_value.length) {
      shortcutListValue.value = JSON.parse(res.data.quick_doc_content_value) || [];
    }else{
      shortcutListValue.value = []
    }

    let editStatus = res.data.edit_user && res.data.edit_user != 0;

    if (editStatus && res.data.edit_user != user_id) {
      editUser.show = true
      getEditUser(res.data.edit_user)
    }

    nextTick(() => {
      if (editStatus && res.data.edit_user == user_id) {
        setHomeEditStatus(true)
      } else {
        setHomeEditStatus(false)
      }
    })
  })
}

// 知识库列表
const libraryList = computed(() => libraryStore.libraryList)

const fetchLibraryList = () => {
  libraryStore.getLibraryList()
}

const handleHomeEdit = async () => {
  try {
    let doc = await getLibDocInfo({
      doc_id: state.doc_id,
      library_key: libraryInfo.value.library_key
    })

    let editStatus = doc.data.edit_user && doc.data.edit_user != 0

    if (editStatus && doc.data.edit_user != user_id) {
      let user = await getEditUser(doc.data.edit_user)
      let name = user.data.nick_name || user.data.user_name
      editUser.show = true
      Modal.warning({
        title: t('editing_warning_title', { name }),
        content: t('editing_warning_content'),
        okText: t('editing_warning_ok')
      })

      return
    }

    saveDraftDoc()

    // setHomeEditStatus(true)
    isEdit.value = true
    editUser.show = false
  } catch (error) {
    return
  }
}

const setHomeEditStatus = (type) => {
  if (!iframeRef.value) {
    return
  }
  isEdit.value = type

  const iframeWindow = iframeRef.value.contentWindow

  iframeWindow.postMessage({ action: 'setEditStatus', data: { type: isEdit.value } }, '*')
}

const { copyShareUrl } = useCopyShareUrl()

const handleCopyShareUrl = async () => {
  const docUrl = OPEN_BOC_BASE_URL + '/home/' + libraryInfo.value.library_key

  await copyShareUrl(docUrl)
}

const initHomeDoc = () => {

}

const handleEditTitle = (data) => {
  state.title = data.title
  state.content = data.content
  editTitleRef.value.open(data)
}

const handleEditDesc = (data) => {
  state.title = data.title
  state.content = data.content
  editDescRef.value.open(data)
}

const handleEditQuestionGuide = (data) => {
  questionGuide.id = data.id
  questionGuide.question = data.question
  editQuestionGuideRef.value.open(data)
}

const handleDeleteQuestionGuide = (data) => {
  questionGuide.id = data.id
  questionGuide.question = data.question

  // 弹窗确认删除
  Modal.confirm({
    title: t('delete_question_guide_title'),
    content: t('delete_question_guide_content'),
    okText: t('delete_question_guide_ok'),
    cancelText: t('delete_question_guide_cancel'),
    onOk() {
      deleteQuestionGuide({
        library_key: libraryInfo.value.library_key,
        id: questionGuide.id
      }).then(() => {
        message.success(t('delete_success'))
        time.value = new Date().getTime()
      })
    }
  })
}

const saveBannerImage = (url) => {
  let data = {
    banner_img_url: url,
    doc_id: state.doc_id,
    library_key: libraryInfo.value.library_key,
  }

  saveLibDocBannerImg(data).then(() => {
    time.value = new Date().getTime()
    message.success(t('save_success'))
  })
}

const onQuestionMessage = (message) => {
  if (message.action == 'edit' || message.action == 'add') {
    handleEditQuestionGuide(message.data)
  } else if (message.action == 'delete') {
    handleDeleteQuestionGuide(message.data)
  }
}

const onShortcutMessage = (message) => {
  if (message.action == 'edit' || message.action == 'add') {
    handleOpenShortcutSelector(message.data)
  }else if (message.action == 'delete') {
    handleDeleteShorcut(message.data)
  }else if(message.action == 'dragEnd'){
    handleShortcutSort(message.data)
  }
}

const onBannerMessage = (message) => {
  let data = message.data;

  uploadFile({
    file: data.file,
    category: 'library_doc_image'
  }).then((res) => {
    let url = res.data.link;
    saveBannerImage(url)
  })
}

const onPreviewPageMessage = (event) => {
  const message = event.data

  if (message.key == 'init') {
    initHomeDoc(message.data)

    if (isEdit.value) {
      setHomeEditStatus(true)
    }
  } else if (message.key == 'content') {
    handleEditDesc(message.data);
  } else if (message.key == 'title') {
    handleEditTitle(message.data);
  } else if (message.key == 'banner') {
    onBannerMessage(message);
  } else if (message.key == 'question') {
    onQuestionMessage(message);
  } else if (message.key == 'shortcut') {
    onShortcutMessage(message)
  }

  // 好像没用了，验证无用后可以删除
  // if (message.action === 'check_preview') {
  //   const iframeWindow = iframeRef.value.contentWindow

  //   iframeWindow.postMessage({ action: 'setPreview', token: getToken }, '*')
  // }
}

const toSettingPage = () => {
  router.push({
    path: '/public-library/config',
    query: {
      library_id: libraryId.value,
      library_key: libraryKey.value
    }
  })
}

const handleChangeLibrary = async (data) => {
  await router.replace({
    path: '/public-library/home',
    query: {
      library_id: data.id,
      library_key: data.library_key
    }
  })

  window.location.reload()
}

onMounted(async () => {
  window.addEventListener('message', onPreviewPageMessage)

  await getHomeInfo();

  fetchLibraryList();
})

onBeforeUnmount(() => {
  window.removeEventListener('message', onPreviewPageMessage)
})
</script>
