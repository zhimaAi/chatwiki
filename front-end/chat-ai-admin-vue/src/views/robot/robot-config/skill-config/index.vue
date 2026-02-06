<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  background-color: #f2f4f7;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 16px;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }

  .list-wrapper {
    background: #fff;
    flex: 1;
    overflow-x: hidden;
    overflow-y: auto;
  }
  .content-wrapper {
    padding: 16px;
    padding-top: 0;
  }
  .actions-box {
    margin: 16px 0;
    display: flex;
  }
  .avatar-box {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    img {
      width: 32px;
      height: 32px;
    }
  }

  .status-block{
    display: flex;
    gap: 4px;
    align-items: center;
    color: #8c8c8c;
    span{
      width: 8px;
      height: 8px;
      border-radius: 8px;
      display: block;
      background: #8c8c8c;
    }
    &.success{
      color: #52C41A;
      span{
        background: #52C41A;
      }
    }
  }
}
</style>

<template>
  <div class="user-model-page">
    <div class="page-title">{{ t('title_workflow') }}</div>
    <div class="list-wrapper">
      <div class="content-wrapper">
        <a-alert show-icon style="align-items: baseline">
          <template #message>
            <div>
              {{ t('msg_workflow_tip_1') }}
            </div>
            <div>{{ t('msg_workflow_tip_2') }}</div>
          </template>
        </a-alert>
        <div class="actions-box">
          <a-button type="primary" :icon="h(PlusOutlined)" @click="handleOpenSelectLibraryAlert"
            >{{ t('btn_add_skill') }}</a-button
          >
        </div>
        <a-table :data-source="selectedLibraryRows">
          <a-table-column key="robot_name" :title="t('col_workflow')">
            <template #default="{ record }">
              <div class="avatar-box">
                <img :src="record.robot_avatar" alt="" />
                <div>{{ record.robot_name }}</div>
              </div>
            </template>
          </a-table-column>
          <a-table-column key="robot_intro" data-index="robot_intro" :title="t('col_workflow_desc')" >
            <template #default="{ record }">
              {{ record.robot_intro || '--' }}
            </template>
          </a-table-column>
          <a-table-column key="start_node_key" :title="t('col_status')">
            <template #default="{ record }">
              <div v-if="!record.start_node_key" class="status-block"><span></span>{{ t('status_unpublished') }}</div>
              <div v-else class="status-block success"><span></span>{{ t('status_published') }}</div>
            </template>
          </a-table-column>
          <a-table-column key="start_node_key" :title="t('col_actions')">
            <template #default="{ record }">
              <a @click="handleRemoveCheckedLibrary(record)">{{ t('action_remove') }}</a>
            </template>
          </a-table-column>
        </a-table>
      </div>
    </div>
    <RobotSelectAlert ref="robotSelectAlertRef" @change="onChangeLibrarySelected" />
  </div>
</template>

<script setup>
import { getRobotList, relationWorkFlow } from '@/api/robot/index.js'
import { storeToRefs } from 'pinia'
import { useRoute } from 'vue-router'
import { CloseCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import RobotSelectAlert from '../basic-config/components/skill/robot-select-alert.vue'
import { message } from 'ant-design-vue'
import { useRobotStore } from '@/stores/modules/robot'
import { reactive, ref, computed, watchEffect, toRaw, h } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.skill-config.index')
const query = useRoute().query
const robotStore = useRobotStore()

const { robotInfo } = storeToRefs(robotStore)
const { getRobot } = robotStore

const robotList = ref([])

const formState = reactive({
  work_flow_ids: []
})

// 知识库
const robotSelectAlertRef = ref(null)
const selectedLibraryRows = computed(() => {
  return robotList.value.filter((item) => {
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
    id: query.id,
    ...formData
  }).then((res) => {
    message.success(t('msg_save_success'))
    getRobot(query.id)
  })
}

function getRobotData() {
  getRobotList().then((res) => {
    robotList.value = res.data || []
  })
}

getRobotData()

watchEffect(() => {
  formState.work_flow_ids = robotInfo.value.work_flow_ids.split(',')
})
</script>
