<template>
  <a-modal
    width="746px"
    v-model:open="show"
    :confirmLoading="saveLoading"
    :title="`${formState.id ? '编辑' : '新建'}生效时间段`"
    @ok="handleSave"
    @cancel="onCancel"
  >
    <a-alert class="mt24" type="info" show-icon>
      <template #message>
        <div>优先匹配指定日期，然后匹配每周。比如设置了如下两条规则：</div>
        <div>
          规则一，每周一至周五，9：00-18:00，触发留资机器人规则二，2023年5月1日至2023年5月3日，13:30-19:00，触发留资机器人
        </div>
        <div>则访客2023年5月1日（周一）早上10：00咨询，不会触发留资机器人</div>
        <div>如果想要某天整天都不触发，可以将自动回复时间段设置为00:00-00:00</div>
      </template>
    </a-alert>
    <div class="form-box">
      <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-item label="类型" class="line-height-nomal">
          <a-radio-group v-model:value="formState.time_type">
            <a-radio :value="1">每周</a-radio>
            <a-radio :value="2">指定日期</a-radio>
          </a-radio-group>
          <div class="form-tip">节假日可使用指定时间，单独设置生效时间</div>
        </a-form-item>
        <a-form-item
          ref="name"
          label="指定每周"
          v-if="formState.time_type == 1"
          v-bind="validateInfos.week_numbers"
        >
          <a-checkbox-group v-model:value="formState.week_numbers" :options="weeekOptions" />
        </a-form-item>
        <a-form-item
          ref="name"
          label="指定日期"
          v-if="formState.time_type == 2"
          v-bind="validateInfos.dates"
        >
          <a-range-picker
            valueFormat="YYYY-MM-DD"
            style="width: 428px"
            v-model:value="formState.dates"
          />
        </a-form-item>
        <a-form-item ref="name" label="自动回复时间段" v-bind="validateInfos.time_list">
          <div class="keyword-list-box">
            <div class="list-item" v-for="(item, index) in formState.time_list" :key="index">
              <a-time-range-picker
                valueFormat="HH:mm"
                format="HH:mm"
                :allowClear="false"
                v-model:value="item.value"
              />
              <div class="btn-hover-wrap">
                <CloseCircleOutlined style="color: #595959" @click="delKeywordItem(index)" />
              </div>
            </div>
            <div class="add-btn-box">
              <a-button type="dashed" block @click="addKeywordItem">
                <template #icon>
                  <PlusOutlined />
                </template>
                添加时间段
              </a-button>
            </div>
          </div>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, h } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import {
  CloseCircleFilled,
  CloseCircleOutlined,
  LoadingOutlined,
  PlusOutlined,
} from '@ant-design/icons-vue'
import { saveField } from '@/api/retention/index.js'
import { useRoute } from 'vue-router'
const query = useRoute().query

const emit = defineEmits(['ok'])

const props = defineProps({
  tableData: {
    default: () => [],
    type: Array,
  },
})

const useForm = Form.useForm

const labelCol = {
  span: 5,
}
const wrapperCol = {
  span: 19,
}

const show = ref(false)

const saveLoading = ref(false)

const formState = reactive({
  time_type: 1,
  dates: [],
  week_numbers: [],
  time_list: [{ value: [] }],
})

const weeekOptions = [
  { label: '周一', value: 1 },
  { label: '周二', value: 2 },
  { label: '周三', value: 3 },
  { label: '周四', value: 4 },
  { label: '周五', value: 5 },
  { label: '周六', value: 6 },
  { label: '周日', value: 7 },
]

const rules = reactive({
  dates: [
    {
      required: true,
      validator: async (rule, value) => {
        if (formState.time_type == 2) {
          if (!value) {
            return Promise.reject('请选择指定日期')
          }
        }
        return Promise.resolve()
      },
    },
  ],
  week_numbers: [
    {
      required: true,
      validator: async (rule, value) => {
        if (formState.time_type == 1) {
          if (value.length == 0) {
            return Promise.reject('请至少选择一个选项')
          }
        }
        return Promise.resolve()
      },
    },
  ],
  time_list: [
    {
      required: true,
      validator: async (rule, value) => {
        if (value.length == 0) {
          return Promise.reject('请设置一个时间段')
        }
        let list = value.map((item) => item.value.join(''))
        if (list.join('') == '') {
          return Promise.reject('请设置一个时间段')
        }
        if (isTimeRangesOverlapping(value)) {
          return Promise.reject('时间范围存在重叠')
        }
        return Promise.resolve()
      },
    },
  ],
})

const addKeywordItem = () => {
  formState.time_list.push({ value: [] })
}
const delKeywordItem = (index) => {
  // 移除
  if (formState.time_list.length == 1) {
    return message.error('至少保留一个时间段')
  }
  formState.time_list.splice(index, 1)
}

const { validate, validateInfos, resetFields } = useForm(formState, rules)

const saveForm = () => {
  let data = {
    time_type: formState.time_type,
  }
  if (formState.time_type == 1) {
    data.week_numbers = formState.week_numbers.join(',')
  } else {
    data.date_start = formState.dates[0]
    data.date_end = formState.dates[1]
  }
  let time_list = formState.time_list.map((item) => {
    return {
      start: item.value[0],
      end: item.value[1],
    }
  })
  data.time_list = time_list

  //  校验一下 有没有重复的时间范围
  let tableData = props.tableData || []
  if (timeIndex >= 0) {
    tableData = tableData.filter((item, index) => index != timeIndex)
  }
  let checkLists = tableData.filter((item) => item.time_type == formState.time_type)
  if (formState.time_type == 1) {
    // 指定周
    let week_all = checkLists.map((item) => (item.week_numbers ? item.week_numbers.split(',') : []))
    week_all = [...new Set(week_all.flat())]
    week_all = week_all.map((item) => +item)
    let errorFlag = ''
    formState.week_numbers.forEach((item) => {
      if (week_all.includes(item) && errorFlag == '') {
        errorFlag = '检查到你已经设置指定周' + item + ',请先删除后再添加'
      }
    })
    if (errorFlag) {
      return message.error(errorFlag)
    }
  } else {
    const startDate = new Date(data.date_start)
    const endDate = new Date(data.date_end)
    const isOverlapping = hasOverlap(checkLists, startDate, endDate)
    if (isOverlapping) {
      return message.error('检测到你当前设置的指定日期和已经设置的指定日期存在重叠,请先修改后再试')
    }
  }
  resetFields()
  show.value = false
  emit('ok', data)
}
let timeIndex = null
const onShow = async (data, index) => {
  timeIndex = index
  formState.time_type = +data.time_type || 1
  formState.dates = []
  formState.week_numbers = []
  if (data.date_start || data.date_end) {
    formState.dates = [data.date_start, data.date_end]
  }
  if (data.week_numbers) {
    formState.week_numbers = data.week_numbers.split(',').map((item) => +item)
  }
  if (data.time_list && data.time_list.length > 0) {
    formState.time_list = []
    data.time_list.forEach((item) => {
      formState.time_list.push({
        value: [item.start, item.end],
      })
    })
  } else {
    formState.time_list = [{ value: [] }]
  }

  show.value = true
}

const onCancel = () => {
  resetFields()
  show.value = false
}

const handleSave = () => {
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

const addOptions = () => {
  formState.options.push({
    value: '',
  })
}
const deleteOptions = (index) => {
  formState.options.splice(index, 1)
}

defineExpose({
  onShow,
})

function isTimeRangesOverlapping(lists) {

  let timeRanges = []
  lists.forEach((item) => {
    timeRanges.push(item.value)
  })
  timeRanges = timeRanges.filter((item) => item.length > 0)
  if (timeRanges.length == 1 && timeRanges[0] == null) {
    new Error('请选择时间范围')
    return false
  }
  const hasFullTimeRange = timeRanges.some((range) => range[0] === '00:00' && range[1] === '00:00')

  if (hasFullTimeRange && timeRanges.length > 1 && timeRanges[1].length) {
    // 如果存在 [00:00, 00:00]，则认为所有时间段都重叠
    return true
  }
  // 辅助函数，将时间字符串转换为分钟数
  function timeToMinutes(time) {
    const [hours, minutes] = time.split(':').map(Number)
    return hours * 60 + minutes
  }

  // 遍历时间范围数组
  for (let i = 1; i < timeRanges.length; i++) {
    if (timeRanges[i].length == 0) {
      continue
    }
    const currentStart = timeToMinutes(timeRanges[i][0])
    const previousEnd = timeToMinutes(timeRanges[i - 1][1])

    // 如果当前时间段的开始时间小于或等于前一个时间段的结束时间，则存在重叠
    if (currentStart < previousEnd) {
      return true
    }
  }

  // 如果没有找到重叠的时间段，则返回false
  return false
}

function hasOverlap(timeRanges, startDate, endDate) {
  for (let i = 0; i < timeRanges.length; i++) {
    const rangeStartDate = new Date(timeRanges[i].date_start)
    const rangeEndDate = new Date(timeRanges[i].date_end)

    // 检查重叠条件
    if (!(endDate <= rangeStartDate || startDate >= rangeEndDate)) {
      return true // 有重叠
    }
  }
  return false // 无重叠
}
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
  .line-height-nomal {
    &::v-deep(.ant-form-item-label label) {
      height: 22px;
    }
    &::v-deep(.ant-form-item-control-input) {
      min-height: 22px;
    }
  }
}
.keyword-list-box {
  width: 452px;
  .list-item {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
    .ant-picker {
      flex: 1;
    }
    .anticon-close-circle {
      color: #595959;
      font-size: 16px;
      cursor: pointer;
    }
  }
  .add-btn-box {
    width: 412px;
  }
}
.form-tip {
  color: #8c8c8c;
  line-height: 22px;
  font-size: 14px;
  font-weight: 400;
  margin-top: 2px;
}
.btn-hover-wrap {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease-in;
  &:hover {
    background: #e4e6eb;
  }
}
</style>
