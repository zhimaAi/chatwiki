<style lang="less" scoped>
.add-data-form {
  .node-box-content {
    margin-top: 16px;
    overflow: hidden;
    border-radius: 6px;
    background: #f2f4f7;
  }

  .setting-label {
    line-height: 22px;
    margin-bottom: 8px;
    font-size: 14px;
    color: #262626;
    .tip {
      color: #8c8c8c;
    }
  }

  .setting-box {

  }
}
</style>

<template>
  <div class="add-data-form">
    <div class="node-box-content">
      <div class="setting-label">
        <span>插入数据</span>
        <span class="tip">（不填写字段值则为空）</span>
      </div>

      <div class="setting-box">
        <FieldListSelect
          :showEmptyFieldRow="!tableId"
          :list="dataItems"
          :fields="fields"
          :showAdd="false"
          :showDelete="false"
          @change="onChangeFields"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, toRaw } from 'vue'
import FieldListSelect from "@/views/workflow/components/feishu-table/field-selector/index.vue";

const emit = defineEmits(['update'])
const props = defineProps({
  tableId: {
    type: [String, Number]
  },
  fields: {
    type: Array,
    default: () => ([])
  },
})

const dataItems = ref([])

function init(nodeParams=null) {
  dataItems.value = []
  if (Array.isArray(nodeParams?.plugin?.params?.arguments?.fields) && nodeParams.plugin.params.arguments.fields.length) {
    nodeParams.plugin.params.arguments.fields.map(item => {
      item.value = item.value.toString()
    })
    dataItems.value = nodeParams.plugin.params.arguments.fields
  } else {
    props.fields.map(item => {
      dataItems.value.push({
        value: '',
        atTags: [],
        ...item
      })
    })
  }
}
const onChangeFields = (selectedRows) => {
  dataItems.value = selectedRows
  update()
}

const update = () => {
  emit('update', {
    fields: toRaw(dataItems.value)
  })
}

defineExpose({
  init
})
</script>

