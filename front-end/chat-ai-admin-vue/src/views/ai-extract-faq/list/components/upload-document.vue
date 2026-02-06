<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :title="t('modal_title')"
      width="746px"
      @ok="handleSave"
    >
      <a-form
        class="url-add-form"
        layout="vertical"
        ref="formRef"
        :model="formState"
        :rules="rules"
      >
        <a-form-item name="faq_files" :label="null">
          <UploadFilesInput
            :maxCount="10"
            v-model:value="formState.faq_files"
            @change="onFilesChange"
          />
        </a-form-item>
        <a-form-item name="chunk_model_config_id" :label="t('chunk_model_label')">
          <ModelSelect
            modelType="LLM"
            :placeholder="t('chunk_model_placeholder')"
            v-model:modeName="formState.chunk_model"
            v-model:modeId="formState.chunk_model_config_id"
            :modeName="formState.chunk_model"
            :modeId="formState.chunk_model_config_id"
            style="width: 100%"
            @loaded="onVectorModelLoaded"
          />
        </a-form-item>
        <a-form-item :label="t('chunk_type_label')">
          <a-radio-group v-model:value="formState.chunk_type">
            <a-radio :value="1">{{ t('chunk_type_by_length') }}</a-radio>
            <a-radio :value="2">{{ t('chunk_type_by_separator') }}</a-radio>
          </a-radio-group>
          <div class="mt8" v-if="formState.chunk_type == 1">
            <a-input-number
              v-model:value="formState.chunk_size"
              :min="1"
              :max="10000"
              :step="1"
              :placeholder="t('chunk_size_placeholder')"
              style="width: 300px"
            />
          </div>
          <div v-if="formState.chunk_type == 2" class="mt8">
            <a-select
              :placeholder="t('separator_placeholder')"
              style="width: 300px"
              mode="tags"
              v-model:value="formState.separators_no"
            >
              <a-select-option :value="item.no" v-for="item in props.separatorsOptions" :key="item.no">{{
                item.name
              }}</a-select-option>
            </a-select>
          </div>
        </a-form-item>
        <a-form-item :label="t('max_length_label')" v-if="formState.chunk_type == 2">
            <a-input-number
              v-model:value="formState.chunk_size"
              :min="1"
              :max="10000"
              :step="1"
              :placeholder="t('chunk_size_placeholder')"
              style="width: 300px"
            />
        </a-form-item>
        <a-form-item :label="t('faq_prompt_label')">
          <a-textarea
            v-model:value="formState.chunk_prompt"
            style="height: 150px"
            :placeholder="t('faq_prompt_placeholder')"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, nextTick, onMounted } from 'vue'
import {addFAQFile, getFAQConfig } from '@/api/library/index'
import UploadFilesInput from './upload-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { formatSeparatorsNo } from '@/utils/index'
import { message } from 'ant-design-vue'


const { t } = useI18n('views.ai-extract-faq.list.components.upload-document')

const emit = defineEmits(['ok'])

const defaultPromt = t('default_prompt')

const props = defineProps({
  separatorsOptions: {
    type: Array,
    default: () => []
  },
})



const formState = reactive({
  faq_files: [],
  chunk_model: '',
  chunk_model_config_id: '',
  chunk_type: 1,
  chunk_size: 8000,
  chunk_prompt: defaultPromt,
  separators_no: [12, 11]
})
const rules = ref({
  chunk_model_config_id: [
    {
      message: t('confirm_validate_msg_model'),
      required: true
    }
  ],
  faq_files: [
    {
      message: t('confirm_validate_msg_file'),
      required: true
    }
  ]
})

const open = ref(false)
const show = () => {
  open.value = true
  getFAQConfig().then((res) => {
    let data = res.data || {
      faq_files: [],
      chunk_model: '',
      chunk_model_config_id: '',
      chunk_type: 1,
      chunk_size: 8000,
      chunk_prompt: defaultPromt,
      separators_no: [12, 11]
    }
    formState.faq_files = []
    formState.chunk_model = data.chunk_model
    formState.chunk_model_config_id = data.chunk_model_config_id
    formState.chunk_type = 1
    formState.chunk_size = data.chunk_size || 8000
    formState.chunk_prompt = data.chunk_prompt || defaultPromt
    formState.separators_no = formatSeparatorsNo(data.separators_no, [12, 11])
    setTimeout(() => {
      if (!formState.chunk_model_config_id && vectorModelList.value.length > 0) {
        if (vectorModelList.value[0].children && vectorModelList.value[0].children.length) {
          formState.chunk_model_config_id = vectorModelList.value[0].model_config_id
          formState.chunk_model = vectorModelList.value[0].children[0].name
        }
      }
    }, 500)
  })
}
const confirmLoading = ref(false)

const formRef = ref(null)
const handleSave = () => {
  formRef.value
    .validate()
    .then(() => {
      saveForm()
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

const saveForm = () => {
  let formData = new FormData()
  formState.faq_files.forEach((file) => {
    formData.append('faq_files', file)
  })
  formData.append('chunk_type', formState.chunk_type)
  formData.append('chunk_size', formState.chunk_size)
  formData.append('chunk_model', formState.chunk_model)
  formData.append('chunk_model_config_id', formState.chunk_model_config_id)
  formData.append('chunk_prompt', formState.chunk_prompt)
  formData.append('separators_no', JSON.stringify(formState.separators_no))
  confirmLoading.value = true
  addFAQFile(formData)
    .then((res) => {
      message.success(t('upload_success_msg'))
      emit('ok')
      open.value = false
    })
    .finally(() => {
      confirmLoading.value = false
    })
}


const onFilesChange = (files) => {
  formState.faq_files = files
  formRef.value.validate(['faq_files'])
}
const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  nextTick(() => {})
  // handleEdit()
}



defineExpose({
  show
})
</script>

<style lang="less" scoped>
.url-add-form {
  margin: 24px 0;
}
.mt8 {
  margin-top: 8px;
}
</style>
