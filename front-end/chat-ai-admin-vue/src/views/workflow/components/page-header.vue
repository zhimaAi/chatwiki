<style lang="less" scoped>
.ml8 {
  margin-left: 8px;
}
.page-header {
  display: flex;
  align-items: center;
  height: 56px;
  padding: 0 16px;
  padding-left: 0;
  box-shadow: 1px 1px 4px 0px rgba(0, 0, 0, 0.1);
  z-index: 10;
  background: #f0f2f5;
  // border-bottom: 1px solid #e8eaec;
  .header-left {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 6px;
    color: #000;
    cursor: pointer;
    transition: all 0.2s ease-in;
    &:hover {
      background: #e4e6eb;
    }
  }
  .header-content {
    flex: 1;
    display: flex;
    align-items: center;
    // padding-left: 20px;
  }

  .edit-box {
    margin-left: 4px;

    .edit-box-icon {
      font-size: 18px;
      cursor: pointer;

      &:hover {
        color: #2475fc;
      }
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .back-btn {
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
  }

  .robot-avatar {
    width: 32px;
    height: 32px;
    margin-right: 8px;
    border-radius: 8px;
  }

  .robot-name {
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
  .robot-status-box {
    display: flex;
    gap: 2px;
    align-items: center;
    height: 22px;
    width: fit-content;
    font-size: 14px;
    font-weight: 500;
    border-radius: 6px;
    padding: 0 6px;
    .robot-status-icon {
      font-size: 16px;
    }

    &.status-1 {
      background: #ccdfff;
      color: #2475fc;
    }

    &.status-0 {
      background: #bfbfbf;
      color: #fff;
    }
  }

  .last-save-time {
    font-size: 14px;
    color: #8c8c8c;
    display: flex;
    flex-direction: column;
    align-items: center;
  }

}
.version-name {
  color: #262626;
}
.lock-tip {
  background:#FFEBCC;
  position: absolute;
  height: 40px;
  width: 100%;
  top: 56px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: #000000;

  .lock-info {
    color: #8c8c8c;
  }
}
.auto-save-hint {
  color: #FAAD14; /* Ant Design 的警示黄色 */
  font-size: 12px;
}
</style>

<style lang="less">
.popover-material {
  .ant-popover-arrow {
    bottom: 2px !important;
  }
  .ant-popover-inner {
    box-shadow:
      0 -3px 6px rgba(0, 0, 0, 0.1),
      0 2px 8px rgba(0, 0, 0, 0.15),
      -3px 0 6px rgba(0, 0, 0, 0.1),
      3px 0 6px rgba(0, 0, 0, 0.1);
  }
  .ant-popover-inner-content {
    padding: 0;
  }
  .lock-config {
    display: flex;
    align-items:center;
    margin-top: 24px;
    padding-left: 24px;
  }
}
</style>

<template>
  <div class="page-header">
    <!-- <div class="header-left">
      <div class="back-btn" @click="onBack"><LeftOutlined /></div>
    </div> -->
    <div class="header-content">
      <!-- <a-image :width="32" :src="robotInfo.robot_avatar_url" />
      <div class="robot-name ml8">{{ robotInfo.robot_name }}</div>
      <div class="edit-box" @click="onEdit">
        <svg-icon class="edit-box-icon" name="edit"></svg-icon>
      </div> -->
      <TopHeader />
      <template v-if="robotInfo.start_node_key != ''">
        <div class="robot-status-box status-1" @click="onEdit">
          <CheckCircleFilled class="robot-status-icon" />
          <div class="robot-status-text">已发布</div>
        </div>
      </template>
      <template v-else>
        <div class="robot-status-box status-0">
          <ExclamationCircleFilled class="robot-status-icon" />
          <div class="robot-status-text">未发布</div>
        </div>
      </template>
    </div>
    <div class="header-right" v-if="props.showRight">
      <template v-if="currentVersion == '' && isEditing">
        <div class="last-save-time" v-if="robotInfo.draft_save_time && robotInfo.draft_save_time != 0">
          {{ formatTime(robotInfo.draft_save_time, 'MM/DD HH:mm:ss') }} {{ robotInfo.draft_save_type == 'handle' ? '手动' : '自动' }}保存草稿
          <span v-if="isEditing && !isLockedByOther && (!autoSaveEnabled || !isLeader)" class="auto-save-hint">你可能已在其他页面编辑，已停止自动保存</span>
        </div>
        <template v-if="!isLockedByOther">
          <a-button class="save-draft" @click="handleSave('handle')">保存草稿</a-button>
          <a-button type="primary" :loading="props.saveLoading"  class="publish-robot" @click="handleRelease">发布机器人</a-button>
        </template>
      </template>
      <!-- <template v-else>
        <div class="last-save-time" v-if="currentVersion && currentVersion.version">
         <strong class="version-name">v{{ currentVersion.version }}</strong> {{ formatTime(currentVersion?.create_time) }}
        </div>
      </template> -->
      <a-button v-if="!isEditing && !currentVersion || isLockedByOther" :disabled="isLockedByOther" type="primary" class="ml8" @click="handleEdit">
        <template #icon>
          <EditOutlined />
        </template>
        编辑
      </a-button>
      <a-tooltip title="历史发布详情">
        <a-button @click="getVersionRecord()"><ClockCircleOutlined /></a-button>
      </a-tooltip>
      <!-- 管理员显示锁定图标 -->
      <a-popover
        v-if="isAdmin"
        trigger="click"
        placement="bottomRight"
        v-model:open="lockPopoverOpen"
        :overlay-style="{
          width: '372px',
          padding: '16px'
        }"
        overlay-class-name="popover-material"
        @onOpenChange="visibleChange"
      >
        <template #title>
          <div style="display:flex;align-items:center;gap:8px;" >
            <svg-icon
              class="lock-icon"
              name="lock-icon"
              style="font-size: 14px; color: #262626"
            ></svg-icon>
            <span style="color: #262626;">编辑锁设置</span>
            <span style="color:#8c8c8c; font-weight: normal;">不同的IP和UserAgent都会限制</span>
          </div>
        </template>
        <template #content>
          <div class="lock-config">
            <span style="color: #595959;">编辑结束后锁</span>
            <a-input-number
              v-model:value="lockMinutes"
              :min="10"
              :max="60"
              :precision="0"
              style="margin:0 4px;width:86px"
            />
            <span style="color: #595959;">分钟后才能编辑</span>
          </div>
          <div style="margin-top: 24px;text-align:right;">
            <a-button class="ml8" @click="lockPopoverOpen=false">取 消</a-button>
            <a-button type="primary" class="ml8" :loading="lockLoading" @click="saveLockConfig">保 存</a-button>
          </div>
        </template>
        <a-button class="" @click="lockPopoverOpen=true" style="padding:5px 9px;">
          <svg-icon
            class="lock-icon"
            name="lock-icon"
            style="font-size: 14px;"
          ></svg-icon>
        </a-button>
      </a-popover>
    </div>
    <div class="lock-tip" v-if="isLockedByOther">
      <div>当前已有用户正在编辑中，无法编辑保存草稿或发布工作流</div>
      <div class="lock-info">（{{ '用户：' + loginUserName || '--' }}  {{ 'IP：' + lockRemoteAddr || '--' }}  {{ 'User Agent：' + lockUserAgent || '--' }}）</div>
    </div>
    <RunTest ref="runTestRef" :start_node_params="start_node_params" @getGlobal="getGlobal" @save="handleSave" :isLockedByOther="isLockedByOther" />
  </div>
</template>

<script setup>
import { ExclamationCircleFilled, CheckCircleFilled, ClockCircleOutlined, EditOutlined } from '@ant-design/icons-vue'
import { computed, ref, onMounted } from 'vue'
import dayjs from 'dayjs'
// front-end\chat-ai-admin-vue\src\views\robot\robot-config\components\top-header.vue
import TopHeader from '@/views/robot/robot-config/components/top-header.vue'
import RunTest from './run-test/index.vue'
import { useRobotStore } from '@/stores/modules/robot'
import { useUserStore } from '@/stores/modules/user'
import { saveDraftExTime, getAdminConfig } from '@/api/robot/index'
import { message } from 'ant-design-vue'
const robotStore = useRobotStore()
const robotInfo = computed(() => {
  return robotStore.robotInfo
})

// const router = useRouter()

const emit = defineEmits(['save', 'release', 'edit', 'getGlobal', 'getVersionRecord'])
const props = defineProps({
  saveLoading: {
    default: false,
    type: Boolean
  },
  showRight: {
    default: true,
    type: Boolean
  },
  start_node_params: {
    default: () => ({}),
    type: Object
  },
  currentVersion:{
    default: '',
    type: String
  },
  isEditing: { default: true, type: Boolean },
  isLockedByOther: { default: false, type: Boolean },
  lockRemoteAddr: { default: '', type: String },
  lockUserAgent: { default: '', type: String },
  loginUserName: { default: '', type: String },
  autoSaveEnabled: { default: true, type: Boolean },
  isLeader: { default: true, type: Boolean }
})
const runTestRef = ref(null)
// 锁定设置仅超管/管理员显示
const userStore = useUserStore()
const isAdmin = computed(() => {
  const role = userStore.userInfo?.role_type
  // 约定：1=超级管理员，2=管理员（若后端不同可调整）
  return role == 1 || role == 2
})

// 编辑锁设置弹层
const lockPopoverOpen = ref(false)
const lockMinutes = ref(10)
const lockLoading = ref(false)

const saveLockConfig = async () => {
  lockLoading.value = true
  try {
    await saveDraftExTime({ draft_exptime: lockMinutes.value })
    message.success('编辑锁设置已保存')
    lockPopoverOpen.value = false
  } catch (e) {
    message.error('保存失败，请稍后重试')
  } finally {
    lockLoading.value = false
  }
}

onMounted(async () => {
  try {
    const res = await getAdminConfig()
    const mins = +res?.data?.draft_exptime
    if (typeof mins === 'number') {
      lockMinutes.value = Math.max(10, Math.min(60, Math.round(mins)))
    }
  } catch (e) {
    // console.warn('getAdminConfig failed', e)
  }
})

const visibleChange = () => {}

// const onBack = () => {
//   router.push('/')
// }

const onEdit = () => {
  emit('edit')
}

const handleSave = (type = "handle") => {
  emit('save', type)
}

const handleRelease = () => {
  emit('release')
}

const getGlobal = () => {
  emit('getGlobal')
}

const getVersionRecord = () => {
  emit('getVersionRecord')
}

const handleEdit = () => {
  emit('edit')
}
const openRunTest = () => {
  runTestRef.value.open()
}
function formatTime(time, formatType = 'YY-MM-DD HH:mm:ss') {
  if (!time || time <= 0) {
    return '--'
  }
  return dayjs(time * 1000).format(formatType)
}
defineExpose({
  openRunTest
})
</script>
