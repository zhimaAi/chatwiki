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
      <div class="node-desc">{{ t('msg_trigger_node_end_return_result') }}</div>
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
  setup() {
    const { t } = useI18n('views.workflow.components.nodes.end-node.index')
    return { t }
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
  mounted() {
    this.$nextTick(() => {
      this.resetSize()
    })
  },
  computed: {},
  methods: {
    handleMenu() {

    }
  }
}
</script>
