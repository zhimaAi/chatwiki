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
  .list-group-box {
    display: flex;
    gap: 16px;
    height: 388px;
    overflow: hidden;
  }

  .group-list-box {
    width: 256px;
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    margin-top: 8px;
    padding: 16px;
    padding-right: 0;
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
      margin-bottom: 8px;
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
  <a-modal width="946px" v-model:open="show" title="关联知识库" @ok="saveCheckedList">
    <div class="library-checkbox-box">
      <a-alert
        message="请选择关联知识库，机器人会根据知识库内上传的文档回复用户的提问。每个机器人最多关联5个知识库"
        type="info"
      />
      <div class="list-tools">
        <div>
          <a-input
            style="width: 282px"
            v-model:value="searchKeyword"
            placeholder="请输入知识库名称搜索"
            @change="onSearch"
          >
            <template #suffix>
              <SearchOutlined style="color: rgba(0, 0, 0, 0.25)" />
            </template>
          </a-input>
        </div>
        <div>
          <a-button style="margin-right: 8px" @click="onRefresh"> <SyncOutlined /> 刷新 </a-button>
          <a-button type="primary" ghost @click="openAddLibrary">新建知识库</a-button>
        </div>
      </div>
      <div class="list-group-box">
        <div class="group-list-box" :class="{ 'hide-group': isHideGroup }">
          <cu-scroll style="padding-right: 16px;">
            <div class="group-head-box" v-if="!isHideGroup">
              <div class="head-title">
                <div>知识库分组</div>
              </div>
            </div>
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
          <a-tooltip placement="right" :title="isHideGroup ? '展开分组' : '收起分组'">
            <div class="hide-group-box" @click="handleChangeHideGroup">
              <LeftOutlined v-if="!isHideGroup" />
              <RightOutlined v-else />
            </div>
          </a-tooltip>
        </div>
        <div style="flex: 1; overflow: hidden">
          <a-spin :spinning="isRefresh" :delay="100">
            <a-checkbox-group :value="state.checkedList" style="width: 100%">
              <div class="list-box" ref="scrollContainer">
                <div class="list-item-wraapper" v-for="item in showOptions" :key="item.id">
                  <a-checkbox
                    class="list-item"
                    :value="item.id"
                    @change="handleChangeChecked"
                    :disabled="robotInfo.default_library_id == item.id"
                  >
                    <div class="library-name">{{ item.library_name }}</div>
                    <div class="library-desc">{{ item.library_intro }}</div>
                  </a-checkbox>
                </div>
              </div>
            </a-checkbox-group>
          </a-spin>
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
let hideGroupLocalKey = 'library-list-hide-group-key'
const isHideGroup = ref(localStorage.getItem(hideGroupLocalKey) == 1)

const handleChangeHideGroup = () => {
  isHideGroup.value = !isHideGroup.value
  localStorage.setItem(hideGroupLocalKey, isHideGroup.value ? 1 : 0)
}
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
    return options.value.filter(item => item.type != 3)
  }
})

const triggerChange = () => {
  emit('change', [...state.checkedList])
}

const getList = async () => {
  const res = await getLibraryList({
    library_name: searchKeyword.value,
    type: '',
    show_open_docs: 1
  })
  if (res) {
    let list = res.data || []
    if(group_id.value != ''){
      list = list.filter(item => item.group_id == group_id.value)
    }
    options.value = list
  }
}

const handleChangeChecked = (e) => {
  let value = e.target.value
  if(state.checkedList.includes(value)){
    state.checkedList = state.checkedList.filter(item => item != value)
  }else{
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
