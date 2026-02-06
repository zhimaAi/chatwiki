<style lang="less" scoped>
.form-box {
  margin-top: 32px;
  margin-bottom: 24px;
}

.form-item-tip {
  color: #999;
  color: #8c8c8c;
  margin-top: 8px;
  text-align: center;
}

.form-alert-tip {
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 400;
  line-height: 14px;
  margin: 2px 0 6px;
  white-space: nowrap;
}

.hight-set-text {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  justify-content: space-between;
  height: 56px;
  border-radius: 4px;
  background-color: #f2f4f7;
  padding: 0 16px;
  .title-block {
    display: flex;
    gap: 8px;
    color: #262626;
    font-size: 16px;
    font-weight: 600;
  }
  .right-btn-box {
    display: flex;
    gap: 4px;
    border-radius: 6px;
    font-size: 14px;
    padding: 0 8px;
    cursor: pointer;
    transition: all 0.3s ease-in-out;
    &:hover {
      background: var(--07, #e4e6eb);
      color: #2475fc;
    }
  }
}
.ml4 {
  margin-left: 4px;
}
</style>

<template>
  <div class="form-box">
    <a-form layout="vertical">
      <a-form-item ref="name" :label="false" v-bind="validateInfos.avatar">
        <AvatarInput v-model:value="formState.avatar" @change="onAvatarChange" />
        <div class="form-item-tip">{{ t('avatar_tip') }}</div>
      </a-form-item>
      <a-form-item ref="name" :label="t('library_name_label')" v-bind="validateInfos.library_name">
        <a-input
          v-model:value="formState.library_name"
          :placeholder="t('library_name_placeholder')"
          :maxlength="20"
        />
      </a-form-item>

      <a-form-item :label="t('library_intro_label')" v-bind="validateInfos.library_intro">
        <a-textarea
          :maxlength="1000"
          v-model:value="formState.library_intro"
          :placeholder="t('library_intro_placeholder')"
        />
      </a-form-item>

      <div class="hight-set-text">
        <div class="title-block"><SettingOutlined />{{ t('advanced_settings') }}</div>
        <div class="right-btn-box" @click="isHide = !isHide">
          <template v-if="isHide">
            {{ t('expand') }}
            <DownOutlined />
          </template>
          <template v-else>
            {{ t('collapse') }}
            <UpOutlined />
          </template>
        </div>
      </div>
      <a-form-item v-show="!isHide">
        <template #label
          >{{ t('document_vector_model') }}
          <a-tooltip>
            <template #title
              >{{ t('vector_model_tooltip') }}</template
            >
            <QuestionCircleOutlined class="ml4" />
          </a-tooltip>
        </template>
        <ModelSelect
          modelType="TEXT EMBEDDING"
          v-model:modeName="formState.use_model"
          v-model:modeId="formState.model_config_id"
          @loaded="onTextModelLoaded"
          style="width: 100%"
        />
      </a-form-item>
    </a-form>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { Form, message } from 'ant-design-vue'
import { createLibrary, getSeparatorsList } from '@/api/library/index'
import AvatarInput from './avatar-input.vue'
import { LIBRARY_QA_AVATAR } from '@/constants/index'
import { transformUrlData } from '@/utils/validate.js'
import {
  SettingOutlined,
  QuestionCircleOutlined,
  DownOutlined,
  UpOutlined
} from '@ant-design/icons-vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('views.ai-extract-faq.list.components.add-qa-library')
const emit = defineEmits('ok')

const isHide = ref(true)

const useForm = Form.useForm
const default_ai_chunk_prumpt = t('ai_chunk_prumpt_default')
const formState = reactive({
  type: '0',
  access_rights: 0,
  library_name: '',
  library_intro: '',
  use_model: '',
  model_config_id: '',
  library_files: undefined,
  avatar: '',
  avatar_file: undefined,
  is_offline: false,
  urls: '',
  doc_type: 1,
  file_name: '',
  is_qa_doc: 0, // 0 普通文档 1QA文档
  qa_index_type: 1, // 1问题与答案一起生成索引 2仅对问题生成索引
  doc_auto_renew_frequency: 1,
  chunk_type: 1,
  normal_chunk_default_separators_no: [12, 11],
  normal_chunk_default_not_merged_text: 0,
  normal_chunk_default_chunk_size: 512,
  normal_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_chunk_size: 512,
  semantic_chunk_default_chunk_overlap: 50,
  semantic_chunk_default_threshold: 90,
  graph_switch: false,
  graph_model_config_id: void 0,
  graph_use_model: '',
  ai_chunk_size: 5000, // ai大模型分段最大字符数
  ai_chunk_model: '', // ai大模型分段模型名称
  ai_chunk_model_config_id: '', // ai大模型分段模型配置id
  ai_chunk_prumpt: default_ai_chunk_prumpt // ai大模型分段提示词设置
})

const rules = reactive({
  library_name: [{ required: true, message: t('library_name_required'), trigger: 'change' }]
})

const { validate, validateInfos } = useForm(formState, rules)

const onAvatarChange = (data) => {
  formState.avatar = data.imageUrl
  formState.avatar_file = data.file
}

const segmentationTags = ref([])
getSeparatorsList().then((res) => {
  segmentationTags.value = res.data.map((item) => {
    return {
      label: item.name,
      value: item.no
    }
  })
})

const onTextModelLoaded = (list) => {
  if (list.length) {
    formState.use_model = list[0].children[0].name
    formState.model_config_id = list[0].model_config_id
  }
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
  formData.append('chunk_type', formState.chunk_type)
  formData.append(
    'normal_chunk_default_separators_no',
    formState.normal_chunk_default_separators_no.join(',')
  )
  formData.append('normal_chunk_default_chunk_size', formState.normal_chunk_default_chunk_size)
  formData.append('normal_chunk_default_not_merged_text', formState.normal_chunk_default_not_merged_text)
  formData.append(
    'normal_chunk_default_chunk_overlap',
    formState.normal_chunk_default_chunk_overlap
  )
  formData.append('semantic_chunk_default_chunk_size', formState.semantic_chunk_default_chunk_size)
  formData.append(
    'semantic_chunk_default_chunk_overlap',
    formState.semantic_chunk_default_chunk_overlap
  )
  formData.append('semantic_chunk_default_threshold', formState.semantic_chunk_default_threshold)
  formData.append('graph_switch', formState.graph_switch ? 1 : 0)
  formData.append('graph_model_config_id', formState.graph_model_config_id)
  formData.append('graph_use_model', formState.graph_use_model)

  formData.append('ai_chunk_size', formState.ai_chunk_size)
  formData.append('ai_chunk_model', formState.ai_chunk_model)
  formData.append('ai_chunk_model_config_id', formState.ai_chunk_model_config_id)
  formData.append('ai_chunk_prumpt', formState.ai_chunk_prumpt)
  formData.append('is_default', 2)

  createLibrary(formData).then((res) => {
    message.success(t('create_success'))
    // res.data.id
    emit('ok', res.data.id)
  })
}

const show = () => {
  formState.type = 2
  formState.avatar = LIBRARY_QA_AVATAR
  formState.avatar_file = LIBRARY_QA_AVATAR
  formState.library_name = ''
  formState.library_intro = ''
  formState.chunk_type = 1
  formState.normal_chunk_default_separators_no = [12, 11]
  formState.normal_chunk_default_not_merged_text = 0
  formState.normal_chunk_default_chunk_size = 512
  formState.normal_chunk_default_chunk_overlap = 50
  formState.semantic_chunk_default_chunk_size = 512
  formState.semantic_chunk_default_chunk_overlap = 50
  formState.semantic_chunk_default_threshold = 50
  formState.graph_switch = false

  formState.ai_chunk_size = 5000
  formState.ai_chunk_model = ''
  formState.ai_chunk_model_config_id = ''
  formState.ai_chunk_prumpt = default_ai_chunk_prumpt
}

defineExpose({
  show,
  handleOk
})

onMounted(() => {})
</script>
