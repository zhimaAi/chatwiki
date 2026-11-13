<template>
  <div class="slider-box two-weight-slider" :style="cssVars">
    <div class="form-label-block">
      <div class="label-title">
        {{ t('label_weight') }}
        <a-tooltip>
          <template #title>
            {{ t('tooltip_weight') }}
          </template>
          <QuestionCircleOutlined />
        </a-tooltip>
      </div>
      <div class="item-list-box">
        <div
          class="list-item"
          v-for="item in listItems"
          :key="item.label"
          :style="{ color: item.color }"
        >
          <span class="dot" :style="{ background: item.color }"></span>
          <span class="text">{{ item.label }}：{{ item.value }}</span>
        </div>
      </div>
    </div>
    <a-slider
      @change="handleChange"
      v-model:value="rrf_value"
      :tip-formatter="formatter"
    />
  </div>
</template>

<script setup>
import { normalizeRrfWeight } from '@/utils/rrf-weight'
import { ref, watch, computed } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('components.weight-select.index')

const props = defineProps({
  rrf_weight: {
    type: Object
  }
})
const emit = defineEmits(['update:rrf_weight', 'save'])
const rrf_value = ref(null)

const listItems = computed(() => {
  const { vector, search } = normalizeRrfWeight(props.rrf_weight)
  return [
    {
      label: t('label_vector'),
      value: formatter(vector),
      color: '#2475fc'
    },
    {
      label: t('label_fulltext'),
      value: formatter(search),
      color: '#03b615'
    }
  ]
})

const handleChange = () => {
  const vector = rrf_value.value
  const search = 100 - vector
  emit('update:rrf_weight', { vector, search })
  emit('save', { vector, search })
}

watch(
  () => props.rrf_weight,
  (val) => {
    const next = normalizeRrfWeight(val)
    rrf_value.value = next.vector
    if (
      Number(val?.vector) !== next.vector ||
      Number(val?.search) !== next.search
    ) {
      emit('update:rrf_weight', next)
    }
  },
  { immediate: true, deep: true }
)

const backgroundLiner = computed(() => {
  return `linear-gradient(to right, #2475fc ${rrf_value.value}%, #03B615 ${100 - rrf_value.value}%)`
})

// 计算属性返回 CSS 变量对象
const cssVars = computed(() => ({
  '--background-liner': backgroundLiner.value
}))

function formatter(value) {
  if (value <= 0) {
    return 0
  }
  if (value >= 100) {
    return 1
  }
  return (value / 100).toFixed(2)
}
</script>

<style lang="less" scoped>
.two-weight-slider {
  &::v-deep(.ant-slider) {
    .ant-slider-track {
      background: #2475fc;
    }
  }
}
.slider-box {
  &::v-deep(.ant-slider) {
    .ant-slider-rail {
      background: var(--background-liner);
    }
    .ant-slider-track-1 {
      background: #03b615;
    }
  }
}
.form-label-block {
  display: flex;
  align-items: center;
  justify-content: space-between;
  .label-title {
    display: flex;
    align-items: center;
    gap: 2px;
    margin-right: 8px;
  }
  .item-list-box {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 13px;
    .list-item {
      display: flex;
      align-items: center;
      gap: 4px;
      .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
      }
    }
  }
}
</style>
