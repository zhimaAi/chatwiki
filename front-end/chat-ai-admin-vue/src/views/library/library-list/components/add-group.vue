<template>
  <div>
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :zIndex="2000">
      <div class="form-item">
        <div class="form-label">分组名称：</div>
        <div class="form-content">
          <a-input
            :maxLength="15"
            v-model:value="formState.group_name"
            placeholder="请输入分组名称"
          ></a-input>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { saveLibraryListGroup } from '@/api/library'
import { message } from 'ant-design-vue'
const open = ref(false)
const emit = defineEmits(['ok'])
const modalTitle = ref('新建分组')
const formState = reactive({
  group_name: '',
  id: ''
})
const show = (data) => {
  formState.group_name = data.group_name || ''
  formState.id = data.id || ''
  modalTitle.value = data.id ? '编辑分组' : '新建分组'
  open.value = true
}
const handleOk = () => {
  if (!formState.group_name) {
    return message.error('请输入分组名称')
  }
  saveLibraryListGroup({
    ...formState
  }).then((res) => {
    message.success(`${modalTitle.value}成功`)
    open.value = false
    emit('ok')
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.form-item {
  display: flex;
  align-items: center;
  margin: 24px 0;
  .form-content {
    flex: 1;
  }
}
</style>
