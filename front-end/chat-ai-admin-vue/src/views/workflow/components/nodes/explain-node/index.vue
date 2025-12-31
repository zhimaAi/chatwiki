<template>
  <node-common
    :title="properties.node_name"
    :icon-url="properties.node_icon"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    style="width: 420px"
  >
    <div class="explain-node">
      <a-textarea
        auto-size
        :bordered="false"
        class="explain-textarea"
        v-model:value="formState.content"
        placeholder=""
      />
    </div>
  </node-common>
</template>

<script setup>
import NodeCommon from '../base-node.vue'
import { nextTick, onMounted, inject, watch, reactive } from 'vue'
const resetSize = inject('resetSize')
const props = defineProps({
  properties: {
    type: Object,
    default() {
      return {}
    }
  },
  isSelected: { type: Boolean, default: false },
  isHovered: { type: Boolean, default: false }
})

watch(
  () => props.properties,
  (newVal, oldVal) => {
    const newDataRaw = newVal.node_params || '{}'
    const oldDataRaw = oldVal.node_params || '{}'
    if (newDataRaw != oldDataRaw) {
      reset()
    }
  },
  { deep: true }
)

const formState = reactive({
  content: ''
})

const reset = () => {
  if (!props.properties || !props.properties.node_params) {
    return
  }

  let node_params = JSON.parse(props.properties.node_params)

  Object.assign(formState, node_params)

  nextTick(() => {
    resetSize()
  })
}

onMounted(() => {
  reset()
})
</script>
