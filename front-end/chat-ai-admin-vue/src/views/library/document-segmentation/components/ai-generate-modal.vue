<template>
  <div>
    <a-modal
      v-model:open="open"
      @ok="handleOk"
      :maskClosable="false"
      :width="isGenerate ? '944px' : '472px'"
      :ok-button-props="{ style: { display: 'none' } }"
    >
      <!-- 自定义标题 -->
      <template #title>
        <div class="custom-modal-header">
          <span>AI生成提示词</span>
        </div>
      </template>
      <!-- AI生成表单 -->
      <div class="ai-generate-wrapper">
        <div class="form-box-wrapper ai-generate-wrapper-box" :class="{ 'show-generate': isGenerate }">
          <div class="ai-generate-wrapper-left">
            <a-alert class="alert-box" message="根据导入的文档的内容生成AI分段的提示词，支持设置导入文档的字数。" type="info" />
            <!-- 模型选择 -->
            <div class="form-box-wrapper-options">
              <div class="form-item flex1">
                <div class="form-label required">选择模型：</div>
                <div class="form-content">
                  <!-- 自定义选择器 -->
                  <ModelSelect
                    modelType="LLM"
                    v-model:modeName="formState.use_model"
                    v-model:modeId="formState.model_config_id"
                  />
                </div>
              </div>

              <!-- 根据文档内容前和生成按钮 -->
              <div class="form-item flex1">
                <div class="form-content dual-inputs">
                  <div class="input-group">
                    <div class="form-label required">根据文档内容前：</div>
                    <div class="form-label-input-box">
                      <a-input-number 
                        v-model:value="formState.count" 
                        :min="1" 
                        :max="10000" 
                        style="width: 137px;"
                      />个字，自动生成提示词
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div class="form-item generate-btn-box">
              <a-button 
                @click="handleGenerate"
                :loading="generating"
                class="generate-btn"
              >
                <svg-icon name="ai-generate-white" style="font-size: 16px;"></svg-icon>
                开始生成
              </a-button>
            </div>
          </div>
          <div class="ai-generate-wrapper-right" :class="{ 'show-result': isGenerate }">
            <!-- 生成结果区域 -->
            <div class="result-nav">
              <svg-icon name="generate-icon" style="font-size: 14px;"></svg-icon>
              <span class="tip" v-if="generating">提示词生成中...</span>
              <div class="result-nav-text" v-else>提示词生成完毕</div>
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
              <a-button class="result-footer-btn" :disabled="aiResult == ''" type="primary" @click="handleOk">使用该提示词</a-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部自定义footer处理步骤切换 -->
      <template #footer>
      </template>
    </a-modal>
  </div>
</template>
<script setup>
import { ref, reactive, nextTick, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import { generateAiPrompt } from '@/api/library'
import ModelSelect from '@/components/model-select/model-select.vue'
import MarkdownIt from 'markdown-it';
import { useLibraryStore } from '@/stores/modules/library'

const libraryStore = useLibraryStore()
const props = defineProps({
  detailsInfo: {
    type: Object
  }
})
// 新增状态
const isError = ref(false)
const showResult = ref(false)
const generating = ref(false)
const isGenerate = ref(false)
const formState = reactive({
  use_model: '',
  use_model_icon: '', // 新增图标字段
  use_model_name: '', // 新增系统名称
  is_offline: '',
  model_config_id: '',
  count: 1000
})
const md = new MarkdownIt({
  html: true,        // 启用 HTML 标签
  linkify: true,     // 自动转换 URL 为链接
  typographer: true // 启用一些排版替换
});
const renderedMarkdown = ref('');
const messageObj = reactive({
  content: '',
})
const documentFragmentList = computed(() => libraryStore.initDocumentFragmentList)
watch(() => messageObj.content, () => {
  renderedMarkdown.value = md.render(messageObj.content);
})

const emit = defineEmits(['handleEdit'])
const currentModelDefine = ref('')
const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent']
const open = ref(false)
const showModal = () => {
  // 重置AI生成状态
  showResult.value = false
  isGenerate.value = false
  open.value = true
}

const handleOk = () => {
  message.success('操作成功')
  showResult.value = false
  isGenerate.value = false
  open.value = false
  emit('handleEdit', aiResult.value)
}

const aiResult = ref('')

// JS实现数组内容拼接并截取
function concatContentsByLength(arr, maxLength) {
  let result = '';
  let currentLength = 0;
  for (const item of arr) {
      if (currentLength >= maxLength) break;
      const content = item.content.trim() || '';
      const remaining = maxLength - currentLength;
      if (content.length <= remaining) {
          result += content;
          currentLength += content.length;
      } else {
          result += content.slice(0, remaining);
          currentLength = maxLength;
          break;
      }
  }
  return result;
}
// 处理生成请求
const handleGenerate = async () => {
  messageObj.content = ''
  aiResult.value = ''
  generating.value = true
  isGenerate.value = true
  try {
    const ai_prompt_question = await concatContentsByLength(documentFragmentList.value, formState.count)
    let params = {
      ai_prompt_model_config_id: formState.model_config_id,
      ai_prompt_model: formState.use_model,
      ai_prompt_question: ai_prompt_question
    }
    const res = await generateAiPrompt(params)
    // 接口调用
    aiResult.value = res.data
    messageObj.content = `<div class="dialog-message">${ res.data }</div>`
    showResult.value = true
  } finally {
    generating.value = false
  }
}

defineExpose({ showModal })
</script>
<style lang="less" scoped>
.custom-modal-header {
  display: flex;
  align-items: center;
  gap: 4px;

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
    flex-direction: column;
    gap: 16px;
    align-items: center;
    width: 408px;
    margin-top: 16px;
  }
  .form-item {
    width: 100%;
  }

  .form-label {
    color: #262626;
    font-size: 14px;
    line-height: 22px;
    &.required::before {
      content: '*';
      display: inline-block;
      color: #fb363f;
      margin-right: 2px;
    }
  }
  .form-content {
    margin-top: 4px;

    .form-content-tip {
      align-self: stretch;
      color: #8c8c8c;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }
  }
}

.generate-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: var(---ai, linear-gradient(94deg, #2475FC 0.65%, #3C01FF 53.2%, #C20CFF 100%));
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
  min-height: 274px;
  // min-height: 534px;

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
    // border-left: 1px solid #f0f0f0;
    padding: 0 32px;
    max-height: 0;
    overflow: hidden;
    transition: all 0.3s ease;

    &.show-result {
      max-height: 534px;
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

.show-generate {
  min-height: 534px;
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
  background: linear-gradient(94deg, #00BFFF 2.9%, #C1C4FF 63.43%, #FFF 98.28%);
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
