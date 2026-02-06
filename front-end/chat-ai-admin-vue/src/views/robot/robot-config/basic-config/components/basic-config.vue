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

    .form-item {
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
</style>

<template>
  <edit-box
    class="basic-config"
    :title="t('title_basic_config')"
    icon-name="jibenpeizhi"
    v-model:isEdit="isEdit"
    @save="onSave"
    @edit="handleEdit"
  >
    <div class="robot-form-box" v-show="isEdit">
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
    <div class="robot-info-box" v-show="!isEdit">
      <div class="robot-name-box">
        <img class="robot-avatar" :src="robotInfo.robot_avatar_url" alt="" />
        <span class="robot-name">{{ robotInfo.robot_name }}</span>
      </div>
      <div class="robot-intro">
        <span>{{ t('label_group') }}{{ groupName }}</span>
        <span>{{ t('label_intro_colon') }}{{ robotInfo.robot_intro || t('msg_no_intro') }}</span>
      </div>
    </div>
  </edit-box>
</template>
<script setup>
import { getBase64 } from '@/utils/index'
import { ref, reactive, inject, toRaw, computed } from 'vue'
import EditBox from './edit-box.vue'
import CuUpload from '@/components/cu-upload/cu-upload.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.basic-config')
const robotStore = useRobotStore()
const isEdit = ref(false)
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
  formState.robot_avatar = file

  getBase64(file).then((base64Url) => {
    formState.robot_avatar_url = base64Url
  })
  return false
}

const handleEdit = () => {
  formState.robot_name = robotInfo.robot_name
  formState.robot_intro = robotInfo.robot_intro
  formState.robot_avatar = robotInfo.robot_avatar
  formState.robot_avatar_url = robotInfo.robot_avatar_url
}
</script>
