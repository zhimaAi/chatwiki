<style lang="less" scoped>
.prompt-form {
  margin-top: 24px;

  .prompt-form-item {
    margin: 24px 0;
  }

  .prompt-form-label {
    line-height: 22px;
    margin-bottom: 8px;
    font-size: 14px;
    color: #262626;
  }

  .question-icon {
    color: #8c8c8c;
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
.model-icon {
  height: 18px;
}
</style>

<template>
  <div class="page-right">
    <div class="page-right-body">
      <div class="prompt-tips">
        通过修改提示词可知识库可以优化机器人回答的效果。<span style="color: #ed744a"
          >测试时可以直接修改提示词，不用保存即可测试优化效果。提示词优化好之后再点击保存，才会对外生效</span
        >
      </div>

      <div class="prompt-form">
        <div class="prompt-form-label">提示词</div>
        <div clas="prompt-form-item">
          <a-textarea
            :value="formState.prompt"
            :auto-size="{ minRows: 10, maxRows: 10 }"
            placeholder="请输入内容"
            @input="handlePromptChange"
          ></a-textarea>
        </div>
      </div>

      <div class="prompt-form">
        <div class="prompt-form-label">模型设置</div>

        <div clas="prompt-form-item">
          <div class="prompt-form-item-label">
            <span>LLM模型</span>
          </div>
          <div class="prompt-form-item-body">
            <a-select
              v-model:value="formState.use_model"
              @change="handleChangeModel"
              style="width: 100%"
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
            </a-select>
          </div>
        </div>

        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span>最大token</span>
            <a-tooltip>
              <template #title>
                <div>问题+答案的最大token数，如果出现回答被截断，可调高此值</div>
              </template>
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
                <a-input-number v-model:value="formState.max_token" :min="0" :max="100 * 1024" />
              </div>
            </div>
          </div>
        </div>

        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span>温度</span>
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

        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span>上下文数量</span>
            <a-tooltip>
              <template #title
                >提示词中携带的历史聊天记录轮次。设置为0则不携带聊天记录。最多设置50轮。注意，携带的历史聊天记录越多，消耗的token相应也就越多。</template
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
                  :max="50"
                />
              </div>
              <div class="number-input-box">
                <a-input-number v-model:value="formState.context_pair" :min="0" :max="50" />
              </div>
            </div>
          </div>
        </div>

        <div class="prompt-form-item">
          <div class="prompt-form-item-label">
            <span>显示推理过程</span>
            <a-tooltip>
              <template #title>
                <span>开启后，API接口、聊天测试页、 webAPP页会显示或返回推理模型的推理过程</span>
              </template>
              <QuestionCircleOutlined class="question-icon" />
            </a-tooltip>
          </div>
          <div class="form-item-body" style="padding-top: 5px">
            <a-switch
              v-model:checked="formState.think_switch"
              :checkedValue="1"
              :unCheckedValue="0"
            />
          </div>
        </div>
      </div>
    </div>
    <div class="page-right-footer">
      <a-button type="primary" block @click="handleSaveRobotPrompt">保存并应用</a-button>
    </div>
  </div>
</template>

<script setup>
import { getModelConfigOption } from '@/api/model/index'
import { ref, toRaw, onMounted, reactive } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { duplicateRemoval, removeRepeat } from '@/utils/index'

const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const oldModelDefineList = ['azure']

const props = defineProps({
  robotInfo: {
    type: Object,
    default: () => {}
  }
})

const emit = defineEmits(['promptChange', 'saveRobotPrompt'])

const formState = reactive({
  use_model: undefined,
  model_config_id: '',
  temperature: 0,
  max_token: 0,
  prompt: props.robotInfo.prompt,
  context_pair: 0,
  think_switch: 1
})
const currentModelDefine = ref('')

// 获取LLM模型
const modelList = ref([])

const updateRobotInfo = (val) => {
  let newState = JSON.parse(JSON.stringify(val))
  newState.max_token = parseFloat(newState.max_token)
  newState.temperature = parseFloat(newState.temperature)
  newState.context_pair = parseFloat(newState.context_pair)

  Object.keys(newState).forEach((key) => {
    formState[key] = newState[key]
  })
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
    let list = res.data || []
    let children = []

    modelList.value = list.map((item) => {
      children = []
      for (let i = 0; i < item.model_info.llm_model_list.length; i++) {
        const ele = item.model_info.llm_model_list[i]
        if (
          modelDefine.indexOf(item.model_info.model_define) > -1 &&
          props.robotInfo.model_config_id == item.model_config.id
        ) {
          formState.use_model = item.model_config.deployment_name
          currentModelDefine.value = item.model_info.model_define
        }
        children.push({
          name: ele,
          deployment_name: item.model_config.deployment_name,
          id: item.model_config.id,
          model_define: item.model_info.model_define
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

const handlePromptChange = (e) => {
  formState.prompt = e.currentTarget.value
  emit('promptChange', e)
}

const handleSaveRobotPrompt = () => {
  let isDefault = false
  if (oldModelDefineList.indexOf(currentModelDefine.value) > -1) {
    // 传给后端的是默认，渲染的是真实名称
    isDefault = true
  }
  emit('saveRobotPrompt', formState, isDefault)
}
// const emit = defineEmits(['openChat', 'onScrollEnd'])

onMounted(() => {
  // 获取llm
  getModelList()
  updateRobotInfo({ ...toRaw(props.robotInfo) })
})
</script>
