<template>
  <a-modal
    v-model:open="visible"
    :title="t(isChatRobot ? 'title_chat_export_confirm' : 'title_export_confirm')"
    width="746px"
    @ok="ok"
  >
    <a-alert class="zm-alert-info mt24" type="info" :message="t(isChatRobot ? 'msg_chat_export_info' : 'msg_export_info')"/>
    <div class="main">
      <div v-if="libraryList.length" class="item">
        <a-checkbox v-model:checked="libraryChecked">{{ t('label_library_content') }}</a-checkbox>
        <div class="list">
          <a-tag v-for="item in libraryList" :key="item.id">{{ item.library_name }}</a-tag>
        </div>
        <div v-show="libraryChecked" class="warning">{{ t(isChatRobot ? 'msg_chat_library_warning' : 'msg_library_warning') }}</div>
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
import {computed, onBeforeUnmount, ref} from 'vue'
import {useI18n} from '@/hooks/web/useI18n'
import {robotImportDataInfo} from "@/api/robot/index.js";

const {t} = useI18n('views.robot.robot-list.components.workflow-export-modal')

const emit = defineEmits(['ok'])
const visible = ref(false)
const robotInfo = ref({})
const isChatRobot = computed(() => String(robotInfo.value.application_type) === '0')
const libraryChecked = ref(false)
const dbChecked = ref(false)
const dbList = ref([])
const libraryList = ref([])
let requestSequence = 0

function handle(r = {}) {
  const request = ++requestSequence
  visible.value = false
  reset()
  // 请求错误由全局拦截器提示。
  return robotImportDataInfo({id: r.id}).catch(() => null).then(res => {
    if (request !== requestSequence || !res?.data) return

    robotInfo.value = r
    dbList.value = res.data.databaseList || []
    libraryList.value = res.data.library || []
    if (!libraryList.value.length && !dbList.value.length) {
      emit('ok', robotInfo.value)
    } else {
      visible.value = true
    }
  })
}

onBeforeUnmount(() => { requestSequence++ })

function reset() {
  libraryList.value = []
  dbList.value = []
  libraryChecked.value = false
  dbChecked.value = false
}

function ok() {
  if (!visible.value) return

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
