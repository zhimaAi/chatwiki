<template>
    <div class="recent-list">
        <div v-if="!workbenchStore.recentRobots || workbenchStore.recentRobots.length === 0" class="empty-tip">
            {{ t('emptyTip') }}
        </div>
        <div v-else v-for="robot in workbenchStore.recentRobots" :key="robot.robot_key" class="recent-item" :class="{ active: activeRobotId === robot.robot_id }" @click="handleRecentClick(robot)">
            <img :src="robot.robot_avatar" class="recent-avatar" alt="">
            <span class="recent-text">{{ robot.robot_name }}</span>
        </div>
    </div>
</template>

<script setup>
import { onMounted, onActivated, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useWorkbenchStore } from '@/stores/modules/workbench'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workbench.components.recent-list')
const router = useRouter()
const route = useRoute()
const workbenchStore = useWorkbenchStore()

/**
 * 当前选中的机器人 ID（仅在工作台聊天页面且匹配选中模式时有效）
 */
const activeRobotId = computed(() => {
    return route.name === 'workbenchChat' ? workbenchStore.selectedMode : null
})

/**
 * 点击最近使用项目

 * @param {Object} robot - 机器人对象
 */
const handleRecentClick = async (robot) => {
    // 设置选中模式为最近使用项（使用 id 而不是 robot_id）
    workbenchStore.selectRecent(robot.robot_key, robot.robot_id)

    // 记录使用机器人（会自动刷新列表）
    if (robot.robot_id) {
        await workbenchStore.recordRobotVisit(robot.robot_id)
    }

    // 导航到聊天页面
    if(router.path != '/workbench/chat'){
        router.push('/workbench/chat')
    }
}

// 加载历史访问列表
const loadRecentList = () => {
    workbenchStore.fetchRobotHistoryVisit()
}

// 组件挂载时加载
onMounted(() => {
    loadRecentList()
})

// 组件激活时刷新（支持 keep-alive）
onActivated(() => {
    loadRecentList()
})
</script>

<style lang="less" scoped>
.recent-list {
    .empty-tip {
        text-align: center;
        padding: 20px;
        font-size: 14px;
        color: #999;
    }

    .recent-item {
        display: flex;
        align-items: center;
        flex-wrap: nowrap;
        overflow: hidden;
        width: 100%;
        line-height: 22px;
        padding: 8px 9px;
        margin-bottom: 2px;
        border-radius: 6px;
        color: #595959;
        cursor: pointer;
        transition: background-color 0.2s;

        &:hover {
            background-color: #E4E6EB;
        }

        &.active {
            color: #595959;
            background-color: #D8DDE6;
        }

        .recent-avatar {
            width: 18px;
            height: 18px;
            border-radius: 4px;
            margin-right: 8px;
            object-fit: cover;
        }

        .recent-icon {
            width: 18px;
            height: 18px;
            font-size: 12px;
            margin-right: 8px;
            display: flex;
            justify-content: center;
        }

        .recent-text {
            font-size: 14px;
        }
    }
}
</style>
