<style lang="less" scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  border-radius: 6px;

  .nav-menu {
    display: flex;
    position: relative;
    padding: 9px 16px;
    margin-right: 4px;
    line-height: 22px;
    font-size: 14px;
    font-weight: 700;
    border-radius: 6px;
    color: #262626;
    cursor: pointer;
    transition: all 0.2s;
    &:hover {
      background: #e4e6eb;
    }
    &.active {
      color: #fff;
      background: #2475fc;
    }

    .nav-icon {
      margin-right: 8px;
      font-size: 14px;
    }
    .down-icon {
      margin-left: 16px;
      font-size: 14px;
    }
  }
}
.robot-ment-item .anticon-check {
  opacity: 0;
}
.robot-ment-item.robot-active-item .anticon-check {
  opacity: 1;
}
.robot-active-item {
  color: #2475fc;
}
</style>

<template>
  <div class="navbar-wrapper">
    <div class="navbar">
      <div class="nav-menu" v-if="role_type == 1" :class="{active: activeMenu == 'guide'}" @click="handleToGuide">
        <svg-icon class="nav-icon" name="guide"></svg-icon>
         <span class="nav-name">新手指引 ({{ total_process }}%)</span>
      </div>
      <template v-for="item in navs">
        <template v-if="checkRole(item.permission)" :key="item.key">
          <div
            class="nav-menu"
            v-if="item.key === 'robot'"
            :class="{ active: item.key === rootPath || item.key === activeMenu }"
            @click="handleRobotMenuClick(item)"
          >
            <svg-icon class="nav-icon" :name="item.icon"></svg-icon>
            <span class="nav-name">{{ item.title }}</span>
            <a-dropdown>
              <template #overlay>
                <a-menu>
                  <a-menu-item key="1" @click="handleChangeRobotmenuItem('robot_detail', item)">
                    <div
                      class="robot-ment-item"
                      :class="{ 'robot-active-item': activeRobotMenu == 'robot_detail' }"
                    >
                      <CheckOutlined class="anticon-check" /> 进入机器人详情
                    </div>
                  </a-menu-item>
                  <a-menu-item key="2" @click="handleChangeRobotmenuItem('robot_list', item)">
                    <div
                      class="robot-ment-item"
                      :class="{ 'robot-active-item': activeRobotMenu == 'robot_list' }"
                    >
                      <CheckOutlined /> 进入机器人管理
                    </div>
                  </a-menu-item>
                </a-menu>
              </template>
              <div class="down-icon">
                <DownOutlined />
              </div>
            </a-dropdown>
          </div>
          <div
            v-else
            class="nav-menu"
            :class="{ active: item.key === rootPath || item.key === activeMenu }"
            @click="handleClickNav(item)"
          >
            <svg-icon class="nav-icon" :name="item.icon"></svg-icon>
            <span class="nav-name">{{ item.title }}</span>
          </div>
        </template>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed, watch, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { DownOutlined, CheckOutlined } from '@ant-design/icons-vue'
import { checkRole } from '@/utils/permission'
import { getRobotList, robotAutoAdd } from '@/api/robot/index.js'
import { useGuideStore } from '@/stores/modules/guide'
import { useCompanyStore } from '@/stores/modules/company'
import { usePermissionStore } from '@/stores/modules/permission'
const companyStore = useCompanyStore()
const guideStore = useGuideStore()

const permissionStore = usePermissionStore()

const role_type = computed(() => {
  return permissionStore.role_type
})

const total_process = computed(() => {
  return +guideStore.total_process
})

const router = useRouter()
const roure = useRoute()

const activeMenu = computed(() => {
  return roure.meta.activeMenu || ''
})

const rootPath = computed(() => {
  return roure.path.split('/')[1]
})

const activeRobotMenu = ref(localStorage.getItem('local__robot_menu_key') || 'robot_detail')

const getActiveMenu = () => {}

const baseNavs = [
  {
    id: 0,
    key: 'explore',
    label: 'explore',
    title: '探索',
    icon: 'nav-explore',
    path: '/explore/index',
    permission: ['AbilityCenter']
  },
  {
    id: 1,
    key: 'robot',
    label: 'robot',
    title: '机器人',
    icon: 'nav-robot',
    path: '/robot/list',
    permission: ['RobotManage']
  },
  {
    id: 2,
    key: 'library',
    label: 'library',
    title: '知识库',
    icon: 'nav-library',
    path: '/library/list',
    permission: ['LibraryManage']
  },
  {
    id: 3,
    key: 'PublicLibrary',
    label: 'PublicLibrary',
    title: '文档',
    icon: 'nav-doc',
    path: '/public-library/list',
    permission: ['OpenLibDocManage']
  },
  {
    id: 4,
    key: 'library-search',
    label: 'library-search',
    title: '搜索',
    icon: 'search',
    path: '/library-search/index',
    permission: ['SearchManage']
  },
  {
    id: 5,
    key: 'chat-monitor',
    label: 'chat-monitor',
    title: '会话',
    icon: 'nav-chat',
    path: '/chat-monitor/index',
    permission: ['ChatSessionManage']
  }
  // {
  //   id: 6,
  //   key: 'user',
  //   label: 'user',
  //   title: '系统管理',
  //   path: '/user/model',
  //   permission: [
  //     'ModelManage',
  //     'TokenManage',
  //     'TeamManage',
  //     'AccountManage',
  //     'CompanyManage',
  //     'ClientSideManage'
  //   ]
  // }
]

const top_navigate = computed(() => {
  return companyStore.top_navigate
})

const navs = computed(() => {
  const openList = top_navigate.value.filter((item) => item.open) // 获取所有打开的菜单项

  return openList
    .map((item) => baseNavs.find((nav) => nav.key === item.id)) // 查找匹配的菜单项
    .filter(Boolean) // 过滤掉未找到的菜单项（undefined）
})

const handleClickNav = (item) => {
  guideStore.getUseGuideProcess()
  router.push(item.path)
  // window.open(`/#${item.path}`, "_blank", "noopener") // 建议添加 noopener 防止安全漏洞
}

const handleToGuide = () => {
  guideStore.getUseGuideProcess()
  router.push('/guide')
}

const handleChangeRobotmenuItem = (type, item) => {
  localStorage.setItem('local__robot_menu_key', type)
  activeRobotMenu.value = type
  handleRobotMenuClick(item)
}

const handleRobotMenuClick = (item) => {
  if (activeRobotMenu.value == 'robot_detail') {
    handleToRobotDetail()
  } else {
    handleClickNav(item)
  }
}

const handleToRobotDetail = async () => {
  try {
    const { data: lists = [] } = await getRobotList()

    if (lists.length === 0) {
      // 云版这里需要创建一个新的机器人
      router.push('/robot/list')
      // robotAutoAdd().then(res=>{
      //   window.open(`/#/robot/config/basic-config?id=${res.data.id}&robot_key=${res.data.robot_key}`)
      // })
      return
    }

    const localRobotId = localStorage.getItem('last_local_robot_id')
    let toDetailRobot

    if (localRobotId) {
      toDetailRobot = lists.find((item) => item.id == localRobotId)
    }

    toDetailRobot = toDetailRobot || lists[0]

    const { id, robot_key, application_type } = toDetailRobot
    router.push({
      path: `/robot/config/${application_type == 1 ? 'workflow' : 'basic-config'}`,
      query: {
        id: id,
        robot_key: robot_key
      }
    })
    // window.open(`/#/robot/config/basic-config?id=${id}&robot_key=${robot_key}`)
  } catch (error) {
    router.push('/robot/list')
  }
}

watch(
  () => roure.path,
  () => {
    getActiveMenu()
  },
  {
    immediate: true
  }
)

onMounted(()=>{
  guideStore.getUseGuideProcess()
})
</script>
