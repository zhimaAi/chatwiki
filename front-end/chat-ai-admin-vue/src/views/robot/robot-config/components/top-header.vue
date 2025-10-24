<template>
  <div class="robot-name-box" ref="robotNameBoxRef">
    <addRobotAlert @addRobot="handleAddRobot" @editRobot="handleEditRobot" ref="addRobotAlertRef" />
    <div class="left-icon hover-block" @click="goBack">
      <HomeOutlined />
    </div>
    <div class="select-robot-block">
      <a-popover placement="bottomLeft" trigger="click" :overlayInnerStyle="{ 'padding-right': 0 }">
        <template #content>
          <RobotGroup ref="robotGroupRef" v-model:value="robotId" @change="handleChangeRobot" />
        </template>

        <div class="robot-option-block show-block-item">
          <img class="robot-avatar" :src="currentRobotItem.robot_avatar" alt="" />
          <div class="robot-name">{{ currentRobotItem.robot_name }}</div>
          <DownOutlined class="down-icon" />
        </div>
      </a-popover>
    </div>
    <a-tooltip v-if="routeName == 'robotWorkflow'">
      <template #title>编辑机器人</template>
      <div class="hover-block" @click="toEditRobot()"><EditOutlined /></div>
    </a-tooltip>
    <a-dropdown v-if="robotCreate">
      <a-tooltip>
        <template #title>新建机器人</template>
        <div class="hover-block" @click.prevent=""><PlusCircleOutlined /></div>
      </a-tooltip>

      <template #overlay>
        <a-menu>
          <a-menu-item @click.prevent="toAddRobot(0)">
            <span class="create-action">
              <img class="icon" :src="DEFAULT_ROBOT_AVATAR" alt="" />
              <span>聊天机器人</span>
            </span>
          </a-menu-item>
          <a-menu-item @click.prevent="toAddRobot(1)">
            <span class="create-action">
              <img class="icon" :src="DEFAULT_WORKFLOW_AVATAR" alt="" />
              <span>工作流</span>
            </span>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import addRobotAlert from '../../robot-list/components/add-robot-alert.vue'
import { DEFAULT_ROBOT_AVATAR, DEFAULT_WORKFLOW_AVATAR } from '@/constants/index.js'
import { HomeOutlined, EditOutlined, PlusCircleOutlined, DownOutlined } from '@ant-design/icons-vue'
import { useRouter, useRoute } from 'vue-router'
import { usePermissionStore } from '@/stores/modules/permission'
import { useRobotStore } from '@/stores/modules/robot'
import RobotGroup from './robot-group.vue'

const router = useRouter()
const route = useRoute()

const robotGroupRef = ref(null)

const routeName = computed(() => {
  return route.name
})

const robotStore = useRobotStore()
const { robotInfo } = storeToRefs(robotStore)
const goBack = () => {
  router.push('/robot/list')
}

const addRobotAlertRef = ref(null)

const permissionStore = usePermissionStore()
let { role_permission } = permissionStore
const robotCreate = computed(() => role_permission.includes('RobotCreate'))

const robotList = computed(() => {
  return robotStore.robotList || []
})

const robotId = ref(robotInfo.value.id)
const robotNameBoxRef = ref(null)

const currentRobotItem = computed(() => {
  return robotList.value.find((item) => item.id == robotId.value) || {}
})

const toAddRobot = (val) => {
  addRobotAlertRef.value.open(val)
  let group_id = robotGroupRef.value ? robotGroupRef.value.group_id : '0'
  addRobotAlertRef.value.setGroupId(group_id)
}

const toEditRobot = () => {
  let application_type = +robotInfo.value.application_type
  addRobotAlertRef.value.open(application_type, true)
  addRobotAlertRef.value.setGroupId(currentRobotItem.value?.group_id)
}

const handleAddRobot = () => {
  getRobotData()
  getGroupList()
}
const handleEditRobot = () => {
  getRobotData()
  getGroupList()
}

const handleChangeRobot = () => {
  let path = route.path
  let selectItem = robotList.value.find((item) => item.id == robotId.value)
  let { id, robot_key } = selectItem
  localStorage.setItem('last_local_robot_id', id)
  window.location.href = `/#${path}?id=${id}&robot_key=${robot_key}`
  window.location.reload()
}

function getRobotData() {
  robotStore.getRobotLists()
}

function getGroupList() {
  robotStore.getGroupList()
}
onMounted(() => {})
</script>

<style lang="less" scoped>
.robot-name-box {
  display: flex;
  align-items: center;
  padding: 16px 0;
  overflow: hidden;
  padding-left: 16px;
  padding-right: 4px;
  .hover-block {
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    width: fit-content;
    height: 36px;
    padding: 0 8px;
    border-radius: 6px;
    &:hover {
      background: var(--07, #e4e6eb);
    }
    transition: all 0.3s ease-in-out;
  }
  .left-icon {
    cursor: pointer;
    font-size: 14px;
  }
  .select-robot-block {
    position: relative;
    height: 36px;
    padding-left: 8px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    transition: all 0.3s ease-in-out;
    flex: 1;
    overflow: hidden;
    cursor: pointer;
    .show-block-item {
      width: 100%;
      .robot-name {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
    .down-icon {
      font-size: 12px;
      color: #8c8c8c;
      margin-right: 2px;
    }
    &:hover {
      background: var(--07, #e4e6eb);
    }
    &::v-deep(.ant-select) {
      .ant-select-selector {
        padding: 0;
      }
      .ant-select-selection-item {
        padding-right: 30px;
      }
      .ant-select-arrow {
        margin-top: -4px;
        color: #262626;
      }
    }
  }
  .robot-option-block {
    display: flex;
    align-items: center;
    gap: 8px;
    .robot-avatar {
      width: 32px;
      height: 32px;
      border-radius: 8px;
    }
    .robot-name {
      color: #262626;
      font-weight: 600;
      line-height: 24px;
    }
  }
}

.create-action {
  display: flex;
  align-items: center;
  .icon {
    width: 20px;
    height: 20px;
    margin-right: 8px;
  }
}
</style>
