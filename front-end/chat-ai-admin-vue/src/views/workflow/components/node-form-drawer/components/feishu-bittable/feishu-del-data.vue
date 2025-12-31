<template>
  <div class="options-item is-required">
    <div class="options-item-tit">
      <div class="option-label">record_id</div>
      <div class="option-type">string</div>
    </div>
    <div>
      <AtInput
        type="textarea"
        inputStyle="height: 64px;"
        :options="variableOptions"
        :defaultSelectedList="state.tags"
        :defaultValue="state.value"
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
    <div class="desc">待删除的记录ID；示例值："reCWNXZPQv"</div>
  </div>
</template>

<script setup>
import {ref, reactive, toRaw} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";

const props = defineProps({
  variableOptions: {
    type: Array,
  }
})
const emit = defineEmits(['update', 'updateVar'])
const state = reactive({
  tags: [],
  value: '',
})

function init(nodeParams=null) {
  if (nodeParams && nodeParams?.plugin?.params?.arguments?.record_id) {
    state.value = nodeParams.plugin.params.arguments.record_id
    state.tags = nodeParams.plugin.params.arguments.record_tags || []
  }
}

function changeValue(val, tags) {
  state.value = val
  state.tags = tags
  emit('update', {
    record_id: val,
    record_tags: tags,
  })
}

defineExpose({
  init,
})
</script>

<style scoped lang="less">
@import "common";
</style>
