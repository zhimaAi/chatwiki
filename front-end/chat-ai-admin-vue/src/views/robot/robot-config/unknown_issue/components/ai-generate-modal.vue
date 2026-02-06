<template>
  <div>
    <a-modal
      v-model:open="open"
      :maskClosable="false"
      :width="'472px'"
      :ok-button-props="{ style: { display: 'none' } }"
    >
      <!-- 自定义标题 -->
      <template #title>
        <div class="custom-modal-header">
          <span>{{ t('title') }}</span>
        </div>
      </template>
      <!-- AI生成表单 -->
      <div class="ai-generate-wrapper">
        <div class="form-box-wrapper ai-generate-wrapper-box">
          <div class="ai-generate-wrapper-left">
            <!-- 模型选择 -->
            <div class="form-box-wrapper-options">
              <div class="form-item flex1">
                <div class="form-label required">{{ t('model_selection_label') }}</div>
                <div class="form-content">
                  <ModelSelect
                    modelType="LLM"
                    v-model:modeName="formState.use_model"
                    v-model:modeId="formState.model_config_id"
                  />
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
                {{ t('generate_button') }}
              </a-button>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部自定义footer处理步骤切换 -->
      <template #footer>
      </template>
    </a-modal>

    <!-- 生成中 -->
    <a-modal v-model:open="resultOpen" :title="t('generating_title')" :footer="null" :width="746">
      <div class="progress-box" v-if="percent < 100">
        <a-progress type="circle" :percent="percent" :status="progressStatus" :stroke-color="strokeColor"/>
        <div class="tip">{{ t('generating_tip') }}</div>
      </div>
      <a-result v-else status="success" :title="t('complete_title')">
        <template #subTitle>
          {{ t('complete_subtitle') }}
        </template>
        <template #extra>
          <a-button
            @click="downFailData"
            type="primary"
            >{{ t('download_button') }}</a-button
          >
        </template>
      </a-result>
    </a-modal>
  </div>
</template>
<script setup>
import { ref, reactive, nextTick, computed, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import { generateSimilarQuestions } from '@/api/library'
import ModelSelect from '@/components/model-select/model-select.vue'
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('views.robot.robot-config.unknown-issue.components.ai-generate-modal')

// 新增状态
const percent = ref(10)
const resultOpen = ref(false)
const generating = ref(false)
const formState = reactive({
  use_model: '',
  use_model_icon: '', // 新增图标字段
  use_model_name: '', // 新增系统名称
  is_offline: '',
  model_config_id: '',
  count: 1000
})

const props = defineProps({
  detailsInfo: {
    type: Object
  }
})
const emit = defineEmits(['handleDownload'])
const open = ref(false)
const showModal = () => {
  // 重置AI生成状态
  open.value = true
}

const hideModal = () => {
  resultOpen.value = false
}
// 处理生成请求
const handleGenerate = async () => {
  generating.value = true
  try {
    let params = {
      library_id: props.detailsInfo?.library_id,
      model_config_id: formState.model_config_id,
      use_model: formState.use_model,
      num: formState.count
    }
    // generateSimilarQuestions 换成正式接口
    // const res = await generateSimilarQuestions(params)

    open.value = false
    resultOpen.value = true
    startGenerate()
    // 接口调用
  } finally {
    generating.value = false
  }
}

// 是否正在生成
const isGenerating = ref(false)
// 是否已完成
const isCompleted = ref(false)
// 进度条定时器
let progressTimer = null

// 动态进度条状态
const progressStatus = computed(() => {
  return percent.value >= 100 ? 'success' : 'active'
})

// 动态渐变色（示例）
const strokeColor = computed(() => ({
  '0%': '#108ee9',
  '100%': '#87d068'
}))

// 开始生成
const startGenerate = () => {
  isGenerating.value = true
  isCompleted.value = false
  percent.value = 0

  // 模拟进度更新（真实项目替换为实际任务）
  progressTimer = setInterval(() => {
    if (percent.value < 100) {
      percent.value += Math.floor(Math.random() * 10) + 1
      if (percent.value > 100) percent.value = 100
    } else {
      clearInterval(progressTimer)
      isGenerating.value = false
      isCompleted.value = true
    }
  }, 800)
}


const downFailData = () => {
  emit('handleDownload')
}


// 页面关闭拦截
const handleBeforeUnload = (e) => {
  if (isGenerating.value) {
    e.preventDefault()
    e.returnValue = t('before_unload_tip')
    return t('before_unload_tip')
  }
}

// 组件卸载前清理
onBeforeUnmount(() => {
  clearInterval(progressTimer)
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

defineExpose({
  showModal,
  hideModal
})
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
  min-height: 142px;
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
}

.progress-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  margin: 100px 0;
  color: #8c8c8c;
}
</style>
