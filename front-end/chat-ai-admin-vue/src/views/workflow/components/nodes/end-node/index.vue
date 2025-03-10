<style lang="less" scoped>
.end-node {
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }
}
</style>

<template>
  <node-common
    :title="properties.node_name"
    :menus="menus"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    @handleMenu="handleMenu"
  >
    <div class="end-node">
      <div class="node-desc">触发该节点结束调用,并返回流程执行结果</div>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'

export default {
  name: 'EndNode',
  inject: ['getNode', 'getGraph', 'setData'],
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
    return {
      menus: [{ name: '删除', key: 'delete', color: '#fb363f' }]
    }
  },
  computed: {},

  methods: {
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()
        this.getGraph().deleteNode(node.id)
      }
    }
  }
}
</script>
