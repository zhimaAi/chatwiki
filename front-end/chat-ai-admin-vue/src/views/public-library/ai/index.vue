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
              :options="processedModelList"
              placeholder="请选择嵌入模型"
              label-key="use_model_name"
              value-key="value"
              :model-define="modelDefine"
              :model-config-id="formState.model_config_id"
              @change="handleModelChange"
            />
            <!-- <a-select
              @change="handleChangeModel"
              v-model:value="formState.use_model"
              placeholder="请选择嵌入模型"
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
            </a-select> -->
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
              :options="processedModelList2"
              placeholder="请选择总结模型"
              label-key="summary_model_name"
              value-key="value"
              :model-define="modelDefine"
              :model-config-id="formState.summary_model_config_id"
              @change="handleModelChange2"
            />
            <!-- <a-select
              v-model:value="formState.ai_summary_model"
              placeholder="请选择总结模型"
              @change="handleChangeModel2"
              style="width: 100%"
            >
              <a-select-opt-group v-for="item in modelList2" :key="item.id">
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
            </a-select> -->
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
import { duplicateRemoval, removeRepeat } from '@/utils/index'
import { getModelConfigOption } from '@/api/model/index'
import { getLibraryInfo, editLibrary } from '@/api/library/index'
import { useRoute } from 'vue-router'
import { ref, reactive, computed, onMounted, toRaw } from 'vue'
import { InfoCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import ConfigPageMenu from '../components/config-page-menu.vue'
import CustomSelector from '@/components/custom-selector/index.vue'

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

// 获取嵌入模型列表
const modelList = ref([])

// 获取LLM模型
const modelList2 = ref([])

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']
const currentModelDefine = ref('')
const currentModelDefine2 = ref('')

// 处理原始数据格式
const processedModelList = computed(() => {
  return modelList.value.map(group => ({
    groupLabel: group.name,
    groupIcon: group.icon,
    children: group.children.map(child => ({
      icon: child.icon,
      use_model_name: child.use_model_name,
      value: modelDefine.includes(child.model_define) && child.deployment_name ? child.deployment_name : child.name,
      rawData: child // 保留原始数据
    }))
  }))
})

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

// 处理原始数据格式
const processedModelList2 = computed(() => {
  return modelList2.value.map(group => ({
    groupLabel: group.name,
    groupIcon: group.icon,
    children: group.children.map(child => ({
      icon: child.icon,
      summary_model_name: child.summary_model_name,
      value: modelDefine.includes(child.model_define) && child.deployment_name ? child.deployment_name : child.name,
      rawData: child // 保留原始数据
    }))
  }))
})

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

const getModelList = (is_offline) => {
  getModelConfigOption({
    model_type: 'TEXT EMBEDDING',
    is_offline
  }).then((res) => {
    let list = res.data || []
    let children = []

    modelList.value = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.vector_model_list.length; i++) {
        const ele = item.model_info.vector_model_list[i]
        children.push({
          name: ele,
          deployment_name: item.model_config.deployment_name,
          id: item.model_config.id,
          model_define: item.model_info.model_define,
          icon: item.model_info.model_icon_url, // 添加图标字段
          use_model_name: item.model_info.model_name // 添加系统名称字段
        })
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

    // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
    modelList.value = uniqueArr(
      duplicateRemoval(modelList.value, 'model_define'),
      modelList.value,
      'model_define'
    )
  })
}

const getModelList2 = () => {
  getModelConfigOption({
    model_type: 'LLM'
  }).then((res) => {
    // currentUseModel.value = robotInfo.use_model
    let list = res.data || []
    let children = []

    modelList2.value = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.llm_model_list.length; i++) {
        const ele = item.model_info.llm_model_list[i]
        // if (modelDefine.indexOf(item.model_info.model_define) > -1 && robotInfo.model_config_id == item.model_config.id) {
        //   currentUseModel.value = item.model_config.deployment_name
        //   currentModelDefine.value = item.model_info.model_define
        //   oldModelDefine.value = item.model_info.model_define
        // }
        children.push({
          name: ele,
          deployment_name: item.model_config.deployment_name,
          id: item.model_config.id,
          model_define: item.model_info.model_define,
          icon: item.model_info.model_icon_url, // 添加图标字段
          summary_model_name: item.model_info.model_name // 添加系统名称字段
        })
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

    // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
    modelList2.value = uniqueArr(
      duplicateRemoval(modelList2.value, 'model_define'),
      modelList2.value,
      'model_define'
    )
  })
}

onMounted(() => {
  getData()
  getModelList()
  getModelList2()
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
