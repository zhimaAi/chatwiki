<template>
  <div>
    <div class="options-item is-required">
      <div class="options-item-tit">
        <div class="option-label">多维表格的名称</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <AtInput
          type="textarea"
          inputStyle="height: 64px;"
          :options="variableOptions"
          :defaultSelectedList="state?.tag_map?.name || []"
          :defaultValue="state.name"
          ref="atInputRef"
          @open="emit('updateVar')"
          @change="(val, tags) => changeValue('name', val, tags)"
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
      <div class="desc">更新多维表格的名称</div>
    </div>
    <div class="options-item is-required">
      <div class="flex-between">
        <div class="options-item-tit ">
          <div class="option-label">开启高级权限</div>
          <div class="option-type">string</div>
        </div>
        <div>
          <a-switch :checked="state.is_advanced > 0" @change="advancedChange"></a-switch>
        </div>
      </div>
      <div class="min-input">
        <a-radio-group v-model:value="state.is_advanced" :disabled="state.is_advanced < 1" @change="update">
          <a-radio value="1">开启</a-radio>
          <a-radio value="2">关闭</a-radio>
        </a-radio-group>
      </div>
      <div class="desc">多维表格是否已开启高级权限，开关关闭则不更新设置</div>
    </div>
  </div>
</template>

<script setup>
import {ref, reactive, watch} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {runPlugin} from "@/api/plugins/index.js";

const props = defineProps({
  variableOptions: {
    type: Array,
  }
})
const emit = defineEmits(['update', 'updateVar'])

const atInputRef = ref(null)
const state = reactive({
  name: "",
  is_advanced: '0',
  tag_map: {},
})

function init(nodeParams = null) {
  if (nodeParams && nodeParams?.plugin?.params?.arguments) {
    Object.assign(state, nodeParams.plugin.params.arguments)
  }
}

function changeValue(field, val, tags) {
  state[field] = val
  state.tag_map[field] = tags
  update()
}

function advancedChange(val) {
  state.is_advanced = val ? '1' : '0'
  update()
}

function update() {
  emit('update', JSON.parse(JSON.stringify(state)))
}

defineExpose({
  init,
  update
})
</script>

<style scoped lang="less">
@import "common";

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
