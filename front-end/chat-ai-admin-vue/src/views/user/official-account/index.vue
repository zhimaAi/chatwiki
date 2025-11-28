<style lang="less" scoped>
._container {
  padding: 24px;
  height: 100%;
  width: 100%;

  .page-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #000000;
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 24px;
  }
}

.empty-box {
  display: flex;
  height: 100%;
  width: 100%;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  img {
    width: 200px;
    height: 200px;
  }

  .title {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }

  .btn {
    margin-top: 24px;
  }
}

.loading-box {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
}

.wechat-app-list {
  display: flex;
  gap: 24px;
  flex-flow: row wrap;
  margin-top: 16px;
}

.add-btn-block {
  display: flex;
  align-items: center;
  gap: 16px;
}

.drag-item {
  cursor: move;
}
</style>

<template>
  <div class="_container">
    <div class="page-title">
      <span>公众号管理</span>
      <a-button @click="showAddModal" :icon="createVNode(PlusOutlined)" type="primary">绑定公众号</a-button>
    </div>
    <div v-if="loading" class="loading-box">
      <a-spin/>
    </div>
    <template v-else-if="list.length">
      <draggable
        class="wechat-app-list"
        v-model="list"
        item-key="id"
        @end="onDragEnd"
        handle=".drag-item"
      >
        <template #item="{ element, index }">
          <WechatAppItem
            class="drag-item"
            :key="element.id"
            :item="element"
            app_type="official_account"
            @edit="handleEdit"
            @delete="handleDelete"
            @refresh="handleRefresh"
          />
        </template>
      </draggable>
    </template>
    <div v-else class="empty-box">
      <img src="@/assets/empty.png"/>
      <div class="title">暂未绑定公众号</div>
      <a-button @click="showAddModal" class="btn" type="primary">立即绑定</a-button>
    </div>

    <StoreModal ref="storeRef" @ok="getList"/>
  </div>
</template>

<script setup>
import {ref, onMounted, createVNode} from 'vue'
import draggable from 'vuedraggable'
import {message, Modal} from 'ant-design-vue'
import {ExclamationCircleOutlined, PlusOutlined} from '@ant-design/icons-vue'
import {getWechatAppList, deleteWechatApp, refreshAccountVerify, sortWechatApp} from '@/api/robot'
import WechatAppItem from './components/wechat-app-item.vue'
import StoreModal from './components/store-modal.vue'

const storeRef = ref()
const list = ref([])
const loading = ref(true)

onMounted(() => {
  getList()
})

const getList = () => {
  getWechatAppList({
    app_type: 'official_account',
    app_name: ''
  }).then((res) => {
    list.value = res.data
  }).finally(() => {
    loading.value = false
  })
}

const showAddModal = () => {
  storeRef.value.open()
}

const handleDelete = (item) => {
  let secondsToGo = 3

  const modal = Modal.confirm({
    title: `确定移除${item.app_name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '移除后，机器人已关联绑定或其他已配置该公众号的功能都会失效',
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    cancelText: '取消',
    okButtonProps: {
      disabled: true
    },
    onOk() {
      deleteWechatApp({id: item.id}).then(() => {
        getList()
        message.success('删除成功')
      })
    },
    onCancel() {
    }
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: '确定',
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' 确 定',
        okButtonProps: {
          disabled: true
        }
      })
    }
  }, 1000)
}

const handleEdit = (item) => {
  storeRef.value.open({...item})
}

const handleRefresh = (item) => {
  refreshAccountVerify({
    id: item.id
  }).then(() => {
    message.success('刷新成功')
    getList()
  })
}

const onDragEnd = () => {
  let filter_sort = []
  list.value.forEach((item, index) => {
    filter_sort.push(item.id)
  })
  sortWechatApp({id_list: filter_sort.toString()}).then((res) => {
    message.success('排序成功')
  })
}
</script>
