<template>
  <a-modal
    v-model:open="visible"
    :title="formState.server_id ? '编辑MCP' : '添加MCP'"
    width="472px"
    :confirm-loading="saving"
    @ok="save">
    <div class="avatar-box">
      <AvatarInput v-model:value="formState.avatar" @change="avatarChange"/>
      <div class="tip-info">点击替换，建议尺寸为100*100px，大小不超过100kb</div>
    </div>
    <a-form class="form-box" labelAlign="left">
      <a-form-item label="MCP插件名称" required :colon="false">
        <a-input v-model:value="formState.name" placeholder="请输入MCP插件名称，最多20个字" :maxlength="20"/>
      </a-form-item>
      <a-form-item label="描述" :colon="false">
        <a-textarea
          v-model:value="formState.description"
          :auto-size="{ minRows: 2, maxRows: 5 }"
          placeholder="请输入描述" :maxlength="20"/>
      </a-form-item>
      <template v-if="!formState.server_id">
        <div class="tit-box">配置信息 <span class="desc">以下配置创建完成后自动生成</span></div>
        <a-form-item label="插件URL" :colon="false">
          <a-input disabled placeholder="创建完成后默认生成"/>
        </a-form-item>
        <a-form-item label="授权方式" :colon="false">
          <a-input disabled placeholder="Service token / API key"/>
        </a-form-item>
        <a-form-item label="请求头" :colon="false">
          <a-input disabled placeholder="Authorization"/>
        </a-form-item>
        <a-form-item label="API KEY" :colon="false">
          <a-input disabled placeholder="创建完成后默认生成"/>
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
    if (!formState.name) throw '请输入MCP插件名称'
    if (!formState.description) throw '请输入MCP插件描述'
    let data = {...formState}
    if (avatarData.value) data.avatar = avatarData.value?.file
    saveMcpServer(data).then(res => {
      emit('ok', formState)
      message.success('已保存')
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
