<template>
  <div class="zm-date-quick">
    <div class="btns">
      <template v-for="item in tagDateArr" :key="item.key">
        <span
          class="btn"
          :class="{ active: item.key === create_time_type }"
          @click="handleTagChange(item.key)"
        >
          {{ t(item.label) }}
          <DownOutlined v-if="item.key === 5" class="down-icon" :class="{ rotate: showCustomPanel }" />
        </span>
      </template>
    </div>
    <div class="custom-panel" v-if="showCustomPanel">
      <a-range-picker
        v-model:value="range"
        format="YYYY-MM-DD"
        :disabledDate="disabledDate"
        :allowClear="false"
        @change="handleDateChange"
      />
    </div>
  </div>
</template>
<script setup>
import dayjs from 'dayjs'
import { ref, reactive, onMounted, watch } from 'vue'
import { DownOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat-monitor.components.date-filter')

const props = defineProps({
  datekey: {
    type: String,
    default: '2'
  }
})

const emit = defineEmits(['dateChange'])

const create_time_type = ref(2)
const showCustomPanel = ref(false)
const range = ref([])

const tagDateArr = reactive([
  {
    label: 'label_today',
    key: 2,
    start_date: dayjs().startOf('day'),
    end_date: dayjs().endOf('day')
  },
  {
    label: 'label_yesterday',
    key: 3,
    start_date: dayjs().subtract(1, 'day').startOf('day'),
    end_date: dayjs().startOf('day').subtract(1, 'millisecond')
  },
  {
    label: 'label_this_week',
    key: 6,
    start_date: dayjs().startOf('week'),
    end_date: dayjs().endOf('day').subtract(1, 'millisecond')
  },
  {
    label: 'label_last_3_months',
    key: 7,
    start_date: dayjs().subtract(3, 'month').startOf('day'),
    end_date: dayjs().endOf('day').subtract(1, 'millisecond')
  },
  {
    label: 'label_custom',
    key: 5,
    start_date: dayjs().startOf('day'),
    end_date: dayjs().endOf('day')
  }
])

const disabledDate = (current) => {
  return current && current.valueOf() > dayjs().endOf('day').valueOf()
}

const emitRange = (start, end) => {
  emit('dateChange', {
    start_time: start.startOf('day').format('YYYYMMDD'),
    end_time: end.endOf('day').format('YYYYMMDD')
  })
}

const handleTagChange = (key) => {
  create_time_type.value = key
  const dateItem = tagDateArr.find((i) => i.key === key)
  if (!dateItem) return

  if (key === 5) {
    range.value = [dateItem.start_date, dateItem.end_date]
    showCustomPanel.value = !showCustomPanel.value
    create_time_type.value = 5
    return
  }

  range.value = [dateItem.start_date, dateItem.end_date]
  emitRange(dateItem.start_date, dateItem.end_date)
}

const handleDateChange = (dates) => {
  if (!dates || dates.length < 2) return
  const [start, end] = dates
  create_time_type.value = 5
  tagDateArr.forEach((i) => {
    if (i.key === 5) {
      i.start_date = start
      i.end_date =end
    }
  })
  emitRange(start, end)
}

onMounted(() => {
  const initKey = parseInt(props.datekey.split('-')[0]) || 2
  handleTagChange(initKey)
})

watch(
  () => props.datekey,
  (newVal) => {
    const key = parseInt(String(newVal).split('-')[0]) || 2
    handleTagChange(key)
  }
)
</script>
<style lang="less" scoped>
.zm-date-quick {
  .btns {
    display: flex;
    align-items: center;
    gap: 16px;
    .btn {
      font-size: 14px;
      color: #595959;
      cursor: pointer;
      transition: color 0.2s ease;
      display: inline-flex;
      align-items: center;
      .down-icon {
        margin-left: 4px;
        font-size: 12px;
        color: #8c8c8c;
      }
      &:hover {
        color: #2475fc;
      }
      &.active {
        color: #2475fc;
      }
    }
  }
  .custom-dropdown {
    padding: 8px;
    background: #fff;
    border-radius: 6px;
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
  }
  .custom-panel {
    margin-top: 8px;
    padding: 8px;
    background: #fff;
    border-radius: 6px;
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
  }
  .down-icon {
    transition: transform 0.2s ease;
    &.rotate {
      transform: rotate(180deg);
    }
  }
}
</style>
