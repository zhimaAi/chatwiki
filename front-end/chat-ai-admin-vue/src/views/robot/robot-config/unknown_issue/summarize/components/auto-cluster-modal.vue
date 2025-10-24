<template>
  <div>
    <a-modal v-model:open="open" title="自动聚类设置" @ok="handleOk">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item v-bind="validateInfos.unknown_summary_model_config_id">
            <template #label>请设置使用的向量模型 </template>
            <ModelSelect
              modelType="TEXT EMBEDDING"
              v-model:modeName="formState.unknown_summary_use_model"
              v-model:modeId="formState.unknown_summary_model_config_id"
              @loaded="onTextModelLoaded"
              style="width: 100%"
            />
          </a-form-item>
          <a-form-item v-bind="validateInfos.unknown_summary_similarity">
            <template #label
              >相似度阈值
              <a-tooltip>
                <template #title>语意相似度大于指定阈值的问题会自动聚类</template>
                <QuestionCircleOutlined class="ml4" />
              </a-tooltip>
            </template>
            <a-input-number
              placeholder="请输入相似度阈值"
              style="width: 100%"
              v-model:value="formState.unknown_summary_similarity"
              :min="0"
              :max="1"
              :precision="2"
              :step="0.1"
            />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { setUnknownIssueSummary } from '@/api/robot/index.js'
import ModelSelect from '@/components/model-select/model-select.vue'

const emit = defineEmits(['ok'])
const query = useRoute().query

const open = ref(false)
const formState = reactive({
  unknown_summary_model_config_id: '',
  unknown_summary_use_model: '',
  unknown_summary_similarity: '',
  unknown_summary_status: '',
  id: query.id
})

const useForm = Form.useForm
const saveLoading = ref(false)
const show = (data) => {
  console.log(data, '===')
  formState.unknown_summary_model_config_id = data.unknown_summary_model_config_id || ''
  formState.unknown_summary_use_model = data.unknown_summary_use_model || ''
  formState.unknown_summary_similarity = data.unknown_summary_similarity || ''
  formState.unknown_summary_status = data.unknown_summary_status || ''
  console.log(formState, '==')
  open.value = true
}

const rules = reactive({
  unknown_summary_similarity: [{ required: true, message: '请输入相似度阈值', trigger: 'change' }],
  unknown_summary_model_config_id: [
    { required: true, message: '请选择使用的向量模型', trigger: 'change' }
  ]
})

const { validate, validateInfos } = useForm(formState, rules)

const onTextModelLoaded = (list) => {
  if (list.length) {
    // formState.use_model = list[0].children[0].name
    // formState.model_config_id = list[0].model_config_id
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
  setUnknownIssueSummary({
    ...formState
  }).then((res) => {
    message.success('保存成功')
    open.value = false
    emit('ok')
  })
}
defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-box {
  margin: 24px 0;
}
.ml4 {
  margin-left: 4px;
}
</style>
