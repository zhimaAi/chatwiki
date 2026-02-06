<template>
  <div>
    <a-modal
      v-model:open="visible"
      :confirm-loading="loading"
      :maskClosable="false"
      :title="t('modal_title')"
      @ok="handleOk"
    >
      <a-form style="margin: 32px 0" layout="vertical" :model="formState">
        <a-form-item name="doc_auto_renew_frequency" :label="t('label_update_frequency')" required>
          <a-select v-model:value="formState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">{{ t('option_not_auto_update') }}</a-select-option>
            <a-select-option :value="2">{{ t('option_daily') }}</a-select-option>
            <a-select-option :value="3">{{ t('option_every_3_days') }}</a-select-option>
            <a-select-option :value="4">{{ t('option_every_7_days') }}</a-select-option>
            <a-select-option :value="5">{{ t('option_every_30_days') }}</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { editLibFile } from '@/api/library'
import { message } from 'ant-design-vue'

const { t } = useI18n('views.library.library-preview.components.update-frequency')
const emit = defineEmits(['ok'])
const visible = ref(false)
const loading = ref(false)
const formState = reactive({
  doc_auto_renew_frequency: 1,
  id: ''
})
const open = (data) => {
  let { doc_auto_renew_frequency, id } = data
  formState.doc_auto_renew_frequency = +doc_auto_renew_frequency
  formState.id = id
  visible.value = true
}
const handleOk = () => {
  loading.value = true
  editLibFile({
    ...formState
  })
    .then((res) => {
      emit('ok', formState.doc_auto_renew_frequency)
      visible.value = false
      message.success(t('msg_set_success'))
    })
    .finally(() => {
      loading.value = false
    })
}
defineExpose({ open })
</script>

<style lang="less" scoped>
</style>
