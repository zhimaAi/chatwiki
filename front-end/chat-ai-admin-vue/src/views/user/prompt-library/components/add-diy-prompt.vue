<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :title="modalTitle"
      width="766px"
      @ok="handleOk"
    >
      <cu-scroll style="max-height: 600px; margin-top: 24px">
        <a-form
          :label-col="{ span: 3 }"
          :wrapper-col="{ span: 21 }"
          ref="formRef"
          :model="formState"
        >
          <a-form-item
            name="title"
            :label="t('title_label')"
            :rules="[{ required: true, message: t('title_required') }]"
          >
            <a-input
              v-model:value="formState.title"
              :placeholder="t('title_placeholder')"
              :maxLength="10"
            ></a-input>
          </a-form-item>
          <a-form-item name="group_id" :label="t('group_label')">
            <a-select v-model:value="formState.group_id" style="width: 100%">
              <a-select-option v-for="item in props.groupList" :value="item.id">{{
                item.group_name
              }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item
            name="prompt"
            :label="t('prompt_label')"
            v-if="formState.prompt_type == 0"
            :rules="[{ required: true, message: t('prompt_required') }]"
          >
            <a-textarea
              style="height: 280px"
              v-model:value="formState.prompt"
              :placeholder="t('prompt_placeholder')"
            />
          </a-form-item>
          <div class="prompt-list-box" v-else>
            <div class="header-box">
              <div class="header-left">{{ t('prompt_section_title') }}</div>
              <div class="ai-mark-box" @click="onShowAiCreateModal">
                <svg-icon name="ai-mark" />
                {{ t('ai_auto_generate') }}
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

            <div class="prompt-list">
              <div class="prompt-header">
                <div class="prompt-title">{{ formState.prompt_struct.skill.subject }}</div>
              </div>
              <div class="prompt-content">
                <a-textarea
                  style="height: 80px"
                  :bordered="false"
                  v-model:value="formState.prompt_struct.skill.describe"
                  :placeholder="placeholderMap.skill"
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
                    :placeholder="t('enter_subject')"
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
                  :placeholder="t('prompt_placeholder')"
                />
              </div>
            </div>
            <div class="add-theme-block">
              <a-button @click="handleAddTheme" block :icon="h(PlusOutlined)">{{ t('add_theme') }}</a-button>
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
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.prompt-library.components.add-diy-prompt')
const emit = defineEmits(['ok'])
const props = defineProps({
  groupList: {
    type: Array,
    default: () => []
  }
})
let placeholderMap = {
  role: t('role_placeholder'),
  task: t('task_placeholder'),
  constraints: t('constraints_placeholder'),
  output: t('output_placeholder'),
  tone: t('tone_placeholder'),
  skill: t('skill_placeholder')
}

let prompt_struct_default = {
  role: {
    subject: t('role_subject'),
    describe: ''
  },
  task: {
    subject: t('task_subject'),
    describe: ''
  },
  constraints: {
    subject: t('constraints_subject'),
    describe: ''
  },
  skill: {
    subject: t('skill_subject'),
    describe: ''
  },
  output: {
    subject: t('output_subject'),
    describe: ''
  },
  tone: {
    subject: t('tone_subject'),
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

const show = (data = {}, type) => {
  formRef.value && formRef.value.resetFields()
  formState.title = data.title || ''
  formState.group_id = data.group_id > 0 ? data.group_id : 0
  formState.id = type == 'copy' ? '' : data.id || ''
  formState.prompt_type = data.prompt_type
  formState.prompt = data.prompt || ''
  formState.prompt_struct = data.prompt_struct || JSON.parse(JSON.stringify(prompt_struct_default))
  const action = data.id ? 'edit' : 'create'
  const promptType = data.prompt_type == 0 ? 'custom' : 'structured'
  modalTitle.value = t(`${action}_${promptType}_prompt`)
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
      message.success(t('operation_success'))
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
