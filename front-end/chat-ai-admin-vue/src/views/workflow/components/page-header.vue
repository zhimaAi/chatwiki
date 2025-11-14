<style lang="less" scoped>
.page-header {
  display: flex;
  align-items: center;
  height: 56px;
  padding: 0 16px;
  padding-left: 0;
  box-shadow: 1px 1px 4px 0px rgba(0, 0, 0, 0.1);
  z-index: 10;
  background: #f0f2f5;
  .ml8 {
    margin-left: 8px;
  }
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
      <template v-if="currentVersion == '' ">
        <div class="last-save-time" v-if="draftSaveTime && draftSaveTime.time">
          {{ draftSaveTime.time }} {{ draftSaveTime.type == 'handle' ? '手动' : '自动' }}保存草稿
        </div>
        <a-button class="save-draft" @click="handleSave">保存草稿</a-button>
        <a-button type="primary" :loading="props.saveLoading"  class="publish-robot" @click="handleRelease">发布机器人</a-button>
      </template>
      <a-tooltip title="历史发布详情">
        <a-button @click="getVersionRecord()"><ClockCircleOutlined /></a-button>
      </a-tooltip>
    </div>
    <RunTest ref="runTestRef" :start_node_params="start_node_params" @getGlobal="getGlobal" @save="handleSave" />
  </div>
</template>

<script setup>
import { ExclamationCircleFilled, CheckCircleFilled, ClockCircleOutlined } from '@ant-design/icons-vue'
import { ref,  computed } from 'vue'
// front-end\chat-ai-admin-vue\src\views\robot\robot-config\components\top-header.vue
import TopHeader from '@/views/robot/robot-config/components/top-header.vue'
import { useRouter } from 'vue-router'
import RunTest from './run-test/index.vue'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()
const robotInfo = computed(() => {
  return robotStore.robotInfo
})
const draftSaveTime = computed(() => {
  return robotStore.draftSaveTime
})

const router = useRouter()

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
  }
})

const runTestRef = ref(null)

const onBack = () => {
  router.push('/')
}

const onEdit = () => {
  emit('edit')
}

const handleSave = (type="handle") => {
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

const openRunTest = () => {
  runTestRef.value.open()
}

defineExpose({
  openRunTest
})

</script>
