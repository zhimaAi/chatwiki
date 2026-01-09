<style lang="less" scoped>
._container {
  padding: 24px;
  height: 100%;
  width: 100%;
  overflow-y: auto;

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
      <span>{{ t('title') }}</span>
      <a-button @click="showAddModal" :icon="createVNode(PlusOutlined)" type="primary">{{ t('bind_btn') }}</a-button>
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
      <div class="title">{{ t('empty_title') }}</div>
      <a-button @click="showAddModal" class="btn" type="primary">{{ t('bind_now_btn') }}</a-button>
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
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.user.official-account.index')

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
    title: t('confirm_remove', { app_name: item.app_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('remove_warning'),
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    cancelText: '取消',
    okButtonProps: {
      disabled: true
    },
    onOk() {
      deleteWechatApp({id: item.id}).then(() => {
        getList()
        message.success(t('delete_success'))
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
    message.success(t('refresh_success'))
    getList()
  })
}

const onDragEnd = () => {
  let filter_sort = []
  list.value.forEach((item, index) => {
    filter_sort.push(item.id)
  })
  sortWechatApp({id_list: filter_sort.toString()}).then((res) => {
    message.success(t('sort_success'))
  })
}
</script>
