<template>
  <a-modal
    :title="t('title_content')"
    width="880px"
    v-model:open="visible"
    @ok="save"
  >
    <a-alert type="info" class="zm-alert-info">
      <template #message>
        {{ t('msg_rich_text_content_intro') }}
      </template>
    </a-alert>
    <AtInput
      style="margin-top: 16px;"
      type="textarea"
      inputStyle="height: 420px;"
      :options="variableOptions"
      :defaultSelectedList="tags || []"
      :defaultValue="value"
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
  </a-modal>
</template>

<script setup>
import {ref, shallowRef, onBeforeUnmount, toRaw, watch} from 'vue'
import '@wangeditor/editor/dist/css/style.css'
import {Editor, Toolbar} from '@wangeditor/editor-for-vue'
import {Editor as SlateEditor, Range, Transforms, Text} from 'slate'
import CascadePanel from "@/views/workflow/components/at-input/cascade-panel.vue";
import {runPlugin} from "@/api/plugins/index.js";
import {getBase64, getTreeOptions} from "@/utils/index.js";
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.official-draft.components.draft-newspic-cont')


const emit = defineEmits(['change', 'updateVar'])
const props = defineProps({
  appData: {
    type: Object,
    default: () => ({})
  },
  variableOptions: {
    type: Array,
  },
  defaultSelectedList: {
    type: Array,
    default: () => [],
  },
})

const visible = ref(false)
const value = ref('')
const tags = ref([])

function show(val = '', _tags = []) {
  value.value = val
  tags.value = _tags
  visible.value = true
}

function changeValue(val, _tags ) {
  value.value = val
  tags.value = _tags
}

function save() {
  emit('change', toRaw(value.value), toRaw(tags.value))
  visible.value = false
}

defineExpose({
  show,
})
</script>

<style scoped lang="less">

</style>
