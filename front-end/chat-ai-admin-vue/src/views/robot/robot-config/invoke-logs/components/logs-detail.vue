<template>
  <a-modal
    v-model:open="open"
    :title="t('title_log_detail')"
    :footer="null"
    :width="880"
    wrapClassName="no-padding-modal"
    :bodyStyle="{ 'max-height': '600px', 'padding-right': '24px', 'overflow-y': 'auto' }"
  >
    <div class="flex-content-box">
      <div class="test-model-box">
        <div class="top-title">{{ t('label_start_node_params') }}</div>
        <a-form
          :model="formState"
          ref="formRef"
          layout="vertical"
          :wrapper-col="{ span: 24 }"
          autocomplete="off"
        >
          <a-form-item name="question" :rules="[{ required: true, message: t('msg_input_question') }]">
            <template #label>
              <a-flex :gap="4">question <a-tag style="margin: 0">string</a-tag> </a-flex>
            </template>
            <a-input readonly :placeholder="t('ph_input')" v-model:value="formState.question" />
          </a-form-item>
          <a-form-item name="openid" :rules="[{ required: true, message: t('msg_input_openid') }]">
            <template #label>
              <a-flex :gap="4">openid <a-tag style="margin: 0">string</a-tag> </a-flex>
            </template>
            <a-input readonly :placeholder="t('ph_input')" v-model:value="formState.openid" />
          </a-form-item>
        </a-form>

        <div class="result-list-box" v-if="resultList.length > 0">
          <div
            class="list-item-block"
            :class="{ active: currentNodeKey == item.node_key }"
            v-for="(item, index) in resultList"
            @click="handleChangeNodeKey(item)"
            :key="index"
          >
            <div class="status-block">
              <CheckCircleFilled v-if="item.is_success" style="color: #138b1b" />
              <CloseCircleFilled v-else style="color: #d81e06" />
            </div>
            <div class="icon-name-box">
              <img :src="item.node_icon" alt="" />
              <div class="node-name">{{ item.node_name }}</div>
            </div>
            <div class="time-tag" v-if="item.is_success">{{ item.use_time }}ms</div>
            <div class="right-active-icon"><RightCircleOutlined /></div>
          </div>
        </div>
      </div>
      <div class="preview-box">
        <template v-if="cuttentItem">
          <div class="preview-title">
            <div class="title-text">{{ t('title_log_detail') }}</div>
            <div class="icon-name-box">
              <img :src="cuttentItem.node_icon" alt="" />
              <div class="node-name">{{ cuttentItem.node_name }}</div>
            </div>
            <div class="time-tag" v-if="cuttentItem.is_success">{{ cuttentItem.use_time }}ms</div>
          </div>
          <div class="preview-content-block">
            <div class="title-block">{{ t('label_input') }}<CopyOutlined @click="handleCopy('input')" /></div>
            <div class="preview-code-box">
              <template v-if="isHugeData(cuttentItem.input)">
                <div class="large-data-tip">{{ t('msg_data_too_large') }}</div>
                <textarea class="large-data-textarea" readonly :value="getJsonText(cuttentItem.input)" />
              </template>
              <vue-json-pretty v-else :data="cuttentItem.input" v-bind="getJsonViewerProps(cuttentItem.input)" />
            </div>
          </div>
          <div class="preview-content-block">
            <div class="title-block">{{ t('label_output') }}<CopyOutlined @click="handleCopy('node_output')" /></div>
            <div class="preview-code-box">
              <template v-if="isHugeData(cuttentItem.node_output)">
                <div class="large-data-tip">{{ t('msg_data_too_large') }}</div>
                <textarea class="large-data-textarea" readonly :value="getJsonText(cuttentItem.node_output)" />
              </template>
              <vue-json-pretty v-else :data="cuttentItem.node_output" v-bind="getJsonViewerProps(cuttentItem.node_output)" />
            </div>
          </div>
          <div class="preview-content-block">
            <div class="title-block">{{ t('label_run_log') }}<CopyOutlined @click="handleCopy('output')" /></div>
            <div class="preview-code-box">
              <template v-if="isHugeData(cuttentItem.output)">
                <div class="large-data-tip">{{ t('msg_data_too_large') }}</div>
                <textarea class="large-data-textarea" readonly :value="getJsonText(cuttentItem.output)" />
              </template>
              <vue-json-pretty v-else :data="cuttentItem.output" v-bind="getJsonViewerProps(cuttentItem.output)" />
            </div>
          </div>
        </template>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {
  CheckCircleFilled,
  CloseCircleFilled,
  RightCircleOutlined,
  CopyOutlined
} from '@ant-design/icons-vue'
import VueJsonPretty from 'vue-json-pretty'
import 'vue-json-pretty/lib/styles.css'
import { reactive, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getImageUrl } from '@/views/workflow/components/util.js'
import { message } from 'ant-design-vue'
import { copyText } from '@/utils/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.invoke-logs.components.logs-detail')
const emit = defineEmits(['save', 'getGlobal'])
const query = useRoute().query

const open = ref(false)
const currentNodeKey = ref('')
const resultList = ref([])

const cuttentItem = computed(() => {
  if (!currentNodeKey.value) {
    return null
  }
  return resultList.value.filter((item) => item.node_key == currentNodeKey.value)[0]
})

const formState = reactive({
  robot_key: query.robot_key,
  question: '',
  openid: ''
})

const show = (record) => {
  formState.question = record.question || ''
  formState.openid = record.openid || ''
  resultList.value = []
  currentNodeKey.value = ''
  formateData(record.node_logs)
  open.value = true
}

const formateData = (data) => {
  resultList.value = data.map((item) => {
    return {
      ...item,
      is_success: item.error_msg === '<nil>',
      node_icon: getImageUrl(item.node_type)
    }
  })
  currentNodeKey.value = resultList.value[0]?.node_key
}
const handleChangeNodeKey = (item) => {
  currentNodeKey.value = item.node_key
}

const handleCopy = (key) => {
  copyText(JSON.stringify(cuttentItem.value[key]))
  message.success(t('msg_copy_success'))
}

// ---- JSON 数据量安全渲染辅助 ----
const JSON_COLLAPSE_THRESHOLD = 5 * 1024    // 5 KB  → 折叠到 depth=2
const JSON_VIRTUAL_THRESHOLD  = 100 * 1024  // 100 KB → 开启虚拟滚动
const JSON_HUGE_THRESHOLD     = 500 * 1024  // 500 KB → 降级为纯文本

function getJsonSize(data) {
  if (data == null) return 0
  try { return JSON.stringify(data).length } catch { return 0 }
}

function isHugeData(data) {
  return getJsonSize(data) > JSON_HUGE_THRESHOLD
}

/** 根据数据大小返回 vue-json-pretty 的合适 props */
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
// ---- end ----

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.flex-content-box {
  display: flex;
}
.test-model-box {
  flex: 1;
  margin: 24px 24px 0 0;
  .top-title {
    font-weight: 600;
    margin-bottom: 16px;
  }
  .save-btn-box {
    margin: 32px 0;
    margin-top: 50px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
.tooltip-content {
  white-space: pre-wrap;
}
.loading-box {
  height: 100px;
  justify-content: center;
}
.result-list-box {
  margin: 24px 0;
  width: 100%;
  border: 1px solid #ebebeb;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  padding: 8px;
  .list-item-block {
    display: flex;
    align-items: center;
    overflow: hidden;
    gap: 8px;
    padding: 8px;
    color: #333;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    .right-active-icon {
      margin-left: auto;
      color: #2475fc;
      opacity: 0;
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
      font-size: 20px;
    }
    .icon-name-box {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 14px;
      font-weight: 600;
      img {
        width: 24px;
        height: 24px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
    .out-put-box {
      flex: 1;
      margin-left: 24px;
      overflow: hidden;
      .out-text-box {
        background: #f2f2f2;
        border-radius: 6px;
        padding: 8px;
        width: 100%;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.preview-box {
  flex: 1;
  border-left: 1px solid #d9d9d9;
  padding: 16px;
  padding-right: 0;
  .preview-title {
    display: flex;
    align-items: center;
    gap: 8px;
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
      img {
        width: 16px;
        height: 16px;
      }
    }
    .time-tag {
      width: fit-content;
      border-radius: 4px;
      height: 22px;
      background: #d2f1dc;
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0 4px;
      font-size: 12px;
    }
  }
  .preview-content-block {
    margin-top: 16px;
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
      width: fit-content;
      min-width: 100%;
      margin-top: 16px;
      padding: 8px;
      border-radius: 8px;
      border: 1px solid #d9d9d9;

      &::v-deep(.vjs-tree) {
        width: fit-content;
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
