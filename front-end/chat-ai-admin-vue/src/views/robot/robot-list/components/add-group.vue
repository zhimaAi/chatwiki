<template>
  <div>
    <a-modal v-model:open="open" :title="t(modalTitle)" @ok="handleOk" :zIndex="2000">
      <div class="form-item">
        <div class="form-label">{{ t('label_group_name') }}：</div>
        <div class="form-content">
          <a-input
            :maxLength="15"
            v-model:value="formState.group_name"
            :placeholder="t('ph_input_group_name')"
          ></a-input>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { saveRobotGroup } from '@/api/robot'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-list.components.add-group')

const open = ref(false)
const emit = defineEmits(['ok'])
const modalTitle = ref('title_create_group')
const formState = reactive({
  group_name: '',
  id: ''
})
const show = (data) => {
  formState.group_name = data.group_name || ''
  formState.id = data.id || ''
  modalTitle.value = data.id ? 'title_edit_group' : 'title_create_group'
  open.value = true
}
const handleOk = () => {
  if (!formState.group_name) {
    return message.error(t('msg_group_name_required'))
  }
  saveRobotGroup({
    ...formState
  }).then((res) => {
    message.success(t('msg_operation_success', { operation: t(modalTitle.value) }))
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
