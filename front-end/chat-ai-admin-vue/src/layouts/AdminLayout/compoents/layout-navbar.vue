<style lang="less" scoped>
.navbar {
  display: flex;
  justify-content: center;
  align-items: center;

  .nav-menu {
    position: relative;
    margin: 0 24px;
    cursor: pointer;

    .nav-menu-name {
      display: inline-block;
      line-height: 22px;
      padding: 4px 0 6px 0;
      font-size: 14px;
      font-weight: 400;
      color: #595959;
    }
  }

  .nav-menu.active {
    &::after {
      content: '';
      display: block;
      position: absolute;
      bottom: 0;
      left: 50%;
      margin-left: -20px;
      width: 40px;
      height: 2px;
      background-color: #2475fc;
    }

    .nav-menu-name {
      font-size: 14px;
      font-weight: 600;
      color: #2475fc;
    }
  }
}
</style>

<template>
  <div class="navbar-wrapper">
    <div class="navbar">
      <div
        class="nav-menu"
        :class="{ active: item.key === rootPath || item.key === activeMenu }"
        v-for="item in items"
        :key="item.key"
      >
        <router-link :to="item.path" class="nav-menu-name">{{ item.title }}</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
// front-end\chat-ai-admin-vue\src\utils\permission.js
const roure = useRoute()

const activeMenu = computed(() => {
  return roure.meta.activeMenu || ''
})

const rootPath = computed(() => {
  return roure.path.split('/')[1]
})

const getActiveMenu = () => {}

const items = ref([
  {
    key: 'robot',
    label: 'robot',
    title: '机器人管理',
    path: '/robot/list'
  },
  {
    key: 'library',
    label: 'library',
    title: '知识库管理',
    path: '/library/list'
  },
  {
    key: 'user',
    label: 'user',
    title: '系统管理',
    path: '/user/model'
  }
])

watch(
  () => roure.path,
  () => {
    getActiveMenu()
  },
  {
    immediate: true
  }
)
</script>
