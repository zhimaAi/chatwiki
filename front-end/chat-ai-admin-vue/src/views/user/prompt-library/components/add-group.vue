<template>
  <div>
    <a-modal
      v-model:open="open"
      :confirm-loading="confirmLoading"
      :title="modalTitle"
      width="476px"
      @ok="handleOk"
    >
      <a-form style="margin-top: 24px" layout="vertical" ref="formRef" :model="formState">
        <a-form-item
          name="group_name"
          :label="t('group_name_label')"
          :rules="[{ required: true, message: t('group_name_required') }]"
        >
          <a-input
            v-model:value="formState.group_name"
            :placeholder="t('group_name_placeholder')"
            :maxLength="10"
          ></a-input>
        </a-form-item>
        <a-form-item name="group_desc" :label="t('group_desc_label')">
          <a-textarea
            style="height: 80px"
            :maxLength="100"
            v-model:value="formState.group_desc"
            :placeholder="t('group_desc_placeholder')"
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { savePromptLibraryGroup } from '@/api/user/index.js'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.prompt-library.components.add-group')
const emit = defineEmits(['ok'])
const formState = reactive({
  id: '',
  group_name: '',
  group_desc: ''
})
const confirmLoading = ref(false)
const modalTitle = ref(t('create_group'))
const open = ref(false)
const formRef = ref(null)

const show = (data = {}) => {
  formRef.value && formRef.value.resetFields()
  formState.group_name = data.group_name || ''
  formState.group_desc = data.group_desc || ''
  formState.id = data.id || ''
  modalTitle.value = data.id ? t('edit_group') : t('create_group')
  open.value = true
}

const handleOk = () => {
  formRef.value.validate().then(() => {
    savePromptLibraryGroup({
      ...formState
    }).then(() => {
      emit('ok')
      message.success(t('operation_success'))
      open.value = false
    })
  })
}
defineExpose({
  show
})
</script>

<style lang="less" scoped></style>
