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
      <a-range-picker v-model:value="range" @change="handleDateChange" format="YYYY-MM-DD" :disabledDate=disabledDate></a-range-picker>
    </div>
  </div>
</template>
<script setup>
import dayjs from 'dayjs';
import { ref, onMounted, reactive } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.usetoken')

const props = defineProps({
  dateType: {
    type: Number,
    default: 1
  },
  rangeData: {
    type: Array,
    default: null
  }
})

const create_time_type = ref(1)

const emit = defineEmits(['dateChange'])

let tagDateArr = reactive([
  {
    label: t('date.today'),
    value: false,
    key: 2,
    start_date: dayjs(),
    end_date: dayjs(),
  },
  {
    label: t('date.yesterday'),
    value: false,
    key: 3,
    start_date: dayjs().subtract(1, 'day'),
    end_date: dayjs().subtract(1, 'day'),
  },
  {
    label: t('date.last_seven_days'),
    value: true,
    key: 4,
    start_date: dayjs().subtract(6, 'day'),
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
  create_time_type.value = 1;
  range.value = dates
  start_time.value = dateStrings[0]
  end_time.value = dateStrings[1]
  search();
}
const search = () => {
  emit("dateChange", {start_date: start_time.value, end_date: end_time.value})
}

onMounted(() => {
  create_time_type.value = props.dateType
  range.value = props.rangeData
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

  .date-content {
    width: 280px;
  }
}


</style>