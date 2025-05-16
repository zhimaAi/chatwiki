<template>
  <div class="details-library-page">
    <div class="between-content-box">
      <div class="layout-left">
        <div class="library-name-box">
          <img class="avatar" :src="avatar || LIBRARY_NORMAL_AVATAR" alt="" />
          <div class="name">
            {{ library_name }}
          </div>
        </div>
        <div class="left-menu-wrapper">
          <LeftMenus></LeftMenus>
        </div>
      </div>

      <div
        class="right-content-box"
        :class="['page-' + rotue.name]"
        :style="{ overflow: rotue.name == 'knowledgeDocument' ? 'hidden' : '' }"
      >
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import LeftMenus from './components/left-menus.vue'
import { LIBRARY_NORMAL_AVATAR } from '@/constants/index'
import { useLibraryStore } from '@/stores/modules/library'

const libraryStore = useLibraryStore()
const { getLibraryInfo } = libraryStore

const library_name = computed(() => libraryStore.library_name)
const avatar = computed(() => libraryStore.avatar)

const rotue = useRoute()
const query = rotue.query

const getInfo = () => {
  getLibraryInfo(query.id)
}

onMounted(() => {
  // getInfo()
})
</script>

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: #fff;
  border-radius: 2px;
}
.layout-left {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-right: 1px solid rgba(5, 5, 5, 0.06);
  .library-name-box {
    display: flex;
    align-items: center;
    padding: 24px 24px 16px 24px;
    .avatar {
      width: 32px;
      height: 32px;
      margin-right: 8px;
      border-radius: 4px;
    }
    .name {
      line-height: 22px;
      font-size: 14px;
      font-weight: 600;
      color: #262626;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      max-width: 200px;
    }
  }
}
.between-content-box {
  display: flex;
  flex: 1;
  overflow: hidden;
  .left-menu-wrapper {
    flex: 1;
    overflow: hidden;
    padding: 0 4px;
  }
  .right-content-box {
    flex: 1;
    padding: 24px 10px 0 24px;
  }
}
</style>
<style>
.page-knowledgeGraph{
  padding: 0 !important;
}
</style>
