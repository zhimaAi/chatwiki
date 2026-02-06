<template>
  <div>
    <a-modal
      v-model:open="open"
      @ok="handleOk"
      wrapClassName="no-padding-modal"
      :bodyStyle="{ 'max-height': '670px', 'overflow-y': 'auto', 'padding-right': '12px' }"
      :width="!isAIGenerate ? '746px' : isGenerate ? '944px' : '472px'"
      :ok-button-props="{ style: { display: isAIGenerate ? 'none' : 'block' } }"
    >
      <!-- 自定义标题 -->
      <template #title>
        <div class="custom-modal-header">
          <a-button
            v-if="isAIGenerate"
            type="text"
            class="back-button"
            @click="isAIGenerate = false"
          >
            <template #icon><LeftOutlined style="font-size: 12px" /></template>
          </a-button>
          <span v-if="isAIGenerate">{{ t('title_ai_generate') }}</span>
          <span v-else>{{ t('title_import_library') }}</span>
        </div>
      </template>

      <div class="form-box-wrapper" v-show="!isAIGenerate">
        <div class="form-item">
          <div class="form-label required">{{ t('label_library') }}</div>
          <div class="form-content">
            <a-select
              @change="handleChangeLibrary"
              v-model:value="formState.library_id"
              style="width: 100%"
              :placeholder="t('placeholder_select')"
            >
              <a-select-option v-for="item in qaLists" :value="item.id" :key="item.id">{{
                item.library_name
              }}</a-select-option>
            </a-select>
          </div>
        </div>
        <div class="form-item">
          <div class="form-label required">{{ t('label_category') }}</div>
          <div class="form-content">
            <a-select v-model:value="group_id" style="width: 100%">
              <a-select-option v-for="item in groupLists" :value="item.id">{{
                item.group_name
              }}</a-select-option>
            </a-select>
          </div>
        </div>
        <div class="form-item">
          <div class="form-label required">{{ t('label_question') }}</div>
          <div class="form-content">
            <a-textarea
              :placeholder="t('placeholder_question')"
              v-model:value="question"
              style="height: 100px"
            ></a-textarea>
          </div>
        </div>
        <div class="form-item">
          <div class="form-label-box">
            <div class="form-label required">{{ t('label_similar_questions') }}</div>
            <div class="ai-generate" @click="handleAIGenerate">
              <svg-icon name="ai-generate" style="font-size: 14px"></svg-icon>
              <div class="ai-generate-text">{{ t('button_ai_generate') }}</div>
            </div>
          </div>
          <div class="form-content">
            <a-textarea
              :placeholder="t('placeholder_similar_questions')"
              v-model:value="similar_questions"
              style="height: 100px"
              @blur="onProcessText"
            ></a-textarea>
            <div class="form-content-tip">{{ t('tip_similar_questions') }}</div>
          </div>
        </div>
        <div class="form-item">
          <div class="form-label required">{{ t('label_answer') }}</div>
          <div class="form-content">
            <a-textarea
              :placeholder="t('placeholder_answer')"
              v-model:value="answer"
              :maxlength="10000"
              style="height: 100px"
            ></a-textarea>
            <div v-if="answer.length > 10000" class="error-tip">{{ t('error_answer_length') }}</div>
          </div>
        </div>
        <div class="form-item">
          <div class="form-label">{{ t('label_attachments') }}</div>
          <div class="form-content">
            <div class="upload-box-wrapper">
              <a-tabs v-model:activeKey="activeKey" size="small">
                <a-tab-pane key="1">
                  <template #tab>
                    <span>
                      <svg-icon name="img-icon" style="font-size: 14px; color: #2475fc"></svg-icon>
                      {{ t('tab_images') }}
                      <span v-if="images.length">({{ images.length }})</span>
                    </span>
                  </template>
                </a-tab-pane>
              </a-tabs>
              <UploadImg v-model:value="images"></UploadImg>
            </div>
          </div>
        </div>
      </div>

      <!-- AI生成表单 -->
      <div class="ai-generate-wrapper" v-show="isAIGenerate">
        <div class="form-box-wrapper ai-generate-wrapper-box">
          <div class="ai-generate-wrapper-left">
            <!-- 模型选择 -->
            <div class="form-box-wrapper-options">
              <div class="form-item flex1">
                <div class="form-label required">{{ t('label_model') }}</div>
              <div class="form-content">
                <!-- 自定义选择器 -->
                <ModelSelect
                  modelType="LLM"
                  v-model:modeName="formState.use_model"
                  v-model:modeId="formState.model_config_id"
                />
              </div>
            </div>

          <!-- 问题个数和生成按钮 -->
          <div class="form-item flex1">
            <div class="form-content dual-inputs">
              <div class="input-group">
                <div class="form-label required">{{ t('label_question_count') }}</div>
                <div class="form-label-input-box">
                  <a-input-number
                    v-model:value="formState.count"
                    :min="1"
                    :max="20"
                    style="width: 100%; padding: 2px 11px"
                  />个
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="form-box-wrapper-content">
          <div class="problem-box">
            <div class="problem">{{ t('label_question') }}</div>
            <div class="problem-content">{{ question }}</div>
          </div>
          <div class="problem-box">
            <div class="problem">{{ t('label_answer') }}</div>
            <div class="problem-content">{{ answer }}</div>
          </div>
        </div>

        <div class="form-item generate-btn-box">
          <a-button @click="handleGenerate" :loading="generating" class="generate-btn">
            <svg-icon name="ai-generate-white" style="font-size: 16px"></svg-icon>
            {{ t('button_generate') }}
          </a-button>
        </div>
          </div>
          <div class="ai-generate-wrapper-right" :class="{ 'show-result': isGenerate }">
            <!-- 生成结果区域 -->
            <div class="result-nav">
              <svg-icon name="generate-icon" style="font-size: 14px"></svg-icon>
              <span class="tip" v-if="generating">{{ t('generating') }}</span>
              <div class="result-nav-text" v-else>{{ t('generate_complete') }}</div>
            </div>

            <div class="container" v-if="!messageObj.content && generating && !isError">
              <div class="rect-container rect1">
                <div class="rect"></div>
              </div>
              <div class="rect-container rect2">
                <div class="rect"></div>
              </div>
              <div class="rect-container rect3">
                <div class="rect"></div>
              </div>
            </div>

            <div class="result-container">
              <div class="dialog-content" v-html="renderedMarkdown"></div>
            </div>
            <div class="result-footer">
              <a-button class="result-footer-btn" type="primary" @click="handleSaveSimilar"
                >{{ t('button_save') }}</a-button
              >
            </div>
          </div>
        </div>
      </div>

      <!-- 底部自定义footer处理步骤切换 -->
      <template #footer>
        <template v-if="!isAIGenerate">
          <a-button @click="open = false">{{ t('button_cancel') }}</a-button>
          <a-button type="primary" @click="handleOk">{{ t('button_confirm') }}</a-button>
        </template>
      </template>
    </a-modal>
  </div>
</template>
<script setup>
import { LeftOutlined } from '@ant-design/icons-vue'
import { ref, reactive, nextTick, watch } from 'vue'
import { message } from 'ant-design-vue'
import {
  getLibraryList,
  editParagraph,
  generateSimilarQuestions,
  getLibraryGroup
} from '@/api/library'
import { unknownIssueSummaryImport } from '@/api/robot/index.js'
import UploadImg from '@/components/upload-img/index.vue'
import { isArray } from 'ant-design-vue/lib/_util/util.js'
import ModelSelect from '@/components/model-select/model-select.vue'
import MarkdownIt from 'markdown-it'
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('views.robot.robot-config.unknown-issue.summarize.components.import-library-modal')

// 新增状态
const isError = ref(false)
const isAIGenerate = ref(false)
const showResult = ref(false)
const generating = ref(false)
const isGenerate = ref(false)
const formState = reactive({
  use_model: '',
  use_model_icon: '', // 新增图标字段
  use_model_name: '', // 新增系统名称
  is_offline: '',
  model_config_id: '',
  count: 5,
  library_id: void 0
})
const md = new MarkdownIt({
  html: true, // 启用 HTML 标签
  linkify: true, // 自动转换 URL 为链接
  typographer: true // 启用一些排版替换
})
const renderedMarkdown = ref('')
const messageObj = reactive({
  finish: '',
  content: ''
})

watch(
  () => messageObj.content,
  () => {
    renderedMarkdown.value = md.render(messageObj.content)
  }
)

const emit = defineEmits(['ok'])

const activeKey = ref('1')
const open = ref(false)
const title = ref('')
const content = ref('')
const answer = ref('')
const question = ref('')
const similar_questions = ref('')
const images = ref([])
const id = ref('')

const group_id = ref('0')
const qaLists = ref([])
const getQaLibraryList = () => {
  getLibraryList({
    type: 2
  }).then((res) => {
    group_id.value = '0'
    qaLists.value = res.data || []
  })
}

getQaLibraryList()

const handleChangeLibrary = () => {
  getGroupLists()
}
const showModal = (data) => {
  // 重置AI生成状态
  getQaLibraryList()
  isAIGenerate.value = false
  showResult.value = false
  isGenerate.value = false

  title.value =  ''
  content.value = ''
  id.value = data.id || ''
  answer.value = data.answer || ''
  images.value = data.images || []
  question.value = data.question || ''
  similar_questions.value = data.similar_questions ? data.similar_questions.join('\n') : ''
  group_id.value = data.group_id || '0'
  formState.library_id = data.library_id || ''
  if (!formState.library_id) {
    formState.library_id = qaLists.value[0]?.id
  }
  getGroupLists()
  open.value = true
}

const groupLists = ref([])
const getGroupLists = () => {
  getLibraryGroup({
    library_id: formState.library_id
  }).then((res) => {
    groupLists.value = res.data || []
  })
}

const handleOk = () => {
  // 格式相似问法内容
  onProcessText()

  if (!question.value) {
    return message.error(t('validation_question'))
  }
  if (!answer.value) {
    return message.error(t('validation_answer'))
  }
  if (answer.value.length > 10000) {
    return
  }

  let data = {
    title: title.value,
    content: content.value,
    question: question.value,
    answer: answer.value,
    images: images.value,
    library_id: formState.library_id,
  }
  data.group_id = group_id.value || 0
  let similarQuestions = similar_questions.value.trim()
  if (similarQuestions) {
    similarQuestions = similarQuestions.split('\n')
    data.similar_questions = JSON.stringify(similarQuestions)
  } else {
    data.similar_questions = '[]'
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

  editParagraph(formData).then((res) => {
    unknownIssueSummaryImport({
      id: id.value,
      to_library_id: formState.library_id
    }).then((res) => {
      message.success(t('success_import'))
      open.value = false
      emit('ok')
    })
  })
}

// 处理AI生成按钮点击
const handleAIGenerate = () => {
  if (!formState.library_id) {
    return message.error(t('validation_library'))
  }
  isAIGenerate.value = true
}

const aiResult = ref('')
// 处理生成请求
const handleGenerate = async () => {
  messageObj.content = ''
  aiResult.value = ''
  generating.value = true
  isGenerate.value = true
  try {
    let params = {
      library_id: formState.library_id,
      model_config_id: formState.model_config_id,
      use_model: formState.use_model,
      num: formState.count,
      question: question.value,
      answer: answer.value
    }
    const res = await generateSimilarQuestions(params)
    // 接口调用
    aiResult.value = res.data.map((item) => item).join('\n')
    messageObj.content = res.data
      .map((text, i) => `<div class="dialog-message">${i + 1}. ${text}</div>`)
      .join('')
    showResult.value = true
  } finally {
    generating.value = false
  }
}

// 处理保存
const handleSaveSimilar = () => {
  // 这里添加保存逻辑
  if (!similar_questions.value) {
    similar_questions.value += aiResult.value
  } else {
    similar_questions.value += `\n${aiResult.value}`
  }
  message.success(t('success_save'))
  showResult.value = false
  isAIGenerate.value = false
  isGenerate.value = false
}

function processText(input) {
  return input
    .split('\n') // 按换行符分割成数组
    .map((line) => line.trim()) // 去除每行首尾空格
    .filter((line) => line) // 过滤掉空行
    .join('\n') // 重新用换行符连接
}

const onProcessText = () => {
  similar_questions.value = processText(similar_questions.value)
}

defineExpose({ showModal })
</script>
<style lang="less" scoped>
.custom-modal-header {
  display: flex;
  align-items: center;
  gap: 4px;

  .back-button {
    width: 24px;
    height: 24px;
    display: flex;
    padding: 4px;
    justify-content: center;
    align-items: center;
    gap: 4px;
    border-radius: 6px;
    color: rgba(0, 0, 0, 0.45);
    &:hover {
      color: rgba(0, 0, 0, 0.85);
      background: none;
    }
  }

  span {
    font-weight: 500;
    font-size: 16px;
  }
}

.dual-inputs {
  display: flex;
  gap: 16px;
  align-items: flex-end;

  .input-group {
    flex: 1;
  }

  .form-label-input-box {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

.result-container {
  position: relative;
  padding-bottom: 50px;

  .dialog-content {
    margin-top: 16px;
    max-height: 400px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 16px;
    color: #262626;
    text-align: left;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
}

.ai-generate-wrapper {
  :deep(.ant-modal-body) {
    padding: 24px;
  }
}

.flex1 {
  flex: 1;
}

.form-box-wrapper {
  .form-box-wrapper-options {
    display: flex;
    gap: 16px;
    align-items: center;
    width: 408px;
  }
  .form-item {
    margin-top: 16px;
  }

  .form-label-box {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .ai-generate {
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      padding: 1px 8px;
      justify-content: center;
      align-items: center;
      gap: 4px;
      border-radius: 6px;

      &:hover {
        background: #e4e6eb;
      }

      .ai-generate-text {
        color: #6524fc;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
      }
    }
  }

  .form-label {
    color: #262626;
    font-size: 14px;
    line-height: 22px;
    padding-top: 5px;
    &.required::before {
      content: '*';
      display: inline-block;
      color: #fb363f;
      margin-right: 2px;
    }
  }
  .form-content {
    margin-top: 8px;

    .form-content-tip {
      align-self: stretch;
      color: #8c8c8c;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }
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
}
.ant-segmented-item-selected .star-item-box {
  color: #2475fc;
}
.star-item-box {
  display: flex;
  align-items: center;
  gap: 4px;
  .anticon {
    font-size: 16px;
  }
}
.error-tip {
  margin-top: 4px;
  color: #fb363f;
}

.form-box-wrapper-content {
  overflow-x: hidden;
  overflow-y: auto;
  max-height: 304px;
  margin-top: 16px;
  display: flex;
  width: 422px;
  flex-direction: column;
  align-items: flex-start;
  gap: 16px;
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
  padding-bottom: 50px;

  .problem {
    color: #262626;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
}

.generate-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: var(---ai, linear-gradient(94deg, #2475fc 0.65%, #3c01ff 53.2%, #c20cff 100%));
  color: #ffffff;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 32px;
  width: 160px;

  &:hover {
    color: white !important;
  }
}

.ai-generate-wrapper-box {
  display: flex;
  min-height: 450px;

  .ai-generate-wrapper-left {
    position: relative;
    flex: 1;
    padding-right: 12px;

    .generate-btn-box {
      display: flex;
      justify-content: center;
      position: absolute;
      left: 50%;
      bottom: 0px;
      transform: translateX(-50%);
    }
  }

  .ai-generate-wrapper-right {
    position: relative;
    flex: 1;
    border-left: 1px solid #f0f0f0;
    padding: 0 32px;
    max-height: 0;
    overflow: hidden;
    transition: all 0.3s ease;

    &.show-result {
      max-height: 500px;
    }

    .result-nav {
      width: 100%;
      display: flex;
      align-items: center;
      gap: 8px;

      .result-nav-text {
        color: #6524fc;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
      }
    }

    .result-footer {
      position: absolute;
      left: 50%;
      bottom: 0px;
      transform: translateX(-50%);

      .result-footer-btn {
        width: 160px;
      }
    }
  }
}

.container {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 12px;
}

.rect-container {
  overflow: hidden;
  height: 20px;

  &.rect1 {
    width: 206px;
  }

  &.rect2,
  &.rect3 {
    width: 408px;
  }
}

.rect {
  height: 100%;
  width: 100%;
  border-radius: 24px;
  background: linear-gradient(94deg, #00bfff 2.9%, #c1c4ff 63.43%, #fff 98.28%);
  transform-origin: left;
  animation: slide 1s infinite;

  .rect2 & {
    animation-delay: 0.1s;
  }

  .rect3 & {
    animation-delay: 0.2s;
  }
}

@keyframes slide {
  from {
    transform: scaleX(0);
  }

  to {
    transform: scaleX(1);
  }
}
</style>
