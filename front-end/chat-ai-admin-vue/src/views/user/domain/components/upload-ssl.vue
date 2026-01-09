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
      <div class="domain">{{ t('domain_label') }}{{ domain }}</div>
      <a-form ref="formRef" :model="formState" :rules="rules" autocomplete="off" layout="vertical">
        <a-form-item :label="t('certificate_file_label')" name="ssl_certificate">
          <a-textarea :rows="5" v-model:value="formState.ssl_certificate" />
        </a-form-item>

        <a-form-item :label="t('private_key_label')" name="ssl_certificate_key">
          <a-textarea :rows="5" v-model:value="formState.ssl_certificate_key" />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const emit = defineEmits(['ok'])

const { t } = useI18n('views.user.domain.components.upload-ssl')

const props = defineProps({
  confirmLoading: {
    type: Boolean,
    default: false
  }
})

const formRef = ref()

const formState = reactive({
  id: '',
  domain: '',
  ssl_certificate: '',
  ssl_certificate_key: ''
})

const rules = {
  ssl_certificate: [{ required: true, message: t('certificate_required') }],
  ssl_certificate_key: [{ required: true, message: t('private_key_required') }]
}

const show = ref(false)

const open = (record) => {
  formState.id = record.id
  formState.domain = record.url
  formState.ssl_certificate = ''
  formState.ssl_certificate_key = ''
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
