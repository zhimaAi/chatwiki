<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;
  .list-wrapper {
    background: #fff;
    height: calc(100% - 47px);
    overflow-x: hidden;
    overflow-y: auto;
  }
}
</style>

<template>
  <div class="user-model-page">
    <PageTabs v-model:value="activeTab" @change="onChangeTab" />
    <div class="list-wrapper">
      <TeamMembers v-if="activeTab == 1"></TeamMembers>
      <RolePermissions v-else></RolePermissions>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { message, Modal } from 'ant-design-vue'
import PageTabs from './components/page-tabs.vue'
import TeamMembers from './team-members.vue'
import RolePermissions from './role-permissions.vue'

const activeTab = ref(1)

const onChangeTab = () => {
  console.log(activeTab.value)
}
</script>
