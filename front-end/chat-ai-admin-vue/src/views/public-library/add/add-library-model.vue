<style lang="less" scoped>
.form-box {
  width: 568px;
  margin: 0 auto;
}

.form-item-tip {
  color: #999;
}
</style>

<template>
  <a-modal
    v-model:open="visible"
    title="添加知识库"
    @ok="handleOk"
    :confirmLoading="saveLoading"
    @cancel="handleCancel"
    width="640px"
  >
    <div class="form-box">
      <a-form :label-col="{ span: 5 }">
        <a-form-item ref="name" label="知识库名称" v-bind="validateInfos.library_name">
          <a-input
            v-model:value="formState.library_name"
            type="text"
            placeholder="请输入知识库名称，最多20个字"
            :maxlength="20"
          />
        </a-form-item>

        <a-form-item label="知识库简介">
          <a-textarea
            :maxlength="1000"
            v-model:value="formState.library_intro"
            placeholder="请输入知识库介绍"
          />
        </a-form-item>

        <a-form-item ref="name" label="知识库封面" v-bind="validateInfos.avatar">
          <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
          <div class="form-item-tip">请上传知识库封面，建议尺寸为100*100px.大小不超过100kb</div>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { LIBRARY_NORMAL_AVATAR, LIBRARY_OPEN_AVATAR } from '@/constants/index'
import { transformUrlData } from '@/utils/validate.js'
import { saveDraftLibDoc } from '@/api/public-library'
import { reactive, ref, onMounted, computed } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { createLibrary } from '@/api/library/index'
import AvatarInput from './components/avatar-input.vue'

// 设置全局默认的duration为（2秒）
message.config({
  duration: 2
})

const router = useRouter()
const visible = ref(false)

const type = computed(() => {
  return 1
})

let library_key = ''
let library_id = ''

const useForm = Form.useForm
const saveLoading = ref(false)

const defaultAvatar = type.value == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_OPEN_AVATAR
const formState = reactive({
  type: type.value,
  access_rights: 0,
  library_name: '',
  library_intro: '',
  use_model: '',
  model_config_id: '',
  library_files: undefined,
  avatar: defaultAvatar,
  avatar_file: undefined,
  is_offline: false,
  urls: '',
  doc_type: 1,
  file_name: '',
  is_qa_doc: type.value == 2 ? 1 : 0, // 0 普通文档 1QA文档
  qa_index_type: 1, // 1问题与答案一起生成索引 2仅对问题生成索引
  doc_auto_renew_frequency: 1
})

const rules = reactive({
  library_name: [{ required: true, message: '请输入库名称', trigger: 'change' }]
})

const { validate, validateInfos } = useForm(formState, rules)

const onAvatarChange = (data) => {
  formState.avatar = data.imageUrl
  formState.avatar_file = data.file
}

const handleCancel = () => {
  visible.value = false
}

const handleOk = () => {
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log(err, 'err')
    })
}

const saveForm = () => {
  let formData = new FormData()

  let newFormState = JSON.parse(JSON.stringify(formState)) // 深拷贝，不能改变原对象

  formData.append('type', formState.type)
  formData.append('access_rights', formState.access_rights)
  formData.append('library_name', formState.library_name)
  formData.append('library_intro', formState.library_intro)
  formData.append('use_model', newFormState.use_model)
  formData.append('model_config_id', formState.model_config_id)
  formData.append('avatar', formState.avatar_file ? formState.avatar_file : '')
  formData.append('is_offline', formState.is_offline)
  formData.append('urls', JSON.stringify(transformUrlData(formState.urls)))
  formData.append('doc_type', formState.doc_type)
  formData.append('file_name', formState.file_name)
  formData.append('is_qa_doc', formState.is_qa_doc)
  formData.append('qa_index_type', formState.qa_index_type)
  formData.append('doc_auto_renew_frequency', formState.doc_auto_renew_frequency)

  saveLoading.value = true

  createLibrary(formData)
    .then((res) => {
      message.success('创建成功')
      library_id = res.data.id
      library_key = res.data.library_key
      toHome(res.data.id)
    })
    .catch(() => {
      saveLoading.value = false
    })
}

const addDoc = async () => {
  let docName = '无标题文档'

  let data = {
    library_key: library_key,
    doc_id: '',
    doc_type: '4',
    pid: 0,
    title: docName,
    content: '# 无标题文档',
    sort: 0
  }

  const res = await saveDraftLibDoc(data)

  return res
}

const toHome = () => {
  addDoc().then((res) => {
    let doc_id = res.data.doc_id

    router.replace({
      path: '/public-library/editor',
      query: {
        library_id: library_id,
        library_key: library_key,
        doc_id: doc_id
      }
    })
  })
}

const show = () => {
  visible.value = true
}

defineExpose({
  show
})

onMounted(() => {})
</script>
