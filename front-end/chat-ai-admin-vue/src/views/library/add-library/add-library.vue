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
        <a-form-item ref="name" :label="t('label_library_name')" v-bind="validateInfos.library_name">
          <a-input
            v-model:value="formState.library_name"
            type="text"
            :placeholder="t('ph_library_name')"
            :maxlength="20"
          />
        </a-form-item>

        <a-form-item :label="t('label_library_intro')">
          <a-textarea v-model:value="formState.library_intro" :placeholder="t('ph_library_intro')" />
        </a-form-item>

        <a-form-item ref="name" :label="t('label_library_avatar')" v-bind="validateInfos.avatar">
          <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
          <div class="form-item-tip">{{ t('msg_avatar_tip') }}</div>
        </a-form-item>
        <template v-if="type != 1">
          <a-form-item :label="t('label_doc_type')" required>
            <div class="upload-document-type-box">
              <div
                class="type-item"
                :class="{ active: formState.doc_type == item.value }"
                v-for="item in documentTypeList"
                :key="item.value"
                @click="handleChangeUrlType(item.value)"
              >
                <div class="right-block">
                  <div class="title-block">
                    <svg-icon
                      :name="formState.doc_type == item.value ? item.iconNameActive : item.iconName"
                    ></svg-icon>
                    <div class="title-text">{{ item.title }}</div>
                  </div>
                  <div class="desc">{{ item.desc }}</div>
                </div>
                <svg-icon
                  class="check-arrow"
                  name="check-arrow-filled"
                  style="font-size: 24px; color: #fff"
                  v-if="formState.doc_type == item.value"
                ></svg-icon>
              </div>
            </div>
          </a-form-item>

          <a-form-item
            v-show="formState.doc_type == 1"
            ref="name"
            :label="t('label_library_files')"
            v-bind="validateInfos.library_files"
          >
            <UploadFile :type="type" @change="handleFileChange" />
          </a-form-item>

          <a-form-item
            v-show="formState.doc_type == 2"
            ref="urls"
            :label="t('label_urls')"
            v-bind="validateInfos.urls"
          >
            <a-textarea
              style="height: 120px"
              v-model:value="formState.urls"
              :placeholder="t('ph_urls')"
            />
          </a-form-item>

          <a-form-item
            v-show="formState.doc_type == 2"
            ref="doc_auto_renew_frequency"
            :label="t('label_update_frequency')"
            required
          >
            <a-select v-model:value="formState.doc_auto_renew_frequency" style="width: 100%">
              <a-select-option :value="1">{{ t('option_no_auto_update') }}</a-select-option>
              <a-select-option :value="2">{{ t('option_every_day') }}</a-select-option>
              <a-select-option :value="3">{{ t('option_every_3_days') }}</a-select-option>
              <a-select-option :value="4">{{ t('option_every_7_days') }}</a-select-option>
              <a-select-option :value="5">{{ t('option_every_30_days') }}</a-select-option>
            </a-select>
          </a-form-item>

          <div v-show="formState.doc_type == 3">
            <a-form-item ref="file_name" :label="t('label_file_name')" v-bind="validateInfos.file_name">
              <a-input :placeholder="t('ph_file_name')" v-model:value="formState.file_name"></a-input>
            </a-form-item>
            <a-form-item :label="t('label_doc_category')" required v-if="false">
              <a-radio-group v-model:value="formState.is_qa_doc">
                <a-radio :value="0">{{ t('option_normal_doc') }}</a-radio>
                <a-radio :value="1">{{ t('option_qa_doc') }}</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item :label="t('label_index_method')" required v-if="formState.is_qa_doc == 1">
              <div class="upload-document-type-box">
                <div
                  class="type-item"
                  :class="{ active: formState.qa_index_type == item.value }"
                  v-for="item in qaIndexTypeList"
                  :key="item.value"
                  @click="handleChangeQaIndexType(item.value)"
                >
                  <div class="right-block">
                    <div class="title-block">
                      <svg-icon
                        :name="
                          formState.qa_index_type == item.value
                            ? item.iconNameActive
                            : item.iconName
                        "
                      ></svg-icon>
                      <div class="title-text">{{ item.title }}</div>
                    </div>
                    <div class="desc">{{ item.desc }}</div>
                  </div>
                  <svg-icon
                    class="check-arrow"
                    name="check-arrow-filled"
                    style="font-size: 24px; color: #fff"
                    v-if="formState.qa_index_type == item.value"
                  ></svg-icon>
                </div>
              </div>
            </a-form-item>
          </div>
        </template>

        <a-form-item :label="t('label_embedding_model')" v-bind="validateInfos.use_model">
          <div class="card-box" v-if="false">
            <div
              class="use-model-item"
              :class="{ active: isActive == item.value }"
              v-for="item in libraryModeList"
              :key="item.value"
              @click="handleSelectLibrary(item)"
            >
              <div class="use-model-item-top" :class="{ active: isActive == item.value }">
                <svg-icon
                  class="use-model-item-top-icon"
                  style="color: red"
                  :name="isActive == item.value ? item.iconNameActive : item.iconName"
                ></svg-icon>
                <p>{{ item.title }}</p>
              </div>
              <p>{{ item.desc }}</p>
              <svg-icon
                class="check-arrow"
                name="check-arrow-filled"
                style="font-size: 24px; color: #fff"
                v-if="isActive == item.value"
              ></svg-icon>
            </div>
          </div>
          <a-select
            v-model:value="formState.use_model"
            :placeholder="t('ph_select_model')"
            @change="handleChangeModel"
          >
            <a-select-opt-group v-for="item in modelList" :key="item.id">
              <template #label>
                <a-flex align="center" :gap="6">
                  <img class="model-icon" :src="item.icon" alt="" />{{ item.name }}
                </a-flex>
              </template>
              <a-select-option
                :value="
                  modelDefine.indexOf(item.model_define) > -1 && val.deployment_name
                    ? val.deployment_name
                    : val.name + val.id
                "
                :model_config_id="item.id"
                :current_obj="val"
                v-for="val in item.children"
                :key="val.name + val.id"
              >
                <span v-if="modelDefine.indexOf(item.model_define) > -1 && val.deployment_name">{{
                  val.deployment_name
                }}</span>
                <span v-else>{{ val.name }}</span>
              </a-select-option>
            </a-select-opt-group>
          </a-select>
        </a-form-item>

        <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
          <a-button @click="resetForm">{{ t('btn_cancel') }}</a-button>
          <a-button
            type="primary"
            style="margin-left: 16px"
            :loading="saveLoading"
            @click.prevent="onSubmit"
            >{{ t('btn_next') }}</a-button
          >
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted, computed } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useRouter, useRoute } from 'vue-router'
import UploadFile from './components/upload-input.vue'
import { createLibrary } from '@/api/library/index'
import { getModelConfigOption } from '@/api/model/index'
import AvatarInput from './components/avatar-input.vue'
import { LIBRARY_NORMAL_AVATAR, LIBRARY_OPEN_AVATAR } from '@/constants/index'
import { transformUrlData } from '@/utils/validate.js'
import { useStorage } from '@/hooks/web/useStorage'
import { duplicateRemoval, removeRepeat } from '@/utils/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.add-library.add-library')

const { getStorage, setStorage } = useStorage('localStorage')

// 设置全局默认的duration为（2秒）
message.config({
  duration: 2
})

const router = useRouter()
const route = useRoute()

const type = computed(() => {
  return route.query.type || 0
})

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']

const libraryModeList = computed(() => [
  {
    iconName: 'high',
    iconNameActive: 'high-active',
    title: t('title_high_quality'),
    value: 1,
    is_offline: false,
    desc: t('desc_high_quality')
  }
])

const documentTypeList = computed(() => {
  const list = [
    {
      iconName: 'doc-icon',
      iconNameActive: 'doc-icon-active',
      title: t('title_local_doc'),
      value: 1,
      desc: type.value == 2 ? t('desc_local_doc_open') : t('desc_local_doc_normal')
    },
    {
      iconName: 'link-icon',
      iconNameActive: 'link-icon-active',
      title: t('title_online_data'),
      value: 2,
      desc: t('desc_online_data')
    },
    {
      iconName: 'cu-doc-icon',
      iconNameActive: 'cu-doc-active',
      title: t('title_custom_doc'),
      value: 3,
      desc: t('desc_custom_doc')
    }
  ]
  if (type.value == 2) {
    return list.filter((item) => item.value != 2)
  }
  return list
})

const qaIndexTypeList = computed(() => [
  {
    iconName: 'file-search',
    iconNameActive: 'file-search',
    title: t('title_qa_index_both'),
    value: 1,
    desc: t('desc_qa_index_both')
  },
  {
    iconName: 'comment-search',
    iconNameActive: 'comment-search',
    title: t('title_qa_index_question_only'),
    value: 2,
    desc: t('desc_qa_index_question_only')
  }
])

const useForm = Form.useForm
const saveLoading = ref(false)

const lastEmbeddedModel = computed(() => getStorage('lastEmbeddedModel') || {})
const isActive = ref(1)
const defaultAvatar = type.value == 0 ? LIBRARY_NORMAL_AVATAR : LIBRARY_OPEN_AVATAR
const formState = reactive({
  type: type.value,
  access_rights: 0,
  library_name: '',
  library_intro: '',
  use_model: 'text-embedding-v2',
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
const currentModelDefine = ref('')

const validateFiles = (_rule, value) => {
  if (type.value == 1) {
    return Promise.resolve()
  }

  if ((value && value.length > 0) || formState.doc_type != 1) {
    return Promise.resolve()
  } else {
    return Promise.reject(new Error(t('msg_upload_file')))
  }
}

const validateUrl = (_rule, value) => {
  if (type.value == 1) {
    return Promise.resolve()
  }

  if (formState.doc_type != 2) {
    return Promise.resolve()
  }
  if (transformUrlData(value) === false) {
    return Promise.reject(new Error(t('msg_invalid_url')))
  }
  return Promise.resolve()
}

const rules = reactive({
  library_name: [{ required: true, message: t('msg_enter_library_name'), trigger: 'change' }],
  use_model: [{ required: true, message: t('msg_select_model'), trigger: 'change' }],
  library_files: [
    {
      required: type.value == 0,
      message: t('msg_select_file'),
      trigger: 'change',
      validator: validateFiles
    }
  ],
  urls: [
    {
      required: type.value == 0,
      validator: validateUrl
    }
  ],
  file_name: [
    {
      required: type.value == 0,
      validator: (_rule, value) => {
        if (type.value == 1) {
          return Promise.resolve()
        }
        if (formState.doc_type != 3 || value) {
          return Promise.resolve()
        } else {
          return Promise.reject(new Error(t('msg_enter_doc_name')))
        }
      }
    }
  ]
})

const { validate, validateInfos } = useForm(formState, rules)

const handleFileChange = (fileList) => {
  formState.library_files = fileList
}

const handleChangeModel = (val, option) => {
  const self = option.current_obj
  formState.use_model =
    modelDefine.indexOf(self.model_define) > -1 && self.deployment_name
      ? self.deployment_name
      : self.name
  currentModelDefine.value = self.model_define
  formState.model_config_id = self.id || option.model_config_id
  console.log(formState, '==')
}

const onAvatarChange = (data) => {
  formState.avatar = data.imageUrl
  formState.avatar_file = data.file
}

const handleSelectLibrary = (item) => {
  formState.is_offline = item.is_offline
  isActive.value = item.value
  formState.use_model = undefined
  getModelList(item.is_offline)
}

const handleChangeUrlType = (type) => {
  formState.doc_type = type
}

const handleChangeQaIndexType = (type) => {
  formState.qa_index_type = type
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

  if (oldModelDefineList.indexOf(currentModelDefine.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    newFormState.use_model = '默认'
  }

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

  // “嵌入模型”记住用户上次选择
  setStorage('lastEmbeddedModel', {
    // model_config_id: formState.model_config_id,
    // is_offline: formState.is_offline,
    // use_model: formState.use_model,
    // isActive: isActive.value
  })
  let isTableType = false

  if (type.value == 0 || type.value == 2) {
    if (formState.doc_type == 1) {
      formState.library_files.forEach((file) => {
        if (file.name.includes('.xlsx') || file.name.includes('.csv')) {
          isTableType = true
        }
        formData.append('library_files', file)
      })
    }
  }

  saveLoading.value = true

  createLibrary(formData)
    .then((res) => {
      message.success(t('msg_create_success'))
      // 对外知识库
      if (type.value == 1) {
        router.replace({
          path: '/public-library/config',
          query: {
            library_id: res.data.id
          }
        })

        saveLoading.value = false

        return
      }

      let path = '/library/details/knowledge-document'
      let query = {
        id: res.data.id
      }

      if (isTableType) {
        path = '/library/document-segmentation'
        query = {
          document_id: res.data.file_ids[0]
        }
      }
      if (formState.doc_type == 3) {
        path = '/library/preview'
        query = {
          id: res.data.file_ids[0]
        }
      }
      router.replace({
        path,
        query
      })
      saveLoading.value = false
    })
    .catch(() => {
      saveLoading.value = false
    })
}

// 获取嵌入模型列表
const modelList = ref([])

function uniqueArr(arr, arr1, key) {
  const keyVals = new Set(arr.map((item) => item.model_define))
  arr1.filter((obj) => {
    let val = obj[key]
    if (keyVals.has(val)) {
      arr.filter((obj1) => {
        if (obj1.model_define == val) {
          obj1.children = removeRepeat(obj1.children, obj.children)
          return false
        }
      })
    }
  })
  return arr
}

const getModelList = (is_offline) => {
  getModelConfigOption({
    model_type: 'TEXT EMBEDDING',
    is_offline
  }).then((res) => {
    let list = res.data || []
    let children = []
    // let isCheckId = false
    // 没有模型选项则不用缓存中的记录上传模型选择
    // if (!list.length && lastEmbeddedModel.value.model_config_id) {
    //   formState.use_model = undefined
    //   removeStorage('lastEmbeddedModel')
    // }

    modelList.value = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.vector_model_list.length; i++) {
        const ele = item.model_info.vector_model_list[i]
        children.push({
          name: ele,
          deployment_name: item.model_config.deployment_name,
          id: item.model_config.id,
          model_define: item.model_info.model_define
        })
      }
      if (item.model_config.id == lastEmbeddedModel.value.model_config_id) {
        // 有缓存的id
        // isCheckId = true
      }
      return {
        id: item.model_config.id,
        name: item.model_info.model_name,
        model_define: item.model_info.model_define,
        icon: item.model_info.model_icon_url,
        children: children,
        deployment_name: item.model_config.deployment_name
      }
    })

    // if (!isCheckId) {
    //   // 如果获取的模型列表没有缓存的这个,则清除缓存的模型
    //   formState.use_model = undefined
    //   removeStorage('lastEmbeddedModel')
    // }

    // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
    modelList.value = uniqueArr(
      duplicateRemoval(modelList.value, 'model_define'),
      modelList.value,
      'model_define'
    )
    modelList.value.forEach((item) => {
      if (item.model_define == 'chatwiki') {
        formState.model_config_id = item.id
      }
    })
  })
}

onMounted(() => {
  getModelList(false)
})
</script>
