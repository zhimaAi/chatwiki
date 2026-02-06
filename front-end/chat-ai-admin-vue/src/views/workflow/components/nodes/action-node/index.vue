<style lang="less" scoped>
.question-node {
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }

  .q-title {
    line-height: 22px;
    padding: 8px 12px;
    margin-top: 12px;
    border-radius: 4px;
    font-size: 14px;
    color: #595959;
    background-color: #f2f4f7;
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
    @handleMenu="handleMenu"
  >
    <div class="question-node">
      <div class="node-desc">{{ properties.content || t('msg_click_to_set_content') }}</div>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'QuestionNode',
  components: {
    NodeCommon,
  },
  inject: ['getNode', 'getGraph'],
  props: {
    properties: {
      type: Object,
      default() {
        return {}
      },
    },
    isSelected: { type: Boolean, default: false },
    isHovered: { type: Boolean, default: false },
  },
  setup() {
    const { t } = useI18n('views.workflow.components.nodes.action-node.index')
    return { t }
  },
  data() {
    return {
      menus: [],
    }
  },
  mounted() {},
  methods: {
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()

        this.getGraph().deleteNode(node.id)
      }
    },
  },
}
</script>
