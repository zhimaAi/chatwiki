<template>
  <div class="center-display">
    <img class="brand-logo" :src="avatar" alt="" srcset="" v-if="avatar">
    <h1 class="brand-title" v-if="title">{{ title }}</h1>
    <p class="brand-desc" v-if="description">{{ description }}</p>
  </div>
</template>

<script lang="ts" setup>
interface Props {
  title?: string
  description?: string
  avatar?: string
}

withDefaults(defineProps<Props>(), {
  title: '',
  description: '',
  avatar: ''
})
</script>

<style lang="less" scoped>
.center-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
}
.brand-logo{
  width: 80px;
  height: 80px;
}
.brand-title {
  line-height: 44px;
  font-size: 36px;
  font-weight: 600;
  color: #0a0a0a;
  margin: 32px 0 24px 0;
}

.brand-desc {
  line-height: 24px;
  max-width: 800px;
  margin: 0;
  font-size: 16px;
  color: #7a8699;
  text-align: center;
}
</style>
