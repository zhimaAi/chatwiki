<template>
  <a-modal
    v-model:open="configOpen"
    :confirm-loading="configSaving"
    title="API Key授权配置"
    width="620px"
    @ok="saveConfig"
  >
    <a-alert type="info" show-icon message="配置凭据后，工作流可直接调用此插件"/>
    <a-form
      class="mt16"
      :model="config"
      :label-col="{ span: 4 }"
      :wrapper-col="{ span: 20 }"
    >
      <a-form-item
        label="凭证名称"
        name="name"
        :rules="[{ required: true, message: '请输入凭证名称!' }]"
      >
        <a-input v-model:value.trim="config.name" :maxlength="16" placeholder="请输入凭证名称"/>
        <div class="tip-desc">仅用于区分多个授权配置的名称，不用于接口调用</div>
      </a-form-item>
      <a-form-item
        label="APP ID"
        name="appid"
        :rules="[{ required: true, message: '请输入APP ID!' }]"
      >
        <a-input v-model:value.trim="config.appid" placeholder="请输入APP ID"/>
      </a-form-item>
      <a-form-item
        label="APP Secret"
        name="app_secret"
        :rules="[{ required: true, message: '请输入APP Secret!' }]"
      >
        <a-input v-model:value.trim="config.app_secret" placeholder="请输入APP Secret"/>
        <div class="tip-desc">如何获取APP ID和APP Secret</div>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue';
import {setPluginConfig} from "@/api/plugins/index.js";
import {message} from 'ant-design-vue';

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
    if (!config.name) throw '请输入凭证名称'
    if (configData.value[config.name]) throw '凭证名称已存在'
    if (!config.appid) throw '请输入APP ID'
    if (!config.app_secret) throw '请输入APP Secret'
    setPluginConfig({
      name: 'feishu_bitable',
      data: JSON.stringify({
        ...configData.value,
        [config.name]: config
      })
    }).then(res => {
      emit('change')
      configOpen.value = false
      message.success('已保存')
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
