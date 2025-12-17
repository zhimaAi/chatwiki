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
    title="模型列表"
    :footer="null"
  >
    <div class="form-wrapper">
      <div class="tools-wrapper">
        <img class="model-logo" :src="modelConfig.model_icon_url" alt="" />
        <span>{{ modelConfig.model_name }}</span>
      </div>
      <a-table :data-source="use_model_configs" :pagination="false" :scroll="{ y: 500 }">
        <a-table-column key="show_model_name" title="模型名称" data-index="show_model_name" />
        <a-table-column key="model_type" title="模型类型" data-index="model_type"> </a-table-column>
      </a-table>
    </div>
  </a-modal>
</template>
<script setup>
import { ref, reactive } from 'vue'

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
