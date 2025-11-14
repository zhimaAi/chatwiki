<style lang="less" scoped>

</style>

<template>
  <span>{{ text }}</span>
</template>

<script setup>
import { ref, inject, watch } from 'vue'

const props = defineProps({
  value: {
    type: Array,
    default: () => [],
  },
})

const getNode = inject('getNode')

const variableOptionsSelect = ref([])
const text = ref('')

function getOptions() {
  let list = getNode().getAllParentVariable()

  variableOptionsSelect.value = handleOptions(list)
}

function handleOptions(options) {
  options.forEach((item) => {
    if (item.typ == 'node') {
      if (item.node_type == 1) {
        item.value = 'global'
      } else {
        item.value = item.node_id
      }
    } else {
      item.value = item.key
    }

    if (item.children && item.children.length > 0) {
      item.children = handleOptions(item.children)
    }
  })

  return options
}

// 递归查找函数，根据value在树结构中查找对应的label
const findLabelByValue = (options, value) => {
  for (const item of options) {
    if (item.value === value) {
      return item.label
    }
    if (item.children && item.children.length > 0) {
      const foundLabel = findLabelByValue(item.children, value)
      if (foundLabel) {
        return foundLabel
      }
    }
  }
  return null
}

const setText = () => {
  if (!props.value || props.value.length === 0) {
    return ''
  }

  getOptions()

  let labelArr = []

  // 根据valueArr中的每个值，递归查找对应的label
  props.value.forEach((value) => {
    const label = findLabelByValue(variableOptionsSelect.value, value)
    if (label) {
      labelArr.push(label)
    }
  })

  // 返回拼接后的label字符串，用"."连接
  text.value = labelArr.join('/')
}

const refresh = () => {
  setText()
}

watch(() => props.value, () => {
 setText()
}, {
  immediate: true,
})

defineExpose({
  refresh,
})
</script>
