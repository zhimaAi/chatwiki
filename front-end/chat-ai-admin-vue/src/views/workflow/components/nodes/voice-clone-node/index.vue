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
          <div class="static-field-item-label">{{ t('label_clone_audio') }}</div>
          <div class="static-field-item-content">
            <div class="static-field-value">
              <at-text
                v-if="formState.file_url.length > 0"
                :options="valueOptions"
                :defaultSelectedList="formState.file_url_tags"
                :defaultValue="formState.file_url"
                ref="atInputRef"
                class="text"
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
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'SpecifyReplyNode',
  inject: ['getNode', 'getGraph', 'setData', 'resetSize'],
  components: {
    NodeCommon,
    AtText
  },
  setup() {
    const { t } = useI18n('views.workflow.components.nodes.voice-clone-node.index')
    return { t }
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
        file_url: '',
        file_url_tags: []
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
      let args = node_params?.voice_clone?.arguments || {}
      this.formState.file_url = args.file_url || '--'
      this.formState.file_url_tags = args?.tag_map?.file_url || []
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
  .static-field-item-content{
    flex: 1;
    overflow: hidden;
  }
  .static-field-value {
    height: 22px;
    line-height: 20px;
    font-size: 12px;
    overflow: hidden;
    .text{
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
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
