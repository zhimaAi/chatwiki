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
      <a-button @click="handleShowVerifiedConfig" :icon="createVNode(SettingOutlined)">{{ t('btn_mini_program_reply_settings') }}</a-button>
    </div>
    <div class="wechat-app-list">
      <AddWechatApp :label="t('label_bind_wechat_mini_program')" @click="showAddAlert" />
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
import { useI18n } from '@/hooks/web/useI18n'
import WechatAppItem from './wechat-app-item.vue'
import AddWechatApp from './add-wechat-app.vue'
import AddWechatAppAlert from './add-wechat-mini-program-alert.vue'
import MiniProgramReplyConfig from './mini-program-reply-config.vue'

const { t } = useI18n('views.robot.robot-config.external-service.components.wechat-mini-program')


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
    title: t('title_confirm_remove', { app_name: item.app_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('msg_remove_warning'),
    okText: secondsToGo + ' ' + t('btn_confirm'),
    okType: 'danger',
    cancelText: t('btn_cancel'),
    okButtonProps: {
      disabled: true
    },
    onOk() {
      deleteWechatApp({
        id: item.id
      }).then(() => {
        getList()

        message.success(t('msg_delete_success'))
      })
    },
    onCancel() {}
  })

  let interval = setInterval(() => {
    if (secondsToGo == 1) {
      modal.update({
        okText: t('btn_confirm'),
        okButtonProps: {
          disabled: false
        }
      })

      clearInterval(interval)
      interval = undefined
    } else {
      secondsToGo -= 1

      modal.update({
        okText: secondsToGo + ' ' + t('btn_confirm'),
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
    message.success(t('msg_refresh_success'))
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
