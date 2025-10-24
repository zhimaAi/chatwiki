<template>
  <div class="main-content-block">
    <a-tabs class="tab-wrapper" v-model:activeKey="activeKey" @change="handleChangeTab">
      <a-tab-pane :key="1" tab="默认知识库"></a-tab-pane>
      <a-tab-pane :key="2" tab="关联知识库"></a-tab-pane>
    </a-tabs>
    <div class="body-content-box" style="padding-right: 16px" v-if="activeKey == 1">
      <DefaultLibrary />
    </div>
    <div class="body-content-box" v-else>
      <cu-scroll>
        <RelatedLibrary />
      </cu-scroll>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import DefaultLibrary from './default-library.vue'
import RelatedLibrary from './related-library.vue'

const activeLocalKey = '/robot/config/library-config/activeKey'
const activeKey = ref(+localStorage.getItem(activeLocalKey) || 1)
const handleChangeTab = () => {
  localStorage.setItem(activeLocalKey, activeKey.value)
}
</script>

<style lang="less" scoped>
.main-content-block {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  ::v-deep(.ant-tabs-nav-wrap) {
    padding-left: 24px;
  }
  .tab-wrapper ::v-deep(.ant-tabs-nav) {
    margin-bottom: 0;
  }
}
.body-content-box {
  flex: 1;
  overflow: hidden;
}
</style>
