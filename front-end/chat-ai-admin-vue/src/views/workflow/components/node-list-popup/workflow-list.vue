<style lang="less" scoped>
.trigger-list-box {
  .search-box {
    margin: 0 8px 8px;
  }

  .trigger-list {
    .trigger-item {
      &:hover {
        background: #e4e6eb;
        border-radius: 6px;
      }
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

      .right-icon {
        margin-left: auto;
        color: #8c8c8c;
        font-size: 12px;
        font-weight: 400;
      }
    }
  }

  .sub-menu-list {
    margin-left: 34px;
    border-left: 1px solid #d9d9d9;
    padding-left: 6px;
    display: flex;
    flex-direction: column;
    gap: 2px;

    .sub-item {
      height: 26px;
      display: flex;
      align-items: center;
      cursor: pointer;
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
      <a-input-search v-model:value.trim="keyword" allowClear :placeholder="t('ph_search_by_name')"/>
    </div>
    <template v-if="showList.length">
      <div class="trigger-list">
        <div
          class="trigger-item"
          v-for="item in showList"
          @click="handleAddNode(item)"
          :key="item.id"
        >
          <div class="trigger-info">
            <img class="avatar" :src="item.robot_avatar"/>
            <div class="info">
              <span class="name">{{ item.robot_name }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>
    <div v-else class="empty-box">
      <img style="height: 200px" src="@/assets/empty.png"/>
      <div>{{ t('msg_no_workflows_available') }}</div>
      <a href="/#/robot/list?active=1" target="_blank">{{ t('btn_add_now') }}
        <RightOutlined/>
      </a>
    </div>
  </div>
</template>

<script setup>
import {ref, computed, onMounted} from 'vue'
import {createWorkflowNode} from '../node-list'
import {RightOutlined, DownOutlined} from '@ant-design/icons-vue'
import {getRobotList} from "@/api/robot/index.js";
import {useI18n} from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-list-popup.workflow-list')

const emit = defineEmits(['add'])
const list = ref([])
const keyword = ref('')

onMounted(() => {
  loadData()
})

function loadData() {
  getRobotList({application_type: 1}).then(res => {
    let _list = res.data || []
    list.value = _list.filter(i => i.has_published == 1)
  })
}

const showList = computed(() => {
  if (keyword.value) {
    return list.value.filter((item) => item.robot_name.includes(keyword.value))
  }
  return list.value
})

const handleAddNode = (item) => {
  let node = createWorkflowNode(item)
  emit('add', node)
}
</script>
