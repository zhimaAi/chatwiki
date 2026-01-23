<template>
  <div>
    <a-modal
      v-model:open="open"
      :title="t('batch_edit_department')"
      @ok="handleOk"
      :confirmLoading="saveLoading"
    >
      <a-form layout="vertical" style="margin-top: 24px;">
        <a-form-item :label="t('select_department')" v-bind="validateInfos.department_ids">
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
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { Form, message } from 'ant-design-vue'
import { getDepartmentList, batchUpdateUserDepartment } from '@/api/department/index.js'
import { formateDepartmentCascaderData } from '@/utils/index.js'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.manage.components.batch-edit-department')

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
      message: t('please_select_department'),
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
        message.success(t('batch_edit_department_success'))
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