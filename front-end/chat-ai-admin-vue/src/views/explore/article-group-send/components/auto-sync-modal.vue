<template>
    <a-modal
      :title="t('title_auto_sync_history_article')"
      v-model:open="open"
      width="540px"
      :confirm-loading="confirmLoading"
      @ok="handleOk"
    >
      <div class="modal-content">
        <a-form
          :model="formState"
          :label-col="{ span: 6 }"
          :wrapper-col="{ span: 18 }"
        >
          <a-form-item :label="t('label_auto_sync')">
            <a-switch v-model:checked="formState.auto_sync_switch" :checked-children="t('switch_on')" :un-checked-children="t('switch_off')"/>
          </a-form-item>
          <a-form-item :label="t('label_sync_frequency')">
            {{ t('text_daily') }} <a-time-picker v-model:value="formState.sync_time" value-format="HH:mm" format="HH:mm"/>
          </a-form-item>
          <a-form-item :label="t('label_ai_comment')">
            <a-checkbox v-model:checked="formState.ai_comment_switch" >
              {{ t('checkbox_auto_enable') }}
            </a-checkbox>
            <div class="tip-info">{{ t('msg_ai_comment_tip') }}</div>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </template>
  
  <script setup>
  import { ref, reactive, onMounted, toRefs } from 'vue'
  import { message } from 'ant-design-vue'
  import { getSyncHisArticleTask, saveSyncHisArticleTask } from '@/api/robot'
  import { useI18n } from '@/hooks/web/useI18n'

  const { t } = useI18n('views.explore.article-group-send.components.auto-sync-modal')
  
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
    auto_sync_switch: false,
    ai_comment_switch: false,
    sync_time: '09:00'
  })

  onMounted(() => {
    loadData()
  })

  function loadData() {
    getSyncHisArticleTask({
      app_id: props.appId
    }).then(res => {
      formState.auto_sync_switch = res.data.auto_sync_switch == 1
      formState.ai_comment_switch = res.data.ai_comment_switch == 1
      formState.sync_time = res.data.sync_time || '09:00'
      emit('updated', toRefs(formState))
    })
  }

  function show() {
    open.value = true
    loadData()
  }
  
  function handleOk() {
    confirmLoading.value = true
    saveSyncHisArticleTask({
      app_id: props.appId,
      sync_time: formState.sync_time,
      ai_comment_switch: formState.ai_comment_switch ? 1 : 0,
      auto_sync_switch: formState.auto_sync_switch ? 1 : 0,
    }).then(res => {
      message.success(t('msg_save_success'))
      emit('updated', toRefs(formState))
      open.value = false
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