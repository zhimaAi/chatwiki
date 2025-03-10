<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="476px">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item label="登录账号" v-bind="validateInfos.user_name">
            <a-input
              :disabled="formState.id != ''"
              type="text"
              placeholder="账号只能为字母、数字、“-”、“_”,“.”的组合"
              v-model:value="formState.user_name"
            ></a-input>
          </a-form-item>
          <a-form-item label="成员昵称" v-bind="validateInfos.nick_name">
            <a-input
              :maxlength="100"
              type="text"
              placeholder="请输入昵称"
              v-model:value="formState.nick_name"
            ></a-input>
          </a-form-item>
          <a-form-item label="成员头像">
            <AvatarInput v-model:value="formState.avatar_url" @change="onAvatarChange" />
          </a-form-item>
          <a-form-item label="成员角色" v-bind="validateInfos.user_roles">
            <a-select
              v-model:value="formState.user_roles"
              style="width: 100%"
              placeholder="请选择成员角色"
            >
              <a-select-option v-for="item in roleLists" :key="item.id" :value="item.id">{{
                item.name
              }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="登录密码" v-bind="validateInfos.password" v-if="!formState.id">
            <a-input-password
              v-model:value="formState.password"
              placeholder="密码必须包含字母、数字或者字符中的两种，6-32位"
            />
          </a-form-item>
          <a-form-item label="确认密码" v-bind="validateInfos.check_password" v-if="!formState.id">
            <a-input-password
              v-model:value="formState.check_password"
              placeholder="请重新输入密码"
            />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { validatePassword } from '@/utils/validate.js'
import { ref, reactive, toRaw } from 'vue'
import { Form, message } from 'ant-design-vue'
import AvatarInput from './avatar-input.vue'
import { saveUser, getRoleList } from '@/api/manage/index.js'
import defaultAvatar from "@/assets/img/role_avatar.png"
const emit = defineEmits(['ok'])


const useForm = Form.useForm
const show = ref(false)
const modalTitle = ref('添加团队成员')
const formState = reactive({
  user_name: '',
  nick_name: '',
  avatar: '',
  avatar_url: '',
  user_roles: '3',
  password: '',
  check_password: '',
  id: ''
})
const onAvatarChange = (data) => {
  formState.avatar = data.file
}
const roleLists = ref([])
getRoleList({
  page: 1,
  size: 200,
  search: ''
}).then((res) => {
  let lists = res.data.list || []
  roleLists.value = lists.filter((item) => item.id != '1')
})
const formRules = reactive({
  user_name: [
    {
      message: '请输入登录账号',
      required: true
    },
    {
      validator: async (rule, value) => {
        if (!/^[a-zA-Z0-9_.-]+$/.test(value) && value) {
          return Promise.reject('账号只能为字母、数字、“-”、“_”,“.”的组合')
        }
        return Promise.resolve()
      }
    }
  ],
  nick_name: [
    {
      message: '请输入昵称',
      required: true
    }
  ],
  avatar_url: [
    {
      message: '请上传头像',
      required: true
    }
  ],
  user_roles: [
    {
      message: '请选择成员角色',
      required: true
    }
  ],
  password: [
    {
      message: '请输入登录密码',
      required: true
    },
    {
      validator: async (rule, value) => {
        if (!validatePassword(value) && value) {
          return Promise.reject('密码必须包含字母、数字或者字符中的两种，6-32位')
        }
        return Promise.resolve()
      }
    }
  ],
  check_password: [
    {
      message: '请输入确认密码',
      required: true
    },
    {
      validator: async (rule, value) => {
        if (value != formState.password && value) {
          return Promise.reject('两次输入的密码不一致')
        }
        return Promise.resolve()
      }
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const add = () => {
  modalTitle.value = '添加团队成员'
  show.value = true
  resetFields()
  formState.user_name = ''
  formState.nick_name = ''
  formState.avatar = ''
  formState.avatar_url = defaultAvatar
  formState.user_roles = '3'
  formState.password = ''
  formState.check_password = ''
  formState.id = ''
  
}

const edit = (data) => {
  modalTitle.value = '编辑团队成员'
  formState.user_name = data.user_name
  formState.nick_name = data.nick_name
  formState.avatar_url = data.avatar || defaultAvatar;
  formState.avatar = ''
  formState.user_roles = data.user_roles
  formState.password = ''
  formState.check_password = ''
  formState.id = data.id
  delete formRules.password
  delete formRules.check_password
  show.value = true
}

const saveLoading = ref(false)
const handleOk = () => {
  validate().then(() => {
    let formData = {
      ...toRaw(formState)
    }
    delete formData.avatar_url
    console.log(formData, formState)
    saveLoading.value = true
    saveUser(formData)
      .then((res) => {
        message.success(`${modalTitle.value}成功`)
        show.value = false
        emit('ok')
      })
      .finally(() => {
        saveLoading.value = false
      })
  })
}

defineExpose({
  add,
  edit
})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
}
</style>
