<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        @changeTitle="handleTitleChange"
        @deleteNode="handleDeleteNode"
        :desc="t('desc_timing_trigger')"
      >
      </NodeFormHeader>
    </template>

    <div class="variable-node">
      <div class="node-form-content">
        <div class="gray-block">
          <div class="output-label">
            <img src="@/assets/svg/alarm-clock.svg" alt="" class="output-label-icon" />
            <span class="output-label-text">{{ t('label_trigger_time') }}</span>
          </div>
          <div class="form-item">
            <a-radio-group v-model:value="formState.type">
              <a-radio value="select_time">{{ t('label_select_time') }}</a-radio>
              <a-radio value="linux_crontab">{{ t('label_linux_crontab') }}</a-radio>
            </a-radio-group>
          </div>
          <div class="form-item" v-if="formState.type === 'select_time'">
            <a-flex :gap="8">
              <a-cascader
                v-model:value="formState.time_value"
                :options="options"
                :placeholder="t('ph_select')"
                style="width: 194px"
              />
              <a-time-picker
                :allowClear="false"
                v-model:value="formState.hour_minute"
                valueFormat="HH:mm"
                format="HH:mm"
                style="width: 124px"
              />
            </a-flex>
          </div>
          <div class="form-item custom-textarea" v-else>
            <a-textarea
              v-model:value="formState.linux_crontab"
              :placeholder="t('ph_input_code')"
              style="min-height: 142px; background: #262626; color: #bfbfbf"
            />
          </div>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { ref, onMounted, inject, reactive, computed, watch } from 'vue'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import CodeEditBox from './code-edit-box.vue'

const { t } = useI18n('views.workflow.components.node-form-drawer.timing-trigger-node.timing-trigger-node-form')

const emit = defineEmits(['update-node'])
const props = defineProps({
  lf: {
    type: Object,
    default: null
  },
  nodeId: {
    type: String,
    default: ''
  },
  node: {
    type: Object,
    default: () => ({})
  }
})

function getMonthDays() {
  let lists = []
  for (let i = 0; i < 31; i++) {
    lists.push({
      label: i + 1 + t('label_day'),
      value: i + 1 + ''
    })
  }
  return lists
}

const options = computed(() => [
  {
    value: 'day',
    label: t('label_every_day')
  },
  {
    value: 'week',
    label: t('label_every_week'),
    children: [
      {
        value: '0',
        label: t('label_sunday')
      },
      {
        value: '1',
        label: t('label_monday')
      },
      {
        value: '2',
        label: t('label_tuesday')
      },
      {
        value: '3',
        label: t('label_wednesday')
      },
      {
        value: '4',
        label: t('label_thursday')
      },
      {
        value: '5',
        label: t('label_friday')
      },
      {
        value: '6',
        label: t('label_saturday')
      }
    ]
  },
  {
    value: 'month',
    label: t('label_every_month'),
    children: getMonthDays()
  }
])

const formState = reactive({
  type: 'select_time',
  time_value: [],
  every_type: 'day',
  week_number: '',
  month_day: '',
  hour_minute: '',
  linux_crontab: ''
})

watch(
  () => formState,
  () => {
    update()
  },
  {
    deep: true
  }
)

const getNode = inject('getNode')
const getGraph = inject('getGraph')

const handleTitleChange = () => {
  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change', { ...props.node })
  }, 10)
}

const handleDeleteNode = () => {
  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change', null)
  }, 10)
}

const update = () => {
  let node_params = JSON.parse(props.node.node_params)
  let every_type = ''
  let week_number = ''
  let month_day = ''
  if (formState.time_value && formState.time_value.length) {
    every_type = formState.time_value[0]
    if (formState.time_value[0] == 'week') {
      week_number = formState.time_value[1]
    }
    if (formState.time_value[0] == 'month') {
      month_day = formState.time_value[1]
    }
  }

  node_params.trigger.cron_config = {
    ...formState,
    every_type,
    week_number,
    month_day
  }

  let data = { ...props.node, node_params: JSON.stringify(node_params) }

  emit('update-node', data)

  setTimeout(() => {
    getGraph().eventCenter.emit('custom:trigger-change', data)
  }, 10)
}

const init = () => {
  try {
    let dataRaw = props.node.dataRaw || props.node.node_params || '{}'

    dataRaw = JSON.parse(dataRaw)

    let cron_config = dataRaw.trigger.cron_config
    formState.type = cron_config.type
    formState.linux_crontab = cron_config.linux_crontab
    formState.every_type = cron_config.every_type
    formState.week_number = cron_config.week_number
    formState.month_day = cron_config.month_day
    formState.hour_minute = cron_config.hour_minute
    let every_type = cron_config.every_type
    if (every_type == 'day') {
      formState.time_value = ['day']
    }
    if (every_type == 'week') {
      formState.time_value = ['week', cron_config.week_number]
    }
    if (every_type == 'month') {
      formState.time_value = ['month', cron_config.month_day]
    }
  } catch (error) {
    console.log(error)
  }
}

onMounted(() => {
  init()
})
</script>

<style lang="less" scoped>
@import '../form-block.less';
.variable-node {
  .output-label {
    margin-bottom: 16px;
  }
  .form-item {
    margin-bottom: 4px;
  }
  .custom-textarea {
    &::v-deep(.ant-input::placeholder) {
      color: #bfbfbf !important;
    }
  }
}
</style>
