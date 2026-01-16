<style lang="less" scoped>
.beginner-guide-tag {
  position: relative;
  display: flex;
  position: relative;
  line-height: 20px;
  padding: 5px 12px;
  margin-left: 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 400;
  color: #595959;
  transition: all .2s;
  cursor: pointer;

  &:hover {
    background: #e4e6eb;
  }

  &.active {
    color: #fff;
    background: #2475fc;
  }

  .nav-icon {
    margin-right: 4px;
    font-size: 16px;
  }

  .red-dot{
    position: absolute;
    top: 4px;
    left: 24px;
    display: inline-block;
    width: 8px;
    height: 8px;
    background: #ff4d4f;
    border-radius: 50%;
    margin-right: 4px;
  }
}
</style>

<template>
  <div class="beginner-guide-tag" v-if="role_type == 1 && total_process != 100" @click="handleToGuide">
    <span class="red-dot"></span>
    <svg-icon class="nav-icon" name="guide"></svg-icon>
    <span class="nav-name">新手指引</span>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useGuideStore } from '@/stores/modules/guide'
import { usePermissionStore } from '@/stores/modules/permission'

const guideStore = useGuideStore()
const permissionStore = usePermissionStore()

const role_type = computed(() => {
  return permissionStore.role_type
})

const total_process = computed(() => {
  return +guideStore.total_process
})

const router = useRouter()

const handleToGuide = () => {
  guideStore.getUseGuideProcess()
  router.push('/guide')
}

onMounted(() => {
  guideStore.getUseGuideProcess()
})
</script>
