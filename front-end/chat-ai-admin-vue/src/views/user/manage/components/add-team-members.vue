<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="746px">
      <div class="form-box">
        <div class="avatar-box">
          <AvatarInput v-model:value="formState.avatar_url" @change="onAvatarChange" />
          <div class="tip">{{ t('click_to_replace_tip') }}</div>
        </div>
        <a-form layout="vertical" autocomplete="off">
          <div class="flex-item-box">
            <a-form-item :label="t('member_nickname')" v-bind="validateInfos.nick_name">
              <a-input
                :maxlength="100"
                type="text"
                :placeholder="t('input_nickname')"
                v-model:value="formState.nick_name"
              ></a-input>
            </a-form-item>
            <a-form-item :label="t('login_account')" v-bind="validateInfos.user_name">
              <a-input
                :disabled="formState.id != ''"
                type="text"
                :placeholder="t('account_format_tip')"
                v-model:value="formState.user_name"
              ></a-input>
            </a-form-item>
          </div>
          <div class="flex-item-box">
            <a-form-item :label="t('member_role')" v-bind="validateInfos.user_roles">
              <a-select
                v-model:value="formState.user_roles"
                style="width: 100%"
                :placeholder="t('select_member_role')"
              >
                <a-select-option v-for="item in roleLists" :key="item.id" :value="item.id">{{
                  item.name
                }}</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item
              :label="t('select_department')"
              v-bind="validateInfos.department_ids"
              name="department_ids"
            >
              <a-tree-select
                v-model:value="formState.department_ids"
                show-search
                style="width: 100%"
                :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
                :placeholder="t('please_select')"
                allow-clear
                multiple
                tree-default-expand-all
                :tree-data="gData"
                tree-node-filter-prop="label"
              >
                <template #title="{ value: val, label }">
                  <div>{{ label }}</div>
                </template>
              </a-tree-select>
            </a-form-item>
          </div>
          <div class="flex-item-box">
            <a-form-item :label="t('login_password')" v-bind="validateInfos.password" v-if="!formState.id">
              <a-input-password
                v-model:value="formState.password"
                :placeholder="t('password_format_tip')"
              />
            </a-form-item>
            <a-form-item
              :label="t('confirm_password')"
              v-bind="validateInfos.check_password"
              v-if="!formState.id"
            >
              <a-input-password
                v-model:value="formState.check_password"
                :placeholder="t('reinput_password')"
              />
            </a-form-item>
          </div>
          <div class="flex-item-box">
            <a-form-item :label="t('expire_time')" v-bind="validateInfos.expire_time" required>
              <a-radio-group v-model:value="formState.expire_time_type">
                <a-radio :value="0">{{ t('permanent_valid') }}</a-radio>
                <a-radio :value="1">{{ t('specify_time') }}</a-radio>
              </a-radio-group>
              <div style="margin-top: 4px" v-if="formState.expire_time_type == 1">
                <a-date-picker
                  :showNow="false"
                  format="YYYY-MM-DD HH:mm"
                  :disabled-date="disabledDate"
                  v-model:value="formState.expire_time"
                  :show-time="{ defaultValue: dayjs('00:00', 'HH:mm') }"
                />
              </div>
            </a-form-item>
          </div>
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
import defaultAvatar from '@/assets/img/role_avatar.png'
import { getDepartmentList } from '@/api/department/index.js'
import { formateDepartmentCascaderData } from '@/utils/index.js'
import dayjs from 'dayjs'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.manage.add-team-members')
const emit = defineEmits(['ok'])

const useForm = Form.useForm
const show = ref(false)
const modalTitle = ref(t('add_team_member'))
const formState = reactive({
  user_name: '',
  nick_name: '',
  avatar: '',
  avatar_url: '',
  user_roles: '3',
  password: '',
  check_password: '',
  department_ids: [],
  expire_time_type: 0,
  expire_time: '',
  id: ''
})

const gData = ref([])
const getLists = () => {
  getDepartmentList({}).then((res) => {
    let treeData = res.data || []
    gData.value = formateDepartmentCascaderData(treeData)
  })
}

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
  roleLists.value = lists.filter((item) => item.id != '1' && item.role_type != 1)
})
const formRules = reactive({
  user_name: [
    {
      message: t('please_input_login_account'),
      required: true
    },
    {
      validator: async (rule, value) => {
        if (!/^[a-zA-Z0-9_.-]+$/.test(value) && value) {
          return Promise.reject(t('account_format_tip'))
        }
        return Promise.resolve()
      }
    }
  ],
  nick_name: [
    {
      message: t('input_nickname'),
      required: true
    }
  ],
  avatar_url: [
    {
      message: t('please_upload_avatar'),
      required: true
    }
  ],
  user_roles: [
    {
      message: t('please_select_member_role'),
      required: true
    }
  ],
  password: [
    {
      message: t('please_input_login_password'),
      required: true
    },
    {
      validator: async (rule, value) => {
        if (!validatePassword(value) && value) {
          return Promise.reject(t('password_format_tip'))
        }
        return Promise.resolve()
      }
    }
  ],
  check_password: [
    {
      message: t('please_input_confirm_password'),
      required: true
    },
    {
      validator: async (rule, value) => {
        if (value != formState.password && value) {
          return Promise.reject(t('passwords_not_match'))
        }
        return Promise.resolve()
      }
    }
  ],
  expire_time: [
    {
      validator: async (rule, value) => {
        if (formState.expire_time_type == 1 && !value) {
          return Promise.reject(t('please_select_specify_date'))
        }
        return Promise.resolve()
      }
    }
  ],
  department_ids: [
    {
      message: t('please_select_department'),
      required: true
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const add = (department_id) => {
  modalTitle.value = t('add_team_member')
  show.value = true
  resetFields()
  formState.user_name = ''
  formState.nick_name = ''
  formState.avatar = ''
  formState.avatar_url = defaultAvatar
  formState.user_roles = '3'
  formState.password = ''
  formState.check_password = ''
  formState.expire_time_type = 0
  formState.expire_time = ''
  formState.department_ids = department_id ? [department_id] : []
  formState.id = ''
  getLists()
}

const edit = (data) => {
  modalTitle.value = t('edit_team_member')
  formState.user_name = data.user_name
  formState.nick_name = data.nick_name
  formState.avatar_url = data.avatar || defaultAvatar
  formState.avatar = ''
  formState.user_roles = data.user_roles
  formState.password = ''
  formState.check_password = ''
  formState.id = data.id
  if (data.expire_time > 0) {
    formState.expire_time_type = 1
    formState.expire_time = dayjs(data.expire_time * 1000)
  } else {
    formState.expire_time_type = 0
    formState.expire_time = ''
  }

  formState.department_ids = data.departments.map((item) => +item.id)

  delete formRules.password
  delete formRules.check_password
  show.value = true
  getLists()
}

function setDepartmentIds(data) {
  let list = []
  data.forEach((item) => {
    let filterItem = findNodeByValue(gData.value, item.id)
    if (filterItem && filterItem.path) {
      list.push(filterItem.path)
    }
  })
  return list
}

function findNodeByValue(data, targetValue) {
  for (const node of data) {
    if (node.value == targetValue) {
      return { ...node } // 返回匹配的节点及其子节点
    }
    if (node.children && node.children.length > 0) {
      const foundNode = findNodeByValue(node.children, targetValue)
      if (foundNode) {
        return foundNode
      }
    }
  }
  return null // 如果没有找到匹配的节点，返回 null
}

const saveLoading = ref(false)
const handleOk = () => {
  validate().then(() => {
    let formData = {
      ...toRaw(formState)
    }

    if (formData.expire_time_type == 1) {
      formData.expire_time = dayjs(formData.expire_time).unix()
    } else {
      formData.expire_time = 0
    }

    formData.department_ids = formState.department_ids.join(',')
    delete formData.avatar_url
    console.log(formData, formState)
    saveLoading.value = true
    saveUser(formData)
      .then((res) => {
        message.success(`${modalTitle.value}${t('success')}`)
        show.value = false
        emit('ok')
      })
      .finally(() => {
        saveLoading.value = false
      })
  })
}

const disabledDate = (current) => {
  return current && current < dayjs().startOf('day')
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
.flex-item-box {
  display: flex;
  align-items: center;
  gap: 32px;
  .ant-form-item {
    flex: 1;
  }
}
.avatar-box {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 24px;
  .tip {
    color: #8c8c8c;
    line-height: 22px;
  }
}
</style>