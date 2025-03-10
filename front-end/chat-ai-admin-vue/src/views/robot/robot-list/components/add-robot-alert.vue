<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}
.form-item-tip{
  color: #8c8c8c;
}
.robot-type-box {
  display: flex;
  align-items: center;
  gap: 18px;
  .robot-item {
    padding: 16px;
    flex: 1;
    height: 128px;
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
    .desc {
      margin-top: 12px;
      color: #595959;
      font-size: 14px;
      line-height: 22px;
    }
    .avatar-box {
      display: flex;
      align-items: center;
      gap: 16px;

      img {
        width: 40px;
        height: 40px;
      }
    }
    .robot-item-title {
      font-size: 16px;
      color: #262626;
      font-weight: 600;
    }

    &.active {
      border: 1px solid #2475fc;
      .robot-item-title {
        color: #2475fc;
      }
    }

    .check-arrow {
      position: absolute;
      display: block;
      right: -1px;
      bottom: -1px;
      width: 24px;
      height: 24px;
      font-size: 24px;
      color: #fff;
    }
  }
}
</style>

<template>
  <a-modal
    width="746px"
    v-model:open="show"
    :confirmLoading="saveLoading"
    title="新建机器人"
    @ok="handleSave"
    @cancel="onCancel"
  >
    <div class="form-box">
      <a-form layout="vertical">
        <a-form-item label="应用类型" required>
          <div class="robot-type-box">
            <div class="robot-item" :class="{ active: formState.application_type == 0 }" @click="formState.application_type = 0">
              <svg-icon v-if="formState.application_type == 0" class="check-arrow" name="check-arrow-filled"></svg-icon>
              <div class="avatar-box">
                <img src="@/assets/svg/chat-robot-icon.svg" alt="" />
                <div class="robot-item-title">聊天机器人</div>
              </div>
              <div class="desc">简易的基于LLM和知识库的聊天应用，配置简单，适合新人上手</div>
            </div>
            <div class="robot-item" :class="{ active: formState.application_type == 1 }" @click="formState.application_type = 1">
              <svg-icon v-if="formState.application_type == 1" class="check-arrow" name="check-arrow-filled"></svg-icon>
              <div class="avatar-box">
                <img src="@/assets/svg/workflow-robot-icon.svg" alt="" />
                <div class="robot-item-title">工作流</div>
              </div>
              <div class="desc">可编排的任务工作流，将任务拆解成多个步骤逐步执行</div>
            </div>
          </div>
        </a-form-item>
        <a-form-item ref="name" label="应用名称" v-bind="validateInfos.robot_name">
          <a-input v-model:value="formState.robot_name" placeholder="请输入应用名称" />
        </a-form-item>

        <a-form-item label="简介" v-bind="validateInfos.robot_intro">
          <a-textarea
            :rows="3"
            v-model:value="formState.robot_intro"
            placeholder="请输入机器人简介，比如ChatWiki产品帮助机器人，可以通过提问获取ChatWiki的使用帮助，比如如何创建机器人，如何新建知识库，如何添加模型等。"
          />
        </a-form-item>

        <a-form-item ref="name" label="机器人头像" v-bind="validateInfos.robot_avatar_url">
          <AvatarInput v-model:value="formState.robot_avatar_url" @change="onAvatarChange" />
          <div class="form-item-tip">请上传机器人头像，建议尺寸为100*100px,大小不超过100KB</div>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, reactive, h } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import { CloseCircleFilled } from '@ant-design/icons-vue'
import { getModelConfigOption } from '@/api/model/index'
import { saveRobot } from '@/api/robot/index'
import { useRouter } from 'vue-router'
import AvatarInput from './avatar-input.vue'
import { DEFAULT_ROBOT_AVATAR } from '@/constants/index'

const router = useRouter()
const useForm = Form.useForm

const labelCol = {
  span: 4
}
const wrapperCol = {
  span: 20
}

const show = ref(false)

const saveLoading = ref(false)

const formState = reactive({
  robot_name: '',
  robot_intro: '',
  robot_avatar: undefined,
  robot_avatar_url: '',
  application_type: 0
})

const rules = reactive({
  robot_name: [
    {
      required: true,
      message: '请输入机器人名称',
      trigger: 'change'
    },
    {
      min: 1,
      max: 20,
      message: '最多20个字',
      trigger: 'change'
    }
  ]
})

const { validate, validateInfos, clearValidate } = useForm(formState, rules)

const onAvatarChange = (data) => {
  formState.robot_avatar = data.file
}

const saveForm = () => {
  let formData = {
    robot_name: formState.robot_name,
    robot_intro: formState.robot_intro,
    robot_avatar: formState.robot_avatar || DEFAULT_ROBOT_AVATAR
  }

  saveLoading.value = true

  saveRobot(formData, formState.application_type)
    .then((res) => {
      saveLoading.value = false

      if (res.res != 0) {
        return message.error(res.msg)
      }

      message.success('机器人创建成功')
      if(formState.application_type == 0) {
        router.push('/robot/config/basic-config?id=' + res.data.id + '&robot_key=' + res.data.robot_key)
      } else {
        router.push('/robot/config/workflow?id=' + res.data.id + '&robot_key=' + res.data.robot_key)
      }
    })
    .catch(() => {
      saveLoading.value = false
    })
}

// 验证是否有LLM
const checkLLM = async () => {
  let res = await getModelConfigOption({
    model_type: 'LLM'
  })

  if (!res || res.data.length == 0) {
    Modal.confirm({
      title: '请先添加LLM模型服务商?',
      icon: h(CloseCircleFilled, {
        style: {
          color: 'red'
        }
      }),
      content: '机器人聊天基于LLM(大语言模型)，请先完成LLM服务商配置',
      onOk() {
        router.push('/user/model')
      },
      onCancel() {}
    })

    return false
  }

  return true
}

const open = async () => {
  // 验证是否有LLM
  let state = await checkLLM()

  if (state) {
    formState.robot_avatar = ''
    formState.robot_avatar_url = DEFAULT_ROBOT_AVATAR
    formState.robot_name = ''
    formState.robot_intro = ''
    show.value = true
  }
}

const onCancel = () => {
  clearValidate()
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

defineExpose({
  open
})
</script>
