<template>
    <div class="custom-select" @click.stop="toggleOptions">
      <div class="selected-item">
        <template v-if="selectedItem">
          <img 
            v-if="selectedItem[iconKey]" 
            :src="selectedItem[iconKey]" 
            class="model-icon"
          />
          <span>
            {{ selectedItem[labelKey] }}
            <span v-if="showValue && selectedItem[valueKey]">
              {{ selectedItem[valueKey] }}
            </span>
          </span>
        </template>
        <span v-else class="placeholder">{{ placeholder }}</span>
      </div>
  
      <div class="custom-options" v-show="showOptions" @wheel.prevent="handleWheel">
        <div v-for="(group, gi) in options" :key="gi">
          <div v-if="group.groupLabel" class="option-group">
            <img 
              v-if="group[groupIconKey]" 
              :src="group[groupIconKey]" 
              class="model-icon"
            />
            {{ group[groupLabelKey] }}
          </div>
          <div 
            v-for="(item, ii) in group.children"
            :key="ii"
            class="option-item"
            @click.stop="handleSelect(item)"
          >
            <!-- <img 
              v-if="item[iconKey]" 
              :src="item[iconKey]" 
              class="model-icon"
            />
            <span>
              {{ item[labelKey] }}
              <span v-if="showValue && item[valueKey]">({{ item[valueKey] }})</span>
            </span> -->
            <span>{{ item[valueKey] }}</span>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
  
  const props = defineProps({
    modelValue: [String, Number],
    options: {
      type: Array,
      default: () => []
    },
    placeholder: {
      type: String,
      default: '请选择'
    },
    // 字段配置
    groupLabelKey: {
      type: String,
      default: 'groupLabel'
    },
    groupIconKey: {
      type: String,
      default: 'groupIcon'
    },
    iconKey: {
      type: String,
      default: 'icon'
    },
    labelKey: {
      type: String,
      default: 'label'
    },
    valueKey: {
      type: String,
      default: 'value'
    },
    showValue: {
      type: Boolean,
      default: true
    },
    modelDefine: {
      type: Array,
      default: () => []
    },
    modelConfigId: {
      type: String,
      default: ''
    }
  })

  const emit = defineEmits(['update:modelValue', 'change'])

  const showOptions = ref(false)

  const selectedItem = computed(() => {
    const { modelConfigId, modelDefine, modelValue } = props
    let newItem = {}
    flattenOptions.value.find(item => {
      const { rawData } = item
      if (modelDefine.includes(rawData.model_define) && rawData.deployment_name) {
        if (rawData.deployment_name === modelValue || rawData.id === modelConfigId) {
          newItem = item
        }
      } else {
        if (rawData.name === modelValue && rawData.id === modelConfigId) {
          newItem = item
        }
      }
    })
    return newItem
  })
  
  // 平铺选项用于查找
  const flattenOptions = computed(() => {
    return props.options.reduce((acc, group) => {
      return acc.concat(group.children || [])
    }, [])
  })
  
  // 新增滚动控制逻辑
const scrollTop = ref(0)
const isScrolling = ref(false)

// 鼠标滚轮处理
const handleWheel = (event) => {
  const element = event.currentTarget
  const delta = event.deltaY
  
  // 计算可滚动范围
  const maxScroll = element.scrollHeight - element.clientHeight
  const willScrollPastTop = scrollTop.value + delta <= 0
  const willScrollPastBottom = scrollTop.value + delta >= maxScroll

  // 阻止默认滚动传播
  event.stopPropagation()
  
  if (!willScrollPastTop && !willScrollPastBottom) {
    // 优先处理内部滚动
    element.scrollTop += delta
    event.preventDefault()
  } else if (
    (scrollTop.value === 0 && delta < 0) || 
    (scrollTop.value === maxScroll && delta > 0)
  ) {
    // 到达边界时允许外部滚动
    isScrolling.value = false
  } else {
    // 处理边界情况
    element.scrollTop = Math.max(0, Math.min(maxScroll, scrollTop.value + delta))
    event.preventDefault()
  }

  // 更新滚动位置
  scrollTop.value = element.scrollTop
}

  const toggleOptions = () => {
    showOptions.value = !showOptions.value
  }
  
  const handleSelect = (item) => {
    emit('update:modelValue', item[props.valueKey])
    emit('change', item)
    showOptions.value = false
  }
  
  // 点击外部关闭
  const clickOutsideHandler = (e) => {
    if (!e.target.closest('.custom-select')) {
      showOptions.value = false
    }
  }
  
  onMounted(() => {
    document.addEventListener('click', clickOutsideHandler)
  })
  
  onBeforeUnmount(() => {
    document.removeEventListener('click', clickOutsideHandler)
  })
  </script>
  
<style lang="less" scoped>
  .custom-select {
    position: relative;
    width: 100%;
    border: 1px solid #d9d9d9;
    border-radius: 4px;
    padding: 5px 11px;
    cursor: pointer;
    background: white;
    min-height: 32px;
  }
  
  .selected-item {
    display: flex;
    align-items: center;
    gap: 8px;
    min-height: 22px;
  }
  
  .placeholder {
    color: rgba(0, 0, 0, 0.25);
  }
  
  .model-icon {
    height: 18px;
    object-fit: contain;
  }
  
  .custom-options {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    background: white;
    border: 1px solid #d9d9d9;
    border-radius: 4px;
    z-index: 1000;
    max-height: 300px;
    overflow-y: auto;
    overscroll-behavior: contain; /* 防止滚动链 */
    box-shadow: 0 3px 6px -4px rgba(0, 0, 0, 0.12), 
              0 6px 16px 0 rgba(0, 0, 0, 0.08),
              0 9px 28px 8px rgba(0, 0, 0, 0.05);

    .option-group {
      background: white;
      padding: 8px 12px;
      font-size: 14px;
      color: rgba(0, 0, 0, 0.45);
      display: flex;
      align-items: center;
      gap: 8px;

      .model-icon {
        height: 18px;
      }
    }

    .option-item {
      padding: 8px 12px 8px 38px;
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      transition: all 0.3s;
      font-size: 14px;

      &:hover {
        background: #f5f5f5;
      }

      .model-icon {
        height: 16px;
      }

      span > span {
        color: rgba(0, 0, 0, 0.45);
        font-size: 12px;
        margin-left: 4px;
      }
    }
  }
</style>