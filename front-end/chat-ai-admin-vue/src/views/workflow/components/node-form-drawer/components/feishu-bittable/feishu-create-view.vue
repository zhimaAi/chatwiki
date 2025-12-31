<template>
  <div>
    <div class="options-item is-required">
      <div class="options-item-tit">
        <div class="option-label">视图名称</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <AtInput
          type="textarea"
          inputStyle="height: 33px;"
          :options="variableOptions"
          :defaultSelectedList="state?.tag_map?.view_name || []"
          :defaultValue="state.view_name"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(val, tags) => changeValue('view_name', val, tags)"
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
      <div class="desc">视图名称。名称不能包含特殊字符；不可以包含 [ ] 等特殊字符；1 字符 ～ 100 字符</div>
    </div>
    <div class="options-item">
      <div class="options-item-tit">
        <div class="option-label">视图类型</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <a-select
          v-model:value="state.view_type"
          :options="VIEW_SELECT_OPTIONS"
          placeholder="选择视图"
          style="width: 100%;"
        />
      </div>
      <div class="desc">视图类型，默认为表格视图。</div>
    </div>
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
  view_name: '',
  view_type: 'grid',
  tag_map: {},
})

const VIEW_SELECT_OPTIONS = [
  { label: '表格视图', value: 'grid' },
  { label: '看板视图', value: 'kanban' },
  { label: '画册视图', value: 'gallery' },
  { label: '甘特视图', value: 'gantt' },
  { label: '表单视图', value: 'form' }
]

function init(nodeParams=null) {
  if (nodeParams && nodeParams?.plugin?.params?.arguments) {
    state.name = nodeParams.plugin.params.arguments?.name || ""
    state.view_name = nodeParams.plugin.params.arguments?.view_name || ""
    state.tag_map = nodeParams.plugin.params.arguments?.tag_map || {}
  }
}

function changeValue(field, val, tags) {
  state[field] = val
  state.tag_map[field] = tags
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
