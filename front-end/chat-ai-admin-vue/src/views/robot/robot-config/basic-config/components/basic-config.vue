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
    title="基本配置"
    icon-name="jibenpeizhi"
    v-model:isEdit="isEdit"
    @save="onSave"
    @edit="handleEdit"
  >
    <div class="robot-form-box" v-show="isEdit">
      <div class="form-item">
        <div class="form-item-label">机器人名称/头像</div>
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
            <a-input v-model:value="formState.robot_name" placeholder="请输入机器人名称" />
          </div>
        </div>
      </div>

      <div class="form-item">
        <div class="form-item-label">简介</div>
        <div class="form-item-body">
          <a-textarea
            :rows="3"
            v-model:value="formState.robot_intro"
            placeholder="请输入机器人简介，比如 ZHIMA CHATAI 基于大语言模型提供ZHIMA CHATAI 产品帮助。"
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
        {{ robotInfo.robot_intro || '还没有设置机器人介绍，快来写一段简介吧' }}
      </div>
    </div>
  </edit-box>
</template>
<script setup>
import { getBase64 } from '@/utils/index'
import { ref, reactive, inject, toRaw } from 'vue'
import EditBox from './edit-box.vue'
import CuUpload from '@/components/cu-upload/cu-upload.vue'

const isEdit = ref(false)
const { robotInfo, updateRobotInfo } = inject('robotInfo')

const formState = reactive({
  robot_name: '',
  robot_avatar: '',
  robot_avatar_url: '',
  robot_intro: ''
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
