<style lang="less" scoped>
.add-wechat-app-alert {
  padding-top: 16px;
  .tip-alert {
    margin-bottom: 30px;
  }
  .alert-body {
    display: flex;
    align-items: center;
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
    width: 299px;
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

  .config-preview-box {
    position: relative;
    width: 842px;
    min-height: 378px;
    margin: 0 auto;

    .config-items {
      position: absolute;
      right: 82px;
      top: 186px;
      width: 432px;
      .config-item {
        display: flex;
        height: 24px;
        margin-bottom: 10px;
        background-color: #fff;
        box-shadow:
          0 1px 10px 0 #0000000d,
          0 4px 5px 0 #00000014,
          0 2px 4px -1px #0000001f;

        .config-value {
          flex: 1;
          height: 24px;
          line-height: 22px;
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
          height: 24px;
          line-height: 24px;
          font-size: 14px;
          border-radius: 2px;
          color: #fff;
          text-align: center;
          background-color: #2475fc;
          cursor: pointer;
        }
      }
    }

    .config-preview-img {
      width: 100%;
    }
  }
}
</style>

<template>
  <a-modal width="1168px" v-model:open="show" :title="title" @cancel="handleCancel">
    <template #footer>
      <a-button key="back" @click="handleCancel" v-if="step === 1">{{ t('btn_cancel') }}</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="handleSave" v-if="step === 1"
        >{{ t('btn_save_and_next') }}</a-button
      >
      <a-button key="submit" type="primary" :loading="loading" @click="handleOk" v-if="step === 2"
        >{{ t('btn_config_completed') }}</a-button
      >
    </template>
    <div class="add-wechat-app-alert" v-if="step === 1">
      <a-alert
        class="tip-alert"
        :message="t('tip_step1')"
        type="info"
        show-icon
      />

      <div class="alert-body">
        <div class="form-box">
          <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
            <a-form-item :label="t('label_mini_program_avatar_name')" name="app_name">
              <PageTitleInput
                :autoUpload="false"
                v-model:avatar="formState.app_avatar_url"
                v-model:value.trim="formState.app_name"
                :placeholder="t('ph_enter_mini_program_name')"
                @changeAvatar="onChangeAvatar"
              />
            </a-form-item>
            <a-form-item :label="t('label_app_id')" name="app_id">
              <a-input v-model:value.trim="formState.app_id" :placeholder="t('ph_enter_app_id')" />
            </a-form-item>
            <a-form-item :label="t('label_app_secret')" name="app_secret">
              <a-input
                v-model:value.trim="formState.app_secret"
                :placeholder="t('ph_enter_app_secret')"
              />
            </a-form-item>
            <a-form-item
              :label="t('label_add_ip_to_whitelist')"
              name="wechat_ip"
              required
            >
              <div class="my-ip">
                <span class="ip-text">{{ robotInfo.wechat_ip }}</span>
                <a class="copy-btn" @click="copyIp">{{ t('btn_copy') }}</a>
              </div>
            </a-form-item>
          </a-form>
        </div>
        <div class="preview-box">
          <img class="preview-img" src="@/assets/img/robot/preview_02.png" />
        </div>
      </div>
    </div>
    <div class="add-wechat-app-alert" v-if="step === 2">
      <a-alert
        class="tip-alert"
        :message="t('tip_step2')"
        type="info"
        show-icon
      />
      <div class="config-preview-box">
        <div class="config-items">
          <div class="config-item">
            <span class="config-value">{{ step2Info.push_url }}</span>
            <span class="copy-btn" @click="handleCopy(step2Info.push_url)">{{ t('btn_copy') }}</span>
          </div>

          <div class="config-item">
            <span class="config-value">{{ step2Info.push_token }}</span>
            <span class="copy-btn" @click="handleCopy(step2Info.push_token)">{{ t('btn_copy') }}</span>
          </div>

          <div class="config-item">
            <span class="config-value">{{ step2Info.push_aeskey }}</span>
            <span class="copy-btn" @click="handleCopy(step2Info.push_aeskey)">{{ t('btn_copy') }}</span>
          </div>
        </div>
        <img class="config-preview-img" src="@/assets/img/robot/config_preview_02.png" />
      </div>
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

const { t } = useI18n('views.robot.robot-config.external-service.components.add-wechat-mini-program-alert')

const emit = defineEmits(['ok'])
const defaultAvatar = '/upload/default/mini_program_avatar.png'

const { robotInfo } = inject('robotInfo')

const show = ref(false)
const step = ref(1)
const loading = ref(false)
const formRef = ref()
const formState = reactive({
  id: undefined,
  robot_id: robotInfo.id,
  app_name: '',
  app_id: '',
  app_secret: '',
  app_avatar: '',
  app_avatar_url: defaultAvatar,
  app_type: 'mini_program'
})

const formRules = {
  app_name: [
    {
      required: true,
      message: t('msg_enter_official_account_name'),
      trigger: 'change'
    },
    {
      trigger: 'change',
      validator: () => {
        if (!formState.app_avatar_url) {
          return Promise.reject(t('msg_upload_official_account_avatar'))
        } else {
          return Promise.resolve()
        }
      }
    }
  ],
  app_id: [
    {
      required: true,
      message: t('msg_enter_app_id'),
      trigger: 'change'
    }
  ],
  app_secret: [
    {
      required: true,
      message: t('msg_enter_app_secret'),
      trigger: 'change'
    }
  ]
}

const title = computed(() => {
  let text = !formState.id ? t('title_add_wechat_mini_program') : t('title_edit_wechat_mini_program')
  let text2 = step.value == 1 ? t('title_step1') : t('title_step2')

  return text + text2
})

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

const copyIp = () => {
  handleCopy(robotInfo.wechat_ip)
}

const open = (data) => {
  if (!data) {
    data = {
      id: undefined,
      robot_id: robotInfo.id,
      app_name: '',
      app_id: '',
      app_secret: '',
      app_avatar: '',
      app_avatar_url: defaultAvatar,
      app_type: 'mini_program',
      wechat_ip: robotInfo.wechat_ip
    }
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
