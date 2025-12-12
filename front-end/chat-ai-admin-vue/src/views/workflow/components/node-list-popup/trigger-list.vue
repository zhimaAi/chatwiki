<style lang="less" scoped>
.trigger-list-box {
  .search-box {
    margin: 0 8px 8px;
  }

  .trigger-list {
    .trigger-item {
      .trigger-info {
        display: flex;
        align-items: center;
        padding: 4px 8px;
        border-radius: 6px;
        cursor: pointer;
      }

      .avatar {
        width: 20px;
        height: 20px;
        flex-shrink: 0;
        border-radius: 4px;
        margin-right: 8px;
      }
    }
  }
  .empty-box {
    text-align: center;
  }
}
</style>

<template>
  <div class="trigger-list-box">
    <div class="search-box" @click.stop="">
      <a-input-search v-model:value.trim="keyword" allowClear placeholder="请输入名称查询" />
    </div>
    <template v-if="triggerList.length">
      <div class="trigger-list">
        <div
          class="trigger-item"
          v-for="item in triggerList"
          @click="handleAddNode(item)"
          :key="item.type"
        >
          <div class="trigger-info">
            <img class="avatar" :src="item.trigger_icon" />
            <div class="info">
              <span class="name">{{ item.trigger_name }}</span>
            </div>
          </div>
        </div>
      </div>
      <!-- <a class="more-link" href="/#/plugins/index?active=2" target="_blank">更多插件
        <RightOutlined />
      </a> -->
    </template>
    <div v-else class="empty-box">
      <img style="height: 200px" src="@/assets/empty.png" />
      <div>暂无可用插件</div>
      <a href="/#/plugins/index?active=2" target="_blank"
        >去添加
        <RightOutlined />
      </a>
    </div>
  </div>
</template>

<script setup>
import { createTriggerNode } from '../node-list'
import { ref, computed } from 'vue'
import { useWorkflowStore } from '@/stores/modules/workflow'
import { RightOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['add'])

const workflowStore = useWorkflowStore()

// 触发器列表
const keyword = ref('')
const triggerList = computed(() => {
  let list = workflowStore.triggerList
  if (keyword.value) {
    return list.filter((item) => item.trigger_name.includes(keyword.value))
  }
  return list
})

const handleAddNode = (item) => {
  let node = createTriggerNode(item)
  emit('add', node)
}
</script>
