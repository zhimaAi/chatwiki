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
      box-shadow: 0px 5px 5px -3px rgba(0, 0, 0, 0.1), 0px 8px 10px 1px rgba(0, 0, 0, 0.06),
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
            placeholder="请输入知识库名称，最多20个字"
          />
        </a-form-item>

        <a-form-item label="知识库简介" v-bind="validateInfos.library_intro">
          <a-textarea v-model:value="formState.library_intro" placeholder="请输入知识库介绍" />
        </a-form-item>

        <a-form-item ref="name" label="知识库封面" v-bind="validateInfos.robot_avatar_url">
          <AvatarInput v-model:value="formState.robot_avatar_url" @change="onAvatarChange" />
          <div class="form-item-tip">请上传知识库封面，建议尺寸为100*100px.大小不超过100kb</div>
        </a-form-item>

        <a-form-item label="导入文档类型" required>
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
          label="知识库文档"
          v-bind="validateInfos.library_files"
        >
          <UploadFile @change="handleFileChange" />
        </a-form-item>

        <a-form-item
          v-show="formState.doc_type == 2"
          ref="urls"
          label="网页链接"
          v-bind="validateInfos.urls"
        >
          <a-textarea
            style="height: 120px"
            v-model:value="formState.urls"
            placeholder="请输入网页链接,形式：一行标题一行网页链接"
          />
        </a-form-item>

        <a-form-item
          v-show="formState.doc_type == 2"
          ref="doc_auto_renew_frequency"
          label="更新频率"
          required
        >
          <a-select v-model:value="formState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">不自动更新</a-select-option>
            <a-select-option :value="2">每天</a-select-option>
            <a-select-option :value="3">每3天</a-select-option>
            <a-select-option :value="4">每7天</a-select-option>
            <a-select-option :value="5">每30天</a-select-option>
          </a-select>
        </a-form-item>

        <div v-show="formState.doc_type == 3">
          <a-form-item ref="file_name" label="文档名称" v-bind="validateInfos.file_name">
            <a-input placeholder="请输入文档名称" v-model:value="formState.file_name"></a-input>
          </a-form-item>
          <a-form-item label="文档类型" required>
            <a-radio-group v-model:value="formState.is_qa_doc">
              <a-radio :value="0">普通文档</a-radio>
              <a-radio :value="1">QA文档</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="索引方式" required v-if="formState.is_qa_doc == 1">
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
                        formState.qa_index_type == item.value ? item.iconNameActive : item.iconName
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

        <a-form-item label="嵌入模型 " v-bind="validateInfos.use_model">
          <div class="card-box">
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
            placeholder="请选择嵌入模型"
            @change="handleChangeModel"
          >
            <a-select-opt-group v-for="item in modelList" :key="item.id">
              <template #label>
                <a-flex align="center" :gap="6">
                  <img class="model-icon" :src="item.icon" alt="" />{{ item.name }}
                </a-flex>
              </template>
              <a-select-option
                :value="modelDefine.indexOf(item.model_define) > -1 && val.deployment_name ? val.deployment_name : val.name + val.id"
                :model_config_id="item.id"
                :current_obj="val"
                v-for="val in item.children"
                :key="val.name + val.id"
              >
                <span v-if="modelDefine.indexOf(item.model_define) > -1 && val.deployment_name">{{val.deployment_name}}</span>
                <span v-else>{{ val.name }}</span>
              </a-select-option>
            </a-select-opt-group>
          </a-select>
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
import UploadFile from './components/upload-input.vue'
import { createLibrary } from '@/api/library/index'
import { getModelConfigOption } from '@/api/model/index'
import AvatarInput from './avatar-input.vue'
import { DEFAULT_LIBRARY_AVATAR } from '@/constants/index'
import { transformUrlData } from '@/utils/validate.js'
import { useStorage } from '@/hooks/web/useStorage'
import { duplicateRemoval, removeRepeat } from '@/utils/index'

const { getStorage, setStorage, removeStorage } = useStorage('localStorage')

// 设置全局默认的duration为（2秒）
message.config({
  duration: 2,
});

const router = useRouter()

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']

const libraryModeList = ref([
  {
    iconName: 'high',
    iconNameActive: 'high-active',
    title: '高质量',
    value: 1,
    is_offline: false,
    desc: '使用在线的嵌入模型，在召回时具有更高的准确度，但需要花费token'
  },
  {
    iconName: 'economic',
    iconNameActive: 'economic-active',
    title: '经济',
    value: 2,
    is_offline: true,
    desc: '使用离线的向量模型，较在线模型准确度稍低，但是不需要消耗token'
  }
])

const documentTypeList = ref([
  {
    iconName: 'doc-icon',
    iconNameActive: 'doc-icon-active',
    title: '本地文档',
    value: 1,
    desc: '上传本地 text/pdf/doc/docx/xlsx 等格式文件'
  },
  {
    iconName: 'link-icon',
    iconNameActive: 'link-icon-active',
    title: '在线数据',
    value: 2,
    desc: '获取在线网页内容'
  },
  {
    iconName: 'cu-doc-icon',
    iconNameActive: 'cu-doc-active',
    title: '自定义文档',
    value: 3,
    desc: '自定义一个空文档，手动添加或编辑内容'
  }
])
const qaIndexTypeList = ref([
  {
    iconName: 'file-search',
    iconNameActive: 'file-search',
    title: '问题与答案一起生成索引',
    value: 1,
    desc: '回答用户提问时，将用户提问与导入的问题和答案一起对比相似度，根据相似度高的问题和答案回复'
  },
  {
    iconName: 'comment-search',
    iconNameActive: 'comment-search',
    title: '仅对问题生成索引',
    value: 2,
    desc: '回答用户提问时，将用户提问与导入的问题一起对比相似度，再根据相似度高的问题和对应的答案来回复'
  }
])

const useForm = Form.useForm
const saveLoading = ref(false)

const lastEmbeddedModel = computed(() => getStorage('lastEmbeddedModel') || {})
const isActive = ref(lastEmbeddedModel.value.isActive ? lastEmbeddedModel.value.isActive : 2)
const formState = reactive({
  library_name: '',
  library_intro: '',
  use_model: lastEmbeddedModel.value.use_model ? lastEmbeddedModel.value.use_model : undefined,
  model_config_id: lastEmbeddedModel.value.model_config_id ? lastEmbeddedModel.value.model_config_id : '',
  library_files: undefined,
  avatar: DEFAULT_LIBRARY_AVATAR,
  robot_avatar_url: DEFAULT_LIBRARY_AVATAR,
  is_offline: Object.prototype.hasOwnProperty.call(lastEmbeddedModel.value, 'is_offline') ? lastEmbeddedModel.value.is_offline : false,
  urls: '',
  doc_type: 1,
  file_name: '',
  is_qa_doc: 0, // 0 普通文档 1QA文档
  qa_index_type: 1, // 1问题与答案一起生成索引 2仅对问题生成索引
  doc_auto_renew_frequency: 1
})
const currentModelDefine = ref('')

const validateFiles = (_rule, value) => {
  if ((value && value.length > 0) || formState.doc_type != 1) {
    return Promise.resolve()
  } else {
    return Promise.reject(new Error('请上传文件'))
  }
}

const validateUrl = (_rule, value) => {
  if (formState.doc_type != 2) {
    return Promise.resolve()
  }
  if(transformUrlData(value) === false){
    return Promise.reject(new Error('网页地址不合法'))
  }
  return Promise.resolve()
}

const rules = reactive({
  library_name: [{ required: true, message: '请输入库名称', trigger: 'change' }],
  library_intro: [{ required: true, message: '请输入库简介', trigger: 'change' }],
  use_model: [{ required: true, message: '请选择嵌入模型', trigger: 'change' }],
  library_files: [
    { required: true, message: '请选择文件', trigger: 'change', validator: validateFiles }
  ],
  urls: [
    {
      required: true,
      validator: validateUrl
    }
  ],
  file_name: [
    {
      required: true,
      validator: (_rule, value) => {
        if (formState.doc_type != 3 || value) {
          return Promise.resolve()
        } else {
          return Promise.reject(new Error('请输入文档名称'))
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
  formState.use_model = modelDefine.indexOf(self.model_define) > -1 && self.deployment_name ? self.deployment_name : self.name
  currentModelDefine.value = self.model_define
  formState.model_config_id = self.id || option.model_config_id
}

const onAvatarChange = (data) => {
  formState.avatar = data.file
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

  let newFormState = JSON.parse(JSON.stringify(formState)); // 深拷贝，不能改变原对象

  if (oldModelDefineList.indexOf(currentModelDefine.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    newFormState.use_model = '默认'
  }

  formData.append('library_name', formState.library_name)
  formData.append('library_intro', formState.library_intro)
  formData.append('use_model', newFormState.use_model)
  formData.append('model_config_id', formState.model_config_id)
  formData.append('avatar', formState.avatar || DEFAULT_LIBRARY_AVATAR)
  formData.append('is_offline', formState.is_offline)
  formData.append('urls', JSON.stringify(transformUrlData(formState.urls)))
  formData.append('doc_type', formState.doc_type)
  formData.append('file_name', formState.file_name)
  formData.append('is_qa_doc', formState.is_qa_doc)
  formData.append('qa_index_type', formState.qa_index_type)
  formData.append('doc_auto_renew_frequency', formState.doc_auto_renew_frequency)

  // “嵌入模型”记住用户上次选择
  setStorage('lastEmbeddedModel', {
    model_config_id: formState.model_config_id,
    is_offline: formState.is_offline,
    use_model: formState.use_model,
    isActive: isActive.value
  })
  let isTableType = false
  if (formState.doc_type == 1) {
    formState.library_files.forEach((file) => {
      if (file.name.includes('.xlsx') || file.name.includes('.csv')) {
      isTableType = true
      }
      formData.append('library_files', file)
    })
  }

  saveLoading.value = true

  createLibrary(formData)
    .then((res) => {
      message.success('创建成功')
      let path = '/library/details/knowledge-document'
      let query = {
        id: res.data.id
      }
      if(isTableType){
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
    const keyVals = new Set(arr.map(item => item.model_define));
    arr1.filter(obj => {
        let val = obj[key];
        if (keyVals.has(val)) {
          arr.filter(obj1 => {
            if (obj1.model_define == val) {
              obj1.children = removeRepeat(obj1.children, obj.children)
              return false
            }
          })
        }
    });
    return arr
}

const getModelList = (is_offline) => {
  getModelConfigOption({
    model_type: 'TEXT EMBEDDING',
    is_offline
  }).then((res) => {
    let list = res.data || []
    let children = []
    let isCheckId = false
    // 没有模型选项则不用缓存中的记录上传模型选择
    if (!list.length && lastEmbeddedModel.value.model_config_id) {
      formState.use_model = undefined
      removeStorage('lastEmbeddedModel')
    }

    modelList.value = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.vector_model_list.length; i++) {
        const ele = item.model_info.vector_model_list[i];
        children.push({
          name: ele,
          deployment_name: item.model_config.deployment_name,
          id: item.model_config.id,
          model_define: item.model_info.model_define
        })
      }
      if (item.model_config.id == lastEmbeddedModel.value.model_config_id) {
        // 有缓存的id
        isCheckId = true
      }
      return {
        id: item.model_config.id,
        name: item.model_info.model_name,
        model_define: item.model_info.model_define,
        icon: item.model_info.model_icon_url,
        children: children,
        deployment_name: item.model_config.deployment_name,
      }
    })

    if (!isCheckId) {
      // 如果获取的模型列表没有缓存的这个,则清除缓存的模型
      formState.use_model = undefined
      removeStorage('lastEmbeddedModel')
    }
    
    // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
    modelList.value = uniqueArr(duplicateRemoval(modelList.value, 'model_define'), modelList.value, 'model_define')
  })
}

onMounted(() => {
  if (Object.prototype.hasOwnProperty.call(lastEmbeddedModel.value, 'is_offline')) {
    getModelList(lastEmbeddedModel.value.is_offline)
  } else {
    getModelList(true)
  }
})

</script>
