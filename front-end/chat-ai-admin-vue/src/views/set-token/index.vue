<template>
  <div></div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'
const userStore = useUserStore()
const router = useRouter()
const query = useRoute().query
userStore.setToken({
  token: query.token,
  exp: query.exp,
  ttl: query.ttl
})
userStore.setUserInfo({
  ...query
})
let redirect_url = query.redirect_url || '/'
router.push(redirect_url)
</script>

<style lang="less" scoped></style>
