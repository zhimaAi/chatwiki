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
  <div class="custom-line-popup" v-show="isDeleteConfirmShow">
    <button class="delete-btn" @click.stop="handleDelete">{{ t('btn_delete') }}</button>
  </div>
</template>

<script>
import { useI18n } from '@/hooks/web/useI18n'

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
  data() {
    return {
      isDeleteConfirmShow: false,
    }
  },
  computed: {
    t() {
      return useI18n('views.workflow.components.edges.custom-line.index').t
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
    },
  },
}
</script>
