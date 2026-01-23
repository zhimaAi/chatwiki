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
  margin-bottom: 16px;
}
</style>

<template>
  <div class="">
    <div class="add-btn-block">
      <a-button @click="handleShowVerifiedConfig" :icon="createVNode(SettingOutlined)">小程序回复设置</a-button>
    </div>
    <div class="wechat-app-list">
      <AddWechatApp label="绑定微信小程序" @click="showAddAlert" />
      <WechatAppItem
        v-for="item in list"
        :key="item.id"
        app_type="mini_program"
        :item="item"
        @edit="handleEdit"
        @delete="handleDelete"
        @refresh="handleRefresh"
      />
    </div>
    <AddWechatAppAlert ref="addAppAlertRef" @ok="onSaveSuccess" />
     <MiniProgramReplyConfig ref="replyConfigRef" @change="onSaveReplyConfigSuccess" />
  </div>
</template>

<script setup>
import { ref, inject, onMounted, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined, SettingOutlined } from '@ant-design/icons-vue'
import { getWechatAppList, deleteWechatApp, refreshAccountVerify } from '@/api/robot'
import WechatAppItem from './wechat-app-item.vue'
import AddWechatApp from './add-wechat-app.vue'
import AddWechatAppAlert from './add-wechat-mini-program-alert.vue'
import MiniProgramReplyConfig from './mini-program-reply-config.vue'


const { robotInfo, getRobot } = inject('robotInfo')

const addAppAlertRef = ref()
const list = ref([])

const getList = () => {
  getWechatAppList({
    robot_id: robotInfo.id,
    app_type: 'mini_program',
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

const handleRefresh = (item) => {
  refreshAccountVerify({
    id: item.id
  }).then(() => {
    message.success('刷新成功')
    getList()
  })
}

const onSaveSuccess = () => {
  getList()
}

const replyConfigRef = ref(null)

const handleShowVerifiedConfig = () => {
  let config = {
    robotId: robotInfo.id,
    aiGenerated: robotInfo.show_ai_msg_mini,
    typingIndicator: robotInfo.show_typing_mini
  }

  replyConfigRef.value.open(config)
}

const onSaveReplyConfigSuccess = () => {
  getRobot(robotInfo.id)
}

onMounted(() => {
  getList()
})
</script>
