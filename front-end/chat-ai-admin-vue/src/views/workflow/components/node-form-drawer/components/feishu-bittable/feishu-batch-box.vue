<template>
  <div class="add-data-form">
    <div class="node-box-content">
      <div
        v-for="(item, key) in formState"
        :key="key"
        class="options-item is-required"
      >
        <div class="options-item-tit">
          <div class="option-label">{{ item.name || key }}</div>
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
            @change="(text, selectedList) => changeValue(item, text, selectedList)"
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
        <div class="desc">{{item.desc}}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {reactive, onMounted} from 'vue';
import AtInput from "@/views/workflow/components/at-input/at-input.vue";
import {jsonDecode} from "@/utils/index.js";
import {getBatchActionParams} from "@/constants/feishu-table.js";

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

onMounted(() => {
  init()
})

function init(nodePrams=null) {
  Object.assign(formState, getBatchActionParams(props.action?.params || {}))
  if (nodePrams) {
    let args = nodePrams?.plugin?.params?.arguments || {}
    let tag_map = args?.tag_map || {}
    for (let field in formState) {
      formState[field].value = args[field] || ""
      formState[field].tags = tag_map[field] || []
    }
  }
}

function changeValue(item, val, tags) {
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
defineExpose({
  init,
})
</script>

<style scoped>

</style>
