<template>
  <node-common
    :properties="properties"
    :title="properties.node_name"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
  >
    <div class="specify-reply-node">
      <div class="static-field-list">
        <div class="static-field-item">
          <div class="static-field-item-label">{{ t('label_message_content') }}</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              <AtText
                :options="valueOptions"
                :defaultSelectedList="formState.content_tags"
                :defaultValue="formState.content"
                ref="atInputRef"
                class="text-input"
                @resize="onInputResize"
                v-if="formState.content.length > 0"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'
import AtText from '../../at-input/at-text.vue'
import { haveOutKeyNode } from '@/views/workflow/components/util.js'
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'ImmediatelyReplyNode',
  inject: ['getNode', 'getGraph', 'setData', 'resetSize'],
  components: {
    NodeCommon,
    AtText
  },
  setup() {
    const { t } = useI18n('views.workflow.components.nodes.immediately-reply-node.index')
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
      menus: [],
      valueOptions: [],
      formState: {
        content: '',
        content_tags: []
      }
    }
  },
  computed: {},
  watch: {
    properties: {
      handler(newVal, oldVal) {
        if (newVal.dataRaw !== oldVal.dataRaw) {
          this.init()
        }
      },
      deep: true
    }
  },
  mounted() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.on('custom:setNodeName', this.onUpatateNodeName)
    this.init()
  },
  onBeforeUnmount() {
    const graphModel = this.getGraph()
    graphModel.eventCenter.off('custom:setNodeName', this.onUpatateNodeName)
  },
  methods: {
    init() {
      this.getValueOptions()

      let dataRaw = this.properties.dataRaw || this.properties.node_params || '{}'
      let node_params = JSON.parse(dataRaw)

      const reply = node_params.immediately_reply || {
        content: '',
        content_tags: []
      }

      this.formState.content = reply.content

      this.formState.content_tags = reply.content_tags || []

      this.$nextTick(() => {
        this.resetSize()
      })
    },
    onInputResize() {
      this.$nextTick(() => {
        this.resetSize()
      })
    },
    onUpatateNodeName (data) {
      if (!this.formState.content) return
      if (!haveOutKeyNode.includes(data.node_type)) {
        return
      }

      this.getValueOptions()

      this.$nextTick(() => {
        if (this.formState.content_tags && this.formState.content_tags.length > 0) {
          this.formState.content_tags.forEach((tag) => {
            if (tag.node_id == data.node_id) {
              let arr = tag.label.split('/')
              arr[0] = data.node_name
              tag.label = arr.join('/')
              tag.node_name = data.node_name
            }
          })
        }

        this.$refs[`atInputRef`].refresh()

        // this.update()
      })
    },
    getValueOptions() {
      let options = this.getNode().getAllParentVariable()

      this.valueOptions = options || []
    },
    update() {
      let node_params = JSON.parse(this.properties.node_params)

      node_params.immediately_reply = { ...this.formState }

      this.setData({
        node_params: JSON.stringify(node_params)
      })
    },
  }
}
</script>

<style lang="less" scoped>
@import '../form-block.less';
.specify-reply-node {
  position: relative;
  z-index: 2;
  .static-field-value{
    width: 100%;
    height: 58px;
    line-height: 18px;
    font-size: 13px;
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    white-space: pre-wrap;
  }
  &:deep(.j-mention-at) {
    padding: 0 8px;
    border-radius: 4px;
    font-size: 12px;
    line-height: 16px;
    font-weight: 400;
    background: #f2f4f5;
  }
}
</style>
