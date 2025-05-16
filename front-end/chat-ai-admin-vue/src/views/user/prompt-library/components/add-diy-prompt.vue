<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :title="modalTitle"
      width="746px"
      @ok="handleOk"
    >
      <cu-scroll style="max-height: 600px">
        <a-form style="margin-top: 24px" layout="vertical" ref="formRef" :model="formState">
          <a-flex :gap="24">
            <div style="flex: 1">
              <a-form-item
                name="title"
                label="提示词标题"
                :rules="[{ required: true, message: '请输入提示词标题' }]"
              >
                <a-input
                  v-model:value="formState.title"
                  placeholder="请输入提示词标题"
                  :maxLength="10"
                ></a-input>
              </a-form-item>
            </div>
            <div style="flex: 1">
              <a-form-item name="group_id" label="分组">
                <a-select v-model:value="formState.group_id" style="width: 100%">
                  <a-select-option v-for="item in props.groupList" :value="item.id">{{
                    item.group_name
                  }}</a-select-option>
                </a-select>
              </a-form-item>
            </div>
          </a-flex>
          <a-form-item
            name="prompt"
            label="提示词"
            v-if="formState.prompt_type == 0"
            :rules="[{ required: true, message: '请输入提示词' }]"
          >
            <a-textarea
              style="height: 280px"
              v-model:value="formState.prompt"
              placeholder="请输入"
            />
          </a-form-item>
          <div class="prompt-list-box" v-else>
            <div class="header-box">
              <div class="header-left">提示词</div>
              <div class="ai-mark-box" @click="onShowAiCreateModal">
                <svg-icon name="ai-mark" />
                AI自动生成
              </div>
            </div>
            <!-- 角色 -->
            <div class="prompt-list">
              <div class="prompt-header">
                <div class="prompt-title">{{ formState.prompt_struct.role.subject }}</div>
              </div>
              <div class="prompt-content">
                <a-textarea
                  style="height: 80px"
                  :bordered="false"
                  v-model:value="formState.prompt_struct.role.describe"
                  :placeholder="placeholderMap.role"
                />
              </div>
            </div>

            <!-- 任务 -->
            <div class="prompt-list">
              <div class="prompt-header">
                <div class="prompt-title">{{ formState.prompt_struct.task.subject }}</div>
              </div>
              <div class="prompt-content">
                <a-textarea
                  style="height: 80px"
                  :bordered="false"
                  v-model:value="formState.prompt_struct.task.describe"
                  :placeholder="placeholderMap.task"
                />
              </div>
            </div>

            <!-- 要求 -->
            <div class="prompt-list">
              <div class="prompt-header">
                <div class="prompt-title">{{ formState.prompt_struct.constraints.subject }}</div>
              </div>
              <div class="prompt-content">
                <a-textarea
                  style="height: 80px"
                  :bordered="false"
                  v-model:value="formState.prompt_struct.constraints.describe"
                  :placeholder="placeholderMap.constraints"
                />
              </div>
            </div>

            <!-- 输出格式 -->
            <div class="prompt-list">
              <div class="prompt-header">
                <div class="prompt-title">{{ formState.prompt_struct.output.subject }}</div>
              </div>
              <div class="prompt-content">
                <a-textarea
                  style="height: 80px"
                  :bordered="false"
                  v-model:value="formState.prompt_struct.output.describe"
                  :placeholder="placeholderMap.output"
                />
              </div>
            </div>

            <!-- 风格 -->
            <div class="prompt-list">
              <div class="prompt-header">
                <div class="prompt-title">{{ formState.prompt_struct.tone.subject }}</div>
              </div>
              <div class="prompt-content">
                <a-textarea
                  style="height: 80px"
                  :bordered="false"
                  v-model:value="formState.prompt_struct.tone.describe"
                  :placeholder="placeholderMap.tone"
                />
              </div>
            </div>

            <!-- 自定义 -->
            <div
              class="prompt-list"
              v-for="(item, index) in formState.prompt_struct.custom"
              :key="index + item.key ? item.key : ''"
            >
              <div class="prompt-header">
                <div class="prompt-title" style="flex: 1">
                  <a-input
                    :bordered="false"
                    style="width: 100%"
                    v-model:value="item.subject"
                    placeholder="请输入主题"
                  ></a-input>
                </div>
                <div class="btn-wrapper-box">
                  <div class="hover-btn-box" @click="handleDeleteTheme(index)">
                    <CloseCircleOutlined />
                  </div>
                </div>
              </div>
              <div class="prompt-content">
                <a-textarea
                  style="height: 80px"
                  :bordered="false"
                  v-model:value="item.describe"
                  placeholder="请输入"
                />
              </div>
            </div>
            <div class="add-theme-block">
              <a-button @click="handleAddTheme" block :icon="h(PlusOutlined)">添加主题</a-button>
            </div>
          </div>
        </a-form>
      </cu-scroll>
    </a-modal>

    <AiCreatePrompt @handleAiSave="handleAiSave" ref="aiCreatePromptRef" />
  </div>
</template>

<script setup>
import { ref, reactive, h } from 'vue'
import { CloseCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { savePromptLibraryItems } from '@/api/user/index.js'
import { message } from 'ant-design-vue'
import AiCreatePrompt from './ai-create-prompt.vue'
const emit = defineEmits(['ok'])
const props = defineProps({
  groupList: {
    type: Array,
    default: () => []
  }
})

let placeholderMap = {
  role: '请输入，告知大模型要扮演的身份、职责、沟通风格等，比如“你扮演一名经验丰富的电商行业售后客服AI助手，具备良好的沟通能力和解决问题的能力。”',
  task: '请输入，比如“根据提供的知识库资料，找到对应的售后知识（每个知识点之间使用⧼-split_line-⧽进行分割），快速准确回答用户的问题。”',
  constraints:
    '请输入对大模型在回复时的要求，比如“你的回答应该使用自然的对话方式，简单直接地回答，不要解释你的答案；当用户问题没有找到相关知识点时，直接告诉用户问题暂时无法回答，不能胡编乱造，否则你将受到惩罚。”',
  output: '请输入对大模型输出格式的要求，比如“请使用markdown格式输出”',
  tone: '请告知语言风格要求，比如“专业而不失亲切，适当使用emoji增强可读性”'
}

let prompt_struct_default = {
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
  id: '',
  title: '',
  group_id: 0,
  prompt_type: 0,
  prompt: '',
  prompt_struct: JSON.parse(JSON.stringify(prompt_struct_default))
})
const confirmLoading = ref(false)
const modalTitle = ref('')
const open = ref(false)
const formRef = ref(null)

const show = (data = {}) => {
  formRef.value && formRef.value.resetFields()
  formState.title = data.title || ''
  formState.group_id = data.group_id || 0
  formState.id = data.id || ''
  formState.prompt_type = data.prompt_type
  formState.prompt = data.prompt || ''
  formState.prompt_struct = data.prompt_struct || JSON.parse(JSON.stringify(prompt_struct_default))
  modalTitle.value = `${data.id ? '编辑' : '新建'}${data.prompt_type == 0 ? '自定义' : '结构化'}提示词`
  console.log(data,formState.prompt_struct,prompt_struct_default,'===')
  open.value = true
}

const handleOk = () => {
  formRef.value.validate().then(() => {
    let parmas = {
      ...formState
    }
    if (formState.prompt_type == 0) {
      delete parmas.prompt_struct
    } else {
      parmas.prompt_struct = JSON.stringify(formState.prompt_struct)
    }
    savePromptLibraryItems({
      ...parmas
    }).then(() => {
      emit('ok')
      message.success(`${modalTitle.value}成功`)
      open.value = false
    })
  })
}

const aiCreatePromptRef = ref(null)
const onShowAiCreateModal = () => {
  aiCreatePromptRef.value.show()
}

const handleAiSave = (prompt_struct) => {
  formState.prompt_struct = prompt_struct
}

const handleAddTheme = () => {
  formState.prompt_struct.custom.push({
    subject: '',
    describe: '',
    key: Math.random() * 10000
  })
}
const handleDeleteTheme = (index) => {
  formState.prompt_struct.custom.splice(index, 1)
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.prompt-list-box {
  background: var(--09, #f2f4f7);
  border-radius: 6px;
  padding: 16px;
  .header-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #262626;
    margin-bottom: 12px;
    .ai-mark-box {
      cursor: pointer;
      display: flex;
      align-items: center;
      gap: 4px;
      color: #6524fc;
      font-size: 14px;
      line-height: 22px;
      padding: 0 8px;
      border-radius: 6px;
      cursor: pointer;
      transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
      &:hover {
        background: #e4e6eb;
      }
    }
  }
  .prompt-list {
    border: 1px solid var(--06, #d9d9d9);
    background: #fff;
    border-radius: 6px;
    margin-bottom: 8px;
    padding: 9px 12px;
    &:focus-within {
      border: 1px solid #2475fc;
      box-shadow: 0 0 0 2px rgba(5, 145, 255, 0.1);
    }
    textarea {
      padding: 0;
    }
    .ant-input {
      padding: 0;
    }
    .ant-input[disabled] {
      cursor: auto;
      color: rgba(0, 0, 0, 0.88);
    }
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
      .ant-input {
        font-weight: 600;
      }
    }
    .btn-wrapper-box {
      display: flex;
      align-items: center;
      gap: 12px;
      color: #262626;
      .swich-item {
        display: flex;
        align-items: center;
      }
    }

    .hover-btn-box {
      width: 24px;
      height: 24px;
      border-radius: 6px;
      padding: 4px;
      font-size: 16px;
      display: flex;
      align-items: center;
      color: #595959;
      cursor: pointer;
      transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
      &:hover {
        background: #e4e6eb;
      }
    }
  }
}
</style>
