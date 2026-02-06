<template>
  <a-modal
    v-model:open="configOpen"
    :confirm-loading="configSaving"
    :title="t('title_api_key_auth_config')"
    width="620px"
    @ok="saveConfig"
  >
    <a-alert type="info" show-icon :message="t('msg_config_workflow_call')"/>
    <a-form
      class="mt16"
      :model="config"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
    >
      <a-form-item
        :label="t('label_credential_name')"
        name="name"
        :rules="[{ required: true, message: t('msg_input_credential_name') }]"
      >
        <a-input v-model:value.trim="config.name" :maxlength="16" :placeholder="t('ph_input_credential_name')"/>
        <div class="tip-desc">{{ t('msg_credential_name_desc') }}</div>
      </a-form-item>
      <a-form-item
        :label="t('label_app_id')"
        name="appid"
        :rules="[{ required: true, message: t('msg_input_app_id') }]"
      >
        <a-input v-model:value.trim="config.appid" :placeholder="t('ph_input_app_id')"/>
      </a-form-item>
      <a-form-item
        :label="t('label_app_secret')"
        name="app_secret"
        :rules="[{ required: true, message: t('msg_input_app_secret') }]"
      >
        <a-input v-model:value.trim="config.app_secret" :placeholder="t('ph_input_app_secret')"/>
        <div class="tip-desc">{{ t('msg_how_to_get_app_id_secret') }}</div>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue';
import {setPluginConfig} from "@/api/plugins/index.js";
import {message} from 'ant-design-vue';
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.components.feishu-config-modal')
const emit = defineEmits(['change'])
const configData = ref({})
const configOpen = ref(false)
const configSaving = ref(false)
const config = reactive({})

function show(cData = {}) {
  configData.value = cData
  Object.assign(config, {
    name: '',
    appid: '',
    app_secret: '',
    is_default: false
  })
  // 首次添加设置默认
  if (!configData.value || !Object.keys(configData.value).length) config.is_default = true
  configOpen.value = true
}

function saveConfig() {
  try {
    configSaving.value = true
    if (!config.name) throw t('msg_input_credential_name')
    if (configData.value[config.name]) throw t('msg_credential_name_exists')
    if (!config.appid) throw t('msg_input_app_id')
    if (!config.app_secret) throw t('msg_input_app_secret')
    setPluginConfig({
      name: 'feishu_bitable',
      data: JSON.stringify({
        ...configData.value,
        [config.name]: config
      })
    }).then(res => {
      emit('change')
      configOpen.value = false
      message.success(t('msg_saved'))
    }).finally(() => {
      configSaving.value = false
    })
  } catch (e) {
    console.log('Err:', e)
    configSaving.value = false
    message.error(e)
  }
}

defineExpose({
  show
})
</script>

<style scoped lang="less">
.tip-desc {
  color: #8C8C8C;
  margin-top: 4px;
}

.mt16 {
  margin-top: 16px;
}
</style>
