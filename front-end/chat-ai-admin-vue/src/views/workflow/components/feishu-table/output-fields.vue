<template>
  <div class="field-tree-wrapper" @wheel.stop @touchmove.stop>
    <a-tree
      v-model:selectedKeys="selectedKeys"
      blockNode
      :defaultExpandAll="true"
      :height="224"
      :tree-data="props.treeData"
      :fieldNames="{children: 'subs'}"
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
          <span class="field-name">{{ dataRef.key }}</span>
          <span class="field-type">{{ dataRef.typ }}</span>
        </div>
      </template>
    </a-tree>
  </div>
</template>

<script setup>
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
    .field-name{
      line-height: 22px;
      font-size: 14px;
      color: rgb(38, 38, 38);
    }
    .field-type{
      height: 22px;
      line-height: 20px;
      padding: 0 8px;
      margin-left: 8px;
      border-radius: 6px;
      background: #FFFFFF;
      border: 1px solid rgba(0, 0, 0, 0.15);
    }
  }
}
</style>
