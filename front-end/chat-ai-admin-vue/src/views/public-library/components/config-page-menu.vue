<template>
  <div class="page-menus-box">
    <a-menu :selectedKeys="current" mode="horizontal" :items="menus" @click="handleMenu" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.public-library.components.config-page-menu')

const router = useRouter()
const route = useRoute()

const current = computed(() => [route.name])

// const current = ref(['permissions'])
const menus = ref([
  {
    key: 'PublicLibraryConfig',
    label: t('knowledge_base_settings'),
    title: t('knowledge_base_settings'),
    path: '/public-library/config'
  },
  {
    key: 'PublicLibraryPermissions',
    label: t('access_permissions'),
    title: t('access_permissions'),
    path: '/public-library/permissions'
  },
  {
    key: 'PublicLibraryAi',
    label: t('document_ai'),
    title: t('document_ai'),
    path: '/public-library/ai'
  },
  {
    key: 'PublicLibraryWebStatistics',
    label: t('statistics_settings'),
    title: t('statistics_settings'),
    path: '/public-library/web-statistics'
  }
])

const handleMenu = ({ item }) => {
  router.push({
    path: item.path,
    query: {
      library_id: route.query.library_id,
      library_key: route.query.library_key
    }
  })
}
</script>

<style lang="less" scoped>
.page-menus-box {
  margin-bottom: 16px;
  padding-left: 24px;
}
</style>
