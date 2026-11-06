<template>
  <div class="date-panel">
    <div class="panel-head">
      <span class="panel-title">{{ t('date_range') }}</span>
      <button type="button" class="panel-close" @click="$emit('close')">
        <CloseOutlined />
      </button>
    </div>

    <div class="date-input">
      <span :class="{ placeholder: !start }">{{ start ? start.format('YYYY-MM-DD') : t('start_date') }}</span>
      <span class="date-dash">—</span>
      <span :class="{ placeholder: !end }">{{ end ? end.format('YYYY-MM-DD') : t('end_date') }}</span>
    </div>

    <div class="date-helper" :class="{ invisible: start && end }">{{ helper }}</div>

    <div class="date-months">
      <div v-for="(month, index) in months" :key="month.valueOf()" class="date-month">
        <div class="month-head">
          <button
            v-if="index === 0"
            type="button"
            class="month-nav"
            :aria-label="t('prev_month')"
            @click="shiftMonth(-1)"
          >
            <LeftOutlined />
          </button>
          <span v-else class="month-nav-placeholder"></span>
          <span class="month-title">{{ monthTitle(month) }}</span>
          <button
            v-if="index === 1"
            type="button"
            class="month-nav"
            :aria-label="t('next_month')"
            @click="shiftMonth(1)"
          >
            <RightOutlined />
          </button>
          <span v-else class="month-nav-placeholder"></span>
        </div>
        <div class="weekdays">
          <span v-for="label in weekdayLabels" :key="label">{{ label }}</span>
        </div>
        <div class="days">
          <button
            v-for="cell in monthCells(month)"
            :key="cell.key"
            type="button"
            class="day-cell"
            :class="cellClasses(cell)"
            :disabled="!cell.inMonth"
            @click="pick(cell.date)"
            @mouseenter="onHover(cell)"
          >
            <span v-if="cell.inMonth" class="day-inner">{{ cell.date.date() }}</span>
          </button>
        </div>
      </div>
    </div>

    <div class="date-footer">
      <a-button class="date-footer-button" @click="clear">{{ t('clear') }}</a-button>
      <a-button
        class="date-footer-button"
        type="primary"
        :disabled="!start || !end"
        @click="confirm"
      >
        {{ t('confirm') }}
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import dayjs from 'dayjs'
import { CloseOutlined, LeftOutlined, RightOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const props = defineProps({
  value: {
    type: Array,
    default: null
  },
  active: {
    type: Boolean,
    default: false
  }
})
const emit = defineEmits(['close', 'confirm', 'clear'])
const { t, locale } = useI18n('views.customer-service.index')

const isZh = computed(() => String(locale.value || '').toLowerCase().startsWith('zh'))
const weekdayLabels = computed(() => (
  isZh.value ? ['一', '二', '三', '四', '五', '六', '日'] : ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
))
const englishMonths = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

const start = ref(props.value?.[0] || null)
const end = ref(props.value?.[1] || null)
const hoverDate = ref(null)
const viewMonth = ref((start.value || dayjs()).startOf('month'))

const resetDraft = () => {
  start.value = props.value?.[0] || null
  end.value = props.value?.[1] || null
  hoverDate.value = null
  viewMonth.value = (start.value || dayjs()).startOf('month')
}

watch(
  () => props.active,
  (active) => {
    if (active) resetDraft()
  }
)

const months = computed(() => [viewMonth.value, viewMonth.value.add(1, 'month')])
const helper = computed(() => (start.value && !end.value ? t('select_end_date') : t('select_start_date')))
const rangeStart = computed(() => start.value)
const rangeEnd = computed(() => {
  if (end.value) return end.value
  if (start.value && hoverDate.value && hoverDate.value.isAfter(start.value, 'day')) return hoverDate.value
  return null
})

const monthTitle = (month) => (
  isZh.value ? month.format('YYYY年M月') : `${englishMonths[month.month()]} ${month.year()}`
)

const monthCells = (month) => {
  const first = month.startOf('month')
  const offset = (first.day() + 6) % 7
  const daysInMonth = month.daysInMonth()
  const cells = []
  for (let i = 0; i < 42; i += 1) {
    const dayNumber = i - offset + 1
    const inMonth = dayNumber >= 1 && dayNumber <= daysInMonth
    cells.push({
      key: `${month.valueOf()}-${i}`,
      inMonth,
      date: inMonth ? first.date(dayNumber) : null
    })
  }
  return cells
}

const sameDay = (a, b) => Boolean(a && b && a.isSame(b, 'day'))
const cellClasses = (cell) => {
  if (!cell.inMonth) return {}
  const date = cell.date
  const startDay = rangeStart.value
  const endDay = rangeEnd.value
  const isStart = sameDay(date, start.value)
  const isEnd = sameDay(date, end.value)
  const inRange = startDay && endDay && date.isAfter(startDay, 'day') && date.isBefore(endDay, 'day')
  return {
    'range-start': isStart,
    'range-end': isEnd,
    'in-range': inRange,
    selected: isStart || isEnd
  }
}

const pick = (date) => {
  if (!date) return
  hoverDate.value = null
  if (!start.value || (start.value && end.value)) {
    start.value = date
    end.value = null
    return
  }
  if (date.isBefore(start.value, 'day')) {
    end.value = start.value
    start.value = date
    return
  }
  end.value = date
}

const onHover = (cell) => {
  if (cell.inMonth) hoverDate.value = cell.date
}

const shiftMonth = (step) => {
  viewMonth.value = viewMonth.value.add(step, 'month')
}

const clear = () => {
  start.value = null
  end.value = null
  hoverDate.value = null
  emit('clear')
}

const confirm = () => {
  if (!start.value || !end.value) return
  emit('confirm', [start.value, end.value])
}
</script>

<style lang="less" scoped>
.date-panel {
  display: flex;
  flex-direction: column;
  width: 568px;
  padding: 20px;
  font-family: 'PingFang SC', 'Noto Sans SC', sans-serif;
}
.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.panel-title {
  color: #242424;
  font-size: 16px;
  font-weight: 500;
}
.panel-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  padding: 0;
  border: 0;
  color: #4e5969;
  background: transparent;
  font-size: 14px;
  cursor: pointer;
}
.date-input {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  padding: 0 16px;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
  color: #242424;
  font-size: 14px;
}
.date-input .placeholder {
  color: #8a8f99;
}
.date-dash {
  color: #c9cdd4;
}
.date-helper {
  margin: 12px 0;
  color: #888;
  font-size: 12px;
  line-height: 18px;
}
.date-helper.invisible {
  visibility: hidden;
}
.date-months {
  display: flex;
  gap: 16px;
}
.date-month {
  flex: 1;
  min-width: 0;
}
.month-head {
  display: flex;
  align-items: center;
  height: 40px;
}
.month-title {
  flex: 1;
  color: #242424;
  font-size: 16px;
  font-weight: 500;
  text-align: center;
}
.month-nav {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  border: 0;
  border-radius: 6px;
  color: #4e5969;
  background: transparent;
  font-size: 14px;
  cursor: pointer;

  &:hover {
    background: #f5f6f8;
  }
}
.month-nav-placeholder {
  width: 32px;
  flex-shrink: 0;
}
.weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  align-items: center;
  height: 40px;
  color: #888;
  font-size: 14px;
  text-align: center;
}
.days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
}
.day-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 34px;
  padding: 0;
  border: 0;
  background: transparent;
  cursor: pointer;

  &:disabled {
    cursor: default;
  }

  &.in-range {
    background: #eef3ff;
  }

  &.range-start {
    border-radius: 50% 0 0 50%;
    background: #eef3ff;
  }

  &.range-end {
    border-radius: 0 50% 50% 0;
    background: #eef3ff;
  }

  &.range-start.range-end {
    border-radius: 50%;
    background: transparent;
  }
}
.day-inner {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  color: #242424;
  font-size: 14px;
  line-height: 32px;
  transition: background 150ms;

  .day-cell:hover & {
    background: #f0f1f3;
  }

  .day-cell.selected & {
    color: #fff;
    background: #1456f0;
  }
}
.date-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}
.date-footer-button {
  height: 34px;
  padding: 0 18px;
  border-radius: 8px;
  font-size: 14px;
}
</style>
