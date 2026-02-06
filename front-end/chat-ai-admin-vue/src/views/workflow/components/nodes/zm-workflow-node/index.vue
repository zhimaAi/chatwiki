<style lang="less" scoped>
.start-node {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 8px;

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
    :properties="properties"
    :title="properties.node_name"
    :menus="menus"
    :icon-url="robotAvatar"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    @handleMenu="handleMenu"
  >
    <div class="start-node">
      <div class="start-node-options">
        <div class="options-title">{{ t('label_input_params') }}</div>
        <div class="options-list" style="max-height: 82px;overflow: hidden;">
          <div
            class="options-item"
            :class="{ 'is-required': item.required }"
            v-for="item in options"
            :key="item.key"
          >
            <div class="option-label">{{ item.key }}</div>
            <div class="option-type">{{ item.typ }}</div>
          </div>
        </div>
      </div>
      <div class="start-node-options">
        <div class="options-title">{{ t('label_exception_handling') }}</div>
        <div class="options-list">
          <div class="options-item">
            <div class="option-label">{{ t('msg_error_handling') }}</div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script>
import {jsonDecode} from '@/utils/index'
import NodeCommon from '../base-node.vue'
import {getRobotStartNode} from "@/api/robot/index.js";
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'ZmWorkflowNode',
  components: {
    NodeCommon,
  },
  inject: ['getNode', 'getGraph', 'resetSize', 'setData'],
  setup() {
    const { t } = useI18n('views.workflow.components.nodes.zm-workflow-node.index')
    return { t }
  },
  props: {
    properties: {
      type: Object,
      default() {
        return {}
      },
    },
    isSelected: {type: Boolean, default: false},
    isHovered: {type: Boolean, default: false},
  },
  data() {
    return {
      menus: [],
      robotAvatar: '',
      nodeParams: {},
      options: []
    }
  },
  mounted() {
   this.init()
  },
  methods: {
    init() {
      this.nodeParams = jsonDecode(this.properties.node_params)
      this.robotAvatar = this.nodeParams?.workflow?.robot_info?.robot_avatar
      this.loadRobotNode()
    },
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()

        this.getGraph().deleteNode(node.id)
      }
    },
    loadRobotNode() {
      getRobotStartNode({
        robot_id: this.nodeParams?.workflow?.robot_id
      }).then(res => {
        let node_params = jsonDecode(res?.data?.node_params)
        this.options = node_params?.start?.diy_global || []
        this.update()
        this.$nextTick(() => {
          this.resetSize()
        })
      })
    },
    update() {
      let nodeParams = this.nodeParams
      let _data = nodeParams?.workflow?.params || []
      if (!Array.isArray(_data) || !_data.length) {
        // 首次添加时参数字段添加
        let _state = JSON.parse(JSON.stringify(this.options))
        for (let item of _state) {
          item.variable = String(item.variable || "")
        }
        nodeParams.workflow.params = _state
        this.setData({
          node_params: JSON.stringify(nodeParams)
        })
      }
    }
  },
}
</script>
