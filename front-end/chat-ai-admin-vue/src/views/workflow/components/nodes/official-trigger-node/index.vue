<style lang="less" scoped>
.start-node {
  position: relative;

  .start-node-options {
    display: flex;
    gap: 4px;
    margin-top: 12px;
  }
  .options-title {
    line-height: 22px;
    margin-right: 8px;
    font-size: 14px;
    color: #262626;
  }
  .options-list {
    flex: 1;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }
  .options-item {
    display: flex;
    align-items: center;
    height: 22px;
    padding: 2px 2px 2px 4px;
    border-radius: 4px;
    border: 1px solid #d9d9d9;

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
    :properties="properties"
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
        <div class="options-title">{{ t('label_trigger_event') }}</div>
        <div class="options-list">
          <div class="options-item">
            <div class="option-label">{{ t(msg_type_map[formState.msg_type]) }}</div>
          </div>
        </div>
      </div>
      <div class="start-node-options">
        <div class="options-title" @click="test">{{ t('label_official_account') }}</div>
        <div class="options-list">
          <div class="options-item" v-for="item in selectAppItems.slice(0, 1)" :key="item.app_id">
            <div class="option-label">{{ item.app_name }}</div>
          </div>
          <div class="options-item" v-if="selectAppItems.length > 1">
            <div class="option-label">...</div>
          </div>
        </div>
      </div>
      <div class="start-node-options">
        <div class="options-title">{{ t('label_output') }}</div>
        <div class="options-list">
          <div
            class="options-item"
            :class="{ 'is-required': item.required }"
            v-for="item in options.slice(0, 2)"
            :key="item.key"
          >
            <div class="option-label">{{ item.key }}</div>
            <div class="option-type">{{ item.typ }}</div>
          </div>
          <div class="options-item" v-if="options.length > 2">
            <div class="option-label">...</div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import NodeCommon from '../base-node.vue'
import { nextTick, onMounted, inject, watch, reactive, ref, computed } from 'vue'
import { useWorkflowStore } from '@/stores/modules/workflow'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.nodes.official-trigger-node.index')

const workflowStore = useWorkflowStore()

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

const officialList = computed(() => workflowStore.officialList)
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

let msg_type_map = {
  message: 'text_private_message',
  subscribe_unsubscribe: 'text_subscribe_unsubscribe',
  qrcode_scan: 'text_qrcode_scan',
  menu_click: 'text_menu_click'
}

const formState = reactive({
  msg_type: '',
  app_ids: []
})

const selectAppItems = computed(() => {
  return formState.app_ids.map((item) => officialList.value.find((it) => it.app_id == item)).filter(Boolean)
})

const outputs = ref([])

const options = computed(() => {
  return [...outputs.value]
})

const reset = () => {
  if (!props.properties || !props.properties.node_params) {
    return
  }

  let node_params = JSON.parse(props.properties.node_params)

  outputs.value = node_params.trigger.outputs || []

  let trigger_official_config = node_params.trigger.trigger_official_config

  formState.msg_type = trigger_official_config.msg_type
  formState.app_ids = trigger_official_config.app_ids
    ? trigger_official_config.app_ids.split(',')
    : []

  nextTick(() => {
    resetSize()
  })
}

const test = ()=>{
  let aa = officialList.value.find((it) => it.app_id == 'sdwdw')
  console.log(selectAppItems.value,aa, '==')
}

onMounted(() => {
  reset()
})
</script>
