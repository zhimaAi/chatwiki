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
    .config-preview-img {
      display: block;
      height: 378px;
    }
    .config-items {
      position: absolute;
      right: 61px;
      top: 124px;
      width: 432px;
      .config-item {
        display: flex;
        height: 24px;
        margin-bottom: 9px;
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
  }
}
</style>

<template>
  <a-modal width="1168px" v-model:open="show" :title="title" @cancel="handleCancel">
    <template #footer>
      <a-button key="back" @click="handleCancel" v-if="step === 1">取 消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="handleSave" v-if="step === 1"
        >保存并进入下一步</a-button
      >
      <a-button key="submit" type="primary" :loading="loading" @click="handleOk" v-if="step === 2"
        >已完成配置</a-button
      >
    </template>
    <div class="add-wechat-app-alert" v-if="step === 1">
      <a-alert
        class="tip-alert"
        message="进入微信公众号后台（mp.weixin.qq.com），在【设置与开发-基本配置】获取开发者信息。"
        type="info"
        show-icon
      />

      <div class="alert-body">
        <div class="form-box">
          <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
            <a-form-item label="公众号头像/名称" name="app_name">
              <PageTitleInput
                :autoUpload="false"
                v-model:avatar="formState.app_avatar_url"
                v-model:value.trim="formState.app_name"
                placeholder="请输入您的公众号名称"
                @changeAvatar="onChangeAvatar"
              />
            </a-form-item>
            <a-form-item label="AppId" name="app_id">
              <a-input v-model:value.trim="formState.app_id" placeholder="请输入您的AppId" />
            </a-form-item>
            <a-form-item label="AppSecret" name="app_secret">
              <a-input
                v-model:value.trim="formState.app_secret"
                placeholder="请输入您的AppSecret"
              />
            </a-form-item>
            <a-form-item label="将以下IP添加到公众号IP白名单中" name="wechat_ip" required>
              <div class="my-ip">
                <span class="ip-text">{{ robotInfo.wechat_ip }}</span>
                <a class="copy-btn" @click="copyIp">复 制</a>
              </div>
            </a-form-item>
          </a-form>
        </div>
        <div class="preview-box">
          <img class="preview-img" src="@/assets/img/robot/preview_01.png" />
        </div>
      </div>
    </div>
    <div class="add-wechat-app-alert" v-if="step === 2">
      <a-alert
        class="tip-alert"
        message="进入微信公众号后台（mp.weixin.qq.com），在【设置与开发-基本配置】中配置并启用服务器配置。请按照下图配置"
        type="info"
        show-icon
      />
      <div class="config-preview-box">
        <div class="config-items">
          <div class="config-item">
            <span class="config-value">{{ step2Info.push_url }}</span>
            <span class="copy-btn" @click="handleCopy(step2Info.push_url)">复制</span>
          </div>

          <div class="config-item">
            <span class="config-value">{{ step2Info.push_token }}</span>
            <span class="copy-btn" @click="handleCopy(step2Info.push_token)">复制</span>
          </div>

          <div class="config-item">
            <span class="config-value">{{ step2Info.push_aeskey }}</span>
            <span class="copy-btn" @click="handleCopy(step2Info.push_aeskey)">复制</span>
          </div>
        </div>
        <img class="config-preview-img" src="@/assets/img/robot/config_preview_01.png" />
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { saveWechatApp } from '@/api/robot'
import { ref, reactive, computed, toRaw, inject } from 'vue'
import { copyText } from '@/utils/index'
import { message } from 'ant-design-vue'
import PageTitleInput from './page-title-input.vue'

const emit = defineEmits(['ok'])
const defaultAvatar = '/upload/default/official_account_avatar.png'

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
  app_type: 'official_account'
})

const formRules = {
  app_name: [
    {
      required: true,
      message: '请输入您的公众号名称',
      trigger: 'change'
    },
    {
      trigger: 'change',
      validator: () => {
        if (!formState.app_avatar_url) {
          return Promise.reject('请上传公众号头像')
        } else {
          return Promise.resolve()
        }
      }
    }
  ],
  app_id: [
    {
      required: true,
      message: '请输入您的AppId',
      trigger: 'change'
    }
  ],
  app_secret: [
    {
      required: true,
      message: '请输入您的AppSecret',
      trigger: 'change'
    }
  ]
}

const title = computed(() => {
  let text = !formState.id ? '添加微信公众号' : '编辑微信公众号'
  let text2 = step.value == 1 ? '-填写公众号开发信息' : '-服务器配置'

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
    message.success('保存成功')

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
  message.success('复制成功')
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
      app_type: 'official_account',
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
