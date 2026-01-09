<style lang="less">
.correlation-library-modal.ant-modal-wrap .ant-modal-content {
  padding: 0;
}
</style>
<style lang="less" scoped>
.library-checkbox-box {
  .list-tools {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 16px;
    margin-bottom: 8px;
  }
  .list-group-box {
    display: flex;
    height: 608px;
    overflow: hidden;
  }

  .group-list-box {
    width: 248px;
    border-right: 1px solid #d9d9d9;
    position: relative;
    padding: 24px 0 0 16px;
    .main-title-block {
      color: #262626;
      font-size: 16px;
      font-weight: 600;
      line-height: 24px;
      padding-left: 8px;
      padding-bottom: 24px;
    }

    .classify-box {
      flex: 1;
      overflow: hidden;
      font-size: 14px;
      .classify-item {
        height: 40px;
        padding: 0 16px;
        display: flex;
        align-items: center;
        justify-content: space-between;
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
        &:hover {
          background: #f2f4f7;
        }
        &.active {
          color: #2475fc;
          background: #e6efff;
        }
      }
    }
  }
  .list-item-box {
    flex: 1;
    height: 100%;
    overflow: hidden;
    padding-top: 24px;
    .alert-box {
      width: 600px;
      padding-left: 16px;
    }
    .btn-box {
      display: flex;
      align-items: center;
      gap: 8px;
      margin-top: 16px;
      margin-bottom: 8px;
      padding-left: 16px;
    }
  }
  .library-list-item {
    min-height: 97px;
    padding: 24px 12px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;
    &:hover {
      background: var(--09, #f2f4f7);
    }
    .avatar-box {
      width: 48px;
      height: 48px;
      img {
        width: 100%;
        height: 100%;
        border-radius: 12px;
      }
    }
    .info-content-box {
      flex: 1;
      .title-info-block {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #262626;
        font-weight: 600;
        line-height: 22px;
        .type-tag {
          height: 18px;
          line-height: 16px;
          padding: 0 4px;
          font-size: 12px;
          font-weight: 400;
          border-radius: 6px;
          color: #2475fc;
          border: 1px solid #99bffd;
          &.gray-tag {
            border: 1px solid var(--06, #d9d9d9);
            color: #8c8c8c;
          }
        }
      }
      .desc-info-block {
        color: #8c8c8c;
        font-size: 14px;
        line-height: 22px;
        font-weight: 400;
        margin-top: 4px;
      }
    }
  }
  .footer-box {
    height: 52px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 8px;
    border-top: 1px solid var(--06, #d9d9d9);
    padding: 0 24px;
  }
}
</style>

<template>
  <a-modal
    :width="916"
    v-model:open="show"
    wrapClassName="correlation-library-modal"
    :title="null"
    :footer="null"
  >
    <div class="library-checkbox-box">
      <div class="list-group-box">
        <div class="group-list-box">
          <div class="main-title-block">关联知识库</div>
          <div style="margin-bottom: 16px">
            <a-input
              style="width: 216px"
              v-model:value="searchKeyword"
              placeholder="请输入知识库名称搜索"
              @change="onSearch"
            >
              <template #suffix>
                <SearchOutlined style="color: rgba(0, 0, 0, 0.25)" />
              </template>
            </a-input>
          </div>
          <cu-scroll style="padding-right: 16px; height: 480px">
            <div class="classify-box">
              <div
                class="classify-item"
                @click="handleChangeGroup(item)"
                :class="{ active: item.id == group_id }"
                v-for="item in groupLists"
                :key="item.id"
              >
                <div class="classify-title">{{ item.group_name }}</div>
                <div class="right-content">
                  <div class="num" :class="{ 'num-block': item.id <= 0 }">{{ item.total }}</div>
                </div>
              </div>
            </div>
          </cu-scroll>
        </div>
        <div class="list-item-box">
          <div class="alert-box">
            <a-alert
              class="zm-alert-info"
              message="请选择关联知识库，机器人会根据知识库内上传的文档回复用户的提问。每个机器人最多关联5个知识库"
              type="info"
            />
          </div>
          <div class="btn-box">
            <a-button @click="onRefresh"> <SyncOutlined /> 刷新 </a-button>
            <a-button type="primary" ghost @click="openAddLibrary">新建知识库</a-button>
          </div>
          <a-spin :spinning="isRefresh" :delay="100">
            <cu-scroll style="padding: 0 16px; height: 414px">
              <div
                class="library-list-item"
                @click="handleChangeChecked(item)"
                v-for="item in showOptions"
                :key="item.id"
              >
                <div class="avatar-box">
                  <img :src="item.avatar" alt="" />
                </div>
                <div class="info-content-box">
                  <div class="title-info-block">
                    {{ item.library_name }}
                    <span class="type-tag" :class="{ 'gray-tag': item.graph_switch == 0 }"
                      >Graph</span
                    >
                    <span class="type-tag" v-if="item.type == 0">普通知识库</span>
                    <span class="type-tag" v-if="item.type == 1">对外知识库</span>
                    <span class="type-tag" v-if="item.type == 2">问答知识库</span>
                    <span class="type-tag" v-if="item.type == 3">公众号知识库</span>
                  </div>
                  <div class="desc-info-block">{{ item.library_intro }}</div>
                </div>
                <div class="check-block">
                  <a-checkbox
                    :disabled="robotInfo.default_library_id == item.id"
                    :checked="state.checkedList.includes(item.id)"
                  ></a-checkbox>
                </div>
              </div>
            </cu-scroll>
          </a-spin>
          <div class="footer-box">
            <a-button @click="show = false">取消</a-button>
            <a-button type="primary" @click="saveCheckedList">
              ({{ state.checkedList.length }}) 确定</a-button
            >
          </div>
        </div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, SyncOutlined, RightOutlined, LeftOutlined } from '@ant-design/icons-vue'
import { getLibraryList, getLibraryListGroup } from '@/api/library/index'
import { useRobotStore } from '@/stores/modules/robot'
const robotStore = useRobotStore()
const { robotInfo } = robotStore

const emit = defineEmits(['change'])
const props = defineProps({
  showWxType: {
    type: Boolean,
    default: false
  }
})

const state = reactive({
  indeterminate: false,
  checkAll: false,
  checkedList: []
})

const group_id = ref('')

const groupLists = ref([])

const getGroupList = () => {
  getLibraryListGroup().then((res) => {
    let lists = res.data || []
    let totalNumber = 0
    // 计算每个分组的机器人数量
    lists.forEach((group) => {
      totalNumber += +group.total
    })
    groupLists.value = [
      {
        group_name: '全部',
        total: totalNumber,
        id: ''
      },
      ...lists
    ]
  })
}

const handleChangeGroup = (item) => {
  group_id.value = item.id
  getList()
}

const show = ref(false)

const open = (checkedList) => {
  group_id.value = ''
  getList()
  getGroupList()

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

const showOptions = computed(() => {
  if (props.showWxType) {
    return options.value
  } else {
    return options.value.filter((item) => item.type != 3)
  }
})

const triggerChange = () => {
  emit('change', [...state.checkedList])
}

const getList = async () => {
  isRefresh.value = true
  const res = await getLibraryList({
    library_name: searchKeyword.value,
    type: '',
    show_open_docs: 1
  })
  isRefresh.value = false
  if (res) {
    let list = res.data || []
    if (group_id.value != '') {
      list = list.filter((item) => item.group_id == group_id.value)
    }
    options.value = list
  }
}

const handleChangeChecked = (e) => {
  let value = e.id
  if (robotInfo.default_library_id == value) {
    return
  }
  if (state.checkedList.includes(value)) {
    state.checkedList = state.checkedList.filter((item) => item != value)
  } else {
    state.checkedList.push(value)
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
  window.open('/#/library/list')
}

defineExpose({
  open
})
</script>
