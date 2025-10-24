<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}
.form-item-tip {
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

.robot-avatar-box {
  display: flex;
  flex-direction: column;
  align-items: center;

  .form-item-tip {
    color: #8c8c8c;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
}
</style>

<template>
  <a-modal
    :width="isEdit ? '472px' : '746px'"
    v-model:open="show"
    :confirmLoading="saveLoading"
    :title="title"
    @ok="handleSave"
    @cancel="onCancel"
  >
    <div class="form-box">
      <a-form layout="vertical">
        <!-- <a-form-item label="应用类型" required>
          <div class="robot-type-box">
            <div
              class="robot-item"
              :class="{ active: formState.application_type == 0 }"
              @click="formState.application_type = 0"
            >
              <svg-icon
                v-if="formState.application_type == 0"
                class="check-arrow"
                name="check-arrow-filled"
              ></svg-icon>
              <div class="avatar-box">
                <img src="@/assets/svg/chat-robot-icon.svg" alt="" />
                <div class="robot-item-title">聊天机器人</div>
              </div>
              <div class="desc">简易的基于LLM和知识库的聊天应用，配置简单，适合新人上手</div>
            </div>
            <div
              class="robot-item"
              :class="{ active: formState.application_type == 1 }"
              @click="formState.application_type = 1"
            >
              <svg-icon
                v-if="formState.application_type == 1"
                class="check-arrow"
                name="check-arrow-filled"
              ></svg-icon>
              <div class="avatar-box">
                <img src="@/assets/svg/workflow-robot-icon.svg" alt="" />
                <div class="robot-item-title">工作流</div>
              </div>
              <div class="desc">可编排的任务工作流，将任务拆解成多个步骤逐步执行</div>
            </div>
          </div>
        </a-form-item> -->
        <a-form-item ref="name" label="" v-bind="validateInfos.robot_avatar_url" v-if="isEdit">
          <div class="robot-avatar-box">
            <AvatarInput :listType="'picture'" :style="{ width: '62px', height: '62px', borderRadius: '16px' }" v-model:value="formState.robot_avatar_url" @change="onAvatarChange" />
            <div class="form-item-tip">请上传机器人头像，建议尺寸为100*100px,大小不超过100KB</div>
          </div>
        </a-form-item>

        <a-form-item  ref="name" label="应用名称" v-bind="validateInfos.robot_name">
          <a-input v-model:value="formState.robot_name" placeholder="请输入应用名称" />
        </a-form-item>

        <a-form-item  ref="group_id" label="分组" v-bind="validateInfos.group_id">
          <a-select v-model:value="formState.group_id" style="width: 100%" placeholder="请选择分组">
            <a-select-option v-for="item in groupLists" :value="item.id">{{ item.group_name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item label="简介" v-bind="validateInfos.robot_intro">
          <a-textarea
            v-if="!isEdit"
            :rows="3"
            v-model:value="formState.robot_intro"
            placeholder="请输入机器人简介，比如ChatWiki产品帮助机器人，可以通过提问获取ChatWiki的使用帮助，比如如何创建机器人，如何新建知识库，如何添加模型等。"
          />
          <a-textarea
            v-else
            :rows="3"
            v-model:value="formState.robot_intro"
            placeholder="请输入机器人简介，比如 ZHIMA CHATAI 基于大预言模型提供 ZHIMA CHATAI 产品帮助"
          />
        </a-form-item>

        <a-form-item ref="name" label="机器人头像" v-bind="validateInfos.robot_avatar_url" v-if="!isEdit">
          <AvatarInput v-model:value="formState.robot_avatar_url" @change="onAvatarChange" />
          <div class="form-item-tip">请上传机器人头像，建议尺寸为100*100px,大小不超过100KB</div>
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { ref, reactive, h, onMounted, computed } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import { CloseCircleFilled } from '@ant-design/icons-vue'
import { getModelConfigOption } from '@/api/model/index'
import { saveRobot, editBaseInfo, getRobotGroupList } from '@/api/robot/index'
import { useRoute, useRouter } from 'vue-router'
import AvatarInput from './avatar-input.vue'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_WORKFLOW_AVATAR } from '@/constants/index'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()

const { setStorage } = useStorage('localStorage')

const emit = defineEmits(['addRobot', 'editRobot'])
let default_avatar = ''

const route = useRoute()
const router = useRouter()
const useForm = Form.useForm

const isEdit = ref(false)
const show = ref(false)

const saveLoading = ref(false)
const groupLists = ref([])

const formState = reactive({
  robot_name: '',
  robot_intro: '',
  robot_avatar: undefined,
  robot_avatar_url: '',
  application_type: 0,
  group_id: '0'
})

const title = computed(() => {
  if(formState.application_type == 0){
    if(isEdit.value){
      return '编辑机器人基本信息'
    }else{
      return '新建聊天机器人'
    }

  }else{
    if(isEdit.value){
      return '编辑工作流基本信息'
    }else{
      return '新建工作流'
    }
  }
})

const robotInfo = computed(() => {
  return robotStore.robotInfo
})

const { getRobot } = robotStore

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
    robot_avatar: formState.robot_avatar || default_avatar,
    group_id: formState.group_id || '0'
  }

  saveLoading.value = true

  let requertUrl = saveRobot
  let tipVal = '机器人创建成功'
  if (isEdit.value) {
    requertUrl = editBaseInfo
    formData.id = route.query.id
    tipVal = '机器人编辑成功'
    if (formData.robot_avatar == default_avatar) {
      delete formData.robot_avatar
    }
  }

  requertUrl(formData, formState.application_type)
    .then((res) => {
      saveLoading.value = false

      if (res.res != 0) {
        return message.error(res.msg)
      }

      message.success(tipVal)
      if(!isEdit.value && route.query.robot_key){
        emit('addRobot')
        show.value = false
        return
      }
      if (!isEdit.value) {
        if (formState.application_type == 0) {
          router.push(
            '/robot/config/basic-config?id=' + res.data.id + '&robot_key=' + res.data.robot_key
          )
        } else {
          router.push('/robot/config/workflow?id=' + res.data.id + '&robot_key=' + res.data.robot_key)
        }
      } else {
        // 弹窗关闭 刷新数据
        getRobot(route.query.id)
        emit('editRobot')
        show.value = false
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

const open = async (type, is_edit) => {
  if(type == 0){
    default_avatar = DEFAULT_ROBOT_AVATAR
  } else {
    default_avatar = DEFAULT_WORKFLOW_AVATAR
  }
  isEdit.value = is_edit
  // 验证是否有LLM
  let state = await checkLLM()
  getGroupList()
  if (state) {
    formState.robot_avatar = ''
    formState.robot_avatar_url = default_avatar
    formState.robot_name = ''
    formState.robot_intro = ''
    formState.application_type = type
    show.value = true

    if (isEdit.value && robotInfo.value) {
      formState.robot_name = robotInfo.value.robot_name
      formState.robot_intro = robotInfo.value.robot_intro
      formState.robot_avatar_url = robotInfo.value.robot_avatar_url
    }
  }
}

const setGroupId = (id) => {
  formState.group_id = id || '0'
}

const getGroupList = () => {
  getRobotGroupList().then((res) => {
    groupLists.value  = res.data || []
  })
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

onMounted(() => {

})

defineExpose({
  open,
  setGroupId,
})
</script>
