<template>
  <EmptyBox>
    <template #image>
      <img src="@/assets/empty-network.png" width="200px"/>
    </template>
    <template #title>网络异常，无法连接服务器</template>
    <template #desc>
      <div class="empty-desc-box">
        <div>相关功能已自动关闭请检查网络后重试</div>
        <a-button class="btn" type="primary" :loading="loading" @click="reload">
          {{loading ? '重新连接中，请稍候...' : '重新连接'}}
        </a-button>
      </div>
    </template>
  </EmptyBox>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from 'vue'
import EmptyBox from "@/components/common/empty-box.vue";
import {useCompanyStore} from "@/stores/modules/company.js";

const companyStore = useCompanyStore()
const reloadCount = ref(1)
const timer = ref(null)
const loading = ref(false)

onMounted(() => {
  autoReload()
})

onUnmounted(() => {
  timer.value && clearInterval(timer.value)
})

function autoReload() {
  reload()
  timer.value = setInterval(() => {
    if (reloadCount.value > 2) {
      clearInterval(timer.value)
    } else {
      reloadCount.value += 1
      reload()
    }
  }, 4000)
}

function reload() {
  if (loading.value) return
  loading.value = true
  companyStore.getCompanyInfo().finally(() => {
    setInterval(() => {
      loading.value = false
    }, 1500)
  })
}
</script>

<style scoped lang="less">
.empty-desc-box {
  text-align: center;
  .btn {
    min-width: 148px;
    margin-top: 16px;
  }
}
</style>
