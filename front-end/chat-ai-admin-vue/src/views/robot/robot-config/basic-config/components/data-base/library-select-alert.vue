<style lang="less" scoped>
.library-checkbox-box {
  padding-top: 16px;

  .list-tools {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;
    margin-bottom: 8px;
  }

  .list-box {
    display: flex;
    flex-flow: row wrap;
    height: 388px;
    width: 100%;
    overflow-y: auto;
    align-content: flex-start;
    margin: 0 -8px;

    .list-item-wraapper {
      padding: 8px;
      width: 50%;
    }

    .list-item {
      width: 100%;
      padding: 14px 12px;
      border: 1px solid #f0f0f0;
      border-radius: 2px;

      &:hover {
        cursor: pointer;
        box-shadow: 0 4px 16px 0 #1b3a6929;
      }

      .library-name {
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #262626;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }

      .library-desc {
        line-height: 20px;
        margin-top: 2px;
        font-size: 12px;
        font-weight: 400;
        color: #8c8c8c;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
    }

    .list-item :deep(span:last-child) {
      flex: 1;
      overflow: hidden;
    }
  }
}
.empty-box {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  padding-top: 40px;
  padding-bottom: 40px;
  color: #8c8c8c;
  img {
    width: 150px;
  }
}
</style>

<template>
  <a-modal width="746px" v-model:open="show" :title="t('title_link_data_table')" @ok="saveCheckedList">
    <div class="library-checkbox-box">
      <!-- <a-alert message="请选择关联知识库，机器人会根据知识库内上传的文档回复用户的提问。每个机器人最多关联5个知识库" type="info" /> -->
      <!-- <div class="list-tools">
        <div>
          <a-input style="width: 282px" v-model:value="searchKeyword" placeholder="请输入知识库名称搜索" @change="onSearch">
            <template #suffix>
              <SearchOutlined style="color: rgba(0, 0, 0, 0.25)" />
            </template>
          </a-input>
        </div>
        <div>
          <a-button style="margin-right: 8px" @click="onRefresh">
            <SyncOutlined /> 刷新
          </a-button>
          <a-button type="primary" ghost @click="openAddLibrary">新建知识库</a-button>
        </div>
      </div> -->
      <a-spin :spinning="isRefresh" :delay="100">
        <a-checkbox-group v-model:value="state.checkedList" style="width: 100%">
          <div class="list-box" ref="scrollContainer">
            <div class="list-item-wraapper" v-for="item in options" :key="item.id">
              <a-checkbox class="list-item" :value="item.id">
                <div class="library-name">{{ item.name }}</div>
                <div class="library-desc">{{ item.description }}</div>
              </a-checkbox>
            </div>
          </div>
        </a-checkbox-group>
      </a-spin>
      <div class="empty-box" v-if="!isRefresh && !options.length">
        <img src="@/assets/img/library/preview/empty.png" alt="" />
        <div>
          {{ t('msg_no_data_table_add_first') }}
          <a @click="openAddLibrary"> {{ t('btn_go_add') }}</a>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, SyncOutlined } from '@ant-design/icons-vue'
import { getFormList } from '@/api/database/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.basic-config.components.data-base.library-select-alert')

const emit = defineEmits(['change'])

const state = reactive({
  indeterminate: false,
  checkAll: false,
  checkedList: []
})

const show = ref(false)

const open = (checkedList) => {
  getList()

  state.checkedList = checkedList
  show.value = true
}

const saveCheckedList = () => {
  show.value = false
  triggerChange()
}

const searchKeyword = ref('')

const onSearch = () => {
  getList()
}

const options = ref([])

const triggerChange = () => {
  emit('change', [...state.checkedList])
}

const getList = async () => {
  const res = await getFormList()
  if (res) {
    let list = res.data || []
    options.value = list
  }
}

const isRefresh = ref(false)
const onRefresh = async () => {
  isRefresh.value = true

  await getList()

  isRefresh.value = false

  message.success(t('msg_refresh_success'))
}

const openAddLibrary = () => {
  window.open('/#/database/list')
}

defineExpose({
  open
})
</script>
