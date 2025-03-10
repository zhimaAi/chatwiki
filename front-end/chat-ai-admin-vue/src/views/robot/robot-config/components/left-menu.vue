<template>
  <div class="left-menu-box">
    <a-menu
      class="left-menu"
      :selectedKeys="selectedKeys"
      @click="handleChangeMenu"
    >
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
    path: '/robot/config/workflow'
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
    path: '/robot/config/basic-config'
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
    path: '/robot/config/external-services'
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
    path: '/robot/config/qa-feedbacks'
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
    }
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
    path: '/robot/config/api-key-manage'
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
    path: '/robot/config/statistical_analysis'
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
    path: '/robot/config/export-record'
  },
]

const items = computed(()=>{
  if(props.robotInfo.application_type == 0){
    return baseItems.filter(item => item.id != 'workflow')
  }else{
    return baseItems.filter(item => item.id != 'basic-config')
  }
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
.default-color{
  color: inherit;
}
.left-menu-box {
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
  }
}
</style>
