<style lang="less" scoped></style>

<template>
  <PopupMenu @addNode="handleAddNode" type="node" />
</template>

<script>
import PopupMenu from '../../node-list-popup.vue'

export default {
  name: 'QuestionNode',
  components: {
    PopupMenu,
  },
  emits: ['click-item'],
  inject: ['getNode', 'getGraph'],
  props: {
    properties: {
      type: Object,
      default() {
        return {}
      },
    },
  },
  data() {
    return {}
  },
  mounted() {},
  methods: {
    handleAddNode(item) {
      this.$emit('click-item', item)
    },
  },
}
</script>
