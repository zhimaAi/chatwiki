<style lang="less" scoped>
.wechat-app-list {
  display: flex;
  gap: 24px;
  flex-flow: row wrap;
  .add-wechat-app{
    min-height: 112px;
  }
}
.tip-alert {
  margin-bottom: 24px;
}
.messenger-embedded {
  .messenger-toolbar {
    margin-bottom: 24px;
  }

  .wechat-app-list {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    gap: 16px;
  }

  .tip-alert {
    padding: 16px;
    border: 0;
    border-radius: 8px;
    background: var(--app-selected);
  }
}
@media (max-width: 1200px) {
  .messenger-embedded .wechat-app-list {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
@media (max-width: 760px) {
  .messenger-embedded .wechat-app-list {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>

<template>
  <div :class="{ 'messenger-embedded': embedded }">
    <a-alert class="tip-alert zm-alert-info" type="info" show-icon>
      <template #message>{{ customerT('messenger_intro') }}</template>
    </a-alert>
    <div v-if="embedded" class="messenger-toolbar">
      <a-button type="primary" @click="showAddAlert"><PlusOutlined />{{ customerT('create_channel') }}</a-button>
    </div>
    <div class="wechat-app-list">
      <AddWechatApp v-if="!embedded" label="创建渠道" @click="showAddAlert" />
      <WechatAppItem
        v-for="item in list"
        :key="item.id"
        :item="item"
        :compact="embedded"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </div>
    <AddMessagentAlert ref="addAppAlertRef" :compact="embedded" @ok="onSaveSuccess" />
  </div>
</template>

<script setup>
import { ref, inject, onMounted, createVNode } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue'
import { getWechatAppList, deleteWechatApp } from '@/api/robot'
import { useI18n } from '@/hooks/web/useI18n'
import WechatAppItem from './wechat-app-item.vue'
import AddWechatApp from './add-wechat-app.vue'
import AddMessagentAlert from "./add-messenget-alert.vue";

const { t } = useI18n('views.robot.robot-config.external-service.components.dingding-robot')
const { t: customerT } = useI18n('views.customer-service.index')
defineProps({ embedded: Boolean })
const { robotInfo } = inject('robotInfo')

const addAppAlertRef = ref()
const list = ref([])

const getList = () => {
  getWechatAppList({
    robot_id: robotInfo.id,
    app_type: 'messenger',
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
