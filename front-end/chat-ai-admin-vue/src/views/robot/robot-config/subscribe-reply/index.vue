<template>
  <div class="auto-reply-home">
    <template v-if="hasApps">
      <ReceivedReply />
    </template>
    <template v-else>
      <div class="user-model-page">
        <!-- <div class="page-title">关注后自动回复</div> -->
        <div class="breadcrumb-wrap">
          <svg-icon @click="goBack" name="back" style="font-size: 20px;" />
          <div @click="goBack" class="breadcrumb-title">{{ t('subscribe_reply') }}</div>
          <a-switch v-model:checked="enabled_status" :checkedValue="'1'" :un-checkedValue="'0'" :checkedChildren="t('switch_on')" :un-checkedChildren="t('switch_off')" @change="handleSwitchChange" />
          <span class="switch-tip">{{ t('switch_tip') }}</span>
        </div>
        <a-alert show-icon>
          <template #message>
            <span v-html="t('alert_message', { highlight: `<span style='color: #FF4D4F;'>${t('alert_highlight')}</span>` })"></span>
          </template>
        </a-alert>
        <div class="empty-wrap">
          <ListEmpty size="180">
            <div class="empty-default">{{ t('empty_default') }}</div>
            <div class="empty-sub">{{ t('empty_sub') }}</div>
          </ListEmpty>
          <div class="empty-actions">
            <a-button type="primary" @click="goBind">{{ t('btn_bind_account') }}</a-button>
          </div>
        </div>
      </div>
    </template>
  </div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute,useRouter } from 'vue-router'
import { getWechatAppList } from '@/api/robot'
import ReceivedReply from './received-reply.vue'
import ListEmpty from '@/views/robot/robot-config/function-center/components/list-empty.vue'
import { getSpecifyAbilityConfig, saveUserAbility } from '@/api/explore/index.js'
import { message, Modal } from 'ant-design-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.subscribe-reply.index')

const route = useRoute()
const router = useRouter()
const hasApps = ref(true)
const enabled_status = ref('0')

const fetchApps = async () => {
  try {
    const res = await getWechatAppList({ app_type: 'official_account', app_name: '' })
    // 只需要account_is_verify为true的公众号
    const list = Array.isArray(res?.data) ? res.data.filter((it) => it.account_is_verify == 'true') : []
    hasApps.value = list.length > 0
  } catch (_e) {
    hasApps.value = false
  }
}

onMounted(async () => {
  try {
    const res = await getSpecifyAbilityConfig({ ability_type: 'robot_subscribe_reply' })
    const item = res?.data
    const status = String(item?.user_config?.switch_status ?? '0')
    enabled_status.value = status
  } catch (_) { enabled_status.value = '0' }
  fetchApps()
})

const goBind = () => {
  const url = router.resolve({ path: '/user/official-account' })
  window.open(url.href, '_blank')
}

const goBack = () => {
  if (route.query.id && route.query.robot_key) {
    router.push({ path: '/robot/config/function-center', query: { id: route.query.id, robot_key: route.query.robot_key } })
  } else {
    router.push({ path: '/explore/index' })
  }
}

const handleSwitchChange = (checked) => {
  const prev = enabled_status.value
  const next = checked
  if (next === '0') {
    Modal.confirm({
      title: t('modal_title'),
      content: t('modal_close_content'),
      onOk: () => {
        saveUserAbility({ ability_type: 'robot_subscribe_reply', switch_status: next }).then((res) => {
          if (res && res.res == 0) {
            enabled_status.value = next
            message.success(t('msg_operation_success'))
          } else {
            enabled_status.value = prev
          }
        }).catch(() => { enabled_status.value = prev })
      },
      onCancel: () => { enabled_status.value = '1' }
    })
    return
  }
  saveUserAbility({ ability_type: 'robot_subscribe_reply', switch_status: next }).then((res) => {
    if (res && res.res == 0) {
      enabled_status.value = next
      message.success(t('msg_operation_success'))
    } else {
      enabled_status.value = prev
    }
  }).catch(() => { enabled_status.value = prev })
}
</script>

<style lang="less" scoped>
.auto-reply-home {
  height: 100%;
  width: 100%;
  padding: 16px 48px;
  overflow-y: auto;
}
.empty-wrap {
  margin-top: 100px;
  height: calc(100% - 32px);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  .empty-default {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }
  .empty-sub {
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
}
.empty-actions {
  display: flex;
  align-items: center;
  justify-content: center;
}
.page-title {
  display: flex;
  align-items: center;
  gap: 24px;
  padding-bottom: 16px;
  background-color: #fff;
  color: #000000;
  font-size: 16px;
  font-style: normal;
  font-weight: 600;
  line-height: 24px;
}
.subManage-breadcrumb {
  display: flex;
  align-items: center;
  color: #000000;
  font-family: "PingFang SC";
  font-size: 14px;
  font-style: normal;
  line-height: 22px;
  padding-bottom: 16px;
}

.breadcrumb-wrap {
  width: fit-content;
  display: flex;
  align-items: center;
  cursor: pointer;
  margin-bottom: 16px;
}
.breadcrumb-title {
  margin: 0 12px 0 2px;
  color: #262626;
  font-size: 16px;
  font-style: normal;
  font-weight: 600;
  line-height: 24px;
}
.switch-tip {
  margin-left: 4px;
  color: #8c8c8c;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}
</style>
