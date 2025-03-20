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
    title="上传证书"
    @ok="handleOk"
    :confirmLoading="props.confirmLoading"
  >
    <div class="upload-ssl-form">
      <div class="domain">域名：{{ domain }}</div>
      <a-form ref="formRef" :model="formState" :rules="rules" autocomplete="off" layout="vertical">
        <a-form-item label="证书文件(.crt或.pem)" name="ssl_certificate">
          <a-textarea :rows="5" v-model:value="formState.ssl_certificate" />
        </a-form-item>

        <a-form-item label="私钥文件(.key)" name="ssl_certificate_key">
          <a-textarea :rows="5" v-model:value="formState.ssl_certificate_key" />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'

const emit = defineEmits(['ok'])

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
  ssl_certificate: [{ required: true, message: '请输入证书文件(.crt或.pem)' }],
  ssl_certificate_key: [{ required: true, message: '请输入私钥文件(.key)' }]
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
