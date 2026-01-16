<template>
  <div class="field-tree-wrapper" @wheel.stop @touchmove.stop>
    <a-tree
      v-model:selectedKeys="selectedKeys"
      blockNode
      :defaultExpandAll="true"
      :height="224"
      :tree-data="props.treeData"
      :fieldNames="{children: 'subs', key: 'cu_key'}"
    >
      <template #switcherIcon="{ switcherCls }">
        <span :class="switcherCls">
          <svg-icon
            class="switcher-btn"
            name="arrow-down"
            style="font-size: 14px; color: #333"
          ></svg-icon>
        </span>
      </template>

      <template #title="{ dataRef }">
        <div class="field-item">
          <span class="field-key">{{ dataRef.key }}</span>
          <a-tooltip v-if="dataRef.tip" :overlayStyle="{ maxWidth: '800px' }">
            <template #title>
              <div class="tip-content">{{ dataRef.tip }}</div>
            </template>
            <QuestionCircleOutlined style="margin-left: 4px;" />
          </a-tooltip>
          <span class="field-type">{{ dataRef.typ }}</span>
          <span class="field-name">{{ dataRef.name }}</span>
        </div>
      </template>
    </a-tree>
  </div>
</template>

<script setup>
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
import { ref } from 'vue'
// const expandedKeys = ref(['0-0-0'])

const props = defineProps({
  treeData: {
    type: Array,
    default: () => []
  }
})

const selectedKeys = ref([])
</script>

<style lang="less" scoped>
.field-tree-wrapper {
  position: relative;
  height: 224px;
  :deep(.ant-tree) {
    background: none;
  }
  .switcher-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    font-size: 16px;
    border-radius: 6px;
    color: #000;

    &:hover {
      cursor: pointer;
      background-color: #d9d9d9;
    }
  }

  .field-item{
    display: flex;
    align-items: center;
    .field-key{
      color: #262626;
      text-align: right;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }
    .field-name{
      color: #595959;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
    }
    .field-type{
      display: flex;
      padding: 1px 8px;
      align-items: center;
      gap: 4px;
      border-radius: 6px;
      border: 1px solid #00000026;
      background: #FFF;
      color: #595959;
      font-size: 12px;
      font-style: normal;
      font-weight: 400;
      line-height: 20px;
      margin: 0 12px 0 8px;
    }
  }
}

.tip-content {
  white-space: pre-wrap;
}
</style>
