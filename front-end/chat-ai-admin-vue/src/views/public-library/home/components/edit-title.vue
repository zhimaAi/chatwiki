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
    :title="t('edit_title')"
    @ok="handleOk"
    :confirmLoading="props.confirmLoading"
    :ok-text="t('save')"
  >
    <div class="upload-ssl-form">
      <a-form ref="formRef" :model="formState" :rules="rules" autocomplete="off" layout="vertical">
        <a-form-item label="" name="title">
          <a-textarea :rows="5" v-model:value="formState.title" />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.home.components.edit-title')

const emit = defineEmits(['ok'])

const props = defineProps({
  confirmLoading: {
    type: Boolean,
    default: false
  }
})

const formRef = ref()

const formState = reactive({
  title: ''
})

const rules = {
  title: [{ required: true, message: t('please_enter_title') }]
}

const show = ref(false)

const open = (record) => {
  formState.title = record.title
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
