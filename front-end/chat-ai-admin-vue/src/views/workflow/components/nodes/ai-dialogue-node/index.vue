<style lang="less" scoped>
.ai-dialogue-node {
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
    <div class="ai-dialogue-node">
      <div class="node-desc">调用大模型，生成回复。</div>
      <FormBlock @setData="handleSetData" :properties="properties" />
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'
import FormBlock from './form-block.vue'
export default {
  name: 'AiDialogueNode',
  components: {
    NodeCommon,
    FormBlock
  },
  inject: ['getNode', 'getGraph', 'setData'],
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
  methods: {
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()

        this.getGraph().deleteNode(node.id)
      }
    },
    handleSetData(data) {
      this.setData({
        ...data
      })
    }
  },
  mounted() {}
}
</script>
