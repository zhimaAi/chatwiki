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
          <div class="model-title">未认证公众号回复设置</div>
          <a-alert
            class="zm-alert-info"
            message="由于微信接口限制,未认证公众号只能在用户提问后立即回复答案。若答案生成较慢,系统将提示用户手动获取回复"
            type="info"
          />
          <a-form
            style="margin-top: 24px"
            ref="formRef"
            layout="vertical"
            :model="formState"
            :rules="formRules"
          >
            <a-form-item label="手动获取回复提示语" name="wechat_not_verify_hand_get_reply">
              <a-textarea
                v-model:value="formState.wechat_not_verify_hand_get_reply"
                placeholder="请输入"
                :maxLength="100"
                @blur="handleBlur"
                :auto-size="{ minRows: 3, maxRows: 3 }"
              />
            </a-form-item>
            <a-form-item label="手动获取回复蓝字文案" name="wechat_not_verify_hand_get_word">
              <a-input
                :maxLength="100"
                @blur="handleBlur"
                v-model:value="formState.wechat_not_verify_hand_get_word"
                placeholder="请输入"
              />
            </a-form-item>
            <a-form-item
              label="内容超500字截断，获取下文蓝字文案"
              name="wechat_not_verify_hand_get_next"
            >
              <a-input
                :maxLength="100"
                @blur="handleBlur"
                v-model:value="formState.wechat_not_verify_hand_get_next"
                placeholder="请输入"
              />
            </a-form-item>
          </a-form>
        </div>
        <div class="footer-box">
          <a-button @click="handleCancel">取 消</a-button>
          <a-button type="primary" @click="handleSave">确 定</a-button>
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

const emit = defineEmits(['ok'])

const { robotInfo, getRobot } = inject('robotInfo')

const show = ref(false)
const formRef = ref()

const wechat_not_verify_hand_get_reply_default = '正在思考中，请稍后点击下方蓝字\r\n获取回复👇👇👇'
const wechat_not_verify_hand_get_word_default = '👉👉点我获取回复👈👈'
const wechat_not_verify_hand_get_next_default = '内容较多，点此查看下文'

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
      message: '请输入手动获取回复提示语',
      trigger: 'change'
    }
  ],
  wechat_not_verify_hand_get_word: [
    {
      required: true,
      message: '请输入手动获取回复蓝字文案',
      trigger: 'change'
    }
  ],
  wechat_not_verify_hand_get_next: [
    {
      required: true,
      message: '请输入内容超500字截断，获取下文蓝字文案',
      trigger: 'change'
    }
  ]
}

const submitForm = () => {
  let data = { ...toRaw(formState) }

  setWechatNotVerifyConfig(data).then((res) => {
    message.success('保存成功')
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
