<style lang="less" scoped>
.direct-box {
    display: flex;
    flex-direction: column;
    padding: 16px;
  }
  .direct-item {
    flex: 1;
    border: 1px solid #d9d9d9;
    padding: 15px;
    border-radius: 6px;

    .direct-title {
      font-weight: bolder;
    }

    .direct-desc {
      color: #595959;
    }

    .qa {
        display: flex;
        flex-direction: column;
        margin-top: 20px;

        .prompt-form-item-label {
            display: flex;
            align-items: center;
        }

        .qa-title {
            margin-right: 5px;
        }

        .direct-switch {
            margin-left: 20px;
        }
    }

    .similarity {
        display: flex;
        align-items: center;
        margin-top: 20px;

        .number-input-box {
            margin: 0 5px;
        }
    }
  }

  .direct-item:nth-child(2) {
    margin: 10px 0px;
  }

  .direct-item.active {
    border: 1px solid #2475fc;
  }
</style>

<template>
    <a-modal width="746px" v-model:open="show" :title="t('title_chat_mode_settings')" @ok="handleSave">
        <a-radio-group class="direct-box" v-model:value="chat_typeValue">
            <a-radio
            class="direct-item"
            :class="{active: chat_typeValue == item.direct_id}"
            :value="item.direct_id"
            v-for="item in directState"
            :key="item.direct_id"
            >
            <div class="direct-title">{{ item.direct_title }}</div>
            <div class="direct-desc">{{ item.direct_desc }}</div>
            <div class="qa" v-show="item.isQaDirectReply && chat_typeValue == item.direct_id">
                <div class="prompt-form-item-label">
                    <span class="qa-title">{{ t('label_reply_directly_when_qa_matched') }}</span>
                    <a-tooltip>
                        <template #title v-if="item.chat_type == '1'">{{ t('tooltip_library_mode_qa_direct_reply') }}</template>
                        <template #title v-else-if="item.chat_type == '3'">{{ t('tooltip_mixture_mode_qa_direct_reply') }}</template>
                        <QuestionCircleOutlined class="question-icon" />
                    </a-tooltip>
                    <div class="direct-switch">
                        <a-switch
                            v-if="item.chat_type == '1'"
                            :checkedValue="1"
                            :unCheckedValue="0"
                            v-model:checked="item.library_qa_direct_reply_switch"
                        />
                        <a-switch
                            v-else-if="item.chat_type == '3'"
                            :checkedValue="1"
                            :unCheckedValue="0"
                            v-model:checked="item.mixture_qa_direct_reply_switch"
                        />
                    </div>
                </div>
            </div>

            <div class="similarity" v-show="item.isQaDirectReply && chat_typeValue == item.direct_id">
                <span>{{ t('label_similarity_exceeds') }}</span>
                <div class="number-input-box">
                    <a-input-number
                    v-if="item.chat_type == '1'"
                    v-model:value="item.library_qa_direct_reply_score"
                    :disabled="!item.library_qa_direct_reply_switch"
                    :min="0"
                    :max="1"
                    :step="0.01"
                    />
                    <a-input-number
                    v-else-if="item.chat_type == '3'"
                    v-model:value="item.mixture_qa_direct_reply_score"
                    :disabled="!item.mixture_qa_direct_reply_switch"
                    :min="0"
                    :max="1"
                    :step="0.01"
                    />
                </div>
                <span>{{ t('label_reply_directly_when_exceeds') }}</span>
            </div>

            </a-radio>
        </a-radio-group>
    </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw, onMounted, watch } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.chat-mode.chat-edit-alert')

const emit = defineEmits(['save'])

const chat_typeValue = ref(1)

const directStateData = () => ([
  {
    direct_id: 1,
    direct_title: t('title_library_mode'),
    direct_desc:
      t('desc_library_mode'),
    chat_type: '1',
    library_qa_direct_reply_score: '0.900',
    library_qa_direct_reply_switch: 0,
    isQaDirectReply: true
  },
  {
    direct_id: 3,
    direct_title: t('title_mixture_mode'),
    direct_desc: t('desc_mixture_mode'),
    chat_type: '3',
    mixture_qa_direct_reply_score: '0.900',
    mixture_qa_direct_reply_switch: 0,
    isQaDirectReply: true
  },
  {
    direct_id: 2,
    direct_title: t('title_direct_mode'),
    direct_desc: t('desc_direct_mode'),
    chat_type: '2',
    isQaDirectReply: false
  }
])

// 定义响应式
const directState = reactive(directStateData())

const show = ref(false)

const formState = reactive({
  chat_type: '',
  library_qa_direct_reply_score: '',
  library_qa_direct_reply_switch: 'false',
  mixture_qa_direct_reply_score: '',
  mixture_qa_direct_reply_switch: 'false'
})

const handleSave = () => {

    formState.chat_type = chat_typeValue.value.toString()
    directState.forEach((item) => {
    if (item.chat_type == chat_typeValue.value) {
        if (formState.chat_type == '3') {
            formState.mixture_qa_direct_reply_switch = item.mixture_qa_direct_reply_switch == 1 ? 'true' : 'false'
            formState.mixture_qa_direct_reply_score = item.mixture_qa_direct_reply_score.toString()
        } else if (formState.chat_type == '1') {
            formState.library_qa_direct_reply_switch = item.library_qa_direct_reply_switch == 1 ? 'true' : 'false'
            formState.library_qa_direct_reply_score = item.library_qa_direct_reply_score.toString()
        }
    }
    })

    triggerSave()
    // 初始化
    Object.assign(directState, directStateData())
    show.value = false
}

const triggerSave = () => {
    emit('save', toRaw(formState))
}

onMounted(() => {
})

const open = (data) => {
    chat_typeValue.value = parseInt(data.chat_type)
    formState.chat_type = data.chat_type
    formState.library_qa_direct_reply_score = data.library_qa_direct_reply_score
    formState.library_qa_direct_reply_switch = data.library_qa_direct_reply_switch
    formState.mixture_qa_direct_reply_score = data.mixture_qa_direct_reply_score
    formState.mixture_qa_direct_reply_switch = data.mixture_qa_direct_reply_switch

    directState.forEach((item) => {
        if (item.chat_type == '3') {
            item.mixture_qa_direct_reply_score = data.mixture_qa_direct_reply_score
            item.mixture_qa_direct_reply_switch = data.mixture_qa_direct_reply_switch == 'true' ? 1 : 0
        } else if (item.chat_type == '1') {
            item.library_qa_direct_reply_score = data.library_qa_direct_reply_score
            item.library_qa_direct_reply_switch = data.library_qa_direct_reply_switch == 'true' ? 1 : 0
        }
    })

    show.value = true
}

watch(
    () => show.value,
    (val) => {
        // 
        if (val === false) {
            // 初始化
            Object.assign(directState, directStateData())
        }
    }, {
        immediate: true
    }
)

defineExpose({
  open
})
</script>
