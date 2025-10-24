<style lang="less" scoped>
.wechat-app-list {
  display: flex;
  gap: 24px;
  flex-flow: row wrap;
}
.tip-alert {
  margin-bottom: 24px;
}
</style>

<template>
  <div class="">
    <a-alert
      class="tip-alert"
      message="绑定微信客服，用户通过微信客服咨询时由机器人自动回复"
      type="info"
      show-icon
    />
    <div class="wechat-app-list">
      <AddWechatApp label="绑定微信客服" @click="showAddAlert" />
      <WechatAppItem
        v-for="item in list"
        :key="item.id"
        :item="item"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </div>
    <AddWechatAppAlert ref="addAppAlertRef" @ok="onSaveSuccess" />
  </div>
</template>

<script setup>
import { ref, inject, onMounted, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getWechatAppList, deleteWechatApp } from '@/api/robot'
import WechatAppItem from './wechat-app-item.vue'
import AddWechatApp from './add-wechat-app.vue'
import AddWechatAppAlert from './add-wechat-customer-service-alert.vue'

const { robotInfo } = inject('robotInfo')

const addAppAlertRef = ref()
const list = ref([])

const getList = () => {
  getWechatAppList({
    robot_id: robotInfo.id,
    app_type: 'wechat_kefu',
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
        okText: '确 定',
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

const onSaveSuccess = () => {
  getList()
}

onMounted(() => {
  getList()
})
</script>
