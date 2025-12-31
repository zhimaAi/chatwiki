<template>
  <node-common
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
          <div class="static-field-item-label">语音模型</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              {{ formState.use_model }}
            </div>
          </div>
        </div>
        <div class="static-field-item">
          <div class="static-field-item-label">文本</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              <at-text
                v-if="formState.text.length > 0"
                :options="valueOptions"
                :defaultSelectedList="formState.text_tags"
                :defaultValue="formState.text"
                ref="atInputRef"
              />
              <span v-else>--</span>
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
import {haveOutKeyNode} from '@/views/workflow/components/util.js'

export default {
  name: 'SpecifyReplyNode',
  inject: ['getNode', 'getGraph', 'setData', 'resetSize'],
  components: {
    NodeCommon,
    AtText
  },
  props: {
    properties: {
      type: Object,
      default() {
        return {}
      }
    },
    isSelected: {type: Boolean, default: false},
    isHovered: {type: Boolean, default: false}
  },
  data() {
    return {
      menus: [],
      valueOptions: [],
      formState: {
        use_model: '',
        text: '',
        text_tags: []
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
    this.init()
  },
  methods: {
    init() {
      this.getValueOptions()

      let dataRaw = this.properties.dataRaw || this.properties.node_params || '{}'
      let node_params = JSON.parse(dataRaw)
      let args = node_params?.text_to_audio?.arguments || {}
      this.formState.use_model = args.use_model || '--'
      this.formState.text = args.text || ''
      this.formState.text_tags = args?.tag_map?.text || []
      this.$nextTick(() => {
        this.resetSize()
      })
    },
    getValueOptions() {
      let options = this.getNode().getAllParentVariable()
      this.valueOptions = options || []
    },
  }
}
</script>

<style lang="less" scoped>
@import '../form-block.less';

.specify-reply-node {
  position: relative;
  z-index: 2;

  .static-field-value {
    height: auto;
    line-height: 16px;
    font-size: 12px;
    word-break: break-all;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
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
