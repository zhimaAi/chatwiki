<template>
    <div class="workbench-team-tools">
        <!-- <div class="page-title">{{ t('title_team_apps') }}</div> -->
        <!-- 顶部导航标签 -->
        <div class="top-nav">
            <div
                class="nav-item"
                :class="{ active: currentGroupId === '' }"
                @click="handleGroupClick('')"
            >
                {{ t('nav_all') }}
            </div>
            <div
                v-for="group in groupList"
                :key="group.id"
                class="nav-item"
                :class="{ active: currentGroupId === group.id }"
                @click="handleGroupClick(group.id)"
            >
                {{ group.group_name }}
            </div>
            <!-- <div class="search-box">
                <input
                    v-model="searchKeyword"
                    type="text"
                    :placeholder="t('ph_search_tools')"
                    @input="handleSearch"
                >
                <SearchOutlined />
            </div> -->
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
            {{ t('msg_loading') }}
        </div>

        <!-- 空状态 -->
        <div v-else-if="filteredRobots.length === 0" class="empty-state">
            {{ t('msg_no_data') }}
        </div>

        <!-- 工具网格布局 -->
        <div v-else class="tools-grid">
            <div
                v-for="robot in filteredRobots"
                :key="robot.id"
                class="tool-card"
                :class="{ active: robot.top_id > 0}"
                @click="handleToolClick(robot)"
            >   
                <cu-tooltip :title="robot.top_id > 0 ? t('tooltip_cancel_top') : t('tooltip_top')" :disabled="tooltipDisabledMap.get(robot.id)" v-if="role_type == 1">
                    <span class="handle-top" @click.stop="handleTop(robot)"></span>
                </cu-tooltip>

                <div class="tool-header">
                    <img
                        v-if="robot.robot_avatar"
                        :src="robot.robot_avatar"
                        class="tool-icon"
                        alt=""
                    >
                    <div
                        v-else
                        class="tool-icon"
                        :style="{ backgroundColor: getRandomColor(robot.id) }"
                    ></div>
                    <div>
                        <div class="tool-name">{{ robot.robot_name }}</div>
                        <div class="tool-category">{{ getApplicationType(robot.application_type) }}</div>
                    </div>
                </div>
                <div class="tool-info">
                    <div class="tool-description">
                        {{ robot.robot_intro || t('msg_no_description') }}
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
// import { SearchOutlined } from '@ant-design/icons-vue'
import { getRobotGroupList, getWorkbenchTeamRobotList, topWorkbenchRobot } from '@/api/workbench'
import CuTooltip from '@/components/cu-tooltip/index.vue'
import { usePermissionStore } from '@/stores/modules/permission'

import { useWorkbenchStore } from '@/stores/modules/workbench'
import { useI18n } from '@/hooks/web/useI18n'

const router = useRouter()
const workbenchStore = useWorkbenchStore()
const permissionStore = usePermissionStore()
const { t } = useI18n('views.workbench.team-tools.index')

// 响应式数据
const groupList = ref([])
const robotList = ref([])
const currentGroupId = ref('')
// const searchKeyword = ref('')
const loading = ref(false)

// tooltip 禁用状态 Map（key: robot.id, value: boolean）
const tooltipDisabledMap = ref(new Map())

const role_type = computed(() => {
    return permissionStore.role_type
})

// 计算属性：根据分组过滤后的机器人列表
const filteredRobots = computed(() => {
    const list = currentGroupId.value === ''
        ? robotList.value
        : robotList.value.filter(robot => robot.group_id === currentGroupId.value)

    return [...list].sort((a, b) => {
        const aTopId = a?.top_id * 1
        const bTopId = b?.top_id * 1

        // 如果都有 top_id，按 top_id 降序排列
        if (aTopId != null && bTopId != null) {
            return bTopId - aTopId
        }

        // 如果只有 a 有 top_id，a 排在前面
        if (aTopId != null) {
            return -1
        }

        // 如果只有 b 有 top_id，b 排在前面
        if (bTopId != null) {
            return 1
        }

        // 都没有 top_id，保持原顺序
        return 0
    })
})


// 获取随机颜色（用于默认图标）
const colors = ['#FFD93D', '#E84855', '#C06C84', '#FF6B6B', '#FF8CC6', '#7367F0', '#F87979', '#4DABF7', '#F03E3E']
const getRandomColor = (id) => {
    const index = id ? id.toString().charCodeAt(0) % colors.length : 0
    return colors[index]
}

// 获取应用类型名称
const getApplicationType = (type) => {
    const types = {
        '0': t('type_chat_bot'),
        '1': t('type_workflow')
    }
    return types[type] || t('type_unknown')
}

// 获取分组列表
const fetchGroupList = async () => {
    try {
        const res = await getRobotGroupList()
        if (res && res.data) {
            groupList.value = res.data
        }
    } catch (error) {
        console.error('获取分组列表失败:', error)
    }
}

// 获取机器人列表
const fetchRobotList = async () => {
    try {
        const res = await getWorkbenchTeamRobotList()
        if (res && res.data) {
            robotList.value = res.data
        }
    } catch (error) {
        console.error('获取机器人列表失败:', error)
    }
}

// 处理分组点击
const handleGroupClick = (groupId) => {
    currentGroupId.value = groupId
}

// 处理搜索（预留）
// const handleSearch = () => {
//     // 搜索逻辑已通过计算属性自动处理
// }

// 处理置顶点击
const handleTop = async (robot) => {
    if (!robot?.id) {
        console.warn('机器人数据不完整:', robot)
        return
    }

    // 禁用该卡片的 tooltip
    tooltipDisabledMap.value.set(robot.id, true)

    try {
        await topWorkbenchRobot({ robot_id: robot.id })
        // robot.is_top = robot.is_top == 1 ? 0 : 1;

        fetchRobotList()
    } catch (error) {
        console.error('置顶失败:', error)
    }

    // 300ms 后恢复 tooltip
    setTimeout(() => {
        tooltipDisabledMap.value.set(robot.id, false)
    }, 350)


}

// 处理工具卡片点击
const handleToolClick = async (robot) => {
    if (!robot.id || !robot.robot_key) {
        console.warn('机器人数据不完整:', robot)
        return
    }

    // 记录使用机器人（生成最近使用记录）
    await workbenchStore.recordRobotVisit(robot.id)

    // 设置选中模式为最近使用项
    workbenchStore.selectRecent(robot.robot_key, robot.id, true)

    // 跳转到聊天页面
    router.push('/workbench/chat')
}


// 页面加载时获取数据
onMounted(async () => {
    loading.value = true
    await Promise.all([fetchGroupList(), fetchRobotList()])
    loading.value = false
})
</script>

<style lang="less" scoped>
.workbench-team-tools {
    padding: 24px;
    background-color: #f5f7fa;
    min-height: 100vh;

    .page-title {
        line-height: 28px;
        margin-bottom: 24px;
        font-size: 20px;
        font-weight: 600;
        color: #2475fc;
    }

    .top-nav {
        display: flex;
        align-items: center;
        margin-bottom: 16px;
        gap: 8px;
        flex-wrap: wrap;

        .nav-item {
            line-height: 22px;
            padding: 5px 16px;
            border-radius: 6px;
            font-size: 14px;
            cursor: pointer;
            color: #595959;
            transition: all 0.3s;
        }

        .nav-item:hover {
            background-color: #E4E6EB;
        }

        .nav-item.active {
            background-color: #fff;
            color: #1890ff;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
        }

        .search-box {
            margin-left: auto;
            position: relative;
            display: flex;
            align-items: center;

            input {
                padding: 6px 12px;
                padding-right: 32px;
                border: 1px solid #d9d9d9;
                border-radius: 4px;
                width: 200px;
                font-size: 14px;
            }

            .anticon {
                position: absolute;
                right: 10px;
                color: #999;
                font-size: 14px;
            }
        }
    }

    .loading-state,
    .empty-state {
        text-align: center;
        padding: 60px 20px;
        font-size: 14px;
        color: #999;
    }

    .tools-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr);
        gap: 20px;

        .tool-card {
            position: relative;
            background-color: #fff;
            border-radius: 8px;
            padding: 16px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
            transition: all 0.3s;
            cursor: pointer;
            overflow: hidden;
            &:hover{
                transform: translateY(-2px);
                box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
            }

            &.active,
            &:hover {
                .handle-top {
                    opacity: 1;
                }
            }

            &.active{
                 .handle-top{
                    background-image: url('../../../assets//img/workbench/handle_top_active.png');
                 }
            }

            .handle-top{
                position: absolute;
                width: 36px;
                height: 36px;
                top: 0;
                right: 0;
                opacity: 0;
                background: url('../../../assets//img/workbench/handle_top.png') 0 0 no-repeat;
                background-size: cover;
                transition: all 0.2s;

                &:hover,
                &.active{
                    background-image: url('../../../assets//img/workbench/handle_top_active.png');
                }
            }

            .tool-header {
                display: flex;
                align-items: center;
                margin-bottom: 12px;

                .tool-icon {
                    width: 40px;
                    height: 40px;
                    border-radius: 8px;
                    margin-right: 12px;
                    object-fit: cover;
                }

                .tool-name {
                    line-height: 24px;
                    font-size: 16px;
                    font-weight: 600;
                    color: #262626;
                }

                .tool-category {
                    line-height: 16px;
                    font-size: 12px;
                    font-weight: 400;
                    color: #7a8699;
                }
            }

            .tool-info {
                .tool-description {
                    height: 80px;
                    line-height: 20px;
                    font-size: 12px;
                    font-weight: 400;
                    color: #7a8699;
                    text-overflow: ellipsis;
                    overflow: hidden;
                }
            }
        }
    }

    .tools-grid {
        grid-template-columns: repeat(1, 1fr);
    }

    /* 默认（包含 1200 宽）：3 列 */
    @media (min-width: 990px) {
        .tools-grid {
            grid-template-columns: repeat(2, 1fr);
        }
    }

    /* 默认（包含 1200 宽）：3 列 */
    @media (min-width: 1200px) {
        .tools-grid {
            grid-template-columns: repeat(3, 1fr);
        }
    }

    /* 1920 宽度时显示 4 列 */
    @media (min-width: 1440px) {
        .tools-grid {
            grid-template-columns: repeat(4, 1fr);
        }
    }

    /* 宽度大于 1920 时显示 5 列 */
    @media (min-width: 1920px) {
        .tools-grid {
            grid-template-columns: repeat(5, 1fr);
        }
    }
}


</style>
