<template>
  <div>
    <a-modal v-model:open="show" :title="t('title')" @ok="handleOk" width="746px">
      <a-form layout="vertical">
        <a-form-item :label="t('label_title')">
          <a-input
            type="text"
            :placeholder="t('ph_title')"
            v-model:value="formState.title"
          ></a-input>
        </a-form-item>
        <template v-if="isExcelQa">
          <a-form-item :label="t('label_question')" v-bind="validateInfos.question">
            <a-textarea
              :placeholder="t('ph_question')"
              v-model:value="formState.question"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>
          <a-form-item :label="t('label_similar')">
            <a-textarea
              :placeholder="t('ph_similar')"
              v-model:value="formState.similar_question_list"
              style="height: 100px"
            ></a-textarea>
          </a-form-item>
          <a-form-item :label="t('label_answer')" v-bind="validateInfos.answer">
            <a-textarea
              :placeholder="t('ph_answer')"
              v-model:value="formState.answer"
              style="height: 140px"
            ></a-textarea>
          </a-form-item>
        </template>
        <a-form-item v-else :label="t('label_content')" v-bind="validateInfos.content">
          <a-textarea
            :placeholder="t('ph_content')"
            v-model:value="formState.content"
            style="height: 150px"
          ></a-textarea>
        </a-form-item>
        <a-form-item :label="t('label_attachment')">
          <div class="upload-box-wrapper">
            <a-tabs size="small">
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
          <div class="upload-box-wrapper" style="margin-top: 16px;" v-if="showMiniCard">
            <a-tabs size="small">
              <a-tab-pane key="1">
                <template #tab>
                  <span>
                    <svg-icon name="applet" style="font-size: 14px; color: #2475fc"></svg-icon>
                    {{ t('tab_mini_card') }}
                    <span v-if="formState.mini_cards.length">({{ formState.mini_cards.length }})</span>
                  </span>
                </template>
              </a-tab-pane>
            </a-tabs>
            <MiniCardTabContent v-model="formState.mini_cards" />
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
import MiniCardTabContent from '@/components/mini-card/mini-card-tab-content.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.document-segmentation.components.edit-fragment-alert')
defineProps({
  showMiniCard: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['ok'])

const useForm = Form.useForm
const show = ref(false)

const formState = reactive({
  title: '',
  content: '',
  question: '',
  answer: '',
  similar_question_list: '',
  images: [],
  mini_cards: []
})

const formRules = reactive({
  content: [
    {
      message: t('validation_content'),
      validator: async (rule, value) => {
        if (!isExcelQa.value) {
          if (!value) {
            return Promise.reject(t('validation_content'))
          }
          return Promise.resolve()
        }
        return Promise.resolve()
      }
    }
  ],
  question: [
    {
      message: t('validation_question'),
      validator: async (rule, value) => {
        if (isExcelQa.value) {
          if (!value) {
            return Promise.reject(t('validation_question'))
          }
          return Promise.resolve()
        }
        return Promise.resolve()
      }
    }
  ],
  answer: [
    {
      message: t('validation_answer'),
      validator: async (rule, value) => {
        if (isExcelQa.value) {
          if (!value) {
            return Promise.reject(t('validation_answer'))
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

const open = ({ title, content, question, answer, images, mini_card, similar_question_list }) => {
  formState.title = title
  formState.content = content
  formState.question = question
  formState.answer = answer
  formState.images = images || []
  formState.mini_cards = mini_card || []
  formState.similar_question_list = similar_question_list? similar_question_list.join('\n') : ''
  isExcelQa.value = question != ''

  show.value = true
}

const handleOk = () => {
  validate().then(() => {
    show.value = false
    const raw = toRaw(formState)
    const { mini_cards, ...rest } = raw
    emit('ok', {
      ...rest,
      mini_card: mini_cards.map(item => item.id).join(','),
      mini_cards
    })
  })
}

defineExpose({
  open
})
</script>

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
