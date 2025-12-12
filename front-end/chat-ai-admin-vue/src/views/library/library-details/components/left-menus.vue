<template>
  <div class="left-menu-box">
    <a-menu class="left-menu" :selectedKeys="selectedKeys" @click="handleChangeMenu">
      <router-link
        class="default-color"
        :to="{ path: item.path, query: query }"
        v-for="item in menus"
        :key="item.key"
      >
        <a-menu-item :icon="item.icon" :path="item.path" :key="item.key" v-if="!item.hidden">
          {{item.label}}
          <span v-if="item.syncStatus == 2" class="sync-tag run">同步中</span>
          <span v-else-if="item.syncStatus == 3" class="sync-tag fail">同步失败</span>
        </a-menu-item>
      </router-link>
    </a-menu>
  </div>
</template>

<script setup>
import { h, computed } from 'vue'
import { useRoute } from 'vue-router'
import SvgIcon from '@/components/svg-icon/index.vue'
import { useUserStore } from '@/stores/modules/user'
import { useLibraryStore } from '@/stores/modules/library'

const libraryStore = useLibraryStore()

const graph_switch = computed(() => {
  return libraryStore.graph_switch
})

const userStore = useUserStore()
const emit = defineEmits(['changeMenu'])
const route = useRoute()
const query = route.query

const robot_nums = computed(() => {
  return userStore.getRobotNums
})

const selectedKeys = computed(() => {
  return [route.path.split('/')[3]]
})

const librarySyncStatus = computed(() => {
  return libraryStore.sync_official_content_status
})

const menus = computed(() => {

  let lists = [
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
      label: libraryStore.type == 2 ? '知识管理' : '知识库文档' ,
      title: '知识库文档',
      path: '/library/details/knowledge-document',
      syncStatus: Number(librarySyncStatus.value)
    },
    {
      key: 'knowledge-graph',
      id: 'knowledge-graph',
      hidden: graph_switch.value == 0,
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
      label: '知识图谱',
      title: '知识图谱',
      path: '/library/details/knowledge-graph'
    },
    {
      key: 'categary-manage',
      id: 'categary-manage',
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
      label: '精选',
      title: '精选',
      path: '/library/details/categary-manage'
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
    },
    {
      key: 'related-robots',
      id: 'related-robots',
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
      label: `关联机器人${robot_nums.value > 0 ? ` (${robot_nums.value})` : ''}`,
      path: '/library/details/related-robots'
    }
  ]
  if (libraryStore.type == 2) {
    lists.push({
      key: 'import-record',
      id: 'import-record',
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
      label: '导入记录',
      title: '导入记录',
      path: '/library/details/import-record'
    })
    lists.push({
      key: 'export-record',
      id: 'export-record',
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
      label: '导出记录',
      title: '导出记录',
      path: '/library/details/export-record'
    })
  }else{
    lists.push({
      key: 'recycle-bin-record',
      id: 'recycle-bin-record',
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
      label: '回收站',
      title: '回收站',
      path: '/library/details/recycle-bin-record'
    })
  }
  return lists
})

const handleChangeMenu = ({ item }) => {
  if (selectedKeys.value.includes(item.id)) {
    return
  }

  emit('changeMenu', item)
}
</script>

<style lang="less" scoped>
.left-menu-box {
  height: 100%;
  width: 232px;
  .default-color {
    color: inherit;
  }
  .left-menu {
    height: 100%;
    border-right: 0 !important;

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

.sync-tag {
  width: 100%;
  background: rgba(0,0,0,0.5);
  color: #FFF;
  font-size: 10px;
  border-radius: 12px;
  padding: 2px 6px;
  &.run {
    background: #2475FC;
  }
  &.fail {
    background: #FB363F;
  }
}
</style>
