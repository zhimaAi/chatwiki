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
</style>

<template>
  <a-modal width="746px" v-model:open="show" title="关联知识库" @ok="saveCheckedList">
    <div class="library-checkbox-box">
      <a-alert message="请选择关联知识库，机器人会根据知识库内上传的文档回复用户的提问。每个机器人最多关联5个知识库" type="info" />
      <div class="list-tools">
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
      </div>
      <a-spin :spinning="isRefresh" :delay="100">
        <a-checkbox-group v-model:value="state.checkedList" style="width: 100%">
          <div class="list-box" ref="scrollContainer">
            <div class="list-item-wraapper" v-for="item in options" :key="item.id">
              <a-checkbox class="list-item" :value="item.id">
                <div class="library-name">{{ item.library_name }}</div>
                <div class="library-desc">{{ item.library_intro }}</div>
              </a-checkbox>
            </div>
          </div>
        </a-checkbox-group>
      </a-spin>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, SyncOutlined } from '@ant-design/icons-vue'
import { getLibraryList } from '@/api/library/index'

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
  const res = await getLibraryList({ library_name: searchKeyword.value })
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

  message.success('刷新完成')
}

const openAddLibrary = () => {
  window.open('/#/library/add')
}

defineExpose({
  open
})
</script>
