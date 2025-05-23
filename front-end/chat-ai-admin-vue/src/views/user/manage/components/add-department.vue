<template>
  <div>
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :width="472">
      <div class="form-box">
        <div class="form-item">
          <div class="form-label">部门名称</div>
          <div class="form-content">
            <a-input
              :maxLength="10"
              v-model:value="formState.department_name"
              placeholder="请输入部门名称"
            />
          </div>
        </div>
        <div class="form-item" v-if="!formState.id">
          <div class="form-label">所属部门</div>
          <div class="form-content">
            <a-tree-select
              v-model:value="formState.pid"
              show-search
              style="width: 100%"
              :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
              placeholder="请选择"
              allow-clear
              tree-default-expand-all
              :tree-data="gData"
              tree-node-filter-prop="label"
            >
              <template #title="{ value: val, label }">
                <div>{{ label }}</div>
              </template>
            </a-tree-select>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { saveDepartment, getDepartmentList } from '@/api/department/index.js'
import { formateDepartmentCascaderData } from '@/utils/index.js'
import { message } from 'ant-design-vue'
const open = ref(false)
const emit = defineEmits(['ok'])
const formState = reactive({
  department_name: '',
  id: '',
  pid: ''
})

const modalTitle = ref('')
const show = (data) => {
  formState.id = ''
  formState.pid = data.id > 0 ? data.id || '' : void 0
  formState.department_name = ''
  modalTitle.value = '添加部门'
  open.value = true
  getLists()
}

const rename = (data) => {
  formState.id = data.id || ''
  formState.pid = data.pid || ''
  formState.department_name = data.title || ''
  modalTitle.value = '重命名'
  open.value = true
}
const gData = ref([])
const getLists = () => {
  getDepartmentList({}).then((res) => {
    let treeData = res.data || []
    gData.value = formateDepartmentCascaderData(treeData)
  })
}

const handleOk = () => {
  if (!formState.department_name) {
    return message.error('请输入部门名称')
  }
  if (!formState.id && !formState.pid) {
    return message.error('请选择所属部门')
  }
  saveDepartment({
    ...formState
  }).then((res) => {
    message.success('修改成功')
    open.value = false
    emit('ok')
  })
}

defineExpose({
  show,
  rename
})
</script>

<style lang="less" scoped>
.form-box {
  margin-top: 24px;
  .form-item {
    margin-bottom: 24px;
    line-height: 22px;
    font-size: 14px;
    .form-label {
      color: #262626;
      margin-bottom: 4px;
      &::before {
        content: '*';
        color: #ff0000;
        margin-right: 4px;
        font-size: 14px;
      }
    }
  }
}
</style>
