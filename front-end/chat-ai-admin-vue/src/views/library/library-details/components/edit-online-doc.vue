<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :maskClosable="false"
      :title="t('title_add_online_data')"
      width="746px"
      @ok="handleSaveUrl"
    >
      <a-form
        class="url-add-form"
        layout="vertical"
        ref="urlFormRef"
        :model="formState"
        :rules="rules"
      >
        <a-form-item name="file_name" :label="t('label_file_name')">
          <a-input v-model:value="formState.file_name" :placeholder="t('ph_input')" :disabled="disableFileName" />
        </a-form-item>
        <a-form-item name="doc_auto_renew_frequency" :label="t('label_update_frequency')">
          <a-select v-model:value="formState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">{{ t('option_no_auto_update') }}</a-select-option>
            <a-select-option :value="2">{{ t('option_daily') }}</a-select-option>
            <a-select-option :value="3">{{ t('option_every_3_days') }}</a-select-option>
            <a-select-option :value="4">{{ t('option_every_7_days') }}</a-select-option>
            <a-select-option :value="5">{{ t('option_every_30_days') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
          v-if="formState.doc_auto_renew_frequency > 1"
          name="doc_auto_renew_minute"
          :label="t('label_update_time')"
        >
          <a-time-picker
            valueFormat="HH:mm"
            v-model:value="formState.doc_auto_renew_minute"
            format="HH:mm"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { editLibFile } from '@/api/library'
import { convertTime } from '@/utils/index'

const { t } = useI18n('views.library.library-details.components.edit-online-doc')
const emit = defineEmits(['ok'])
const formState = reactive({
  id: '',
  file_name: '',
  doc_auto_renew_frequency: 1,
  doc_auto_renew_minute: ''
})
const disableFileName = ref(false)
const rules = ref({
  file_name: [
    {
      message: t('msg_input_file_name'),
      required: true
    }
  ],
  doc_auto_renew_frequency: [
    {
      message: t('msg_select_update_frequency'),
      required: true
    }
  ]
})

const open = ref(false)
const show = (data) => {
  data = JSON.parse(JSON.stringify(data))
  formState.id = data.id || ''
  formState.file_name = data.file_name || ''
  formState.doc_auto_renew_frequency = +data.doc_auto_renew_frequency || 1
  if (data.doc_auto_renew_minute > 0) {
    formState.doc_auto_renew_minute = convertTime(data.doc_auto_renew_minute)
  }
  disableFileName.value = !!data.from_update_info
  open.value = true
}
const confirmLoading = ref(false)

const urlFormRef = ref(null)
const handleSaveUrl = () => {
  urlFormRef.value
    .validate()
    .then(() => {
      confirmLoading.value = true
      editLibFile({
        ...formState,
        doc_auto_renew_minute: formState.doc_auto_renew_minute
          ? convertTime(formState.doc_auto_renew_minute)
          : 0
      }).then(() => {
        open.value = false
        confirmLoading.value = false
        emit('ok')
      })
    })
    .catch(() => {
      confirmLoading.value = false
    })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped></style>
