<style lang="less" scoped>

</style>

<template>
  <span>{{ text }}</span>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { ref, watch } from 'vue'
import { useRobotStore } from '@/stores/modules/robot'

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

const robotStore = useRobotStore()
const { modelList } = storeToRefs(robotStore)

const text = ref('')

const setText = () => {
  if (!props.useModel) {
    text.value = ''
    return
  }

  for (const item of modelList.value) {
      if (item.children) {
        for (const child of item.children) {
          if (props.useModel == child.name && props.modelConfigId == child.id) {
            text.value = child.show_model_name || child.deployment_name ||  child.name;
            return // Exit loops once found
          }
        }
      }
    }
}

watch(() => [props.useModel, props.modelConfigId], () => {
  setText()
}, {
  immediate: true
})
</script>
