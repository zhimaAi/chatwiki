<template>
  <div class="page-container-box">
    <a-tabs v-model:activeKey="activeKey" @change="handleChangeTabs">
      <a-tab-pane :key="1" tab="知识库配置"></a-tab-pane>
      <a-tab-pane :key="2" tab="角色权限"></a-tab-pane>
    </a-tabs>
    <div class="content-box">
      <KnowledgeConfig v-if="activeKey == 1" />
      <rolePermission v-if="activeKey == 2" />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import KnowledgeConfig from '../knowledge-config.vue'
import rolePermission from './role-permission.vue'
import { useRoute, useRouter } from 'vue-router'
const route = useRoute()
const router = useRouter()
const query = route.query
const activeKey = ref(+query.activeKey || 1)
const handleChangeTabs = () => {
  let queryParmas = {
    ...query
  }
  if (activeKey.value > 1) {
    queryParmas.activeKey = activeKey.value
  }else{
    delete queryParmas.activeKey
  }
  router.push({
    query: {
      ...queryParmas
    }
  })
}
</script>

<style lang="less" scoped>
.page-container-box {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  ::v-deep(.ant-tabs-nav) {
    margin-bottom: 0;
    .ant-tabs-nav-wrap {
      padding-left: 24px;
    }
  }
  .content-box {
    flex: 1;
    overflow: hidden;
  }
}
</style>
