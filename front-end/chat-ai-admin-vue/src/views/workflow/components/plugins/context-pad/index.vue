<style lang="less" scoped></style>

<template>
  <PopupMenu @addNode="handleAddNode" :type="type" :excludedNodeTypes="excludedNodeTypes" v-model:active="nodeListTabActive" />
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
      excludedNodeTypes: [],
      nodeListTabActive: 1,
    }
  },
  mounted() {
    const mode = this.getGraph()
    mode.eventCenter.on('custom:showPopupMenu', this.onShowPopupMenu)
  },
  methods: {
    onShowPopupMenu(data){
      const {type, loop_parent_key} = data.model.properties

      if(loop_parent_key){
        let excludedNodeTypes = ['custom-group', 'batch-group', 'end-node'];
        let parentNode = this.getGraph().getNodeModelById(loop_parent_key)

        if(parentNode.type == 'batch-group'){
          excludedNodeTypes.push('terminate-node')
        }

        this.excludedNodeTypes = excludedNodeTypes;
      }else{
        this.excludedNodeTypes = ['explain-node', 'terminate-node']
      }
    },
    handleAddNode(item) {
      this.$emit('click-item', item)
    },
  },
}
</script>
