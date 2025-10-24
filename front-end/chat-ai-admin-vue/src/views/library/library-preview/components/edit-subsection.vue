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
          <span v-if="isAIGenerate">AI生成相似问法</span>
          <span v-else>{{ id ? '编辑分段' : '新增分段' }}</span>
        </div>
      </template>

      <div class="form-box-wrapper" v-show="!isAIGenerate">
        <div class="form-item" v-if="!isQaDocment">
          <div class="form-label">分段标题：</div>
          <div class="form-content">
            <a-input :maxLength="25" v-model:value="title" placeholder="请输入分段标题" />
          </div>
        </div>
        <div class="form-item" v-if="!isQaDocment">
          <div class="form-label">分段分类标记</div>
          <div class="form-content">
            <a-segmented v-model:value="category_id" :options="startLists">
              <template #label="{ payload }">
                <div class="star-item-box">
                  <StarFilled v-if="payload.id > 0" :style="{ color: payload.color }" />
                  <StarOutlined v-else />
                  <div>{{ payload.name || '-' }}</div>
                  <span v-if="payload.data_count > 0">({{ payload.data_count }})</span>
                </div>
              </template>
            </a-segmented>
          </div>
        </div>
        <template v-if="isQaDocment">
          <div class="form-item">
            <div class="form-label required">所属分组：</div>
            <div class="form-content">
              <a-select v-model:value="group_id" style="width: 100%">
                <a-select-option v-for="item in groupLists" :value="item.id">{{ item.group_name }}</a-select-option>
              </a-select>
            </div>
          </div>
          <div class="form-item">
            <div class="form-label required">分段问题：</div>
            <div class="form-content">
              <a-textarea
                placeholder="请输入分段问题"
                v-model:value="question"
                style="height: 100px"
              ></a-textarea>
            </div>
          </div>
          <div class="form-item">
            <div class="form-label required">分段答案：</div>
            <div class="form-content">
              <a-textarea
                placeholder="请输入分段答案"
                v-model:value="answer"
                style="height: 100px"
              ></a-textarea>
              <div v-if="answer.length > 10000" class="error-tip">分段答案最多支持10000个字符</div>
            </div>
          </div>
        </template>
        <div class="form-item" v-else>
          <div class="form-label required">分段内容：</div>
          <div class="form-content">
            <a-textarea
              style="height: 150px"
              v-model:value="content"
              placeholder="请输入分段内容"
            />
            <div v-if="content.length > 10000" class="error-tip">分段内容最多支持10000个字符</div>
          </div>
        </div>
        <div class="form-item">
          <!-- <div class="form-label">附件</div> -->
          <div class="form-content">
            <div class="upload-box-wrapper">
              <a-tabs v-model:activeKey="activeKey" size="small">
                <a-tab-pane key="1">
                  <template #tab>
                    <span>
                      <svg-icon name="img-icon" style="font-size: 14px; color: #2475fc"></svg-icon>
                      答案附图
                      <span v-if="images.length">({{ images.length }})</span>
                    </span>
                  </template>
                </a-tab-pane>
              </a-tabs>
              <UploadImg v-model:value="images"></UploadImg>
            </div>
          </div>
        </div>

        <div class="form-item" v-if="isQaDocment">
          <div class="form-label-box">
            <div class="form-label">相似问法</div>
            <div class="ai-generate" @click="handleAIGenerate">
              <svg-icon name="ai-generate" style="font-size: 14px"></svg-icon>
              <div class="ai-generate-text">AI自动生成</div>
            </div>
          </div>
          <div class="form-content">
            <a-textarea
              placeholder="请输入相似问法"
              v-model:value="similar_questions"
              style="height: 100px"
              @blur="onProcessText"
            ></a-textarea>
            <div class="form-content-tip">一行一个，最多可添加100个相似问法</div>
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
                <div class="form-label required">相似问法模型：</div>
                <div class="form-content">
                  <!-- 自定义选择器 -->
                  <CustomSelector
                    v-model="formState.use_model"
                    label-key="use_model_name"
                    value-key="value"
                    :modelType="'LLM'"
                    :model-config-id="formState.model_config_id"
                    @change="handleModelChange"
                    @loaded="onVectorModelLoaded"
                  />
                </div>
              </div>

              <!-- 问题个数和生成按钮 -->
              <div class="form-item flex1">
                <div class="form-content dual-inputs">
                  <div class="input-group">
                    <div class="form-label required">问题个数：</div>
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
                <div class="problem">问题</div>
                <div class="problem-content">{{ question }}</div>
              </div>
              <div class="problem-box">
                <div class="problem">答案</div>
                <div class="problem-content">{{ answer }}</div>
              </div>
            </div>

            <div class="form-item generate-btn-box">
              <a-button @click="handleGenerate" :loading="generating" class="generate-btn">
                <svg-icon name="ai-generate-white" style="font-size: 16px"></svg-icon>
                开始生成
              </a-button>
            </div>
          </div>
          <div class="ai-generate-wrapper-right" :class="{ 'show-result': isGenerate }">
            <!-- 生成结果区域 -->
            <div class="result-nav">
              <svg-icon name="generate-icon" style="font-size: 14px"></svg-icon>
              <span class="tip" v-if="generating">相似问法生成中...</span>
              <div class="result-nav-text" v-else>相似问法生成完毕</div>
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
                >保存相似问法</a-button
              >
            </div>
          </div>
        </div>
      </div>

      <!-- 底部自定义footer处理步骤切换 -->
      <template #footer>
        <template v-if="!isAIGenerate">
          <a-button @click="open = false">取消</a-button>
          <a-button type="primary" @click="handleOk">确定</a-button>
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
  editParagraph,
  getCategoryList,
  saveCategoryParagraph,
  generateSimilarQuestions,
  getLibraryGroup
} from '@/api/library'
import { StarOutlined, StarFilled } from '@ant-design/icons-vue'
import { useRoute } from 'vue-router'
import UploadImg from '@/components/upload-img/index.vue'
import { isArray } from 'ant-design-vue/lib/_util/util.js'
import colorLists from '@/utils/starColors.js'
import CustomSelector from '@/components/custom-selector/index.vue'
import MarkdownIt from 'markdown-it'

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
  count: 5
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

const route = useRoute()
const query = route.query
const pathName = route.name
const props = defineProps({
  detailsInfo: {
    type: Object
  }
})
const emit = defineEmits(['handleEdit', 'handleStatrList'])
const oldModelDefineList = ['azure']
const currentModelDefine = ref('')
const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']

const activeKey = ref('1')
const open = ref(false)
const title = ref('')
const content = ref('')
const answer = ref('')
const question = ref('')
const similar_questions = ref('')
const images = ref([])
const category_id = ref(0)
const id = ref('')

const group_id = ref('0')

const isQaDocment = ref(false)
const showModal = (data) => {
  // 重置AI生成状态
  isAIGenerate.value = false
  showResult.value = false
  isGenerate.value = false

  title.value = data.title || ''
  content.value = data.content || ''
  id.value = data.id || ''
  answer.value = data.answer || ''
  images.value = data.images || []
  question.value = data.question || ''
  similar_questions.value = data.similar_questions ? data.similar_questions.join('\n') : ''
  isQaDocment.value = props.detailsInfo.is_qa_doc == '1'
  category_id.value = data.category_id || 0
  group_id.value = data.group_id || '0'
  getCategoryLists()
  if(isQaDocment) {
    getGroupLists()
  }
  open.value = true
}

const groupLists = ref([])
const getGroupLists = () => {
  getLibraryGroup({
    library_id: props.detailsInfo.library_id
  }).then((res) => {
    groupLists.value = res.data || []
  })
}


const startLists = ref([])
const getCategoryLists = () => {
  getCategoryList({ file_id: query.id }).then((res) => {
    let list = res.data || []
    list = list.map((item) => {
      return {
        value: item.id,
        payload: {
          ...item,
          color: colorLists[item.type]
        }
      }
    })
    if (pathName == 'categaryManages') {
      startLists.value = [...list]
    } else {
      startLists.value = [
        {
          value: 0,
          payload: {
            id: 0,
            name: '未标记'
          }
        },
        ...list
      ]
    }
    emit('handleStatrList', res.data || [])
  })
}

const handleOk = () => {
  // 格式相似问法内容
  onProcessText()

  if (!content.value && !isQaDocment.value) {
    return message.error('请输入分段内容')
  }
  if (!question.value && isQaDocment.value) {
    return message.error('请输入分段问题')
  }
  if (!answer.value && isQaDocment.value) {
    return message.error('请输入分段答案')
  }
  if (isQaDocment.value && answer.value.length > 10000) {
    return
  }
  if (!isQaDocment.value && content.value.length > 10000) {
    return
  }
  if (pathName == 'categaryManages' && category_id.value <= 0) {
    return message.error('请选择分类标记')
  }
  let data = {
    title: title.value,
    content: content.value,
    question: question.value,
    answer: answer.value,
    images: images.value,
    category_id: category_id.value
  }
  if(isQaDocment.value){
    data.group_id = group_id.value || 0
  }
  let similarQuestions = similar_questions.value.trim()
  if (similarQuestions) {
    similarQuestions = similarQuestions.split('\n')
    data.similar_questions = JSON.stringify(similarQuestions)
  } else {
    data.similar_questions = '[]'
  }
  if (id.value) {
    data.id = id.value
  } else {
    data.file_id = route.query.id
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
  if(pathName == 'knowledgeDocument'){
    formData.delete('file_id')
    formData.append('library_id', route.query.id)
  }
  if(pathName == 'libraryConfig'){
    formData.delete('file_id')
    formData.append('library_id', props.detailsInfo.library_id)
  }
  if (pathName == 'categaryManages') {
    formData.delete('file_id')
    formData.append('library_id', route.query.id)
    saveCategoryParagraph(formData).then((res) => {
      message.success(data.id ? '修改成功' : '添加成功')
      open.value = false
      emit('handleEdit', {
        ...data
      })
      getCategoryLists()
    })
  } else {
    editParagraph(formData).then((res) => {
      message.success(data.id ? '修改成功' : '添加成功')
      open.value = false
      emit('handleEdit', {
        ...data
      })
      getCategoryLists()
    })
  }
}

// 处理AI生成按钮点击
const handleAIGenerate = () => {
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
      library_id: props.detailsInfo.library_id,
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
  message.success('保存成功')
  showResult.value = false
  isAIGenerate.value = false
  isGenerate.value = false
}

// 处理选择事件
const handleModelChange = (item) => {
  formState.use_model =
    modelDefine.includes(item.rawData.model_define) && item.rawData.deployment_name
      ? item.rawData.deployment_name
      : item.rawData.name
  formState.use_model_icon = item.icon
  formState.use_model_name = item.use_model_name
  formState.model_config_id = item.rawData.id
  currentModelDefine.value = item.rawData.model_define
}

const vectorModelList = ref([])
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  nextTick(() => {
    if (!formState.ai_chunk_model || !Number(formState.ai_chunk_model_config_id)) {
      setDefaultModel()
    }
  })
}

const setDefaultModel = () => {
  if (vectorModelList.value.length > 0) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    // 云版默认选中qwen-max
    for (let item of vectorModelList.value) {
      if (item.model_define === 'tongyi') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'qwen-max') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig) {
      modelConfig = vectorModelList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig && model) {
      formState.use_model = model.name
      formState.model_config_id = model.model_config_id
      formState.ai_chunk_model = model.name
      formState.ai_chunk_model_config_id = model.model_config_id
    }
  }
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
