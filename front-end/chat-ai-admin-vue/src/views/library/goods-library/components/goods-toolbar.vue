<template>
  <div class="goods-toolbar">
    <div class="toolbar-search">
      <a-input
        v-model:value="keywordValue"
        :placeholder="t('toolbar.search_placeholder')"
        allowClear
        @change="handleInputChange"
        @pressEnter="handleSearch"
      >
        <template #suffix>
          <SearchOutlined />
        </template>
      </a-input>
    </div>
    <div class="toolbar-actions">
      <a-button @click="handleImport">
        <template #icon>
          <UploadOutlined />
        </template>
        {{ t('toolbar.import') }}
      </a-button>
      <a-button @click="handleExport">
        <template #icon>
          <DownloadOutlined />
        </template>
        {{ t('toolbar.export') }}
      </a-button>
      <a-button type="primary" @click="handleCreate">
        <template #icon>
          <PlusOutlined />
        </template>
        {{ t('toolbar.create') }}
      </a-button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { DownloadOutlined, PlusOutlined, SearchOutlined, UploadOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.goods-library.index')

const props = defineProps({
  keyword: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:keyword', 'search', 'import', 'export', 'create'])

const keywordValue = ref(props.keyword)

watch(
  () => props.keyword,
  (value) => {
    keywordValue.value = value || ''
  }
)

const handleSearch = () => {
  emit('update:keyword', keywordValue.value)
  emit('search', keywordValue.value)
}

const handleInputChange = (event) => {
  const value = event?.target?.value ?? ''
  keywordValue.value = value
  emit('update:keyword', value)
  if (!value) {
    emit('search', '')
  }
}

const handleImport = () => {
  emit('import')
}

const handleExport = () => {
  emit('export')
}

const handleCreate = () => {
  emit('create')
}
</script>

<style lang="less" scoped>
.goods-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;

  .toolbar-search {
    width: 240px;
  }

  .toolbar-actions {
    display: flex;
    align-items: center;
    gap: 8px;
    flex-shrink: 0;
  }
}
</style>
