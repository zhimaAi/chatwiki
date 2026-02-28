<template>
    <div v-if="!isCollapsed" class="user-info-wrapper" @click="handleSettingClick">
        <div class="user-info-left">
            <img class="user-avatar" :src="userStore.avatar" alt="avatar" />
            <span class="user-id">{{ userStore.user_name }}</span>
        </div>
        <div class="user-info-right">
            <svg-icon class="setting-icon" name="reply-settings"></svg-icon>
        </div>
    </div>
    <Tooltip v-else :title="t('tooltip_home_config')" placement="right">
        <div class="nav-item" @click="handleSettingClick">
            <svg-icon class="nav-icon" name="reply-settings"></svg-icon>
        </div>
    </Tooltip>
</template>

<script setup>
import { Tooltip } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workbench.components.user-info')

const props = defineProps({
    isCollapsed: {
        type: Boolean,
        default: false
    }
})

const userStore = useUserStore()
const router = useRouter()

const handleSettingClick = () => {
    router.push('/workbench/home-config')
}
</script>

<style lang="less" scoped>
.user-info-wrapper {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    background-color: #F2F4F7;
    border-radius: 6px;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
        background-color: #E4E6EB;
    }


    .user-info-left {
        display: flex;
        align-items: center;
        gap: 8px;
        flex: 1;
        min-width: 0;

        .user-avatar {
            width: 24px;
            height: 24px;
            border-radius: 50%;
            flex-shrink: 0;
        }

        .user-id {
            font-size: 14px;
            font-weight: 400;
            color: #595959;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    }

    .user-info-right {
        display: flex;
        align-items: center;
        flex-shrink: 0;

        .setting-icon {
            width: 16px;
            height: 16px;
            font-size: 16px;
            color: #262626;
            cursor: pointer;
            transition: color 0.2s;

            &:hover {
                color: #595959;
            }
        }
    }
}

// 收起状态下的导航样式
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

    &:hover {
        background-color: #E4E6EB;
    }

    .nav-icon {
        width: 18px;
        height: 18px;
        font-size: 18px;
    }
}
</style>
