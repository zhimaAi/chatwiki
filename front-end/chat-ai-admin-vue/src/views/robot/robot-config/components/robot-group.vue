<template>
  <div class="robot-group-box">
    <div class="left-group-box" :class="{ 'hide-group': isHideGroup }">
      <div class="popover-scroll-box">
        <div class="group-head-box">
          <div class="head-title">
            <div>机器人分组</div>
            <a-tooltip title="新建分组">
              <div class="hover-btn-wrap" @click="openGroupModal({})"><PlusOutlined /></div>
            </a-tooltip>
          </div>
          <div class="search-box">
            <a-input
              v-model:value="groupSearchKey"
              allowClear
              placeholder="搜索分组"
              style="width: 100%"
            >
              <template #suffix>
                <SearchOutlined @click.stop="" />
              </template>
            </a-input>
          </div>
        </div>
        <div class="classify-box">
          <div
            class="classify-item"
            @click="handleChangeGroup(item)"
            :class="{ active: item.id == group_id }"
            v-for="item in filterGroupLists"
            :key="item.id"
          >
            <div class="classify-title">{{ item.group_name }}</div>
            <div class="right-content">
              <div class="num" :class="{ 'num-block': item.id <= 0 }">{{ item.total }}</div>
              <div class="btn-box" v-if="item.id > 0">
                <a-dropdown placement="bottomRight">
                  <div class="hover-btn-wrap">
                    <EllipsisOutlined />
                  </div>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item>
                        <div @click.stop="openGroupModal(item)">重命名</div>
                      </a-menu-item>
                      <a-menu-item>
                        <div @click.stop="handleDelGroup(item)">删 除</div>
                      </a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </div>
            </div>
          </div>
        </div>
      </div>
      <a-tooltip placement="right" :title="isHideGroup ? '展开分组' : '收起分组'">
        <div class="hide-group-box" @click="handleChangeHideGroup">
          <LeftOutlined v-if="!isHideGroup" />
          <RightOutlined v-else />
        </div>
      </a-tooltip>
    </div>
    <div class="list-box popover-scroll-box">
      <div class="empty-box" v-if="filterRobotList.length == 0">
        <a-empty></a-empty>
      </div>
      <div
        class="list-item-wrapper"
        :class="{ 'show-more': isHideGroup }"
        v-for="item in filterRobotList"
        :key="item.id"
        @click="handleChangeRobot(item)"
      >
        <div class="list-item" :class="{ active: item.id == props.value }">
          <div class="robot-info">
            <img class="robot-avatar" :src="item.robot_avatar" alt="" />
            <div class="robot-info-content">
              <div class="robot-name">{{ item.robot_name }}</div>
              <div class="robot-type-tag">
                {{ item.application_type == 0 ? '聊天机器人' : '工作流' }}
              </div>
            </div>
          </div>
          <div class="robot-desc">{{ item.robot_intro }}</div>
        </div>
      </div>
    </div>
    <AddGroup ref="addGroupRef" @ok="initData" />
  </div>
</template>

<script setup>
import { computed, ref, createVNode } from 'vue'
import { storeToRefs } from 'pinia'
import {
  ExclamationCircleOutlined,
  PlusOutlined,
  SearchOutlined,
  EllipsisOutlined,
  LeftOutlined,
  RightOutlined
} from '@ant-design/icons-vue'
import { deleteRobotGroup } from '@/api/robot/index.js'
import { Modal, message } from 'ant-design-vue'
import AddGroup from '@/views/robot/robot-list/components/add-group.vue'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()

const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:value', 'change'])

let hideGroupLocalKey = 'robot-list-hide-group-key'

const isHideGroup = ref(localStorage.getItem(hideGroupLocalKey) == 1)

const handleChangeHideGroup = () => {
  isHideGroup.value = !isHideGroup.value
  localStorage.setItem(hideGroupLocalKey, isHideGroup.value ? 1 : 0)
}

const group_id = ref('')

const robotList = computed(() => {
  return robotStore.robotList || []
})

const filterRobotList = computed(() => {
  let lists = robotList.value
  if (group_id.value != '') {
    lists = lists.filter((item) => item.group_id == group_id.value)
  }
  return lists
})

const robotGroupList = computed(() => {
  let lists = robotStore.robotGroupList || []
  let totalNumber = 0
  // 计算每个分组的机器人数量
  lists.forEach((group) => {
    totalNumber += +group.total
  })
  return [
    {
      group_name: '全部',
      total: totalNumber,
      id: ''
    },
    ...lists
  ]
})

const groupSearchKey = ref('')

const filterGroupLists = computed(() => {
  return robotGroupList.value.filter((item) => item.group_name.includes(groupSearchKey.value))
})

const initData = () => {
  robotStore.getRobotLists()
  robotStore.getGroupList()
}

const handleChangeGroup = (item) => {
  group_id.value = item.id
}

const addGroupRef = ref(null)
const openGroupModal = (data) => {
  addGroupRef.value.show({
    ...data
  })
}

const handleDelGroup = (item) => {
  Modal.confirm({
    title: `确认删除分组${item.group_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: '确认',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      deleteRobotGroup({
        id: item.id
      }).then(() => {
        message.success('删除成功')
        robotStore.getGroupList()
        if (group_id.value == item.id) {
          group_id.value = ''
          robotStore.getRobotLists()
        }
      })
    }
  })
}

const handleChangeRobot = (item) => {
  emit('update:value', item.id)
  emit('change', item.id)
}

defineExpose({
  group_id
})
</script>

<style lang="less" scoped>
.robot-group-box {
  width: 800px;
  display: flex;
}
.left-group-box {
  width: 200px;
  margin-right: 16px;
  border: 1px solid #d9d9d9;
  padding: 12px;
  border-radius: 6px;
  position: relative;
  &.hide-group {
    width: 0;
    padding: 0;
    border-left: 0;
    border-top: 0;
    border-bottom: 0;
  }
  .hide-group-box {
    position: absolute;
    right: -8px;
    top: 40%;
    height: 50px;
    width: 13px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #e5e5ea;
    color: #8c8c8c;
    cursor: pointer;
    opacity: 0.78;
    &:hover {
      opacity: 1;
    }
  }
  .head-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
  .search-box {
    margin-top: 16px;
  }

  .classify-box {
    flex: 1;
    overflow: hidden;
    font-size: 14px;
    .classify-item {
      height: 32px;
      padding: 0 8px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-top: 4px;
      cursor: pointer;
      border-radius: 6px;
      color: #595959;
      transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);

      .classify-title {
        flex: 1;
        text-overflow: ellipsis;
        overflow: hidden;
        white-space: nowrap;
      }
      .num {
        display: block;
      }

      .btn-box {
        display: none;
      }
      &:hover {
        background: #f2f4f7;
        .num {
          display: none;
        }
        .num.num-block {
          display: block;
        }
        .btn-box {
          display: block;
        }
      }
      &.active {
        color: #2475fc;
        background: #e6efff;
      }
    }
  }

  .hover-btn-wrap {
    width: fit-content;
    height: 24px;
    border-radius: 6px;
    padding: 0 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
    &:hover {
      background: #e4e6eb;
    }
  }
}

.list-box {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding-right: 12px;
  height: fit-content;
  .list-item-wrapper {
    width: calc(50% - 4px);
    &.show-more {
      width: calc(33% - 5px);
    }
  }
  .list-item {
    position: relative;
    width: 100%;
    padding: 24px;
    border: 1px solid #e4e6eb;
    border-radius: 12px;
    background-color: #fff;
    transition: all 0.25s;
    cursor: pointer;
    &.active {
      border: 1px solid #2475fc;
    }
    &:hover {
      box-shadow: 0 4px 16px 0 rgba(0, 0, 0, 0.12);
    }

    .robot-info {
      display: flex;
      align-items: center;
    }

    .robot-avatar {
      width: 38px;
      height: 38px;
      border-radius: 14px;
      overflow: hidden;
    }

    .robot-info-content {
      flex: 1;
      padding-left: 12px;
      overflow: hidden;
    }

    .robot-name {
      height: 24px;
      line-height: 24px;
      margin-bottom: 4px;
      font-size: 16px;
      font-weight: 600;
      color: rgb(38, 38, 38);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .robot-desc {
      height: 44px;
      line-height: 22px;
      margin-top: 12px;
      font-size: 14px;
      font-weight: 400;
      color: rgb(89, 89, 89);
      // 超出2行显示省略号
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
    }

    .robot-type-tag {
      display: inline-block;
      height: 22px;
      line-height: 20px;
      padding: 0 8px;
      font-size: 12px;
      font-weight: 400;
      border-radius: 6px;
      color: rgb(36, 117, 252);
      border: 1px solid #cde0ff;
    }
  }
}

.empty-box {
  width: 100%;
  padding-top: 70px;
}

.popover-scroll-box {
  max-height: 400px;
  min-height: 180px;
  margin-top: 4px;
  overflow: hidden;
  overflow-y: auto;
  /* 整个页面的滚动条 */
  &::-webkit-scrollbar {
    width: 6px; /* 垂直滚动条宽度 */
    height: 6px; /* 水平滚动条高度 */
  }

  /* 滚动条轨道 */
  &::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 10px;
  }

  /* 滚动条滑块 */
  &::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 10px;
    transition: background 0.3s ease;
  }

  /* 滚动条滑块悬停状态 */
  &::-webkit-scrollbar-thumb:hover {
    background: #555;
  }

  /* 滚动条角落 */
  &::-webkit-scrollbar-corner {
    background: #f1f1f1;
  }
}
</style>
