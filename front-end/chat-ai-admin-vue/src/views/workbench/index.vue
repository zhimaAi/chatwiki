<template>
    <div class="workbench-layout">
        <LeftSidebar />
        <div class="workbench-layout-container">
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup>
import { onMounted } from 'vue'
import LeftSidebar from './components/left-sidebar.vue'
import { useWorkbenchStore } from '@/stores/modules/workbench'

const workbenchStore = useWorkbenchStore()

onMounted(() => {
    // 初始化时加载工作台配置
    workbenchStore.initConfig()
    // 初始化时加载历史访问列表
    workbenchStore.fetchRobotHistoryVisit()
})
</script>

<style lang="less" scoped>
.workbench-layout {
    display: flex;
    height: 100vh;
    width: 100vw;
    .workbench-layout-container {
        flex: 1;
        height: 100%;
        overflow: auto;
    }
}
</style>