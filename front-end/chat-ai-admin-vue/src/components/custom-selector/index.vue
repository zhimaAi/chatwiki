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
        <div v-for="(group, gi) in processedModelList" :key="gi">
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
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { duplicateRemoval, removeRepeat } from '@/utils/index'
import { getModelConfigOption } from '@/api/model/index'

const props = defineProps({
  modelValue: [String, Number],
  modelType: {
    type: String,
    validator: (value) => {
      return ['TEXT EMBEDDING', 'RERAN', 'LLM'].includes(value)
    },
    required: true
  },
  isOffline: {
    type: Boolean,
    default: false
  },
  placeholder: {
    type: String,
    default: '请选择模型'
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
  modelConfigId: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'change', 'loaded'])


const showOptions = ref(false)

const placeholder = ref(props.placeholder)

const modelList = ref([])
const modelDefine = ['azure', 'ollama', 'xinference', 'openaiAgent', 'doubao']
const modelShowModelName = ['doubao']
// 处理原始数据格式
const processedModelList = computed(() => {
  return modelList.value.map(group => ({
    groupLabel: group.name,
    groupIcon: group.icon,
    children: group.children.map(child => ({
      icon: child.icon,
      use_model_name: child.use_model_name,
      show_model_name: child.show_model_name,
      value: child.show_model_name ? child.show_model_name : modelDefine.includes(child.model_define) && child.deployment_name ? child.deployment_name : child.name,
      rawData: child // 保留原始数据
    }))
  }))
})

const selectedItem = computed(() => {
  const { modelConfigId, modelValue } = props
  let newItem = {}
  if (!modelConfigId && !modelValue) return
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
  return processedModelList.value.reduce((acc, group) => {
    return acc.concat(group.children || [])
  }, [])
})

function uniqueArr(arr, arr1, key) {
  const keyVals = new Set(arr.map((item) => item.model_define))
  arr1.filter((obj) => {
    let val = obj[key]
    if (keyVals.has(val)) {
      arr.filter((obj1) => {
        if (obj1.model_define == val) {
          obj1.children = removeRepeat(obj1.children, obj.children)
          return false
        }
      })
    }
  })
  return arr
}

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

  const getModelList = () => {
    getModelConfigOption({
      model_type: props.modelType,
      is_offline: props.isOffline // 0 1 区分线上线下
    }).then((res) => {
      let list = res.data || []

      let newList = list.map((item) => {
        let { model_define, deployment_name, id } = item.model_config
        let key = `${model_define}-${deployment_name || ''}-${id}`

        let subModelList = []
        let children = []

        if (props.modelType == 'TEXT EMBEDDING') {
          subModelList = item.model_info.vector_model_list
        } else if (props.modelType == 'RERANK') {
          subModelList = item.model_info.rerank_model_list
        } else if (props.modelType == 'LLM') {
          subModelList = item.model_info.llm_model_list
        }

        for (let i = 0; i < subModelList.length; i++) {
          let ele = subModelList[i]
          let label = ele
          let value = id + '-' + ele // 这里的值使用父级id加子模型名称的方式来避免重复（因为有多个叫’默认的模型‘）
          // 部分模型显示的模型名称要特殊处理
          if (
            modelDefine.indexOf(item.model_config.model_define) > -1 &&
            item.model_config.deployment_name
          ) {
            label = item.model_config.deployment_name
          }
          let show_model_name = ''
          if(modelShowModelName.includes(item.model_config.model_define)){
            show_model_name = item.model_config.show_model_name
          }

          children.push({
            key: id + '_' + ele,
            label: label,
            name: ele,
            deployment_name: deployment_name,
            model_config_id: id,
            model_define: model_define,
            value: value,
            id: item.model_config.id,
            show_model_name,
            icon: item.model_info.model_icon_url, // 添加图标字段
            use_model_name: item.model_info.model_name // 添加系统名称字段
          })
        }

        return {
          key: key,
          model_define: model_define,
          children: children,
          model_config_id: id,
          deployment_name: deployment_name,
          name: item.model_info.model_name,
          icon: item.model_info.model_icon_url
        }
      })

      // 如果modelList存在两个相同model_define情况就合并到一个对象的children中去
      newList = uniqueArr(duplicateRemoval(newList, 'model_define'), newList, 'model_define')

      if (newList.length === 0) {
        placeholder.value = '无可用模型'
      }

      modelList.value = [...newList]

      emit('loaded', JSON.parse(JSON.stringify(newList)))
    })
  }

  watch(
    () => props.isOffline,
    () => {
      getModelList()
    }
  )
  
  onMounted(() => {
    getModelList()
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
    border-radius: 6px;
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
    border-radius: 6px;
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