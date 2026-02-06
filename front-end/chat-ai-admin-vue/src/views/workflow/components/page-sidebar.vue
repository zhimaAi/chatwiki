<style lang="less" scoped>
.sidebar-wrapper {
  position: relative;
  height: 100%;
  width: 256px;
  border-radius: 6px;
  z-index: 10;
  background: none;
  transition: width 0.2s ease;
  background-color: #fff;
  box-shadow: 4px 0 8px 0 rgba(0, 0, 0, 0.04);

  &:hover {
    .sidebar-handle {
      opacity: 1;
    }
  }

  .sidebar-handle-wrapper {
    position: absolute;
    top: 0;
    right: -18px;
    width: 18px;
    height: 100%;
    z-index: 100;
  }
  .sidebar-handle {
    position: absolute;
    right: 0;
    top: 50%;
    width: 12px;
    height: 26px;
    transform: translateY(-50%);
    cursor: pointer;
    transition: all 0.2s ease;
    opacity: 0;

    .handle-line {
      position: absolute;
      width: 4px;
      height: 13px;
      left: 4px;
      position: absolute;
      transition: all 0.2s ease;
      background-color: #bfbfbf;
    }

    .handle-line01 {
      top: 0;
      border-top-left-radius: 4px;
      border-top-right-radius: 4px;
      transform-origin: 50% 0;
    }

    .handle-line02 {
      bottom: 0;
      border-bottom-left-radius: 4px;
      border-bottom-right-radius: 4px;
      transform-origin: 50% 100%;
    }
  }

  .sidebar-handle:hover {
    .handle-line01 {
      background-color: #595959;
      transform: rotate(18deg) translateY(0);
      border-top-left-radius: 4px;
      border-top-right-radius: 4px;
      border-bottom-left-radius: 10px;
      height: 16px;
    }

    .handle-line02 {
      background-color: #595959;
      transform: rotate(-18deg) translateY(0);
      border-bottom-left-radius: 4px;
      border-bottom-right-radius: 4px;
      border-top-left-radius: 10px;
      height: 16px;
    }
  }

  .sidebar-container {
    position: relative;
    width: 100%;
    height: 100%;
    padding: 16px;
    border-radius: 6px;
    overflow: hidden;
  }

  .sidebar-menus {
    .sidebar-menu {
      display: flex;
      align-items: center;
      height: 40px;
      line-height: 40px;
      padding: 0 12px;
      border-radius: 4px;
      margin-bottom: 4px;

      &:last-child {
        margin-bottom: 0;
      }

      &:hover {
        cursor: pointer;
        background-color: #f0f0f0;
      }

      &.active {
        background-color: #e6efff;

        .menu-name,
        .menu-icon {
          color: #2475fc;
        }
      }

      .menu-icon {
        width: 16px;
        height: 16px;
        color: #a1a7b3;
      }

      .menu-name {
        flex: 1;
        margin-left: 12px;
        font-size: 14px;
        color: #595959;
        white-space: nowrap;
        overflow: hidden;
      }
    }
  }
}

.sidebar-wrapper.is-hide {
  width: 72px;

  .sidebar-handle {
    opacity: 1 !important;
  }

  .sidebar-handle:hover {
    .handle-line01 {
      background-color: #595959;
      transform: rotate(-18deg) translateY(0);
      border-bottom-left-radius: 4px;
      border-bottom-right-radius: 4px;
      border-top-left-radius: 10px;
      height: 16px;
    }

    .handle-line02 {
      background-color: #595959;
      transform: rotate(18deg) translateY(0);
      border-top-left-radius: 4px;
      border-top-right-radius: 4px;
      border-bottom-left-radius: 10px;
      height: 16px;
    }
  }
}
</style>

<template>
  <div class="sidebar-wrapper" :class="{ 'is-hide': sidebarHide }">
    <div class="sidebar-handle-wrapper" v-if="showSidebarBtn">
      <a-tooltip
        :mouseEnterDelay="0"
        :mouseLeaveDelay="0"
        placement="right"
        v-model:open="handleTooltipShow"
        :arrow="false"
      >
        <template #title>{{ sidebarHide ? t('tooltip_expand') : t('tooltip_collapse') }}</template>
        <span class="sidebar-handle" @click="onHandleClick">
          <span class="handle-line handle-line01"></span>
          <span class="handle-line handle-line02"></span>
        </span>
      </a-tooltip>
    </div>

    <div class="sidebar-container">
      <div class="sidebar-menus">
        <template v-for="item in menus" :key="item.value">
          <cu-tooltip :disabled="!sidebarHide" :title="item.label" placement="bottom">
            <div
              class="sidebar-menu"
              :class="{ active: activeMenuKey == item.value }"
              @click="onMenuClick(item)"
            >
              <svg-icon :name="item.iconName" class="menu-icon" size="16"></svg-icon>
              <span class="menu-name">{{ item.label }}</span>
           </div>
          </cu-tooltip>
        </template>
        
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGlobalStore } from '@/stores/modules/global'
import CuTooltip from '@/components/cu-tooltip/index.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.page-sidebar')

const globalStore = useGlobalStore()

const router = useRouter()
const route = useRoute()
const activeMenuKey = computed(() => route.path.split('/')[3])
const sidebarHide = ref(true)


const handleTooltipShow = ref(true)
function onHandleClick() {
  sidebarHide.value = !sidebarHide.value
  handleTooltipShow.value = false
  globalStore.setHideLayoutTopAndBottom(sidebarHide.value)
}

const props = defineProps({
  showSidebarBtn: {
    default: true,
    type: Boolean
  }
})

const menus = computed(() => [
  {
    label: t('menu_workflow'),
    value: 'workflow',
    path: '/robot/config/workflow',
    iconName: 'workflow'
  },
  // {
  //   label: '基础配置',
  //   value: 'basic-config',
  //   path: '/robot/config/basic-config',
  //   iconName: 'jichupeizhi'
  // },
  {
    label: t('menu_external_services'),
    value: 'external-services',
    path: '/robot/config/external-services',
    iconName: 'duiwaifuwu'
  },
  {
    label: t('menu_chat_test'),
    value: 'test',
    path: '/robot/test',
    iconName: 'liaotianceshi',
    isNewWindowOpen: true
  },
  // {
  //   label: '功能中心',
  //   value: 'function-center',
  //   path: '/robot/config/function-center',
  //   iconName: 'function-center'
  // },
  {
    label: t('menu_qa_feedback'),
    value: 'qa-feedbacks',
    path: '/robot/config/qa-feedbacks',
    iconName: 'qa-feedback'
  },
  {
    label: t('menu_session_record'),
    value: 'session-record',
    path: '/robot/config/session-record',
    iconName: 'session-record'
  },
  {
    label: t('menu_api_key_manage'),
    value: 'api-key-manage',
    path: '/robot/config/api-key-manage',
    iconName: 'duiwaifuwu'
  },
  // {
  //   label: '未知问题',
  //   value: 'unknown_issue',
  //   path: '/robot/config/unknown_issue',
  //   iconName: 'unknown-issue'
  // },
  {
    label: t('menu_statistical_analysis'),
    value: 'statistical_analysis',
    path: '/robot/config/statistical_analysis',
    iconName: 'statistical-analysis'
  },
  {
    label: t('menu_export_record'),
    value: 'export-record',
    path: '/robot/config/export-record',
    iconName: 'export-record'
  },
  {
    label: t('menu_invoke_logs'),
    value: 'invoke-logs',
    path: '/robot/config/invoke-logs',
    iconName: 'doc-file'
  }
])

const onMenuClick = (item) => {
  if(item.isNewWindowOpen){
    window.open(`/#${item.path}?robot_key=${route.query.robot_key}&id=${route.query.id}`)
    return
  }
  router.push({
    path: item.path,
    query: route.query
  })
}

onMounted(()=>{
  setTimeout(()=>{
    handleTooltipShow.value = false
  },3000)
})
</script>
