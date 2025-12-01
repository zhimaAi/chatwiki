<style lang="less" scoped></style>

<template>
  <PopupMenu @addNode="handleAddNode" :type="type" />
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
    return {
      type: 'node'
    }
  },
  mounted() {
    const mode = this.getGraph()
    mode.eventCenter.on('custom:showPopupMenu', this.onShowPopupMenu)
  },
  methods: {
    onShowPopupMenu(data){
      if(data.model.properties.loop_parent_key){
        this.type = 'loop-node'
      }else{
        this.type = 'node'
      }
    },
    handleAddNode(item) {
      this.$emit('click-item', item)
    },
  },
}
</script>
