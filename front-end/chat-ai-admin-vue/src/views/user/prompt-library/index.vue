<template>
  <div class="user-model-page">
    <div class="page-title">
      提示词模板库
      <a-divider type="vertical" />
      <div class="desc">
        添加提示词模板，在应用配置和工作流的AI对话节点，可以从提示词模板库选择提示词并应用
      </div>
    </div>
    <div class="list-wrapper">
      <div class="left-group">
        <cu-scroll>
          <GroupList @change="handleChangeGroup" @load="handleLoadGroup" />
        </cu-scroll>
      </div>
      <div class="right-cotnetn">
        <cu-scroll>
          <div class="content-box">
            <div class="btn-block">
              <a-dropdown>
                <a-button type="primary" :icon="createVNode(PlusOutlined)">提示词</a-button>
                <template #overlay>
                  <a-menu>
                    <a-menu-item>
                      <div @click="handleAddWord(1)">添加结构化提示词</div>
                    </a-menu-item>
                    <a-menu-item>
                      <div @click="handleAddWord(0)">添加自定义提示词</div>
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>

              <a-button @click="handleAllExpend" :icon="createVNode(DatabaseOutlined)"
                >{{ hideStatus ? '全部展开' : '全部收起' }}
              </a-button>
            </div>
            <div class="loading-box" v-if="isLoading">
              <a-spin></a-spin>
            </div>
            <div class="prompt-list-box">
              <div class="prompt-list" v-for="item in lists" :key="item.id">
                <div class="prompt-header">
                  <div class="prompt-title">
                    <div class="prompt-type diy" v-if="item.prompt_type == 0">自定义</div>
                    <div class="prompt-type" v-if="item.prompt_type == 1">结构化</div>
                    <div class="title">{{ item.title }}</div>
                    <a-divider type="vertical" />
                    <div>{{ getGroupName(item.group_id) }}</div>
                  </div>
                  <div class="right-btn-box">
                    <!-- <UpOutlined /> -->
                    <a-tooltip :title="item.isHide ? '展开' : '收起'">
                      <div class="hover-btn-box" @click="handleHide(item)">
                        <DownOutlined v-if="item.isHide" />
                        <UpOutlined v-else />
                      </div>
                    </a-tooltip>

                    <div class="hover-btn-box" @click="handleEditWord(item)"><EditOutlined /></div>
                    <a-dropdown>
                      <div class="hover-btn-box"><EllipsisOutlined /></div>
                      <template #overlay>
                        <a-menu>
                          <a-menu-item>
                            <div @click="handleEditWord(item, 'copy')">复制</div>
                          </a-menu-item>
                          <a-menu-item>
                            <div @click="handleOpenGroup(item)">修改分组</div>
                          </a-menu-item>
                          <a-menu-item>
                            <div @click="handleDel(item)">删除</div>
                          </a-menu-item>
                        </a-menu>
                      </template>
                    </a-dropdown>
                  </div>
                </div>
                <div class="prompt-content-box">
                  <div v-if="item.prompt_type == 0" class="prompt-content">
                    <template v-if="item.isHide">
                      {{ item.prompt.slice(0, 200) }}
                      <span
                        @click="handleHide(item)"
                        style="cursor: pointer"
                        v-if="item.prompt.length > 200"
                        >...</span
                      >
                    </template>
                    <template v-else>{{ item.prompt }}</template>
                  </div>
                  <div v-else class="structure-list-box">
                    <div class="structure-list">
                      <div class="structure-title">{{ item.prompt_struct.role.subject }}</div>
                      <div class="structure-content">{{ item.prompt_struct.role.describe }}</div>
                    </div>
                    <div class="structure-list">
                      <div class="structure-title">{{ item.prompt_struct.task.subject }}</div>
                      <div class="structure-content">
                        {{ item.prompt_struct.task.describe }}
                        <span @click="handleHide(item)" style="cursor: pointer" v-if="item.isHide"
                          >...</span
                        >
                      </div>
                    </div>
                    <template v-if="!item.isHide">
                      <div class="structure-list">
                        <div class="structure-title">
                          {{ item.prompt_struct.constraints.subject }}
                        </div>
                        <div class="structure-content">
                          {{ item.prompt_struct.constraints.describe }}
                        </div>
                      </div>
                      <div class="structure-list">
                        <div class="structure-title">
                          {{ item.prompt_struct.skill.subject }}
                        </div>
                        <div class="structure-content">
                          {{ item.prompt_struct.skill.describe }}
                        </div>
                      </div>
                      <div class="structure-list">
                        <div class="structure-title">{{ item.prompt_struct.output.subject }}</div>
                        <div class="structure-content">
                          {{ item.prompt_struct.output.describe }}
                        </div>
                      </div>
                      <div class="structure-list">
                        <div class="structure-title">{{ item.prompt_struct.tone.subject }}</div>
                        <div class="structure-content">{{ item.prompt_struct.tone.describe }}</div>
                      </div>
                      <div
                        class="structure-list"
                        v-for="custom in item.prompt_struct.custom"
                        :key="custom.subject"
                      >
                        <div class="structure-title">{{ custom.subject }}</div>
                        <div class="structure-content">{{ custom.describe }}</div>
                      </div>
                    </template>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <a-empty style="padding-top: 100px" v-if="lists.length == 0"></a-empty>
        </cu-scroll>
      </div>
    </div>
    <AddDiyPrompt :groupList="groupLists" @ok="getLists" ref="addDiyPromptRef" />
    <a-modal v-model:open="moveGroupModal" title="修改分组" @ok="handleMoveGroup">
      <div class="move-group-box">
        <div class="label-text">修改分组为</div>
        <a-select v-model:value="currentItem.group_id" style="width: 100%">
          <a-select-option v-for="item in groupLists" :value="item.id">{{
            item.group_name
          }}</a-select-option>
        </a-select>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import {
  getPromptLibraryItems,
  deletePromptLibraryItems,
  movePromptLibraryItems
} from '@/api/user/index.js'
import {
  ExclamationCircleOutlined,
  PlusOutlined,
  DatabaseOutlined,
  DownOutlined,
  UpOutlined,
  EditOutlined,
  EllipsisOutlined
} from '@ant-design/icons-vue'
import { reactive, ref, createVNode } from 'vue'
import { Modal, message } from 'ant-design-vue'
import GroupList from './components/group-list.vue'
import AddDiyPrompt from './components/add-diy-prompt.vue'

let group_id = -1
const handleChangeGroup = (item) => {
  group_id = item.id
  getLists()
}

let groupLists = ref([])
const handleLoadGroup = (group) => {
  groupLists.value = group.filter((item) => item.id >= 0)
}

const lists = ref([])
const isLoading = ref(false)
const getLists = () => {
  isLoading.value = true
  getPromptLibraryItems({
    group_id,
    page: 1,
    size: 9999
  })
    .then((res) => {
      let data = res.data.list || []
      data = data.map((item) => {
        return {
          ...item,
          isHide: true,
          prompt_struct: item.prompt_struct ? JSON.parse(item.prompt_struct) : {}
        }
      })
      lists.value = data
    })
    .finally(() => {
      isLoading.value = false
    })
}

getLists()

const addDiyPromptRef = ref(null)
const handleAddWord = (prompt_type) => {
  addDiyPromptRef.value.show({
    group_id: group_id == -1 ? 0 : group_id,
    prompt_type
  })
}

const handleEditWord = (data, type) => {
  addDiyPromptRef.value.show(
    {
      ...data
    },
    type
  )
}

function getGroupName(id) {
  return groupLists.value.filter((item) => item.id == id)[0]?.group_name
}

const handleHide = (item) => {
  item.isHide = !item.isHide
}

const hideStatus = ref(true)
const handleAllExpend = () => {
  hideStatus.value = !hideStatus.value
  lists.value = lists.value.map((item) => {
    return { ...item, isHide: hideStatus.value }
  })
}

const handleDel = (record) => {
  Modal.confirm({
    title: '删除确认',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定要删除提示词[${record.title}]吗？`,
    okText: '删除',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      deletePromptLibraryItems({ id: record.id }).then((res) => {
        message.success('删除成功')
        getLists()
      })
    },
    onCancel() {}
  })
}

const moveGroupModal = ref(false)
let currentItem = reactive({
  id: '',
  group_id: 0
})
const handleOpenGroup = (item) => {
  item = { ...item }
  currentItem.id = item.id
  currentItem.group_id = item.group_id > 0 ? item.group_id : 0
  moveGroupModal.value = true
}
const handleMoveGroup = () => {
  if (currentItem.id < 0) {
    return message.error('请选择要移动的分组')
  }
  movePromptLibraryItems({
    id: currentItem.id,
    group_id: currentItem.group_id
  }).then(() => {
    message.success('移动成功')
    moveGroupModal.value = false
    getLists()
  })
}
</script>

<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  .page-title {
    display: flex;
    align-items: center;
    padding: 12px 24px;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-weight: 600;
    border-bottom: 1px solid var(--07, #f0f0f0);
    .desc {
      font-size: 14px;
      color: #8c8c8c;
      font-weight: 400;
    }
  }
  .list-wrapper {
    flex: 1;
    background: #fff;
    overflow: hidden;
    display: flex;
  }
  .left-group {
    height: 100%;
    overflow: hidden;
    width: 280px;
    border-right: 1px solid var(--07, #f0f0f0);
  }
  .right-cotnetn {
    flex: 1;
    height: 100%;
    overflow: hidden;
    position: relative;
  }
}
.loading-box {
  position: absolute;
  top: 160px;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.content-box {
  padding: 24px;
  .btn-block {
    display: flex;
    align-items: center;
    gap: 8px;
  }
}

.prompt-list-box {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
  .prompt-list {
    padding: 16px;
    background: var(--09, #f2f4f7);
    border-radius: 6px;
  }
  .prompt-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #8c8c8c;
    .prompt-title {
      display: flex;
      align-items: center;
      gap: 4px;
      line-height: 22px;
      .ant-divider-vertical {
        margin: 0;
      }
    }
    .title {
      color: #262626;
      font-weight: 600;
    }
    .prompt-type {
      display: flex;
      align-items: center;
      height: 22px;
      border-radius: 6px;
      padding: 0 6px;
      width: fit-content;
      background: #e4d2fa;
      color: #7000ff;
      &.diy {
        background: #d4e3fc;
        color: #2475fc;
      }
    }
    .right-btn-box {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }
}

.prompt-content-box {
  margin-top: 6px;
  color: #3a4559;
  line-height: 22px;
}
.prompt-content {
  white-space: pre-line;
}

.structure-list-box {
  .structure-list {
    margin-bottom: 6px;
    border-bottom: 1px solid #d9d9d9;
    .structure-title {
      color: #242933;
      line-height: 22px;
    }
    .structure-content {
      margin-top: 4px;
      margin-bottom: 6px;
      color: #3a4559;
      font-size: 14px;
      line-height: 22px;
      white-space: pre-line;
    }
  }
}

.hover-btn-box {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 24px;
  width: 24px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
  &:hover {
    background: #e3e5ea;
  }
}
.text-over-box {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
}

.move-group-box {
  margin: 24px 0;
  .label-text {
    color: #262626;
    line-height: 22px;
    margin-bottom: 4px;
    &::before {
      content: '*';
      color: red;
      margin-right: 4px;
      font-size: 12px;
      font-weight: 600;
    }
  }
}
</style>
