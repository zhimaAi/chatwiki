<template>
  <div>
    <a-modal v-model:open="open" :title="t(modalTitle)" @ok="handleOk">
      <div class="form-item">
        <div class="form-label">{{ t('label_group_name') }}：</div>
        <div class="form-content">
          <a-select v-model:value="formState.group_id" style="width: 100%" :placeholder="t('ph_select_group')">
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
import { relationRobotGroup, getRobotGroupList } from '@/api/robot'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-list.components.edit-group')

const open = ref(false)
const emit = defineEmits(['ok'])
const modalTitle = ref('title_edit_group')
const formState = reactive({
  group_id: '',
  robot_id: ''
})

const groupLists = ref([])

const show = (data) => {
  formState.group_id = data.group_id || '0'
  formState.robot_id = data.id || ''
  getGroupList()
  open.value = true
}
const handleOk = () => {
  if (!formState.group_id) {
    return message.error(t('msg_select_group'))
  }
  relationRobotGroup({
    ...formState
  }).then((res) => {
    message.success(t('msg_operation_success', { title: t(modalTitle.value) }))
    open.value = false
    emit('ok')
  })
}
const getGroupList = () => {
  getRobotGroupList().then((res) => {
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
