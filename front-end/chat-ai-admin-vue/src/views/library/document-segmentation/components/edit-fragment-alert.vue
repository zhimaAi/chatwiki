<template>
  <div>
    <a-modal v-model:open="show" title="编辑分段" @ok="handleOk" width="746px">
      <a-form layout="vertical">
        <a-form-item label="分段标题：">
          <a-input type="text" placeholder="请输入分段标题" v-model:value="formState.title"></a-input>
        </a-form-item>
        <template v-if="isExcelQa">
          <a-form-item label="分段问题：" v-bind="validateInfos.question">
            <a-textarea placeholder="请输入分段问题" v-model:value="formState.question" style="height: 100px"></a-textarea>
          </a-form-item>
          <a-form-item label="分段答案：" v-bind="validateInfos.answer">
            <a-textarea placeholder="请输入分段答案" v-model:value="formState.answer" style="height: 250px"></a-textarea>
          </a-form-item>
        </template>
        <a-form-item v-else label="分段内容：" v-bind="validateInfos.content">
          <a-textarea placeholder="请输入分段内容" v-model:value="formState.content" style="height: 450px"></a-textarea>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { Form } from 'ant-design-vue'
const emit = defineEmits(['ok'])

const useForm = Form.useForm
const show = ref(false)

const formState = reactive({
  title: '',
  content: '',
  question: '',
  answer: '',
})

const formRules = reactive({
  content: [
    {
      message: '请输入分段内容',
      validator: async (rule, value) => {
        if (!isExcelQa.value) {
          if (!value) {
            return Promise.reject('请输入分段内容')
          }
          return Promise.resolve()
        }
        return Promise.resolve()
      }
    }
  ],
  question: [
    {
      message: '请输入分段问题',
      validator: async (rule, value) => {
        if (isExcelQa.value) {
          if (!value) {
            return Promise.reject('请输入分段问题')
          }
          return Promise.resolve()
        }
        return Promise.resolve()
      }
    }
  ],
  answer: [
    {
      message: '请输入分段答案',
      validator: async (rule, value) => {
        if (isExcelQa.value) {
          if (!value) {
            return Promise.reject('请输入分段答案')
          }
          return Promise.resolve()
        }
        return Promise.resolve()
      }
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)
const isExcelQa = ref(false)

const open = ({ title, content, question, answer }) => {
  formState.title = title
  formState.content = content
  formState.question = question
  formState.answer = answer
  isExcelQa.value = question != '';
  show.value = true
}

const handleOk = () => {
  validate().then(() => {
    show.value = false
    emit('ok', toRaw(formState))
  })
}

defineExpose({
  open
})
</script>
