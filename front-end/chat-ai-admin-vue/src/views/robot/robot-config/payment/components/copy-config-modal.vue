<template>
  <a-modal
    title="复制设置到其他应用"
    v-model:open="visible"
    :confirm-loading="saving"
    width="669px"
    @ok="ok"
  >
    <a-alert type="info" class="zm-alert-info" message="将当前页面的收费设置以及功能开关状态，复制到选中的应用"/>
    <div class="data-list">
      <div
        v-for="(item, i) in apps"
        :key="i"
        :class="['data-item', {disabled: item.has_published != 1}]"
        @click="item.checked = !item.checked"
      >
        <img class="avatar" :src="item.robot_avatar"/>
        <div class="info">
          <div>
            <div class="tit">{{ item.robot_name }}</div>
            <div class="desc">{{ item.robot_intro }}</div>
          </div>
          <a-checkbox v-if="item.has_published == 1" v-model:checked="item.checked"/>
          <span v-else class="status-tag waiting"><ExclamationCircleFilled/> 未发布</span>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {ref} from 'vue';
import {getRobotList} from "@/api/robot/index.js";
import {message} from 'ant-design-vue';
import {copyPaymentSetting} from "@/api/robot/payment.js";
import {ExclamationCircleFilled} from '@ant-design/icons-vue';

const props = defineProps({
  robotId: {
    type: [Number, String]
  }
})
const visible = ref(false)
const loading = ref(false)
const saving = ref(false)
const apps = ref([])

function show() {
  loadApps()
  visible.value = true
}

function loadApps() {
  loading.value = true
  getRobotList().then(res => {
    const list = res?.data || []
    apps.value = list.filter(i => i.id != props.robotId)
  }).finally(() => {
    loading.value = false
  })
}

function ok() {
  let ids = apps.value.filter(i => i.checked).map(i => i.id)
  if (!ids.length) return message.error('请选择应用')
  saving.value = true
  copyPaymentSetting({
    from_robot_id: props.robotId,
    to_robot_id: ids.toString()
  }).then(res => {
    message.success('操作完成')
    visible.value = false
  }).finally(() => {
    saving.value = false
  })
}

defineExpose({
  show
})
</script>

<style scoped lang="less">
.data-list {
  max-height: 500px;
  overflow-y: auto;
  margin-top: 8px;

  .data-item {
    display: flex;
    align-items: center;
    padding: 24px 12px;
    border-bottom: 1px solid #F0F0F0;
    border-radius: 6px;
    cursor: pointer;

    &.disabled {
      opacity: 0.4;
      cursor: not-allowed;
    }

    &:hover {
      background: #F2F4F7;
    }

    .avatar {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      flex-shrink: 0;
      margin-right: 12px;
    }

    .info {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-size: 14px;

      .tit {
        color: #262626;
        font-weight: 600;
      }

      .desc {
        color: #8c8c8c;
      }
    }
  }
}

.mt24 {
  margin-top: 24px;
}

.status-tag {
  display: inline-block;
  padding: 2px 6px;
  align-items: center;
  gap: 4px;
  border-radius: 6px;
  background: #CAFCE4;
  white-space: nowrap;
  font-size: 14px;
  font-weight: 500;

  &.finished {
    background-color: #CAFCE4;
    color: #21A665;
  }

  &.waiting {
    background-color: #F0F0F0;
    color: #595959;
  }
}
</style>
