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
import { relationLibraryGroup, getLibraryGroup } from '@/api/library'
import { message } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.qa-knowledge-document.components.edit-group')
const open = ref(false)
const emit = defineEmits(['ok'])
const modalTitle = ref('title_edit_group')

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
    return message.error(t('msg_select_group'))
  }
  relationLibraryGroup({
    ...formState,
    library_id: props.libraryId,
    sense: props.sense,
  }).then((res) => {
    message.success(`${t(modalTitle.value)}${t('msg_operation_success')}`)
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
