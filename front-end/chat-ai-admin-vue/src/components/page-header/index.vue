<template>
  <div class="page-header">
    <div class="page-breadcrumb" v-if="routes.length > 0">
      <a-breadcrumb :routes="routes" separator=">">
        <template #itemRender="{ route }">
          <span v-if="routes.indexOf(route) === routes.length - 1">
            {{ route.title }}
          </span>
          <a @click="onClick(route)" v-else>
            {{ route.title }}
          </a>
        </template>
      </a-breadcrumb>
    </div>
    <div class="page-name" v-else>{{ props.title }}</div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const props = defineProps({
  basePath: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  routes: {
    type: Array,
    default: () => []
  }
})

const onClick = (route) => {
  router.replace({
    path: route.path,
    query: route.query || {}
  })
}
</script>

<style lang="less" scoped>
.page-header {
  padding: 11px 24px;
  border-radius: 2px;
  background-color: #fff;
  .page-name {
    line-height: 24px;
    font-size: 16px;
    font-size: 16px;
    font-weight: 600;
    color: #000000;
  }
  .page-breadcrumb {
    padding: 1px 0;
  }
}
</style>
