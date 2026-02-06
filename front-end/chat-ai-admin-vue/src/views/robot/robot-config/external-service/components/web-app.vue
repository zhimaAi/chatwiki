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
        <card-box :title="t('h5_link_title')">
          <template #icon>
            <svg-icon name="phone" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <div class="web-app-info">
            <div class="web-app-link">
              <a :href="h5_src" target="_blank">{{ h5_src }}</a>
            </div>
            <div class="link-action">
              <a-button class="action-btn" type="primary" ghost @click="copyH5WebSite(h5_src)"
                >{{ t('copy_btn') }}</a-button
              >
              <a-button class="action-btn" @click="handlePreview(h5_src)">{{ t('preview_btn') }}</a-button>
              <a-tooltip color="#fff" placement="top">
                <template #title>
                  <img style="width: 180px" :src="previewQrcodeH5" alt="" />
                </template>
                <a-button class="action-btn">{{ t('qrcode_btn') }}</a-button>
              </a-tooltip>
            </div>

            <div class="card-title">
              <svg-icon name="circularNeedle" style="font-size: 16px; color: #262626"></svg-icon>
              <div class="title-text">{{ t('pc_web_link_title') }}</div>
            </div>

            <div class="web-app-link">
              <a :href="pc_src" target="_blank">{{ pc_src }}</a>
            </div>
            <div class="link-action">
              <a-button class="action-btn" type="primary" ghost @click="copyH5WebSite(pc_src)"
                >{{ t('copy_btn') }}</a-button
              >
              <a-button class="action-btn" @click="handlePreview(pc_src)">{{ t('preview_btn') }}</a-button>
              <a-tooltip color="#fff" placement="top">
                <template #title>
                  <img style="width: 180px" :src="previewQrcodePc" alt="" />
                </template>
                <a-button class="action-btn">{{ t('qrcode_btn') }}</a-button>
              </a-tooltip>
            </div>

            <div class="access-restrictions form-box">
              <div class="form-item">
                <div class="form-item-label">{{ t('access_restrictions_label') }}</div>
                <div class="form-item-body">
                  <a-radio-group
                    v-model:value="accessRestrictionsType"
                    @change="saveAccessRestrictionsType"
                  >
                    <a-radio :value="1"><span class="default-text-color">{{ t('no_restriction') }}</span></a-radio>
                    <a-radio :value="2"
                      ><span class="default-text-color">{{ t('login_required') }}</span></a-radio
                    >
                    <a-radio v-if="false" :value="3"
                      ><span class="default-text-color">{{ t('authorized_login_required') }}</span>
                      <a-tooltip :title="t('authorized_tooltip')">
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
        <card-box :title="t('style_settings_title')">
          <template #icon>
            <svg-icon name="phone" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <template #action>
            <a-button @click="saveForm" size="small" type="primary">{{ t('save_btn') }}</a-button>
          </template>
          <div class="web-app-style form-box">
            <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
              <a-form-item class="form-item" :label="t('show_navbar_label')" name="navbarShow">
                <a-radio-group v-model:value="formState.navbarShow" name="navbarShow">
                  <a-radio :value="1">{{ t('show_navbar_yes') }}</a-radio>
                  <a-radio :value="2">{{ t('show_navbar_no') }}</a-radio>
                </a-radio-group>
              </a-form-item>

              <a-form-item class="form-item" :label="t('show_history_and_new_session_btn_label')" name="new_session_btn_show">
                <a-radio-group v-model:value="formState.new_session_btn_show" name="new_session_btn_show">
                  <a-radio :value="1">{{ t('show_navbar_yes') }}</a-radio>
                  <a-radio :value="2">{{ t('show_navbar_no') }}</a-radio>
                </a-radio-group>
              </a-form-item>

              <a-form-item class="form-item" :label="t('page_title_label')" name="pageTitle">
                <PageTitleInput
                  v-model:avatar="formState.logo"
                  v-model:value="formState.pageTitle"
                />
              </a-form-item>

              <a-form-item
                class="form-item"
                :label="t('navbar_color_label')"
                :name="['pageStyle', 'navbarBackgroundColor']"
              >
                <ColorPicker v-model:value="formState.pageStyle.navbarBackgroundColor" />
              </a-form-item>

              <a-form-item class="form-item" :label="t('language_label')" name="lang">
                <a-select
                  style="width: 180px"
                  v-model:value="formState.lang"
                  :placeholder="t('language_placeholder')"
                >
                  <a-select-option value="zh-CN">{{ t('language_zh_cn') }}</a-select-option>
                  <a-select-option value="en-US">{{ t('language_en_us') }}</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item class="form-item" :label="t('url_open_type_label')" name="open_type" required>
                <a-radio-group v-model:value="formState.open_type">
                  <a-radio :value="1">{{ t('open_new_tab') }}</a-radio>
                  <a-radio :value="2">{{ t('open_new_window') }}
                    <a-tooltip :title="t('open_new_window_tooltip')">
                      <template #title>prompt text</template>
                      <QuestionCircleOutlined />
                    </a-tooltip>
                  </a-radio>
                </a-radio-group>
                <a-form-item-rest v-if="formState.open_type == 2">
                  <div class="window-size-box">
                    <a-flex align="center" :gap="8">
                      <div>{{ t('window_height_label') }}</div>
                      <a-input-number  v-model:value="formState.window_height" :min="500" :max="2000" />
                      PX
                    </a-flex>
                    <a-flex align="center" :gap="8">
                      <div>{{ t('window_width_label') }}</div>
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
import { useI18n } from '@/hooks/web/useI18n'
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

const { t } = useI18n('views.robot.robot-config.external-service.components.web-app')
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
      message: t('please_select_language'),
      trigger: 'change'
    }
  ],
  pageTitle: [
    {
      required: true,
      message: t('please_input_title'),
      trigger: 'input'
    },
    {
      trigger: 'input',
      validator: () => {
        if (!formState.logo) {
          return Promise.reject(t('please_upload_logo'))
        } else {
          return Promise.resolve()
        }
      }
    }
  ],
  navbarShow: [
    {
      required: true,
      message: t('please_select_navbar_show'),
      trigger: 'change'
    }
  ],
  pageStyle: {
    navbarBackgroundColor: [
      {
        required: true,
        message: t('please_select_navbar_color'),
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
    message.success(t('save_success'))
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
  message.success(t('copy_success'))
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
