<template>
  <a-modal
    v-model:open="open"
    :confirm-loading="saving"
    :maskClosable="false"
    :title="step == 1 ? '导入飞书知识库' : '请选择需要同步的文件'"
    width="646px"
    @ok="save"
  >
    <template v-if="step == 1">
      <a-alert class="zm-alert-info" show-icon type="info">
        <template #message>
          请先在<a href="https://open.feishu.cn/document/server-docs/api-call-guide/terminology" target="_blank">飞书开发者后台</a>创建自建应用获取App ID 和 App Secret，并将需要同步的 <a
          href="https://open.feishu.cn/document/server-docs/docs/wiki-v2/wiki-qa#b5da330b" target="_blank">知识库授权给对应应用</a>；
          <a-popover>
            <template #title>
              <a-image width="600px" :src="RedirectExampleImg"/>
            </template>
            在飞书应用中填写重定向URL：<a @click="handleCopy(redirectUrl)" title="点击复制">{{redirectUrl}}</a>
          </a-popover>
          。完成后，系统将自动同步知识库目录下.docx格式的云文档
        </template>
      </a-alert>
      <a-form
        class="mt16"
        layout="vertical"
        ref="formRef"
        :model="formState"
        :rules="rules"
      >
        <a-form-item name="feishu_app_id" label="AppId">
          <a-input v-model:value.trim="formState.feishu_app_id" placeholder="请输入飞书 AppId"/>
        </a-form-item>
        <a-form-item name="feishu_app_secret" label="AppSecret">
          <a-input v-model:value.trim="formState.feishu_app_secret" placeholder="请输入飞书 AppSecret"/>
        </a-form-item>
      </a-form>
    </template>
    <template v-else>
      <a-form
        layout="vertical"
        ref="formRef"
        :model="formState"
        :rules="rules"
      >
<!--        <a-form-item name="doc_auto_renew_frequency" label="更新频率">-->
<!--          <a-select v-model:value="formState.doc_auto_renew_frequency" style="width: 100%">-->
<!--            <a-select-option :value="1">不自动更新</a-select-option>-->
<!--            <a-select-option :value="2">每天</a-select-option>-->
<!--            <a-select-option :value="3">每3天</a-select-option>-->
<!--            <a-select-option :value="4">每7天</a-select-option>-->
<!--            <a-select-option :value="5">每30天</a-select-option>-->
<!--          </a-select>-->
<!--        </a-form-item>-->
<!--        <a-form-item-->
<!--          v-if="formState.doc_auto_renew_frequency > 1"-->
<!--          name="doc_auto_renew_minute"-->
<!--          label="更新时间"-->
<!--        >-->
<!--          <a-time-picker-->
<!--            v-model:value="formState.doc_auto_renew_minute"-->
<!--            format="HH:mm"-->
<!--            valueFormat="HH:mm"-->
<!--          />-->
<!--        </a-form-item>-->
        <a-form-item label="文件范围">
          <a-radio-group v-model:value="formState.file_type">
            <a-radio :value="1">权限范围内所有文件</a-radio>
            <a-radio :value="2">部分文件</a-radio>
          </a-radio-group>
          <div class="file-list" v-if="formState.file_type == 2">
            <a-directory-tree
              v-model:checkedKeys="formState.feishu_document_id_list"
              multiple
              checkable
              :field-names="{key: 'token'}"
              :tree-data="filesTreeData"
            >
              <template #title="{name, type}">
                <template v-if="type === 'docx'">{{name}}.docx</template>
                <template v-else>{{name}}</template>
              </template>
            </a-directory-tree>
          </div>
        </a-form-item>
      </a-form>
    </template>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue'
import {addLibraryFile, getFeishuDocFileList} from "@/api/library/index.js"
import {message} from 'ant-design-vue'
import {convertTime, copyText, objectToQueryString, strToBase64} from "@/utils/index.js";
import {useUserStore} from '@/stores/modules/user'
import RedirectExampleImg from '@/assets/img/library/feishu-redirct-example.png'

const emit = defineEmits(['ok'])
const props = defineProps({
  libraryId: {
    type: [Number, String],
    default: ''
  },
})

const userStore = useUserStore()
const redirectUrl = `${window.location.origin}/manage/feishuUserAuthLogin/callback`

const formRef = ref(null)
const formStateDefault = {
  feishu_app_id: '',
  feishu_app_secret: '',
  // doc_auto_renew_frequency: 1,
  // doc_auto_renew_minute: '',
  user_access_token: '',
  file_type: 1,
  feishu_document_id_list: []
}
const open = ref(false)
const saving = ref(false)
const step = ref(1)
const formState = reactive({
  feishu_document_id_list: []
})
const rules = reactive({
  feishu_app_id: {
    message: '请输入飞书 AppId',
    required: true
  },
  feishu_app_secret: {
    message: '请输入飞书 AppSecret',
    required: true
  },
})
const filesTreeData = ref([])

function show(p = null) {
  Object.assign(formState, JSON.parse(JSON.stringify(formStateDefault)))
  if (p) {
    Object.assign(formState, p)
    searchFiles()
    step.value = 2
  } else {
    step.value = 1
  }
  open.value = true
}

function searchFiles() {
  const {feishu_app_id, feishu_app_secret, user_access_token} = formState
  getFeishuDocFileList({feishu_app_id, feishu_app_secret, user_access_token}).then(res => {
    filesTreeData.value = filterDocxTree(res?.data || [])
  })
}

function filterDocxTree(list = []) {
  return list
    .map(item => {
      // 如果有 children，先递归处理
      if (Array.isArray(item.children)) {
        item.checkable = false
        const children = filterDocxTree(item.children)
        return { ...item, children }
      }
      item.checkable = true
      return item
    })
    .filter(item => {
      // 保留 docx
      if (item.type === 'docx') return true

      // 保留 children 中还有 docx 的 folder
      if (item.type === 'folder' && item.children?.length) return true

      return false
    })
}

function getAllDocxTokens(list = []) {
  const result = []

  function dfs(nodes) {
    for (const item of nodes) {
      if (item.type === 'docx' && item.token) {
        result.push(item.token)
      }
      if (Array.isArray(item.children)) {
        dfs(item.children)
      }
    }
  }

  dfs(list)
  return result
}

const handleCopy = (text) => {
  copyText(text)
  message.success('复制成功')
}

function save() {
  formRef.value.validate().then(() => {
    if (step.value == 1) {
      const hide = message.loading('正在进行验证...')
      let host = ''
      if (import.meta.env.MODE !== 'production') {
        host = `http://${import.meta.env.MODE}.zhima_chat_ai.applnk.cn`
      }
      const baseData = {
        feishu_app_id: formState.feishu_app_id,
        feishu_app_secret: formState.feishu_app_secret,
      }
      const redirectQuery = {
        ...baseData,
        id: props.libraryId
      }
      let query = {
        ...baseData,
        token: userStore.getToken ?? '',
        feishu_frontend_auth_redirect_url: strToBase64(`${host}/#/library/details/knowledge-document?${objectToQueryString(redirectQuery)}`)
      }
      window.location.href = `${host}/manage/feishuUserAuthLogin/redirect?${objectToQueryString(query)}`
    } else {
      // if (formState.doc_auto_renew_frequency > 1 && !formState.doc_auto_renew_minute) return message.error('请选择更新时间')
      if (formState.file_type == 2) {
        if (!formState.feishu_document_id_list.length) return message.error('请选择同步文件')
      } else {
        formState.feishu_document_id_list = getAllDocxTokens(filesTreeData.value)
        if (!formState.feishu_document_id_list.length) return message.error('暂无可同步文件')
      }
      let params = {
        library_id: props.libraryId,
        doc_type: 6,
        ...formState,
        doc_auto_renew_minute: convertTime(formState.doc_auto_renew_minute),
        feishu_document_id_list: formState.feishu_document_id_list.toString(),
      }
      delete params.file_type
      delete params.user_access_token
      saving.value = true
      addLibraryFile(params).then(res => {
        emit('ok')
        message.success('添加完成')
        open.value = false
      }).finally(() => {
        saving.value = false
      })
    }
  })
}

defineExpose({
  show
})
</script>

<style scoped lang="less">
.mt16 {
  margin-top: 16px;
}

.file-list {
  margin-top: 8px;

  :deep(.ant-checkbox-group) {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }
}
</style>
