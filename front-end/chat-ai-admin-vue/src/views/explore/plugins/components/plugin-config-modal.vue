<template>
  <a-modal
    v-model:open="configOpen"
    :confirm-loading="configSaving"
    :title="t('title_api_key_auth_config')"
    width="620px"
    @ok="saveConfig"
  >
    <a-alert type="info" show-icon :message="t('msg_config_workflow_tip')"/>
    <a-form
      class="mt16"
      :model="config"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
    >
      <a-form-item
        :label="t('label_credential_name')"
        name="name"
        :rules="[{ required: true, message: t('msg_please_input_credential_name') }]"
      >
        <a-input v-model:value.trim="config.name" :maxlength="16" :placeholder="t('ph_input_credential_name')"/>
        <div class="tip-desc">{{ t('msg_credential_name_desc') }}</div>
      </a-form-item>
      <template v-for="field in fields" :key="field.key">
        <a-form-item
          :label="field.name || field.key"
          :name="field.key"
          :rules="field.required ? [{ required: true, message: t('msg_please_input_field', { field: field.name || field.key }) }] : []"
        >
          <div style="display:flex; align-items:center; gap:8px;">
            <a-input v-model:value.trim="config[field.key]" :placeholder="field.placeholder || ''"/>
            <a style="flex-shrink:0;" v-if="field.click_url" :href="field.click_url" target="_blank" rel="noopener noreferrer">
              {{ field.click_title || t('btn_get') }}
            </a>
          </div>
          <div class="tip-desc" v-if="field.tip">{{ field.tip }}</div>
          <div class="tip-desc" v-else-if="field.desc">{{ field.desc }}</div>
        </a-form-item>
      </template>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue';
import {setPluginConfig} from "@/api/plugins/index.js";
import {message} from 'ant-design-vue';
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.explore.plugins.components.plugin-config-modal');

const emit = defineEmits(['change'])
const configData = ref({})
const pluginName = ref('')
const schemaData = ref({})
const configOpen = ref(false)
const configSaving = ref(false)
const config = reactive({})
const fields = ref([])

function sanitize(s) {
  return String(s || '').replace(/`/g, '').trim()
}
function buildFields(schema) {
  const list = []
  Object.keys(schema || {}).forEach((act) => {
    const meta = schema[act] || {}
    if (meta && meta.use_plugin_config === true) {
      const params = meta.params || {}
      Object.keys(params).forEach((k) => {
        const p = params[k] || {}
        if (p.plugin_config_component === true) {
          list.push({
            key: k,
            name: p.name,
            desc: p.desc,
            tip: p.tip,
            required: !!p.required,
            placeholder: p.placeholder,
            click_title: sanitize(p.click_title),
            click_url: sanitize(p.click_url)
          })
        }
      })
    }
  })
  // 去重
  const seen = new Set()
  fields.value = list.filter((f) => {
    if (seen.has(f.key)) return false
    seen.add(f.key)
    return true
  }).sort((a, b) => (+(schema?.[a.key]?.sort_num || 0)) - (+(schema?.[b.key]?.sort_num || 0)))
}

function show(cData = {}, name = '', schema = {}) {
  configData.value = cData
  pluginName.value = name
  schemaData.value = schema
  Object.keys(config).forEach((k) => delete config[k])
  Object.assign(config, {
    name: '',
    is_default: false
  })
  if (!configData.value || !Object.keys(configData.value).length) config.is_default = true
  buildFields(schemaData.value)
  fields.value.forEach((f) => {
    config[f.key] = ''
  })
  configOpen.value = true
}

function saveConfig() {
  try {
    configSaving.value = true
    if (!config.name) throw t('msg_please_input_credential_name')
    if (configData.value[config.name]) throw t('msg_credential_name_exists')
    // 校验必填字段
    fields.value.forEach((f) => {
      const v = String(config[f.key] || '').trim()
      if (f.required && !v) {
        throw t('msg_please_input_field', { field: f.name || f.key })
      }
    })
    setPluginConfig({
      name: pluginName.value,
      data: JSON.stringify({
        ...configData.value,
        [config.name]: config
      })
    }).then(() => {
      emit('change')
      configOpen.value = false
      message.success(t('msg_saved'))
    }).finally(() => {
      configSaving.value = false
    })
  } catch (e) {
    configSaving.value = false
    message.error(String(e || t('msg_config_save_failed')))
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
