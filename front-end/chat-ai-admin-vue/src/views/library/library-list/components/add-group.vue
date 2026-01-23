<template>
  <div>
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk" :zIndex="2000">
      <div class="form-item">
        <div class="form-label">{{ t('group_name_label') }}</div>
        <div class="form-content">
          <a-input
            :maxLength="15"
            v-model:value="formState.group_name"
            :placeholder="t('group_name_placeholder')"
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
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-list.components.add-group')

const open = ref(false)
const emit = defineEmits(['ok'])
const modalTitle = ref(t('create_group_title'))
const formState = reactive({
  group_name: '',
  id: ''
})
const show = (data) => {
  formState.group_name = data.group_name || ''
  formState.id = data.id || ''
  modalTitle.value = data.id ? t('edit_group_title') : t('create_group_title')
  open.value = true
}
const handleOk = () => {
  if (!formState.group_name) {
    return message.error(t('group_name_required'))
  }
  saveLibraryListGroup({
    ...formState
  }).then((res) => {
    message.success(t('operation_success', { title: modalTitle.value }))
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
