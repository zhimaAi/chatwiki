<template>
  <a-modal
    v-model:open="visible"
    :title="t('title_export_confirm')"
    width="746px"
    @ok="ok"
  >
    <a-alert class="zm-alert-info mt24" type="info" :message="t('msg_export_info')"/>
    <div class="main">
      <div v-if="libraryList.length" class="item">
        <a-checkbox v-model:checked="libraryChecked">{{ t('label_library_content') }}</a-checkbox>
        <div class="list">
          <a-tag v-for="item in libraryList" :key="item.id">{{ item.library_name }}</a-tag>
        </div>
        <div v-show="libraryChecked" class="warning">{{ t('msg_library_warning') }}</div>
      </div>
      <div v-if="dbList.length" class="item">
        <a-checkbox v-model:checked="dbChecked">{{ t('label_database_content') }}</a-checkbox>
        <div class="list">
          <a-tag v-for="item in dbList" :key="item.form_id">{{ item.form_name }}</a-tag>
        </div>
        <div v-show="dbChecked" class="warning">{{ t('msg_database_warning') }}</div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {ref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {robotImportDataInfo} from "@/api/robot/index.js";

const {t} = useI18n('views.robot.robot-list.components.workflow-export-modal')

const emit = defineEmits(['ok'])
const visible = ref(false)
const robotInfo = ref({})
const libraryChecked = ref(false)
const dbChecked = ref(false)
const dbList = ref([])
const libraryList = ref([])

function handle(r = {}) {
  reset()
  robotInfo.value = r
  robotImportDataInfo({id: r?.id}).then(res => {
    let data = res?.data || {}
    dbList.value = data?.databaseList || []
    libraryList.value = data?.library || []
    if (!libraryList.value.length && !dbList.value.length) {
      emit('ok', robotInfo.value)
    } else {
      visible.value = true
    }
  })
}

function reset() {
  libraryList.value = []
  dbList.value = []
  libraryChecked.value = false
  dbChecked.value = false
}

function ok() {
  let data = {}
  if (libraryChecked.value && libraryList.value.length) {
    data.library_id = libraryList.value.map(i => i.id).toString()
  }
  if (dbChecked.value && dbList.value.length) {
    data.form_id = dbList.value.map(i => i.form_id).toString()
  }
  visible.value = false
  emit('ok', robotInfo.value, data)
}

defineExpose({
  handle,
})
</script>

<style scoped lang="less">
.mt24 {
  margin-top: 24px;
}

.main {
  .item {
    margin-top: 12px;
    padding: 16px;
    background: #F2F4F7;

    .list {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      gap: 8px;
      margin-top: 12px;

      :deep(.ant-tag) {
        color: #595959;
        font-size: 14px;
        font-weight: 400;
        padding: 5px 16px;
        border-radius: 6px;
        border: 1px solid #D9D9D9;
        background: #F5F5F5;
      }
    }

    .warning {
      color: #fb363f;
      margin-top: 8px;
    }
  }
}
</style>
