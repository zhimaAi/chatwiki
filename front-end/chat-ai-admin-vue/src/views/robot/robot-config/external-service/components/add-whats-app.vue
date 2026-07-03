<style lang="less" scoped>
.add-wechat-app-alert {
  padding-top: 16px;

  :deep(.slick-slider) {
    width: 620px;
    overflow: hidden;

    .slick-dots li {
      background: red;
      border-radius: 4px;
      &.slick-active button {
        opacity: 0;
        background: unset;
      }
    }
  }

  .tip-alert {
    margin-bottom: 24px;
  }

  .alert-body {
    align-items: center;
    margin-top: 16px;
  }

  .step-box {
    margin: 12px 0;
    display: flex;
    flex-direction: column;
    gap: 8px;
    .step-title {
      color: #595959;
      font-size: 16px;
      font-weight: 600;
    }
    .step-desc {
      color: #8c8c8c;
      font-size: 14px;
      font-weight: 400;
    }
  }

  .form-box {
    .my-ip {
      display: flex;
      align-items: center;
      height: 32px;
      padding: 0 12px;
      border-radius: 2px;
      border: 1px solid #d9d9d9;
      background: #f0f0f0;

      .ip-text {
        flex: 1;
      }
    }
  }

  .config-items {
    width: 600px;

    .config-item {
      display: flex;
      height: 32px;
      margin-bottom: 9px;
      background-color: #fff;
      box-shadow:
        0 1px 10px 0 #0000000d,
        0 4px 5px 0 #00000014,
        0 2px 4px -1px #0000001f;

      .config-value {
        flex: 1;
        height: 32px;
        line-height: 32px;
        padding: 0 8px;
        font-size: 14px;
        border-radius: 2px;
        color: #595959;
        border: 1px dashed #2475fc;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      .copy-btn {
        width: 50px;
        height: 32px;
        line-height: 32px;
        font-size: 14px;
        border-radius: 2px;
        color: #fff;
        text-align: center;
        background-color: #2475fc;
        cursor: pointer;
      }
    }
  }

  .config-preview-box {
    position: relative;
    width: 842px;
    min-height: 378px;
    margin: 24px auto 0;

    .config-preview-img {
      display: block;
      height: 378px;
    }
  }
}
</style>

<template>
  <a-modal width="700px" v-model:open="show" title="添加账号" @cancel="handleCancel">
    <template #footer>
      <a-button key="back" @click="handleCancel" v-if="step === 1">{{ t('btn_cancel') }}</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="handleSave" v-if="step === 1"
        >确定
      </a-button>
      <a-button
        key="submit"
        type="primary"
        :loading="loading"
        @click="handleOk"
        v-if="step === 2"
        >{{ t('btn_config_completed') }}</a-button
      >
    </template>
    <div class="add-wechat-app-alert" v-if="step === 1">
      <a-alert class="zm-alert-info" type="info" show-icon>
        <template #message>
          <div>
            通过阿里云ChatApp连接官方WhatsApp商业API，配置前请阅读
            <a href="https://help.chatwiki.com/zh/docs/configuring-WhatsApp" target="_blank"
              >帮助文档</a
            >
          </div>
        </template>
      </a-alert>

      <div class="alert-body">
        <div class="form-box">
          <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
            <div class="step-box">
              <div class="step-title">第一步</div>
              <div class="step-desc">
                <a href="https://chatapp.console.aliyun.com/Overview" target="_blank">登录</a>
                阿里云ChatApp后台，创建WhatsApp通道，绑定WABA账号和电话号，如果已绑定，跳过此步骤
              </div>
            </div>
            <div class="step-box">
              <div class="step-title">第二步</div>
              <div class="step-desc">绑定成功后，查看通道ID和绑定的号码</div>
            </div>
            <a-form-item label="头像/昵称" name="app_name">
              <PageTitleInput
                :autoUpload="false"
                v-model:avatar="formState.app_avatar_url"
                v-model:value.trim="formState.app_name"
                placeholder="请输入"
                @changeAvatar="onChangeAvatar"
              />
            </a-form-item>
            <a-form-item label="通道ID" name="cust_space_id">
              <a-input v-model:value.trim="formState.cust_space_id" placeholder="请输入" />
            </a-form-item>
            <a-form-item label="电话号码" name="app_id">
              <a-input v-model:value.trim="formState.app_id" placeholder="请输入" />
            </a-form-item>

            <div class="step-box">
              <div class="step-title">第三步</div>
              <div class="step-desc">在阿里云管理控制台获取AccessKeyId和AccessKeySecret</div>
            </div>

            <a-form-item label="AccessKeyId" name="access_key_id">
              <a-input v-model:value.trim="formState.access_key_id" placeholder="请输入" />
            </a-form-item>

            <a-form-item label="AccessKeySecret" name="access_key_secret">
              <a-input v-model:value.trim="formState.access_key_secret" placeholder="请输入" />
            </a-form-item>
          </a-form>
        </div>
      </div>
    </div>
    <div class="add-wechat-app-alert" v-if="step === 2">
      <a-alert
        class="tip-alert"
        message="复制下方URL，粘贴到阿里云ChatApp控制台 - [通道管理] - [WABA管理] - [通道Webhook设置]界面的通知回调地址里，并且开启http协议"
        type="info"
        show-icon
      />
      <div class="config-items">
        <div style="margin-bottom: 8px;">传入消息的URL</div>
        <div class="config-item">
          <span class="config-value">{{ step2Info.push_url }}</span>
          <span class="copy-btn" @click="handleCopy(step2Info.push_url)">{{ t('btn_copy') }}</span>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { saveWechatApp } from '@/api/robot'
import { ref, reactive, toRaw, inject } from 'vue'
import { copyText } from '@/utils/index'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import PageTitleInput from './page-title-input.vue'

const { t } = useI18n('views.robot.robot-config.external-service.components.add-dingding-alert')
const emit = defineEmits(['ok'])
const defaultAvatar = '/upload/default/whatsapp_avatar.png'

const { robotInfo } = inject('robotInfo')

const show = ref(false)
const step = ref(1)
const loading = ref(false)
const formRef = ref()
const formStateStruct = {
  id: undefined,
  robot_id: robotInfo.id,
  app_type: 'whatsapp',
  app_avatar: '',
  app_avatar_url: defaultAvatar,
  app_name: '',
  app_id: '',
  access_key_id: '',
  access_key_secret: '',
  cust_space_id: ''
}
const formState = reactive({})

const formRules = {
  app_name: [
    {
      required: true,
      message: '请输入头像/昵称',
      trigger: 'change'
    }
  ],
  cust_space_id: [
    {
      required: true,
      message: '请输入通道ID',
      trigger: 'change'
    }
  ],
  app_id: [
    {
      required: true,
      message: '请输入电话号码',
      trigger: 'change'
    }
  ],
  access_key_id: [
    {
      required: true,
      message: '请输入AccessKeyId',
      trigger: 'change'
    }
  ],
  access_key_secret: [
    {
      required: true,
      message: '请输入AccessKeySecret',
      trigger: 'change'
    }
  ]
}

const onChangeAvatar = ({ file, url }) => {
  formState.app_avatar = file
  formState.app_avatar_url = url
}

const step2Info = reactive({
  push_aeskey: '',
  push_token: '',
  push_url: ''
})

const submitForm = () => {
  let data = { ...toRaw(formState) }

  delete data.app_avatar_url

  saveWechatApp(data).then((res) => {
    Object.assign(step2Info, res.data)
    step.value = 2
    message.success(t('msg_save_success'))

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

const handleOk = () => {
  show.value = false
}

const handleCancel = () => {
  formRef.value.clearValidate()
  formRef.value.resetFields()
  show.value = false
}

const handleCopy = (text) => {
  copyText(text)
  message.success(t('msg_copy_success'))
}

const open = (data) => {
  if (!data) {
    data = JSON.parse(JSON.stringify(formStateStruct))
  } else {
    data.app_avatar_url = data.app_avatar
    data.app_avatar = ''
  }
  step.value = 1
  Object.assign(formState, data)
  show.value = true
}

defineExpose({
  open
})
</script>
