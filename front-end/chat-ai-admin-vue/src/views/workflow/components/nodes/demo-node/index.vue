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
  <base-node
    :title="properties.node_name"
    :menus="menus"
    :icon-name="properties.node_icon_name"
    :node-key="properties.node_key"
    @handleMenu="handleMenu"
  >
    <template v-slot:icon>
      <img src="../../../../../assets/img/workflow/start-node-icon.svg" alt="" />
    </template>
    <div class="question-node">
      <div class="node-desc">{{ properties.content }}</div>
    </div>
  </base-node>
</template>

<script>
import BaseNode from '../base-node.vue'
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'QuestionNode',
  components: {
    BaseNode,
  },
  inject: ['getNode', 'getGraph'],
  setup() {
    const { t } = useI18n('views.workflow.components.nodes.demo-node.index')
    return { t }
  },
  props: {
    properties: {
      type: Object,
      default() {
        return {}
      },
    },
  },
  data() {
    return {
      menus: [{ name: this.t('btn_delete'), key: 'delete', color: '#fb363f' }],
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
