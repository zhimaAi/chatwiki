<template>
  <NodeFormLayout>
    <template #header>
      <NodeFormHeader
        :title="node.node_name"
        :iconName="node.node_icon_name"
        @changeTitle="handleTitleChange"
        @deleteNode="handleDeleteNode"
        desc="在设定的时间点，自动执行工作流"
      >
      </NodeFormHeader>
    </template>

    <div class="variable-node">
      <div class="node-form-content">
        <div class="gray-block">
          <div class="output-label">
            <img src="@/assets/svg/alarm-clock.svg" alt="" class="output-label-icon" />
            <span class="output-label-text">触发时间</span>
          </div>
          <div class="form-item">
            <a-radio-group v-model:value="formState.type">
              <a-radio value="select_time">选择触发时间</a-radio>
              <a-radio value="linux_crontab">Linux Crontab 代码</a-radio>
            </a-radio-group>
          </div>
          <div class="form-item" v-if="formState.type === 'select_time'">
            <a-flex :gap="8">
              <a-cascader
                v-model:value="formState.time_value"
                :options="options"
                placeholder="请选择"
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
              placeholder="请输入代码"
              style="min-height: 142px; background: #262626; color: #bfbfbf"
            />
          </div>
        </div>
      </div>
    </div>
  </NodeFormLayout>
</template>

<script setup>
import { ref, onMounted, inject, reactive, computed, watch } from 'vue'
import NodeFormLayout from '../node-form-layout.vue'
import NodeFormHeader from '../node-form-header.vue'
import CodeEditBox from './code-edit-box.vue'

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
      label: i + 1 + '号',
      value: i + 1 + ''
    })
  }
  return lists
}

const options = [
  {
    value: 'day',
    label: '每天执行'
  },
  {
    value: 'week',
    label: '每周执行',
    children: [
      {
        value: '0',
        label: '星期日'
      },
      {
        value: '1',
        label: '星期一'
      },
      {
        value: '2',
        label: '星期二'
      },
      {
        value: '3',
        label: '星期三'
      },
      {
        value: '4',
        label: '星期四'
      },
      {
        value: '5',
        label: '星期五'
      },
      {
        value: '6',
        label: '星期六'
      }
    ]
  },
  {
    value: 'month',
    label: '每月执行',
    children: getMonthDays()
  }
]

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
