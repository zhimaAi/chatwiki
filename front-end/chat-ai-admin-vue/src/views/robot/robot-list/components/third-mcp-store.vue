<template>
  <a-modal
    v-model:open="visible"
    :zIndex="1002"
    :title="formState.provider_id ? t('title_edit_mcp') : t('title_add_mcp')"
    width="472px"
    :confirm-loading="saving"
    @ok="save">
    <div v-if="checking" class="checking-box">
      <div class="cont">
        <a-spin/>
        {{ t('msg_checking') }}
      </div>
    </div>
    <div class="avatar-box">
      <AvatarInput v-model:value="formState.avatar" @change="avatarChange"/>
      <div class="tip-info">{{ t('tip_avatar') }}</div>
    </div>
    <a-form class="form-box" labelAlign="left">
      <a-form-item :label="t('label_mcp_plugin_name')" required :colon="false">
        <a-input v-model:value="formState.name" :placeholder="t('ph_mcp_plugin_name')" :maxlength="20"/>
      </a-form-item>
      <a-form-item :label="t('label_description')" :colon="false">
        <a-textarea
          v-model:value="formState.description"
          :auto-size="{ minRows: 2, maxRows: 5 }"
          :placeholder="t('ph_description')"
          :maxlength="120"
        />
      </a-form-item>
      <div class="tit-box">{{ t('title_config_info') }}</div>
      <a-form-item :label="t('label_plugin_url')" :colon="false">
        <a-input v-model:value="formState.url" :placeholder="t('ph_plugin_url')"/>
        <div v-if="formState.provider_id" class="cFB363F">{{ t('msg_url_warning') }}</div>
      </a-form-item>
      <a-form-item :label="t('label_timeout')" :colon="false">
        <a-input-number v-model:value="formState.request_timeout" :precision="0" style="width: 100%;"
                        :placeholder="t('ph_timeout')"/>
      </a-form-item>
      <a-form-item :label="t('label_headers')" :colon="false">
        <div class="tip-info">{{ t('tip_headers') }}</div>
        <div class="req-head-box">
          <div class="req-head-item">
            <div class="tit">{{ t('label_header_name') }}</div>
            <div class="tit">{{ t('label_header_value') }}</div>
          </div>
          <div v-for="(item, i) in formState.headers" :key="i" class="req-head-item">
            <a-input v-model:value="item.key" :placeholder="t('ph_input')"/>
            <a-input v-model:value="item.value" :placeholder="t('ph_input')"/>
            <CloseCircleOutlined @click="delHeader(i)"/>
          </div>
          <a-button class="add-btn" type="dashed" :icon="h(PlusOutlined)" @click="addHeader">{{ t('btn_add_header') }}</a-button>
        </div>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive, h} from 'vue';
import {message} from 'ant-design-vue';
import {CloseCircleOutlined, PlusOutlined} from '@ant-design/icons-vue';
import AvatarInput from "@/views/robot/robot-list/components/avatar-input.vue";
import {DEFAULT_MCP_AVATAR} from "@/constants/index.js";
import {authTMcpProvider, saveTMcpProvider} from "@/api/robot/thirdMcp.js";
import {setShowReqError} from "@/utils/http/axios/config.js";
import {useI18n} from '@/hooks/web/useI18n';

const {t} = useI18n('views.robot.robot-list.components.third-mcp-store');

const emit = defineEmits(['ok'])

const providerId = ref(0)
const visible = ref(false)
const checking = ref(false)
const saving = ref(false)
const avatarData = ref(null)
const authDataBack = ref('')
const formStateStruct = {
  avatar: DEFAULT_MCP_AVATAR,
  name: '',
  description: '',
  provider_id: '',
  url: '',
  request_timeout: '30',
  headers: []
}
const formState = reactive({})

function show(info = null) {
  Object.assign(formState, info || formStateStruct)
  avatarData.value = null
  formState.headers = []
  if (info) {
    providerId.value = info.id
    formState.provider_id = info.id
    let headers = JSON.parse(info.headers)
    for (let key in headers) {
      formState.headers.push({key, value: headers[key]})
    }
  }
  authDataBack.value = JSON.stringify({headers: formState.headers, url: formState.url})
  visible.value = true
}

function avatarChange(data) {
  avatarData.value = data
}

function addHeader() {
  formState.headers.push({
    key: '',
    value: ''
  })
}

function delHeader(index) {
  formState.headers.splice(index, 1)
}

function hasAuthDataChange() {
  return authDataBack.value != JSON.stringify({headers: formState.headers, url: formState.url})
}

function authConfig() {
  const hide = message.loading(t('msg_saving_auth'), 0)
  checking.value = true
  setTimeout(() => {
    setShowReqError(false)
    authTMcpProvider({provider_id: providerId.value}).then(res => {
      message.success(t('msg_auth_success'))
    }).catch(err => {
      message.warning(t('msg_auth_failed') + err.message)
    }).finally(() => {
      hide()
      checking.value = false
      visible.value = false
      emit('ok', formState)
      setTimeout(() => {
        setShowReqError(true)
      }, 500)
    })
  }, 1800)
}

function save() {
  try {
    saving.value = true
    formState.name = formState.name.trim()
    formState.description = formState.description.trim()
    if (!formState.name) throw t('err_plugin_name_required')
    if (!formState.description) throw t('err_plugin_desc_required')
    if (!formState.url) throw t('err_url_required')
    if (!formState.request_timeout) throw t('err_timeout_required')
    let data = {...formState}
    if (avatarData.value) data.avatar = avatarData.value?.file
    data.headers = {}
    for (let item of formState.headers) {
      item.key = item.key.trim()
      item.value = item.value.trim()
      if (item.key && item.value) {
        data.headers[item.key] = item.value
      }
    }
    data.headers = JSON.stringify(data.headers)
    delete data.toools
    saveTMcpProvider(data).then(res => {
      let {provider_id} = res?.data || {}
      providerId.value = provider_id
      if (hasAuthDataChange()) {
        authConfig()
      } else {
        emit('ok', formState)
        message.success(t('msg_saved'))
        visible.value = false
      }
    }).finally(() => {
      saving.value = false
    })
  } catch (e) {
    saving.value = false
    message.error(e)
  }
}

defineExpose({
  show,
})
</script>

<style scoped lang="less">
.checking-box {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 99;

  .cont {
    padding: 40px 80px;
    border-radius: 6px;
    background: #FFF;
    box-shadow: 0 4px 16px 0 #1b3a6929;
    display: flex;
    gap: 8px;
    font-size: 16px;
    font-weight: 500;
  }
}

.avatar-box {
  display: flex;
  flex-direction: column;
  align-items: center;

  :deep(.ant-upload-select) {
    width: 62px !important;
    height: 62px !important;
    border: none !important;
    border-radius: 8px !important;
    overflow: hidden;
  }

  .tip-info {
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
    margin-top: 8px;
  }
}

.form-box {
  margin-top: 24px;

  :deep(.ant-form-item) {
    margin-bottom: 12px;

    .ant-row {
      display: block;

      .ant-form-item-control-input {
        min-height: unset;
      }
    }
  }

  .tit-box {
    color: #262626;
    font-size: 14px;
    font-weight: 600;
    margin: 16px 0 8px;

    .desc {
      color: #595959;
      font-weight: 400;
      margin-left: 12px;
    }
  }

  .tip-info {
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
  }

  .req-head-box {
    color: #595959;
    font-size: 14px;
    font-weight: 400;
    margin-top: 8px;

    .req-head-item {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-bottom: 8px;

      &:first-child {
        margin-top: 4px;
      }

      .tit {
        width: 198px;
      }
    }

    .add-btn {
      width: 100%;
    }
  }
}

.cFB363F {
  color: #FB363F;
}
</style>
