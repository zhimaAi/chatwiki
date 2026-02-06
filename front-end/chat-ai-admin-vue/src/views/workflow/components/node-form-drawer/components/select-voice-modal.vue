<template>
  <a-modal
    v-model:open="visible"
    :title="t('title_select_voice')"
    width="670px"
    @ok="save"
  >
    <a-input v-model:value.trim="keyword" :placeholder="t('ph_search_voice')" allow-clear style="width: 240px;">
      <template #suffix>
        <SearchOutlined/>
      </template>
    </a-input>
    <a-table
      class="mt16"
      :loading="loading"
      :columns="columns"
      :data-source="showList"
      :pagination="false"
      row-key="voice_id"
      :row-selection="{ type: 'radio', selectedRowKeys, onChange }"
      :scroll="{y: 500}"
    >
      <template #bodyCell="{column, record}">
        <template v-if="'description' === column.dataIndex">
          <a-tooltip :title="record.description">
            <div class="zm-line2">{{record.description}}</div>
          </a-tooltip>
        </template>
      </template>
    </a-table>
  </a-modal>
</template>

<script setup>
import {ref, computed, toRaw} from 'vue';
import {getVoiceList} from "@/api/node/index.js";
import {SearchOutlined} from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-form-drawer.components.select-voice-modal')
const emit = defineEmits(['change'])
const props = defineProps({
  modelConfigId: {
    type: [Number, String]
  },
  voiceType: {
    type: String,
    default: 'all'
  }
})
const visible = ref(false)
const loading = ref(false)
const keyword = ref('')
const list = ref([])
const selectedRowKeys = ref([])
const selectedRows = ref([])

const columns = computed(() => [
  {
    title: t('label_index'),
    dataIndex: 'index',
    key: 'index',
    width: 80,
  },
  {
    title: t('label_name'),
    dataIndex: 'voice_name',
    width: 120,
  },
  {
    title: t('label_id'),
    dataIndex: 'voice_id',
    width: 80,
  },
  {
    title: t('label_description'),
    dataIndex: 'description',
    width: 200,
  },
])
const showList = computed(() => {
  if (keyword.value) {
    return list.value.filter(item => {
      return item.voice_name.indexOf(keyword.value) > -1 || item.voice_id == keyword.value
    })
  }
  return list.value
})

function show(keys = []) {
  selectedRowKeys.value = keys
  loadData()
  visible.value = true
}

function loadData() {
  loading.value = true
  getVoiceList({
    model_config_id: props.modelConfigId,
    voice_type: props.voiceType,
  }).then(res => {
    list.value = res?.data || []
    list.value.forEach((item, index) => {
      item.index = index + 1
      item.description = item?.description?.[0] || ''
    })
    selectedRows.value = list.value.filter(i => selectedRowKeys.value.includes(i.voice_id))
  }).finally(() => {
    loading.value = false
  })
}

function onChange(keys, rows) {
  selectedRowKeys.value = keys
  selectedRows.value = rows
}

function save() {
  emit('change', toRaw(selectedRowKeys.value),toRaw(selectedRows.value))
  visible.value = false
}

defineExpose({
  show,
})
</script>

<style scoped lang="less">
.mt16 {
  margin-top: 16px;
}
</style>
