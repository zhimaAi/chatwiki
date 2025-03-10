<style lang="less" scoped>
.start-node {
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }

  .start-node-options {
    background: #f2f4f7;
    border-radius: 6px;
    padding: 12px;
    margin-top: 16px;
    .options-title {
      color: var(--wf-color-text-1);
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 600;
      line-height: 22px;
      font-size: 14px;
    }
    .options-item{
      margin-top: 12px;
      height: 22px;
      line-height: 22px;
      display: flex;
      align-items: center;
      gap: 8px;
      .option-label{
        color: var(--wf-color-text-1);
        font-size: 14px;
        &::before{
          content: '*';
          color: #FB363F;
          display: inline-block;
          margin-right: 2px;
        }
      }
      .option-type{
        height: 22px;
        width: fit-content;
        padding: 0 8px;
        border-radius: 6px;
        border: 1px solid rgba(0, 0, 0, 0.15);
        background-color: #fff;
        color: var(--wf-color-text-3);
        font-size: 12px;
        display: flex;
        align-items: center;
      }
    }
  }
}
</style>

<template>
  <node-common
    :title="properties.node_name"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
  >
    <div class="start-node">
      <div class="node-desc">流程开始节点，点击可以设置全局变量、触发方式等</div>

      <div class="start-node-options">
        <div class="options-title">输出</div>
        <div class="options-item">
          <div class="option-label">用户消息</div>
          <div class="option-type">string</div>
        </div>
        <div class="options-item">
          <div class="option-label">open_id</div>
          <div class="option-type">string</div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'

export default {
  name: 'StartNode',
  components: {
    NodeCommon
  },
  props: {
    properties: {
      type: Object,
      default() {
        return {}
      }
    },
    isSelected: { type: Boolean, default: false },
    isHovered: { type: Boolean, default: false }
  },
  data() {
    return {}
  },
  computed: {},
  created() {
  },
  methods: {}
}
</script>
