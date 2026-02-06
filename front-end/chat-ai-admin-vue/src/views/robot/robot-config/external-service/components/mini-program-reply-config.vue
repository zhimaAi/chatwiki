<template>
    <a-modal v-model:open="visible" :title="t('title_mini_program_reply_settings')" :width="480" :footer="null">
        <div class="reply-config-content">
            <!-- 选项1 -->
            <div class="setting-item">
                <div class="setting-label">
                    <span>{{ t('label_ai_generated_content') }}</span>
                    <a-tooltip>
                        <template #title>
                            <img class="tip-img" src="@/assets/img/robot/official_account_reply_config01.png" alt="">
                        </template>
                        <QuestionCircleOutlined class="info-icon" />
                    </a-tooltip>
                </div>
                <div class="switch-wrapper">
                    <a-switch v-model:checked="aiGenerated"
                        @change="() => handleSave('show_ai_msg_mini', aiGenerated)" />
                    <span class="switch-text">{{ aiGenerated ? t('switch_on') : t('switch_off') }}</span>
                </div>
            </div>

            <!-- 选项2 -->
            <div class="setting-item">
                <div class="setting-label">
                    <span>{{ t('label_typing_indicator') }}</span>
                    <a-tooltip>
                        <template #title>
                            <img class="tip-img" src="@/assets/img/robot/official_account_reply_config02.png" alt="">
                        </template>
                        <QuestionCircleOutlined class="info-icon" />
                    </a-tooltip>
                </div>
                <div class="switch-wrapper">
                    <a-switch v-model:checked="typingIndicator"
                        @change="() => handleSave('show_typing_mini', typingIndicator)" />
                    <span class="switch-text">{{ typingIndicator ? t('switch_on') : t('switch_off') }}</span>
                </div>
            </div>

        </div>
    </a-modal>
</template>

<script setup>
import { setWechatConfigSwitch } from '@/api/robot'
import { ref } from 'vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.external-service.components.mini-program-reply-config')

const emit = defineEmits(['change'])

const visible = ref(false)
const aiGenerated = ref(true)
const typingIndicator = ref(true)
const robotId = ref(null)

const open = (config = {}) => {
    if (config.aiGenerated !== undefined) {
        aiGenerated.value = config.aiGenerated == 1 ? true : false
    }

    if (config.typingIndicator !== undefined) {
        typingIndicator.value = config.typingIndicator == 1 ? true : false
    }

    robotId.value = config.robotId

    visible.value = true
}

const handleSave = async (key, value) => {
    try {
        const res = await setWechatConfigSwitch({
            id: robotId.value,
            key: key,
            val: value ? 1 : 0
        })

        if (res.code == 0) {
            message.success(t('msg_save_success'))
            //   visible.value = false
            emit('change')
        }
    } catch (error) {
        console.error(t('msg_save_failed'), error)
        message.error(t('msg_save_failed'))
    }
}



defineExpose({
    open
})
</script>

<style lang="less" scoped>
.tip-img {
    width: 230px;
    border-radius: 6px;
    display: block;
}

.reply-config-content {
    padding: 8px 0 24px;
}

.setting-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 0;
    border-bottom: 1px solid #f0f0f0;

    &:last-of-type {
        border-bottom: none;
    }
}

.setting-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);

    .info-icon {
        color: rgba(0, 0, 0, 0.45);
        font-size: 14px;
        cursor: help;

        &:hover {
            color: #1890ff;
        }
    }
}

.switch-wrapper {
    display: flex;
    align-items: center;
    gap: 12px;

    .switch-text {
        font-size: 14px;
        color: rgba(0, 0, 0, 0.65);
        min-width: 14px;
    }
}

.progress-section {
    margin-top: 24px;
    padding: 16px;
    background: #fafafa;
    border-radius: 8px;

    .progress-label {
        font-size: 14px;
        color: rgba(0, 0, 0, 0.85);
        margin-bottom: 12px;
    }

    .progress-value {
        font-size: 18px;
        color: #52c41a;
        font-weight: 600;
        margin-top: 8px;
    }
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    border-top: 1px solid #f0f0f0;
    padding-top: 16px;
}
</style>
