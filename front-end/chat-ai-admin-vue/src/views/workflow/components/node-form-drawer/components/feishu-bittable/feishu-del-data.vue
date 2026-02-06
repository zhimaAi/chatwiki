<template>
  <div class="options-item is-required">
    <div class="options-item-tit">
      <div class="option-label">{{ t('label_record_id') }}</div>
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
    <div class="desc">{{ t('desc_record_id') }}</div>
  </div>
</template>

<script setup>
import {ref, reactive, toRaw} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.feishu-bittable.feishu-del-data')

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
