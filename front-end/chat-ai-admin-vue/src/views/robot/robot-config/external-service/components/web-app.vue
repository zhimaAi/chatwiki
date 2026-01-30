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
  .ml4 {
    margin-left: 4px;
  }
}
.box-right {
  margin: 0 96px 0 48px;
  .demo-box {
    position: relative;
  }
  iframe {
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

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 16px 0;
  height: 24px;
  .title-text {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
}

.window-size-box{
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 14px;
}

</style>

<template>
  <div class="web-app-box">
    <div class="box-left">
      <div class="box-wrapper">
        <card-box title="H5链接">
          <template #icon>
            <svg-icon name="phone" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <div class="web-app-info">
            <div class="web-app-link">
              <a :href="h5_src" target="_blank">{{ h5_src }}</a>
            </div>
            <div class="link-action">
              <a-button class="action-btn" type="primary" ghost @click="copyH5WebSite(h5_src)"
                >复 制</a-button
              >
              <a-button class="action-btn" @click="handlePreview(h5_src)">预 览</a-button>
              <a-tooltip color="#fff" placement="top">
                <template #title>
                  <img style="width: 180px" :src="previewQrcodeH5" alt="" />
                </template>
                <a-button class="action-btn">二维码</a-button>
              </a-tooltip>
            </div>

            <div class="card-title">
              <svg-icon name="circularNeedle" style="font-size: 16px; color: #262626"></svg-icon>
              <div class="title-text">pc网页链接</div>
            </div>

            <div class="web-app-link">
              <a :href="pc_src" target="_blank">{{ pc_src }}</a>
            </div>
            <div class="link-action">
              <a-button class="action-btn" type="primary" ghost @click="copyH5WebSite(pc_src)"
                >复 制</a-button
              >
              <a-button class="action-btn" @click="handlePreview(pc_src)">预 览</a-button>
              <a-tooltip color="#fff" placement="top">
                <template #title>
                  <img style="width: 180px" :src="previewQrcodePc" alt="" />
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
                    <a-radio v-if="false" :value="3"
                      ><span class="default-text-color">有权限的账号登录后才可访问</span>
                      <a-tooltip title="可在基础配置→ 权限管理处添加协作者。">
                        <QuestionCircleOutlined class="ml4" />
                      </a-tooltip>
                    </a-radio>
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

              <a-form-item class="form-item" label="是否显示历史对话和新增对话按钮" name="new_session_btn_show">
                <a-radio-group v-model:value="formState.new_session_btn_show" name="new_session_btn_show">
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
                  <a-select-option value="en-US">English</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item class="form-item" label="网址打开方式" name="open_type" required>
                <a-radio-group v-model:value="formState.open_type">
                  <a-radio :value="1">新标签页打开</a-radio>
                  <a-radio :value="2">新窗口弹窗打开
                    <a-tooltip title="仅管控PC端新窗口打开,移动端依然用新标签页打开">
                      <template #title>prompt text</template>
                      <QuestionCircleOutlined />
                    </a-tooltip>
                  </a-radio>
                </a-radio-group>
                <a-form-item-rest v-if="formState.open_type == 2">
                  <div class="window-size-box">
                    <a-flex align="center" :gap="8">
                      <div>弹窗高度</div>
                      <a-input-number  v-model:value="formState.window_height" :min="500" :max="2000" />
                      PX
                    </a-flex>
                    <a-flex align="center" :gap="8">
                      <div>弹窗宽度</div>
                      <a-input-number  v-model:value="formState.window_width" :min="500" :max="2000" />
                      PX
                    </a-flex>
                  </div>
                </a-form-item-rest>
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
import { QuestionCircleOutlined } from '@ant-design/icons-vue'

const router = useRouter()
const robotStore = useRobotStore()
const { robotInfo, external_config_h5 } = storeToRefs(robotStore)
const { h5_website, h5_domain, robot_key } = robotInfo.value
const { getRobot } = robotStore

const previewQrcodePc = ref('')
const previewQrcodeH5 = ref('')

const formRef = ref()
const formState = reactive({
  logo: external_config_h5.value.logo,
  pageTitle: external_config_h5.value.pageTitle,
  navbarShow: external_config_h5.value.navbarShow,
  lang: external_config_h5.value.lang,
  pageStyle: external_config_h5.value.pageStyle,
  open_type: external_config_h5.value.open_type,
  window_width: external_config_h5.value.window_width,
  window_height: external_config_h5.value.window_height,
  new_session_btn_show: external_config_h5.value.new_session_btn_show
})

const pc_src = computed(() => {
  return `${h5_domain}/#/chat/pc?robot_key=${robot_key}`
})

const h5_src = computed(() => {
  return `${h5_domain}/#/chat/h5?robot_key=${robot_key}`
})

const previewIframeSrc = computed(() => {
  return h5_website
})
watch(formState, (val) => {
  updatePreview(val)
})
const updataQuickComand = (data) => {
  updatePreview(data, 'updataQuickComand')
}
const updatePreview = (data, type) => {
  let iframe = document.getElementById('mobile-preview')
  iframe.contentWindow.postMessage(
    {
      type: type || 'onPreview',
      data: JSON.parse(JSON.stringify(data))
    },
    '*'
  )
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

      formData.window_width = +formData.window_width || 1200
      formData.window_height = +formData.window_height || 650

      formData.accessRestrictionsType = accessRestrictionsType.value

      saveWebAppInfo(formData)
    })
    .catch((error) => {
      console.log('error', error)
    })
}

const handlePreview = (src) => {
  window.open(src)
}

const copyH5WebSite = (text) => {
  copyText(text)
  message.success('复制成功')
}

const generateQR = async () => {
  try {
    previewQrcodePc.value = await QRCode.toDataURL(pc_src.value)
    previewQrcodeH5.value = await QRCode.toDataURL(h5_src.value)
  } catch (err) {
    console.error(err)
  }
}

setTimeout(() => {
  generateQR()
}, 300)
</script>
