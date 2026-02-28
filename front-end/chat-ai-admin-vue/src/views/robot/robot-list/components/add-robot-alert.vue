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
    :width="isEdit ? '472px' : '746px'"
    v-model:open="show"
    :confirmLoading="saveLoading"
    :title="title"
    @ok="handleSave"
    @cancel="onCancel"
  >
    <div class="form-box">
      <a-form layout="vertical">
          <div class="avatar-wrap">
            <AvatarInput v-model:value="formState.robot_avatar_url" @change="onAvatarChange" />
            <div class="form-item-tip">{{ t('ph_upload_avatar') }}</div>
          </div>
          <a-form-item ref="name" :label="t('label_default_avatar')" v-bind="validateInfos.robot_avatar_url">
            <div class="avatar-list">
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
          </a-form-item>

<!--        <a-form-item ref="name" label="" v-bind="validateInfos.robot_avatar_url" v-if="isEdit">-->
<!--          <div class="robot-avatar-box">-->
<!--            <AvatarInput :listType="'picture'" :style="{ width: '62px', height: '62px', borderRadius: '16px' }" v-model:value="formState.robot_avatar_url" @change="onAvatarChange" />-->
<!--            <div class="form-item-tip">{{ t('ph_upload_avatar') }}</div>-->
<!--          </div>-->
<!--        </a-form-item>-->

        <a-form-item required v-if="formState.application_type != 0" ref="en_name" :label="t('label_workflow_name')" v-bind="validateInfos.en_name">
          <a-input v-model:value="formState.en_name" :maxLength="50" :placeholder="t('ph_input')" />
          <div class="form-item-tip">{{ t('msg_workflow_name_format') }}</div>
        </a-form-item>

        <a-form-item  ref="name" :label="formState.application_type == 0 ? t('label_app_name') : t('label_remark_name')" v-bind="validateInfos.robot_name">
          <a-input v-model:value="formState.robot_name" :maxLength="20" :placeholder="t('ph_input')" />
          <div  v-if="formState.application_type != 0"  class="form-item-tip">{{ t('msg_remark_name_tip') }}</div>
        </a-form-item>

        <a-form-item  ref="group_id" :label="t('label_group')" v-bind="validateInfos.group_id">
          <a-select v-model:value="formState.group_id" style="width: 100%" :placeholder="t('ph_select_group')">
            <a-select-option v-for="item in groupLists" :value="item.id">{{ item.group_name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item :label="t('label_intro')" v-bind="validateInfos.robot_intro">
          <a-textarea
            v-if="!isEdit"
            :rows="3"
            v-model:value="formState.robot_intro"
            :placeholder="t('ph_intro_chat')"
          />
          <a-textarea
            v-else
            :rows="3"
            v-model:value="formState.robot_intro"
            :placeholder="t('ph_intro_workflow')"
          />
        </a-form-item>
      </a-form>
    </div>
  </a-modal>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { useI18n } from '@/hooks/web/useI18n'
import { ref, reactive, h, onMounted, computed } from 'vue'
import { Form, message, Modal } from 'ant-design-vue'
import { CloseCircleFilled } from '@ant-design/icons-vue'
import { getModelConfigOption } from '@/api/model/index'
import { saveRobot, editBaseInfo, getRobotGroupList } from '@/api/robot/index'
import { useRoute, useRouter } from 'vue-router'
import AvatarInput from './avatar-input.vue'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_WORKFLOW_AVATAR, DEFAULT_ROBOT_AVATAR_LIST } from '@/constants/index'
import { useRobotStore } from '@/stores/modules/robot'

const { t } = useI18n('views.robot.robot-list.components.add-robot-alert')
const robotStore = useRobotStore()

const { setStorage } = useStorage('localStorage')

const emit = defineEmits(['addRobot', 'editRobot', 'ok'])
let default_avatar = ''

const route = useRoute()
const router = useRouter()
const useForm = Form.useForm

const isEdit = ref(false)
const show = ref(false)

const saveLoading = ref(false)
const groupLists = ref([])
const showDefaultAvatars = ref([])

const formState = reactive({
  robot_name: '',
  en_name: '',
  robot_intro: '',
  robot_avatar: undefined,
  robot_avatar_url: '',
  application_type: 0,
  group_id: '0'
})

const title = computed(() => {
  if(formState.application_type == 0){
    if(isEdit.value){
      return t('title_edit_robot_info')
    }else{
      return t('title_create_chat_robot')
    }

  }else{
    if(isEdit.value){
      return t('title_edit_workflow_info')
    }else{
      return t('title_create_workflow')
    }
  }
})

const robotInfo = computed(() => {
  return robotStore.robotInfo
})

const { getRobot } = robotStore

const rules = computed(() => ({
  robot_name: [
    {
      required: true,
      message: t('msg_robot_name_required'),
      trigger: 'change'
    },
    {
      min: 1,
      max: 20,
      message: t('msg_max_20_chars'),
      trigger: 'change'
    }
  ],
  en_name:[
  ]
}))

const { validate, validateInfos, clearValidate } = useForm(formState, rules)

const onAvatarChange = (data) => {
  formState.robot_avatar = data.file
}

const saveForm = () => {
  let formData = {
    robot_name: formState.robot_name,
    en_name: formState.en_name,
    robot_intro: formState.robot_intro,
    robot_avatar: formState.robot_avatar || default_avatar,
    group_id: formState.group_id || '0'
  }
  if (typeof formData.robot_avatar === "string") {
    formData.robot_avatar_url = formData.robot_avatar
    delete formData.robot_avatar
  }

  saveLoading.value = true

  let requertUrl = saveRobot
  let tipVal = 'msg_robot_created'
  if (isEdit.value) {
    requertUrl = editBaseInfo
    formData.id = route.query.id
    tipVal = 'msg_robot_edited'
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

      message.success(t(tipVal))
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
          // message.info('按住Shift 滚动鼠标可左右移动画布，按住Ctrl 滚动鼠标可放大缩小画布', 6)
          window.open('/#/robot/config/workflow?id=' + res.data.id + '&robot_key=' + res.data.robot_key + '&show_tips=1')
          show.value = false
           emit('ok')
          // router.push('/robot/config/workflow?id=' + res.data.id + '&robot_key=' + res.data.robot_key)
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
      title: t('msg_add_llm_provider_title'),
      icon: h(CloseCircleFilled, {
        style: {
          color: 'red'
        }
      }),
      content: t('msg_add_llm_provider_content'),
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
  showDefaultAvatars.value = [default_avatar, ...DEFAULT_ROBOT_AVATAR_LIST]
  isEdit.value = is_edit
  // 验证是否有LLM
  let state = await checkLLM()
  getGroupList()
  if (state) {
    formState.robot_avatar = ''
    formState.robot_avatar_url = default_avatar
    formState.robot_name = ''
    formState.en_name = ''
    formState.robot_intro = ''
    formState.application_type = type
    show.value = true

    if (isEdit.value && robotInfo.value) {
      formState.robot_name = robotInfo.value.robot_name
      formState.en_name = robotInfo.value.en_name
      formState.robot_intro = robotInfo.value.robot_intro
      formState.robot_avatar_url = robotInfo.value.robot_avatar_url
    }
  }
  if(type == 0) {
    rules.value.en_name = []
  }else{
  rules.value.en_name = [
    {
      required: true,
      message: t('msg_workflow_name_required'),
      trigger: 'change'
    },
    {
      validator: async (rule, value, callback) => {
        if(!value){
          return Promise.resolve()
        }
        //只能输入英文数字和字符“_”、“”、“_”
        if (!/^[a-zA-Z0-9_\.\-]+$/.test(value)) {
           return Promise.reject(t('msg_workflow_name_format_error'))
        }else{
          return Promise.resolve()
        }
      },
    }
  ]
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
  if(!isEdit.value && formState.application_type != 0){
    if(!formState.robot_name){
      formState.robot_name = formState.en_name
    }
  }
  validate()
    .then(() => {
      saveForm()
    })
    .catch((err) => {
      console.log('error', err)
    })
}

const selectDefaultAvatar = (val) => {
  formState.robot_avatar_url = val
  formState.robot_avatar = val
}

onMounted(() => {

})

defineExpose({
  open,
  setGroupId,
})
</script>
