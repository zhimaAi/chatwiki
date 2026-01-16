<style lang="less" scoped>
.start-node {
  position: relative;
  margin-bottom: 8px;

  .start-node-options {
    display: flex;
    gap: 4px;
    overflow: hidden;
  }
  .options-title {
    line-height: 22px;
    margin-right: 8px;
    font-size: 14px;
    color: #262626;
    white-space: nowrap;
  }
  .options-list {
    flex: 1;
    display: flex;
    gap: 8px;
    overflow: hidden;
  }
  .options-item {
    display: flex;
    align-items: center;
    height: 22px;
    padding: 2px 2px 2px 4px;
    border-radius: 4px;
    border: 1px solid #d9d9d9;
    width: 100%;
    overflow: hidden;

    &.is-required .option-label::before {
      vertical-align: middle;
      content: '*';
      color: #fb363f;
      margin-right: 2px;
    }

    .option-label {
      color: var(--wf-color-text-3);
      font-size: 12px;
      margin-right: 4px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .option-type {
      height: 18px;
      line-height: 18px;
      padding: 0 8px;
      border-radius: 4px;
      font-size: 12px;
      background-color: #e4e6eb;
      color: var(--wf-color-text-3);
    }
  }
}
</style>

<template>
  <node-common
    :title="properties.node_name"
    :icon-url="properties.node_icon"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    style="width: 420px"
  >
    <div class="start-node">
      <div class="start-node-options">
        <div class="options-title">请求地址</div>
        <div class="options-list">
          <div class="options-item">
            <div class="option-label">{{ formState.url }}</div>
          </div>
        </div>
      </div>
    </div>
    <div class="start-node">
      <div class="start-node-options">
        <div class="options-title">请求方式</div>
        <div class="options-list">
          <div class="options-item">
            <div class="option-label">{{ formState.method }}</div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import NodeCommon from '../base-node.vue'
import { nextTick, onMounted, inject, watch, reactive } from 'vue'
const resetSize = inject('resetSize')
const props = defineProps({
  properties: {
    type: Object,
    default() {
      return {}
    }
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

watch(
  () => props.properties,
  (newVal, oldVal) => {
    const newDataRaw = newVal.node_params || '{}'
    const oldDataRaw = oldVal.node_params || '{}'
    if (newDataRaw != oldDataRaw) {
      reset()
    }
  },
  { deep: true }
)

const formState = reactive({
  url: '',
  method: 'GET',
  switch_verify: '1',
  switch_allow_ip: '0',
  allow_ips: '',

  params: [],
  form: [],

  x_form: [],

  json: [],

  request_content_type: 'none',

  response_type: 'now',
  response_now: ''
})

const reset = () => {
  if (!props.properties || !props.properties.node_params) {
    return
  }

  let node_params = JSON.parse(props.properties.node_params)

  let trigger_web_hook_config = node_params.trigger.trigger_web_hook_config
  // formState.url = trigger_web_hook_config.url
  // formState.method = trigger_web_hook_config.method || 'GET'
  // formState.switch_verify = trigger_web_hook_config.switch_verify || '1'
  // formState.switch_allow_ip = trigger_web_hook_config.switch_allow_ip || '0'
  // formState.allow_ips = trigger_web_hook_config.allow_ips
  //   ? trigger_web_hook_config.allow_ips.split(',').join('\n')
  //   : ''
  // formState.params = trigger_web_hook_config.params || []
  // formState.form = trigger_web_hook_config.form || []
  // formState.x_form = trigger_web_hook_config.x_form || []
  // formState.json = trigger_web_hook_config.json || []
  // formState.request_content_type = trigger_web_hook_config.request_content_type || 'none'
  // formState.response_type = trigger_web_hook_config.response_type || 'now'
  // formState.response_now = trigger_web_hook_config.response_now || ''

  Object.assign(formState, trigger_web_hook_config)

  nextTick(() => {
    resetSize()
  })
}

onMounted(() => {
  reset()
})
</script>
