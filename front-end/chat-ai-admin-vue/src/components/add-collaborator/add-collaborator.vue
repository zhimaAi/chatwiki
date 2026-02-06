<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="t('title_add_collaborator')"
      @ok="handleOk"
      :width="746"
      wrapClassName="no-padding-modal"
      :bodyStyle="{ 'max-height': '500px', 'overflow-y': 'auto' }"
      :confirmLoading="saveLoading"
    >
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item :label="t('label_collaboration_permission')" required>
            <a-radio-group v-model:value="formState.operate_rights">
              <a-radio value="4">{{ t('btn_manage') }}</a-radio>
              <a-radio value="2">{{ t('btn_edit') }}</a-radio>
              <a-radio value="1">{{ t('btn_view') }}</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item :label="t('label_collaborator')" required>
            <a-radio-group v-model:value="formState.identity_type">
              <a-radio value="1">{{ t('btn_select_by_member') }}</a-radio>
              <a-radio value="2">{{ t('btn_select_by_org') }}</a-radio>
            </a-radio-group>
            <div class="form-tip" v-if="role_permission && role_permission.includes('TeamManage')">
              {{ t('msg_no_account_tip') }} <a @click="toManagePage">{{ t('btn_go_to_add') }}</a>
            </div>
          </a-form-item>
          <div class="list-box" v-if="formState.identity_type == 1">
            <div
              class="type-item"
              v-for="item in userLists"
              :key="item.value"
              @click="handleChangeUser(item)"
            >
              <div class="user-block">
                <a-avatar shape="square" :src="item.avatar" :size="40">
                  <template #icon><UserOutlined v-if="!item.avatar" /></template>
                </a-avatar>
                <div class="name-box">
                  <div class="name">{{ item.user_name }}</div>
                  <div class="remark-text">{{ item.nick_name }}</div>
                </div>
              </div>
              <svg-icon
                class="check-arrow"
                name="check-arrow-filled"
                style="font-size: 24px; color: #fff"
                v-if="userCheckList.includes(item.id)"
              ></svg-icon>
            </div>
          </div>
          <div class="checkbox-list" v-if="formState.identity_type == 2">
            <a-checkbox-group v-model:value="departmentCheckList" style="width: 100%">
              <div class="checklist-item">
                <a-checkbox v-for="item in departmentLists" :key="item.id" :value="item.id">
                  <div>{{ item.department_name }}</div>
                </a-checkbox>
              </div>
            </a-checkbox-group>
          </div>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { Form, message } from 'ant-design-vue'
import { UserOutlined } from '@ant-design/icons-vue'
import { getUserList } from '@/api/manage/index.js'
import { getAllDepartment, savePermissionManage } from '@/api/department/index.js'
import { usePermissionStore } from '@/stores/modules/permission'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('components.add-collaborator.add-collaborator')

const role_permission = computed(() => usePermissionStore().role_permission || [])

const emit = defineEmits(['ok'])
const useForm = Form.useForm
const open = ref(false)

const props = defineProps({
  libraryInfo: {
    type: Object,
    default: () => {}
  }
})

const formState = reactive({
  identity_type: '1',
  operate_rights: '4'
})

const userCheckList = ref([])
const departmentCheckList = ref([])

const handleChangeUser = (item) => {
  if (item.role_type == 1) {
    // if (!userCheckList.value.includes(item.id)) {
    //   userCheckList.value.push(item.id)
    // }
    return
  }
  if (userCheckList.value.includes(item.id)) {
    userCheckList.value = userCheckList.value.filter((id) => id != item.id)
  } else {
    userCheckList.value.push(item.id)
  }
}
const formRules = reactive({})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

let object_array = []

let tempCheckList = []
const show = (data, list) => {
  getAllDepartList()
  userCheckList.value = []
  departmentCheckList.value = []
  tempCheckList = []
  data.forEach((item) => {
    if (item.identity_type == 1) {
      // 用户
      // userCheckList.value.push(item.identity_id)
      tempCheckList.push(item.identity_id)
    }
    if (item.identity_type == 2) {
      // 部门
      departmentCheckList.value.push(item.identity_id)
    }
  })
  getAllUserList()
  formState.operate_rights = '4'
  formState.identity_type = '1'
  object_array = list
  open.value = true
}

const saveLoading = ref(false)
const handleOk = () => {
  object_array = object_array.map((item) => {
    return {
      object_id: +item.object_id,
      object_type: +item.object_type,
      operate_rights: +formState.operate_rights
    }
  })
  let identity_ids = []
  if (formState.identity_type == 1) {
    identity_ids = userCheckList.value
  } else {
    identity_ids = departmentCheckList.value
  }
  if (identity_ids.length == 0) {
    message.error(t('msg_select_member_or_org'))
    return
  }
  let parmas = {
    identity_type: formState.identity_type,
    object_array: JSON.stringify(object_array),
    identity_ids: identity_ids.join(',')
  }
  saveLoading.value = true
  savePermissionManage(parmas)
    .then((res) => {
      message.success(t('msg_save_success'))
      open.value = false
      emit('ok')
    })
    .finally(() => {
      saveLoading.value = false
    })
}

const userLists = ref([])
const getAllUserList = () => {
  getUserList({
    page: 1,
    size: 200
  }).then((res) => {
    // userLists.value = res.data.list
    userLists.value = res.data.list.filter(item => !tempCheckList.includes(item.id))
  })
}

const departmentLists = ref([])
const getAllDepartList = () => {
  getAllDepartment().then((res) => {
    let data = res.data || []
    departmentLists.value = data
  })
}

const toManagePage = () => {
  window.open(`/#/user/manage`)
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
  padding-right: 16px;
  .ant-form-item {
    margin-bottom: 16px;
  }
  ::v-deep(.ant-form-item-label) {
    padding: 0 0 4px;
  }
  .form-tip {
    color: #8c8c8c;
    font-size: 14px;
    line-height: 22px;
    font-weight: 400;
    margin-top: 4px;
  }
}

.list-box {
  display: flex;
  margin-top: -8px;
  flex-wrap: wrap;
  gap: 8px;
  .type-item {
    position: relative;
    width: 100%;
    cursor: pointer;
    padding: 16px;
    display: flex;
    border: 1px solid #e8e8e8;
    border-radius: 6px;
    box-shadow: none;
    transition: box-shadow 1s;
    width: calc(50% - 4px);
    &:hover {
      box-shadow:
        0px 5px 5px -3px rgba(0, 0, 0, 0.1),
        0px 8px 10px 1px rgba(0, 0, 0, 0.06),
        0px 3px 14px 2px rgba(0, 0, 0, 0.05);
    }
    &.active {
      border: 2px solid #2475fc;
    }
    .check-arrow {
      position: absolute;
      bottom: 0;
      right: -1px;
    }
  }
  .user-block {
    display: flex;
    align-items: center;
    gap: 8px;
    .name-box {
      flex: 1;
      color: #262626;
      line-height: 22px;
      font-size: 14px;
    }
    .remark-text {
      margin-top: 2px;
      font-size: 12px;
      color: #8c8c8c;
      line-height: 20px;
    }
  }
}

.checklist-item {
  display: flex;
  width: 100%;
  flex-direction: column;
  .ant-checkbox-wrapper {
    width: 100%;
    padding: 12px 16px;
    border-bottom: 1px solid #e8e8e8;
    ::v-deep(.ant-checkbox + span) {
      flex: 1;
    }
  }
}
</style>
