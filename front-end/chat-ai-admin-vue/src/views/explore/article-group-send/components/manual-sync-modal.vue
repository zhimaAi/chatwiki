<template>
  <a-modal
    :title="t('title_sync_history_articles')"
    v-model:open="open"
    width="660px"
    :confirm-loading="confirmLoading"
    @ok="handleOk"
  >
    <div class="modal-content">
      <a-form
        :model="formState"
        :label-col="{ span: 8 }"
        :wrapper-col="{ span: 12 }"
      >
        <a-form-item :label="t('label_history_publish_time')">
          <a-select v-model:value="formState.sync_type">
            <a-select-option :value="2">{{ t('within_one_year') }}</a-select-option>
            <a-select-option :value="1">{{ t('within_half_year') }}</a-select-option>
            <a-select-option :value="12">{{ t('within_three_months') }}</a-select-option>
            <a-select-option :value="11">{{ t('within_one_month') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="t('label_sync_comments')">
          <a-checkbox  v-model:checked="formState.sync_comment_switch" >
            {{ t('sync_comment_data') }}
          </a-checkbox>
        </a-form-item>
        <a-form-item :label="t('label_ai_comment')">
          <a-checkbox v-model:checked="formState.ai_comment_switch" >
            {{ t('auto_enable') }}
          </a-checkbox>
          <div class="tip-info">{{ t('msg_ai_comment_tip') }}</div>
        </a-form-item>
        <a-form-item :label="t('label_history_comment_handling')">
          <a-radio-group v-model:value="formState.replay_his_comment_switch">
            <a-radio :value="1">{{ t('yes') }}</a-radio>
            <a-radio :value="2">{{ t('no') }}</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { syncHisArticle } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.explore.article-group-send.components.manual-sync-modal')

const props = defineProps({
  appId: {
    type: [Number, String],
    default: null
  }
})
const emit = defineEmits(['updated'])
const open = ref(false)
const confirmLoading = ref(false)
const formState = reactive({
  // 同步时间类型：11：一个月，12：三个月，1:半年，2:一年
  sync_type: 12,
  // 是否同步评论，0:不同步，1:同步
  sync_comment_switch: true,
  // 开启AI评论，0:不开启，1:开启
  ai_comment_switch: true,
  // 是否处理历史评论，2:不处理，1:处理
  replay_his_comment_switch: 2
})

function show() {
  open.value = true
}

function handleOk() {
  confirmLoading.value = true
  syncHisArticle({
    app_id: props.appId,
    ...formState,
    sync_comment_switch: formState.sync_comment_switch ? 1 : 0,
    ai_comment_switch: formState.ai_comment_switch ? 1 : 0,
  }).then(res => {
    message.success(t('msg_sync_success'))
    open.value = false
    emit('updated')
  }).finally(() => {
    confirmLoading.value = false
  })
}

defineExpose({
  show
})
</script>

<style scoped>
.modal-content {
    padding-top: 16px;
    :deep(.ant-form-item) {
        margin-bottom: 16px;
    }
}
.tip-info {
  font-size: 14px;
  color: #8c8c8c;
}
</style>
