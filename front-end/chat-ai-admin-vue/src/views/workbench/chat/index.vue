<template>
    <div class="workbench-chat">
        <!-- 初始化中状态 -->
        <div v-if="isInitializing" class="loading-state">
            <div class="loading-spinner"></div>
            <div class="loading-text">{{ t('msg_loading') }}</div>
        </div>

        <!-- 空状态 -->
        <div v-else-if="showEmptyState" class="empty-state">
            <a-empty description="暂无应用" />
        </div>

        <!-- iframe 加载中状态 -->
        <div v-else-if="!iframeSrc" class="loading-state">
            <div class="loading-spinner"></div>
            <div class="loading-text">{{ t('msg_loading') }}</div>
        </div>

        <!-- iframe 容器 -->
        <div v-else class="iframe-container" :class="{ 'iframe-loading': iframeLoading }">
            <div v-if="iframeLoading" class="iframe-overlay">
                <div class="loading-spinner"></div>
            </div>
            <iframe
                :key="iframeKey"
                class="iframe"
                :src="iframeSrc"
                frameborder="0"
                @load="workbenchStore.handleIframeLoad"
            ></iframe>
        </div>
    </div>
</template>

<script setup>
import { computed, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useWorkbenchStore } from '@/stores/modules/workbench'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workbench.chat.index')
const workbenchStore = useWorkbenchStore()

// 从 store 获取响应式状态
const { iframeLoading, iframeKey, showEmptyState, isInitializing } = storeToRefs(workbenchStore)

// 使用 store 中的 getters
const iframeSrc = computed(() => workbenchStore.iframeSrc)

// 初始化时如果有有效的 iframeSrc，设置 loading 状态
if (workbenchStore.iframeSrc) {
    workbenchStore.iframeLoading = true
}

// 监听 iframeSrc 变化，刷新 iframe
watch(iframeSrc, (newVal, oldVal) => {
    if (newVal && newVal !== oldVal) {
        workbenchStore.refreshIframe()
    }
})
</script>

<style lang="less" scoped>
.workbench-chat {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background-color: #f5f7fa;

    .empty-state {
        display: flex;
        align-items: center;
        justify-content: center;
        flex: 1;
    }

    .loading-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        flex: 1;
        gap: 16px;

        .loading-spinner {
            width: 40px;
            height: 40px;
            border: 3px solid #e4e6eb;
            border-top-color: #1890ff;
            border-radius: 50%;
            animation: spin 1s linear infinite;
        }

        .loading-text {
            font-size: 14px;
            color: #999;
        }
    }

    .iframe-container {
        position: relative;
        flex: 1;
        overflow: hidden;

        &.iframe-loading {
            .iframe {
                opacity: 0.5;
            }
        }

        .iframe-overlay {
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            background-color: rgba(255, 255, 255, 0.8);
            z-index: 10;

            .loading-spinner {
                width: 40px;
                height: 40px;
                border: 3px solid #e4e6eb;
                border-top-color: #1890ff;
                border-radius: 50%;
                animation: spin 1s linear infinite;
            }
        }

        .iframe {
            width: 100%;
            height: 100%;
            border: 0;
            transition: opacity 0.3s ease;
        }
    }
}

@keyframes spin {
    to {
        transform: rotate(360deg);
    }
}
</style>
