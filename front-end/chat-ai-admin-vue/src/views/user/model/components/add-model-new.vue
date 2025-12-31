<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="`${formState.id ? '编辑' : '添加'}模型`"
      :width="700"
      @ok="handleOk"
      :confirmLoading="isLoading"
    >
      <a-form class="form-box" ref="formRef" layout="vertical" :model="formState">
        <a-form-item name="model_type" label="模型类型" required>
          <a-radio-group
            @change="handleTypeChange"
            v-model:value="formState.model_type"
            name="radioGroup"
          >
            <a-radio v-if="supported_type.includes('LLM')" value="LLM">大语言模型</a-radio>
            <a-radio v-if="supported_type.includes('TEXT EMBEDDING')" value="TEXT EMBEDDING"
              >嵌入模型</a-radio
            >
            <a-radio v-if="supported_type.includes('RERANK')" value="RERANK">重排序模型</a-radio>
            <a-radio v-if="supported_type.includes('TTS')" value="TTS">语音模型</a-radio>
            <a-radio v-if="supported_type.includes('IMAGE')" value="IMAGE">图片生成模型</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item
          name="use_model_name"
          :label="model_name"
          :rules="[{ required: true, message: `请输入${model_name}` }]"
        >
          <a-input
            v-model:value="formState.use_model_name"
            @blur="handleModelNameBlur"
            placeholder="请输入，为调用模型时参数"
          />
        </a-form-item>
        <a-form-item
          name="show_model_name"
          label="模型名称"
          :rules="[{ required: true, message: '请输入模型名称' }]"
        >
          <a-input
            v-model:value="formState.show_model_name"
            :maxLength="50"
            placeholder="请输入模型名称，仅后台展示用"
          />
        </a-form-item>
        <a-form-item
          name="vector_dimension_list"
          label="向量维度"
          v-if="formState.model_type == 'TEXT EMBEDDING'"
        >
          <a-input
            v-model:value="formState.vector_dimension_list"
            placeholder="请输入支持的向量维度，多个向量维度用英文逗号分割"
          />
          <div class="form-tip">向量维度必须为正整数，比如1024</div>
        </a-form-item>
        <a-form-item
          name="thinking_type"
          label="深度思考"
          required
          v-if="
            showThinkTypeList.includes(model_info.model_define) && formState.model_type == 'LLM'
          "
        >
          <a-radio-group class="thiing-radio-box" v-model:value="formState.thinking_type">
            <a-radio value="1">支持（类似deepseek r1这种只支持思考模式的模型）</a-radio>
            <a-radio value="2"
              >可选（比如doubao 1.6，可以通过接口控制是否启用深度思考模式）</a-radio
            >
            <a-radio value="0">不支持（比如deepseek v3，模型本身不支持深度思考模式）</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item
          name="function_call"
          label="工具调用"
          required
          v-if="formState.model_type == 'LLM'"
        >
          <a-radio-group v-model:value="formState.function_call">
            <a-radio value="1">支持</a-radio>
            <a-radio value="0">不支持</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item
          name="image_sizes"
          label="图片比列支持"
          required
          :rules="[{ required: true, message: '请选择图片比列支持' }]"
          v-if="formState.model_type == 'IMAGE'"
        >
          <a-checkbox-group
            class="image-size-options"
            v-model:value="formState.image_sizes"
            name="checkboxgroup"
            :options="sizeOptions"
          />
        </a-form-item>
        <a-form-item
          name="image_max"
          label="生成图片最大数量"
          :rules="[{ required: true, message: '请输入生成图片最大数量' }]"
          v-if="formState.model_type == 'IMAGE'"
        >
          <a-input-number
            style="width: 100%"
            v-model:value="formState.image_max"
            :max="10"
            :precision="0"
          />
        </a-form-item>
        <a-form-item
          name="image_watermark"
          label="支持图片水印"
          v-if="formState.model_type == 'IMAGE'"
        >
          <a-switch
            v-model:checked="formState.image_watermark"
            checked-children="开"
            un-checked-children="关"
          />
        </a-form-item>
        <a-form-item
          name="image_optimize_prompt"
          label="支持自动优化提示词"
          v-if="formState.model_type == 'IMAGE'"
        >
          <a-switch
            v-model:checked="formState.image_optimize_prompt"
            checked-children="开"
            un-checked-children="关"
          />
        </a-form-item>
        <a-form-item label="支持输入类型" required v-if="formState.model_type != 'RERANK'">
          <a-flex :gap="8" v-if="formState.model_type == 'LLM'">
            <a-checkbox v-model:checked="formState.input_text">文本</a-checkbox>
            <a-checkbox v-model:checked="formState.input_voice">语音</a-checkbox>
            <a-checkbox v-model:checked="formState.input_image">图片</a-checkbox>
            <a-checkbox v-model:checked="formState.input_video">视频</a-checkbox>
            <a-checkbox v-model:checked="formState.input_document">文档</a-checkbox>
          </a-flex>
          <a-flex :gap="8" v-else-if="formState.model_type == 'TTS'">
            <a-checkbox v-model:checked="formState.input_text">文本</a-checkbox>
          </a-flex>
          <a-flex :gap="8" align="center" v-else>
            <a-checkbox v-model:checked="formState.input_text">文本</a-checkbox>
            <a-checkbox v-model:checked="formState.input_image">图片</a-checkbox>
            <a-flex align="center" :gap="4" v-if="formState.model_type == 'IMAGE'">
              最多
              <a-input-number
                style="width: 120p"
                v-model:value="formState.image_inputs_image_max"
                :precision="0"
                :min="1"
                :max="4"
              />
              张图
            </a-flex>
          </a-flex>
        </a-form-item>
        <a-form-item label="支持输出类型" required v-if="formState.model_type == 'LLM'">
          <a-flex :gap="8">
            <a-checkbox v-model:checked="formState.output_text">文本</a-checkbox>
            <a-checkbox v-model:checked="formState.output_voice">语音</a-checkbox>
            <a-checkbox v-model:checked="formState.output_image">图片</a-checkbox>
            <a-checkbox v-model:checked="formState.output_video">视频</a-checkbox>
          </a-flex>
        </a-form-item>
        <a-form-item label="支持输出类型" required v-else-if="formState.model_type == 'TTS'">
          <a-flex :gap="8">
            <a-checkbox v-model:checked="formState.output_voice">语音</a-checkbox>
          </a-flex>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, h, reactive, computed } from 'vue'
import {} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { saveUseModelConfig } from '@/api/model/index'
import { getSizeOptions } from '@/views/workflow/components/util.js'
const open = ref(false)

const emit = defineEmits('ok')

const showThinkTypeList = ['siliconflow', 'doubao', 'tongyi']

const sizeOptions = getSizeOptions()

const formState = reactive({
  model_config_id: '',
  id: '',
  model_type: 'LLM',
  use_model_name: '',
  show_model_name: '',
  thinking_type: '0',
  function_call: '1',
  input_text: true,
  input_voice: false,
  input_image: false,
  input_video: false,
  input_document: false,
  output_text: true,
  output_voice: false,
  output_image: false,
  output_video: false,
  vector_dimension_list: '',
  image_max: 1,
  image_inputs_image_max: 1,
  image_sizes: sizeOptions.map((item) => item.value),
  image_optimize_prompt: true,
  image_watermark: true
})

const model_info = ref({})

const supported_type = ref([])

const model_name = computed(() => {
  if (model_info.value.model_define == 'doubao') {
    return '接入点id'
  }
  if (model_info.value.model_define == 'azure') {
    return '部署名称'
  }
  return 'model'
})

const formRef = ref(null)
const show = (data, record) => {
  model_info.value = data
  supported_type.value = data.supported_type || []
  resetData()
  formState.model_config_id = data.config_info.id
  if (record) {
    formState.id = record.id
    formState.model_type = record.model_type
    formState.use_model_name = record.use_model_name
    formState.show_model_name = record.show_model_name || record.use_model_name
    formState.thinking_type = record.thinking_type + ''
    formState.function_call = record.function_call + ''
    formState.input_text = record.input_text == 1
    formState.input_voice = record.input_voice == 1
    formState.input_image = record.input_image == 1
    formState.input_video = record.input_video == 1
    formState.input_document = record.input_document == 1
    formState.output_text = record.output_text == 1
    formState.output_voice = record.output_voice == 1
    formState.output_image = record.output_image == 1
    formState.output_video = record.output_video == 1
    formState.vector_dimension_list = record.vector_dimension_list
    if (record.image_generation) {
      let image_generation = JSON.parse(record.image_generation)
      formState.image_max = image_generation.image_max
      formState.image_inputs_image_max = image_generation.image_inputs_image_max
      formState.image_sizes = image_generation.image_sizes
        ? image_generation.image_sizes.split(',')
        : []
      formState.image_optimize_prompt = image_generation.image_optimize_prompt == 1
      formState.image_watermark = image_generation.image_watermark == 1
    }
  }
  open.value = true
}

const isLoading = ref(false)

const handleOk = () => {
  formRef.value.validate().then(() => {
    let image_generation = {
      image_sizes: formState.image_sizes.join(','),
      image_max: formState.image_max + '',
      image_watermark: formState.image_watermark ? '1' : '0',
      image_optimize_prompt: formState.image_optimize_prompt ? '1' : '0',
      image_inputs_image_max: formState.image_inputs_image_max + ''
    }
    let parmas = {
      ...formState,
      input_text: formState.input_text ? 1 : 0,
      input_voice: formState.input_voice ? 1 : 0,
      input_image: formState.input_image ? 1 : 0,
      input_video: formState.input_video ? 1 : 0,
      input_document: formState.input_document ? 1 : 0,
      output_text: formState.output_text ? 1 : 0,
      output_voice: formState.output_voice ? 1 : 0,
      output_image: formState.output_image ? 1 : 0,
      output_video: formState.output_video ? 1 : 0,
      image_generation: JSON.stringify(image_generation)
    }
    isLoading.value = true
    saveUseModelConfig({
      ...parmas
    })
      .then(() => {
        message.success(`保存成功`)
        emit('ok')
        open.value = false
      })
      .finally(() => {
        isLoading.value = false
      })
  })
}

const handleTypeChange = () => {
  formState.input_voice = false
  formState.input_image = false
  formState.input_video = false
  formState.input_document = false
  formState.output_text = false
  formState.output_voice = false
  formState.output_image = false
  formState.output_video = false
  if (formState.model_type == 'TTS') {
    formState.output_voice = true
  }
}

function resetData() {
  Object.assign(formState, {
    model_config_id: '',
    id: '',
    model_type: supported_type.value[0],
    use_model_name: '',
    show_model_name: '',
    thinking_type: '0',
    function_call: '1',
    input_text: true,
    input_voice: false,
    input_image: false,
    input_video: false,
    input_document: false,
    output_text: true,
    output_voice: false,
    output_image: false,
    output_video: false,
    vector_dimension_list: '',
    image_max: 1,
    image_inputs_image_max: 1,
    image_sizes: sizeOptions.map((item) => item.value),
    image_optimize_prompt: true,
    image_watermark: true
  })
}

const handleModelNameBlur = () => {
  if (!formState.show_model_name) {
    formState.show_model_name = formState.use_model_name
  }
}

defineExpose({
  show
})
</script>
<style lang="less" scoped>
.form-box {
  margin-top: 24px;
  &::v-deep(.ant-form-item) {
    margin-bottom: 16px;
    .ant-form-item-label {
      padding-bottom: 4px;
    }
  }
  .thiing-radio-box {
    .ant-radio-wrapper {
      margin-bottom: 4px;
    }
  }
  .form-tip {
    color: #8c8c8c;
    font-size: 14px;
    line-height: 22px;
    margin-top: 2px;
  }

  .image-size-options{
    &::v-deep(.ant-checkbox-wrapper) {
      width: 140px;
    }
  }

}
</style>
