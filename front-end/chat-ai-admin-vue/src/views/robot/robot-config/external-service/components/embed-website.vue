<style lang="less" scoped>
.web-app-box {
  display: flex;

  .box-left {
    flex: 1;

    .box-wrapper {
      margin-bottom: 16px;
      &:last-child {
        margin-bottom: 0;
      }
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

  .form-box {
    .form-item {
      margin-bottom: 16px;
    }
  }
}
.box-right {
  width: 343px;
  margin: 0 96px 0 48px;

  .preview-img {
    display: block;
    border-radius: 9px;
    box-shadow: 0 4px 32px 0 rgba(0, 0, 0, 0.16);
  }
  iframe {
    width: 375px;
    height: 720px;
    border-radius: 4px;
    box-shadow: 0 4px 32px 0 rgba(0, 0, 0, 0.16);
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
        <card-box :title="t('style_settings_title')">
          <template #icon>
            <svg-icon name="pc" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <template #action>
            <a-button @click="saveForm" size="small" type="primary">{{ t('save_btn') }}</a-button>
          </template>
          <div class="web-app-style form-box">
            <a-form ref="formRef" layout="vertical" :model="formState" :rules="formRules">
              <a-form-item class="form-item" :label="t('page_title_label')" name="headTitle">
                <PageTitleInput
                  v-model:avatar="formState.headImage"
                  v-model:value="formState.headTitle"
                />
              </a-form-item>
              <a-form-item class="form-item" :label="t('subtitle_label')" name="headSubTitle">
                <a-textarea
                  v-model:value="formState.headSubTitle"
                  :placeholder="t('subtitle_placeholder')"
                />
              </a-form-item>
              <a-form-item
                class="form-item"
                :label="t('color_label')"
                :name="['pageStyle', 'headBackgroundColor']"
              >
                <GradientColorPicker v-model:value="formState.pageStyle.headBackgroundColor" />
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
              <a-form-item class="form-item" :label="t('label_show_history_and_new_session_btn')" name="new_session_btn_show" required>
                <a-radio-group v-model:value="formState.new_session_btn_show" name="new_session_btn_show">
                  <a-radio :value="1">{{ t('label_show') }}</a-radio>
                  <a-radio :value="2">{{ t('label_hide') }}</a-radio>
                </a-radio-group>
              </a-form-item>
            </a-form>
          </div>
        </card-box>
      </div>

      <div class="box-wrapper">
        <card-box :title="t('copy_sdk_code_title')">
          <template #icon>
            <svg-icon name="sdk" style="font-size: 16px; color: #262626"></svg-icon>
          </template>
          <template #action>
            <a-button @click="copySDKCode" size="small">{{ t('copy_btn') }}</a-button>
          </template>
          <div class="sdk-code">
            <pre><code>{{ sdkCode }}</code></pre>
          </div>
        </card-box>
      </div>

      <div class="box-wrapper">
        <QuickInstruction
          :type="robotInfo.app_id_embed"
          @updata="updataQuickComand"
        ></QuickInstruction>
      </div>

      <div class="box-wrapper">
        <FloatIconSetting :form="formState" @save="handleFloatBtnCongiSave"></FloatIconSetting>
      </div>
    </div>
    <div class="box-right">
      <div class="demo-box">
        <iframe id="web-preview" :src="previewIframeSrc" frameborder="0"></iframe>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'
import { getSdkCode } from './sdk-code'
import { ref, reactive, toRaw, watch, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useRobotStore } from '@/stores/modules/robot'
import { message } from 'ant-design-vue'
import { copyText } from '@/utils/index'
import { editExternalConfig } from '@/api/robot/index'
import CardBox from './card-box.vue'
import PageTitleInput from './page-title-input.vue'
import GradientColorPicker from './gradient-color-picker.vue'
import QuickInstruction from './quick-instruction.vue'
import FloatIconSetting from './float-icon-setting.vue'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'

const { t } = useI18n('views.robot.robot-config.external-service.components.embed-website')
const robotStore = useRobotStore()
const { robotInfo, external_config_pc } = storeToRefs(robotStore)
const { getRobot } = robotStore

const sdkCode = getSdkCode(robotInfo.value)

const copySDKCode = () => {
  copyText(sdkCode)
  message.success(t('copy_success'))
}

const formRef = ref()
const formState = reactive({
  headTitle: external_config_pc.value.headTitle,
  headSubTitle: external_config_pc.value.headSubTitle,
  headImage: external_config_pc.value.headImage,
  lang: external_config_pc.value.lang,
  pageStyle: external_config_pc.value.pageStyle,
  floatBtn: external_config_pc.value.floatBtn,
  open_type: external_config_pc.value.open_type,
  window_width: external_config_pc.value.window_width,
  window_height: external_config_pc.value.window_height,
  new_session_btn_show: external_config_pc.value.new_session_btn_show
})

const previewIframeSrc = computed(() => {
  let { pc_domain, robot_key } = robotInfo.value
  return `${pc_domain}/web/#/chat?robot_key=${robot_key}&language=${formState.lang || 'zh-CN'}`
})

watch(formState, (val) => {
  updatePreview(val)
})

const updataQuickComand = (data) => {
  updatePreview(data, 'updataQuickComand')
}

const updatePreview = (data, type) => {
  let iframe = document.getElementById('web-preview')
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
  headTitle: [
    {
      required: true,
      message: t('please_input_title'),
      trigger: 'input'
    },
    {
      trigger: 'input',
      validator: () => {
        if (!formState.headImage) {
          return Promise.reject(t('please_upload_logo'))
        } else {
          return Promise.resolve()
        }
      }
    }
  ],
  pageStyle: {
    headBackgroundColor: [
      {
        required: true,
        message: t('please_select_color'),
        trigger: 'change'
      }
    ]
  }
}

// 保存样式设置
const saveWebSiteInfo = () => {
  const { id } = robotInfo.value
  formState.window_width = +formState.window_width || 1200
  formState.window_height = +formState.window_height || 650
  let formData = { ...toRaw(formState) }

  editExternalConfig({
    id: id,
    external_config_pc: JSON.stringify(formData)
  }).then(() => {
    getRobot(id)
    message.success(t('save_success'))
  })
}

const saveForm = () => {
  formRef.value
    .validate()
    .then(() => {
      saveWebSiteInfo()
    })
    .catch((error) => {
      console.log('error', error)
    })
}

// 保存浮标设置
const handleFloatBtnCongiSave = (data) => {
  console.log(data, 'handleFloatBtnCongiSave')
  formState.floatBtn = { ...data }
  saveForm()
}
</script>
