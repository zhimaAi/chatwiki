<style lang="less" scoped>
.add-library-page {
  padding: 24px;

  .form-box {
    width: 568px;
    margin: 0 auto;
  }
}

.model-icon {
  height: 18px;
}

.form-item-tip {
  color: #999;
}

.card-box {
  display: flex;
  gap: 8px;
}

.use-model-item {
  position: relative;
  width: 226px;
  height: 124px;
  border-radius: 2px;
  border: 2px solid #d9d9d9;
  cursor: pointer;
  padding: 15px;
  margin-bottom: 10px;
}

.use-model-item-top {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  font-size: 14px;
  color: #595959;
}

.use-model-item-top-icon {
  margin-right: 5px;
}

.use-model-item-top.active {
  color: #2475fc;
  font-weight: bolder;
}

.use-model-item.active {
  border: 2px solid #2475fc;

  .check-arrow {
    position: absolute;
    display: block;
    right: -1px;
    bottom: -1px;
    width: 24px;
    height: 24px;
    font-size: 24px;
    color: #fff;
  }

  .retrieval-mode-title {
    color: #2475fc;
  }
}

.upload-document-type-box {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  .type-item {
    position: relative;
    width: calc((100% - 8px) / 2);
    cursor: pointer;
    padding: 16px;
    display: flex;
    border: 1px solid #d9d9d9;
    border-radius: 2px;
    box-shadow: none;
    transition: box-shadow 1s;
    &:hover {
      box-shadow:
        0px 5px 5px -3px rgba(0, 0, 0, 0.1),
        0px 8px 10px 1px rgba(0, 0, 0, 0.06),
        0px 3px 14px 2px rgba(0, 0, 0, 0.05);
    }
    &.active {
      border: 2px solid #2475fc;
      .svg-action {
        color: #2475fc;
      }
      .right-block .title-text {
        color: #2475fc;
      }
    }
    .check-arrow {
      position: absolute;
      bottom: 0;
      right: -1px;
    }
  }
  .right-block {
    .title-block {
      display: flex;
      font-size: 14px;
      line-height: 22px;
      .title-text {
        margin-left: 2px;
        color: #262626;
        font-weight: 600;
      }
    }
    .desc {
      color: #595959;
      margin-top: 4px;
      line-height: 22px;
      word-break: break-all;
    }
  }
}
</style>

<template>
  <div class="add-library-page">
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
          <a-textarea v-model:value="formState.library_intro" placeholder="请输入知识库介绍" />
        </a-form-item>

        <a-form-item ref="name" label="知识库封面" v-bind="validateInfos.avatar">
          <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
          <div class="form-item-tip">请上传知识库封面，建议尺寸为100*100px.大小不超过100kb</div>
        </a-form-item>

        <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
          <a-button @click="resetForm">取 消</a-button>
          <a-button
            type="primary"
            style="margin-left: 16px"
            :loading="saveLoading"
            @click.prevent="onSubmit"
            >下一步</a-button
          >
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, computed } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { createLibrary } from '@/api/library/index'
import AvatarInput from './components/avatar-input.vue'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_OPEN_AVATAR } from '@/constants/index'
import { transformUrlData } from '@/utils/validate.js'
import { getLibraryInfo } from '@/api/library'

// 设置全局默认的duration为（2秒）
message.config({
  duration: 2
})

const router = useRouter()

const type = computed(() => {
  return 1
})

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

const resetForm = () => {
  router.back()
}

const onSubmit = () => {
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
  formData.append('is_default', 2)

  saveLoading.value = true

  createLibrary(formData)
    .then((res) => {
      message.success('创建成功')

      toHome(res.data.id)
    })
    .catch(() => {
      saveLoading.value = false
    })
}

const toHome = (id) => {
  getLibraryInfo({ id: id })
    .then((res) => {
      saveLoading.value = false
      // 对外知识库
      router.replace({
        path: '/public-library/config',
        query: {
          library_id: res.data.id,
          library_key: res.data.library_key
        }
      })
    })
    .catch(() => {
      saveLoading.value = false
    })
}

onMounted(() => {})
</script>
