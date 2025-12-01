<style lang="less" scoped>
.start-node {
  position: relative;
  padding-top: 8px;
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }

  .start-node-options {
    display: flex;
    gap: 4px;
    .options-title {
      line-height: 22px;
      margin-right: 8px;
      font-size: 14px;
      color: #262626;
    }
    .options-list {
      flex: 1;
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
    .options-item {
      display: flex;
      align-items: center;
      height: 22px;
      padding: 2px 2px 2px 4px;
      border-radius: 4px;
      border: 1px solid #d9d9d9;

      &.is-required .option-label::before {
        vertical-align: middle;
        content: '*';
        color: #fb363f;
        margin-right: 2px;
      }

      .option-label {
        color: var(--wf-color-text-3);
        font-size: 12px;
        margin-right: 4px;
      }

      .option-type {
        height: 18px;
        line-height: 18px;
        padding: 0 8px;
        border-radius: 4px;
        font-size: 12px;
        background-color: #e4e6eb;
        color: var(--wf-color-text-3);
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
    :hide-menus="true"
    style="width: 200px"
  >
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'

export default {
  name: 'groupStartNode',
  components: {
    NodeCommon
  },
  inject: ['resetSize'],
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
  watch: {
    properties: {
      handler() {
      },
      deep: true
    }
  },
  mounted() {
  },
  methods: {}
}
</script>
