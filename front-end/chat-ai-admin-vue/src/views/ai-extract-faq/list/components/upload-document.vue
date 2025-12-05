<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      title="上传文档"
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
        <a-form-item name="chunk_model_config_id" label="提取模型">
          <ModelSelect
            modelType="LLM"
            placeholder="请选择AI大模型"
            v-model:modeName="formState.chunk_model"
            v-model:modeId="formState.chunk_model_config_id"
            :modeName="formState.chunk_model"
            :modeId="formState.chunk_model_config_id"
            style="width: 100%"
            @loaded="onVectorModelLoaded"
          />
        </a-form-item>
        <a-form-item label="分块方式">
          <a-radio-group v-model:value="formState.chunk_type">
            <a-radio :value="1">按长度</a-radio>
            <a-radio :value="2">按分隔符</a-radio>
          </a-radio-group>
          <div class="mt8" v-if="formState.chunk_type == 1">
            <a-input-number
              v-model:value="formState.chunk_size"
              :min="1"
              :max="10000"
              :step="1"
              placeholder="请输入"
              style="width: 300px"
            />
          </div>
          <div v-if="formState.chunk_type == 2" class="mt8">
            <a-select
              placeholder="请选择"
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
        <a-form-item label="最大长度" v-if="formState.chunk_type == 2">
            <a-input-number
              v-model:value="formState.chunk_size"
              :min="1"
              :max="10000"
              :step="1"
              placeholder="请输入"
              style="width: 300px"
            />
        </a-form-item>
        <a-form-item label="FAQ提取提示词">
          <a-textarea
            v-model:value="formState.chunk_prompt"
            style="height: 150px"
            placeholder="请输入"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick, onMounted } from 'vue'
import {addFAQFile, getFAQConfig } from '@/api/library/index'
import UploadFilesInput from './upload-input.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { formatSeparatorsNo } from '@/utils/index'
import { message } from 'ant-design-vue'
const emit = defineEmits(['ok'])

const props = defineProps({
  separatorsOptions: {
    type: Array,
    default: () => []
  },
})

const defaultPromt = `根据user角色提供的文本，学习和分析它，并整理学习成果：
- 提出问题并给出每个问题的答案。
- 答案需详细完整，尽可能保留原文描述，可以适当扩展答案描述。
- 答案可以包含普通文字、链接、代码、表格、公示、媒体链接等 Markdown 元素。
- 最多提出 50 个问题。
- 生成的问题和答案和源文本语言相同。`

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
      message: '请选择提取模型',
      required: true
    }
  ],
  faq_files: [
    {
      message: '请选择文件',
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
      message.success('上传成功')
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
