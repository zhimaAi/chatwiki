<template>
  <div class="open-doc-header">
    <div class="library-title">{{ libraryInfo.library_title }}</div>
    <div class="action-box">
      <span class="action-btn" @click="toggleSidebar()">
        <img class="action-btn-icon" src="@/assets/img/open_sidebar.svg" alt="" />
      </span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useOpenDocStore } from '@/stores/open-doc'

const openDocStore = useOpenDocStore()

const libraryInfo = computed(() => {
  return openDocStore.libraryInfo
})

const toggleSidebar = () => {
  openDocStore.toggleSidebar()
}
</script>

<style lang="less" scoped>
.open-doc-header {
  display: flex;
  align-items: center;
  padding: 2px 16px;
  z-index: 100;
  background-color: #fff;
  box-shadow:
    0 1px 2px -2px rgb(0 0 0 / 16%),
    0 2px 3px 0 rgb(0 0 0 / 12%),
    0 3px 4px 0px rgb(0 0 0 / 9%);

  .library-title {
    flex: 1;
    line-height: 22px;
    font-size: 16px;
    font-weight: 400;
    color: #666666;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
.open-doc-header .action-box {
  padding-left: 16px;
}

.open-doc-header .action-box {
  display: flex;
  align-items: center;
  padding-left: 8px;
}
.open-doc-header .action-box .action-btn {
  width: 40px;
  height: 40px;
  line-height: 40px;
  margin-left: 8px;
  border-radius: 8px;
  text-align: center;
  transition: all 0.2s;
}
.open-doc-header .action-box .action-btn:hover {
  cursor: pointer;
  background-color: #f0f2f5;
}

@media (min-width: 992px) {
  .open-doc-header {
    display: none;
  }
}
</style>
