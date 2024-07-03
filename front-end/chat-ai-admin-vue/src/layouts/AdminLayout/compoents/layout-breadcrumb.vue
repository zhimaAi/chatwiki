<template>
  <div class="layout-breadcrumb">
    <a-breadcrumb v-if="props.items.length > 0">
      <template v-for="(item, index) in props.items" :key="item.title">
        <a-breadcrumb-item v-if="index !== props.items.length - 1"
          ><router-link :to="item.path">{{ item.title }}</router-link></a-breadcrumb-item
        >
        <a-breadcrumb-item v-else>{{ item.title }}</a-breadcrumb-item>
      </template>
    </a-breadcrumb>

    <div class="page-title" v-else>{{ props.title }}</div>
  </div>
</template>

<script setup>
const props = defineProps({
  title: {
    type: [String, Boolean],
    default: ''
  },
  items: {
    type: [Array, Boolean],
    default: () => []
  }
})
</script>

<style lang="less" scoped>
.layout-breadcrumb {
  display: flex;
  align-items: center;
  height: 56px;
  padding: 0 22px;

  .page-title {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
    color: #000000;
  }
}
</style>
