<template>
  <a-modal
    v-model:open="open"
    @ok="handleOk"
    title="编辑分段"
    wrapClassName="no-padding-modal"
    :bodyStyle="{ 'max-height': '670px', 'overflow-y': 'auto', 'padding-right': '12px' }"
    :width="944"
  >
    <div class="form-box">
      <a-form layout="vertical">
        <!-- <a-form-item ref="title" label="分段标题" v-bind="validateInfos.title">
          <a-input :maxLength="25" v-model:value="formState.title" placeholder="请输入分段标题" />
        </a-form-item> -->
        <a-form-item ref="question" label="分段问题" v-bind="validateInfos.question">
          <a-textarea
            placeholder="请输入分段问题"
            v-model:value="formState.question"
            style="height: 100px"
          ></a-textarea>
        </a-form-item>
        <a-form-item ref="answer" label="分段答案" v-bind="validateInfos.answer">
          <a-textarea
            placeholder="请输入分段问题"
            v-model:value="formState.answer"
            style="height: 100px"
          ></a-textarea>
        </a-form-item>

        <a-form-item label="附件" v-bind="validateInfos.images">
          <div class="upload-box-wrapper">
            <a-tabs v-model:activeKey="activeKey" size="small">
              <a-tab-pane key="1">
                <template #tab>
                  <span>
                    <svg-icon name="img-icon" style="font-size: 14px; color: #2475fc"></svg-icon>
                    图片
                    <span v-if="formState.images.length">({{ formState.images.length }})</span>
                  </span>
                </template>
              </a-tab-pane>
            </a-tabs>
            <UploadImg v-model:value="formState.images"></UploadImg>
          </div>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { Form, message } from 'ant-design-vue'
import { saveFAQFileQA } from '@/api/library/index'
import { isArray } from 'ant-design-vue/lib/_util/util.js'
import { SettingOutlined } from '@ant-design/icons-vue'
import UploadImg from '@/components/upload-img/index.vue'
const emit = defineEmits(['ok'])

const open = ref(false)
const activeKey = ref('1')
const useForm = Form.useForm
const formState = reactive({
  images: [],
  question: '',
  answer: '',
  id: ''
})

const rules = reactive({
  // title: [{ required: true, message: '请输入分段标题', trigger: 'change' }],
  question: [{ required: true, message: '请输入分段问题', trigger: 'change' }],
  answer: [{ required: true, message: '请输入分段答案', trigger: 'change' }]
  // images: [{ required: true, message: '请上传图片', trigger: 'change' }]
})

const { validate, validateInfos } = useForm(formState, rules)

const handleOk = () => {
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {})
}

const saveForm = () => {
  let data = {
    ...formState
  }
  let formData = new FormData()
  for (let key in data) {
    if (isArray(data[key])) {
      data[key].forEach((v) => {
        formData.append(key, v)
      })
    } else {
      formData.append(key, data[key])
    }
  }
  saveFAQFileQA(formData).then((res) => {
    message.success('编辑成功')
    open.value = false
    emit('ok')
  })
}

const show = (data) => {
  formState.question = data.question
  formState.answer = data.answer
  formState.id = data.id
  formState.title = data.title
  formState.images = data.images || []
  open.value = true
}

defineExpose({
  show
})

onMounted(() => {})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 32px;
  margin-bottom: 24px;
}

.form-item-tip {
  color: #999;
  color: #8c8c8c;
  margin-top: 8px;
  text-align: center;
}

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

.ml4 {
  margin-left: 4px;
}
</style>
