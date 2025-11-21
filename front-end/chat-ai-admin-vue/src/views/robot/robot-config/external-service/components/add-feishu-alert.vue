<style lang="less" scoped>
.add-wechat-app-alert {
  padding-top: 16px;

  :deep(.slick-slider) {
    width: 620px;
    overflow: hidden;

    .slick-dots li{
      background: red;
      border-radius: 4px;
      &.slick-active button{
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

  .config-items {
    width: 600px;

    .config-item {
      display: flex;
      height: 32px;
      margin-bottom: 9px;
      background-color: #fff;
      box-shadow: 0 1px 10px 0 #0000000d,
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
  <a-modal width="1000px" v-model:open="show" title="配置飞书机器人" @cancel="handleCancel">
    <template #footer>
      <a-button key="back" @click="handleCancel" v-if="step === 1">取 消</a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="handleSave" v-if="step === 1">保存并进入下一步
      </a-button>
      <a-button key="submit" type="primary" :loading="loading" @click="handleOk" v-if="step === 2">已完成配置</a-button>
    </template>
    <div class="add-wechat-app-alert" v-if="step === 1">
      <a-alert class="tip-alert" type="info" show-icon>
        <template #message>
          <div>登录飞书开放后台，按照文档完成配置井发布应用 <a href="https://open.feishu.cn" target="_blank">登录后台</a>，<a href="https://www.yuque.com/zhimaxiaoshiwangluo/pggco1/cpoq2kmobhgap70p?singleDoc#" target="_blank">如何配置？</a></div>
          <div>应用添加完成后，复制APPID 和API sercet,加密信息：Encrypt Key , Verification Token，填写到以下配置中</div>
        </template>
      </a-alert>

      <div class="alert-body">
        <div class="form-box">
          <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
            <a-form-item label="机器人头像/名称" name="app_name">
              <PageTitleInput
                :autoUpload="false"
                v-model:avatar="formState.app_avatar_url"
                v-model:value.trim="formState.app_name"
                placeholder="请输入您的机器人名称"
                @changeAvatar="onChangeAvatar"
              />
            </a-form-item>
            <a-form-item label="APPID" name="app_id">
              <a-input v-model:value.trim="formState.app_id" placeholder="请输入您的AppId"/>
            </a-form-item>
            <a-form-item label="API Secret" name="app_secret">
              <a-input v-model:value.trim="formState.app_secret" placeholder="请输入您的APP Secret"/>
            </a-form-item>
            <a-form-item label="Encrypt Key" name="encrypt_key">
              <a-input v-model:value.trim="formState.encrypt_key" placeholder="请输入您的Encrypt Key"/>
            </a-form-item>
            <a-form-item label="Verification Token" name="verification_token">
              <a-input v-model:value.trim="formState.verification_token" placeholder="请输入您的Verification Token"/>
            </a-form-item>
          </a-form>
        </div>
        <div class="preview-box">
          <a-carousel>
            <div><img class="preview-img" src="@/assets/img/robot/feishu-example01.png"/></div>
            <div><img class="preview-img" src="@/assets/img/robot/feishu-example02.png"/></div>
          </a-carousel>
        </div>
      </div>
    </div>
    <div class="add-wechat-app-alert" v-if="step === 2">
      <a-alert class="tip-alert"
               message="复制以下htp地址到应用-事件与回调>订阅方式-选择梅事件发送开发者服务器的请求地址中" type="info"
               show-icon/>
      <div class="config-items">
        <div class="config-item">
          <span class="config-value">{{ step2Info.push_url }}</span>
          <span class="copy-btn" @click="handleCopy(step2Info.push_url)">复制</span>
        </div>
      </div>
      <div class="config-preview-box">
        <img class="config-preview-img" src="@/assets/img/robot/feishu-example03.png"/>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {saveWechatApp} from '@/api/robot'
import {ref, reactive, computed, toRaw, inject} from 'vue'
import {copyText} from '@/utils/index'
import {message} from 'ant-design-vue'
import PageTitleInput from './page-title-input.vue'

const emit = defineEmits(['ok'])
const defaultAvatar = '/upload/default/feishu_robot_avatar.png'

const {robotInfo} = inject('robotInfo')

const show = ref(false)
const step = ref(1)
const loading = ref(false)
const formRef = ref()
const formStateStruct = {
  id: undefined,
  robot_id: robotInfo.id,
  app_type: 'feishu_robot',
  app_avatar: '',
  app_avatar_url: defaultAvatar,
  app_name: '',
  app_id: '',
  app_secret: '',
  encrypt_key: '',
  verification_token: '',
}
const formState = reactive({})

const formRules = {
  app_name: [
    {
      required: true,
      message: '请输入您的机器人名称',
      trigger: 'change'
    },
    {
      trigger: 'change',
      validator: () => {
        if (!formState.app_avatar_url) {
          return Promise.reject('请上传机器人头像')
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
      message: '请输入您的Api Secret',
      trigger: 'change'
    }
  ],
  encrypt_key: [
    {
      required: true,
      message: '请输入您的Encrypt Key',
      trigger: 'change'
    }
  ],
  verification_token: [
    {
      required: true,
      message: '请输入您的Verification Token',
      trigger: 'change'
    }
  ]
}

const onChangeAvatar = ({file, url}) => {
  formState.app_avatar = file
  formState.app_avatar_url = url
}

const step2Info = reactive({
  push_aeskey: '',
  push_token: '',
  push_url: ''
})

const submitForm = () => {
  let data = {...toRaw(formState)}

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
