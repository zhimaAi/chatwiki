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
        <span>{{ t('label_reply_content') }}</span>
        <span class="tip">{{ t('msg_empty_field_tip') }}</span>
      </div>

      <div class="setting-box">
        <FieldListSelect
          ref="fieldRef"
          :showEmptyFieldRow="!tableId"
          :list="dataItems"
          :fields="fields"
          :msgtype="tableId"
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
import { useI18n } from '@/hooks/web/useI18n'
import FieldListSelect from "@/views/workflow/components/data-table/reply-select/index.vue";

const { t } = useI18n('views.workflow.components.node-form-drawer.components.official-send-message.reply-data')

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
const fieldRef = ref(null)

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

function validateAll () {
  return fieldRef.value?.validateAll?.()
}

defineExpose({
  init,
  validateAll
})
</script>

