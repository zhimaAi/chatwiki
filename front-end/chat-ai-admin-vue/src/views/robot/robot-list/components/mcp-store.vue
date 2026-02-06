<template>
  <a-modal
    v-model:open="visible"
    :title="formState.server_id ? t('title_edit_mcp') : t('title_add_mcp')"
    width="472px"
    :confirm-loading="saving"
    @ok="save">
    <div class="avatar-box">
      <AvatarInput v-model:value="formState.avatar" @change="avatarChange"/>
      <div class="tip-info">{{ t('tip_avatar_replace') }}</div>
    </div>
    <a-form class="form-box" labelAlign="left">
      <a-form-item :label="t('label_mcp_plugin_name')" required :colon="false">
        <a-input v-model:value="formState.name" :placeholder="t('ph_mcp_plugin_name')" :maxlength="20"/>
      </a-form-item>
      <a-form-item :label="t('label_description')" :colon="false">
        <a-textarea
          v-model:value="formState.description"
          :auto-size="{ minRows: 2, maxRows: 5 }"
          :placeholder="t('ph_description')" :maxlength="20"/>
      </a-form-item>
      <template v-if="!formState.server_id">
        <div class="tit-box">{{ t('title_config_info') }} <span class="desc">{{ t('desc_config_auto_generate') }}</span></div>
        <a-form-item :label="t('label_plugin_url')" :colon="false">
          <a-input disabled :placeholder="t('ph_plugin_url_default')"/>
        </a-form-item>
        <a-form-item :label="t('label_auth_method')" :colon="false">
          <a-input disabled :placeholder="t('ph_auth_method_default')"/>
        </a-form-item>
        <a-form-item :label="t('label_request_header')" :colon="false">
          <a-input disabled :placeholder="t('ph_request_header_default')"/>
        </a-form-item>
        <a-form-item :label="t('label_api_key')" :colon="false">
          <a-input disabled :placeholder="t('ph_api_key_default')"/>
        </a-form-item>
      </template>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue';
import {message} from 'ant-design-vue';
import AvatarInput from "@/views/robot/robot-list/components/avatar-input.vue";
import {saveMcpServer} from "@/api/robot/mcp.js";
import {DEFAULT_MCP_AVATAR} from "@/constants/index.js";
import {base64ToFile} from "@/utils/index.js";
import {useI18n} from '@/hooks/web/useI18n';

const { t } = useI18n('views.robot.robot-list.components.mcp-store');

const emit = defineEmits(['ok'])

const visible = ref(false)
const saving = ref(false)
const avatarData = ref(null)
const formStateStruct = {
  avatar: DEFAULT_MCP_AVATAR,
  name: '',
  description: '',
  server_id: '',
}
const formState = reactive({})

function show(info = null) {
  Object.assign(formState, info || formStateStruct)
  avatarData.value = null
  visible.value = true
}

function avatarChange(data) {
  avatarData.value = data
}

function save() {
  try {
    saving.value = true
    formState.name = formState.name.trim()
    formState.description = formState.description.trim()
    if (!formState.name) throw t('msg_input_mcp_plugin_name')
    if (!formState.description) throw t('msg_input_mcp_plugin_description')
    let data = {...formState}
    if (avatarData.value) data.avatar = avatarData.value?.file
    saveMcpServer(data).then(res => {
      emit('ok', formState)
      message.success(t('msg_saved'))
      visible.value = false
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
}
</style>
