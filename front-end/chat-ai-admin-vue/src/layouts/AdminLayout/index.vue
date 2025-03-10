<style lang="less" scoped>
.admin-layout {
  height: 100vh;
  position: relative;
  display: flex;
  flex-flow: column nowrap;
  background: linear-gradient(180deg, #e5efff 0%, #f0f2f5 34.43%);

  .layout-header {
    position: relative;
    display: flex;
    align-items: center;
    height: var(--layout-header-height);
    padding: 0 16px;

    .header-left {
      position: absolute;
      left: 16px;
      top: 0;
      bottom: 0;
      display: flex;
      align-items: center;
    }
    .header-right {
      position: absolute;
      right: 16px;
      top: 0;
      bottom: 0;
      display: flex;
      align-items: center;
      gap: 22px;
    }

    .header-body {
      flex: 1;
    }
  }

  .layout-breadcrumb-wrapper {
    padding: 0;
  }

  .layout-body {
    flex: 1;
    padding: 0 16px 0 16px;
    display: flex;
    flex-flow: column nowrap;
    overflow: hidden;
  }

  .page-wrapper {
    display: flex;
    flex-flow: column nowrap;
    flex: 1;
    overflow: hidden;
    border-radius: 2px;

    .page-container {
      flex: 1;
      padding: 0 22px 0 22px;
      overflow-x: hidden;
      overflow-y: auto;
    }
  }
}
</style>

<template>
  <div class="admin-layout">
    <div class="layout-header-wrapper">
      <div class="layout-header">
        <div class="header-left">
          <LayoutLogo />
        </div>
        <div class="header-body">
          <Layoutnavbar />
        </div>
        <div class="header-right">
          <div class="item-box">
            <LocaleDropdown />
          </div>
          <div class="item-box">
            <UserDropdown />
          </div>
        </div>
      </div>
    </div>

    <div class="layout-body">
      <!-- 自定义页面样式 -->
      <template v-if="isCustomPage">
        <router-view></router-view>
      </template>
      <template v-else>
        <div class="page-wrapper" :style="{ 'background-color': bgColor }">
          <div class="layout-breadcrumb-wrapper" v-if="!hideTitle">
            <LayoutBreadcrumb :items="breadcrumb" :title="pageTitle" />
          </div>
          <div class="page-container" :style="{ ...pageStyle }">
            <router-view></router-view>
          </div>
        </div>
      </template>
    </div>
    <div class="layout-footer-wrapper">
      <LayoutFooter />
    </div>
    <ResetPassword v-if="showResetModal"></ResetPassword>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import LayoutLogo from './compoents/layout-logo.vue'
import Layoutnavbar from './compoents/layout-navbar.vue'
import LayoutBreadcrumb from './compoents/layout-breadcrumb.vue'
import LayoutFooter from './compoents/layout-footer.vue'
import UserDropdown from './compoents/user-dropdown.vue'
import LocaleDropdown from './compoents/locale-dropdown.vue'
import ResetPassword from './compoents/reset-password.vue'
import { useUserStore } from '@/stores/modules/user'

const route = useRoute()
const userStore = useUserStore()


const isCustomPage = computed(() => {
  return route.meta.isCustomPage || false
})

const pageStyle = computed(() => {
  return route.meta.pageStyle || {}
})

const bgColor = computed(() => {
  return route.meta.bgColor || '#ffffff'
})

const breadcrumb = computed(() => {
  return route.meta.breadcrumb || []
})
const hideTitle = computed(() => {
  return route.meta.hideTitle || false
})

const pageTitle = computed(() => {
  return route.meta.title || false
})

const showResetModal = computed(() => {
  return userStore.userInfo?.d_pass == 1 && userStore.isShowResetPassModal
})
</script>
