<template>
  <div>
    <a-modal v-model:open="open" title="编辑知识库" @ok="handleOk" width="746px">
      <div class="form-box-wrapper">
        <div class="form-item">
          <div class="form-label required">知识库名称：</div>
          <div class="form-content">
            <a-input v-model:value="library_name" :maxlength="20" placeholder="请输入知识库名称" />
          </div>
        </div>
        <div class="form-item">
          <div class="form-label">知识库简介：</div>
          <div class="form-content">
            <a-textarea
              style="height: 100px"
              v-model:value="library_intro"
              placeholder="请输入知识库简介"
            />
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { editLibrary } from '@/api/library'
import { useRoute } from 'vue-router'
const rotue = useRoute()
const emit = defineEmits(['handleEditLibrary'])
const open = ref(false)
const library_name = ref('')
const library_intro = ref('')
const showModal = (data) => {
  library_name.value = data.library_name
  library_intro.value = data.library_intro
  open.value = true
}
const handleOk = () => {
  if (!library_name.value) {
    return message.error('请输入知识库名称')
  }
  let data = {
    library_name: library_name.value,
    library_intro: library_intro.value,
    id: rotue.query.id
  }
  editLibrary(data).then((res) => {
    message.success('修改成功')
    open.value = false
    emit('handleEditLibrary', data)
  })
}
defineExpose({ showModal })
</script>
<style lang="less" scoped>
.form-box-wrapper {
  .form-item {
    margin-top: 16px;
  }
  .form-label {
    color: #262626;
    font-size: 14px;
    line-height: 22px;
    padding-top: 5px;
    &.required::before {
      content: '*';
      display: inline-block;
      color: #fb363f;
      margin-right: 2px;
    }
  }
  .form-content {
    margin-top: 8px;
  }
}
</style>
