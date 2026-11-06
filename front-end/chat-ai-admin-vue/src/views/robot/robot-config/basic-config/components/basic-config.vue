<style lang="less" scoped>
.basic-config {
  .robot-info-box {
    .robot-name-box {
      display: flex;
      align-items: center;

      .robot-avatar {
        width: 32px;
        height: 32px;
        margin-right: 8px;
        border-radius: 4px;
      }

      .robot-name {
        line-height: 22px;
        font-size: 14px;
        color: #262626;
      }
    }

    .robot-intro {
      line-height: 22px;
      margin-top: 8px;
      font-size: 14px;
      color: #595959;
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 8px;
    }

  }

  .robot-form-box {
    display: flex;
    gap: 50px;

    > div {
      flex: 1;
    }

    .form-item-label {
      line-height: 22px;
      margin-bottom: 4px;
      font-size: 14px;
      color: #262626;
    }

    .robot-name-item {
      display: flex;
      align-items: center;
    }

    .robot-avetar-upload {
      width: 40px;
      height: 40px;
      margin-right: 8px;

      .robot-avetar {
        width: 40px;
        height: 40px;
        border-radius: 4px;
      }
    }

    .robot-name-box {
      flex: 1;
    }
  }
}

.basic-config.profile-summary {
  position: relative;
  padding: 24px;
  border: 1px solid #e8e8e8;
  border-radius: 12px;
  background: #eff7ff url('@/assets/img/customer-service/basic-settings/profile-background.png') center / cover no-repeat;

  :deep(.edit-box-header) {
    position: absolute;
    top: 50%;
    right: calc(20px - 1px);
    height: 32px;
    padding: 0;
    background: transparent;
    transform: translateY(-50%);
    .left-box { display: none; }
  }
  :deep(.edit-box-body) { padding: 0; }
  .inline-summary { display: flex; align-items: center; gap: 12px; padding-right: 80px; }
  .profile-avatar {
    width: 48px;
    height: 48px;
    flex-shrink: 0;
    padding: 0;
    border: 0;
    border-radius: 10px;
    overflow: hidden;
    background: #e8f0fe;
    cursor: pointer;
    img { display: block; width: 48px; height: 48px; object-fit: cover; }
    &:focus-visible { outline: 2px solid #2475fc; outline-offset: 2px; }
  }
  .profile-copy {
    min-width: 0;
    .robot-name { font-size: 16px; line-height: 24px; font-weight: 600; color: #222; }
    .robot-intro { display: block; margin-top: 4px; color: #7a7a7a; }
    .robot-name, .robot-intro { overflow: hidden; white-space: nowrap; text-overflow: ellipsis; }
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

.mt16 {
  margin-top: 16px;
}
</style>

<template>
  <edit-box
    class="basic-config"
    :class="{ 'profile-summary': inlineSummary }"
    :title="t('title_basic_config')"
    icon-name="jibenpeizhi"
    :is-edit="!inlineSummary && isEdit"
    @update:isEdit="value => { if (!inlineSummary) isEdit = value }"
    @save="onSave"
    @edit="handleEdit"
  >
    <div class="robot-form-box" v-if="!inlineSummary" v-show="isEdit">
      <div>
        <div class="form-item">
          <div class="form-item-label">{{ t('label_robot_name_avatar') }}</div>
          <div class="form-item-body robot-name-item">
            <cu-upload accept="image/*" @change="handleChangeAvatar">
              <div class="robot-avetar-upload">
                <img
                  class="robot-avetar"
                  :src="formState.robot_avatar_url"
                  v-if="formState.robot_avatar_url"
                />
              </div>
            </cu-upload>

            <div class="robot-name-box">
              <a-input v-model:value="formState.robot_name" :placeholder="t('ph_input_robot_name')" />
            </div>
          </div>
        </div>
        <div class="form-item mt16">
          <div class="form-item-label">{{ t('label_default_avatar') }}</div>
          <div class="form-item-body avatar-list">
            <div
              v-for="(item, i) in showDefaultAvatars"
              :key="i"
              :class="['avatar-item', {active: formState.robot_avatar_url == item}]"
              @click="selectDefaultAvatar(item)"
            >
              <div class="selected-icon"><img src="@/assets/img/selected-icon.png"/></div>
              <img class="avatar" :src="item"/>
            </div>
          </div>
        </div>
      </div>

      <div class="form-item">
        <div class="form-item-label">{{ t('label_intro') }}</div>
        <div class="form-item-body">
          <a-textarea
            :rows="3"
            v-model:value="formState.robot_intro"
            :placeholder="t('ph_input_robot_intro')"
            allow-clear
          />
        </div>
      </div>
    </div>
    <div class="robot-info-box" :class="{ 'inline-summary': inlineSummary }" v-show="inlineSummary || !isEdit">
      <template v-if="inlineSummary">
        <button type="button" class="profile-avatar" :aria-label="t('label_robot_name_avatar')" @click="handleEdit">
          <img :src="robotInfo.robot_avatar_url || showDefaultAvatars[0]" alt="" />
        </button>
        <div class="profile-copy">
          <div class="robot-name" :title="robotInfo.robot_name">{{ robotInfo.robot_name }}</div>
          <div class="robot-intro" :title="robotInfo.robot_intro || t('msg_no_intro')">{{ robotInfo.robot_intro || t('msg_no_intro') }}</div>
        </div>
      </template>
      <template v-else>
      <div class="robot-name-box">
        <img class="robot-avatar" :src="robotInfo.robot_avatar_url" alt="" />
        <span class="robot-name" :title="inlineSummary ? robotInfo.robot_name : undefined">{{ robotInfo.robot_name }}</span>
      </div>
      <div class="robot-intro" :title="inlineSummary ? robotInfo.robot_intro || t('msg_no_intro') : undefined">
        <span v-if="!inlineSummary">{{ t('label_group') }}{{ groupName }}</span>
        <span>{{ t('label_intro_colon') }}{{ robotInfo.robot_intro || t('msg_no_intro') }}</span>
      </div>
      </template>
    </div>
    <ProfileSettingsDialog v-if="inlineSummary" ref="profileDialog" />
  </edit-box>
</template>
<script setup>
import { getBase64 } from '@/utils/index'
import { ref, reactive, inject, toRaw, computed } from 'vue'
import { message } from 'ant-design-vue'
import EditBox from './edit-box.vue'
import CuUpload from '@/components/cu-upload/cu-upload.vue'
import ProfileSettingsDialog from '@/views/customer-service/basic-settings/components/profile-settings-dialog.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useI18n } from '@/hooks/web/useI18n'
import {DEFAULT_ROBOT_AVATAR, DEFAULT_ROBOT_AVATAR_LIST, DEFAULT_WORKFLOW_AVATAR} from "@/constants/index.js";

const { t } = useI18n('views.robot.robot-config.basic-config.components.basic-config')
const props = defineProps({ inlineSummary: { type: Boolean, default: false } })
const robotStore = useRobotStore()
const isEdit = ref(false)
const profileDialog = ref(null)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  robot_name: '',
  robot_avatar: '',
  robot_avatar_url: '',
  robot_intro: ''
})

const robotList = computed(() => {
  return robotStore.robotList || []
})

const robotGroupList = computed(() => {
  return robotStore.robotGroupList || []
})

const showDefaultAvatars = computed(() => {
  let def_avatar = Number(robotInfo.application_type) === 1 ? DEFAULT_WORKFLOW_AVATAR : DEFAULT_ROBOT_AVATAR
  return [def_avatar, ...DEFAULT_ROBOT_AVATAR_LIST]
})


const groupName = computed(() => {
  let group_name = t('msg_ungrouped')
  let currentRobotItem = robotList.value.find((item) => item.id == robotInfo.id) || {}
  if(currentRobotItem.group_id > 0){
    let groupItem = robotGroupList.value.find((item) => item.id == currentRobotItem.group_id)
    if (groupItem) {
      group_name = groupItem.group_name
    }
  }
  return group_name
})

const onSave = () => {
  updateRobotInfo({ ...toRaw(formState) })
  isEdit.value = false;
}

const handleChangeAvatar = (fileList) => {
  let file = fileList[0]

  if (file.size > 100 * 1024) {
    message.error(t('msg_avatar_size_limit'))
    return false
  }

  formState.robot_avatar = file

  getBase64(file).then((base64Url) => {
    formState.robot_avatar_url = base64Url
  })
  return false
}

const handleEdit = () => {
  if (props.inlineSummary) {
    profileDialog.value?.open()
    return
  }
  formState.robot_name = robotInfo.robot_name
  formState.robot_intro = robotInfo.robot_intro
  formState.robot_avatar = robotInfo.robot_avatar
  formState.robot_avatar_url = robotInfo.robot_avatar_url
}

const selectDefaultAvatar = (val) => {
  formState.robot_avatar_url = val
  formState.robot_avatar = val
}
</script>
