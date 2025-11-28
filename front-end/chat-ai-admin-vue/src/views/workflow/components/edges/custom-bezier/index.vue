<style lang="less" scoped>
.custom-line-popup {
  /* 图标容器：定义圆形轮廓 */
  .minus-circle-outlined {
    width: 24px;
    height: 24px;
    border: 1px solid #d32f2f;
    border-radius: 50%;
    background-color: #ffebee;
    // box-shadow: 0 0 0 2px #ffffff;
    position: relative;
    cursor: pointer;
    &::before {
      content: '';
      width: 50%;
      height: 2px;
      background-color: #d32f2f;
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
    }
  }
}
</style>

<template>
  <div class="custom-line-popup" v-show="isDeleteConfirmShow">
    <div class="minus-circle-outlined" @click.stop="handleDelete"></div>
  </div>
</template>

<script>
export default {
  name: 'CustomLine',
  emits: ['delete'],
  inject: ['getGraphModel', 'getModel'],
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
      isDeleteConfirmShow: false,
    }
  },
  watch: {
    isSelected(newVal) {
      if (newVal) {
        let graphModel = this.getGraphModel()
        let selectedElements = graphModel.getSelectElements(true);

        if(selectedElements.nodes.length === 0 && selectedElements.edges.length === 1) {
          this.isDeleteConfirmShow = true;
        }
      } else {
        this.isDeleteConfirmShow = false;
      }
    }
  },
  methods: {
    handleDelete() {
      const graphModel = this.getGraphModel()
      const model = this.getModel()

      graphModel.eventCenter.emit('custom:edge:delete', model)
      
      this.$emit('delete')
    }
  }
}
</script>