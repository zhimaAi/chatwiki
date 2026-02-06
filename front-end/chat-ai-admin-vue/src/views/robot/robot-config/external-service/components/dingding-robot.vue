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
    <a-alert class="tip-alert" type="info" show-icon>
      <template #message>{{ t('tip_bind_dingding') }}</template>
    </a-alert>
    <div class="wechat-app-list">
      <AddWechatApp :label="t('label_bind_dingding')" @click="showAddAlert" />
      <WechatAppItem
        v-for="item in list"
        :key="item.id"
        :item="item"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </div>
    <AddDingdingAlert ref="addAppAlertRef" @ok="onSaveSuccess" />
  </div>
</template>

<script setup>
import { ref, inject, onMounted, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getWechatAppList, deleteWechatApp } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'
import WechatAppItem from './wechat-app-item.vue'
import AddWechatApp from './add-wechat-app.vue'
import AddDingdingAlert from "./add-dingding-alert.vue";

const { t } = useI18n('views.robot.robot-config.external-service.components.dingding-robot')
const { robotInfo } = inject('robotInfo')

const addAppAlertRef = ref()
const list = ref([])

const getList = () => {
  getWechatAppList({
    robot_id: robotInfo.id,
    app_type: 'dingtalk_robot',
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

const onSaveSuccess = () => {
  getList()
}

onMounted(() => {
  getList()
})
</script>
