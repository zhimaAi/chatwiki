<style lang="less" scoped>
.setting-box {
  .actions-box {
    display: flex;
    align-items: center;
    line-height: 22px;
    font-size: 14px;
    color: #595959;

    .action-btn {
      cursor: pointer;
    }

    .save-btn {
      color: #2475fc;
    }
  }

  .library-list {
    display: flex;
    flex-flow: row wrap;
    gap: 16px;
    padding: 0 16px 16px 16px;

    .library-item {
      position: relative;
      width: 336px;
      padding: 14px 16px;
      border-radius: 2px;
      border: 1px solid #d8dde5;
      background-color: #fff;

      .library-name {
        width: 100%;
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      .library-intro {
        width: 100%;
        line-height: 20px;
        font-size: 12px;
        color: #8c8c8c;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      .close-btn {
        position: absolute;
        top: 0;
        right: 6px;
        font-size: 16px;
        color: #8c8c8c;
        cursor: pointer;
      }
    }
  }
}
</style>

<template>
  <edit-box
    class="setting-box"
    title="技能"
    icon-name="skii"
    v-model:isEdit="isEdit"
    :bodyStyle="{ padding: 0 }"
  >
    <template #tip>
      <a-tooltip placement="top">
        <template #title>
          <span>支持关联工作流</span>
        </template>
        <QuestionCircleOutlined />
      </a-tooltip>
    </template>
    <template #extra>
      <div class="actions-box">
        <a-flex :gap="8">
          <a-button size="small" @click="handleOpenSelectLibraryAlert">添加技能</a-button>
        </a-flex>
      </div>
    </template>
    <div class="library-list" v-if="selectedLibraryRows.length > 0">
      <div class="library-item" v-for="item in selectedLibraryRows" :key="item.id">
        <span class="close-btn" @click="handleRemoveCheckedLibrary(item)">
          <CloseCircleOutlined />
        </span>
        <div class="library-name">{{ item.robot_name }}</div>
        <div class="library-intro">{{ item.robot_intro }}</div>
      </div>
    </div>
    <RobotSelectAlert ref="robotSelectAlertRef" @change="onChangeLibrarySelected" />
  </edit-box>
</template>

<script setup>
import { relationWorkFlow } from '@/api/robot/index'
import { ref, reactive, inject, watchEffect, computed, toRaw } from 'vue'
import { CloseCircleOutlined, QuestionCircleOutlined } from '@ant-design/icons-vue'
import EditBox from '../edit-box.vue'
import RobotSelectAlert from './robot-select-alert.vue'
import { message } from 'ant-design-vue'
const isEdit = ref(false)

const { robotInfo, getRobot } = inject('robotInfo')
const props = defineProps({
  robotList: {
    type: Array,
    default: () => []
  }
})

const formState = reactive({
  work_flow_ids: []
})

// 知识库
const robotSelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  return props.robotList.filter((item) => {
    return formState.work_flow_ids.includes(item.id)
  })
})

// 移除知识库
const handleRemoveCheckedLibrary = (item) => {
  let index = formState.work_flow_ids.indexOf(item.id)

  formState.work_flow_ids.splice(index, 1)

  onSave()
}

const onChangeLibrarySelected = (checkedList) => {
  formState.work_flow_ids = [...checkedList]

  onSave()
}

const handleOpenSelectLibraryAlert = () => {
  robotSelectAlertRef.value.open([...formState.work_flow_ids])
}

const onSave = () => {
  let formData = { ...toRaw(formState) }

  formData.work_flow_ids = formData.work_flow_ids.join(',')

  relationWorkFlow({
    id: robotInfo.id,
    ...formData
  }).then((res) => {
    message.success('保存成功')
    getRobot(robotInfo.id)
  })
}

watchEffect(() => {
  formState.work_flow_ids = robotInfo.work_flow_ids.split(',')
})
</script>
