<template>
  <edit-box class="setting-box" title="模型设置" icon-name="moxingshezhi" v-model:isEdit="isEdit">
    <template #extra>
      <div class="actions-box">
        <a-button @click="handleEdit(true)" size="small">修改</a-button>
      </div>
    </template>
    <div class="setting-info-block">
      <div class="set-item w-100">
        LLM模型：
        <span class="set-item-content" v-if="formState.use_model_icon && formState.use_model_name">
          <img :src="formState.use_model_icon" class="model-icon" />{{ formState.use_model_name }} {{ formState.use_model }}
        </span>
        <span v-else>{{ formState.use_model }}</span>
      </div>
      <div class="set-item">
        温度：
        <span>{{ formState.temperature }}</span>
      </div>
      <div class="set-item">
        最大token：
        <span>{{ formState.max_token }}</span>
      </div>
      <div class="set-item">
        上下文数量：
        <span>{{ formState.context_pair }}</span>
      </div>
      <div class="set-item">
        显示推理过程：
        <span>{{ formState.think_switch == 1 ? '开' : '关' }}</span>
      </div>
    </div>
    <a-modal v-model:open="isEdit" :width="672" title="模型设置" @ok="handleSave">
      <div class="form-box">
        <a-row :gutter="[32, 24]">
          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label mb12">
                <span>LLM模型&nbsp;</span>
              </div>
              <div class="form-item-body">
                <!-- 自定义选择器 -->
                <CustomSelector
                  v-model="formState.use_model"
                  :options="processedModelList"
                  placeholder="请选择LLM模型"
                  label-key="use_model_name"
                  value-key="value"
                  :model-define="modelDefine"
                  :model-config-id="formState.model_config_id"
                  @change="handleModelChange"
                />
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label">
                <span>温度&nbsp;</span>
                <a-tooltip>
                  <template #title>温度越低，回答越严谨。温度越高，回答越发散。</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <div class="form-item-body">
                <div class="number-box">
                  <div class="number-slider-box">
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.temperature"
                      :min="0"
                      :max="2"
                      :step="0.1"
                    />
                  </div>
                  <div class="number-input-box">
                    <a-input-number
                      v-model:value="formState.temperature"
                      :min="0"
                      :max="2"
                      :step="0.1"
                    />
                  </div>
                </div>
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label">
                <span>最大token&nbsp;</span>
                <a-tooltip>
                  <template #title>问题+答案的最大token数，如果出现回答被截断，可调高此值</template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <div class="form-item-body">
                <div class="number-box">
                  <div class="number-slider-box">
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.max_token"
                      :min="0"
                      :max="100 * 1024"
                    />
                  </div>
                  <div class="number-input-box">
                    <a-input-number
                      v-model:value="formState.max_token"
                      :min="0"
                      :max="100 * 1024"
                    />
                  </div>
                </div>
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item">
              <div class="form-item-label">
                <span>上下文数量&nbsp;</span>
                <a-tooltip>
                  <template #title
                    >提示词中携带的历史聊天记录轮次。设置为0则不携带聊天记录。最多设置10轮。注意，携带的历史聊天记录越多，消耗的token相应也就越多。</template
                  >
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <div class="form-item-body">
                <div class="number-box">
                  <div class="number-slider-box">
                    <a-slider
                      class="custom-slider"
                      v-model:value="formState.context_pair"
                      :min="0"
                      :max="10"
                    />
                  </div>
                  <div class="number-input-box">
                    <a-input-number v-model:value="formState.context_pair" :min="0" :max="10" />
                  </div>
                </div>
              </div>
            </div>
          </a-col>

          <a-col v-bind="grid">
            <div class="form-item justify-between">
              <div class="form-item-label">
                <span>显示推理过程&nbsp;</span>
                <a-tooltip>
                  <template #title>
                    <span
                      >开启后，API接口、聊天测试页、 webAPP页会显示或返回推理模型的推理过程</span
                    >
                  </template>
                  <QuestionCircleOutlined class="question-icon" />
                </a-tooltip>
              </div>
              <a-switch
                v-model:checked="formState.think_switch"
                :checkedValue="1"
                :unCheckedValue="0"
              />
            </div>
          </a-col>
        </a-row>
      </div>
    </a-modal>
  </edit-box>
</template>

<script setup>
import { getModelConfigOption } from '@/api/model/index'
import { ref, reactive, inject, toRaw, watchEffect, computed, onMounted, onBeforeUnmount } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import EditBox from './edit-box.vue'
import { duplicateRemoval, removeRepeat } from '@/utils/index'
import CustomSelector from '@/components/custom-selector/index.vue'

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']

const grid = reactive({ sm: 24, md: 24, lg: 12, xl: 24, xxl: 24 })
// 获取LLM模型
const modelList = ref([])

const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')
const formState = reactive({
  use_model: undefined,
  use_model_icon: '', // 新增图标字段
  use_model_name: '', // 新增系统名称
  model_config_id: '',
  temperature: 0,
  max_token: 0,
  context_pair: 0,
  think_switch: 1
})
const currentModelDefine = ref('')
const oldModelDefine = ref('')
const currentUseModel = ref('')

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

// 新增自定义选择器逻辑
const showOptions = ref(false)
const clickOutsideHandler = (e) => {
  if (!e.target.closest('.custom-select')) {
    showOptions.value = false
  }
}

// 事件监听
onMounted(() => {
  // 首次进入初始化回显数据
  if (modelDefine.indexOf(currentModelDefine.value) > -1) {
    formState.use_model = currentUseModel.value
  } else {
    formState.use_model = robotInfo.use_model
  }

  document.addEventListener('click', clickOutsideHandler)
})
onBeforeUnmount(() => {
  document.removeEventListener('click', clickOutsideHandler)
})

const handleSave = () => {
  // 初始化条件
  currentUseModel.value = formState.use_model

  let params = { ...toRaw(formState) }
  if (oldModelDefineList.indexOf(currentModelDefine.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    params.use_model = '默认'
  }

  updateRobotInfo(params)
  isEdit.value = false
}

const handleEdit = (val) => {
  if (!val) {
    // 修改选择的是取消按钮，初始化条件数据
    currentModelDefine.value = oldModelDefine.value
  }
  if (modelDefine.indexOf(currentModelDefine.value) > -1) {
    formState.use_model = currentUseModel.value
  } else {
    formState.use_model = robotInfo.use_model
  }
  formState.model_config_id = robotInfo.model_config_id
  formState.temperature = robotInfo.temperature
  formState.max_token = robotInfo.max_token
  formState.context_pair = robotInfo.context_pair
  formState.think_switch = Number(robotInfo.think_switch)
  isEdit.value = val
}

watchEffect(() => {
  formState.model_config_id = robotInfo.model_config_id
  formState.temperature = robotInfo.temperature
  formState.max_token = robotInfo.max_token
  formState.context_pair = robotInfo.context_pair
  formState.think_switch = Number(robotInfo.think_switch)
})

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

const getModelList = () => {
  getModelConfigOption({
    model_type: 'LLM'
  }).then((res) => {
    currentUseModel.value = robotInfo.use_model
    let list = res.data || []
    let children = []

    modelList.value = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.llm_model_list.length; i++) {
        const ele = item.model_info.llm_model_list[i]
        if (
          modelDefine.indexOf(item.model_info.model_define) > -1 &&
          robotInfo.model_config_id == item.model_config.id
        ) {
          currentUseModel.value = item.model_config.deployment_name
          currentModelDefine.value = item.model_info.model_define
          oldModelDefine.value = item.model_info.model_define
        }
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

    // 初始化当前图标
    if (robotInfo.use_model && robotInfo.model_config_id) {
      modelList.value.some(group => 
        group.children.some(child => {
          const childName = modelDefine.includes(child.model_define) && child.deployment_name 
            ? child.deployment_name 
            : child.name;
          if (modelDefine.includes(child.model_define) && child.deployment_name && child.id === robotInfo.model_config_id) {
            formState.use_model_icon = child.icon;
            formState.use_model_name = child.use_model_name;
            formState.use_model = child.deployment_name
            return true;
          }
          if (childName === robotInfo.use_model && child.id === robotInfo.model_config_id) {
            formState.use_model_icon = child.icon;
            formState.use_model_name = child.use_model_name;
            return true;
          }
        })
      );
    }
  })
}

getModelList()
</script>

<style lang="less" scoped>
.setting-box {
  ::v-deep(.edit-box-body) {
    padding: 0;
  }

  .actions-box {
    display: flex;
    align-items: center;
    line-height: 22px;
    font-size: 14px;
    color: #595959;

    .action-btn {
      cursor: pointer;
    }

    .save-btn {
      color: #2475fc;
    }

    .model-name {
      font-size: 14px;
      line-height: 22px;
      color: #8c8c8c;
    }
  }
}

.set-item-content {
  display: flex;
  align-items: center;
  gap: 6px;
}

.model-icon {
  height: 18px;
  width: 18px;
  object-fit: contain;
  vertical-align: middle;
}

/* 下拉选项对齐优化 */
.ant-select-item-option-content {
  display: flex;
  align-items: center;
  gap: 6px;
}

.setting-info-block {
  padding: 16px;
  padding-top: 0;
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
  color: #595959;
  line-height: 22px;
  .set-item {
    display: flex;
    align-items: center;
  }
  .w-100 {
    width: 100%;
  }
}

.form-box {
  padding: 16px 0 16px 0;

  .justify-between {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .form-item-label {
    line-height: 22px;
    margin-bottom: 4px;
    font-size: 14px;
    color: #262626;

    .question-icon {
      color: #8c8c8c;
    }
  }

  .number-box {
    display: flex;
    align-items: center;

    .number-slider-box {
      flex: 1;
    }

    .number-input-box {
      margin-left: 20px;
    }
  }
}
</style>
