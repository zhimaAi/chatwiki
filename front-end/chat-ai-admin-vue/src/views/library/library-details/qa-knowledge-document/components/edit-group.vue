<template>
  <div>
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk">
      <div class="form-item">
        <div class="form-label">分组名称：</div>
        <div class="form-content">
          <a-select v-model:value="formState.group_id" style="width: 100%" placeholder="请选择分组">
            <a-select-option v-for="item in groupLists" :value="item.id">{{
              item.group_name
            }}</a-select-option>
          </a-select>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { relationLibraryGroup, getLibraryGroup } from '@/api/library'
import { message } from 'ant-design-vue'
const open = ref(false)
const emit = defineEmits(['ok'])
const modalTitle = ref('修改分组')

const props = defineProps({
  libraryId: {
    type: [Number, String],
    default: ''
  },
  sense: {
    type: [Number, String],
    default: 0
  }
})

const formState = reactive({
  group_id: '',
  file_id: '',
})

const groupLists = ref([])

const show = (data) => {
  formState.group_id = data.group_id || '0'
  formState.file_id = data.id || ''
  getGroupList()
  open.value = true
}
const handleOk = () => {
  if (!formState.group_id) {
    return message.error('请选择分组')
  }
  relationLibraryGroup({
    ...formState,
    library_id: props.libraryId,
    sense: props.sense,
  }).then((res) => {
    message.success(`${modalTitle.value}成功`)
    open.value = false
    emit('ok')
  })
}
const getGroupList = () => {
  getLibraryGroup({
    library_id: props.libraryId,
    group_type: 1
  }).then((res) => {
    groupLists.value = res.data || []
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
