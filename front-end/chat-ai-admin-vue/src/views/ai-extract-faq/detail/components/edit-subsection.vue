<template>
  <a-modal
    v-model:open="open"
    @ok="handleOk"
    :title="t('title_edit_segment')"
    wrapClassName="no-padding-modal"
    :bodyStyle="{ 'max-height': '670px', 'overflow-y': 'auto', 'padding-right': '12px' }"
    :width="944"
  >
    <div class="form-box">
      <a-form layout="vertical">
        <!-- <a-form-item ref="title" :label="t('label_segment_title')" v-bind="validateInfos.title">
          <a-input :maxLength="25" v-model:value="formState.title" :placeholder="t('ph_input_segment_title')" />
        </a-form-item> -->
        <a-form-item ref="question" :label="t('label_segment_question')" v-bind="validateInfos.question">
          <a-textarea
            :placeholder="t('ph_input_segment_question')"
            v-model:value="formState.question"
            style="height: 100px"
          ></a-textarea>
        </a-form-item>
        <a-form-item ref="answer" :label="t('label_segment_answer')" v-bind="validateInfos.answer">
          <a-textarea
            :placeholder="t('ph_input_segment_answer')"
            v-model:value="formState.answer"
            style="height: 100px"
          ></a-textarea>
        </a-form-item>

        <a-form-item :label="t('label_attachments')" v-bind="validateInfos.images">
          <div class="upload-box-wrapper">
            <a-tabs v-model:activeKey="activeKey" size="small">
              <a-tab-pane key="1">
                <template #tab>
                  <span>
                    <svg-icon name="img-icon" style="font-size: 14px; color: #2475fc"></svg-icon>
                    {{ t('label_image') }}
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
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.ai-extract-faq.detail.edit-subsection')

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
  // title: [{ required: true, message: 'msg_input_segment_title', trigger: 'change' }],
  question: [{ required: true, message: t('msg_input_segment_question'), trigger: 'change' }],
  answer: [{ required: true, message: t('msg_input_segment_answer'), trigger: 'change' }]
  // images: [{ required: true, message: 'msg_upload_image', trigger: 'change' }]
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
    message.success(t('msg_edit_success'))
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
