<template>
  <div></div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/modules/user'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const query = route.query

const toStringValue = (value, defaultValue = '') => {
  if (Array.isArray(value)) {
    return value[0] ?? defaultValue
  }
  return value == null ? defaultValue : String(value)
}

const redirectToLogin = (redirectUrl = '/') => {
  userStore.setToken(undefined)
  userStore.setUserInfo(undefined)
  userStore.setRoleRouters([])
  router.replace(`/login?redirect=${encodeURIComponent(redirectUrl || '/')}`)
}

const initSetTokenPage = async () => {
  const token = toStringValue(query.token)
  const exp = toStringValue(query.exp)
  const ttl = toStringValue(query.ttl)
  const userId = toStringValue(query.user_id)
  const redirectUrl = toStringValue(query.redirect_url, '/')
  const shouldRefreshUserInfo = toStringValue(query.refresh_user_info) === '1'

  if (!token) {
    redirectToLogin(redirectUrl)
    return
  }

  userStore.setToken({
    token,
    exp,
    ttl
  })
  userStore.setUserInfo({
    ...query
  })

  try {
    if (shouldRefreshUserInfo && userId) {
      await userStore.fetchUserExtraInfo(userId)
    }

    router.replace(redirectUrl || '/')
  } catch (error) {
    console.error(error)
    redirectToLogin(redirectUrl)
  }
}

initSetTokenPage()
</script>

<style lang="less" scoped></style>
