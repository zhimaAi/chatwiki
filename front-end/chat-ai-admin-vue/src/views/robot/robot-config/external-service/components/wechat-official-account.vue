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
      <template #message> {{ t('support_unverified_account') }}，<a @click="handleShowDemoModal">{{ t('view_demo') }}</a> </template>
    </a-alert>
    <LoadingBox v-if="loading"/>
    <template v-else-if="list.length">
      <div class="add-btn-block">
        <a-button @click="showAddAlert" type="primary">{{ t('associate_official_account') }}</a-button>
        <a-button @click="handleAddUnverified" :icon="createVNode(SettingOutlined)">{{ t('unverified_account_reply_settings') }}</a-button>
        <a-button @click="handleShowVerifiedConfig" :icon="createVNode(SettingOutlined)">{{ t('verified_account_reply_settings') }}</a-button>
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
    <EmptyBox v-else style="margin-top: 20vh;" :title="t('no_official_account_associated')">
      <template #desc>
        <a-button @click="showAddAlert" type="primary">{{ t('associate_official_account') }}</a-button>
      </template>
    </EmptyBox>

    <SelectWechatApp ref="selectAppRef" @change="getList" />
    <AddUnverifiedAlert ref="addUnverifiedAlertRef" />
    <DemoPreviewModal ref="demoPreviewModalRef" />
    <OfficialAccountReplyConfig ref="replyConfigRef" @change="onSaveReplyConfigSuccess" />
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
import OfficialAccountReplyConfig from './official-account-reply-config.vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.external-service.components.wechat-official-account')

const { robotInfo, getRobot } = inject('robotInfo')

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
    message.success(t('refresh_success'))
    getList()
  })
}

const handleDelete = (item) => {
  const modal = Modal.confirm({
    title: t('confirm_remove', { app_name: item.app_name }),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('remove_warning'),
    okText: t('confirm'),
    cancelText: t('cancel'),
    onOk() {
      let row = list.value.filter(i => i.app_id != item.app_id)
      let ids = row.map(i => i.app_id).toString()
      robotBindWxApp({
        robot_id: robotInfo.id,
        app_id_list: ids
      }).then(() => {
        getList()
        message.success(t('delete_success'))
      })
    },
  })
}

const addUnverifiedAlertRef = ref(null)

const handleAddUnverified = () => {
  addUnverifiedAlertRef.value.open()
}

const replyConfigRef = ref(null)
const handleShowVerifiedConfig = () => {
  let config = {
    robotId: robotInfo.id,
    aiGenerated: robotInfo.show_ai_msg_gzh,
    typingIndicator: robotInfo.show_typing_gzh
  }

  replyConfigRef.value.open(config)
}

const onSaveReplyConfigSuccess = () => {
  getRobot(robotInfo.id)
}

const demoPreviewModalRef = ref(null)
const handleShowDemoModal = () => {
  demoPreviewModalRef.value.show()
}
</script>
