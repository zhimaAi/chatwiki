<template>
  <div class="options-item is-required">
    <div class="options-item-tit">
      <div class="option-label">多维表格名称</div>
      <div class="option-type">string</div>
    </div>
    <div class="min-input">
      <AtInput
        type="textarea"
        inputStyle="height: 33px;"
        :options="variableOptions"
        :defaultSelectedList="state.tag_map.name || []"
        :defaultValue="state.name"
        ref="atInputRef"
        @open="emit('updateVar')"
        @change="changeValue"
        placeholder="请输入内容，键入“/”可以插入变量"
      >
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
      </AtInput>
    </div>
    <div class="desc">多维表格 App 名称。最长为 255 个字符。 示例值："一篇新的多维表格"</div>
  </div>
</template>

<script setup>
import {reactive} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";

const props = defineProps({
  variableOptions: {
    type: Array,
  }
})
const emit = defineEmits(['update', 'updateVar'])
const state = reactive({
  tag_map: {},
  name: '',
})

function init(nodeParams=null) {
  if (nodeParams && nodeParams?.plugin?.params?.arguments?.name) {
    state.name = nodeParams.plugin.params.arguments.name
    state.tag_map = nodeParams.plugin.params.arguments.tag_map || {}
  }
}

function changeValue(val, tags) {
  state.name = val
  state.tag_map.name = tags
  update()
}

function update() {
  emit('update', JSON.parse(JSON.stringify(state)))
}

defineExpose({
  init,
  update,
})
</script>

<style scoped lang="less">
@import "common";
</style>
