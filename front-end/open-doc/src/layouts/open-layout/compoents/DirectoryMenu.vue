<style lang="less" scoped></style>

<template>
  <div class="directory-menus">
    <DirectoryMenuItem
      v-for="item in props.items"
      :token="token"
      :activeKey="activeKey"
      :previewKey="previewKey"
      :docKey="item.doc_key"
      :icon="item.doc_icon"
      :title="item.title"
      :key="item.doc_key"
      :isDir="item.is_dir"
      :item="item"
      :level="0"
      :children="item.children"
      :forceExpanded="props.forceExpanded"
      :forceCollapsed="props.forceCollapsed"
    />
  </div>
</template>

<script setup>
import DirectoryMenuItem from './DirectoryMenuItem.vue'

const props = defineProps({
  items: {
    type: Array,
    default: () => [],
  },
  forceExpanded: {
    type: Boolean,
    default: false,
  },
  forceCollapsed: {
    type: Boolean,
    default: false,
  },
  activeKey: {
    type: String,
    default: '',
  },
  previewKey: {
    type: String,
    default: '',
  },
  token: {
    type: String,
    default: '',
  },
  iconTemplateConfig: {
    type: Object,
    default: () => ({}),
  }
})
</script>
