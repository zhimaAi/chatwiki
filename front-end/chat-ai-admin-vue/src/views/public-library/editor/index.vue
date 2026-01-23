<style lang="less" scoped>
.editor-page {
  position: relative;
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  overflow: hidden;

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
        color: #21a665;
        background-color: #cafce4;
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

  .loading-wrapper {
    position: absolute;
    left: 0;
    top: 64px;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 999;
    background-color: rgba(255, 255, 255, 1);
  }

  &.is-loading {
    .editor-page-body {
      opacity: 0;
    }
  }

  .editor-page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 60px;
    background: #F2F4F7;
  }

  .editor-page-body {
    flex: 1;
    overflow: hidden;
    overflow-y: auto;
    padding: 0 60px;
  }

  .header-left {
    flex: 1;
    overflow: hidden;

    .title-wrapper {
      display: flex;
      flex: 1;
      align-items: center;
      padding-right: 24px;

      .doc-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 24px;
        height: 24px;
        margin-right: 8px;
        font-size: 20px;
      }

      .doc-title {
        flex: 1;
        line-height: 32px;
        min-height: 32px;
        display: flex;
        align-items: center;

        .doc-title-text {
          line-height: 24px;
          font-size: 24px;
          font-weight: 600;
          color: #262626;
          &.title-highlight{
            color: #2475fc;
            cursor: pointer;
          }
        }
      }
    }

    .doc-last-update {
      line-height: 22px;
      padding-left: 32px;
      margin-top: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #8c8c8c;
    }
  }

  .header-right {
    display: flex;
    align-items: center;

    .action-box {
      margin-right: 8px;

      &:last-child {
        margin-right: 0;
      }

      .action-btn {
        display: flex;
        align-items: center;
      }

      .action-icon {
        margin-right: 4px;
      }
    }

    .last-save-time {
      margin-right: 16px;
      font-size: 14px;
      color: #595959;
    }
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

  .seo-setting-box {
    padding: 0 24px 24px 24px;
  }
}
</style>

<template>
  <div class="editor-page" :class="{ 'is-loading': loading }">
    <div class="toolbar-wrapper">
      <div class="toolbar-left">
        <div class="library-info">
          <div class="library-info-content">
            <wiki-dropdown :libraryInfo="libraryInfo" :libraryList="libraryList"
              @change="handleChangeLibrary"></wiki-dropdown>

            <div class="publish-status" v-if="libraryInfo.access_rights == 0">
              <ExclamationCircleFilled class="status-icon" />
              <span class="status-name">{{ t('not_public') }}</span>
            </div>
            <div class="publish-status is-publish" v-if="libraryInfo.access_rights == 1">
              <CheckCircleFilled class="status-icon" />
              <span class="status-name">{{ t('published') }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="toolbar-right">
        <div class="action-item">
          <a-input :placeholder="t('knowledge_search')" v-model:value="searchValue" @pressEnter="handleSearch" allowClear>
            <template #suffix>
              <SearchOutlined style="color: rgba(0, 0, 0, 0.45)" @click="handleSearch" />
            </template>
          </a-input>
        </div>

        <div class="action-item">
          <a-button type="primary" :loading="publishLoading" @click="handlePublish" v-if="isEdit">
            <CheckOutlined /><span class="action-name">{{ t('publish') }}</span>
          </a-button>
          <a-button type="primary" @click="handleEdit" v-else>
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
            <CheckCircleFilled style="color: #21a665; font-size: 14px" class="status-icon actin-icon" v-if="isSetSeo" />
            <ExclamationCircleFilled  style="color: #ccc; font-size: 14px" class="status-icon actin-icon" v-else />
            <span class="action-name">{{ t('seo') }}</span>
          </a-button>
        </div>

        <div class="action-item">
          <a-button @click="toSettingPage">
            <svg-icon class="actin-icon" name="jibenpeizhi" style="font-size: 14px; color: #595959"></svg-icon>
            <span class="action-name">{{ t('settings') }}</span>
          </a-button>
        </div>
      </div>
    </div>
    <div class="loading-wrapper" v-if="loading"><a-spin tip="Loading..."></a-spin></div>
    <div class="editor-page-header">
      <div class="header-left">
        <div class="title-wrapper" v-if="isEdit">
          <icon-selector @select="handleSelectDocIcon">
            <span class="doc-icon">{{ state.doc_icon || levelDocIcon }}</span>
          </icon-selector>
          <div class="doc-title">
            <DocTitleInput v-model:title="state.title"  @blur="onTitleBlur" @input="onTitleInput"  />
          </div>
        </div>
        <div class="title-wrapper" v-else>
          <span class="doc-icon">{{ state.doc_icon || levelDocIcon }}</span>
          <div class="doc-title">
            <div class="doc-title-text" @click="handleTitleClick" :class="{ 'title-highlight': titleHighlight }" >{{ state.title }}</div>
          </div>
        </div>

        <div class="doc-last-update">{{ t('last_update_time') }}{{ state.update_time_str }}</div>
      </div>
      <div class="header-right">
        <span class="last-save-time">{{ lastSaveTiemDate }} {{ t('auto_save_draft') }}</span>
        <div class="edit-user" v-if="editUser.show && editUser.user_name">
          <i class="line"></i>
          <img class="user-avatar" :src="editUser.avatar" alt="" />
          <span class="user-name">{{ editUser.nick_name || editUser.user_name }}</span>
          <span class="edit-status-text">{{ t('editing') }}</span>
        </div>
        <div class="action-box" v-if="isEdit">
          <a-button class="action-btn" @click="handleSaveDraft(true)">
            <svg-icon class="action-icon" name="save-draft" style="font-size: 16px"></svg-icon>
            <span>{{ t('save_draft') }}</span>
          </a-button>
        </div>

        <div class="action-box" v-if="isEdit">
          <a-button class="action-btn" @click="handlePreview()">
            <svg-icon class="action-icon" name="eye-open" style="font-size: 16px"></svg-icon>
            <span>{{ t('preview') }}</span>
          </a-button>
        </div>
      </div>
    </div>

    <div class="editor-page-body">
      <MdDitor ref="MdDitorRef" @input="onContentInput" @blur="onContentBlur" />
    </div>
    <SeoSetting ref="seoSettingRef" @ok="saveSeoSuccess" />
  </div>
</template>

<script setup>
import { OPEN_BOC_BASE_URL } from '@/constants/index'
import { getUser } from '@/api/manage'
import dayjs from 'dayjs'
import { saveLibDoc, saveDraftLibDoc, getLibDocInfo, getPreviewUrl } from '@/api/public-library'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
import { usePublicLibraryStore } from '@/stores/modules/public-library'
import { ref, reactive, watch, computed, nextTick, onBeforeUnmount, onMounted } from 'vue'
import { CheckOutlined, ExclamationCircleFilled, CheckCircleFilled, SearchOutlined } from '@ant-design/icons-vue'
import { useCopyShareUrl } from '@/hooks/web/useCopyShareUrl'
import { useI18n } from '@/hooks/web/useI18n'
import { message, Modal } from 'ant-design-vue'
import MdDitor from './components/md-editor.vue'
import DocTitleInput from './components/doc-title-input.vue'
import SeoSetting from '../components/seo-setting.vue'
import WikiDropdown from '../components/wiki-dropdown.vue'
import IconSelector from './components/icon-selector.vue'

const { t } = useI18n('views.public-library.editor.index')

const autoSaveTime = 1000 * 15
const emit = defineEmits(['changeDoc'])

const route = useRoute()
const router = useRouter()

const libraryStore = usePublicLibraryStore()

const { user_id } = useUserStore()

const libraryInfo = computed(() => libraryStore.libraryInfo)
const docTreeStateMap = computed(() => libraryStore.docTreeStateMap)
const docId = computed(() => route.query.doc_id);

const levelDocIcon = computed(() => {
  if(!docTreeStateMap.value[docId.value]){
    return '';
  }

  return docTreeStateMap.value[docId.value].doc_icon || '';   
});
const isEdit = ref(false)
const searchValue = ref('')
const state = reactive({
  library_key: '',
  doc_id: docId.value,
  doc_key: '',
  doc_icon: '',
  title: '',
  content: '',
  update_time_str: '',
  is_draft: '',
  pid: '',
  seo_title: '',
  seo_desc: '',
  seo_keywords: ''
})

const isSetSeo = computed(() => { 
  return state.seo_title || state.seo_desc || state.seo_keywords
})

const editUser = reactive({
  show: false,
  user_name: '',
  avatar: '',
  nick_name: ''
})


const handleSelectDocIcon = (icon) => {
  state.doc_icon = icon.content;

  onSaveDoc()
}

// 标题高亮状态计算属性
const titleHighlight = computed(() => {
  return libraryInfo.value.share_url != '' && libraryInfo.value.access_rights == 1
})

const handleTitleClick = () => {
  if (!titleHighlight.value) {
    return
  }

  const docUrl = libraryInfo.value.share_url + OPEN_BOC_BASE_URL + '/doc/' + state.doc_key

  window.open(docUrl)
}

const { copyShareUrl } = useCopyShareUrl()

const handleCopyShareUrl = async () => {
  const docUrl = OPEN_BOC_BASE_URL + '/doc/' + state.doc_key
  await copyShareUrl(docUrl)
}

const seoSettingRef = ref(null)
const handleOpenSeoSetting = () => {
  seoSettingRef.value.open({ ...state })
}

const saveSeoSuccess = (data) => {
  state.seo_title = data.seo_title
  state.seo_desc = data.seo_desc
  state.seo_keywords = data.seo_keywords

  getDoc()
}

let autoSaveTimer = null
// 操作后自动保存草稿
const autoSaveDraft = () => {
  if (autoSaveTimer) {
    clearTimeout(autoSaveTimer)
    autoSaveTimer = null
  }

  autoSaveTimer = setTimeout(() => {
    if (isDocChange.value) {
      onSaveDoc()
    }
  }, autoSaveTime)
}

const clearAutoSaveTimer = () => {
  if (autoSaveTimer) {
    clearTimeout(autoSaveTimer)
    autoSaveTimer = null
  }
}

const loading = ref(true)

const isDocChange = ref(false)

const onTitleInput = (val) => {
  state.title = val
  isDocChange.value = true

  autoSaveDraft()
}

const onTitleBlur = () => {
  if (isDocChange.value) {
    onSaveDoc()
  }
}

const onContentInput = () => {
  isDocChange.value = true

  autoSaveDraft()
}

const onContentBlur = () => {
  if (isDocChange.value) {
    onSaveDoc()
  }
}

const MdDitorRef = ref(null)
const lastSaveTiemDate = ref('')

const setLastSaveTimeDate = () => {
  lastSaveTiemDate.value = dayjs().format('HH:mm')
}
const getDoc = () => {
  getLibDocInfo({
    library_key: libraryInfo.value.library_key,
    doc_id: docId.value
  }).then((res) => {
    res.data.library_key = libraryInfo.value.library_key
    res.data.doc_id = docId.value
    res.data.update_time_str = dayjs(res.data.update_time * 1000).format('YYYY-MM-DD HH:mm:ss')

    Object.assign(state, res.data)
  })
}

const setMDEditor = (type) => {
  isEdit.value = type

  MdDitorRef.value.init(type ? 'editor' : 'preview', state.content)
}
const initDoc = () => {
  loading.value = true
  editUser.show = false

  getLibDocInfo({
    library_key: libraryInfo.value.library_key,
    doc_id: docId.value
  })
    .then((res) => {
      res.data.library_key = libraryInfo.value.library_key
      res.data.doc_id = docId.value
      res.data.update_time_str = dayjs(res.data.update_time * 1000).format('YYYY-MM-DD HH:mm:ss')

      Object.assign(state, res.data)

      let editStatus = res.data.edit_user && res.data.edit_user != 0

      if (editStatus && res.data.edit_user != user_id) {
        editUser.show = true
        getEditUser(res.data.edit_user)
      }

      if (editStatus && res.data.edit_user == user_id) {
        setMDEditor(true)
      } else {
        setMDEditor(false)
      }

      setLastSaveTimeDate(res.data.update_time * 1000)

      loading.value = false
    })
    .catch(() => {
      loading.value = false
    })
}

const saveDraftLoading = ref(false)

const handleSaveDraft = (showTips) => {
  clearAutoSaveTimer()

  let doc = MdDitorRef.value.getDoc()

  state.content = doc.content

  doc.title = state.title

  let data = {
    library_key: libraryInfo.value.library_key,
    pid: state.pid,
    doc_id: docId.value,
    title: doc.title,
    content: doc.content,
    doc_icon: state.doc_icon,
    doc_type: 4,
  }
  if (showTips) {
    saveDraftLoading.value = true
  }

  saveDraftLibDoc(data)
    .then(() => {
      if (showTips) {
        saveDraftLoading.value = false
        message.success(t('save_success'))
      }

      setLastSaveTimeDate()

      state.is_draft = 1

      isDocChange.value = false

      emit('changeDoc', { ...data, id: data.doc_id, is_draft: 1 })
    })
    .catch(() => {
      saveDraftLoading.value = false
      autoSaveDraft()
    })
}

const getEditUser = async (userId) => {
  let res = await getUser({ id: userId })
  editUser.avatar = res.data.avatar
  editUser.user_name = res.data.user_name
  editUser.nick_name = res.data.nick_name
  return res
}

const handleEdit = async () => {
  try {
    let doc = await getLibDocInfo({
      doc_id: docId.value,
      library_key: libraryInfo.value.library_key
    })

    let editStatus = doc.data.edit_user && doc.data.edit_user != 0

    if (editStatus && doc.data.edit_user != user_id) {
            let user = await getEditUser(doc.data.edit_user)
            let name = user.data.nick_name || user.data.user_name
            editUser.show = true
            Modal.warning({
              title: t('editing_document', { name }),
              content: t('need_collaborator_finish'),
              okText: t('got_it')
            })
    
            return
          }
    editUser.show = false

    let data = {
      library_key: libraryInfo.value.library_key,
      pid: state.pid,
      doc_id: docId.value,
      title: state.title,
      content: state.content,
      doc_type: 4,
      doc_icon: state.doc_icon,
    }

    saveDraftLibDoc(data)

    setMDEditor(true)
  } catch (error) {
    console.log(error)
    return
  }
}

const publishLoading = ref(false)
const handlePublish = () => {
  let doc = MdDitorRef.value.getDoc()

  state.content = doc.content

  doc.title = state.title

  let data = {
    library_key: libraryInfo.value.library_key,
    pid: state.pid,
    doc_id: docId.value,
    title: doc.title,
    doc_type: 4,
    content: doc.content,
    doc_icon: state.doc_icon,
  }
  publishLoading.value = true
  saveLibDoc(data)
    .then(() => {
      publishLoading.value = false
      message.success(t('publish_success'))

      setLastSaveTimeDate()

      state.is_draft = 0
      setMDEditor(false)
      emit('changeDoc', { ...data, id: data.doc_id, is_draft: 0 })
    })
    .catch(() => {
      publishLoading.value = false
    })
}

const onSaveDoc = (data, status) => {
  handleSaveDraft(status === 'manuallySave')
}

const onKeydown = (event) => {
  if (event.ctrlKey && event.key === 's') {
    event.preventDefault()

    onSaveDoc('manuallySave')
  }
}

const handlePreview = () => {
  getPreviewUrl({
    library_key: libraryInfo.value.library_key,
    doc_id: docId.value
  }).then((res) => {
    window.open(
      `${libraryInfo.value.share_url}${OPEN_BOC_BASE_URL}/doc/${state.doc_key}?preview=${res.data.preview_key}`,
      '_blank'
    )
  })
}

const toSettingPage = () => {
  console.log(libraryInfo.value)
  router.push({
    path: '/public-library/config',
    query: {
      library_id: libraryInfo.value.id,
      library_key: libraryInfo.value.library_key
    }
  })
}

// 知识库列表
const libraryList = computed(() => libraryStore.libraryList)

const fetchLibraryList = () => {
  libraryStore.getLibraryList()
}

const handleChangeLibrary = async (data) => {
  await router.replace({
    path: '/public-library/home',
    query: {
      library_id: data.id,
      library_key: data.library_key
    }
  })

  window.location.reload();
}

const handleSearch = () => {
  message.warn(t('edit_mode_no_search'))
}

watch(docId, () => {
  nextTick(() => {
    initDoc()
  })
})

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
  fetchLibraryList()
  initDoc()
})

onBeforeUnmount(() => {
  clearAutoSaveTimer()

  document.removeEventListener('keydown', onKeydown)
})
</script>
