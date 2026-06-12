<template>
  <a-modal
    :open="open"
    :title="null"
    :footer="null"
    :width="width"
    :destroyOnClose="destroyOnClose"
    :maskClosable="maskClosable"
    :centered="centered"
    :closable="closable"
    :wrapClassName="`cu-modal-wrap ${wrapClassName}`"
    :zIndex="zIndex"
    @cancel="handleCancel"
  >
    <slot />
  </a-modal>
</template>

<script setup>
defineProps({
  open: {
    type: Boolean,
    default: false
  },
  width: {
    type: [Number, String],
    default: 520
  },
  destroyOnClose: {
    type: Boolean,
    default: true
  },
  maskClosable: {
    type: Boolean,
    default: true
  },
  centered: {
    type: Boolean,
    default: false
  },
  closable: {
    type: Boolean,
    default: true
  },
  wrapClassName: {
    type: String,
    default: ''
  },
  zIndex: {
    type: Number,
    default: 1000
  }
})

const emit = defineEmits(['update:open', 'cancel'])

const handleCancel = (e) => {
  emit('update:open', false)
  emit('cancel', e)
}
</script>

<style lang="less" scoped>
/* 组件自身 DOM 的 scoped 样式槽位（当前为空，留作日后扩展） */
</style>

<style lang="less">
/* a-modal 内容会被 Teleport 到 body 下，scoped 的 data-v 属性带不过去，
   因此这里必须用全局样式。靠 wrapClassName="cu-modal-wrap" 限定作用域，
   该 class 只在本组件被显式添加，不会污染其他弹窗。 */
.cu-modal-wrap {
  .ant-modal-content {
    padding: 0 !important;
    border-radius: 0 !important;
    background: none !important;
    box-shadow: none !important;
  }
  .ant-modal-content {
    padding: 0 !important;
    border-radius: 0 !important;
    background: none !important;
  }
}
</style>
