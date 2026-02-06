<template>
  <div class="prompt-box">
    <div class="tab-header">
      <a-tabs style="width: 400px;" v-model:activeKey="prompt_type">
        <a-tab-pane :key="1" :tab="t('tab_structured')"></a-tab-pane>
        <a-tab-pane :key="0" :tab="t('tab_custom')"></a-tab-pane>
      </a-tabs>
      <div class="opt-block" v-if="isEdit">
        <a-tooltip>
          <template #title>{{ t('tooltip_upload') }}</template>
          <div class="hover-btn-box" style="color: #6524fc" @click="onShowUpPromptModal">
            <ToTopOutlined />{{ t('btn_upload') }}
          </div>
        </a-tooltip>
        <a-tooltip>
          <template #title>{{ t('tooltip_import') }}</template>
          <div class="hover-btn-box" style="color: #6524fc" @click="onShowAddPromptModal">
            <DownloadOutlined />{{ t('btn_import') }}
          </div>
        </a-tooltip>

        <div class="ai-mark-box hover-btn-box" @click="onShowAiCreateModal">
          <svg-icon name="ai-mark" />
          {{ t('btn_ai_generate') }}
        </div>

        <a-flex :gap="8">
          <a-button size="small" type="primary" @click="onSave">{{ t('btn_save') }}</a-button>
          <a-button size="small" @click="onCancel">{{ t('btn_cancel') }}</a-button>
        </a-flex>
      </div>
      <div class="opt-block" v-else>
        <span class="tip-text"
          >{{ t('msg_current_type') }}{{ currentPromptType == 1 ? t('msg_structured_prompt') : t('msg_custom_prompt') }}</span
        >
        <a-button size="small" @click="handleEdit">{{ t('btn_edit') }}</a-button>
      </div>
    </div>
    <div class="prompt-list-box" v-if="prompt_type == 1">
      <!-- 角色 -->
      <div class="prompt-list">
        <div class="prompt-header">
          <div class="prompt-title">{{ formState.prompt_struct.role.subject }}</div>
          <div class="btn-wrapper-box" v-if="isEdit">
            <a @click="handleReset('role')">{{ t('btn_reset_default') }}</a>
          </div>
        </div>
        <div class="prompt-content">
          <a-textarea
            :bordered="false"
            :disabled="!isEdit"
            v-model:value="formState.prompt_struct.role.describe"
            :placeholder="t('ph_role')"
          />
        </div>
      </div>
      <!-- 任务 -->
      <div class="prompt-list">
        <div class="prompt-header">
          <div class="prompt-title">{{ formState.prompt_struct.task.subject }}</div>
          <div class="btn-wrapper-box" v-if="isEdit">
            <a @click="handleReset('task')">{{ t('btn_reset_default') }}</a>
          </div>
        </div>
        <div class="prompt-content">
          <a-textarea
            :bordered="false"
            :disabled="!isEdit"
            v-model:value="formState.prompt_struct.task.describe"
            :placeholder="t('ph_task')"
          />
        </div>
      </div>
      <template v-if="isHide">
        <!-- 要求 -->
        <div class="prompt-list">
          <div class="prompt-header">
            <div class="prompt-title">{{ formState.prompt_struct.constraints.subject }}</div>
            <div class="btn-wrapper-box" v-if="isEdit">
              <a @click="handleReset('constraints')">{{ t('btn_reset_default') }}</a>
            </div>
          </div>
          <div class="prompt-content">
            <a-textarea
              :bordered="false"
              :disabled="!isEdit"
              v-model:value="formState.prompt_struct.constraints.describe"
              :placeholder="t('ph_constraints')"
              style="min-height: 130px"
            />
          </div>
        </div>
        <!-- 技能 -->
        <div class="prompt-list">
          <div class="prompt-header">
            <div class="prompt-title">{{ formState.prompt_struct.skill.subject }}</div>
            <div class="btn-wrapper-box" v-if="isEdit">
              <a @click="handleImportSkill()">{{ t('btn_import_skill') }}</a>
            </div>
          </div>
          <div class="prompt-content">
            <a-textarea
              :bordered="false"
              :disabled="!isEdit"
              v-model:value="formState.prompt_struct.skill.describe"
              :placeholder="t('ph_skill')"
              style="min-height: 80px"
            />
          </div>
        </div>

        <!-- 输出格式 -->
        <div class="prompt-list">
          <div class="prompt-header">
            <div class="prompt-title">{{ formState.prompt_struct.output.subject }}</div>
            <div class="btn-wrapper-box" v-if="isEdit">
              <div class="swich-item">
                {{ t('label_output_markdown') }}
                <a-switch
                  @change="(val) => handleChangeSwitch(val, 'outSwitch')"
                  v-model:checked="outSwitch"
                  :checked-children="t('switch_on')"
                  :un-checked-children="t('switch_off')"
                />
              </div>
              <div class="swich-item">
                {{ t('label_reply_image') }}
                <a-switch
                  @change="(val) => handleChangeSwitch(val, 'imgSwitch')"
                  v-model:checked="imgSwitch"
                  :checked-children="t('switch_on')"
                  :un-checked-children="t('switch_off')"
                />
              </div>
              <a @click="handleReset('output')">{{ t('btn_reset_default') }}</a>
            </div>
          </div>
          <div class="prompt-content">
            <a-textarea
              :bordered="false"
              :disabled="!isEdit"
              v-model:value="formState.prompt_struct.output.describe"
              :placeholder="t('ph_output')"
            />
          </div>
        </div>
        <!-- 风格 -->
        <div class="prompt-list">
          <div class="prompt-header">
            <div class="prompt-title">{{ formState.prompt_struct.tone.subject }}</div>
            <div class="btn-wrapper-box" v-if="isEdit">
              <a @click="handleReset('tone')">{{ t('btn_reset_default') }}</a>
            </div>
          </div>
          <div class="prompt-content">
            <a-textarea
              :bordered="false"
              :disabled="!isEdit"
              v-model:value="formState.prompt_struct.tone.describe"
              :placeholder="t('ph_tone')"
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
                :disabled="!isEdit"
                style="width: 100%"
                v-model:value="item.subject"
                :placeholder="t('ph_input_subject')"
              ></a-input>
            </div>
            <div class="btn-wrapper-box" v-if="isEdit">
              <div class="hover-btn-box" @click="handleDeleteTheme(index)">
                <CloseCircleOutlined />
              </div>
            </div>
          </div>
          <div class="prompt-content">
            <a-textarea
              :bordered="false"
              :disabled="!isEdit"
              v-model:value="item.describe"
              :placeholder="t('ph_input')"
            />
          </div>
        </div>
        <div class="add-theme-block" v-if="isEdit">
          <a-button @click="handleAddTheme" block :icon="h(PlusOutlined)">{{ t('btn_add_theme') }}</a-button>
        </div>
      </template>
      <div class="show-more-block">
        <div class="btn-item" @click="handleShowMore">
          <template v-if="!isHide">{{ t('btn_expand') }}<DownOutlined /></template>
          <template v-else>{{ t('btn_collapse') }}<UpOutlined /></template>
        </div>
      </div>
    </div>
    <div class="diy-prompt-box" :class="{ 'is-disabled': !isEdit }" v-else>
      <a-flex align="center" justify="space-between" v-if="isEdit" style="margin: 8px 0">
        <div style="color: #8c8c8c">
          {{ t('msg_var_insert_tip') }}
        </div>
        <div class="diy-switch-box">
          <div class="swich-item">
            {{ t('label_output_markdown') }}
            <a-switch
              @change="(val) => handleDiyChangeSwitch(val, 'outDiySwitch')"
              v-model:checked="outDiySwitch"
              :checked-children="t('switch_on')"
              :un-checked-children="t('switch_off')"
            />
          </div>
          <div class="swich-item">
            {{ t('label_reply_image') }}
            <a-switch
              @change="(val) => handleDiyChangeSwitch(val, 'imgDiySwitch')"
              v-model:checked="imgDiySwitch"
              :checked-children="t('switch_on')"
              :un-checked-children="t('switch_off')"
            />
          </div>
        </div>
      </a-flex>

      <AtInput
        type="textarea"
        inputStyle="height: 126px"
        :defaultValue="formState.prompt"
        :options="chatVariables"
        :disabled="true"
        :checkAnyLevel="true"
        ref="atinputRef"
        @change="(text, selectedList) => changeValue(text, selectedList)"
      >
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
      </AtInput>
    </div>
    <AiCreatePrompt @handleAiSave="handleAiSave" ref="aiCreatePromptRef" />
    <ImportPrompt @ok="handleSavePrompt" ref="importPromptRef" />
    <UploadPrompt ref="uploadPromptRef" />
  </div>
</template>
<script setup>
import { ref, reactive, computed, inject, toRaw, watch, h } from 'vue'
import {
  PlusOutlined,
  CloseCircleOutlined,
  DownOutlined,
  UpOutlined,
  ToTopOutlined,
  DownloadOutlined
} from '@ant-design/icons-vue'
import AiCreatePrompt from './ai-create-prompt.vue'
import ImportPrompt from '@/components/import-prompt/index.vue'
import UploadPrompt from '@/components/import-prompt/upload-prompt.vue'
import AtInput from './at-input/index.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.system-prompt-words')
const robotStore = useRobotStore()
const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')
const props = defineProps({
  robotList: {
    type: Array,
    default: () => []
  }
})

const PromptDefaultReplyMarkdown = t('prompt_default_reply_markdown')
const PromptDefaultAnswerImage = t('prompt_default_answer_image')

const prompt_struct_default = ref({})

const outSwitch = ref(false)
const imgSwitch = ref(false)

const outDiySwitch = ref(false)
const imgDiySwitch = ref(false)

const atinputRef = ref(null)

const chatVariables = computed(() => {
  let list = robotStore.chatVariables
  return list.map((item) => {
    let label = item.variable_key + '（' + item.variable_name + '）'
    let value = `【chat_variable:${item.variable_key}】`
    return {
      ...item,
      label,
      value
    }
  })
})

watch(
  () => chatVariables,
  () => {
    setTimeout(() => {
      atinputRef.value && atinputRef.value.refresh()
    }, 500)
  },
  {
    deep: true
  }
)

const changeValue = (text, selectedList) => {
  formState.prompt = text
}

const formState = reactive({
  prompt: '',
  prompt_struct: {
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
})

const isHide = ref(localStorage.getItem('system_prompt_words_ishide') == '1')

const handleShowMore = () => {
  isHide.value = !isHide.value
  localStorage.setItem('system_prompt_words_ishide', isHide.value ? 1 : 0)
}

const checkStringContains = (str, searchStr) => {
  // 使用 includes 方法判断
  return str.includes(searchStr)
}

const prompt_type = ref(robotInfo.prompt_type)
const currentPromptType = computed(() => robotInfo.prompt_type)

watch(
  () => robotInfo,
  () => {
    formState.prompt = robotInfo.prompt
    formState.prompt_struct = robotInfo.prompt_struct
    prompt_struct_default.value = robotInfo.prompt_struct_default
    // prompt_type.value = robotInfo.prompt_type
  },
  { immediate: true, deep: true }
)

watch(
  () => formState.prompt_struct.output.describe,
  (val) => {
    outSwitch.value = checkStringContains(val, PromptDefaultReplyMarkdown)
    imgSwitch.value = checkStringContains(val, PromptDefaultAnswerImage)
  },
  { immediate: true, deep: true }
)

watch(
  () => formState.prompt,
  (val) => {
    outDiySwitch.value = checkStringContains(val, PromptDefaultReplyMarkdown)
    imgDiySwitch.value = checkStringContains(val, PromptDefaultAnswerImage)
  },
  { immediate: true, deep: true }
)

const handleChangeSwitch = (val, key) => {
  if (val) {
    if (key == 'outSwitch') {
      if (formState.prompt_struct.output.describe) {
        formState.prompt_struct.output.describe += '\n' + PromptDefaultReplyMarkdown
      } else {
        formState.prompt_struct.output.describe += PromptDefaultReplyMarkdown
      }
    }
    if (key == 'imgSwitch') {
      if (formState.prompt_struct.output.describe) {
        formState.prompt_struct.output.describe += '\n' + PromptDefaultAnswerImage
      } else {
        formState.prompt_struct.output.describe += PromptDefaultAnswerImage
      }
    }
  } else {
    if (key == 'outSwitch') {
      formState.prompt_struct.output.describe = formState.prompt_struct.output.describe.replace(
        new RegExp(PromptDefaultReplyMarkdown, 'g'),
        ''
      )
    }
    if (key == 'imgSwitch') {
      formState.prompt_struct.output.describe = formState.prompt_struct.output.describe.replace(
        new RegExp(PromptDefaultAnswerImage, 'g'),
        ''
      )
    }
    formState.prompt_struct.output.describe = formState.prompt_struct.output.describe.trim()
  }
}

const handleDiyChangeSwitch = (val, key) => {
  if (val) {
    if (key == 'outDiySwitch') {
      if (formState.prompt) {
        formState.prompt += '\n' + PromptDefaultReplyMarkdown
      } else {
        formState.prompt += PromptDefaultReplyMarkdown
      }
    }
    if (key == 'imgDiySwitch') {
      if (formState.prompt) {
        formState.prompt += '\n' + PromptDefaultAnswerImage
      } else {
        formState.prompt += PromptDefaultAnswerImage
      }
    }
  } else {
    if (key == 'outDiySwitch') {
      formState.prompt = formState.prompt.replace(new RegExp(PromptDefaultReplyMarkdown, 'g'), '')
    }
    if (key == 'imgDiySwitch') {
      formState.prompt = formState.prompt.replace(new RegExp(PromptDefaultAnswerImage, 'g'), '')
    }
    formState.prompt = formState.prompt.trim()
  }
}

const handleReset = (key) => {
  formState.prompt_struct[key].describe = prompt_struct_default.value[key].describe
}

const handleImportSkill = () => {
  // {{ t('comment_import_skill') }}
  let work_flow_ids = robotInfo.work_flow_ids.split(',')
  let selectRobotList = props.robotList.filter((item) => {
    return work_flow_ids.includes(item.id)
  })
  let skill_str = []
  skill_str = selectRobotList
    .map((item) => {
      return `- ${item.robot_name} ${item.robot_intro}`
    })
    .join('\n')
  formState.prompt_struct['skill'].describe = skill_str
}

const handleAddTheme = () => {
  formState.prompt_struct.custom.push({
    subject: '',
    describe: '',
    key: Math.random() * 10000
  })
  prompt_type.value = 1
}

const handleDeleteTheme = (index) => {
  formState.prompt_struct.custom.splice(index, 1)
}

const onSave = () => {
  updateRobotInfo({
    ...toRaw(formState),
    prompt_type: prompt_type.value
  })
  isEdit.value = false
}

const handleAiSave = (prompt_struct) => {
  prompt_type.value = 1
  updateRobotInfo({
    prompt_type: prompt_type.value,
    prompt: formState.prompt,
    prompt_struct
  })
  isEdit.value = false
}

const handleEdit = () => {
  // prompt_type.value = robotInfo.prompt_type
  prompt_type.value == 0 && atinputRef.value && atinputRef.value.refresh()
  isEdit.value = true
}

const onCancel = () => {
  isEdit.value = false
}

const aiCreatePromptRef = ref(null)
const onShowAiCreateModal = () => {
  aiCreatePromptRef.value.show()
}

const importPromptRef = ref(null)

const onShowAddPromptModal = () => {
  importPromptRef.value.show()
}

const uploadPromptRef = ref(null)
const onShowUpPromptModal = () => {
  uploadPromptRef.value.show({
    ...formState,
    prompt_type: prompt_type.value
  })
}

const handleSavePrompt = (item) => {
  if (item.prompt_type == 1) {
    prompt_type.value = 1
    formState.prompt_struct = item.prompt_struct
  } else {
    prompt_type.value = 0
    formState.prompt = item.prompt
  }
  onSave()
}
</script>

<style lang="less" scoped>
.prompt-box {
  border-radius: 6px;
  background: var(--09, #f2f4f7);
  padding: 16px;
  .tab-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .opt-block {
    display: flex;
    gap: 8px;
    .tip-text {
      color: #8c8c8c;
      line-height: 22px;
      font-size: 12px;
    }
    .hover-btn-box {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 0 6px;
      width: fit-content;
      cursor: pointer;
      &:hover {
        background: #e4e6eb;
        border-radius: 6px;
      }
    }
  }
  .ai-mark-box {
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 4px;
    color: #6524fc;
    font-size: 14px;
    line-height: 22px;
  }
  ::v-deep(.ant-tabs-nav-wrap) {
    padding-left: 0;
  }
  ::v-deep(.ant-tabs-nav) {
    margin: 0;
    .ant-tabs-tab {
      padding-top: 0;
    }
    &::before {
      border: none;
    }
  }
}

.show-more-block {
  padding-top: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #595959;
  .btn-item {
    width: fit-content;
    display: flex;
    align-items: center;
    gap: 4px;
    height: 24px;
    padding: 0 8px;
    border-radius: 6px;
    cursor: pointer;
    &:hover {
      background: #e4e6eb;
    }
  }
}

.prompt-list-box {
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

.diy-prompt-box {
  position: relative;
  &.is-disabled {
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(255, 255, 255, 0.25);
      z-index: 1;
    }
  }
  .diy-switch-box {
    display: flex;
    align-items: center;
    gap: 16px;
    color: #595959;
    font-size: 14px;
    .swich-item {
      display: flex;
      align-items: center;
    }
  }
  .ant-input {
    padding-top: 24px;
  }
  .ant-input[disabled] {
    padding-top: 4px;
    background-color: #fff;
    color: rgba(0, 0, 0, 0.88);
    cursor: auto;
  }
}
</style>