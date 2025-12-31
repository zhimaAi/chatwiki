<template>
  <div class="slider-box" :class="{ 'neo4j-status-close-slider': !neo4j_status }" :style="cssVars">
    <div class="form-label-block">
      <div class="label-title">
        权重
        <a-tooltip>
          <template #title>
            混合检索模式，结果中各检索方式的加权权重 权重越高，对应检索结果的采纳数量越多
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
      :range="neo4j_status"
    />
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { useCompanyStore } from '@/stores/modules/company'
const companyStore = useCompanyStore()
const neo4j_status = computed(() => {
  return companyStore.companyInfo?.neo4j_status == 'true'
})
const props = defineProps({
  rrf_weight: {
    type: Object
  }
})
const emit = defineEmits(['update:rrf_weight', 'save'])
const rrf_value = ref(null)

const listItems = computed(() => {
  let { vector, search, graph } = props.rrf_weight
  let list = [
    {
      label: '向量',
      value: formatter(vector),
      color: '#2475fc'
    },
    {
      label: '全文',
      value: formatter(search),
      color: '#03b615'
    }
  ]
  if (neo4j_status.value) {
    list.push({
      label: '图',
      value: formatter(graph),
      color: '#F59A23'
    })
  }
  return list
})

const handleChange = () => {
  let vector = 0
  let search = 0
  let graph = 0
  if (neo4j_status.value) {
    vector = rrf_value.value[0]
    search = rrf_value.value[1] - rrf_value.value[0]
    graph = 100 - rrf_value.value[1]
  } else {
    vector = rrf_value.value
    search = 100 - rrf_value.value
  }
  emit('update:rrf_weight', { vector, search, graph })
  emit('save', { vector, search, graph })
}

watch(
  () => props.rrf_weight,
  (val) => {
    let { vector, search, graph } = val
    if (neo4j_status.value) {
      rrf_value.value = [vector, Math.min(vector + search, 100)]
    } else {
      rrf_value.value = vector
    }
  },
  { immediate: true, deep: true }
)

const backgroundLiner = computed(() => {
  if (!neo4j_status.value) {
    return `linear-gradient(to right, #2475fc ${rrf_value.value}%, #03B615 ${100 - rrf_value.value}%)`
  }
  let value1 = rrf_value.value[0]
  let value2 = 100 - value1
  return `linear-gradient(to right, #2475fc ${value1}%, #F59A23 ${value2}%)`
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
.neo4j-status-close-slider {
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
    gap: 8px;
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
