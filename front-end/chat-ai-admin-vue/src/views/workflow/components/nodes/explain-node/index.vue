<template>
  <node-common
    :title="properties.node_name"
    :menus="menus"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    @handleMenu="handleMenu"
  >
    <div class="explain-node" @mousedown.stop="" ref="textareaRef">
      <a-textarea
        :style="{ height: `${formState.height}px` }"
        @wheel.stop=""
        :bordered="false"
        class="explain-textarea"
        v-model:value="formState.content"
        placeholder="请输入注释"
        @input="updateContent($event.target.value)"
      />
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'

export default {
  name: 'EndNode',
  inject: ['getNode', 'getGraph', 'setData'],
  components: {
    NodeCommon
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
      observer: null,
      menus: [{ name: '删除', key: 'delete', color: '#fb363f' }],
      formState: {
        content: '',
        height: 88
      }
    }
  },
  computed: {},
  mounted() {
    let node_params = JSON.parse(this.properties.node_params)
    this.formState = {
      ...this.formState,
      ...node_params
    }

    const textareaEl = this.$refs.textareaRef?.querySelector('textarea')
    if (textareaEl) {
      // 创建观察器
      this.observer = new ResizeObserver((entries) => {
        for (const entry of entries) {
          const { height } = entry.contentRect
          this.formState.height = height + 8
          this.updateContent()
          // 在这里处理高度变化逻辑
        }
      })

      // 开始观察
      this.observer.observe(textareaEl)
    }
  },

  beforeDestroy() {
    // 组件销毁时断开观察
    if (this.observer) {
      this.observer.disconnect()
    }
  },

  methods: {
    getHeight() {
      return 68 + this.formState.height
    },
    updateContent() {
      this.setData({
        node_params: JSON.stringify(this.formState),
        height: this.getHeight()
      })
    },
    handleMenu(item) {
      if (item.key === 'delete') {
        let node = this.getNode()
        this.getGraph().deleteNode(node.id)
      }
    }
  }
}
</script>

<style lang="less" scoped>
.explain-node {
  .explain-textarea {
    color: #242933;
    font-size: 14px;
    background: #ffefd6;
    // min-height: 88px;
    border-color: #ffefd6;
    padding-left: 4px;
    &:hover {
      border-color: #2475fc;
    }

    /* 整个滚动条 */
    &::-webkit-scrollbar {
      width: 10px; /* 垂直滚动条宽度 */
      height: 10px; /* 水平滚动条高度 */
    }

    /* 滚动条轨道 */
    &::-webkit-scrollbar-track {
      background: #f1f1f1;
      border-radius: 5px;
    }

    /* 滚动条滑块 */
    &::-webkit-scrollbar-thumb {
      background: #888;
      border-radius: 5px;
      border: 2px solid transparent;
      background-clip: content-box;
    }

    /* 滚动条滑块悬停状态 */
    &::-webkit-scrollbar-thumb:hover {
      background: #666;
    }
  }
}
</style>
