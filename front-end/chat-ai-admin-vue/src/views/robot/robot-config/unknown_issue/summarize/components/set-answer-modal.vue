<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="t('title')"
      @ok="handleOk"
      :width="746"
      :confirmLoading="saveLoading"
    >
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item :label="t('label_question')" v-bind="validateInfos.question">
            <a-textarea
              :placeholder="t('placeholder_question')"
              v-model:value="formState.question"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>
          <a-form-item
            :label="t('label_similar_questions')"
            v-bind="validateInfos.unknown_list"
          >
            <a-textarea
              :placeholder="t('placeholder_similar_questions')"
              v-model:value="formState.unknown_list"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>
          <a-form-item :label="t('label_answer')" v-bind="validateInfos.answer">
            <a-textarea
              :placeholder="t('placeholder_answer')"
              v-model:value="formState.answer"
              :maxlength="10000"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>

          <a-form-item :label="t('label_attachments')">
            <div class="upload-box-wrapper">
              <a-tabs v-model:activeKey="activeKey" size="small">
                <a-tab-pane key="1">
                  <template #tab>
                    <span>
                      <svg-icon name="img-icon" style="font-size: 14px; color: #2475fc"></svg-icon>
                      {{ t('tab_images') }}
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
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('views.robot.robot-config.unknown-issue.summarize.components.set-answer-modal')

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
  question: [{ required: true, message: t('validation_question'), trigger: 'change' }],
  unknown_list: [{ required: true, message: t('validation_similar_questions'), trigger: 'change' }],
  answer: [{ required: true, message: t('validation_answer'), trigger: 'change' }]
})

const { validate, validateInfos } = useForm(formState, rules)

const handleOk = () => {
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
      message.success(t('save_success'))
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
