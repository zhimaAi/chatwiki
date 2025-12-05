<template>
  <div class="options-item is-required">
    <div class="options-item-tit">
      <div class="option-label">record_id</div>
      <div class="option-type">string</div>
    </div>
    <div>
      <AtInput
        type="textarea"
        inputStyle="height: 64px;"
        :options="variableOptions"
        :defaultSelectedList="state.tags"
        :defaultValue="state.value"
        ref="atInputRef"
        @open="emit('updateVar')"
        @change="changeValue"
        placeholder="请输入内容，键入“/”可以插入变量"
      >
        <template #option="{ label, payload }">
          <div class="field-list-item">
            <div class="field-label">{{ label }}</div>
            <div class="field-type">{{ payload.typ }}</div>
          </div>
        </template>
      </AtInput>
    </div>
    <div class="desc">待删除的记录ID；示例值："reCWNXZPQv"</div>
  </div>
</template>

<script setup>
import {ref, reactive, toRaw} from 'vue'
import AtInput from "@/views/workflow/components/at-input/at-input.vue";

const props = defineProps({
  variableOptions: {
    type: Array,
  }
})
const emit = defineEmits(['update', 'updateVar'])
const state = reactive({
  tags: [],
  value: '',
})

function init(nodeParams=null) {
  if (nodeParams && nodeParams?.plugin?.params?.arguments?.record_id) {
    state.value = nodeParams.plugin.params.arguments.record_id
    state.tags = nodeParams.plugin.params.arguments.record_tags || []
  }
}

function changeValue(val, tags) {
  state.value = val
  state.tags = tags
  emit('update', {
    record_id: val,
    record_tags: tags,
  })
}

defineExpose({
  init,
})
</script>

<style scoped lang="less">
.options-item {
  display: flex;
  flex-direction: column;
  margin-top: 12px;
  line-height: 22px;
  gap: 4px;

  .options-item-tit {
    display: flex;
    align-items: center;
  }

  .option-label {
    color: var(--wf-color-text-1);
    font-size: 14px;
    margin-right: 8px;
  }

  .desc {
    color: var(--wf-color-text-2);
  }


  &.is-required .option-label::before {
    content: '*';
    color: #FB363F;
    display: inline-block;
    margin-right: 2px;
  }

  .option-type {
    height: 22px;
    line-height: 18px;
    padding: 0 8px;
    border-radius: 6px;
    border: 1px solid rgba(0, 0, 0, 0.15);
    background-color: #fff;
    color: var(--wf-color-text-3);
    font-size: 12px;
  }

  .item-actions-box {
    display: flex;
    align-items: center;

    .action-btn {
      margin-left: 12px;
      font-size: 16px;
      color: #595959;
      cursor: pointer;
    }
  }
}
</style>
