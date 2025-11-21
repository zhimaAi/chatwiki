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
    <div class="add-btn-block">
      <a-button @click="showAddAlert" :icon="createVNode(PlusOutlined)" type="primary"
        >绑定公众号</a-button
      >
      <a-button @click="handleAddUnverified" :icon="createVNode(SettingOutlined)"
        >未认证公众号回复设置</a-button
      >
    </div>
    <div class="wechat-app-list">
      <!-- <AddWechatApp label="绑定公众号" @click="showAddAlert" /> -->
      <WechatAppItem
        v-for="item in list"
        :key="item.id"
        :item="item"
        app_type="official_account"
        @edit="handleEdit"
        @delete="handleDelete"
        @refresh="handleRefresh"
      />
    </div>
    <AddWechatOfficialAccountAlert ref="addAppAlertRef" @ok="onSaveSuccess" />
    <AddUnverifiedAlert ref="addUnverifiedAlertRef" />
    <DemoPreviewModal ref="demoPreviewModalRef" />
  </div>
</template>

<script setup>
import { ref, inject, onMounted, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined, PlusOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { getWechatAppList, deleteWechatApp, refreshAccountVerify } from '@/api/robot'
import WechatAppItem from './wechat-app-item.vue'
import AddWechatApp from './add-wechat-app.vue'
import AddWechatOfficialAccountAlert from './add-wechat-official-account-alert.vue'
import AddUnverifiedAlert from './add-unverified-alert.vue'
import DemoPreviewModal from './demo-preview-modal.vue'

const { robotInfo } = inject('robotInfo')

const addAppAlertRef = ref()
const list = ref([])

const getList = () => {
  getWechatAppList({
    robot_id: robotInfo.id,
    app_type: 'official_account',
    app_name: ''
  }).then((res) => {
    list.value = res.data
  })
}

const showAddAlert = () => {
  addAppAlertRef.value.open()
}

const handleDelete = (item) => {
  let secondsToGo = 3

  const modal = Modal.confirm({
    title: `确定移除${item.app_name}吗?`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '移除后，机器人无法继续回复用户消息。',
    okText: secondsToGo + ' 确 定',
    okType: 'danger',
    cancelText: '取消',
    okButtonProps: {
      disabled: true
    },
    onOk() {
      deleteWechatApp({
        id: item.id
      }).then(() => {
        getList()

        message.success('删除成功')
      })
    },
    onCancel() {}
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
  item.wechat_ip = robotInfo.wechat_ip
  addAppAlertRef.value.open({ ...item })
}

const handleRefresh = (item) => {
  refreshAccountVerify({
    id: item.id
  }).then(() => {
    message.success('刷新成功')
    getList()
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

const onSaveSuccess = () => {
  getList()
}

onMounted(() => {
  getList()
})
</script>
