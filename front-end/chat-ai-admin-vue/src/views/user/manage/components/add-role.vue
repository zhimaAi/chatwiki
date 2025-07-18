<template>
  <div>
    <a-modal v-model:open="show" :title="modalTitle" @ok="handleOk" width="746px">
      <div class="form-box">
        <a-form layout="vertical">
          <a-form-item label="角色名称" v-bind="validateInfos.name">
            <a-input
              type="text"
              :maxlength="100"
              :disabled="formState.role_type > 0"
              placeholder="请输入角色名称"
              v-model:value="formState.name"
            ></a-input>
          </a-form-item>
          <a-form-item label="角色备注">
            <a-textarea v-model:value="formState.mark" placeholder="请输入角色备注" />
          </a-form-item>
          <a-form-item label="角色权限">
            <div class="role-check-box" v-for="item in menuOptions" :key="item.uni_key">
              <a-flex class="title-boock" justify="space-between">
                <div class="title-row">{{ item.name }}</div>
                <div class="check-num">
                  {{ robotChecked[item.uni_key].length }}/{{ getChildrenTotal(item.children) }}
                </div>
              </a-flex>
              <a-checkbox-group v-model:value="robotChecked[item.uni_key]" style="width: 100%">
                <a-row :gutter="[0, 12]">
                  <a-col
                    :span="sub.children && sub.children.length ? 12 : 6"
                    v-for="sub in item.children"
                    :key="sub.uni_key"
                  >
                    <a-checkbox
                      @change="handleSubChaneg(sub, item.uni_key)"
                      :disabled="formState.role_type > 0"
                      :value="sub.uni_key"
                      >{{ sub.name }}</a-checkbox
                    >
                    <span v-if="sub.children && sub.children.length > 0"
                      >(
                      <a-checkbox
                        @change="handleSonChange(son, sub, item.uni_key)"
                        v-for="son in sub.children"
                        :key="son.uni_key"
                        :disabled="formState.role_type > 0"
                        :value="son.uni_key"
                        >{{ son.name }}</a-checkbox
                      >
                      )</span
                    >
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
import { ref, reactive, onMounted, nextTick } from 'vue'
import { Form, message } from 'ant-design-vue'
import { getRole, getMenu, saveRole } from '@/api/manage/index.js'
const emit = defineEmits(['ok'])
const useForm = Form.useForm

const menuOptions = ref([])
const robotChecked = ref({})
const show = ref(false)
const modalTitle = ref('添加角色')
const id = ref('')
const formState = reactive({
  name: '',
  mark: '',
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

const getMenuData = async () => {
  await getMenu().then((res) => {
    let data = res.data || []
    // 任何成员 都有账号设置页面
    data.forEach((item) => {
      if (item.uni_key == 'System') {
        item.children = item.children.filter((it) => it.uni_key != 'AccountManage')
      }
    })
    menuOptions.value = data
    for (let i = 0; i < data.length; i++) {
      const item = data[i]
      robotChecked.value[item.uni_key] = []
    }
  })
}

const add = async () => {
  await getMenuData()
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
    robotChecked.value = formatCheckList(menuOptions.value, role_permission)
    show.value = true
  })
}

const formatCheckList = (data, list) => {
  let resultList = {}
  data.forEach((item) => {
    resultList[item.uni_key] = []
    for (let i = 0; i < item.children.length; i++) {
      const sub = item.children[i]
      if (list.includes(sub.uni_key)) {
        resultList[item.uni_key].push(sub.uni_key)
      }
      if (sub.children && sub.children.length > 0) {
        sub.children.forEach((it) => {
          if (list.includes(it.uni_key)) {
            resultList[item.uni_key].push(it.uni_key)
          }
        })
      }
    }
  })
  return resultList
}

const deconstruction = (obj) => {
  let newArr = []
  for (let key in obj) {
    const ele = obj[key]
    ele.map((item) => {
      newArr.push(item)
    })
  }
  return newArr
}

const handleOk = () => {
  validate().then(() => {
    let uni_keys = deconstruction(robotChecked.value)
    let parmas = {
      id: id.value,
      mark: formState.mark
    }

    if (formState.role_type == 0) {
      parmas.name = formState.name
      parmas.role_type = formState.role_type
      parmas.uni_keys = uni_keys.join(',')
    }

    saveRole(parmas).then((res) => {
      show.value = false
      message.success('保存成功')
      emit('ok')
    })
  })
}

const handleSubChaneg = (sub, uni_key) => {
  if (!sub.children) {
    return
  }
  nextTick(() => {
    if (!robotChecked.value[uni_key].includes(sub.uni_key)) {
      // 父级取消勾选 要把子级取消勾选
      let sonKey = sub.children.map((item) => item.uni_key)
      robotChecked.value[uni_key] = robotChecked.value[uni_key].filter(
        (item) => !sonKey.includes(item)
      )
    }
  })
}
const handleSonChange = (son, sub, uni_key) => {
  nextTick(() => {
    if (robotChecked.value[uni_key].includes(son.uni_key)) {
      // 孙子勾选  父级 勾选上去
      robotChecked.value[uni_key].push(sub.uni_key)
      // 数组去重
      robotChecked.value[uni_key] = [...new Set(robotChecked.value[uni_key])]
    }
  })
}

const getChildrenTotal = (data) => {
  let total = data.length
  data.forEach((item) => {
    if (item.children && item.children.length) {
      total = total + item.children.length
    }
  })
  return total
}

onMounted(() => {
  getMenuData()
})

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
  &::v-deep(.ant-checkbox + span) {
    padding-right: 0;
  }
}
</style>
