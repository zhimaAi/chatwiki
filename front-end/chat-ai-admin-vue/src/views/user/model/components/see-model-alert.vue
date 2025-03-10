<style lang="less" scoped>
.add-model-alert {
  .form-wrapper {
    margin-top: 24px;
  }

  .model-logo {
    display: block;
    height: 26px;
  }

  .tools-wrapper {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;
  }
}
</style>

<template>
  <a-modal class="add-model-alert" width="800px" v-model:open="show" :title="t('views.user.model.modelConfigButton')"
    :footer="null" @ok="handleOk">
    <div class="form-wrapper">
      <div class="tools-wrapper">
        <img class="model-logo" :src="modelConfig.model_icon_url" alt="" />
        <div class="tools-box">
          <a-button type="primary" v-if="modelConfig.showAddBtn" @click="handleAdd">{{
            t('views.user.model.addModel')
          }}</a-button>
        </div>
      </div>

      <div class="list-box">
        <a-table :data-source="modelConfig.rows" :scroll="{ y: 400 }" :pagination="false">
          <template v-for="col in modelConfig.columns" :key="col.key">
            <a-table-column :title="t('common.action')" :data-index="col.key" width="140px" v-if="col.key == 'action'">
              <template #default="{ record }">
                <span>
                  <a @click="handleEdit(record)">{{ t('common.edit') }}</a>
                  <a-divider type="vertical" />
                  <a @click="handleRemove(record)">{{ t('common.remove') }}</a>
                </span>
              </template>
            </a-table-column>
            <a-table-column :title="t('views.user.model.' + (modelNames.indexOf(modelConfig.model_define) > -1 && col.key == 'deployment_name' ? 'model_name' : col.key))" :data-index="col.key" v-else></a-table-column>
          </template>
        </a-table>
      </div>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive, markRaw, toRaw } from 'vue'
import { getModelTableConfig } from '../model-config'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()

const emit = defineEmits(['new', 'edit', 'remove'])

const show = ref(false)

const modelNames = ['ollama', 'xinference', 'openaiAgent']

const modelConfig = reactive({
  model_icon_url: '',
  model_define: '',
  columns: [],
  rows: [],
  showAddBtn: false
})

let modelParamsTemp = markRaw({})

const open = (data) => {
  modelParamsTemp = data

  let tableConfig = getModelTableConfig(data)

  for (let key in tableConfig) {
    modelConfig[key] = tableConfig[key]
  }

  show.value = true
}

// 添加模型
const handleAdd = () => {
  show.value = false
  emit('new', { ...modelParamsTemp })
}

const handleEdit = (record) => {
  show.value = false
  emit('edit', { ...modelParamsTemp }, toRaw(record))
}

const handleRemove = (record) => {
  show.value = false
  emit('remove', toRaw(record))
}

const handleOk = () => { }

defineExpose({
  open
})
</script>
