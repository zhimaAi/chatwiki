<template>
  <div>
    <a-modal
      v-model:open="open"
      title="设置答案"
      @ok="handleOk"
      :width="746"
      :confirmLoading="saveLoading"
    >
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item label="问题" v-bind="validateInfos.question">
            <a-textarea
              placeholder="请输入问题"
              v-model:value="formState.question"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>
          <a-form-item
            label="相似问法(一行一个，最多可添加100个相似问法)"
            v-bind="validateInfos.unknown_list"
          >
            <a-textarea
              placeholder="请输入相似问法"
              v-model:value="formState.unknown_list"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>
          <a-form-item label="答案" v-bind="validateInfos.answer">
            <a-textarea
              placeholder="请输入分段答案"
              v-model:value="formState.answer"
              :maxlength="10000"
              style="height: 100px"
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
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { isArray } from 'ant-design-vue/lib/_util/util.js'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { unknownIssueSummaryAnswer } from '@/api/robot/index.js'

const emit = defineEmits(['ok'])
const query = useRoute().query
const activeKey = ref('1')
const open = ref(false)
const formState = reactive({
  id: query.id,
  question: '',
  unknown_list: '',
  answer: '',
  images: []
})

const useForm = Form.useForm
const saveLoading = ref(false)
const show = (data) => {
  Object.assign(formState, JSON.parse(JSON.stringify(data)))
  open.value = true
}

const rules = reactive({
  question: [{ required: true, message: '请输入问题', trigger: 'change' }],
  unknown_list: [{ required: true, message: '请输入相似问法', trigger: 'change' }],
  answer: [{ required: true, message: '请输入答案', trigger: 'change' }]
})

const { validate, validateInfos } = useForm(formState, rules)

const handleOk = () => {
  console.log('---')
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log(err, 'err')
    })
}

const saveForm = () => {
  let data = {
    ...formState
  }

  let similarQuestions = data.unknown_list.trim()
  if (similarQuestions) {
    similarQuestions = similarQuestions.split('\n')
    data.unknown_list = JSON.stringify(similarQuestions)
  } else {
    data.unknown_list = '[]'
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
  saveLoading.value = true
  unknownIssueSummaryAnswer(formData)
    .then((res) => {
      message.success('保存成功')
      open.value = false
      emit('ok')
    })
    .finally(() => {
      saveLoading.value = false
    })
}
defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-box {
  margin: 24px 0;
}
.ml4 {
  margin-left: 4px;
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
</style>
