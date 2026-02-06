<template>
  <div>
    <a-modal v-model:open="open" :title="modalTitle" @ok="handleOk">
      <div class="form-item">
        <div class="form-label">{{ t('label_group_name') }}</div>
        <div class="form-content">
          <a-input :maxLength="15" v-model:value="formState.group_name" :placeholder="t('ph_input_group_name')"></a-input>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { saveLibraryGroup } from '@/api/library'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.qa-knowledge-document.components.add-group')

const query = useRoute().query
const open = ref(false)
const emit = defineEmits(['ok'])
const isEdit = ref(false)
const modalTitle = computed(() => isEdit.value ? t('title_edit_group') : t('title_create_group'))

const props = defineProps({
  group_type:{
    type: [Number, String],
    default: 0
  }
})

const formState = reactive({
  group_name: '',
  id: '',
  library_id: query.id
})
const show = (data) => {
  formState.group_name = data.group_name || ''
  formState.id = data.id || ''
  isEdit.value = !!data.id
  formState.library_id = data.library_id
  open.value = true
}
const handleOk = () => {
  if (!formState.group_name) {
    return message.error(t('msg_group_name_required'))
  }
  saveLibraryGroup({
    ...formState,
    group_type: props.group_type
  }).then((res) => {
    message.success(isEdit.value ? t('msg_edit_group_success') : t('msg_create_group_success'))
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

  .form-label{
    padding-right: 4px;
  }
  .form-content {
    flex: 1;
  }
}
</style>
