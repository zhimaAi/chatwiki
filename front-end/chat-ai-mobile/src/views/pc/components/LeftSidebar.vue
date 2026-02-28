<template>
    <div class="left-sidebar-wrapper" :class="{ 'is-hide': sidebarHide }" v-if="hasSessions">
        <div class="left-sidebar">
            <div class="sidebar-header">
                <!-- <h2 class="sidebar-title">{{ props.title }}</h2> -->
                <div class="header-actions">
                    <button class="btn-new-chat" @click="handleNewChat">
                        <svg-icon name="new-chat-btn" style="font-size: 14px;" />
                        <span>{{ t('btn_new_chat') }}</span>
                    </button>
                    <button class="btn-delete" @click="handleDelete" >
                        <svg-icon name="mini-broom" />
                    </button>
                </div>
            </div>
            <SessionList :groups="sessionGroups" :active-id="activeSessionId" @select="handleSelect" v-if="!sidebarHide" />
        </div>
        <div class="sidebar-handle-wrapper">
            <a-tooltip :title="sidebarHide ? t('tooltip_expand') : t('tooltip_collapse')" placement="right">
                <span class="sidebar-handle" @click="onHandleClick">
                    <span class="handle-line handle-line01"></span>
                    <span class="handle-line handle-line02"></span>
                </span>
            </a-tooltip>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { showConfirmDialog } from 'vant'
import { Tooltip as ATooltip } from 'ant-design-vue'
import SessionList, { type TimeGroup } from './SessionList.vue'
import { useChatStore } from '@/stores/modules/chat'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.pc.components.left-sidebar')

// interface Props {
//     title?: string
// }

// const props = withDefaults(defineProps<Props>(), {
//     title: '聊天'
// })

const emit = defineEmits<{
    newChat: []
    delete: []
    select: [id: string]
}>()

const chatStore = useChatStore()
const { myChatList, dialogue_id } = storeToRefs(chatStore)
const {  delDialogue } = chatStore

// 侧边栏收起状态
let defaultData = localStorage.getItem('sidebar_hide_left') === '1'
const sidebarHide = ref(defaultData)

function onHandleClick() {
    sidebarHide.value = !sidebarHide.value
    localStorage.setItem('sidebar_hide_left', sidebarHide.value ? '1' : '0')
}

// 当前选中的会话ID
const activeSessionId = computed(() => String(dialogue_id.value || ''))

// 是否有会话
const hasSessions = computed(() => (myChatList.value || []).length > 0)

// 将会话列表按时间分组
const sessionGroups = computed<TimeGroup[]>(() => {
    const list = myChatList.value || []
    if (list.length === 0) return []

    const today: any[] = []
    const yesterday: any[] = []
    const thisMonth: any[] = []
    const others: any[] = []

    // 获取当前时间
    const now = new Date()
    const nowYear = now.getFullYear()
    const nowMonth = now.getMonth()
    const nowDate = now.getDate()

    // 今天开始时间戳（秒）
    const todayStart = new Date(nowYear, nowMonth, nowDate).getTime() / 1000
    // 昨天开始时间戳（秒）
    const yesterdayStart = todayStart - 24 * 60 * 60
    // 本月开始时间戳（秒）
    const monthStart = new Date(nowYear, nowMonth, 1).getTime() / 1000

    list.forEach((item) => {
        const session = {
            id: String(item.id),
            title: item.subject || t('label_new_chat')
        }

        const createTime = item.create_time || 0

        if (createTime >= todayStart) {
            today.push(session)
        } else if (createTime >= yesterdayStart) {
            yesterday.push(session)
        } else if (createTime >= monthStart) {
            thisMonth.push(session)
        } else {
            others.push(session)
        }
    })

    const groups: TimeGroup[] = []
    if (today.length) groups.push({ label: t('label_today'), sessions: today })
    if (yesterday.length) groups.push({ label: t('label_yesterday'), sessions: yesterday })
    if (thisMonth.length) groups.push({ label: t('label_this_month'), sessions: thisMonth })
    if (others.length) groups.push({ label: t('label_earlier'), sessions: others })

    return groups
})

// 新建对话
const handleNewChat = () => {
    emit('newChat')
}

// 清空所有记录
const handleDelete = () => {
    showConfirmDialog({
        title: t('title_clear_confirm'),
        message: t('msg_clear_all_records')
    })
        .then(() => {
            delDialogue({ id: -1 })
            emit('delete')
        })
        .catch(() => {})
}

// 选择会话
const handleSelect = (item: any) => {
    emit('select', item)
}
</script>

<style lang="less" scoped>
.left-sidebar-wrapper {
    display: flex;
    height: 100%;
}

.left-sidebar {
    width: 225px;
    height: 100%;
    background: #fff;
    border-right: 1px solid #F0F0F0;
    display: flex;
    flex-direction: column;
    overflow: hidden;
    transition: width 0.3s ease, opacity 0.3s ease;

    .header-actions {
        flex-shrink: 0;
    }
}

.sidebar-handle-wrapper {
    width: 18px;
    height: 100%;
    position: relative;
    background-color: #fff;
}

.sidebar-handle {
    position: absolute;
    right: 0;
    top: 50%;
    width: 12px;
    height: 26px;
    transform: translateY(-50%);
    cursor: pointer;
    transition: all 0.3s ease;

    .handle-line {
        position: absolute;
        width: 4px;
        height: 13px;
        left: 4px;
        transition: all 0.3s ease;
        background-color: #bfbfbf;
    }

    .handle-line01 {
        top: 0;
        border-top-left-radius: 4px;
        border-top-right-radius: 4px;
        transform-origin: 50% 0;
    }

    .handle-line02 {
        bottom: 0;
        border-bottom-left-radius: 4px;
        border-bottom-right-radius: 4px;
        transform-origin: 50% 100%;
    }
}

.sidebar-handle:hover {
    .handle-line01 {
        background-color: #595959;
        transform: rotate(18deg) translateY(0);
        border-top-left-radius: 4px;
        border-top-right-radius: 4px;
        border-bottom-left-radius: 10px;
        height: 16px;
    }

    .handle-line02 {
        background-color: #595959;
        transform: rotate(-18deg) translateY(0);
        border-bottom-left-radius: 4px;
        border-bottom-right-radius: 4px;
        border-top-left-radius: 10px;
        height: 16px;
    }
}

.left-sidebar-wrapper.is-hide {
    .left-sidebar {
        width: 0;
        opacity: 0;
        border-right: none;
    }

    .sidebar-handle {
        opacity: 1 !important;
    }

    .sidebar-handle:hover {
        .handle-line01 {
            background-color: #595959;
            transform: rotate(-18deg) translateY(0);
            border-bottom-left-radius: 4px;
            border-bottom-right-radius: 4px;
            border-top-left-radius: 10px;
            height: 16px;
        }

        .handle-line02 {
            background-color: #595959;
            transform: rotate(18deg) translateY(0);
            border-top-left-radius: 4px;
            border-top-right-radius: 4px;
            border-bottom-left-radius: 10px;
            height: 16px;
        }
    }
}

.sidebar-header {
    padding: 16px 16px 20px 16px;
}

.sidebar-title {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    margin: 0 0 16px 0;
}

.header-actions {
    display: flex;
    align-items: center;
    gap: 8px;
}

.btn-new-chat {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 6px;
    height: 32px;
    padding: 0 16px;
    border-radius: 6px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
    overflow: hidden;
    color: #595959;
    background: #FFF;
    border: 1px solid #D9D9D9;

    &:hover {
        color: #2475FC;
        background: #D6E6FF;
        border: 1px solid #D6E6FF;
    }
}

.btn-delete{
    display: flex;
    width: 32px;
    height: 32px;
    padding: 0;
    font-size: 16px;
    border-radius: 6px;
    justify-content: center;
    align-items: center;
    border: 1px solid #D9D9D9;
    background: #FFF;
    color: #595959;
    cursor: pointer;
    transition: all 0.2s;
    &:hover{
        background: #E4E6EB;
    }
}
</style>
