<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="746px">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item label="角色名称" v-bind="validateInfos.name">
            <a-input
              type="text"
              :maxlength="100"
              placeholder="请输入角色名称"
              v-model:value="formState.name"
            ></a-input>
          </a-form-item>
          <a-form-item label="角色备注">
            <a-textarea v-model:value="formState.mark" placeholder="请输入知识库介绍" />
          </a-form-item>
          <a-form-item label="角色权限">
            <div class="role-check-box">
              <a-flex class="title-boock" justify="space-between">
                <div class="title-row">机器人管理</div>
                <div class="check-num">
                  {{ formState.robotChecked.length }}/{{ robotOptions.length }}
                </div>
              </a-flex>
              <a-checkbox-group v-model:value="formState.robotChecked" style="width: 100%">
                <a-row :gutter="[0, 12]">
                  <a-col :span="6" v-for="item in robotOptions" :key="item.value">
                    <a-checkbox :value="item.value">{{ item.label }}</a-checkbox>
                  </a-col>
                </a-row>
              </a-checkbox-group>
            </div>
            <div class="role-check-box">
              <a-flex class="title-boock" justify="space-between">
                <div class="title-row">知识库管理</div>
                <div class="check-num">
                  {{ formState.libraryChecked.length }}/{{ libraryOptions.length }}
                </div>
              </a-flex>
              <a-checkbox-group v-model:value="formState.libraryChecked" style="width: 100%">
                <a-row :gutter="[0, 12]">
                  <a-col :span="6" v-for="item in libraryOptions" :key="item.value">
                    <a-checkbox :value="item.value">{{ item.label }}</a-checkbox>
                  </a-col>
                </a-row>
              </a-checkbox-group>
            </div>
            <div class="role-check-box">
              <a-flex class="title-boock" justify="space-between">
                <div class="title-row">系统设置</div>
                <div class="check-num">
                  {{ formState.systemChecked.length }}/{{ systemOptions.length }}
                </div>
              </a-flex>
              <a-checkbox-group v-model:value="formState.systemChecked" style="width: 100%">
                <a-row :gutter="[0, 12]">
                  <a-col :span="6" v-for="item in systemOptions" :key="item.value">
                    <a-checkbox :value="item.value">{{ item.label }}</a-checkbox>
                  </a-col>
                </a-row>
              </a-checkbox-group>
            </div>
            <div class="role-check-box">
              <a-flex class="title-boock" justify="space-between">
                <div class="title-row">客户端管理</div>
                <div class="check-num">
                  {{ formState.clientChecked.length }}/{{ clientOptions.length }}
                </div>
              </a-flex>
              <a-checkbox-group v-model:value="formState.clientChecked" style="width: 100%">
                <a-row :gutter="[0, 12]">
                  <a-col :span="6" v-for="item in clientOptions" :key="item.value">
                    <a-checkbox :value="item.value">{{ item.label }}</a-checkbox>
                  </a-col>
                </a-row>
              </a-checkbox-group>
            </div>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, toRaw } from 'vue'
import { Form, message } from 'ant-design-vue'
import { getRole, getMenu, saveRole } from '@/api/manage/index.js'
const emit = defineEmits(['ok'])
const useForm = Form.useForm

let robotOptions = []
let libraryOptions = []
let systemOptions = []
let clientOptions = []

getMenu().then((res) => {
  robotOptions = [
    {
      label: res.data[1].name,
      value: res.data[1].uni_key
    }
  ]
  libraryOptions = [
    {
      label: res.data[2].name,
      value: res.data[2].uni_key
    }
  ]
  systemOptions = [
    {
      label: res.data[3].name,
      value: res.data[3].uni_key
    }
  ]
  clientOptions = [
    {
      label: res.data[4].name,
      value: res.data[4].uni_key
    }
  ]
})

const show = ref(false)
const modalTitle = ref('添加角色')
const id = ref('')
const formState = reactive({
  name: '',
  mark: '',
  robotChecked: [],
  libraryChecked: [],
  systemChecked: [],
  clientChecked: [],
  role_type: ''
})

const formRules = reactive({
  name: [
    {
      message: '请输入角色名称',
      required: true
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const add = () => {
  modalTitle.value = '添加角色'
  id.value = ''
  show.value = true
  formState.name = ''
  formState.mark = ''
  resetFields()
}

const edit = (record) => {
  modalTitle.value = '编辑角色'
  id.value = record.id
  getRole({ id: record.id }).then((res) => {
    let data = res.data
    let role_permission = data.role_permission || []
    formState.name = data.name
    formState.mark = data.mark
    formState.role_type = data.role_type
    formState.robotChecked = formatCheckList(robotOptions, role_permission)
    formState.libraryChecked = formatCheckList(libraryOptions, role_permission)
    formState.systemChecked = formatCheckList(systemOptions, role_permission)
    formState.clientChecked = formatCheckList(clientOptions, role_permission)
    show.value = true
  })
}

const formatCheckList = (data, list) => {
  let resultList = []
  data.forEach((item) => {
    if (list.includes(item.value)) {
      resultList.push(item.value)
    }
  })
  return resultList
}

const handleOk = () => {
  validate().then(() => {
    let uni_keys = [
      ...formState.robotChecked,
      ...formState.libraryChecked,
      ...formState.systemChecked,
      ...formState.clientChecked
    ]
    let parmas = {
      id: id.value,
      name: formState.name,
      mark: formState.mark,
      role_type: formState.role_type,
      uni_keys: uni_keys.join(',')
    }

    saveRole(parmas).then((res) => {
      show.value = false
      message.success('保存成功')
      emit('ok')
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
.role-check-box {
  background: #f2f4f7;
  border-radius: 6px;
  padding: 8px 16px 16px 12px;
  margin-bottom: 8px;
  .title-boock {
    margin-bottom: 8px;
  }
  .title-row {
    color: #262626;
    font-size: 14px;
    font-style: normal;
    font-weight: 600;
    line-height: 22px;
  }
  .check-num {
    color: #8c8c8c;
    font-size: 14px;
    font-weight: 400;
    line-height: 22px;
  }
}
</style>
