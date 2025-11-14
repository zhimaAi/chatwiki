<template>
  <a-modal v-model:open="visible" title="关联工作流" width="669px" @ok="ok">
    <a-alert class="zm-alert-info mt24" type="info" message="选择工作流后，会自动生成唯一标识符，可修改" show-icon/>
    <div class="data-list">
      <div
        v-for="(item, i) in list"
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
import {ref, reactive} from 'vue';
import {ExclamationCircleFilled} from '@ant-design/icons-vue';
import {getRobotList} from "@/api/robot/index.js";

const emit = defineEmits(['ok'])

const visible = ref(false)
const list = ref([])
const selected = ref([])

function show(ids = []) {
  selected.value = ids
  visible.value = true
  loadData()
}

function loadData() {
  getRobotList({application_type: 1}).then(res => {
    list.value = res.data || []
    list.value.forEach(item => {
      item.checked = selected.value.includes(item.id)
    })
  })
}

function ok() {
  visible.value = false
  emit('ok', list.value.filter(i => i.checked).map(i => i.id), list)
}

defineExpose({
  show,
})
</script>

<style scoped lang="less">
.data-list {
  max-height: 60vh;
  overflow-y: auto;

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
