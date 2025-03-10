<template>
  <div class="left-menu-box">
    <a-menu class="left-menu" :selectedKeys="selectedKeys" @click="handleChangeMenu">
      <router-link
        class="default-color"
        :to="{ path: item.path, query: query }"
        v-for="item in items"
        :key="item.key"
      >
        <a-menu-item :icon="item.icon" :path="item.path" :key="item.key">{{
          item.label
        }}</a-menu-item>
      </router-link>
    </a-menu>
  </div>
</template>

<script setup>
import { ref, h, computed } from 'vue'
import { useRoute } from 'vue-router'
import SvgIcon from '@/components/svg-icon/index.vue'

const emit = defineEmits(['changeMenu'])
const route = useRoute()
const query = route.query
const selectedKeys = computed(() => {
  return [route.path.split('/')[3]]
})

const items = ref([
  {
    key: 'field-manage',
    id: 'field-manage',
    icon: () =>
      h('span', {}, [
        h(SvgIcon, {
          name: 'field-menu-icon',
          class: 'menu-icon'
        }),
      ]),
    label: '字段管理',

    title: '字段管理',
    path: '/database/details/field-manage'
  },
  {
    key: 'database-manage',
    id: 'database-manage',
    icon: () =>
      h('span', {}, [
        h(SvgIcon, {
          name: 'database-menu-icon',
          class: 'menu-icon'
        }),
      ]),
    label: '数据管理',
    title: '数据管理',
    path: '/database/details/database-manage'
  }
])

const handleChangeMenu = ({ item }) => {
  if (selectedKeys.value.includes(item.id)) {
    return
  }

  emit('changeMenu', item)
}
</script>

<style lang="less" scoped>
.left-menu-box {
  .default-color {
    color: inherit;
  }
  .left-menu {
    border-right: 0;

    ::v-deep(.menu-icon) {
      color: #a1a7b3;
      font-size: 16px;
      vertical-align: -3px;
    }
    ::v-deep(.ant-menu-item-selected .menu-icon) {
      color: #2475fc;
    }
    ::v-deep(.ant-menu-item-icon +span){
      margin-left: 4px;
    }
    ::v-deep(.ant-menu-item-icon){
      vertical-align: -3px;
    }
  }
}
</style>
