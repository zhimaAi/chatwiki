<template>
  <div class="add-data-form">
    <div class="node-box-content">
      <div
        v-for="(item, key, idx) in formState"
        :key="key"
        class="options-item is-required"
      >
        <div class="options-item-tit">
          <div class="flex-between">
            <div class="option-label">{{ item.name || key }}</div>
            <div class="option-type">{{ item.type }}</div>
          </div>
          <div class="btn-hover-wrap" @click="handleOpenFullAtModal(key, idx)">
            <FullscreenOutlined/>
          </div>
        </div>
        <div>
          <AtInput
            type="textarea"
            inputStyle="height: 64px;"
            :options="variableOptions"
            :defaultSelectedList="item.tags"
            :defaultValue="item.value"
            ref="atInputRef"
            @open="emit('updateVar')"
            @change="(text, selectedList) => changeValue(item, text, selectedList)"
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
        <div class="desc">{{item.desc}}</div>
      </div>
    </div>
    <FullAtInput
      :options="variableOptions"
      :defaultSelectedList="fullDefaultTags"
      :defaultValue="fullDefaultValue"
      :placeholder="t('ph_input_content')"
      type="textarea"
      @open="emit('updateVar')"
      @change="(val, tags) => changeValueByFull(val, tags)"
      @ok="handleRefreshAtInput"
      ref="fullAtInputRef"
    />
  </div>
</template>

<script setup>
import {FullscreenOutlined} from '@ant-design/icons-vue';
import {reactive, onMounted, ref} from 'vue';
import { useI18n } from '@/hooks/web/useI18n'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import FullAtInput from "@/views/workflow/components/at-input/full-at-input.vue";
import {jsonDecode} from "@/utils/index.js";
import {getBatchActionParams} from "@/constants/feishu-table.js";

const { t } = useI18n('views.workflow.components.node-form-drawer.components.feishu-bittable.common')

const emit = defineEmits(['update', 'updateVar'])
const props = defineProps({
  actionName: {
    type: String,
  },
  action: {
    type: Object,
  },
  variableOptions: {
    type: Array,
  }
})
const formState = reactive({})
const fullAtInputRef = ref(null)
const atInputRef = ref([])
const activeKey = ref('')
const activeIdx = ref(-1)
const fullDefaultValue = ref('')
const fullDefaultTags = ref([])

onMounted(() => {
  init()
})

function handleOpenFullAtModal(key, idx) {
  activeKey.value = key
  activeIdx.value = idx
  const cur = formState[key] || {}
  fullDefaultValue.value = cur.value || ''
  fullDefaultTags.value = cur.tags || []
  fullAtInputRef.value.show()
}

function init(nodePrams=null) {
  const base = getBatchActionParams(props.action?.params || {})
  for (let field in base) {
    const old = formState[field] || {}
    formState[field] = {
      ...base[field],
      value: old.value ?? base[field].value ?? "",
      tags: old.tags ?? base[field].tags ?? []
    }
  }
  if (nodePrams) {
    let args = nodePrams?.plugin?.params?.arguments || {}
    let tag_map = args?.tag_map || {}
    for (let field in formState) {
      const cur = formState[field] || {}
      const hasUserInput = !!(String(cur.value || "").length) || (Array.isArray(cur.tags) && cur.tags.length > 0)
      if (!hasUserInput) {
        formState[field].value = args[field] ?? cur.value ?? ""
        formState[field].tags = tag_map[field] ?? cur.tags ?? []
      }
    }
  }
}

function changeValue(item, val, tags) {
  item.value = val
  item.tags = tags
  update()
}

function changeValueByFull(val, tags) {
  const key = activeKey.value
  if (!key) return
  const item = formState[key]
  if (!item) return
  item.value = val
  item.tags = tags
  update()
}


const update = () => {
  let data = {
    tag_map: {}
  }
  for (let field in formState) {
    data[field] = formState[field].value
    data.tag_map[field] = formState[field].tags
  }
  emit('update', data)
}

function handleRefreshAtInput() {
  const idx = activeIdx.value
  if (idx != null && idx > -1) {
    const arr = atInputRef.value || []
    const target = Array.isArray(arr) ? arr[idx] : arr
    target && target.refresh && target.refresh()
  }
}

defineExpose({
  init,
})
</script>

<style scoped lang="less">
.options-item-tit {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-hover-wrap {
  cursor: pointer;
}
</style>
