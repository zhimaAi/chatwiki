<style lang="less" scoped>
.upload-box-wrapper {
  background: #f2f4f7;
  border-radius: 6px;
  &::v-deep(.ant-tabs-nav::before) {
    border-color: #f2f4f7;
  }
  &::v-deep(.ant-tabs-nav) {
    margin: 0;
    margin-left: 16px;
  }
}
</style>

<template>
  <div>
    <a-modal v-model:open="show" title="编辑分段" @ok="handleOk" width="746px">
      <a-form layout="vertical">
        <a-form-item label="分段标题：">
          <a-input
            type="text"
            placeholder="请输入分段标题"
            v-model:value="formState.title"
          ></a-input>
        </a-form-item>
        <template v-if="isExcelQa">
          <a-form-item label="分段问题：" v-bind="validateInfos.question">
            <a-textarea
              placeholder="请输入分段问题"
              v-model:value="formState.question"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>
          <a-form-item label="分段答案：" v-bind="validateInfos.answer">
            <a-textarea
              placeholder="请输入分段答案"
              v-model:value="formState.answer"
              style="height: 140px"
            ></a-textarea>
          </a-form-item>
        </template>
        <a-form-item v-else label="分段内容：" v-bind="validateInfos.content">
          <a-textarea
            placeholder="请输入分段内容"
            v-model:value="formState.content"
            style="height: 150px"
          ></a-textarea>
        </a-form-item>
        <a-form-item label="附件">
          <div class="upload-box-wrapper">
            <a-tabs v-model:activeKey="activeKey" size="small">
              <a-tab-pane key="1">
                <template #tab>
                  <span>
                    <svg-icon name="img-icon" style="font-size: 14px; color: #2475fc"></svg-icon>
                    图片
                    <span v-if="formState.images.length">({{formState.images.length}})</span>
                  </span>
                </template>
              </a-tab-pane>
            </a-tabs>
            <UploadImg v-model:value="formState.images"></UploadImg>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { Form } from 'ant-design-vue'
import UploadImg from '@/components/upload-img/index.vue'
const emit = defineEmits(['ok'])

const activeKey = ref('1')
const useForm = Form.useForm
const show = ref(false)

const formState = reactive({
  title: '',
  content: '',
  question: '',
  answer: '',
  images: []
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

const open = ({ title, content, question, answer, images }) => {
  formState.title = title
  formState.content = content
  formState.question = question
  formState.answer = answer
  formState.images = images || []
  isExcelQa.value = question != ''

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
