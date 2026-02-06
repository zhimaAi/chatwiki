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
    :title="t('title_edit_desc')"
    @ok="handleOk"
    :confirmLoading="props.confirmLoading"
    :ok-text="t('btn_save')"
  >
    <div class="upload-ssl-form">
      <a-form ref="formRef" :model="formState" :rules="rules" autocomplete="off" layout="vertical">
        <a-form-item label="" name="content">
          <a-textarea :rows="5" v-model:value="formState.content" />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.home.components.edit-desc')

const emit = defineEmits(['ok'])

const props = defineProps({
  confirmLoading: {
    type: Boolean,
    default: false
  }
})

const formRef = ref()

const formState = reactive({
  content: ''
})

const rules = {
  // content: [{ required: true, message: '请输入描述' }]
}

const show = ref(false)

const open = (record) => {
  formState.content = record.content
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
