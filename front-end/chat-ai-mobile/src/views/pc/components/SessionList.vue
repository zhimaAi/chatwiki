<template>
    <div class="session-list" ref="scrollContainer" @scroll="handleScroll">
        <div v-for="group in props.groups" :key="group.label" class="time-group">
            <div class="time-label">{{ group.label }}</div>
            <div v-for="session in group.sessions" :key="session.id" class="session-item"
                :class="{ active: props.activeId === session.id }" @click="handleSelect(session)">
                <svg-icon class="session-icon" name="message" />
                <span class="session-text">{{ session.title }}</span>
                <div class="operate-btn-wrapper" @click.stop>
                    <van-popover
                        v-model:show="popoverStates[session.id]"
                        :actions="actions"
                        @select="(action) => onSelect(action, session)"
                    >
                        <template #reference>
                            <div class="operate-btn">
                                <svg-icon name="point-h" />
                            </div>
                        </template>
                    </van-popover>
                </div>
            </div>
        </div>
    </div>
    
    <!-- 重命名弹窗 -->
    <van-dialog
        v-model:show="showDialog"
        :title="t('title_rename')"
        @confirm="handleSave"
        @cancel="showDialog = false"
        show-cancel-button
    >
        <div class="input-box">
            <van-field v-model="textValue" :placeholder="t('ph_input')" />
        </div>
    </van-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue'
import { showConfirmDialog, showToast } from 'vant'
import { useChatStore } from '@/stores/modules/chat'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.pc.components.session-list')

export interface Session {
    id: string
    title: string
}

export interface TimeGroup {
    label: string
    sessions: Session[]
}

interface Props {
    groups: TimeGroup[]
    activeId?: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
    select: [session: Session]
}>()

const chatStore = useChatStore()
const { delDialogue, editDialogueChat, getMyChatList } = chatStore
const popoverStates = reactive<Record<string, boolean>>({})
const showDialog = ref(false)
const textValue = ref('')
let currentSession: Session | null = null

const actions = [
    { text: t('btn_rename'), key: 1 },
    { text: t('btn_delete'), key: 2 }
]

const handleSelect = (session: Session) => {
    emit('select', session)
}

const onSelect = (action: any, session: Session) => {
    currentSession = session
    // 关闭 popover
    popoverStates[session.id] = false
    
    if (action.key === 1) {
        // 重命名
        textValue.value = session.title
        showDialog.value = true
    } else if (action.key === 2) {
        // 删除
        handleDelete(session)
    }
}

const handleDelete = (session: Session) => {
    showConfirmDialog({
        title: t('title_delete_confirm'),
        message: t('msg_delete_record')
    })
        .then(() => {
            delDialogue({ id: session.id })
        })
        .catch(() => {})
}

const handleSave = () => {
    if (!textValue.value.trim()) {
        showToast(t('msg_input_name'))
        return
    }
    if (currentSession) {
        editDialogueChat({
            id: currentSession.id,
            subject: textValue.value
        }).then((res: any) => {
            if (res.res === 0) {
                showToast(t('msg_update_success'))
                showDialog.value = false
            }
        })
    }
}

// 滚动加载更多
const handleScroll = (event: Event) => {
    const target = event.target as HTMLElement
    const { scrollTop, scrollHeight, clientHeight } = target
    // 距离底部 50px 时触发加载
    if (scrollHeight - scrollTop - clientHeight < 50) {
        getMyChatList()
    }
}
</script>

<style lang="less" scoped>
.session-list {
    flex: 1;
    padding: 0 16px;
    overflow-y: auto;
    // 隐藏滚动条但仍可滚动
    -ms-overflow-style: none;
    scrollbar-width: none;
}

.session-list::-webkit-scrollbar {
    width: 0px;
    background: transparent;
}

.time-group {
    margin-bottom: 20px;
}

.time-label {
    line-height: 22px;
    padding: 0 16px;
    font-size: 14px;
    color: #8c8c8c;
}

.session-item {
    display: flex;
    align-items: center;
    gap: 4px;
    line-height: 22px;
    padding: 8px 8px 8px 16px;
    margin: 8px 0;
    border-radius: 6px;
    color: #595959;
    cursor: pointer;
    transition: background 0.2s;


    &:hover {
        background: #E4E6EB;
    }

    &.active {
        background: #E6EFFF;

        .session-text {
            color: #595959;
        }
    }

    .session-icon {
        width: 16px;
        height: 16px;
        font-size: 16px;
    }

    .session-text {
        flex: 1;
        font-size: 14px;
        color: #666;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .operate-btn-wrapper {
        width: 24px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        transition: opacity 0.2s;
    }

    &:hover .operate-btn-wrapper {
        opacity: 1;
    }

    .operate-btn {
        width: 24px;
        height: 24px;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 16px;
        color: #595959;
        border-radius: 6px;
        cursor: pointer;
        &:hover {
            background: #E4E6EB;
        }
    }
}

.input-box {
    border-radius: 6px;
    margin: 24px;
    :deep(.van-cell) {
        background: #f3f3f3;
        border-radius: 6px;
    }
}
</style>
