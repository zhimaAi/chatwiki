<template>
  <div class="zm-date">
    <div class="date-btn">
      <a-radio-group v-model:value="create_time_type" @change="handleTagChange">
        <a-radio-button class="date-btn-button" v-for="item in tagDateArr" :value="item.key" :key="item.key">
          {{item.label}}
        </a-radio-button>
      </a-radio-group>
    </div>
    <div class="date-content">
      <a-range-picker v-model:value="range" @change="handleDateChange" format="YYYY-MM-DD" :disabledDate=disabledDate :allowClear="false"></a-range-picker>
    </div>
  </div>
</template>
<script setup>
import dayjs from 'dayjs';
import { ref, onMounted, reactive, watch } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.statistical-analysis.components.date')

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
    label: t('label_last_7_days'),
    value: true,
    key: 2,
    start_date: dayjs().subtract(6, 'day'),
    end_date: dayjs().subtract(0, 'day'),
  },
  {
    label: t('label_last_14_days'),
    value: false,
    key: 3,
    start_date: dayjs().subtract(13, 'day'),
    end_date: dayjs().subtract(0, 'day'),
  },
  {
    label: t('label_last_30_days'),
    value: false,
    key: 4,
    start_date: dayjs().subtract(29, 'day'),
    end_date: dayjs().subtract(0, 'day'),
  }
])

const start_time = ref(null)
const end_time = ref(null)
const range = ref([])

// 创建一个禁用今天之后日期的函数
const disabledDate = (current) => {
  // 如果当前日期大于或等于今天，则禁用
  return current >= dayjs().subtract(0, 'day');
};

const handleTagChange = () => {
  let key = create_time_type.value

  let dateItem = null;
  tagDateArr.forEach((item) => {
    if (item.key == key) {
      dateItem = item;
    }
  })

  range.value = [dateItem.start_date, dateItem.end_date];
  start_time.value = dateItem.start_date.format('YYYY-MM-DD')
  end_time.value = dateItem.end_date.format('YYYY-MM-DD')
  search()
}

const handleDateChange = (dates, dateStrings) => {
  const startDate = dayjs(dates[0]);
  const endDate = dayjs(dates[1]);
  if (endDate.diff(startDate, 'days') > 29) {
    range.value[0] = dates[0];
    range.value[1] = startDate.add(29, 'days')
    start_time.value = range.value[0].format('YYYY-MM-DD')
    end_time.value = range.value[1].format('YYYY-MM-DD')
    message.error(t('error_max_30_days'))
  } else {
    range.value = dates
    start_time.value = dateStrings[0]
    end_time.value = dateStrings[1]
  }
  create_time_type.value = 1;
  search();
}

const search = () => {
  emit("dateChange", {start_date: start_time.value, end_date: end_time.value})
}

onMounted(() => {
  create_time_type.value = parseInt(props.datekey)
  handleTagChange()
})

watch(() => props.datekey, (newval, oldVla) => {
  create_time_type.value = parseInt(props.datekey.split('-')[0])
  handleTagChange()
})

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
}


</style>