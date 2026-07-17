<style lang="less" scoped>
.add-wechat-app-alert {
  padding-top: 16px;
  .zm-alert-info .title {
    font-size: 15px;
    font-weight: 600;
  }

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
    display: flex;
    align-items: center;
    margin-top: 16px;
  }

  .preview-box {
    flex: 1;
    margin-left: 32px;

    .preview-img {
      display: block;
      width: 100%;
      min-height: 378px;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 2px 16px 0 #00000029;
    }
  }

  .form-box {
    width: 100%;

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
  <a-modal
    width="700px"
    v-model:open="show"
    :title="`${formState.id ? '编辑' : '创建'}Telegram接入渠道`"
    @cancel="handleCancel"
  >
    <template #footer>
      <a-button key="back" @click="handleCancel" v-if="step === 1">{{ t('btn_cancel') }}</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="handleSave" v-if="step === 1"
        >提交
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
      <a-alert class="zm-alert-info" type="info">
        <template #message>
          <div class="title">操作步骤</div>
          <div>第一步：在 Telegram 搜索<a>@BotFather</a>并单击<a>/start</a></div>
          <div>
            第二步：发送 <a>/mybots</a> 并按照说明进行操作，如需帮助，请查看<a
              href="https://help.chatwiki.com/zh/docs/configuring-telegram"
              target="_blank"
              >详细操作文档</a
            >
          </div>
          <div>第三步：创建机器人后，您将收到一条带有令牌的消息。复制令牌并将其粘贴到此处</div>
        </template>
      </a-alert>

      <div class="alert-body">
        <div class="form-box">
          <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
            <a-form-item :label="t('label_robot_avatar_name')" name="app_name">
              <PageTitleInput
                :autoUpload="false"
                v-model:avatar="formState.app_avatar_url"
                v-model:value.trim="formState.app_name"
                :placeholder="t('ph_enter_robot_name')"
                @changeAvatar="onChangeAvatar"
              />
            </a-form-item>
            <a-form-item label="Telegram令牌" name="app_id">
              <a-input v-model:value.trim="formState.app_id" placeholder="请输入Telegram令牌" />
            </a-form-item>
          </a-form>
        </div>
      </div>
    </div>
    <div class="add-wechat-app-alert" v-if="step === 2">
      <div class="config-items">
        <div class="config-item">
          <span class="config-value">{{ step2Info.push_url }}</span>
          <span class="copy-btn" @click="handleCopy(step2Info.push_url)">{{ t('btn_copy') }}</span>
        </div>
      </div>
      <div class="config-preview-box"></div>
    </div>
  </a-modal>
</template>

<script setup>
import { saveWechatApp } from '@/api/robot'
import { ref, reactive, computed, toRaw, inject } from 'vue'
import { copyText } from '@/utils/index'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import PageTitleInput from './page-title-input.vue'

const { t } = useI18n('views.robot.robot-config.external-service.components.add-dingding-alert')
const emit = defineEmits(['ok'])
const defaultAvatar = '/upload/default/telegram_robot.png'

const { robotInfo } = inject('robotInfo')

const show = ref(false)
const step = ref(1)
const loading = ref(false)
const formRef = ref()
const formStateStruct = {
  id: void 0,
  robot_id: robotInfo.id,
  app_type: 'telegram_robot',
  app_avatar: '',
  app_avatar_url: defaultAvatar,
  app_name: '',
  app_id: ''
}
const formState = reactive({
  id: void 0,
  robot_id: robotInfo.id,
  app_type: 'telegram_robot',
  app_avatar: '',
  app_avatar_url: defaultAvatar,
  app_name: '',
  app_id: ''
})

const formRules = {
  app_name: [
    {
      required: true,
      message: t('rule_required_robot_name'),
      trigger: 'change'
    },
    {
      trigger: 'change',
      validator: () => {
        if (!formState.app_avatar_url) {
          return Promise.reject(t('rule_required_avatar'))
        } else {
          return Promise.resolve()
        }
      }
    }
  ],
  app_id: [
    {
      required: true,
      message: '请输入Telegram令牌',
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

  saveWechatApp(data).then((res) => {
    Object.assign(step2Info, res.data)
    // step.value = 2
    handleCancel()
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

const copyIp = () => {
  handleCopy(robotInfo.wechat_ip)
}

const open = (data) => {
  if (data && data.id) {
    data = JSON.parse(JSON.stringify(data))
    formState.app_avatar_url = data.app_avatar
    formState.app_avatar = ''
    formState.id = data.id
    formState.robot_id = data.robot_id
    formState.app_name = data.app_name
    formState.app_id = data.app_id
  } else {
    formState.id = void 0
    Object.assign(formState, JSON.parse(JSON.stringify(formStateStruct)))
  }

  step.value = 1
  show.value = true
}

defineExpose({
  open
})
</script>
