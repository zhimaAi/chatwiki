<template>
  <div class="group-wrapper">
    <div>
      <a-button @click="handldOpenAddModal({})" :icon="h(PlusOutlined)" block>{{ t('add_group') }}</a-button>
    </div>
    <div class="group-box">
      <div
        class="group-list"
        :class="{ active: item.id == groupId }"
        @click="handleChangeGroup(item)"
        v-for="item in groupList"
        :key="item.id"
      >
        <div class="group-name">
          <div class="name-text">{{ item.group_name }}</div>
          <a-tooltip v-if="item.group_desc">
            <template #title>{{ item.group_desc }}</template>
            <InfoCircleOutlined />
          </a-tooltip>
        </div>
        <div class="right-block" @click.stop v-if="item.id > 0">
          <a-dropdown>
            <div class="hover-btn-box">
              <EllipsisOutlined />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item>
                  <div @click="handldOpenAddModal(item)">{{ t('edit') }}</div>
                </a-menu-item>
                <a-menu-item>
                  <div style="color: #fb363f" @click="handleDelGroup(item)">{{ t('delete') }}</div>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
  </div>
  <AddGroup ref="addGroupRef" @ok="getGroupList" />
</template>

<script setup>
import {
  PlusOutlined,
  InfoCircleOutlined,
  EllipsisOutlined,
  ExclamationCircleOutlined
} from '@ant-design/icons-vue'
import { getPromptLibraryGroup, deletePromptLibraryGroup } from '@/api/user/index.js'
import { Modal, message } from 'ant-design-vue'
import AddGroup from './add-group.vue'
import { ref, h } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.prompt-library.components.group-list')

const emit = defineEmits(['change', 'load'])

const groupId = ref(-1)
const groupList = ref([])
const getGroupList = () => {
  getPromptLibraryGroup().then((res) => {
    groupList.value = [
      {
        id: -1,
        group_name: t('all')
      },
      {
        id: 0,
        group_name: t('default_group')
      },
      ...res.data
    ]
    emit('load', groupList.value)
  })
}

getGroupList()

const handleChangeGroup = (item) => {
  groupId.value = item.id
  emit('change', item)
}

const addGroupRef = ref(null)
const handldOpenAddModal = (data) => {
  addGroupRef.value.show({
    ...data
  })
}

const handleDelGroup = (item) => {
  Modal.confirm({
    title: t('delete_confirm'),
    icon: h(ExclamationCircleOutlined),
    content: t('delete_confirm_content', { name: item.group_name }),
    okText: t('delete'),
    okType: 'danger',
    cancelText: t('cancel'),
    onOk() {
      deletePromptLibraryGroup({ id: item.id }).then((res) => {
        if (item.id == groupId.value) {
          handleChangeGroup(groupList.value[0])
        }
        getGroupList()
        message.success(t('delete_success'))
      })
    },
    onCancel() {}
  })
}
</script>

<style lang="less" scoped>
.group-wrapper {
  padding: 24px;
}

.group-box {
  margin-top: 16px;
  .group-list {
    margin-bottom: 4px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 5px 8px;
    font-size: 14px;
    color: #595959;
    cursor: pointer;
    border-radius: 6px;
    transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
    .group-name {
      display: flex;
      align-items: center;
      gap: 4px;
    }
    .name-text {
      max-width: 160px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    &:hover {
      background: #f2f4f7;
      .hover-btn-box {
        display: flex;
      }
    }
    &.active {
      background: #e6efff;
      color: #2475fc;
    }

    .hover-btn-box {
      width: 24px;
      height: 24px;
      display: none;
      align-items: center;
      justify-content: center;
      border-radius: 6px;
      color: #262626;

      &:hover {
        background: #e4e6eb;
      }
    }
  }
}
</style>
