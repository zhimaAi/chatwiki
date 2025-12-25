<template>
  <a-modal
    title="正文内容"
    width="880px"
    v-model:open="visible"
    @ok="save"
  >
    <a-alert type="info" class="zm-alert-info">
      <template #message>
        图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS,涉及图片url必须来源
        "上传图文消息内的图片获取URL"接口获取。外部图片url将被过滤。 图片消息则仅支持纯文本和部分特殊功能标签如商品，商品个数不可超过50个
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
      placeholder="请输入内容，键入“/”可以插入变量"
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
