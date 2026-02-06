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
    :properties="properties"
    :title="properties.node_name"
    :menus="menus"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    @handleMenu="handleMenu"
    style="width: 420px; height: 116px;"
  >
    <div class="end-node">
      <div class="node-desc">{{ t('msg_terminate_desc') }}</div>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'EndNode',
  inject: ['getNode', 'getGraph', 'setData', 'resetSize'],
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
      menus: []
    }
  },
  computed: {
    t() {
      return useI18n('views.workflow.components.nodes.terminate-node.index').t
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.resetSize()
    })
  },
  methods: {
    handleMenu() {

    }
  }
}
</script>
