<template>
  <div>
    <a-modal v-model:open="open" title="AI 生成提示词" :width="746">
      <template #footer>
        <a-button type="primary" :loading="confirmLoading" @click="handleOk">使用该提示</a-button>
      </template>

      <div class="ai-create-box">
        <div class="input-box">
          <ModelSelect
            modelType="LLM"
            v-model:modeName="formState.use_model"
            v-model:modeId="formState.model_config_id"
            style="width: 300px"
            @loaded="onVectorModelLoaded"
          />
        </div>
        <div class="input-box">
          <a-input
            style="width: 100%"
            v-model:value="demand"
            size="large"
            @pressEnter="handleCreatePrompt"
            placeholder="请输入AI 生成提示词例如：创建一个电商行业售后客服"
          ></a-input>
          <div class="btn-box" @click="handleCreatePrompt">
            <svg-icon name="ai-mark" />
            <span>生成</span>
          </div>
        </div>
        <div class="quick-tags-box">
          <div class="quick-label">快捷生成：</div>
          <a-tag
            v-for="item in quickData"
            :key="item.title"
            @click="handleQuickMark(item.title)"
            color="blue"
            >{{ item.title }}</a-tag
          >
        </div>
        <div class="ai-list-box">
          <cu-scroll style="height: 100%">
            <template v-if="isLoading">
              <div class="loading-box">
                <a-spin tip="正在生成中..."></a-spin>
              </div>
            </template>
            <template v-else>
              <div class="prompt-list-box" v-if="hasData">
                <div class="prompt-list">
                  <div class="prompt-header">
                    <div class="prompt-title">{{ formState.promptStruct.role.subject }}</div>
                  </div>
                  <div class="prompt-content">
                    {{ formState.promptStruct.role.describe }}
                  </div>
                </div>
                <!-- 任务 -->
                <div class="prompt-list">
                  <div class="prompt-header">
                    <div class="prompt-title">{{ formState.promptStruct.task.subject }}</div>
                  </div>
                  <div class="prompt-content">
                    {{ formState.promptStruct.task.describe }}
                  </div>
                </div>
                <!-- 要求 -->
                <div class="prompt-list">
                  <div class="prompt-header">
                    <div class="prompt-title">
                      {{ formState.promptStruct.constraints.subject }}
                    </div>
                  </div>
                  <div class="prompt-content">
                    {{ formState.promptStruct.constraints.describe }}
                  </div>
                </div>
                <!-- 技能 -->
                <div class="prompt-list">
                  <div class="prompt-header">
                    <div class="prompt-title">
                      {{ formState.promptStruct.skill.subject }}
                    </div>
                  </div>
                  <div class="prompt-content">
                    {{ formState.promptStruct.skill.describe }}
                  </div>
                </div>
                <!-- 输出格式 -->
                <div class="prompt-list">
                  <div class="prompt-header">
                    <div class="prompt-title">{{ formState.promptStruct.output.subject }}</div>
                  </div>
                  <div class="prompt-content">
                    {{ formState.promptStruct.output.describe }}
                  </div>
                </div>
                <!-- 风格 -->
                <div class="prompt-list">
                  <div class="prompt-header">
                    <div class="prompt-title">{{ formState.promptStruct.tone.subject }}</div>
                  </div>
                  <div class="prompt-content">
                    {{ formState.promptStruct.tone.describe }}
                  </div>
                </div>

                <!-- 自定义 -->
                <div
                  class="prompt-list"
                  v-for="(item, index) in formState.promptStruct.custom"
                  :key="index + item.key ? item.key : ''"
                >
                  <div class="prompt-header">
                    <div class="prompt-title">{{ item.subject }}</div>
                  </div>
                  <div class="prompt-content">{{ item.describe }}</div>
                </div>
              </div>
              <div class="empty-box" v-else>
                <a-empty :image="simpleImage" :description="description" />
              </div>
            </template>
          </cu-scroll>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, nextTick } from 'vue'
import { createPromptByLlm } from '@/api/user/index.js'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { Empty } from 'ant-design-vue'
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const open = ref(false)
const confirmLoading = ref(false)
const query = useRoute().query
const emit = defineEmits(['handleAiSave'])

const demand = ref('')

const isLoading = ref(false)

const quickData = [
  {
    title: '创建一个电商行业售后客服'
  },
  {
    title: '创建一个英语教育机构官网售前机器人'
  }
]

const description = ref('在上方输入框输入要求，点击生成提示词')
const show = () => {
  hasData.value = false
  formState.promptStruct = defaultData
  demand.value = ''
  description.value = '在上方输入框输入要求，点击生成提示词'
  open.value = true
}
const hasData = ref(false)
const defaultData = {
  role: {
    subject: '',
    describe: ''
  },
  task: {
    subject: '',
    describe: ''
  },
  constraints: {
    subject: '',
    describe: ''
  },
  skill: {
    subject: '',
    describe: ''
  },
  output: {
    subject: '',
    describe: ''
  },
  tone: {
    subject: '',
    describe: ''
  },
  custom: []
}
const formState = reactive({
  promptStruct: defaultData,
  use_model: '',
  model_config_id: ''
})
const handleQuickMark = (title) => {
  demand.value = title
}
const handleCreatePrompt = () => {
  if (!demand.value) {
    return message.error('请输入AI 生成提示词')
  }
  if (isLoading.value) {
    return message.error('正在生成中...')
  }
  isLoading.value = true
  createPromptByLlm({
    demand: demand.value,
    use_model: formState.use_model,
    model_config_id: formState.model_config_id
  })
    .then((res) => {
      formState.promptStruct = JSON.parse(res.data.promptStruct)
      hasData.value = true
    })
    .finally(() => {
      isLoading.value = false
    })
    .catch((res) => {
      hasData.value = false
      description.value = res.msg
    })
}

const onVectorModelLoaded = (list) => {
  nextTick(() => {
    formState.model_config_id = list[0]?.model_config_id
    formState.use_model = list[0]?.children[0]?.name
  })
}
const handleOk = () => {
  if (!hasData.value) {
    return message.error('请先生成提示词')
  }
  emit('handleAiSave', formState.promptStruct)
  open.value = false
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.loading-box {
  display: flex;
  align-items: center;
  justify-content: center;
  padding-top: 120px;
}
.empty-box {
  padding-top: 80px;
  ::v-deep(.ant-empty-description) {
    color: #262626;
  }
}
.ai-create-box {
  .input-box {
    position: relative;
    margin-top: 20px;
    .btn-box {
      cursor: pointer;
      position: absolute;
      right: 16px;
      top: 8px;
      display: flex;
      height: 24px;
      align-items: center;
      gap: 4px;
      font-size: 16px;
      line-height: 24px;
      color: #6524fc;
    }
  }
  .ai-list-box {
    margin-top: 8px;
    border-radius: 12px;
    height: 382px;
    padding: 16px 3px 16px 16px;
    border: 1px solid #2475fc;
    background: #f0f5ff;
  }
}

.prompt-list-box {
  .prompt-list {
    border: 1px solid var(--06, #d9d9d9);
    background: #fff;
    border-radius: 6px;
    margin-bottom: 8px;
    padding: 9px 12px;
  }
  .prompt-header {
    margin-bottom: 8px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    .prompt-title {
      color: #262626;
      font-weight: 600;
      font-size: 14px;
    }
  }
  .prompt-content {
    white-space: pre-wrap;
  }
}
.quick-tags-box {
  display: flex;
  align-items: center;
  margin-top: 12px;
  .quick-label {
    color: #333;
    font-weight: 500;
  }
  .ant-tag {
    cursor: pointer;
  }
}
</style>
