<template>
  <div class="left-menu-box">
    <a-menu class="left-menu" :selectedKeys="selectedKeys" @click="handleChangeMenu">
      <router-link
        class="default-color"
        tag="a"
        :target="item.target || '_self'"
        :to="{ path: item.path, query: item.query || query }"
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
import { getRobotPermission } from '@/utils/permission'
const emit = defineEmits(['changeMenu'])
const route = useRoute()
const props = defineProps({
  robotInfo: {
    type: Object,
    default: () => {}
  }
})

const query = route.query
const selectedKeys = computed(() => {
  return [route.path.split('/')[3]]
})

const baseItems = [
  {
    key: 'workflow',
    id: 'workflow',
    icon: () =>
      h(SvgIcon, {
        name: 'workflow',
        class: 'menu-icon'
      }),
    label: '工作流编排',
    title: '工作流编排',
    path: '/robot/config/workflow',
    menuIn: ['1']
  },
  {
    key: 'basic-config',
    id: 'basic-config',
    icon: () =>
      h(SvgIcon, {
        name: 'jichupeizhi',
        class: 'menu-icon'
      }),
    label: '基础配置',
    title: '基础配置',
    path: '/robot/config/basic-config',
    menuIn: ['0']
  },
  {
    key: 'library-config',
    id: 'library-config',
    icon: () =>
      h(SvgIcon, {
        name: 'guanlianzhishiku',
        class: 'menu-icon'
      }),
    label: '知识库',
    title: '知识库',
    path: '/robot/config/library-config',
    menuIn: ['0']
  },
  {
    key: 'skill-config',
    id: 'skill-config',
    icon: () =>
      h(SvgIcon, {
        name: 'skii',
        class: 'menu-icon'
      }),
    label: '工作流',
    title: '工作流',
    path: '/robot/config/skill-config',
    menuIn: ['0']
  },
  {
    key: 'external-services',
    id: 'external-services',
    icon: () =>
      h(SvgIcon, {
        name: 'duiwaifuwu',
        class: 'menu-icon'
      }),
    label: '对外服务',
    title: '对外服务',
    path: '/robot/config/external-services',
    menuIn: ['0', '1']
  },
  {
    key: 'test',
    id: 'test',
    icon: () =>
      h(SvgIcon, {
        name: 'liaotianceshi',
        class: 'menu-icon'
      }),
    label: '聊天测试',
    title: '聊天测试',
    path: '/robot/test',
    query: {
      robot_key: props.robotInfo.robot_key,
      id: props.robotInfo.id
    },
    target: '_blank',
    menuIn: ['0', '1']
  },
  {
    key: 'qa-feedbacks',
    id: 'qa-feedbacks',
    icon: () =>
      h(SvgIcon, {
        name: 'qa-feedback',
        class: 'menu-icon'
      }),
    label: '问答反馈',
    title: '问答反馈',
    path: '/robot/config/qa-feedbacks',
    menuIn: ['0', '1']
  },
  {
    key: 'session-record',
    id: 'session-record',
    icon: () =>
      h(SvgIcon, {
        name: 'session-record',
        class: 'menu-icon'
      }),
    label: '会话记录',
    title: '会话记录',
    path: '/robot/config/session-record',
    query: {
      robot_key: props.robotInfo.robot_key,
      id: props.robotInfo.id
    },
    menuIn: ['0', '1']
  },
  {
    key: 'api-key-manage',
    id: 'api-key-manage',
    icon: () =>
      h(SvgIcon, {
        name: 'duiwaifuwu',
        class: 'menu-icon'
      }),
    label: 'API Key管理',
    title: 'API Key管理',
    path: '/robot/config/api-key-manage',
    menuIn: ['0', '1']
  },
  {
    key: 'unknown_issue',
    id: 'unknown_issue',
    icon: () =>
      h(SvgIcon, {
        name: 'unknown-issue',
        class: 'menu-icon'
      }),
    label: '未知问题',
    title: '未知问题',
    path: '/robot/config/unknown_issue',
    menuIn: ['0', '1']
  },
  {
    key: 'statistical_analysis',
    id: 'statistical_analysis',
    icon: () =>
      h(SvgIcon, {
        name: 'statistical-analysis',
        class: 'menu-icon'
      }),
    label: '统计分析',
    title: '统计分析',
    path: '/robot/config/statistical_analysis',
    menuIn: ['0', '1']
  },
  {
    key: 'export-record',
    id: 'export-record',
    icon: () =>
      h(SvgIcon, {
        name: 'export-record',
        class: 'menu-icon'
      }),
    label: '导出记录',
    title: '导出记录',
    path: '/robot/config/export-record',
    menuIn: ['0', '1']
  }
]

const items = computed(() => {
  let lists = baseItems
  if (getRobotPermission(query.id) == 1) {
    return lists.filter((item) => item.id == 'external-services')
  }
  return lists.filter((item) => item.menuIn.includes(props.robotInfo?.application_type))
})

const handleChangeMenu = ({ item }) => {
  if (selectedKeys.value.includes(item.id)) {
    return
  }
  return
  emit('changeMenu', item)
}
</script>

<style lang="less" scoped>
.default-color {
  color: inherit;
}
.left-menu-box {
  .left-menu {
    border-right: 0 !important;

    ::v-deep(.menu-icon) {
      color: #a1a7b3;
      font-size: 16px;
      vertical-align: -3px;
    }

    ::v-deep(.ant-menu-item-selected .menu-icon) {
      color: #2475fc;
    }
  }
}
</style>
