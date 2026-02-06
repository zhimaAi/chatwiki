<template>
  <a-modal
    :width="749"
    v-model:open="show"
    :title="null"
    wrapClassName="remove-model-padding"
    :footer="null"
  >
    <div class="form-content-box">
      <div class="left-box">
        <div class="form-body">
          <div class="model-title">{{ t('unverified_account_reply_settings') }}</div>
          <a-alert
            class="zm-alert-info"
            :message="t('api_limit_notice')"
            type="info"
          />
          <a-form
            style="margin-top: 24px"
            ref="formRef"
            layout="vertical"
            :model="formState"
            :rules="formRules"
          >
            <a-form-item :label="t('manual_reply_prompt')" name="wechat_not_verify_hand_get_reply">
              <a-textarea
                v-model:value="formState.wechat_not_verify_hand_get_reply"
                :placeholder="t('please_input')"
                :maxLength="100"
                @blur="handleBlur"
                :auto-size="{ minRows: 3, maxRows: 3 }"
              />
            </a-form-item>
            <a-form-item :label="t('manual_reply_text')" name="wechat_not_verify_hand_get_word">
              <a-input
                :maxLength="100"
                @blur="handleBlur"
                v-model:value="formState.wechat_not_verify_hand_get_word"
                :placeholder="t('please_input')"
              />
            </a-form-item>
            <a-form-item
              :label="t('continue_reading_text')"
              name="wechat_not_verify_hand_get_next"
            >
              <a-input
                :maxLength="100"
                @blur="handleBlur"
                v-model:value="formState.wechat_not_verify_hand_get_next"
                :placeholder="t('please_input')"
              />
            </a-form-item>
          </a-form>
        </div>
        <div class="footer-box">
          <a-button @click="handleCancel">{{ t('cancel') }}</a-button>
          <a-button type="primary" @click="handleSave">{{ t('confirm') }}</a-button>
        </div>
      </div>
      <div class="preview-box">
        <img
          src="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/wechat-app/guide1.png"
          alt=""
        />
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, toRaw, inject } from 'vue'
import { setWechatNotVerifyConfig } from '@/api/robot'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const emit = defineEmits(['ok'])

const { robotInfo, getRobot } = inject('robotInfo')
const { t } = useI18n('views.robot.robot-config.external-service.components.add-unverified-alert')

const show = ref(false)
const formRef = ref()

const wechat_not_verify_hand_get_reply_default = t('default_prompt')
const wechat_not_verify_hand_get_word_default = t('default_reply_text')
const wechat_not_verify_hand_get_next_default = t('default_continue_text')

const formState = reactive({
  id: robotInfo.id,
  wechat_not_verify_hand_get_reply: wechat_not_verify_hand_get_reply_default,
  wechat_not_verify_hand_get_word: wechat_not_verify_hand_get_word_default,
  wechat_not_verify_hand_get_next: wechat_not_verify_hand_get_next_default
})

const handleBlur = () => {
  formState.wechat_not_verify_hand_get_reply =
    formState.wechat_not_verify_hand_get_reply || wechat_not_verify_hand_get_reply_default
  formState.wechat_not_verify_hand_get_word =
    formState.wechat_not_verify_hand_get_word || wechat_not_verify_hand_get_word_default
  formState.wechat_not_verify_hand_get_next =
    formState.wechat_not_verify_hand_get_next || wechat_not_verify_hand_get_next_default
  formRef.value.clearValidate()
}

const formRules = {
  wechat_not_verify_hand_get_reply: [
    {
      required: true,
      message: t('error_prompt'),
      trigger: 'change'
    }
  ],
  wechat_not_verify_hand_get_word: [
    {
      required: true,
      message: t('error_reply_text'),
      trigger: 'change'
    }
  ],
  wechat_not_verify_hand_get_next: [
    {
      required: true,
      message: t('error_continue_text'),
      trigger: 'change'
    }
  ]
}

const submitForm = () => {
  let data = { ...toRaw(formState) }

  setWechatNotVerifyConfig(data).then((res) => {
    message.success(t('save_success'))
    handleCancel()
    getRobot(robotInfo.id)
    emit('ok')
  })
}

const handleSave = () => {
  formRef.value
    .validate()
    .then(() => {
      submitForm()
    })
    .catch((error) => {
      console.log('error', error)
    })
}

const handleCancel = () => {
  formRef.value.clearValidate()
  formRef.value.resetFields()
  show.value = false
}

const open = () => {
  formState.wechat_not_verify_hand_get_reply =
    robotInfo.wechat_not_verify_hand_get_reply || wechat_not_verify_hand_get_reply_default
  formState.wechat_not_verify_hand_get_word =
    robotInfo.wechat_not_verify_hand_get_word || wechat_not_verify_hand_get_word_default
  formState.wechat_not_verify_hand_get_next =
    robotInfo.wechat_not_verify_hand_get_next || wechat_not_verify_hand_get_next_default
  show.value = true
}

defineExpose({
  open
})
</script>

<style lang="less" scoped>
.form-content-box {
  display: flex;
  .left-box {
    width: 389px;
    display: flex;
    flex-direction: column;
    .form-body {
      padding: 24px;
      flex: 1;
      .model-title {
        color: #262626;
        font-size: 20px;
        font-weight: 600;
        line-height: 28px;
        margin-bottom: 24px;
      }
    }
    .footer-box {
      width: 100%;
      height: 64px;
      display: flex;
      align-items: center;
      justify-content: flex-end;
      gap: 8px;
      padding: 0 24px;
      box-shadow: 0 -8px 4px 0 #0000000a;
      overflow: hidden;
    }
  }
  .preview-box {
    flex: 1;
    background: #e5efff;
    display: flex;
    justify-content: center;
    padding: 39px 0;
    img {
      width: 240px;
      height: 522px;
    }
  }
}
</style>
