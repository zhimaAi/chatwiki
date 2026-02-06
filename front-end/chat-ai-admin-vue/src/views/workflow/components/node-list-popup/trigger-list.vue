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
      <a-input-search v-model:value.trim="keyword" allowClear :placeholder="t('ph_search_by_name')" />
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
            <div class="right-icon" v-if="item.subMenus && item.subMenus.length">
              <DownOutlined v-if="item.expend" /> <RightOutlined v-else />
            </div>
          </div>
          <div class="sub-menu-list" v-if="item.expend && item.subMenus && item.subMenus.length">
            <div
              class="sub-item"
              @click.stop="handleSubClick(item, sub)"
              v-for="sub in item.subMenus"
              :key="sub.value"
            >
              {{ sub.title }}
            </div>
          </div>
        </div>
      </div>
      <!-- <a class="more-link" href="/#/plugins/index?active=2" target="_blank">{{ t('btn_more_plugins') }}
        <RightOutlined />
      </a> -->
    </template>
    <div v-else class="empty-box">
      <img style="height: 200px" src="@/assets/empty.png" />
      <div>{{ t('msg_no_plugins') }}</div>
      <a href="/#/plugins/index?active=2" target="_blank"
        >{{ t('btn_add_now') }}
        <RightOutlined />
      </a>
    </div>
  </div>
</template>

<script setup>
import { createTriggerNode } from '../node-list'
import { ref, computed } from 'vue'
import { useWorkflowStore } from '@/stores/modules/workflow'
import { RightOutlined, DownOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.workflow.components.node-list-popup.index')

const emit = defineEmits(['add'])

const workflowStore = useWorkflowStore()
const triggerOfficialList = computed(() => workflowStore.triggerOfficialList)

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
  if (item.subMenus && item.subMenus.length) {
    item.expend = !item.expend
    return
  }
  let node = createTriggerNode(item)
  emit('add', node)
}

const handleSubClick = (item, sub) => {
  if (item.trigger_type == 4) {
    item.trigger_official_config.msg_type = sub.value

    item.outputs = triggerOfficialList.value.find((item) => item.msg_type == sub.value).fields
  }
  let node = createTriggerNode(item)

  emit('add', node)
}
</script>
