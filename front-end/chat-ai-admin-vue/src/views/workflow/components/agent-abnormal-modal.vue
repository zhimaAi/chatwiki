<template>
  <a-modal
    v-model:open="open"
    title="Agent节点异常请重新选择"
    :width="520"
    :maskClosable="false"
    :keyboard="false"
    okText="确定"
    cancelText="取消"
    @ok="handleConfirm"
    @cancel="handleCancel"
  >
    <div class="agent-abnormal-list">
      <a-radio-group v-model:value="selectedAgentRobotId">
        <a-radio
          v-for="item in robotOptions"
          :key="item.id"
          :value="Number(item.id)"
        >
          <span class="agent-option">
            <img
              v-if="item.robot_avatar"
              class="agent-avatar"
              :src="item.robot_avatar"
              alt=""
            />
            <span class="agent-name">{{ item.robot_name }}</span>
          </span>
        </a-radio>
      </a-radio-group>
    </div>
  </a-modal>
</template>

<script setup>
import { ref } from 'vue'

const open = ref(false)
const robotOptions = ref([])
const selectedAgentRobotId = ref()
let resolveModal = null

const show = (list = []) => {
  robotOptions.value = list
  selectedAgentRobotId.value = list.length ? Number(list[0].id) : undefined
  open.value = true

  return new Promise((resolve) => {
    resolveModal = resolve
  })
}

const handleConfirm = () => {
  const selectedRobot = robotOptions.value.find(item => Number(item.id) === Number(selectedAgentRobotId.value))
  open.value = false
  resolveModal?.(selectedRobot || null)
  resolveModal = null
}

const handleCancel = () => {
  open.value = false
  resolveModal?.(null)
  resolveModal = null
}

defineExpose({
  open: show
})
</script>

<style lang="less" scoped>
.agent-abnormal-list {
  padding-top: 4px;

  :deep(.ant-radio-wrapper) {
    display: flex;
    align-items: center;
    margin-bottom: 14px;
    color: #595959;
    font-size: 16px;
    font-weight: 500;
  }

  .agent-option {
    display: inline-flex;
    align-items: center;
    min-width: 0;
  }

  .agent-avatar {
    width: 28px;
    height: 28px;
    margin-right: 12px;
    border-radius: 6px;
    object-fit: cover;
  }

  .agent-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
