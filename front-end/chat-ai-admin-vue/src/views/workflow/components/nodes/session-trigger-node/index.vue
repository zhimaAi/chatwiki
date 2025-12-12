<style lang="less" scoped>
.start-node {
  position: relative;
  padding-top: 8px;
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
    :title="properties.node_name"
    :icon-url="properties.node_icon"
    :icon-name="properties.node_icon_name"
    :isSelected="isSelected"
    :isHovered="isHovered"
    :node-key="properties.node_key"
    :node_type="properties.node_type"
    style="width: 420px;"
  >
    <div class="start-node">
      <div class="start-node-options">
        <div class="options-title">输出</div>
        <div class="options-list">
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
    </div>
  </node-common>
</template>

<script>
import NodeCommon from '../base-node.vue'

export default {
  name: 'StartNode',
  components: {
    NodeCommon
  },
  inject: ['resetSize'],
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
      outputs: [],
      show: false,
    }
  },
  computed: {
    options() {
      return [...this.outputs]
    }
  },
  watch: {
    properties: {
      handler(newVal, oldVal) {
        const newDataRaw =  newVal.node_params || '{}'
        const oldDataRaw = oldVal.node_params || '{}'
        
        if(newDataRaw != oldDataRaw) { 
          this.reset()
        }
      },
      deep: true,
    }
  },
  mounted() {
    this.reset()
  },
  methods: {
    reset() {
      if(!this.properties || !this.properties.node_params) {
        return
      }
      
      let node_params = JSON.parse(this.properties.node_params)

      this.outputs = node_params.trigger.outputs

      this.$nextTick(() => {
        this.resetSize()
      })
    },
  }
}
</script>
