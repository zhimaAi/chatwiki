<style lang="less" scoped></style>

<template>
  <span>{{ text }}</span>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { ref, watch } from 'vue'
import { useModelStore } from '@/stores/modules/model'
import { getModelOptionsList } from '@/components/model-select/index.js'

const props = defineProps({
  useModel: {
    type: [String, Number],
    default: ''
  },
  modelConfigId: {
    type: [String, Number],
    default: ''
  }
})

const modelStore = useModelStore()
const { allModelList } = storeToRefs(modelStore)

let modelList = getModelOptionsList(allModelList.value).newList

const text = ref('')

const setText = () => {
  if (!props.useModel) {
    text.value = ''
    return
  }

  for (const item of modelList) {
    if (item.children) {
      for (const child of item.children) {
        if (props.useModel == child.name && props.modelConfigId == item.id) {
          text.value = child.show_model_name || child.deployment_name || child.name
          return // Exit loops once found
        }
      }
    }
  }
}

watch(
  () => [props.useModel, props.modelConfigId],
  () => {
    setText()
  },
  {
    immediate: true
  }
)
</script>
