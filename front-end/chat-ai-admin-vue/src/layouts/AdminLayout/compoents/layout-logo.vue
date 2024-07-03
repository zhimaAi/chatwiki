<template>
  <div class="admin-logo-wrapper">
    <div>
      <img class="admin-logo" :src="imgUrl" alt="logo" />
    </div>
    <div v-if="name" class="system-name">{{ name }}</div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useCompanyStore } from '@/stores/modules/company'
import { DEFAULT_EN_LOGO } from '@/constants/index'

// 显示英文logo
const imgUrl = ref(DEFAULT_EN_LOGO)

const companyStore = useCompanyStore()
const name = computed(() => {
  return companyStore.name
})
</script>

<style lang="less" scoped>
.admin-logo-wrapper {
  display: flex;
  align-items: center;
}
.admin-logo {
  display: block;
  line-height: 1;
  height: 28px;
}
.system-name{
  margin-top: -2px;
  margin-left: 8px;
  font-size: 15px;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 230px;
}
</style>
