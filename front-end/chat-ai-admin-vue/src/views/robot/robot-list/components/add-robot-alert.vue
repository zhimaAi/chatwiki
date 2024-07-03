<style lang="less" scoped>
.form-box {
  margin-top: 24px;
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
      <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-item ref="name" label="机器人名称" v-bind="validateInfos.robot_name">
          <a-input v-model:value="formState.robot_name" placeholder="请输入机器人名称" />
        </a-form-item>

        <a-form-item label="简介" v-bind="validateInfos.robot_intro">
          <a-textarea
            :rows="4"
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
  robot_avatar_url: ''
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

  saveRobot(formData)
    .then((res) => {
      saveLoading.value = false

      if (res.res != 0) {
        return message.error(res.msg)
      }

      message.success('机器人创建成功')
      router.push('/robot/config/basic-config?id=' + res.data.id)
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
