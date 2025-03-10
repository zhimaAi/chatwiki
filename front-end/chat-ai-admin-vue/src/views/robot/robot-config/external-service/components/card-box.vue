<style lang="less" scoped>
.card-box {
  border-radius: 6px;
  background-color: #f2f4f7;
  .card-box-header {
    display: flex;
    align-items: center;
    height: 56px;
    padding: 0 16px;
    .card-icon {
      margin-right: 8px;
    }
    .card-title {
      flex: 1;
      .title-text {
        line-height: 24px;
        font-size: 16px;
        font-weight: 600;
        color: #262626;
      }
    }
  }

  .card-box-body {
    padding: 0 16px 16px 16px;
  }
}
</style>

<template>
  <div class="card-box">
    <div class="card-box-header">
      <div class="card-title">
        <span class="card-icon" v-if="$slots['icon']">
          <slot name="icon"></slot>
        </span>
        <span class="title-text">{{ props.title }}</span>
        <span class="title-text"><slot name="title"></slot></span>
      </div>

      <div class="action">
        <slot name="action"></slot>
      </div>
    </div>

    <div class="card-box-body">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    default: ''
  }
})
</script>
