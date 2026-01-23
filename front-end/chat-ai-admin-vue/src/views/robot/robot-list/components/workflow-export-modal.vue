<template>
  <a-modal
    v-model:open="visible"
    title="导出知识库和数据库确认"
    width="746px"
    @ok="ok"
  >
    <a-alert class="zm-alert-info mt24" type="info" message="导出的工作流导入使用后，会自动创建引用的知识库和数据库，请确认是否同步导出知识库内容和数据库内容"/>
    <div class="main">
      <div v-if="libraryList.length" class="item">
        <a-checkbox v-model:checked="libraryChecked">知识库内容</a-checkbox>
        <div class="list">
          <a-tag v-for="item in libraryList" :key="item.id">{{ item.library_name }}</a-tag>
        </div>
        <div v-show="libraryChecked" class="warning">导出的工作流时会导出引用的知识库，勾选后会同时导出知识库内容，请注意企业内容使用，保证数据安全</div>
      </div>
      <div v-if="dbList.length" class="item">
        <a-checkbox v-model:checked="dbChecked">数据库内容</a-checkbox>
        <div class="list">
          <a-tag v-for="item in dbList" :key="item.form_id">{{ item.form_name }}</a-tag>
        </div>
        <div v-show="dbChecked" class="warning">导出的工作流时会导出引用的数据库，勾选后会同时导出数据库内容，请注意企业内容使用，保证数据安全</div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {ref} from 'vue'
import {robotImportDataInfo} from "@/api/robot/index.js";

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
