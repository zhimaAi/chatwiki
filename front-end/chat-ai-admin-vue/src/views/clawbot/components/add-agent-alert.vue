<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}
.form-item-tip {
  color: #8c8c8c;
}

.avatar-wrap {
  text-align: center;
  margin-bottom: 16px;
  :deep(.ant-upload.ant-upload-select) {
    width: 62px !important;
    height: 62px !important;
    border-radius: 16px !important;
  }
  :deep(.form-item-tip) {
    margin-top: 16px;
  }
}

.avatar-list {
  display: flex;
  align-items: center;
  gap: 4px;

  .avatar-item {
    width: 62px;
    height: 62px;
    border-radius: 16px;
    cursor: pointer;
    overflow: hidden;
    position: relative;

    &.active {
      .selected-icon {
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

    .selected-icon {
      display: none;
      position: absolute;
      width: 100%;
      height: 100%;
      z-index: 99;
      left: 0;
      top: 0;
      background: #00000052;

      img {
        width: 32px;
        height: 32px;
      }
    }

    .avatar {
      width: 100%;
      height: 100%;
    }
  }
}
</style>

<template>
  <a-modal
    :width="'746px'"
    v-model:open="show"
    :confirmLoading="saveLoading"
    :title="t('title_create_agent')"
    @ok="handleSave"
    @cancel="onCancel"
  >
    <div class="form-box">
      <a-form layout="vertical">
        <div class="avatar-wrap">
          <AvatarInput v-model:value="formState.robot_avatar_url" @change="onAvatarChange" />
          <div class="form-item-tip">{{ t('tip_upload_avatar') }}</div>
        </div>
        <a-form-item ref="name" :label="t('label_default_avatar')" v-bind="validateInfos.robot_avatar_url">
          <div class="avatar-list">
            <div
              v-for="(item, i) in defaultAvatars"
              :key="i"
              :class="['avatar-item', { active: formState.robot_avatar_url == item }]"
              @click="selectDefaultAvatar(item)"
            >
              <div class="selected-icon"><img src="@/assets/img/selected-icon.png" /></div>
              <img class="avatar" :src="item" />
            </div>
          </div>
        </a-form-item>

        <a-form-item ref="name" :label="t('label_agent_name')" v-bind="validateInfos.robot_name">
          <a-input v-model:value="formState.robot_name" :maxLength="20" :placeholder="t('ph_input_name')" />
        </a-form-item>

        <a-form-item :label="t('label_intro')" v-bind="validateInfos.robot_intro">
          <a-textarea
            :rows="3"
            v-model:value="formState.robot_intro"
            :placeholder="t('ph_input_intro')"
          />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Form, message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
import { saveRobot } from '@/api/robot/index'
import AvatarInput from '@/views/robot/robot-list/components/avatar-input.vue'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_ROBOT_AVATAR_LIST } from '@/constants/index'

const APPLICATION_TYPE = 2
const { t } = useI18n('views.clawbot.components.add-agent-alert')

const emit = defineEmits(['ok'])
const default_avatar = DEFAULT_ROBOT_AVATAR

const useForm = Form.useForm

const show = ref(false)
const saveLoading = ref(false)
const defaultAvatars = ref([])

const formState = reactive({
  robot_name: '',
  robot_intro: '',
  robot_avatar: undefined,
  robot_avatar_url: '',
  application_type: APPLICATION_TYPE,
  group_id: '0'
})

const rules = reactive({
  robot_name: [
    {
      required: true,
      message: t('validator_input_agent_name'),
      trigger: 'change'
    },
    {
      min: 1,
      max: 20,
      message: t('validator_agent_name_max'),
      trigger: 'change'
    }
  ]
})

const { validate, validateInfos, clearValidate } = useForm(formState, rules)

const onAvatarChange = (data) => {
  formState.robot_avatar = data.file
}

const selectDefaultAvatar = (val) => {
  formState.robot_avatar_url = val
  formState.robot_avatar = val
}

const onCancel = () => {
  clearValidate()
}

const saveForm = () => {
  let formData = {
    robot_name: formState.robot_name,
    robot_intro: formState.robot_intro,
    robot_avatar: formState.robot_avatar || default_avatar,
    group_id: formState.group_id || '0'
  }
  if (typeof formData.robot_avatar === 'string') {
    formData.robot_avatar_url = formData.robot_avatar
    delete formData.robot_avatar
  }

  saveLoading.value = true

  saveRobot(formData, APPLICATION_TYPE)
    .then((res) => {
      saveLoading.value = false

      if (res.res != 0) {
        return message.error(res.msg)
      }

      message.success(t('msg_create_success'))
      show.value = false
      emit('ok', res.data)
    })
    .catch(() => {
      saveLoading.value = false
    })
}

const handleSave = () => {
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

const open = () => {
  defaultAvatars.value = [default_avatar, ...DEFAULT_ROBOT_AVATAR_LIST]
  formState.robot_avatar = ''
  formState.robot_avatar_url = default_avatar
  formState.robot_name = ''
  formState.robot_intro = ''
  formState.group_id = '0'
  show.value = true
}

onMounted(() => {})

defineExpose({
  open
})
</script>
