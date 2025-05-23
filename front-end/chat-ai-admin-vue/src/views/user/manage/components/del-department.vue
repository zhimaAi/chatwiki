<template>
  <div>
    <a-modal v-model:open="open" :title="null" :footer="null" :width="480" :closable="false">
      <div class="del-box">
        <div class="icon-box"><ExclamationCircleFilled /></div>
        <div class="top-title">删除部门</div>
        <div class="content">
          确认删除部门么？直接删除后，该部门成员归属到默认部门或设置新的部门
        </div>
        <div class="footer-btn">
          <a-button @click="open = false">关闭</a-button>
          <a-button danger @click="handleDirectDel">直接删除</a-button>
          <a-button type="primary" @click="handleOpenSetModal">删除并设置</a-button>
        </div>
      </div>
    </a-modal>
    <a-modal v-model:open="setModal" title="设置部门" :width="472" @ok="handleDel">
      <div class="set-box">
        <div class="set-desc">删除之前请先将之前部门的成员换到其他团队</div>
        <div class="set-item">
          <div class="set-label">所属部门</div>
          <div>
            <a-select
              v-model:value="formState.new_department_id"
              style="width: 100%"
              placeholder="请选择"
            >
              <a-select-option v-for="item in departmentLists" :key="item.id" :value="item.id">{{
                item.department_name
              }}</a-select-option>
            </a-select>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { reactive, ref, createVNode } from 'vue'
import { ExclamationCircleFilled, ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getAllDepartment, deleteDepartment } from '@/api/department/index'
import { message, Modal } from 'ant-design-vue'

const open = ref(false)
const setModal = ref(false)

const emit = defineEmits(['ok'])
const formState = reactive({
  id: '',
  new_department_id: void 0
})

let filterIds = []
const show = (data) => {
  // 过滤掉自己 和子辈
  filterIds = []
  filterIds = getSonIds(data)

  formState.new_department_id = void 0
  formState.id = data.id
  if (data.children.length == 0) {
    Modal.confirm({
      title: '提示',
      icon: createVNode(ExclamationCircleOutlined),
      content: '确定删除该部门',
      onOk() {
        handleDel()
      }
    })
    return
  }
  open.value = true
  getAllList()
}

function getSonIds(data) {
  let ids = [data.id]

  if (data.children && data.children.length > 0) {
    for (const child of data.children) {
      ids = ids.concat(getSonIds(child))
    }
  }

  return ids
}

const handleOpenSetModal = () => {
  setModal.value = true
}

const handleDirectDel = () => {
  formState.new_department_id = ''
  handleDel()
}
const handleDel = () => {
  deleteDepartment({
    ...formState
  }).then((res) => {
    setModal.value = false
    open.value = false
    message.success('删除成功')
    emit('ok', formState.id)
  })
}

const departmentLists = ref([])
const getAllList = () => {
  getAllDepartment().then((res) => {
    let data = res.data || []
    departmentLists.value = data.filter((item) => !filterIds.includes(+item.id))
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.del-box {
  position: relative;
  padding-left: 38px;
  padding-top: 12px;
}
.icon-box {
  position: absolute;
  top: 7px;
  left: 0;
  font-size: 20px;
  color: #ff9900;
}
.top-title {
  font-size: 16px;
  font-weight: 600;
  line-height: 24px;
  color: #262626;
}

.content {
  color: #595959;
  line-height: 22px;
  font-size: 14px;
  margin-top: 8px;
}
.footer-btn {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 24px;
}

.set-box {
  margin-top: 24px;
  font-size: 14px;
  line-height: 22px;
  color: #595959;
  margin-bottom: 24px;
  .set-item {
    margin-top: 24px;
    .set-label {
      color: #262626;
      margin-bottom: 4px;
    }
  }
}
</style>
