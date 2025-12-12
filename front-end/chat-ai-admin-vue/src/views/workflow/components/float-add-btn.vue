<style lang="less" scoped>
.add-node-btn {
  position: absolute;
  left: 16px;
  bottom: 16px;
  width: 48px;
  height: 48px;
  cursor: pointer;
  z-index: 100;
}
.add-node-btn.disabled {
  cursor: not-allowed;
  opacity: 0.5;
}
.node-list-fix {
  position: absolute;
  left: 78px;
  bottom: 16px;
  z-index: 100;
}
</style>

<template>
  <div>
    <img
      class="add-node-btn"
      :class="{'disabled': isLockedByOther}"
      src="../../../assets//img/workflow/add-node-btn.svg"
      @click="showMenu"
    />
    <div class="node-list-fix" v-show="isShowMenu">
      <NodeListPopup @addNode="handleAddNode" type="float-btn" />
    </div>
  </div>
</template>

<script>
import NodeListPopup from './node-list-popup/index.vue'
export default {
  name: 'float-add-btn',
  components: { NodeListPopup },
  props: {
    lf: {
      type: Object,
      required: true,
    },
    isLockedByOther: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      isShowMenu: false,
    }
  },
  mounted() {
    document.addEventListener('click', this.documentClick)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.documentClick)
  },
  methods: {
    showMenu() {
      if (this.isLockedByOther) return
      this.isShowMenu = !this.isShowMenu
    },
    // 点击菜单以外的地方，隐藏菜单
    documentClick(e) {
      if (this.isShowMenu) {
        const menus = this.$el.querySelector('.node-list-fix')
        if (!e.composedPath().includes(menus) && !menus.contains(e.target) && e.target !== this.$el.querySelector('.add-node-btn')) {
          this.isShowMenu = false
        }
      }
    },
    handleAddNode(node) {
      this.$emit('addNode', node)
      this.isShowMenu = false
    },
  },
}
</script>
