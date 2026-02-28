<template>
    <div class="left-sidebar-wrapper">
        <!-- 占位 div，撑开文档流 -->
        <div class="sidebar-placeholder" :class="{ collapsed: isCollapsed }"></div>

        <!-- Fixed 定位的 sidebar -->
        <div class="sidebar" :class="{ collapsed: isCollapsed }">
            <!-- Header -->
            <div class="sidebar-header">
                <div class="logo-wrapper">
                    <img class="logo-icon" src="@/assets/logo.svg" alt="">
                    <img class="min-logo-icon" src="@/assets/default_logo.svg" alt="">
                </div>
                <span class="header-actions">
                    <Tooltip :title="t('tooltip_collapse')" placement="right">
                        <svg-icon class="collapse-icon" name="collapse-close" @click="toggleCollapse"></svg-icon>
                    </Tooltip>
                </span>
            </div>

            <!-- Navigation -->
            <nav class="sidebar-nav">
                <Tooltip :title="t('tooltip_expand')" placement="right" v-if="isCollapsed">
                    <div class="nav-item" @click="toggleCollapse">
                        <svg-icon class="nav-icon" name="collapse-open"></svg-icon>
                    </div>
                </Tooltip>
                <Tooltip :title="isCollapsed ? t('nav_home') : ''" placement="right">
                    <div class="nav-item" :class="{ active: isHomeActive }"
                        @click="handleHomeClick">
                        <svg-icon class="nav-icon" name="chat-plus"></svg-icon>
                        <span class="nav-text">{{ t('nav_home') }}</span>
                    </div>
                </Tooltip>

                <Tooltip :title="isCollapsed ? t('nav_team_apps') : ''" placement="right">
                    <div class="nav-item" :class="{ active: activeNav === 'workbenchTeamTools' }"
                        @click="handleNavClick('workbenchTeamTools', '/workbench/team-tools')">
                        <svg-icon class="nav-icon" name="app-outlined"></svg-icon>
                        <span class="nav-text">{{ t('nav_team_apps') }}</span>
                    </div>
                </Tooltip>
                <a-popover title="" placement="rightTop" :arrow="false" v-if="isCollapsed">
                    <template #content>
                        <div class="recent-list-popup">
                            <recent-list></recent-list>
                        </div>
                    </template>
                    <div class="nav-item" @click="toggleCollapse">
                        <svg-icon class="nav-icon" name="history"></svg-icon>
                    </div>
                </a-popover>
            </nav>

            <!-- Recently Used -->
            <div class="recent-section">
                <div class="section-title">{{ t('section_recent') }}</div>
                <recent-list></recent-list>
            </div>

            <div class="user-info-section" v-if="role_type == 1">
                <user-info :is-collapsed="isCollapsed"></user-info>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Tooltip } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import RecentList from './recent-list.vue'
import UserInfo from './user-info.vue'
import { useWorkbenchStore } from '@/stores/modules/workbench'
import { usePermissionStore } from '@/stores/modules/permission'

const { t } = useI18n('views.workbench.components.left-sidebar')
const router = useRouter()
const route = useRoute()
const workbenchStore = useWorkbenchStore()
const permissionStore = usePermissionStore()

const activeNav = ref(null)
const isCollapsed = ref(workbenchStore.sidebarCollapsed || false)

const role_type = computed(() => {
    return permissionStore.role_type
})
// 首页是否高亮：在 chat 页面且 selectedMode 为 'home'
const isHomeActive = computed(() => {
    return route.name === 'workbenchChat' && workbenchStore.selectedMode === 'home'
})

const handleHomeClick = async () => {
    activeNav.value = 'workbenchChat'
    try {
        await workbenchStore.refreshHomeConfig()
    } catch (error) {
        console.error('刷新首页配置失败:', error)
    }

    router.push('/workbench/chat')
}



const handleNavClick = (navItem, path) => {
    activeNav.value = navItem
    router.push(path)
}


const toggleCollapse = () => {
    isCollapsed.value = !isCollapsed.value
    workbenchStore.toggleSidebar()
}



// 监听路由变化更新激活状态
watch(() => route.path, () => {
    activeNav.value = route.name
})

// 监听侧边栏折叠状态变化
watch(() => workbenchStore.sidebarCollapsed, (newValue) => {
    isCollapsed.value = newValue
})

// 组件挂载时初始化激活状态
onMounted(() => {
    activeNav.value = route.name
})
</script>

<style lang="less" scoped>
@sidebar-transition-duration: 0.3s;

.left-sidebar-wrapper {
    position: relative;

    .sidebar-placeholder {
        width: 256px;
        height: 100vh;
        transition: width @sidebar-transition-duration;
    }

    .sidebar-placeholder.collapsed {
        width: 52px;
    }
}

.sidebar {
    position: fixed;
    width: 256px;
    height: 100vh;
    top: 0;
    left: 0;
    z-index: 1000;
    background-color: #F2F4F7;
    border-right: 1px solid #e5e7eb;
    display: flex;
    flex-direction: column;
    padding: 12px 8px;
    transition: width @sidebar-transition-duration;

    &-header {
        display: flex;
        align-items: center;
        padding: 0;
        margin-bottom: 12px;
        display: flex;
        align-items: center;

        .logo-wrapper {
            flex: 1;
            overflow: hidden;
            white-space: nowrap;
        }

        .logo-icon {
            width: 116px;
            height: 32px;
        }

        .min-logo-icon {
            display: none;
            width: 32px;
            height: 32px;
        }

        .header-actions {
            display: flex;
            align-items: center;
        }

        .collapse-icon {
            width: 18px;
            height: 18px;
            font-size: 18px;
            color: #262626;
            cursor: pointer;
        }
    }

    &-nav {
        margin-bottom: 12px;

        .nav-item {
            display: flex;
            align-items: center;
            justify-content: flex-start;
            padding: 7px 9px;
            line-height: 22px;
            margin-bottom: 4px;
            border-radius: 6px;
            color: #262626;
            cursor: pointer;
            transition: background-color 0.2s;

            &:last-child {
                margin-bottom: 0;
            }

            &:hover {
                background-color: #E4E6EB;
            }

            &.active {
                color: #2475FC;
                background-color: #FFFFFF;
            }

            .nav-icon {
                width: 18px;
                height: 18px;
                font-size: 18px;
                margin-right: 8px;
                flex-shrink: 0;
            }


            .nav-text {
                flex: 1;
                font-size: 14px;
                font-weight: 400;
                overflow: hidden;
                white-space: nowrap;
            }
        }
    }
}

.sidebar.collapsed {
    width: 52px;

    .sidebar-header {
        justify-content: start;

        .logo-icon {
            display: none;
        }

        .min-logo-icon {
            display: block;
        }
    }

    .header-actions {
        display: none;
    }

    .nav-text {
        display: none;
    }

    .recent-section {
        display: none;
    }
}

.recent-section {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;

    .section-title {
        width: 100%;
        line-height: 22px;
        padding: 5px 9px;
        font-size: 12px;
        font-weight: 400;
        color: #9ca3af;
        overflow: hidden;
        white-space: nowrap;
    }
}
.recent-list-popup{
    width: 200px;
}

.user-info-section {
    width: 100%;
    overflow: hidden;
    margin-top: auto;
    padding-top: 12px;
}
</style>
