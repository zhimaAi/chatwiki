<style lang="less" scoped>
.upload-ssl-form {
  padding-top: 16px;
  .domain {
    margin-bottom: 16px;
    color: #333;
  }
}
</style>

<template>
  <a-modal
    v-model:open="show"
    :title="t('title')"
    @ok="handleOk"
    :confirmLoading="props.confirmLoading"
  >
    <div class="upload-ssl-form">
      <a-form ref="formRef" :model="formState" :rules="rules" autocomplete="off" layout="vertical">
        <a-form-item :label="t('file_name_label')" name="file_name">
          <a-input v-model:value="formState.file_name" />
        </a-form-item>

        <a-form-item :label="t('file_content_label')" name="file_content">
          <a-textarea :rows="5" v-model:value="formState.file_content" />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const emit = defineEmits(['ok'])

const { t } = useI18n('views.user.domain.components.upload-validation-file')

const props = defineProps({
  confirmLoading: {
    type: Boolean,
    default: false
  }
})

const formRef = ref()

const formState = reactive({
  id: '',
  file_name: '',
  file_content: ''
})

const rules = {
  file_name: [{ required: true, message: t('file_name_required') }],
  file_content: [{ required: true, message: t('file_content_required') }]
}

const show = ref(false)

const open = (record) => {
  formState.id = record.id
  formState.file_name = ''
  formState.file_content = ''
  show.value = true
}

const close = () => {
  show.value = false
}

const handleOk = () => {
  formRef.value
    .validate()
    .then(() => {
      emit('ok', toRaw(formState))
    })
    .catch((error) => {
      console.log('error', error)
    })
}

defineExpose({
  open,
  close
})
</script>
