<style lang="less" scoped>
.custom-line-popup {
  border-radius: 6px;
  background: #fff;

  .delete-btn {
    width: 70px;
    height: 36px;
    line-height: 36px;
    padding: 0;
    border: 0;
    font-size: 14px;
    color: #fb363f;

    text-align: center;
    cursor: pointer;
    background: none;
  }
}
</style>

<template>
  <div class="custom-line-popup" v-show="isSelected">
    <button class="delete-btn" @click.stop="handleDelete">删 除</button>
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
      },
    },
    isSelected: { type: Boolean, default: false },
    isHovered: { type: Boolean, default: false },
  },
  methods: {
    handleDelete() {
      const graphModel = this.getGraphModel()
      const model = this.getModel()

      graphModel.deleteEdgeById(model.id)

      this.$emit('delete')
    },
  },
}
</script>
