<style lang="less" scoped>
.qa-node {
  .node-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: var(--wf-color-text-2);
  }

  .q-title {
    min-height: 38px;
    line-height: 22px;
    padding: 8px 12px;
    margin-top: 12px;
    border-radius: 4px;
    font-size: 14px;
    color: #595959;
    background-color: #f2f4f7;
  }

  .no-data {
    height: 38px;
    line-height: 22px;
    font-size: 14px;
    color: var(--wf-color-text-2);
  }

  .switch-status-box{
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 22px;
    margin-top: 12px;

    .switch-status-name{
      font-size: 14px;
      font-weight: 400;
      color: #262626;
    }
    .switch-status-value{
      display: flex;
      align-items: center;
      height: 22px;
      padding: 0 6px;
      border-radius: 6px;
      font-size: 13px;
      &.is-on{
        color: #2475fc;
        background: #E8EFFC;
      }
      &.is-off{
        color: #595959;
        background: #F5F5F5;
      }
      .status-icon{
        margin-right: 2px;
      }
    }
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
    <div class="qa-node">
      <div class="no-data" v-if="!properties.no_resp_jump_time">点击设置消息内容</div>
    </div>
  </node-common>
</template>

<script>
import { minutesToSeconds, secondsToMinutes} from '@/utils/index'
import NodeCommon from '../base-node.vue'

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
  data() {
    return {
      menus: [],
    }
  },
  methods: {
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()

        this.getGraph().deleteNode(node.id)
      }
    },
    secondsToMinutes(val){
      return secondsToMinutes(val)
    }
  },
}
</script>
