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
    margin-bottom: 12px;
    gap: 12px;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }
}
</style>

<template>
  <a-modal
    class="add-model-alert"
    width="800px"
    v-model:open="show"
    :title="t('model_list')"
    :footer="null"
  >
    <div class="form-wrapper">
      <div class="tools-wrapper">
        <img class="model-logo" :src="modelConfig.model_icon_url" alt="" />
        <span>{{ modelConfig.model_name }}</span>
      </div>
      <a-table :data-source="use_model_configs" :pagination="false" :scroll="{ y: 500 }">
        <a-table-column key="show_model_name" :title="t('model_name')" data-index="show_model_name" />
        <a-table-column key="model_type" :title="t('model_type')" data-index="model_type"> </a-table-column>
      </a-table>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.model.components.see-model-alert')

const modelConfig = ref({})

const show = ref(false)
const use_model_configs = ref([])
const open = (data) => {
  modelConfig.value = data
  use_model_configs.value = data.use_model_configs || []

  show.value = true
}

defineExpose({
  open
})
</script>
