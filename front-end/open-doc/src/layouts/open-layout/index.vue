<template>
  <div class="open-doc-layout">
    <PreviewAlert v-if="previewKey" />
    <LayoutHeader v-if="!hideHeader" />
    <div class="open-doc-container">
      <div class="open-doc-left" v-if="!isEditPage">
        <LayoutSideNav :activeKey="activeKey" :previewKey="previewKey" />
      </div>

      <div class="open-layout-body">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import LayoutHeader from './compoents/LayoutHeader.vue'
import LayoutSideNav from './compoents/LayoutSideNav.vue'
import PreviewAlert from './compoents/preview-alert.vue'
import { useOpenDocStore } from '@/stores/open-doc'

const store = useOpenDocStore()

const route = useRoute()

const previewKey = computed(() => {
  return store.previewKey
})

const activeKey = computed(() => {
  return route.params.id
})

const hideHeader = computed(() => {
  return route.meta.hideHeader
})

const isEditPage = computed(() => { return store.isEditPage })
</script>

<style lang="less" scoped>
.open-doc-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  overflow: hidden;
}
.open-doc-container {
  flex: 1;
  display: flex;
  height: 100vh;
  overflow: hidden;
  flex-direction: row;
}
.open-layout-body {
  flex: 1;
  overflow-y: auto;
}
@media (min-width: 992px) {
  .open-doc-left {
    width: 256px;
    height: 100%;
    border-right: 1px solid #d9d9d9;
  }
}
</style>
