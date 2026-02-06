<template>
  <div>
    <a-modal v-model:open="open" :title="t('modal_title')" @ok="handleOk">
      <div class="form-item">
        <div class="form-label">{{ t('form_label') }}</div>
        <div class="form-content">
          <a-select v-model:value="formState.group_id" style="width: 100%" :placeholder="t('select_placeholder')">
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
import { relationLibraryGroup, getLibraryListGroup } from '@/api/library'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'
const { t } = useI18n('views.library.library-list.components.edit-group')
const open = ref(false)
const emit = defineEmits(['ok'])
const modalTitle = ref(t('modal_title'))
const formState = reactive({
  group_id: '',
  library_id: ''
})

const groupLists = ref([])

const show = (data) => {
  formState.group_id = data.group_id || '0'
  formState.library_id = data.id || ''
  getGroupList()
  open.value = true
}
const handleOk = () => {
  if (!formState.group_id) {
    return message.error(t('error_message'))
  }
  relationLibraryGroup({
    ...formState
  }).then((res) => {
    message.success(`${t('modal_title')}成功`)
    open.value = false
    emit('ok')
  })
}
const getGroupList = () => {
  getLibraryListGroup().then((res) => {
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
