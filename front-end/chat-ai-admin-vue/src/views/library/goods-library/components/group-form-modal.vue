<template>
  <a-modal
    v-model:open="open"
    :title="modalTitle"
    :ok-text="okText"
    :cancel-text="t('confirm.cancel_btn')"
    :confirm-loading="submitting"
    @ok="handleOk"
    :destroyOnClose="true"
    :width="560"
  >
    <a-form layout="vertical" class="group-form">
      <a-form-item :label="t('group_modal.parent_label')">
        <a-tree-select
          v-model:value="formState.parent_id"
          :tree-data="parentTreeOptions"
          :placeholder="t('group_modal.root_group')"
          show-search
          tree-default-expand-all
          tree-node-filter-prop="title"
          style="width: 100%"
        />
      </a-form-item>
      <a-form-item :label="t('group_modal.name_label')" required>
        <a-input
          v-model:value="formState.group_name"
          :placeholder="t('group_modal.name_placeholder')"
          :maxlength="15"
          show-count
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.goods-library.index')

const emit = defineEmits(['ok'])

const open = ref(false)
const submitting = ref(false)
const parentTreeOptions = ref([])

const formState = reactive({
  id: '',
  parent_id: '0',
  parent_name: '',
  group_name: '',
  level: 0
})

const modalTitle = ref(t('group_modal.title_add'))

const okText = computed(() => {
  return formState.id ? t('group_modal.save_edit') : t('group_modal.save_add')
})

const show = (data = {}) => {
  formState.id = data.id || ''
  formState.parent_id = String(data.parent_id ?? data.parentId ?? '0')
  formState.parent_name = data.parent_name || data.parentName || ''
  formState.group_name = data.group_name || ''
  formState.level = data.level ?? 0
  parentTreeOptions.value = data.parentTreeOptions || []
  modalTitle.value = data.id ? t('group_modal.title_edit') : t('group_modal.title_add')
  submitting.value = false
  open.value = true
}

const close = () => {
  submitting.value = false
  open.value = false
}

const handleOk = () => {
  if (!formState.group_name.trim()) {
    return message.error(t('validation.group_name_required'))
  }

  if (submitting.value) {
    return
  }

  submitting.value = true

  emit(
    'ok',
    {
      ...formState
    },
    {
      close,
      setSubmitting: (value) => {
        submitting.value = value
      }
    }
  )
}

defineExpose({
  show,
  close
})
</script>

<style lang="less" scoped>
.group-form {
  margin-top: 4px;
}
</style>
