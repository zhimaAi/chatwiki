<style lang="less" scoped>
.wechat-app-list {
  display: flex;
  gap: 24px;
  flex-flow: row wrap;
}
.add-btn-block {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 16px 0;
}
</style>

<template>
  <div>
    <a-alert banner>
      <template #message> 支持接入未认证公众号，<a @click="handleShowDemoModal">查看效果示例</a> </template>
    </a-alert>
    <LoadingBox v-if="loading"/>
    <template v-else-if="list.length">
      <div class="add-btn-block">
        <a-button @click="showAddAlert" type="primary">关联公众号</a-button>
        <a-button @click="handleAddUnverified" :icon="createVNode(SettingOutlined)">未认证公众号回复设置</a-button>
      </div>
      <div class="wechat-app-list">
        <WechatAppItem
          v-for="item in list"
          :key="item.id"
          :item="item"
          :showMenu="['del']"
          app_type="official_account"
          @delete="handleDelete"
          @refresh="handleRefresh"
        />
      </div>
    </template>
    <EmptyBox v-else style="margin-top: 20vh;" title="暂未关联公众号">
      <template #desc>
        <a-button @click="showAddAlert" type="primary">关联公众号</a-button>
      </template>
    </EmptyBox>

    <SelectWechatApp ref="selectAppRef" @change="getList" />
    <AddUnverifiedAlert ref="addUnverifiedAlertRef" />
    <DemoPreviewModal ref="demoPreviewModalRef" />
  </div>
</template>

<script setup>
import { ref, inject, onMounted, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined, PlusOutlined, SettingOutlined } from '@ant-design/icons-vue'
import {getWechatAppList, deleteWechatApp, refreshAccountVerify, robotBindWxApp} from '@/api/robot'
import WechatAppItem from './wechat-app-item.vue'
import AddUnverifiedAlert from './add-unverified-alert.vue'
import DemoPreviewModal from './demo-preview-modal.vue'
import SelectWechatApp from "@/views/robot/robot-config/external-service/components/select-wechat-app.vue";
import LoadingBox from "@/components/common/loading-box.vue";
import EmptyBox from "@/components/common/empty-box.vue";

const { robotInfo } = inject('robotInfo')

const selectAppRef = ref()
const loading = ref(true)
const list = ref([])

onMounted(() => {
  getList()
})

const getList = () => {
  getWechatAppList({
    robot_id: robotInfo.id,
    app_type: 'official_account',
    app_name: ''
  }).then((res) => {
    list.value = res.data
  }).finally(() => {
    loading.value = false
  })
}

const showAddAlert = () => {
  selectAppRef.value.open(robotInfo, list.value.map(i => i.app_id))
}

const handleRefresh = (item) => {
  refreshAccountVerify({
    id: item.id
  }).then(() => {
    message.success('刷新成功')
    getList()
  })
}

const handleDelete = (item) => {
  const modal = Modal.confirm({
    title: `确定移除${item.app_name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '移除后，机器人无法继续回复用户消息。',
    okText: '确定',
    cancelText: '取消',
    onOk() {
      let row = list.value.filter(i => i.app_id != item.app_id)
      let ids = row.map(i => i.app_id).toString()
      robotBindWxApp({
        robot_id: robotInfo.id,
        app_id_list: ids
      }).then(() => {
        getList()
        message.success('删除成功')
      })
    },
  })
}

const addUnverifiedAlertRef = ref(null)

const handleAddUnverified = () => {
  addUnverifiedAlertRef.value.open()
}

const demoPreviewModalRef = ref(null)
const handleShowDemoModal = () => {
  demoPreviewModalRef.value.show()
}
</script>
