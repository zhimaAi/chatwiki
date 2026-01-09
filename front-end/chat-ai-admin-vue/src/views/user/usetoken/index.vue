<template>
  <div class="user-model-page">
    <!-- <div class="page-title">Token使用</div> -->
    <div class="tabs-box">
      <a-tabs v-model:activeKey="activeKey">
        <a-tab-pane :key="1" :tab="t('views.user.usetoken.tabByApp')"></a-tab-pane>
        <a-tab-pane :key="2" :tab="t('views.user.usetoken.tabByModel')"></a-tab-pane>
        <a-tab-pane v-if="role_type == 1 || role_type == 2" :key="3" :tab="t('views.user.usetoken.tabTokenQuota')"></a-tab-pane>
      </a-tabs>
    </div>

    <div class="list-wrapper">
      <StaticsByApp v-if="activeKey == 1" />
      <UseToken v-if="activeKey == 2"></UseToken>

      <QuotaToken v-if="activeKey == 3"></QuotaToken>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import StaticsByApp from './statics-by-app.vue'
import UseToken from './use-token.vue'
import QuotaToken from './quota-token.vue'
import { usePermissionStore } from '@/stores/modules/permission'

const { t } = useI18n()
const permissionStore = usePermissionStore()
const role_type = computed(() => permissionStore.role_type)

const activeKey = ref(1)
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;

  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 24px 24px 16px;
    background-color: #fff;
    color: #000000;
    font-family: 'PingFang SC';
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }

  .tabs-box {
    background: #fff;
    &::v-deep(.ant-tabs-nav-wrap) {
      padding-left: 24px;
    }
  }

  .list-wrapper {
    background: #fff;
    height: calc(100% - 65px);
    overflow-x: hidden;
    overflow-y: auto;
  }
}
</style>
