<template>
  <div class="left-menu-box">
    <a-menu class="left-menu" mode="inline" :selectedKeys="selectedKeys" :openKeys="openKeys" @click="handleChangeMenu">
      <template v-for="item in items" :key="item.key">
        <a-sub-menu v-if="item.children && item.children.length" :key="item.key" :icon="item.icon" :class="{'submenu-selected': isTopActive(item)}">
          <template #title>
            <router-link
              class="default-color"
              :to="{ path: item.path, query: item.query || query }"
              :target="item.target || '_self'"
              @click.stop
            >{{ item.label }}</router-link>
          </template>
          <router-link
            class="default-color"
            v-for="child in item.children"
            :key="child.key"
            :to="{ path: child.path, query: child.query || query }"
            :target="child.target || '_self'"
          >
            <a-menu-item :icon="child.icon" :path="child.path" :key="child.key">{{ child.label }}</a-menu-item>
          </router-link>
        </a-sub-menu>
        <router-link
          v-else
          class="default-color"
          :target="item.target || '_self'"
          :to="{ path: item.path, query: item.query || query }"
        >
          <a-menu-item :icon="item.icon" :path="item.path" :key="item.key">{{ item.label }}</a-menu-item>
        </router-link>
      </template>
    </a-menu>
  </div>
</template>

<script setup>
import { ref, h, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import SvgIcon from '@/components/svg-icon/index.vue'
import { getRobotPermission } from '@/utils/permission'
import { getRobotAbilityList } from '@/api/explore'
import { useRobotStore } from '@/stores/modules/robot'
const emit = defineEmits(['changeMenu'])
const route = useRoute()
const robotStore = useRobotStore()
const props = defineProps({
  robotInfo: {
    type: Object,
    default: () => {}
  }
})

const query = route.query
const selectedKeys = computed(() => {
  const functionCenter = items.value.find((i) => i.id === 'function-center')
  const children = functionCenter?.children || []
  const autoReplyMenu = children.find((i) => i.id === 'auto-reply')
  if (route.path.split('/')[3] === 'auto-reply' && !autoReplyMenu) {
    return ['function-center']
  }

  const subscribeReplyMenu = children.find((i) => i.id === 'subscribe-reply')
  if (route.path.split('/')[3] === 'subscribe-reply' && !subscribeReplyMenu) {
    return ['function-center']
  }
  const smartMenu = children.find((i) => i.id === 'smart-menu')
  if (route.path.split('/')[3] === 'smart-menu' && !smartMenu) {
    return ['function-center']
  }
  const paymentMenu = children.find((i) => i.id === 'payment')
  if (route.path.split('/')[3] === 'payment' && !paymentMenu) {
    return ['function-center']
  }
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
  // {
  //   key: 'function-center',
  //   id: 'function-center',
  //   icon: () =>
  //     h(SvgIcon, {
  //       name: 'function-center',
  //       class: 'menu-icon'
  //     }),
  //   label: '功能中心',
  //   title: '功能中心',
  //   path: '/robot/config/function-center',
  //   menuIn: ['0', '1']
  // },
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
  // {
  //   key: 'unknown_issue',
  //   id: 'unknown_issue',
  //   icon: () =>
  //     h(SvgIcon, {
  //       name: 'unknown-issue',
  //       class: 'menu-icon'
  //     }),
  //   label: '未知问题',
  //   title: '未知问题',
  //   path: '/robot/config/unknown_issue',
  //   menuIn: ['0', '1']
  // },
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
  },
  {
    key: 'invoke-logs',
    id: 'invoke-logs',
    icon: () =>
      h(SvgIcon, {
        name: 'doc-file',
        class: 'menu-icon'
      }),
    label: '调用日志',
    title: '调用日志',
    path: '/robot/config/invoke-logs',
    menuIn: ['1']
  }
]

const autoReplyMenu = ref(null)
const subscribeReplyMenu = ref(null)
const smartMenu = ref(null)
const paymentMenu = ref(null)

async function refreshAbilityMenu () {
  try {
    const rid = props.robotInfo?.id || query.id
    if (!rid) return
    const res = await getRobotAbilityList({ robot_id: rid })
    const data = res?.data || []

    // 关键词回复
    const autoItem = (data || []).find((it) => it?.ability_type === 'robot_auto_reply')
    if (autoItem) {
      const sw = autoItem?.robot_config?.switch_status ?? autoItem?.user_config?.switch_status ?? '0'
      const ai_reply_status = autoItem?.robot_config?.ai_reply_status ?? autoItem?.user_config?.ai_reply_status ?? '0'
      robotStore.setKeywordReplySwitchStatus(String(sw))
      robotStore.setKeywordReplyAiReplyStatus(String(ai_reply_status))
    } else {
      robotStore.setKeywordReplySwitchStatus('0')
      robotStore.setKeywordReplyAiReplyStatus('0')
    }
    const hit = (data || []).find(
      (it) =>
        it?.ability_type === 'robot_auto_reply' &&
        it?.robot_config?.fixed_menu === '1'
    )
    if (hit) {
      autoReplyMenu.value = {
        key: 'auto-reply',
        id: 'auto-reply',
        icon: () =>
          h(SvgIcon, {
            name: 'auto-reply',
            class: 'menu-icon'
          }),
        label: hit?.menu?.name,
        title: hit?.menu?.name,
        path: hit?.menu?.path || '/robot/ability/auto-reply',
        menuIn: ['0', '1']
      }
    } else {
      autoReplyMenu.value = null
    }

    // 关注后回复
    const subItem = (data || []).find((it) => it?.ability_type === 'robot_subscribe_reply')
    if (subItem) {
      const sw_sub = subItem?.robot_config?.switch_status ?? subItem?.user_config?.switch_status ?? '0'
      const ai_reply_statu_sub = subItem?.robot_config?.ai_reply_status ?? subItem?.user_config?.ai_reply_status ?? '0'
      robotStore.setSubscribeReplySwitchStatus(String(sw_sub))
      robotStore.setSubscribeReplyAiReplyStatus(String(ai_reply_statu_sub))
    } else {
      robotStore.setSubscribeReplySwitchStatus('0')
      robotStore.setSubscribeReplyAiReplyStatus('0')
    }
    const hitSub = false
    // const hitSub = (data || []).find(
    //   (it) =>
    //     it?.ability_type === 'robot_subscribe_reply' &&
    //     it?.robot_config?.fixed_menu === '1'
    //     && it?.robot_config?.switch_status === '1'
    // )
    if (hitSub) {
      subscribeReplyMenu.value = {
        key: 'subscribe-reply',
        id: 'subscribe-reply',
        icon: () =>
          h(SvgIcon, {
            name: 'subscribe-reply',
            class: 'menu-icon'
          }),
        label: hitSub?.menu?.name,
        title: hitSub?.menu?.name,
        path: hitSub?.menu?.path || '/explore/index/subscribe-reply',
        menuIn: ['0', '1']
      }
    } else {
      subscribeReplyMenu.value = null
    }

    const smartItem = (data || []).find((it) => it?.ability_type === 'robot_smart_menu')
    if (smartItem) {
      const sw_smart = smartItem?.robot_config?.switch_status ?? smartItem?.user_config?.switch_status ?? '0'
      const ai_reply_status_smart = smartItem?.robot_config?.ai_reply_status ?? smartItem?.user_config?.ai_reply_status ?? '0'
      robotStore.setSmartMenuSwitchStatus(String(sw_smart))
      robotStore.setSmartMenuAiReplyStatus(String(ai_reply_status_smart))
    } else {
      robotStore.setSmartMenuSwitchStatus('0')
      robotStore.setSmartMenuAiReplyStatus('0')
    }
    const hitSmart = (data || []).find(
      (it) =>
        it?.ability_type === 'robot_smart_menu' &&
        it?.robot_config?.fixed_menu === '1'
    )
    if (hitSmart) {
      smartMenu.value = {
        key: 'smart-menu',
        id: 'smart-menu',
        icon: () =>
          h(SvgIcon, {
            name: 'smart-menu',
            class: 'menu-icon'
          }),
        label: hitSmart?.menu?.name,
        title: hitSmart?.menu?.name,
        path: hitSmart?.menu?.path || '/robot/ability/smart-menu',
        menuIn: ['0', '1']
      }
    } else {
      smartMenu.value = null
    }

    const paymentItem = (data || []).find((it) => it?.ability_type === 'robot_payment')
    if (paymentItem) {
      const sw_payment = paymentItem?.robot_config?.switch_status ?? paymentItem?.user_config?.switch_status ?? '0'
      const ai_reply_status_payment = paymentItem?.robot_config?.ai_reply_status ?? paymentItem?.user_config?.ai_reply_status ?? '0'
      robotStore.setPaymentSwitchStatus(String(sw_payment))
      robotStore.setPaymentAiReplyStatus(String(ai_reply_status_payment))
    } else {
      robotStore.setPaymentSwitchStatus('0')
      robotStore.setPaymentAiReplyStatus('0')
    }
    const hitPayment = (data || []).find(
      (it) =>
        it?.ability_type === 'robot_payment' &&
        it?.robot_config?.fixed_menu === '1'
    )
    if (hitPayment) {
      paymentMenu.value = {
        key: 'payment',
        id: 'payment',
        icon: () =>
          h(SvgIcon, {
            name: 'payment',
            class: 'menu-icon'
          }),
        label: hitPayment?.menu?.name,
        title: hitPayment?.menu?.name,
        path: hitPayment?.menu?.path || '/robot/ability/payment',
        menuIn: ['0', '1']
      }
    } else {
      paymentMenu.value = null
    }
  } catch (e) {
    console.warn('refreshAbilityMenu failed', e)
  }
}

function abilityUpdatedHandler (e) {
  const rid = props.robotInfo?.id || query.id
  const incoming = e?.detail?.robotId
  if (!incoming || !rid || String(incoming) !== String(rid)) return
  refreshAbilityMenu()
}

onMounted(() => {
  refreshAbilityMenu()
  window.addEventListener('robotAbilityUpdated', abilityUpdatedHandler)
})

onUnmounted(() => {
  window.removeEventListener('robotAbilityUpdated', abilityUpdatedHandler)
})

const items = computed(() => {
  let lists = baseItems
  if (getRobotPermission(query.id) == 1) {
    lists = lists.filter((item) => item.id == 'external-services')
  } else {
    lists = lists.filter((item) => item.menuIn.includes(props.robotInfo?.application_type))
  }
  const arr = [...lists]
  const idx = arr.findIndex((i) => i.id === 'function-center')
  if (idx >= 0) {
    const children = []
    if (autoReplyMenu.value) children.push(autoReplyMenu.value)
    if (subscribeReplyMenu.value) children.push(subscribeReplyMenu.value)
    if (smartMenu.value) children.push(smartMenu.value)
    if (paymentMenu.value) children.push(paymentMenu.value)
    arr[idx] = { ...arr[idx], children }
  }
  return arr
})

const openKeys = computed(() => {
  const fc = items.value.find((i) => i.id === 'function-center')
  const hasChildren = fc && Array.isArray(fc.children) && fc.children.length > 0
  if (!hasChildren) return []
  return ['function-center']
})

const isTopActive = (item) => {
  return route.path.split('/')[3] === item.id
}

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
    max-height: calc(100vh - 192px);
    overflow-y: auto;
    scrollbar-width: none;
    -ms-overflow-style: none;
    &::-webkit-scrollbar {
      width: 0;
      height: 0;
    }

    ::v-deep(.menu-icon) {
      color: #a1a7b3;
      font-size: 16px;
      vertical-align: -3px;
    }

    ::v-deep(.ant-menu-item-selected .menu-icon) {
      color: #2475fc;
    }
    ::v-deep(.submenu-selected > .ant-menu-submenu-title) {
      color: #2475fc;
    }
    ::v-deep(.submenu-selected > .ant-menu-submenu-title .menu-icon) {
      color: #2475fc;
    }
    ::v-deep(.ant-menu-submenu-selected >.ant-menu-submenu-title .menu-icon) {
      color: #2475fc;
    }
  }
}
</style>
