<template>
  <div class="_main">
    <div v-if="actionState?.use_plugin_config">
      <span class="label">配置：</span>
      <a-tag>{{ actionState?.config_name || '--' }}</a-tag>
    </div>
    <div v-for="(item, idx) in (actionState?.rendering || [])" :key="item.key || idx">
      <span class="label">{{ item.label }}：</span>
      <a-tag v-if="item.value">
        <at-text :options="valueOptions" :defaultSelectedList="item.tags" :defaultValue="item.value" ref="atTextRef" />
      </a-tag>
      <a-tag v-else>--</a-tag>
    </div>
  </div>
</template>

<script>
import AtText from "@/views/workflow/components/at-input/at-text.vue";
import { haveOutKeyNode } from '@/views/workflow/components/util.js'

export default {
  name: "feishu-table-render",
  inject: ['getGraph', 'setData'],
  components: {AtText},
  props: {
    nodeParams: {
      type: Object
    },
    valueOptions: {
      type: Array
    },
  },
  data() {
    return {
      atTextRef: null,
      actionState: {},
    }
  },
  mounted() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.on('custom:setNodeName', this.onUpatateNodeName)
    this.dataFormat()
  },
  onBeforeUnmount() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.off('custom:setNodeName', this.onUpatateNodeName)
  },
  watch: {
    nodeParams: {
      handler(newVal, oldVal) {
        this.dataFormat()
        this.$emit('updateSize')
      },
      deep: true
    }
  },
  methods: {
    onUpatateNodeName (data) {
      if (!haveOutKeyNode.includes(data.node_type)) {
        return
      }

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

      const rendering = this.actionState?.rendering || []
      rendering.forEach((item) => {
        patchTags(item?.tags)
      })

      this.$nextTick(() => {
        const atTextRef = this.$refs.atTextRef
        const refList = Array.isArray(atTextRef) ? atTextRef : atTextRef ? [atTextRef] : []
        refList.forEach((ref) => {
          ref && ref.refresh && ref.refresh()
        })
        this.update()
        this.dataFormat()
        this.$emit('updateSize')
      })
    },
    update() {
      this.setData({
        node_params: JSON.stringify(this.nodeParams)
      })
    },
    dataFormat() {
      let params = this.nodeParams?.plugin?.params || {}
      let arg = params?.arguments || {}
      let render = params?.rendering || {}

      if (Array.isArray(params?.rendering)) {
        const tagMap = params?.arguments?.tag_map || params?.tag_map || {}
        params.rendering.forEach((item) => {
          if (!item) return
          if (!Array.isArray(item.tags) || item.tags.length === 0) {
            const key = item.key
            const tags = key && tagMap ? tagMap[key] : null
            if (Array.isArray(tags) && tags.length) {
              item.tags = tags
            }
          }
        })
      }
      this.actionState = {
        ...params,
        ...arg,
        ...render
      }
      // console.log('this.actionState', this.actionState)
    },
    getFieldLabel(field) {
      return field?.properties?.[field?.current_properties_key]?.name || field?.name
    },
    getFieldValue(field) {
      if (field?.media_component) {
        const p = field?.properties?.[field?.current_properties_key] || {}
        return p?.value ?? ''
      }
      return field?.value ?? ''
    },
    getFieldTags(field) {
      if (field?.media_component) {
        const p = field?.properties?.[field?.current_properties_key] || {}
        return p?.atTags ?? []
      }
      return field?.atTags ?? []
    }
  }
}
</script>

<style scoped lang="less">
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

.hint {
  position: relative;
  cursor: pointer;
}
.hint-pop {
  position: absolute;
  bottom: 30px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 12px;
  opacity: 0;
  pointer-events: none;
  min-width: max-content;
  min-height: 32px;
  padding: 6px 8px;
  color: #fff;
  text-align: start;
  text-decoration: none;
  word-wrap: break-word;
  background-color: rgba(0, 0, 0, 0.85);
  border-radius: 6px;
  box-shadow: 0 6px 16px 0 rgba(0, 0, 0, 0.08), 0 3px 6px -4px rgba(0, 0, 0, 0.12), 0 9px 28px 8px rgba(0, 0, 0, 0.05);

  &::after {
    content: '';
    position: absolute;
    bottom: -6px;
    left: 50%;
    transform: translateX(-50%);
    width: 0;
    height: 0;
    border-left: 8px solid transparent;
    border-right: 8px solid transparent;
    border-top: 8px solid rgba(0, 0, 0, 0.85);
  }
}
.hint:hover .hint-pop {
  opacity: 1;
  pointer-events: auto;
}
</style>
