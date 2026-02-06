<template>
  <a-modal
    v-model:open="visible"
    :title="formState.id ? t('title_edit_http_tool') : t('title_add_http_tool')"
    width="472px"
    :confirm-loading="saving"
    @ok="save">
    <a-form class="form-box" labelAlign="left">
      <a-form-item :label="t('label_tool_cover')" :colon="false">
        <HttpToolAvatarInput v-model:value="formState.avatar" />
        <div class="form-item-tip">{{ t('tip_cover_size') }}</div>
      </a-form-item>
    </a-form>
    <a-form class="form-box" labelAlign="left">
      <a-form-item :label="t('label_name')" required :colon="false">
        <a-input v-model:value="formState.name" :placeholder="t('ph_input_name')" :maxlength="20"/>
      </a-form-item>
      <a-form-item :label="t('label_description')" :colon="false">
        <a-textarea
          v-model:value="formState.description"
          :auto-size="{ minRows: 2, maxRows: 5 }"
          :placeholder="t('ph_input_desc')" :maxlength="60"/>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import {ref, reactive} from 'vue';
import {message} from 'ant-design-vue';
import { saveHttpTool } from '@/api/robot/http_tool.js';
import HttpToolAvatarInput from './http-tool-avatar-input.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-list.components.http-tool-store')

const emit = defineEmits(['ok'])

const visible = ref(false)
const saving = ref(false)
const formStateStruct = {
  id: '',
  avatar: '/upload/default/http-node.png',
  name: '',
  description: ''
}
const formState = reactive({})

function show(info = null) {
  const base = JSON.parse(JSON.stringify(formStateStruct))
  const source = info || {}
  base.id = source.id || ''
  base.name = source.name || ''
  base.description = source.description || ''
  if (source.avatar) {
    base.avatar = source.avatar
  }
  Object.assign(formState, base)
  visible.value = true
}

function save() {
  try {
    saving.value = true
    formState.name = String(formState.name || '').trim()
    formState.description = String(formState.description || '').trim()
    if (!formState.name) throw t('msg_input_name')
    let data = {...formState}
    const avatarLink = String(formState.avatar || '')
    if (avatarLink) data.avatar = avatarLink
    saveHttpTool(data).then(() => {
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
  show
})
</script>

<style scoped lang="less">
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

.form-item-tip {
  color: #8c8c8c;
  font-size: 12px;
  font-weight: 400;
  line-height: 14px;
  margin: 6px 0 0;
  white-space: nowrap;
}
</style>
