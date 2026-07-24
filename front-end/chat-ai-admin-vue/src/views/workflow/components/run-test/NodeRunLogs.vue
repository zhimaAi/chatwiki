<template>
  <div class="node-run-logs">
    <div class="logs-state-box" v-if="showLoadingState">
      <a-spin :tip="t('msg_generating_test_result')" />
    </div>
    <div class="result-panel-box" v-else :class="{ 'has-result': hasRunResult }">
      <div class="result-list-box customize-scroll-style" v-if="hasRunResult">
        <div
          class="list-item-block"
          :class="{ active: currentNodeKey == getNodeLogKey(item) }"
          v-for="(item, index) in resultList"
          @click="handleChangeNodeKey(item)"
          :key="getNodeLogKey(item) || index"
        >
          <div class="status-block">
            <svg-icon name="status-success" class="success-icon" v-if="item.is_success" style="color: #fff"></svg-icon>
            <CloseCircleFilled v-else style="color: #d81e06" />
          </div>
          <div class="icon-name-box">
            <img :src="item.node_icon" alt="" />
            <div class="node-name">{{ item.node_name }}</div>
          </div>
          <div class="time-tag" v-if="item.is_success">{{ item.use_time }}ms</div>
          <!-- <div class="right-active-icon"><RightCircleOutlined /></div> -->
        </div>
      </div>
      <div class="result-state-box" v-else>
        <a-empty :description="hasRunTested ? t('msg_no_test_result') : t('msg_click_run_test_tip')" />
      </div>
    </div>
    <div class="preview-box customize-scroll-style" v-if="hasRunResult">
      <template v-if="currentItem">
        <!-- <div class="preview-title">
          <div class="title-text">{{ t('title_log_details') }}</div>
          <div class="icon-name-box">
            <img :src="currentItem.node_icon" alt="" />
            <div class="node-name">{{ currentItem.node_name }}</div>
          </div>
          <div class="time-tag" v-if="currentItem.is_success">{{ currentItem.use_time }}ms</div>
        </div> -->
        <div class="preview-content-block" v-if="currentImageList.length > 0">
          <div class="title-block">{{ t('title_generated_image_log') }}</div>
          <div class="preview-img-box">
            <ImageLogs :currentImageList="currentImageList" />
          </div>
        </div>
        <div class="preview-content-block">
          <div class="title-block">{{ t('label_input') }}<CopyOutlined @click="handleCopy('input')" /></div>
          <div class="preview-code-box">
            <template v-if="isHugeData(currentItem.input)">
              <div class="large-data-tip">{{ t('msg_data_too_large') }}</div>
              <textarea class="large-data-textarea" readonly :value="getJsonText(currentItem.input)" />
            </template>
            <vue-json-pretty v-else :data="currentItem.input" v-bind="getJsonViewerProps(currentItem.input)" />
          </div>
        </div>
        <div class="preview-content-block">
          <div class="title-block">{{ t('label_output') }}<CopyOutlined @click="handleCopy('node_output')" /></div>
          <div class="preview-code-box">
            <template v-if="isHugeData(currentItem.node_output)">
              <div class="large-data-tip">{{ t('msg_data_too_large') }}</div>
              <textarea class="large-data-textarea" readonly :value="getJsonText(currentItem.node_output)" />
            </template>
            <vue-json-pretty v-else :data="currentItem.node_output" v-bind="getJsonViewerProps(currentItem.node_output)" />
          </div>
        </div>
        <div class="preview-content-block">
          <div class="title-block">{{ t('label_run_log') }}<CopyOutlined @click="handleCopy('output')" /></div>
          <div class="preview-code-box">
            <template v-if="isHugeData(currentItem.output)">
              <div class="large-data-tip">{{ t('msg_data_too_large') }}</div>
              <textarea class="large-data-textarea" readonly :value="getJsonText(currentItem.output)" />
            </template>
            <vue-json-pretty v-else :data="currentItem.output" v-bind="getJsonViewerProps(currentItem.output)" />
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  CloseCircleFilled,
  CopyOutlined
} from '@ant-design/icons-vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { message } from 'ant-design-vue'
import { copyText } from '@/utils/index'
import ImageLogs from '@/views/workflow/components/image-logs/index.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.run-test.index')

const props = defineProps({
  resultList: {
    type: Array,
    default: () => []
  },
  currentNodeKey: {
    type: String,
    default: ''
  },
  loading: {
    type: Boolean,
    default: false
  },
  hasRunTested: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:currentNodeKey'])

const hasRunResult = computed(() => props.resultList.length > 0)

const showLoadingState = computed(() => props.loading && !hasRunResult.value)

const currentItem = computed(() => {
  if (!props.currentNodeKey) {
    return null
  }
  return props.resultList.find((item) => getNodeLogKey(item) == props.currentNodeKey)
})

const currentImageList = computed(() => {
  let list = []
  if (currentItem.value && currentItem.value.node_type == 33) {
    let output = currentItem.value.output
    for (let key in output) {
      if (key.includes('picture_url_')) {
        list.push(output[key])
      }
    }
  }
  return list
})

const handleChangeNodeKey = (item) => {
  emit('update:currentNodeKey', getNodeLogKey(item))
}

const getNodeLogKey = (item) => item?.log_key || item?.node_key || ''

const handleCopy = (key) => {
  copyText(JSON.stringify(currentItem.value[key]))
  message.success(t('msg_copy_success'))
}

const JSON_COLLAPSE_THRESHOLD = 5 * 1024
const JSON_VIRTUAL_THRESHOLD = 100 * 1024
const JSON_HUGE_THRESHOLD = 500 * 1024

function getJsonSize(data) {
  if (data == null) return 0
  try { return JSON.stringify(data).length } catch { return 0 }
}

function isHugeData(data) {
  return getJsonSize(data) > JSON_HUGE_THRESHOLD
}

function getJsonViewerProps(data) {
  const size = getJsonSize(data)
  if (size > JSON_VIRTUAL_THRESHOLD) {
    return { deep: 1, virtual: true, height: 320, showLength: true }
  }
  if (size > JSON_COLLAPSE_THRESHOLD) {
    return { deep: 2, showLength: true }
  }
  return { deep: 3, showLength: true }
}

function getJsonText(data) {
  if (data == null) return ''
  try { return JSON.stringify(data, null, 2) } catch { return String(data) }
}
</script>

<style lang="less" scoped>
.node-run-logs {
  display: flex;
  flex: 1;
  min-width: 0;
  height: 100%;
}

.result-panel-box {
  width: var(--chat-test-flow-width, 360px);
  min-width: var(--chat-test-flow-width, 360px);
  box-sizing: var(--node-run-logs-panel-box-sizing, content-box);
  background: #fff;
  padding: 0;
  &.has-result {
    border-right: 1px solid #f0f0f0;
  }
}

.result-state-box {
  height: 100%;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: #8c8c8c;
  background: #fff;
}

.logs-state-box {
  flex: 1;
  min-width: 0;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: #8c8c8c;
  background: #fff;
}

.result-list-box {
  height: 100%;
  width: 100%;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  padding: 16px;
  overflow-y: auto;
  background: #fff;
  .list-item-block {
    display: flex;
    align-items: center;
    overflow: hidden;
    gap: 10px;
    height: 32px;
    padding: 0 16px;
    color: #333;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    margin-bottom: 4px;
    .right-active-icon {
      margin-left: auto;
      color: #2475fc;
      opacity: 0;
      font-size: 14px;
    }
    &:hover {
      background: #f2f4f7;
      .right-active-icon {
        opacity: 1;
      }
    }
    &.active {
      color: #2475fc;
      background: #e6efff;
      .right-active-icon {
        opacity: 0;
      }
    }
    .status-block {
      font-size: 16px;
      flex-shrink: 0;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      font-weight: 600;
      min-width: 0;
      .node-name {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
      img {
        width: 20px;
        height: 20px;
        flex-shrink: 0;
      }
    }
    .time-tag {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 16px;
      padding: 0 6px;
      font-size: 12px;
      border-radius: 4px;
      color: #21a665;
      background: #CAFCE4;
    }
  }
}

.preview-box {
  flex: 1;
  width: var(--chat-test-log-width, auto);
  min-width: var(--chat-test-log-width, 0);
  max-width: var(--chat-test-log-width, none);
  box-sizing: var(--node-run-logs-panel-box-sizing, content-box);
  overflow-y: auto;
  padding: 16px;
  background: #fff;
  .preview-title {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
    .title-text {
      font-size: 15px;
      font-weight: 600;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      margin-left: 12px;
      .node-name{
        max-width: 200px;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
      img {
        width: 16px;
        height: 16px;
      }
    }
    .time-tag {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 16px;
      padding: 0 6px;
      font-size: 12px;
      border-radius: 4px;
      color: #21a665;
      background: #CAFCE4;
    }
  }
  .preview-content-block {
    margin-bottom:24px;
    .title-block {
      font-size: 15px;
      color: #262626;
      display: flex;
      align-items: center;
      gap: 4px;
      .anticon-copy {
        cursor: pointer;
        &:hover {
          color: #2475fc;
        }
      }
    }
    .preview-code-box {
      width: 100%;
      margin-top: 12px;
      padding: 10px 12px;
      border-radius: 8px;
      border: 1px solid #d9d9d9;
      overflow-x: auto;
      background: #fff;

      &::v-deep(.vjs-tree) {
        width: fit-content;
        min-width: 100%;
      }

      &::v-deep(.vjs-tree-node) {
        width: calc(100% + 16px);
        padding-right: 16px;
      }
    }
  }
}

.large-data-tip {
  font-size: 12px;
  color: #faad14;
  margin-bottom: 8px;
  padding: 4px 8px;
  background: #fffbe6;
  border: 1px solid #ffe58f;
  border-radius: 4px;
}

.large-data-textarea {
  width: 100%;
  min-height: 280px;
  max-height: 400px;
  resize: vertical;
  font-family: 'Courier New', Courier, monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #333;
  background: #fafafa;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  padding: 8px;
  overflow: auto;
  white-space: pre;
  outline: none;
}
</style>
