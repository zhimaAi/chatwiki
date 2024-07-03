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
    key: 'knowledge-document',
    id: 'knowledge-document',
    icon: () =>
      h('span', {}, [
        h(SvgIcon, {
          name: 'doc-menu-icon',
          class: 'menu-icon'
        }),
        h(SvgIcon, {
          name: 'doc-active-menu-icon',
          class: 'menu-icon-active'
        })
      ]),
    label: '知识库文档',

    title: '知识库文档',
    path: '/library/details/knowledge-document'
  },
  {
    key: 'recall-testing',
    id: 'recall-testing',
    icon: () =>
      h('span', {}, [
        h(SvgIcon, {
          name: 'test-menu-icon',
          class: 'menu-icon'
        }),
        h(SvgIcon, {
          name: 'test-active-menu-icon',
          class: 'menu-icon-active'
        })
      ]),
    label: '召回测试',
    title: '召回测试',
    path: '/library/details/recall-testing'
  },
  {
    key: 'knowledge-config',
    id: 'knowledge-config',
    icon: () =>
      h('span', {}, [
        h(SvgIcon, {
          name: 'knowledge-config-icon',
          class: 'menu-icon'
        }),
        h(SvgIcon, {
          name: 'knowledge-config-active-icon',
          class: 'menu-icon-active'
        })
      ]),
    label: '知识库配置',
    title: '知识库配置',
    path: '/library/details/knowledge-config'
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
    ::v-deep(.ant-menu-item .menu-icon-active) {
      display: none;
    }
    ::v-deep(.ant-menu-item-selected .menu-icon) {
      display: none;
    }
    ::v-deep(.ant-menu-item-selected .menu-icon-active) {
      display: block;
    }
  }
}
</style>
