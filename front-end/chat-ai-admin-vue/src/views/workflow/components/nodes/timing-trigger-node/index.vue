<style lang="less" scoped>
.start-node {
  position: relative;

  .start-node-options {
    display: flex;
    gap: 4px;
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
      <div class="start-node-options" v-if="formState.type == 'select_time'">
        <div class="options-title">{{ t('label_execution_time') }}</div>
        <div class="options-list">
          <div class="options-item">
            <div class="option-label" v-if="formState.every_type == 'day'">
              {{ t('text_every_day') }} {{ formState.hour_minute }}
            </div>
            <div class="option-label" v-if="formState.every_type == 'week'">
              {{ t('text_every_week') }} {{ t(week_number_map[formState.week_number]) }} {{ formState.hour_minute }}
            </div>
            <div class="option-label" v-if="formState.every_type == 'month'">
              {{ t('text_every_month') }}
              <span v-if="formState.month_day">{{ formState.month_day }}{{ t('text_day_suffix') }}</span>
              {{ formState.hour_minute }}
            </div>
          </div>
        </div>
      </div>
      <div class="start-node-options" v-else>
        <div class="options-title">{{ t('label_execution_time') }}</div>
        <div class="options-list">
          <div class="options-item">
            <div class="option-label">{{ t('text_linux_crontab_code') }}</div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script setup>
import NodeCommon from '../base-node.vue'
import { nextTick, onMounted, inject, watch, reactive } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.nodes.timing-trigger-node.index')

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

let week_number_map = {
  1: 'text_monday',
  2: 'text_tuesday',
  3: 'text_wednesday',
  4: 'text_thursday',
  5: 'text_friday',
  6: 'text_saturday',
  0: 'text_sunday'
}
const formState = reactive({
  type: 'select_time',
  every_type: 'day',
  week_number: '',
  month_day: '',
  hour_minute: '',
  linux_crontab: ''
})

const reset = () => {
  if (!props.properties || !props.properties.node_params) {
    return
  }

  let node_params = JSON.parse(props.properties.node_params)
  let cron_config = node_params.trigger.cron_config || {}

  Object.assign(formState, cron_config)

  nextTick(() => {
    resetSize()
  })
}

onMounted(() => {
  reset()
})
</script>
