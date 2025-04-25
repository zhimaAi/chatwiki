<style lang="less" scoped>
.web-app-box {
  display: flex;

  .form-box {
    .form-item {
      margin-bottom: 16px;
    }
    .form-item-label {
      line-height: 22px;
      padding-bottom: 8px;
      font-size: 14px;
      color: rgba(0, 0, 0, 0.88);
    }
  }

  .box-left {
    flex: 1;
    width: calc(100% - 487px);
    .box-wrapper {
      margin-bottom: 16px;
    }
  }
  .web-app-info {
    .web-app-link {
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }
    .link-action {
      margin-top: 8px;
      .action-btn {
        margin-right: 8px;
      }
    }
    .access-restrictions {
      margin-top: 16px;
    }
  }
}
.box-right {
  margin: 0 96px 0 48px;
  .demo-box{
    position: relative;
  }
  iframe{
    width: 375px;
    height: 720px;
    border-radius: 4px;
    box-shadow: 0 4px 32px 0 rgba(0, 0, 0, 0.16);
  }
  .preview-img {
    display: block;
    border-radius: 9px;
    box-shadow: 0 4px 32px 0 rgba(0, 0, 0, 0.16);
  }
}
</style>

<template>
  <div class="web-app-box">
    <div class="box-left">
      <div class="box-wrapper">
        <card-box title="WebApp 链接">
          <template #icon>
            <svg-icon name="circularNeedle" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <div class="web-app-info">
            <div class="web-app-link">
              <a :href="h5_website" target="_blank">{{ h5_website }}</a>
            </div>
            <div class="link-action">
              <a-button class="action-btn" type="primary" ghost @click="copyH5WebSite"
                >复 制</a-button
              >
              <a-button class="action-btn" @click="handlePreview">预 览</a-button>
              <a-tooltip color="#fff" placement="top">
                <template #title>
                  <img style="width: 180px" :src="previewQrcode" alt="" />
                </template>
                <a-button class="action-btn">二维码</a-button>
              </a-tooltip>
            </div>
            <div class="access-restrictions form-box">
              <div class="form-item">
                <div class="form-item-label">访问限制</div>
                <div class="form-item-body">
                  <a-radio-group
                    v-model:value="accessRestrictionsType"
                    @change="saveAccessRestrictionsType"
                  >
                    <a-radio :value="1"><span class="default-text-color">无限制</span></a-radio>
                    <a-radio :value="2"
                      ><span class="default-text-color">登录后才可访问</span></a-radio
                    >
                  </a-radio-group>
                </div>
              </div>
            </div>
          </div>
        </card-box>
      </div>

      <div class="box-wrapper">
        <card-box title="样式设置">
          <template #icon>
            <svg-icon name="phone" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <template #action>
            <a-button @click="saveForm" size="small" type="primary">保存</a-button>
          </template>
          <div class="web-app-style form-box">
            <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
              <a-form-item class="form-item" label="是否显示标题栏" name="navbarShow">
                <a-radio-group v-model:value="formState.navbarShow" name="navbarShow">
                  <a-radio :value="1">显示</a-radio>
                  <a-radio :value="2">不显示</a-radio>
                </a-radio-group>
              </a-form-item>
              <a-form-item class="form-item" label="页面标题" name="pageTitle">
                <PageTitleInput
                  v-model:avatar="formState.logo"
                  v-model:value="formState.pageTitle"
                />
              </a-form-item>

              <a-form-item
                class="form-item"
                label="标题栏颜色"
                :name="['pageStyle', 'navbarBackgroundColor']"
              >
                <ColorPicker v-model:value="formState.pageStyle.navbarBackgroundColor" />
              </a-form-item>

              <a-form-item class="form-item" label="语言" name="lang">
                <a-select
                  style="width: 180px"
                  v-model:value="formState.lang"
                  placeholder="请选择语言"
                >
                  <a-select-option value="zh-CN">简体中文</a-select-option>
                  <!-- <a-select-option value="en-US">English</a-select-option> -->
                </a-select>
              </a-form-item>
            </a-form>
          </div>
        </card-box>
      </div>
      <div class="box-wrapper">
        <QuickInstruction :type="robotInfo.app_id" @updata="updataQuickComand"></QuickInstruction>
      </div>
    </div>
    <div class="box-right">
      <div class="demo-box">
        <iframe id="mobile-preview" :src="previewIframeSrc" frameborder="0"></iframe>
      </div>
    </div>
  </div>
</template>

<script setup>
import QRCode from 'qrcode'
import { ref, reactive, toRaw, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import { storeToRefs } from 'pinia'
import { copyText } from '@/utils/index'
import { useRobotStore } from '@/stores/modules/robot'
import { editExternalConfig } from '@/api/robot/index'
import CardBox from './card-box.vue'
import PageTitleInput from './page-title-input.vue'
import ColorPicker from '@/components/color-picker/index.vue'
import QuickInstruction from './quick-instruction.vue'
import PreviewCommand from './preview-command.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const robotStore = useRobotStore()
const { robotInfo, external_config_h5 } = storeToRefs(robotStore)
const { h5_website } = robotInfo.value
const { getRobot } = robotStore
const previewQrcode = ref('')

const formRef = ref()
const formState = reactive({
  logo: external_config_h5.value.logo,
  pageTitle: external_config_h5.value.pageTitle,
  navbarShow: external_config_h5.value.navbarShow,
  lang: external_config_h5.value.lang,
  pageStyle: external_config_h5.value.pageStyle
})

const previewIframeSrc = computed(()=>{
  return h5_website
})
watch(formState,(val)=>{
  updatePreview(val)
})
const updataQuickComand = (data) =>{
  updatePreview(data, 'updataQuickComand')
}
const updatePreview = (data, type) =>{
  let iframe = document.getElementById('mobile-preview');
  iframe.contentWindow.postMessage(
    {
      type: type || 'onPreview',
      data: JSON.parse(JSON.stringify(data)),
    },
    '*'
  );
}

const formRules = {
  lang: [
    {
      required: true,
      message: '请选择语言',
      trigger: 'change'
    }
  ],
  pageTitle: [
    {
      required: true,
      message: '请输入标题',
      trigger: 'input'
    },
    {
      trigger: 'input',
      validator: () => {
        if (!formState.logo) {
          return Promise.reject('请上传logo')
        } else {
          return Promise.resolve()
        }
      }
    }
  ],
  navbarShow: [
    {
      required: true,
      message: '请选择是否显示标题栏',
      trigger: 'change'
    }
  ],
  pageStyle: {
    navbarBackgroundColor: [
      {
        required: true,
        message: '请选择标题栏颜色',
        trigger: 'change'
      }
    ]
  }
}
// 保存访问限制
const accessRestrictionsType = ref(external_config_h5.value.accessRestrictionsType || 1)
const saveAccessRestrictionsType = () => {
  let formData = { ...toRaw(external_config_h5.value) }

  formData.accessRestrictionsType = accessRestrictionsType.value

  saveWebAppInfo(formData)
}
// 保存样式设置
const saveWebAppInfo = (formData) => {
  const { id } = robotInfo.value

  editExternalConfig({
    id: id,
    external_config_h5: JSON.stringify(formData)
  }).then(() => {
    getRobot(id)
    message.success('保存成功')
    // 刷新一下
    router.go(0)
  })
}

const saveForm = () => {
  formRef.value
    .validate()
    .then(() => {
      let formData = { ...toRaw(formState) }

      formData.accessRestrictionsType = accessRestrictionsType.value

      saveWebAppInfo(formData)
    })
    .catch((error) => {
      console.log('error', error)
    })
}

const handlePreview = () => {
  window.open(h5_website)
}

const copyH5WebSite = () => {
  copyText(h5_website)
  message.success('复制成功')
}

const generateQR = async () => {
  try {
    previewQrcode.value = await QRCode.toDataURL(h5_website)
  } catch (err) {
    console.error(err)
  }
}

generateQR()
</script>
