<template>
  <div class="prompt-box" :class="{ 'workspace-prompt-editor': props.dialogMode }" ref="promptBoxRef" :inert="isSaving ? '' : null">
    <div class="tab-header">
      <a-tabs :style="{ width: props.dialogMode ? 'auto' : '400px' }" v-model:activeKey="prompt_type">
        <a-tab-pane :key="1" :tab="t('tab_structured')"></a-tab-pane>
        <a-tab-pane :key="0" :tab="t('tab_custom')"></a-tab-pane>
      </a-tabs>
      <div v-if="props.dialogMode" class="workspace-prompt-actions">
        <a-tooltip :title="t('tooltip_upload')"><button type="button" @click="onShowUpPromptModal"><img src="@/assets/img/customer-service/basic-settings/prompt/upload.svg" alt="" />{{ t('btn_upload') }}</button></a-tooltip>
        <a-tooltip :title="t('tooltip_import')"><button type="button" @click="onShowAddPromptModal"><img src="@/assets/img/customer-service/basic-settings/prompt/import.svg" alt="" />{{ t('btn_import') }}</button></a-tooltip>
        <button type="button" class="generate-action" @click="onShowAiCreateModal"><img src="@/assets/img/customer-service/basic-settings/prompt/ai.svg" alt="" />{{ t('btn_ai_generate') }}</button>
      </div>
      <div class="opt-block" v-else-if="isEdit">
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

        <a-flex v-if="!props.dialogMode" :gap="8">
          <a-button size="small" type="primary" :loading="isSaving" @click="onSave">{{ t('btn_save') }}</a-button>
          <a-button size="small" @click="onCancel">{{ t('btn_cancel') }}</a-button>
        </a-flex>
      </div>
      <div class="opt-block" v-else>
        <span class="tip-text"
          >{{ t('msg_current_type')
          }}{{ currentPromptType == 1 ? t('msg_structured_prompt') : t('msg_custom_prompt') }}</span
        >
        <a-button size="small" @click="handleEdit">{{ t('btn_edit') }}</a-button>
      </div>
    </div>
    <div :class="{ 'workspace-editor-surface': props.dialogMode }">
    <div v-if="props.dialogMode" class="workspace-editor-toolbar">
      <span>{{ workspaceT('prompt_insert_hint') }}</span>
      <div class="workspace-output-options">
        <label>{{ t('label_output_markdown').replace(/[：:]\s*$/, '') }}<a-switch :checked="prompt_type == 1 ? outSwitch : outDiySwitch" @change="value => prompt_type == 1 ? handleChangeSwitch(value, 'outSwitch') : handleDiyChangeSwitch(value, 'outDiySwitch')" /></label>
        <label>{{ t('label_reply_image').replace(/[：:]\s*$/, '') }}<a-switch :checked="prompt_type == 1 ? imgSwitch : imgDiySwitch" @change="value => prompt_type == 1 ? handleChangeSwitch(value, 'imgSwitch') : handleDiyChangeSwitch(value, 'imgDiySwitch')" /></label>
      </div>
    </div>
    <div class="prompt-list-box" v-if="prompt_type == 1">
      <div v-if="!props.dialogMode" style="color: #8c8c8c;margin: 8px 0;">
        {{ t('msg_var_insert_tip') }}
      </div>
      <!-- 角色 -->
      <div class="prompt-list">
        <div class="prompt-header">
          <div class="prompt-title">{{ formState.prompt_struct.role.subject }}</div>
          <div class="btn-wrapper-box" v-if="isEdit">
            <a @click="handleReset('role')">{{ t('btn_reset_default') }}</a>
          </div>
        </div>
        <div class="prompt-content" :class="{'is-disabled': !isEdit}">
          <component :is="promptInputComponent" v-bind="workspaceInputProps"
            type="textarea"
            :inputStyle="inputStyle"
            :defaultValue="formState.prompt_struct.role.describe"
            :options="chatVariables"
            :noBorder="true"
            :checkAnyLevel="true"
            :placeholder="t('ph_role')"
            :ref="(el) => setAtInputRef(el, 'role', 'describe')"
            @change="(text, selectedList) => changeStructPromptValue(text, 'role', 'describe')"
          >
          </component>
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
        <div class="prompt-content" :class="{'is-disabled': !isEdit}">
          <component :is="promptInputComponent" v-bind="workspaceInputProps"
            type="textarea"
            :inputStyle="inputStyle"
            :defaultValue="formState.prompt_struct.task.describe"
            :options="chatVariables"
            :noBorder="true"
            :checkAnyLevel="true"
            :placeholder="t('ph_task')"
            :ref="(el) => setAtInputRef(el, 'task', 'describe')"
            @change="(text, selectedList) => changeStructPromptValue(text, 'task', 'describe')"
          >
          </component>
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
          <div class="prompt-content" :class="{'is-disabled': !isEdit}">

            <component :is="promptInputComponent" v-bind="workspaceInputProps"
              type="textarea"
              :inputStyle="inputStyle"
              :defaultValue="formState.prompt_struct.constraints.describe"
              :options="chatVariables"
              :noBorder="true"
              :checkAnyLevel="true"
              :placeholder="t('ph_constraints')"
              :ref="(el) => setAtInputRef(el, 'constraints', 'describe')"
              @change="(text, selectedList) => changeStructPromptValue(text, 'constraints', 'describe')"
            >
            </component>
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
          <div class="prompt-content" :class="{'is-disabled': !isEdit}">
            <component :is="promptInputComponent" v-bind="workspaceInputProps"
              type="textarea"
              :inputStyle="inputStyle"
              :defaultValue="formState.prompt_struct.skill.describe"
              :options="chatVariables"
              :noBorder="true"
              :checkAnyLevel="true"
              :placeholder="t('ph_skill')"
              :ref="(el) => setAtInputRef(el, 'skill', 'describe')"
              @change="(text, selectedList) => changeStructPromptValue(text, 'skill', 'describe')"
            >
            </component>
          </div>
        </div>

        <!-- 输出格式 -->
        <div class="prompt-list">
          <div class="prompt-header">
            <div class="prompt-title">{{ formState.prompt_struct.output.subject }}</div>
            <div class="btn-wrapper-box" v-if="isEdit">
              <div v-if="!props.dialogMode" class="swich-item">
                {{ t('label_output_markdown') }}
                <a-switch
                  @change="(val) => handleChangeSwitch(val, 'outSwitch')"
                  v-model:checked="outSwitch"
                  :checked-children="t('switch_on')"
                  :un-checked-children="t('switch_off')"
                />
              </div>
              <div v-if="!props.dialogMode" class="swich-item">
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
          <div class="prompt-content" :class="{'is-disabled': !isEdit}">
            <component :is="promptInputComponent" v-bind="workspaceInputProps"
              type="textarea"
              :inputStyle="inputStyle"
              :defaultValue="formState.prompt_struct.output.describe"
              :options="chatVariables"
              :noBorder="true"
              :checkAnyLevel="true"
              :placeholder="t('ph_output')"
              :ref="(el) => setAtInputRef(el, 'output', 'describe')"
              @change="(text, selectedList) => changeStructPromptValue(text, 'output', 'describe')"
            >
            </component>
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
          <div class="prompt-content" :class="{'is-disabled': !isEdit}">
            <component :is="promptInputComponent" v-bind="workspaceInputProps"
              type="textarea"
              :inputStyle="inputStyle"
              :defaultValue="formState.prompt_struct.tone.describe"
              :options="chatVariables"
              :noBorder="true"
              :checkAnyLevel="true"
              :placeholder="t('ph_tone')"
              :ref="(el) => setAtInputRef(el, 'tone', 'describe')"
              @change="(text, selectedList) => changeStructPromptValue(text, 'tone', 'describe')"
            >
            </component>
          </div>
        </div>

        <!-- 自定义 -->
        <div
          class="prompt-list"
          v-for="(item, index) in formState.prompt_struct.custom"
          :key="index"
        >
          <div class="prompt-header">
            <div class="prompt-title prompt-content" :class="{'is-disabled': !isEdit}" style="flex: 1">
              <component :is="promptInputComponent" v-bind="workspaceInputProps"
                inputStyle="height: 32px"
                :defaultValue="item.subject"
                :options="chatVariables"
                :noBorder="true"
                :checkAnyLevel="true"
                :placeholder="t('ph_input_subject')"
                :ref="(el) => setAtInputRef(el, 'custom', 'subject_' + index)"
                @change="(text, selectedList) => changeCustomPromptValue(text, index, 'subject')"
              >
              </component>
            </div>
            <div class="btn-wrapper-box" v-if="isEdit">
              <div class="hover-btn-box" @click="handleDeleteTheme(index)">
                <CloseCircleOutlined />
              </div>
            </div>
          </div>
          <div class="prompt-content" :class="{'is-disabled': !isEdit}">
            <component :is="promptInputComponent" v-bind="workspaceInputProps"
              type="textarea"
              :inputStyle="inputStyle"
              :defaultValue="item.describe"
              :options="chatVariables"
              :noBorder="true"
              :checkAnyLevel="true"
              :placeholder="t('ph_input')"
              :ref="(el) => setAtInputRef(el, 'custom', 'describe_' + index)"
              @change="(text, selectedList) => changeCustomPromptValue(text, index, 'describe')"
            >
            </component>
          </div>
        </div>
        <div class="add-theme-block" v-if="isEdit">
          <a-button @click="handleAddTheme" block :icon="props.dialogMode ? undefined : h(PlusOutlined)"><template v-if="props.dialogMode" #icon><img src="@/assets/img/customer-service/basic-settings/prompt/add-theme.svg" alt="" /></template>{{
            t('btn_add_theme')
          }}</a-button>
        </div>
      </template>
      <div v-if="!props.dialogMode" class="show-more-block">
        <div class="btn-item" @click="handleShowMore">
          <template v-if="!isHide">{{ t('btn_expand') }}<DownOutlined /></template>
          <template v-else>{{ t('btn_collapse') }}<UpOutlined /></template>
        </div>
      </div>
    </div>
    <div class="diy-prompt-box" :class="{ 'is-disabled': !isEdit }" v-else>
      <a-flex align="center" justify="space-between" v-if="isEdit && !props.dialogMode" style="margin: 8px 0">
        <div style="color: #8c8c8c">
          {{ t('msg_var_insert_tip_custom') }}
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

      <component :is="promptInputComponent" v-bind="workspaceInputProps"
        class="diy-prompt-input"
        type="textarea"
        :defaultValue="formState.prompt"
        :options="chatVariables"
        :cardOptions="miniCardOptions"
        :getPopupContainer="getAtInputPopupContainer"
        :enableCardInsert="true"
        :disabled="true"
        :checkAnyLevel="true"
        ref="atinputRef"
        @change="(text) => changeValue(text)"
      >
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
        <template #dropdown="{ options, select, cardOptions, selectCard, close }">
          <SlashPopupPanel
            :options="options"
            :cardOptions="cardOptions"
            :select="select"
            :selectCard="selectCard"
            :close="close"
            @add-tag="handleOpenModal"
            @add-card="handleAddMiniCard"
            @edit-card="handleEditMiniCard"
          />
        </template>
      </component>
    </div>
    </div>
    <component :is="props.dialogMode ? WorkspaceAiPrompt : AiCreatePrompt" @handleAiSave="handleAiSave" ref="aiCreatePromptRef" />
    <ImportPrompt @ok="handleSavePrompt" ref="importPromptRef" />
    <UploadPrompt ref="uploadPromptRef" />
    <AddVariableModal ref="addVariableModalRef" :compact="props.dialogMode" @ok="handleGetChatVariables" />
    <MiniCardFormModal ref="miniCardFormModalRef" @ok="() => miniCardStore.fetchCardList(true)" />
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
import { message } from 'ant-design-vue'
import AiCreatePrompt from './ai-create-prompt.vue'
import ImportPrompt from '@/components/import-prompt/index.vue'
import UploadPrompt from '@/components/import-prompt/upload-prompt.vue'
import AtInput from './at-input/index.vue'
import SlashPopupPanel from './slash-popup-panel.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useMiniCardStore } from '@/stores/modules/mini-card'
import { useI18n } from '@/hooks/web/useI18n'
import AddVariableModal from './variable-setting/add-variable-modal.vue'
import MiniCardFormModal from '@/components/mini-card/mini-card-form-modal.vue'
import PromptVariableInput from '@/views/customer-service/basic-settings/components/prompt-variable-input.vue'
import WorkspaceAiPrompt from '@/views/customer-service/basic-settings/components/ai-prompt-dialog.vue'
import { getPromptOptions } from '@/views/customer-service/basic-settings/prompt-tokens'

const { t } = useI18n('views.robot.robot-config.basic-config.components.system-prompt-words')
const { t: workspaceT } = useI18n('views.customer-service.index')
const robotStore = useRobotStore()
const miniCardStore = useMiniCardStore()
const isEdit = ref(false)
const isSaving = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')
const props = defineProps({
  dialogMode: { type: Boolean, default: false },
  robotList: {
    type: Array,
    default: () => []
  }
})
if (props.dialogMode) isEdit.value = true

const inputStyle = ref('max-height: 130px;')
const PromptDefaultReplyMarkdown = t('prompt_default_reply_markdown')
const PromptDefaultAnswerImage = t('prompt_default_answer_image')

const prompt_struct_default = ref({})

const outSwitch = ref(false)
const imgSwitch = ref(false)

const outDiySwitch = ref(false)
const imgDiySwitch = ref(false)

const promptBoxRef = ref(null)
const getAtInputPopupContainer = () => props.dialogMode ? document.body : promptBoxRef.value || document.body
const promptInputComponent = computed(() => props.dialogMode ? PromptVariableInput : AtInput)
const workspaceInputProps = computed(() => props.dialogMode ? { cardOptions: miniCardOptions.value, getPopupContainer: () => document.body, onAddTag: handleOpenModal, onAddCard: handleAddMiniCard, onEditCard: handleEditMiniCard } : {})
const workspaceTokenOptions = computed(() => getPromptOptions(robotStore.chatVariables, miniCardStore.cardList))


const atInputRefs = reactive({})
const setAtInputRef = (el, name, index) => {
  if (el) {
    let key = `at_input_${name}_${index}`
    atInputRefs[key] = el
  }
}

const atinputRef = ref(null)

const systemVariables = [
  {
    id: 'system_current_time',
    label: '当前时间-系统变量',
    value: '【system_var:current_time】',
    typ: '系统变量'
  },
  {
    id: 'system_current_date',
    label: '当前日期-系统变量',
    value: '【system_var:current_date】',
    typ: '系统变量'
  }
]

const chatVariables = computed(() => {
  if (props.dialogMode) return workspaceTokenOptions.value.variables
  let list = robotStore.chatVariables
  const customVariables = list.map((item) => {
    let label = item.variable_key + '（' + item.variable_name + '）'
    let value = `【chat_variable:${item.variable_key}】`
    return {
      ...item,
      label,
      value
    }
  })
  return [...systemVariables, ...customVariables]
})

// 小程序卡片选项，映射为 AtInput 可用的 options 格式
const miniCardOptions = computed(() => {
  if (props.dialogMode) return workspaceTokenOptions.value.cards
  let list = miniCardStore.cardList
  return list.map((item) => {
    let label = '小程序卡片：' + item.title
    let value = `【mini_card:${item.id}】`
    return {
      ...item,
      label,
      value
    }
  })
})

// 初始化加载小程序卡片列表
Promise.resolve(miniCardStore.fetchCardList()).catch(() => {})

watch(
  () => chatVariables,
  () => {
    if (props.dialogMode) return
    setTimeout(() => {
      atinputRef.value && atinputRef.value.refresh()
      for(let key in atInputRefs){
        atInputRefs[key] && atInputRefs[key].refresh()
      }
    }, 500)
  },
  {
    deep: true
  }
)

const changeValue = (text) => {
  formState.prompt = text
}

const changeStructPromptValue = (text, key1, key2) =>{
  formState.prompt_struct[key1][key2] = text
}

const changeCustomPromptValue = (text, index, key) => {
  formState.prompt_struct.custom[index][key] = text
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

const isHide = ref(props.dialogMode ? true : localStorage.getItem('system_prompt_words_ishide') == '1')

const handleShowMore = () => {
  isHide.value = !isHide.value
  if (!props.dialogMode) localStorage.setItem('system_prompt_words_ishide', isHide.value ? 1 : 0)
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

const savePrompt = async (data) => {
  // 工作区弹窗只收集草稿，AI 生成和导入也要等待外层统一保存。
  if (props.dialogMode) {
    if (data.prompt !== undefined) formState.prompt = data.prompt
    if (data.prompt_struct) formState.prompt_struct = JSON.parse(JSON.stringify(data.prompt_struct))
    prompt_type.value = data.prompt_type
    return true
  }
  if (isSaving.value) return
  isSaving.value = true
  try {
    const success = await updateRobotInfo(data)
    if (success) {
      isEdit.value = false
    }
  } finally {
    isSaving.value = false
  }
}

const onSave = () => {
  return savePrompt({
    ...toRaw(formState),
    prompt_type: prompt_type.value
  })
}

const handleAiSave = (prompt_struct) => {
  if (isSaving.value) return
  prompt_type.value = 1
  return savePrompt({
    prompt_type: prompt_type.value,
    prompt: formState.prompt,
    prompt_struct
  })
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
  if (isSaving.value) return
  if (item.prompt_type == 1) {
    prompt_type.value = 1
    formState.prompt_struct = item.prompt_struct
  } else {
    prompt_type.value = 0
    formState.prompt = item.prompt
  }
  return onSave()
}

const miniCardFormModalRef = ref(null)

const handleAddMiniCard = () => {
  miniCardFormModalRef.value.show()
}

const handleEditMiniCard = (card) => {
  miniCardFormModalRef.value.edit(card)
}

const addVariableModalRef = ref(null)

const handleOpenModal = () => {
  if(robotStore.chatVariables.length >= 10){
    return message.error(t('msg_max_variables'))
  }
  addVariableModalRef.value.show()
}

const handleGetChatVariables = () => {
  robotStore.fetchChatVariables()
}
defineExpose({
  getDraft: () => {
    if (props.dialogMode) {
      atinputRef.value?.updateValue?.()
      for (const input of Object.values(atInputRefs)) input?.updateValue?.()
    }
    return JSON.parse(JSON.stringify({ prompt: formState.prompt, prompt_struct: formState.prompt_struct, prompt_type: prompt_type.value }))
  }
})

</script>

<style lang="less" scoped>
.prompt-box.workspace-prompt-editor {
  padding: 0;
  background: #fff;
  border-radius: 0;
  .tab-header { position: sticky; top: 0; z-index: 2; min-height: 40px; gap: 16px; padding-right: 40px; background: #fff; }
  :deep(.ant-tabs-nav) { margin: 0; }
  :deep(.ant-tabs-tab) { min-height: 40px; padding: 0 0 8px; color: #5a5a5a; font-size: 14px; line-height: 22px; }
  :deep(.ant-tabs-tab + .ant-tabs-tab) { margin-left: 24px; }
  :deep(.ant-tabs-tab-active .ant-tabs-tab-btn) { color: #146ef5; font-weight: 500; }
  :deep(.ant-tabs-ink-bar) { background: #146ef5; }
  .workspace-prompt-actions { display: flex; align-items: center; justify-content: flex-end; gap: 4px; margin-left: auto; button { display: inline-flex; align-items: center; justify-content: center; gap: 6px; height: 32px; padding: 0 8px; border: 0; border-radius: 6px; background: #fff; color: #5a5a5a; font-size: 14px; cursor: pointer; white-space: nowrap; } button:hover { background: #f5f5f5; } .generate-action { color: #146ef5; } }
  .workspace-editor-surface { margin: 32px 19px 0 8px; padding: 8px; border-radius: 8px; background: #f5f5f5; }
  .workspace-editor-toolbar { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 8px; min-height: 32px; margin-bottom: 8px; color: #7a7a7a; font-size: 14px; line-height: 22px; }
  .workspace-output-options { display: flex; align-items: center; gap: 12px; color: #222; label { display: inline-flex; align-items: center; gap: 8px; } }
  :deep(.ant-switch) { min-width: 36px; width: 36px; height: 20px; line-height: 20px; .ant-switch-handle { width: 16px; height: 16px; top: 2px; } .ant-switch-inner { display: none; } &.ant-switch-checked { background: #146ef5; .ant-switch-handle { inset-inline-start: calc(100% - 18px); } } }
  .prompt-list { padding: 12px; border-color: #e8e8e8; border-radius: 6px; }
  .prompt-list:last-child { margin-bottom: 0; }
  .prompt-header { min-height: 32px; margin-bottom: 8px; }
  .prompt-header .prompt-title { line-height: 22px; }
  .btn-wrapper-box a { color: #7a7a7a; font-size: 14px; }
  .prompt-content :deep(.j-mention.type-textarea) { min-height: 66px; line-height: 22px; }
  .diy-prompt-input :deep(.j-mention.type-textarea) { min-height: 302px; height: 302px; line-height: 22px; }
  .diy-prompt-input :deep(.mention-input-warpper) { padding: 8px 12px; border-color: #e8e8e8; }
  .add-theme-block { margin-top: 16px; :deep(.ant-btn) { height: 32px; border-color: #e8e8e8; border-style: dashed; border-radius: 6px; color: #7a7a7a; font-size: 14px; } }
  @media (max-width: 760px) { .tab-header { flex-wrap: wrap; } .workspace-prompt-actions { flex-wrap: wrap; margin-left: 0; } .workspace-editor-surface { margin-right: 8px; } }
}
</style>
<style lang="less" scoped>
.prompt-box {
  position: relative;
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
  .prompt-content{
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
  }
}

.diy-prompt-box {
  position: relative;

  .diy-prompt-input {
    &::v-deep(.j-mention.type-textarea) {
      height: 126px;
      resize: vertical;
      overflow: auto;

      &::-webkit-scrollbar {
        width: 8px;
        height: 8px;
      }
    }
  }

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
