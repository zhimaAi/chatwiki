<template>
  <div class="_container">
    <div class="header">
      <div class="main-tab-box">
        <MainTab ref="tabRef" @change="tabChange"/>
      </div>
    </div>
    <div class="content">
      <RemotePlugins v-if="active == 3" @installReport="updateTab"/>
      <InstallPlugins v-else @tabChange="showRemote" @installReport="updateTab"/>
    </div>
  </div>
</template>

<script setup>
import { onMounted , ref} from 'vue';
import {useRoute} from 'vue-router'
import RemotePlugins from "./components/remotePlugins.vue";
import InstallPlugins from "./components/installPlugins.vue";
import MainTab from "@/views/explore/components/main-tab.vue";

const route = useRoute()
const tabRef = ref(null)
const active = ref(localStorage.getItem('zm:explore:active') || '2')

onMounted(() => {
  if (route.query.active > 1) {
    tabRef.value.change(route.query.active)
  }
})

function tabChange(val){
  if (val > 1) active.value = val
}

function showRemote(val) {
  tabRef.value.change(val)
}

function updateTab() {
  tabRef.value.update()
}
</script>

<style scoped lang="less">
._container {
  height: 100%;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.header {
  .main-tab-box {
    display: flex;
    margin-bottom: 8px;
  }
  .tabs-box {
    display: flex;
    align-items: center;
    gap: 8px;

    .tab-item {
      padding: 5px 16px;
      border-radius: 6px;
      background: #EDEFF2;
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      cursor: pointer;

      &.active {
        color: #2475fc;
        background: #D6E6FF;
      }
    }
  }
}

.content {
  flex: 1;
  overflow-y: auto;
  margin-top: 16px;
}
</style>
