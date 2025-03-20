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
    height: 80px;
    padding: 0 24px;
    border-bottom: 1px solid #f0f0f0;

    .toolbar-right {
      display: flex;
      align-items: center;
    }
    .action-item {
      margin-right: 8px;
      &:last-child {
        margin-right: 0;
      }
    }
    .action-icon {
      margin-right: 4px;
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

  .preview-box {
    flex: 1;
    overflow: hidden;
    text-align: center;
    display: flex;
    align-items: center;
    justify-content: center;

    .mobile-box {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 100%;
      width: 100%;
      padding: 24px 0;
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
}
</style>

<template>
  <div class="home-preview">
    <div class="toolbar-wrapper">
      <div class="toolbar-left">
        <div class="toolbar-item">
          <a-radio-group v-model:value="styleType">
            <a-radio-button value="pc">PC端样式</a-radio-button>
            <a-radio-button value="mobile">移动端样式</a-radio-button>
          </a-radio-group>
        </div>
      </div>
      <div class="toolbar-right">
        <div class="edit-user" v-if="editUser.show && editUser.user_name">
          <img class="user-avatar" :src="editUser.avatar" alt="" />
          <span class="user-name">{{ editUser.nick_name || editUser.user_name }}</span>
          <span class="edit-status-text">编辑中...</span>
        </div>

        <div class="action-item">
          <SharePopup
            ref="sharePopupRef"
            base-url="/open/home"
            :docKey="state.library_key"
            :libraryId="id"
          >
            <a-button
              ><svg-icon class="action-icon" name="share" style="font-size: 16px"></svg-icon>
              分享</a-button
            >
          </SharePopup>
        </div>
        <div class="action-item" v-if="isEdit">
          <a-button @click="handleSeoSetting">
            <svg-icon
              class="action-icon"
              name="jibenpeizhi"
              style="font-size: 16px; color: #595959"
            ></svg-icon>
            <span>SEO设置</span>
          </a-button>
        </div>
        <div class="action-item">
          <a-button type="primary" :loading="publishLoading" @click="handlePublish" v-if="isEdit"
            ><CheckOutlined /> 发布</a-button
          >
          <a-button type="primary" @click="handleHomeEdit" v-else>
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

    <div class="preview-box">
      <div class="mobile-box" v-if="styleType == 'mobile'">
        <PhoneBox>
          <iframe
            ref="iframeRef"
            class="iframe"
            :src="previewUrl"
            frameborder="0"
            v-if="state.library_key"
          ></iframe>
        </PhoneBox>
      </div>

      <div class="pc-box" v-else>
        <iframe
          ref="iframeRef"
          class="iframe"
          :src="previewUrl"
          frameborder="0"
          v-if="state.library_key"
        ></iframe>
      </div>
    </div>
    <EditTitle ref="editTitleRef" @ok="handleSaveTitle" :confirm-loading="saveTitleLoading" />
    <EditDesc ref="editDescRef" @ok="handleSaveDesc" :confirm-loading="saveDescLoading" />
    <EditQuestionGuide
      ref="editQuestionGuideRef"
      @ok="handleSaveQuestionGuide"
      :confirm-loading="saveQuestionGuideLoading"
    />
    <SeoSetting ref="seoSettingRef" @ok="saveSeoSuccess" />
  </div>
</template>

<script setup>
import { getLibraryInfo } from '@/api/library/index'
import {
  saveDraftLibDoc,
  saveQuestionGuide,
  deleteQuestionGuide,
  saveLibDoc,
  getLibDocInfo
} from '@/api/public-library'
import { getUser } from '@/api/manage'
import { useUserStore } from '@/stores/modules/user'
import { ref, computed, reactive, onMounted, onBeforeUnmount } from 'vue'
import { useRoute } from 'vue-router'
import { CheckOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import EditTitle from './components/edit-title.vue'
import EditDesc from './components/edit-desc.vue'
import EditQuestionGuide from './components/edit-question-guide.vue'
import PhoneBox from './components/phone-box.vue'
import SeoSetting from '../components/seo-setting.vue'
import SharePopup from '../components/share-popup.vue'

const route = useRoute()
const { getToken, user_id } = useUserStore()

const id = computed(() => route.query.library_id)

const iframeRef = ref(null)
const time = ref(new Date().getTime())
const isEdit = ref(false)
const state = reactive({
  doc_id: '',
  title: '',
  content: '',
  library_intro: '',
  library_name: '',
  library_key: ''
})

const docState = reactive({
  doc_id: '',
  is_pub: '',
  doc_key: '',
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

// 保存
const saveDraftDoc = async () => {
  let res = await saveDraftLibDoc({
    library_key: state.library_key,
    doc_id: state.doc_id,
    title: state.title,
    content: state.content,
    is_index: 1
  }).then((res) => {
    state.doc_id = res.data.doc_id
    docState.doc_id = res.data.doc_id
    time.value = new Date().getTime()
    return res
  })

  return res
}
const handleSaveDraft = async () => {
  let res = await saveDraftDoc()

  if (res) {
    message.success('保存成功')
  }
}
// 发布
const publishLoading = ref(false)
const handlePublish = () => {
  let data = {
    library_key: state.library_key,
    doc_id: state.doc_id,
    title: state.title,
    content: state.content,
    is_index: 1
  }

  publishLoading.value = true

  saveLibDoc(data)
    .then((res) => {
      state.doc_id = res.data.doc_id
      docState.doc_id = res.data.doc_id
      publishLoading.value = false
      message.success('发布成功')
      setHomeEditStatus(false)

      time.value = new Date().getTime()
    })
    .catch(() => {
      publishLoading.value = false
    })
}

const seoSettingRef = ref(null)

const handleSeoSetting = () => {
  let params = {
    library_key: state.library_key,
    seo_title: docState.seo_title,
    seo_desc: docState.seo_desc,
    seo_keywords: docState.seo_keywords,
    doc_id: state.doc_id
  }

  seoSettingRef.value.open(params)
}

const saveSeoSuccess = () => {
  getHomeInfo()
  time.value = new Date().getTime()
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
    library_key: state.library_key,
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

const styleType = ref('pc')

let prviewHost = import.meta.env.VITE_OPEN_DOC_HOST || ''

const previewUrl = computed(() => {
  return `${prviewHost}/manage/libDocHome/${state.library_key}?&v=${time.value}&token=${getToken}`
})

const getData = () => {
  getLibraryInfo({ id: id.value }).then((res) => {
    Object.assign(state, res.data)
  })
}

const getHomeInfo = () => {
  getLibDocInfo({ doc_id: state.doc_id, library_key: state.library_key }).then((res) => {
    Object.assign(docState, res.data)

    let editStatus = res.data.edit_user && res.data.edit_user != 0

    if (editStatus && res.data.edit_user != user_id) {
      editUser.show = true
      getEditUser(res.data.edit_user)
    }

    if (editStatus && res.data.edit_user == user_id) {
      setHomeEditStatus(true)
    } else {
      setHomeEditStatus(false)
    }
  })
}

const initHomeDoc = (data) => {
  state.doc_id = data.doc_id
  state.title = data.title
  state.content = data.content

  getHomeInfo()
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
    title: '删除问题引导',
    content: '确定要删除该问题引导吗？',
    okText: '确定',
    cancelText: '取消',
    onOk() {
      deleteQuestionGuide({
        library_key: state.library_key,
        id: questionGuide.id
      }).then(() => {
        message.success('删除成功')
        time.value = new Date().getTime()
      })
    }
  })
}

const handleEditMsg = (message) => {
  if (message.key == 'title') {
    handleEditTitle(message.data)
  }

  if (message.key == 'content') {
    handleEditDesc(message.data)
  }

  if (message.key == 'question') {
    handleEditQuestionGuide(message.data)
  }
}

const handleAddMsg = (message) => {
  if (message.key == 'question') {
    handleEditQuestionGuide(message.data)
  }
}

const handleDeleteMsg = (message) => {
  if (message.key == 'question') {
    handleDeleteQuestionGuide(message.data)
  }
}

const handleHomeEdit = async () => {
  try {
    let doc = await getLibDocInfo({
      doc_id: state.doc_id,
      library_key: state.library_key
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

    saveDraftDoc()

    // setHomeEditStatus(true)
    isEdit.value = true
    editUser.show = false
  } catch (error) {
    return
  }
}

const setHomeEditStatus = (type) => {
  isEdit.value = type

  const iframeWindow = iframeRef.value.contentWindow

  iframeWindow.postMessage({ action: 'setEditStatus', data: { type: isEdit.value } }, '*')
}

const handleMessage = (event) => {
  const message = event.data

  if (message.action === 'check_preview') {
    const iframeWindow = iframeRef.value.contentWindow

    iframeWindow.postMessage({ action: 'setPreview', token: getToken }, '*')
  }

  if (message.action === 'init') {
    initHomeDoc(message.data)

    if (isEdit.value) {
      setHomeEditStatus(true)
    }
  }
  // 处理不同类型的消息
  if (message.action === 'edit') {
    // 在这里处理 data 中的数据
    handleEditMsg(message)
  }

  if (message.action === 'add') {
    handleAddMsg(message)
  }

  if (message.action === 'delete') {
    handleDeleteMsg(message)
  }
}

onMounted(() => {
  getData()

  window.addEventListener('message', handleMessage)
})

onBeforeUnmount(() => {
  window.removeEventListener('message', handleMessage)
})
</script>
