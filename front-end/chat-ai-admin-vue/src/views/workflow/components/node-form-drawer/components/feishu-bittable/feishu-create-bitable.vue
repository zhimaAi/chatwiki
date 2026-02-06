<template>
  <div class="options-item is-required">
    <div class="options-item-tit">
      <div class="option-label">{{ t('label_bitable_name') }}</div>
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
        :placeholder="t('ph_input_content')"
      >
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
      </AtInput>
    </div>
    <div class="desc">{{ t('desc_bitable_app_name') }}</div>
  </div>
</template>

<script setup>
import {reactive} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.feishu-bittable.feishu-create-bitable')

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
