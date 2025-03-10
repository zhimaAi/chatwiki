<template>
  <div class="zm-date">
    <div class="date-btn">
      <a-radio-group v-model:value="create_time_type" @change="handleTagChange">
        <a-radio-button
          class="date-btn-button"
          v-for="item in tagDateArr"
          :value="item.key"
          :key="item.key"
        >
          {{ item.label }}
        </a-radio-button>
      </a-radio-group>
    </div>
    <div class="date-content">
      <a-range-picker
        v-model:value="range"
        @change="handleDateChange"
        format="YYYY-MM-DD"
        :disabledDate="disabledDate"
        :allowClear="false"
      ></a-range-picker>
    </div>
  </div>
</template>
<script setup>
import dayjs from 'dayjs'
import { ref, onMounted, reactive, watch } from 'vue'
import { message } from 'ant-design-vue'

const props = defineProps({
  datekey: {
    type: String,
    default: '1'
  }
})

const create_time_type = ref(1)

const emit = defineEmits(['dateChange'])

let tagDateArr = reactive([
  {
    label: '今日',
    value: true,
    key: 2,
    start_time: dayjs().startOf('day'),
    end_time: dayjs().endOf('day')
  },
  {
    label: '昨日',
    value: false,
    key: 3,
    start_time: dayjs().subtract(1, 'day').startOf('day'),
    end_time: dayjs().startOf('day').subtract(1, 'millisecond')
  },
  {
    label: '近7日',
    value: false,
    key: 4,
    start_time: dayjs().subtract(6, 'day').startOf('day'),
    end_time: dayjs().endOf('day').subtract(1, 'millisecond')
  }
])

const start_time = ref(null)
const end_time = ref(null)
const range = ref([])

// 创建一个禁用今天之后日期的函数
const disabledDate = (current) => {
  // 如果当前日期大于或等于今天，则禁用
  return current >= dayjs().subtract(0, 'day')
}

const handleTagChange = () => {
  let key = create_time_type.value

  let dateItem = null
  tagDateArr.forEach((item) => {
    if (item.key == key) {
      dateItem = item
    }
  })

  if (dateItem) {
    range.value = [dateItem.start_time, dateItem.end_time]
    start_time.value = dateItem.start_time.unix()
    end_time.value = dateItem.end_time.unix()
    // format('YYYY-MM-DD HH:mm:ss')
  }

  search()
}

const handleDateChange = (dates) => {
  const startDate = dayjs(dates[0]).startOf('day')
  const endDate = dayjs(dates[1]).endOf('day')
  if (endDate.diff(startDate, 'days') > 29) {
    range.value[0] = dates[0].startOf('day')
    range.value[1] = startDate.add(29, 'days').endOf('day')
    start_time.value = range.value[0].unix()
    end_time.value = range.value[1].unix()
    message.error(`最多只能选择30天`)
  } else {
    const start = dayjs(dates[0]).startOf('day').unix()
    const end = dayjs(dates[1]).endOf('day').unix()
    range.value = dates
    start_time.value = start
    end_time.value = end
  }
  create_time_type.value = 1
  search()
}

const search = () => {
  emit('dateChange', { start_time: start_time.value, end_time: end_time.value })
}

onMounted(() => {
  create_time_type.value = parseInt(props.datekey)
  handleTagChange()
})

watch(
  () => props.datekey,
  (newval, oldVla) => {
    create_time_type.value = parseInt(props.datekey.split('-')[0])
    handleTagChange()
  }
)
</script>
<style lang="less">
.zm-date {
  display: flex;
  align-items: center;

  .date-btn {
    margin-right: 8px;

    .date-btn-button {
      padding-left: 16px;
      padding-right: 16px;
    }
  }

  .date-content {
    width: 260px;
  }
}
</style>
