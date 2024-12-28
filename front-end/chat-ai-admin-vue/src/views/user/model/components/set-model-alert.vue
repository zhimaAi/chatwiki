<style lang="less" scoped>
.add-model-alert {
  .form-wrapper {
    margin-top: 24px;
  }

  .model-logo {
    height: 26px;
    margin-bottom: 16px;
  }
  .model-logo-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
</style>

<template>
  <a-modal
    class="add-model-alert"
    width="650px"
    v-model:open="show"
    :title="title"
    :confirmLoading="confirmLoading"
    @ok="handleOk"
    @cancel="handleCancel"
    :maskClosable="false"
  >
    <div class="form-wrapper">
      <div class="model-logo-box">
        <img class="model-logo" :src="modelConfig.model_icon_url" alt="" />
        <a
          v-if="modelConfig.model_define == 'doubao'"
          href="https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/hmwl4gveq4xwa4uw?singleDoc# "
          target="_blank"
          >查看使用说明></a
        >
      </div>
      <div class="form-box">
        <a-form layout="horizontal" :label-col="{ span: 6 }" :wrapper-col="{ span: 18 }">
          <template v-for="item in formItems" :key="item.key">
            <a-form-item v-bind="validateInfos[item.key]" v-if="!item.hidden">
              <template #label>
                <div class="form-item-label">
                  <span v-if="modelConfig.model_define == 'doubao' && item.key == 'deployment_name'"
                    >接入点ID</span
                  >
                  <span v-else>{{ t('views.user.model.' + item.key) }}</span>
                </div>
              </template>
              <template v-if="item.componentType == 'radio'">
                <a-radio-group v-model:value="formState[item.key]" :disabled="item.disabled">
                  <a-radio :value="val" v-for="(val, index) in item.options" :key="index"
                    ><span>{{ val }}</span></a-radio
                  >
                </a-radio-group>
              </template>

              <template v-if="item.componentType == 'select'">
                <a-select v-model:value="formState[item.key]" :placeholder="item.placeholder">
                  <a-select-option :value="val" v-for="(val, index) in item.options" :key="index">
                    <span>{{ val }}</span>
                  </a-select-option>
                </a-select>
              </template>

              <template v-if="item.componentType == 'input'">
                <a-input v-model:value="formState[item.key]" :placeholder="item.placeholder" />
                <a :href="item.help_links" target="_blank" v-if="item.key == 'api_key'">{{
                  t('views.user.model.getApikeyText')
                }}</a>
              </template>
            </a-form-item>
          </template>
        </a-form>
      </div>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { getModelFieldConfig } from '../model-config'
import { Form, message } from 'ant-design-vue'
import { addModelConfig, editModelConfig } from '@/api/model/index'

const { t } = useI18n()
const emit = defineEmits(['ok'])
const useForm = Form.useForm

const title = ref('')
const modelConfig = reactive({
  model_icon_url: '',
  model_define: '',
  model_name: ''
})

const show = ref(false)
const formRules = ref({})
const formState = ref({})
const formItems = ref([])
const modelNames = ['ollama', 'xinference', 'openaiAgent']

const open = (data, record) => {
  formRules.value = {}
  formState.value = {}
  formItems.value = []

  const { model_icon_url, model_define, config_params, model_name } = data

  title.value = record ? t('common.edit') + ' ' + model_name : t('common.add') + ' ' + model_name

  modelConfig.model_icon_url = model_icon_url
  modelConfig.model_define = model_define
  modelConfig.model_name = model_name

  let formRulesTemp = {}
  let formStateTemp = {}
  let formItemsTemp = []

  config_params.unshift('id')

  config_params.forEach((key) => {
    if (key == 'model_type') {
      key = 'model_types'
    }

    if (modelNames.indexOf(data.model_define) > -1 && key == 'deployment_name') {
      key = 'model_name'
    }

    let field = getModelFieldConfig(key)

    field.key = key
    field.help_links = data.help_links

    field.disabled = false
    // 编辑时禁止修改模型类型
    if (field.key == 'model_types' && record) {
      field.disabled = true
    }

    let options = field.options

    if (record) {
      // 区分modelNames
      if (modelNames.indexOf(record.model_define) > -1) {
        if (record.deployment_name) {
          record.model_name = record.deployment_name
        }
      }
      formStateTemp[key] = record[key] || undefined
    } else {
      formStateTemp[key] = field.defaultValue
    }

    if (field.optionsKey && data[field.optionsKey]) {
      options = data[field.optionsKey]
    }

    if (!field.hidden) {
      formItemsTemp.push({ ...field, options, key })
    }

    if (field.rules) {
      formRulesTemp[key] = field.rules
    }
  })
  if (data.model_define == 'doubao') {
    // 火山模型特殊处理
    if (!record) {
      formStateTemp.region = 'cn-beijing'
    }
    formItemsTemp.forEach((item) => {
      if (item.key == 'deployment_name') {
        item.placeholder = '请输入接入点ID'
      }
    })
    if (formRulesTemp.deployment_name && formRulesTemp.deployment_name.length) {
      formRulesTemp.deployment_name[0].message = '请输入接入点ID'
    }
  } else {
    if (formRulesTemp.deployment_name && formRulesTemp.deployment_name.length) {
      formRulesTemp.deployment_name[0].message = '请输入部署名称'
    }
  }

  formRules.value = formRulesTemp
  formState.value = formStateTemp
  formItems.value = formItemsTemp
  show.value = true

  setTimeout(() => {
    clearValidate()
  }, 100)
}

const { resetFields, clearValidate, validate, validateInfos } = useForm(formState, formRules)

const handleOk = () => {
  validate()
    .then(() => {
      if (formState.value.id) {
        saveEditModel()
      } else {
        saveAddModel()
      }
    })
    .catch((err) => {
      console.log(err)
    })
}

const handleCancel = () => {
  resetFields()
  clearValidate()
}

const confirmLoading = ref(false)

const saveAddModel = () => {
  let newData = toRaw(formState.value)
  let data = JSON.parse(JSON.stringify(newData)) // 深拷贝，不能改变原对象

  data.model_type = data.model_types
  data.model_define = modelConfig.model_define

  // 区分modelNames
  if (modelNames.indexOf(data.model_define) > -1) {
    if (data.model_name) {
      data.deployment_name = data.model_name
      delete data.model_name
    }
  }

  confirmLoading.value = true

  addModelConfig(data)
    .then(() => {
      confirmLoading.value = false
      show.value = false
      message.success(t('common.saveSuccess'))

      emit('ok')
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

const saveEditModel = () => {
  let newData = toRaw(formState.value)
  let data = JSON.parse(JSON.stringify(newData)) // 深拷贝，不能改变原对象

  data.model_type = data.model_types
  data.model_define = modelConfig.model_define

  // 区分modelNames
  if (modelNames.indexOf(data.model_define) > -1) {
    if (data.model_name) {
      data.deployment_name = data.model_name
      delete data.model_name
    }
  }

  confirmLoading.value = true

  editModelConfig(data)
    .then(() => {
      confirmLoading.value = false
      show.value = false
      message.success(t('common.saveSuccess'))

      emit('ok')
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

defineExpose({
  open
})
</script>
