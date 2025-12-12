<style lang="less" scoped></style>

<template>
  <PopupMenu @addNode="handleAddNode" :type="type" v-model:active="nodeListTabActive" />
</template>

<script>
import PopupMenu from '../../node-list-popup/index.vue'

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
      type: 'node',
      nodeListTabActive: 1,
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
      console.log(item)
      this.$emit('click-item', item)
    },
  },
}
</script>
