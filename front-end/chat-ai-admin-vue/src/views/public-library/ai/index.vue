<template>
  <div class="ai-config-page">
    <ConfigPageMenu />
    <div class="page-container">
      <div class="ai-config-form">
        <a-form
          ref="formRef"
          :model="formState"
          :label-col="{ span: 2 }"
          :wrapper-col="{ span: 8 }"
        >
          <a-form-item>
            <template #label>
              嵌入&nbsp;
              <a-tooltip>
                <template #title>开启嵌入功能后，才能在机器人中关联本知识库。</template>
                <InfoCircleOutlined />
              </a-tooltip>
            </template>
            <a-switch
              v-model:checked="formState.use_model_switch"
              unCheckedValue="0"
              checkedValue="1"
            />
          </a-form-item>

          <a-form-item
            label="嵌入模型"
            name="use_model"
            :rules="[{ required: true, message: '请选择嵌入模型' }]"
            v-if="formState.use_model_switch == 1"
          >
            <!-- 自定义选择器 -->
            <CustomSelector
              v-model="formState.use_model"
              placeholder="请选择嵌入模型"
              label-key="use_model_name"
              value-key="value"
              :modelType="'TEXT EMBEDDING'"
              :model-config-id="formState.model_config_id"
              @change="handleModelChange"
            />
          </a-form-item>

          <a-form-item>
            <template #label>
              AI智能总结&nbsp;
              <a-tooltip>
                <template #title
                  >开启后，搜索知识库时，由大模型自动总结文档，并给出总结后的结果</template
                >
                <InfoCircleOutlined />
              </a-tooltip>
            </template>
            <a-switch v-model:checked="formState.ai_summary" unCheckedValue="0" checkedValue="1" />
          </a-form-item>

          <a-form-item
            label="总结模型"
            name="ai_summary_model"
            :rules="[{ required: true, message: '请选择总结模型' }]"
            v-if="formState.ai_summary == 1"
          >
            <!-- 自定义选择器 -->
            <CustomSelector
              v-model="formState.ai_summary_model"
              placeholder="请选择总结模型"
              label-key="summary_model_name"
              value-key="value"
              :modelType="'LLM'"
              :model-config-id="formState.summary_model_config_id"
              @change="handleModelChange2"
            />
          </a-form-item>

          <a-form-item :wrapper-col="{ offset: 2, span: 8 }">
            <a-button type="primary" :loading="loading" @click="handleSubmit">保存</a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getLibraryInfo, editLibrary } from '@/api/library/index'
import { useRoute } from 'vue-router'
import { ref, reactive, computed, onMounted, toRaw } from 'vue'
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import ConfigPageMenu from '../components/config-page-menu.vue'
import CustomSelector from '@/components/custom-selector/index.vue'

const route = useRoute()
const library_id = computed(() => route.query.library_id)
const formRef = ref()
const loading = ref(false)
let oldState = {}
const formState = reactive({
  library_name: '',
  library_intro: '',
  ai_summary: '',
  ai_summary_model: '',
  access_rights: '',
  model_config_id: '',
  use_model: '',
  use_model_icon: '', // 新增图标字段
  use_model_name: '', // 新增系统名称
  use_model_switch: '',
  statistics_set: '',
  summary_model_icon: '', // 新增图标字段
  summary_model_name: '', // 新增系统名称
  summary_model_config_id: ''
})

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']
const currentModelDefine = ref('')
const currentModelDefine2 = ref('')

// 处理选择事件
const handleModelChange = (item) => {
  formState.use_model = modelDefine.includes(item.rawData.model_define) && item.rawData.deployment_name 
    ? item.rawData.deployment_name 
    : item.rawData.name
  formState.use_model_icon = item.icon
  formState.use_model_name = item.use_model_name
  formState.model_config_id = item.rawData.id
  currentModelDefine.value = item.rawData.model_define
}

// 处理选择事件
const handleModelChange2 = (item) => {
  formState.ai_summary_model = modelDefine.includes(item.rawData.model_define) && item.rawData.deployment_name 
    ? item.rawData.deployment_name 
    : item.rawData.name
  formState.summary_model_icon = item.icon
  formState.summary_model_name = item.summary_model_name
  formState.summary_model_config_id = item.rawData.id
  currentModelDefine2.value = item.rawData.model_define
}

const handleChangeModel = (val, option) => {
  const self = option.current_obj
  formState.use_model =
    modelDefine.indexOf(self.model_define) > -1 && self.deployment_name
      ? self.deployment_name
      : self.name
  currentModelDefine.value = self.model_define
  formState.model_config_id = self.id || option.model_config_id
}

const handleChangeModel2 = (val, option) => {
  const self = option.current_obj
  formState.ai_summary_model =
    modelDefine.indexOf(self.model_define) > -1 && self.deployment_name
      ? self.deployment_name
      : self.name
  currentModelDefine2.value = self.model_define
  formState.summary_model_config_id = self.id || option.model_config_id
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    // TODO: 调用API保存数据
    let isChangeModel = oldState.use_model != formState.use_model
    if (isChangeModel) {
      Modal.confirm({
        title: '确定切换模型为' + formState.use_model + '吗？',
        content: '切换模型后，文档会重新向量化，会消耗token',
        okText: '确定',
        cancelText: '取消',
        onOk() {
          saveData()
        }
      })
    } else {
      saveData()
    }
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const saveData = () => {
  let data = {
    ...toRaw(formState)
  }

  if (oldModelDefineList.indexOf(currentModelDefine.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    data.use_model = '默认'
  }

  if (oldModelDefineList.indexOf(currentModelDefine2.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    data.ai_summary_model = '默认'
  }

  loading.value = true
  editLibrary(data)
    .then(() => {
      message.success('保存成功')
      loading.value = false
      oldState = JSON.parse(JSON.stringify(data))
    })
    .catch(() => {
      loading.value = false
    })
}

const getData = () => {
  getLibraryInfo({ id: library_id.value }).then((res) => {
    let data = res.data || {}
    data.use_model = data.use_model || void 0
    data.ai_summary_model = data.ai_summary_model || void 0
    oldState = JSON.parse(JSON.stringify(data))
    Object.assign(formState, data)
  })
}

onMounted(() => {
  getData()
})
</script>

<style lang="less" scoped>
.ai-config-form {
  padding: 8px 24px 24px;
}
.model-icon {
  height: 18px;
}
</style>
