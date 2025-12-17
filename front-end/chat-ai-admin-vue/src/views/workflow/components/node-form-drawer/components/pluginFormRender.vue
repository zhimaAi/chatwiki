<template>
  <div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/input.svg" class="title-icon"/>输入</div>
      </div>
      <template v-for="(item, key) in formState" :key="key">
        <slot v-if="$slots[key]" :state="formState" :name="key" :item="item" :keyName="key"></slot>
        <div v-else :class="['options-item', {'is-required': item.required}]">
          <div class="options-item-tit">
            <div class="option-label">{{ getLabel(item, key) }}</div>
            <div class="option-type">{{ item.type }}</div>
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
              @change="(text, selectedList) => changeValue(item, text, selectedList, key)"
              :placeholder="getPlaceholder(key)"
            >
              <template #option="{ label, payload }">
                <div class="field-list-item">
                  <div class="field-label">{{ label }}</div>
                  <div class="field-type">{{ payload.typ }}</div>
                </div>
              </template>
            </AtInput>
          </div>
          <div class="desc">{{ item.desc }}</div>
          <div v-if="item.__error" class="desc" style="color: #FB363F">{{ item.__error }}</div>
        </div>
      </template>
    </div>
    <div class="node-options">
      <div class="options-title">
        <div><img src="@/assets/img/workflow/output.svg" class="title-icon"/>输出</div>
      </div>
      <div class="options-item">
        <OutputFields :tree-data="outputData"/>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, reactive, watch, inject} from 'vue';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import OutputFields from "@/views/workflow/components/feishu-table/output-fields.vue";
import {pluginOutputToTree} from "@/constants/plugin.js";

const emit = defineEmits(['updateVar'])
const props = defineProps({
  node: {
    type: Object,
    default: () => ({})
  },
  params: {
    type: Object,
  },
  output: {
    type: Object,
  },
  variableOptions: {
    type: Array,
  },
  keepFormatField: {
    type: Array,
    default: () => ([])
  }
})
const setData = inject('setData')

const formState = reactive({})
const outputData = ref([])
const nodeParams = JSON.parse(props.node.node_params)

const fieldRules = {
  url: {
    label: (item, key) => item?.name || key,
    placeholder: '键入/选择参数，或者输入文本，以https:/开头',
    validate: (item) => {
      const hasVar = Array.isArray(item.tags) && item.tags.length > 0
      if (hasVar) { item.__error = null; return true }
      const v = String(item.value || '').trim()
      const ok = v.startsWith('https:/')
      item.__error = ok ? null : '请输入以https:/开头的网址'
      return ok
    }
  }
}

watch(formState, (newVal) => {
  update()
}, {
  deep: true,
})

watch(() => props.params, (newVal) => {
  Object.assign(formState, JSON.parse(JSON.stringify(newVal || '{}')))
  let _args = nodeParams?.plugin?.params?.arguments || {}
  let val
  for (let key in formState) {
    // 1、是否已设置值 2、是否存在默认值
    val = _args[key] || ''
    if (!val && Object.hasOwn(formState[key], 'default')) {
      val = formState[key].default
    }
    if (!props.keepFormatField.includes(key)) val = String(val)
    if (val === 'null') val = null
    formState[key].value = val
    formState[key].tags = _args?.tag_map?.[key] || []
  }
}, {
  deep: true,
  immediate: true
})

watch(() => props.output, (newVal) => {
  outputData.value = pluginOutputToTree(JSON.parse(JSON.stringify(props.output || '{}')))
}, {
  deep: true,
  immediate: true
})



function changeValue(item, text, selectedList, key) {
  item.value = text
  item.tags = selectedList
  validateField(key, item)
  update()
}

function getLabel(item, key) {
  const r = fieldRules[key]
  if (r && r.label) return r.label(item, key)
  return item?.name || key
}

function getPlaceholder(key) {
  const r = fieldRules[key]
  if (r && r.placeholder) return r.placeholder
  return '请输入内容，键入“/”可以插入变量'
}

function validateField(key, item) {
  const r = fieldRules[key]
  if (r && r.validate) return r.validate(item)
  item.__error = null
  return true
}

function update() {
  let nodeParams = JSON.parse(props.node.node_params)
  nodeParams.plugin.output_obj = outputData.value
  let _args = {
    tag_map: {}
  }
  let value
  for (let key in formState) {
    value = formState[key].value
    if(formState[key].type == 'string'){
      value = String(value)
    } else if(['number', 'integer'].includes(formState[key].type)){
      value = Number(value)
    }
    _args[key] = value
    _args.tag_map[key] = formState[key].tags
  }
  Object.assign(nodeParams.plugin.params, {
    arguments: {
      ..._args
    }
  })
  setData({
    ...props.node,
    node_params: JSON.stringify(nodeParams)
  })
}
</script>

<style lang="less">
@import "node-options";
</style>
