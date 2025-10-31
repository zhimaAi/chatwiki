<style lang="less" scoped>
.admin-layout {
  height: 100vh;
  position: relative;
  display: flex;
  flex-flow: column nowrap;
  background: #fff;

  .layout-header {
    position: relative;
    display: flex;
    align-items: center;
    height: var(--layout-header-height);
    padding: 0 16px;
    border-bottom: 1px solid #D9D9D9;

    .header-left {
      display: flex;
      align-items: center;
    }
    .header-right {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .header-body {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
    }
  }

  .layout-breadcrumb-wrapper {
    padding: 0;
  }

  .layout-body {
    flex: 1;
    padding: 0;
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
      padding: 0 48px;
      overflow-x: hidden;
      overflow-y: auto;
    }
  }
}
</style>

<template>
  <div class="admin-layout">
    <div class="layout-header-wrapper" v-if="!hideLayoutTopAndBottom">
      <div class="layout-header">
        <div class="header-left">
          <LayoutLogo />
        </div>
        <div class="header-body">
          <Layoutnavbar />
        </div>
        <div class="header-right">
          <!-- <div class="item-box">
            <LocaleDropdown />
          </div> -->
          <div class="item-box">
            <UserDropdown />
          </div>
        </div>
      </div>
    </div>

    <div class="layout-body">
      <!-- 自定义页面样式 -->
      <template v-if="isCustomPage">
        <slot><router-view></router-view></slot>
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
    <div class="layout-footer-wrapper" v-if="!hideLayoutTopAndBottom">
      <LayoutFooter />
    </div>
    <ResetPassword v-if="showResetModal"></ResetPassword>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
import LayoutLogo from './compoents/layout-logo.vue'
import Layoutnavbar from './compoents/layout-navbar.vue'
import LayoutBreadcrumb from './compoents/layout-breadcrumb.vue'
import LayoutFooter from './compoents/layout-footer.vue'
import UserDropdown from './compoents/user-dropdown.vue'
// import LocaleDropdown from './compoents/locale-dropdown.vue'
import ResetPassword from './compoents/reset-password.vue'
import { useGlobalStore } from '@/stores/modules/global'
const globalStore = useGlobalStore()



const route = useRoute()
const userStore = useUserStore()

const hideLayoutTopAndBottom = computed(() => {
  if(route.name != 'robotWorkflow'){
    return false
  }
  return globalStore.hideLayoutTopAndBottom || false
})

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
