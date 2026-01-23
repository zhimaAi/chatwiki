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
    :title="t('edit_question_guide')"
    @ok="handleOk"
    :confirmLoading="props.confirmLoading"
    :ok-text="t('save')"
  >
    <div class="upload-ssl-form">
      <a-form ref="formRef" :model="formState" :rules="rules" autocomplete="off" layout="vertical">
        <a-form-item label="" name="question">
          <a-textarea :rows="5" v-model:value="formState.question" />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.home.components.edit-question-guide')

const emit = defineEmits(['ok'])

const props = defineProps({
  confirmLoading: {
    type: Boolean,
    default: false
  }
})

const formRef = ref()

const formState = reactive({
  question: ''
})

const rules = {
  question: [{ required: true, message: t('please_input_question_guide') }]
}

const show = ref(false)

const open = (record) => {
  formState.question = record.question
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
