<template>
    <div class="workbench-home-config">
        <div class="page-title">{{ t('title_page') }}</div>
        <div class="config-card">
            <div class="config-item">
                <div class="item-label">{{ t('label_home_app') }}</div>
                <a-select
                    v-model:value="selectedApp"
                    class="item-control"
                    size="middle"
                    placeholder="请选择应用"
                    :options="appOptions"
                    @change="handleAppChange"
                ></a-select>
            </div>
            <div class="config-item">
                <div class="item-label">{{ t('label_default_last_entry') }}</div>
                <a-switch
                    v-model:checked="enableLastEntry"
                    class="item-switch"
                    @change="handleEnableChange"
                />
            </div>
        </div>



    </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue'

import { message } from 'ant-design-vue'
import { useWorkbenchStore } from '@/stores/modules/workbench'
import { getWorkbenchTeamRobotList, saveWorkbenchConfig } from '@/api/workbench'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workbench.home-config.index')

const workbenchStore = useWorkbenchStore()

const appOptions = ref([])

const selectedApp = ref()
const enableLastEntry = ref(false)
const saving = ref(false)
const isSyncing = ref(true)
const hasInitialized = ref(false)



// 统一下拉选项结构
const normalizeOptions = (list = []) => {

    return list
        .filter(item => item && (item.robot_id || item.id))
        .map(item => ({
            label: item.robot_name || t('label_unnamed_app'),
            value: item.robot_id || item.id,
            robot_key: item.robot_key
        }))
}

// 从配置回填默认首页应用
const syncSelectedApp = () => {
    const robotId = workbenchStore.getWorkbenchConfig?.config?.default_robot_id

    // 验证 robot_id 是否在 appOptions 中存在
    if (robotId) {
        const exists = appOptions.value.some(option => option.value === robotId)
        selectedApp.value = exists ? robotId : undefined
    } else {
        selectedApp.value = undefined
    }
}

// 从配置回填“默认进入上一次访问应用”开关
const syncEnableFromConfig = () => {
    const flag = workbenchStore.getWorkbenchConfig?.config?.enable_last_app_entry
    enableLastEntry.value = flag === '1' || flag === 1
}


// 获取团队应用列表
const fetchAppOptions = async () => {

    try {
        const res = await getWorkbenchTeamRobotList()
        appOptions.value = normalizeOptions(res?.data || [])
    } catch (error) {
        console.error('Failed to fetch home apps:', error)
        appOptions.value = []
    }
}

// 同步回填期间禁用自动保存
const withSyncing = async (fn) => {
    isSyncing.value = true
    try {
        await fn?.()
    } finally {
        await nextTick()
        isSyncing.value = false
    }
}

// 保存配置（支持静默保存）
const handleSave = async (silent = false) => {
    if (saving.value) {
        return
    }

    saving.value = true
    try {
        const payload = {
            default_robot_id: selectedApp.value,
            enable_last_app_entry: enableLastEntry.value ? 1 : 0,
        }

        await saveWorkbenchConfig(payload)

        message.success({
            content: t('msg_save_success'),
            key: 'home-config-save',
            duration: 1.6
        })
        await workbenchStore.fetchWorkbenchConfig()
        await withSyncing(() => {
            syncEnableFromConfig()
            syncSelectedApp(appOptions.value)
        })
    } catch (error) {
        console.error('Failed to save home config:', error)
        message.error(t('msg_save_failed'))
    } finally {
        saving.value = false
    }
}


// 首页应用 change 触发自动保存
const handleAppChange = (value) => {
    selectedApp.value = value
    if (isSyncing.value || !hasInitialized.value) {
        return
    }

    handleSave(true)
}

// 开关 change 触发自动保存
const handleEnableChange = (checked) => {
    enableLastEntry.value = checked
    if (isSyncing.value || !hasInitialized.value) {
        return
    }

    handleSave(true)
}



// 监听配置变更，回填表单
watch(() => workbenchStore.getWorkbenchConfig, async () => {

    await withSyncing(() => {
        syncEnableFromConfig()
        syncSelectedApp()
    })
}, { immediate: true })

// 初始化：回填配置并拉取应用列表
onMounted(async () => {

    await withSyncing(async () => {
        syncEnableFromConfig()
        await fetchAppOptions()
        syncSelectedApp()
    })
    await nextTick()
    hasInitialized.value = true
})


</script>




<style lang="less" scoped>
.workbench-home-config {
    min-height: 100vh;
    padding: 24px 32px;
    background-color: #ffffff;

    .page-title {
        line-height: 26px;
        margin-bottom: 18px;
        font-size: 18px;
        font-weight: 600;
        color: #262626;
    }

    .config-card {
        background-color: transparent;
        border-radius: 12px;
        padding: 0;
        max-width: 620px;
        box-shadow: none;
        display: flex;
        flex-direction: column;
        gap: 12px;

        .config-item {
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 16px;
            padding: 12px 16px;
            border-radius: 8px;
            background-color: #f3f5f7;

            .item-label {
                font-size: 14px;
                color: #262626;
                white-space: nowrap;
                font-weight: 500;
            }

            .item-control {
                min-width: 280px;
            }

            .item-switch {
                margin-left: auto;
            }
        }

    }
}
</style>




