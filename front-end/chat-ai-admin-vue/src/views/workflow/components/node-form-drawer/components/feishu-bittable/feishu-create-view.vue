<template>
  <div>
    <div class="options-item is-required">
      <div class="options-item-tit">
        <div class="option-label">{{ t('label_view_name') }}</div>
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
      <div class="desc">{{ t('desc_view_name') }}</div>
    </div>
    <div class="options-item">
      <div class="options-item-tit">
        <div class="option-label">{{ t('label_view_type') }}</div>
        <div class="option-type">string</div>
      </div>
      <div class="min-input">
        <a-select
          v-model:value="state.view_type"
          :options="VIEW_SELECT_OPTIONS"
          :placeholder="t('ph_select_view')"
          style="width: 100%;"
        />
      </div>
      <div class="desc">{{ t('desc_view_type') }}</div>
    </div>
  </div>
</template>

<script setup>
import {reactive} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.feishu-bittable.feishu-create-view')

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
  { label: t('opt_grid_view'), value: 'grid' },
  { label: t('opt_kanban_view'), value: 'kanban' },
  { label: t('opt_gallery_view'), value: 'gallery' },
  { label: t('opt_gantt_view'), value: 'gantt' },
  { label: t('opt_form_view'), value: 'form' }
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
