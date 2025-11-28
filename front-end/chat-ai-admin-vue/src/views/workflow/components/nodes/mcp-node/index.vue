<style lang="less" scoped>

</style>

<template>
  <node-common
    :title="toolInfo?.name"
    :icon-url="mcpInfo.avatar"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
  >
    <div class="mcp-node">
      MCP插件：
      <a-tag v-if="toolInfo === null" color="red">工具不存在或工具异常</a-tag>
      <a-tag v-else>{{ properties.node_name }}</a-tag>
    </div>
  </node-common>
</template>

<script>
import {jsonDecode} from '@/utils/index'
import NodeCommon from '../base-node.vue'
import {getTMcpProviderInfo} from "@/api/robot/thirdMcp.js";

export default {
  name: 'QuestionNode',
  components: {
    NodeCommon,
  },
  inject: ['getNode', 'getGraph', 'setData'],
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
      menus: [{name: '删除', key: 'delete', color: '#fb363f'}],
      visible: false,
      toolInfo: {},
      mcpInfo: {},
      nodeParams: {
        mcp: {
          provider_id: '',
          tool_name: '',
          arguments: {}
        }
      },
    }
  },
  mounted() {
    this.nodeParams = JSON.parse(this.properties.node_params)
    this.loadProvider()
  },
  methods: {
    loadProvider() {
      getTMcpProviderInfo({
        provider_id: this.nodeParams?.mcp.provider_id || ''
      }).then(res => {
        this.mcpInfo = res?.data || {}
        let tools = jsonDecode(this.mcpInfo?.tools, [])
        this.toolInfo = tools.find(item => item.name == this.nodeParams?.mcp.tool_name)
      })
    },
    update(data) {
      let nodeParams = this.nodeParams
      nodeParams.mcp.arguments = data.arguments
      nodeParams.mcp.tag_map = data.tag_map
      this.setData({
        node_params: JSON.stringify(nodeParams)
      })
    }
  },
}
</script>
