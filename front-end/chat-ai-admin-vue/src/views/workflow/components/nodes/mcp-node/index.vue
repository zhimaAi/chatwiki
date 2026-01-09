<style lang="less" scoped>
._main {
  display: flex;
  flex-direction: column;
  gap: 8px;
  > div {
    display: flex;
    > .label {
      flex-shrink: 0;
      min-width: 70px;
      text-align: right;
    }
    :deep(.ant-tag) {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }
  }
}
</style>

<template>
  <node-common
    :properties="properties"
    :title="properties.node_name"
    :icon-url="mcpInfo.avatar"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    style="min-width: 320px;min-height: 94px;"
  >
    <div class="_main">
      <div v-for="p in propList" :key="p.key">
        <span class="label">{{ p.key }}：</span>
        <a-tag v-if="getArgValue(p.key)">
          <at-text :options="valueOptions" :defaultSelectedList="getArgTags(p.key)" :defaultValue="getArgValue(p.key)" ref="atTextRef" />
        </a-tag>
        <a-tag v-else>--</a-tag>
      </div>
    </div>
  </node-common>
</template>

<script>
import {jsonDecode} from '@/utils/index'
import NodeCommon from '../base-node.vue'
import {getTMcpProviderInfo} from "@/api/robot/thirdMcp.js";
import AtText from "@/views/workflow/components/at-input/at-text.vue";

export default {
  name: 'QuestionNode',
  components: {
    NodeCommon,
    AtText,
  },
  inject: ['getNode', 'getGraph', 'setData', 'resetSize'],
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
      atTextRef: null,
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
      valueOptions: []
    }
  },
  mounted() {
    this.nodeParams = JSON.parse(this.properties.node_params)
    this.valueOptions = this.getNode().getAllParentVariable() || []
    this.loadProvider()

    const graphModel = this.getGraph()
    graphModel.eventCenter.on('custom:setNodeName', this.onUpatateNodeName)

    this.$nextTick(() => {
      this.resetSize()
    })
  },
  onBeforeUnmount() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.off('custom:setNodeName', this.onUpatateNodeName)
  },
  watch: {
    properties: {
      handler(newVal, oldVal) {
        if (newVal.dataRaw !== oldVal.dataRaw) {
          this.nodeParams = JSON.parse(this.properties.node_params || '{}')
          this.$nextTick(() => {
            const atTextRef = this.$refs.atTextRef
            const refList = Array.isArray(atTextRef) ? atTextRef : atTextRef ? [atTextRef] : []
            refList.forEach((ref) => {
              ref && ref.refresh && ref.refresh()
            })
            this.resetSize()
          })
        }
      },
      deep: true
    }
  },
  methods: {
    onUpatateNodeName (data) {
      this.valueOptions = this.getNode().getAllParentVariable() || []
      const patchTags = (tags) => {
        if (!Array.isArray(tags) || tags.length === 0) {
          return
        }
        tags.forEach((tag) => {
          if (tag && tag.node_id == data.node_id) {
            if (tag.label) {
              let arr = tag.label.split('/')
              arr[0] = data.node_name
              tag.label = arr.join('/')
            }
            tag.node_name = data.node_name
          }
        })
      }
      const map = this.nodeParams?.mcp?.tag_map || {}
      Object.keys(map).forEach((k) => {
        patchTags(map[k])
      })
      this.$nextTick(() => {
        const atTextRef = this.$refs.atTextRef
        const refList = Array.isArray(atTextRef) ? atTextRef : atTextRef ? [atTextRef] : []
        refList.forEach((ref) => {
          ref && ref.refresh && ref.refresh()
        })
        this.setData({
          node_params: JSON.stringify(this.nodeParams)
        })
        this.resetSize()
      })
    },
    loadProvider() {
      getTMcpProviderInfo({
        provider_id: this.nodeParams?.mcp.provider_id || ''
      }).then(res => {
        this.mcpInfo = res?.data || {}
        let tools = jsonDecode(this.mcpInfo?.tools, [])
        this.toolInfo = tools.find(item => item.name == this.nodeParams?.mcp.tool_name)
        this.$nextTick(() => {
          this.resetSize()
        })
      })
    },
    update(data) {
      let nodeParams = this.nodeParams
      nodeParams.mcp.arguments = data.arguments
      nodeParams.mcp.tag_map = data.tag_map
      this.setData({
        node_params: JSON.stringify(nodeParams)
      })
    },
    getArgValue(key) {
      const v = this.nodeParams?.mcp?.arguments?.[key]
      return v == null ? '' : String(v)
    },
    getArgTags(key) {
      const tags = this.nodeParams?.mcp?.tag_map?.[key]
      return Array.isArray(tags) ? tags : []
    }
  },
  computed: {
    propList() {
      const schema = this.toolInfo?.inputSchema || {}
      const props = schema?.properties || {}
      const req = Array.isArray(schema?.required) ? schema.required : []
      return Object.keys(props).map(k => {
        const it = props[k] || {}
        return {
          key: k,
          required: req.includes(k)
        }
      })
    }
  }
}
</script>
