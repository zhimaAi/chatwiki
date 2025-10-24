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
  <a-modal width="746px" v-model:open="show" title="添加技能" @ok="saveCheckedList">
    <div class="library-checkbox-box">
      <a-spin :spinning="isRefresh" :delay="100">
        <a-checkbox-group v-model:value="state.checkedList" style="width: 100%">
          <div class="list-box" ref="scrollContainer">
            <div class="list-item-wraapper" v-for="item in options" :key="item.id">
              <a-checkbox class="list-item"  :value="item.id">
                <div class="library-name">{{ item.robot_name }}
                  <span class="library-desc" v-if="item.start_node_key===''"> (未发布)</span></div>
                <div class="library-desc">{{ item.robot_intro || '--' }}</div>
              </a-checkbox>
            </div>
          </div>
        </a-checkbox-group>
      </a-spin>
      <div class="empty-box" v-if="!isRefresh && !options.length">
        <img src="@/assets/img/library/preview/empty.png" alt="" />
        <div>暂无数据，请先添加工作流</div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { message } from 'ant-design-vue'
import { SearchOutlined, SyncOutlined } from '@ant-design/icons-vue'
import { getRobotList } from '@/api/robot/index'

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
  const res = await getRobotList({application_type: 1})
  if (res) {
    let list = res.data || []
    options.value = list
  }
}

const isRefresh = ref(false)

defineExpose({
  open
})
</script>
