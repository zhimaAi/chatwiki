<style lang="less" scoped>
.editor-page {
  position: relative;
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  overflow: hidden;
  .loading-wrapper {
    position: absolute;
    left: 0;
    top: 0;
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
    padding: 24px 24px;
    // border-bottom: 1px solid #f0f0f0;
  }

  .editor-page-body {
    flex: 1;
    overflow: hidden;
    overflow-y: auto;
  }

  .header-left {
    flex: 1;
    overflow: hidden;
    display: flex;
    align-items: center;

    .title-wrapper {
      flex: 1;
      padding-right: 24px;
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
}
</style>

<template>
  <div class="editor-page" :class="{ 'is-loading': loading }">
    <div class="loading-wrapper" v-if="loading"><a-spin tip="Loading..."></a-spin></div>
    <div class="editor-page-header">
      <div class="header-left">
        <div class="title-wrapper">
          <DocTitle v-model:title="state.title" @blur="onTitleBlur" @input="onTitleInput" />
        </div>
      </div>
      <div class="header-right">
        <span class="last-save-time">{{ lastSaveTiemDate }} 自动保存草稿</span>
        <div class="edit-user" v-if="editUser.show && editUser.user_name">
          <i class="line"></i>
          <img class="user-avatar" :src="editUser.avatar" alt="" />
          <span class="user-name">{{ editUser.nick_name || editUser.user_name }}</span>
          <span class="edit-status-text">编辑中...</span>
        </div>
        <div class="action-box" v-if="isEdit">
          <a-button class="action-btn" @click="handleSaveDraft(true)">
            <svg-icon class="action-icon" name="save-draft" style="font-size: 16px"></svg-icon>
            <span>存草稿</span>
          </a-button>
        </div>

        <div class="action-box">
          <SharePopup ref="sharePopupRef" :docKey="state.doc_key">
            <a-button
              ><svg-icon class="action-icon" name="share" style="font-size: 16px"></svg-icon
              >分享</a-button
            >
          </SharePopup>
        </div>

        <div class="action-box" v-if="isEdit">
          <a-button @click="handleSeoSetting">
            <svg-icon
              class="action-icon"
              name="jibenpeizhi"
              style="font-size: 16px; color: #595959"
            ></svg-icon>
            <span>SEO设置</span>
          </a-button>
        </div>

        <div class="action-box" v-if="isEdit">
          <a-button class="action-btn" type="primary" @click="handlePublish">
            <svg-icon class="action-icon" name="gou" style="font-size: 16px"></svg-icon>
            <span>发布</span>
          </a-button>
        </div>

        <div class="action-box" v-if="!isEdit">
          <a-button class="action-btn" type="primary" @click="handleEdit">
            <svg-icon
              class="action-icon"
              name="edit"
              style="font-size: 16px; color: #fff"
            ></svg-icon>
            <span>编辑</span>
          </a-button>
        </div>
      </div>
    </div>
    <div class="editor-page-body">
      <MdDitor ref="MdDitorRef" @input="onContentInput" @blur="onContentBlur" />
    </div>
  </div>
  <ShareModal ref="shareModalRef" />
  <SeoSetting ref="seoSettingRef" @ok="saveSeoSuccess" />
</template>

<script setup>
import { getUser } from '@/api/manage'
import dayjs from 'dayjs'
import { saveLibDoc, saveDraftLibDoc, getLibDocInfo } from '@/api/public-library'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
import { ref, reactive, watch, computed, nextTick, onBeforeUnmount, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import MdDitor from './components/md-editor.vue'
import DocTitle from './components/doc-title.vue'
import ShareModal from '../components/share-modal.vue'
import SeoSetting from '../components/seo-setting.vue'
import SharePopup from '../components/share-popup.vue'

const emit = defineEmits(['changeDoc'])

const route = useRoute()

const { user_id } = useUserStore()
// const libraryId = computed(() => route.query.library_id)
const libraryKey = computed(() => route.query.library_key)
const docId = computed(() => route.query.doc_id)

const isEdit = ref(false)
const state = reactive({
  doc_key: '',
  title: '',
  content: '',
  update_time: '',
  is_draft: '',
  pid: '',
  seo_title: '',
  seo_desc: '',
  seo_keywords: ''
})

const editUser = reactive({
  show: false,
  user_name: '',
  avatar: '',
  nick_name: ''
})

const shareModalRef = ref(null)

const seoSettingRef = ref(null)

const handleSeoSetting = () => {
  let params = {
    library_key: libraryKey.value,
    seo_title: state.seo_title,
    seo_desc: state.seo_desc,
    seo_keywords: state.seo_keywords,
    doc_id: state.id
  }

  seoSettingRef.value.open(params)
}

const saveSeoSuccess = () => {
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
  }, 1000 * 60)
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

const onContentInput = (val) => {
  state.content = val
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
    library_key: libraryKey.value,
    doc_id: docId.value
  }).then((res) => {
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
    library_key: libraryKey.value,
    doc_id: docId.value
  })
    .then((res) => {
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

const handleSaveDraft = (showTips) => {
  clearAutoSaveTimer()

  let doc = MdDitorRef.value.getDoc()

  doc.title = state.title

  let data = {
    library_key: libraryKey.value,
    pid: state.pid,
    doc_id: docId.value,
    title: doc.title,
    content: doc.content,
    doc_type: 4
  }

  saveDraftLibDoc(data)
    .then(() => {
      if (showTips) {
        message.success('保存成功')
      }

      setLastSaveTimeDate()

      state.is_draft = 1

      isDocChange.value = false

      emit('changeDoc', { ...data, id: data.doc_id, is_draft: 1 })
    })
    .catch(() => {
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
      library_key: libraryKey.value
    })

    let editStatus = doc.data.edit_user && doc.data.edit_user != 0

    if (editStatus && doc.data.edit_user != user_id) {
      let user = await getEditUser(doc.data.edit_user)
      let name = user.data.nick_name || user.data.user_name
      editUser.show = true
      Modal.warning({
        title: name + '正在编辑此文档',
        content: '需要其他协作者结束编辑并发布后才可以继续编辑',
        okText: '知道了'
      })

      return
    }

    editUser.show = false

    let data = {
      library_key: libraryKey.value,
      pid: state.pid,
      doc_id: docId.value,
      title: state.title,
      content: state.content,
      doc_type: 4
    }

    saveDraftLibDoc(data)

    setMDEditor(true)
  } catch (error) {
    console.log(error)
    return
  }
}

const handlePublish = () => {
  let doc = MdDitorRef.value.getDoc()

  doc.title = state.title

  let data = {
    library_key: libraryKey.value,
    pid: state.pid,
    doc_id: docId.value,
    title: doc.title,
    doc_type: 4,
    content: doc.content
  }

  saveLibDoc(data).then(() => {
    message.success('发布成功')

    setLastSaveTimeDate()

    state.is_draft = 0
    setMDEditor(false)
    emit('changeDoc', { ...data, id: data.doc_id, is_draft: 0 })
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

watch(docId, () => {
  nextTick(() => {
    initDoc()
  })
})

onMounted(() => {
  document.addEventListener('keydown', onKeydown)

  initDoc()
})

onBeforeUnmount(() => {
  clearAutoSaveTimer()

  document.removeEventListener('keydown', onKeydown)
})
</script>
