<template>
  <div>
    <a-modal
      v-model:open="open"
      title="批量修改部门"
      @ok="handleOk"
      :confirmLoading="saveLoading"
    >
      <a-form layout="vertical" style="margin-top: 24px;">
        <a-form-item label="选择部门" v-bind="validateInfos.department_ids">
          <a-tree-select
            v-model:value="formState.department_ids"
            show-search
            style="width: 100%"
            :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
            placeholder="请选择"
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
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Form, message } from 'ant-design-vue'
import { getDepartmentList, batchUpdateUserDepartment } from '@/api/department/index.js'
import { formateDepartmentCascaderData } from '@/utils/index.js'

const open = ref(false)
const emit = defineEmits(['ok'])

const useForm = Form.useForm
const formState = reactive({
  user_ids: '',
  department_ids: []
})

const formRules = reactive({
  department_ids: [
    {
      message: '请选择部门',
      required: true
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

const show = (ids) => {
  formState.user_ids = ids
  formState.department_ids = []
  open.value = true
  getLists()
}

const gData = ref([])
const getLists = () => {
  getDepartmentList({}).then((res) => {
    let treeData = res.data || []
    gData.value = formateDepartmentCascaderData(treeData)
  })
}

const saveLoading = ref(false)
const handleOk = () => {
  validate().then(() => {
    let formData = {
      ...formState
    }

    formData.department_ids = formState.department_ids.join(',')
    saveLoading.value = true
    batchUpdateUserDepartment(formData)
      .then((res) => {
        message.success(`批量修改部门成功`)
        open.value = false
        emit('ok')
      })
      .finally(() => {
        saveLoading.value = false
      })
  })
}

defineExpose({
  show
})
</script>

<style lang="scss" scoped></style>
