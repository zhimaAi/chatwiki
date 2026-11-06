<template>
  <div class="other-channels-page">
    <template v-if="!activeChannel">
      <CustomerPageHeader :title="t('other_channels')" />
      <div class="channels-body">
        <div class="channel-filters" role="group" :aria-label="t('channel_category')">
          <button
            v-for="category in categories"
            :key="category.key"
            type="button"
            :class="{ selected: selectedCategory === category.key }"
            :aria-pressed="selectedCategory === category.key"
            @click="selectedCategory = category.key"
          >
            {{ t(category.label) }}
          </button>
        </div>
        <div class="channels-grid">
          <button
            v-for="channel in filteredChannels"
            :key="channel.key"
            type="button"
            class="channel-card"
            @click="openChannel(channel)"
          >
            <span class="channel-icon" :class="{ light: channel.lightIcon }" :style="{ background: channel.background }">
              <img :src="channel.icon" alt="" :class="{ large: channel.largeIcon, small: channel.smallIcon }" />
            </span>
            <span class="channel-content">
              <span class="channel-heading">
                <span class="channel-name">{{ t(channel.nameKey) }}</span>
                <span
                  v-if="badgeOf(channel)"
                  class="channel-badge"
                  :class="{ green: badgeOf(channel).active }"
                >
                  <CheckOutlined v-if="channel.key === 'api_key' && badgeOf(channel).active" />
                  {{ badgeOf(channel).text }}
                </span>
              </span>
              <span class="channel-desc">{{ t(channel.descKey) }}</span>
            </span>
          </button>
        </div>
      </div>
    </template>
    <template v-else>
      <CustomerPageHeader class="secondary-header" :title="t(activeTitleKey)" back :back-label="t('back')" @back="backToChannels(); loadCounts()" />
      <div class="secondary-body">
        <div v-if="secondaryLoading" class="secondary-state">
          <a-spin size="large" />
        </div>
        <a-empty v-else-if="secondaryError" :description="t(secondaryError)">
          <a-button type="primary" @click="openChannel(activeChannel)">{{ t('retry') }}</a-button>
        </a-empty>
        <ApiKeyManage v-else-if="activeChannel.key === 'api_key'" embedded />
        <component
          v-else-if="activeComponent && secondaryReady"
          :is="activeComponent"
          :key="activeChannel.key"
          v-bind="activeChannel.key === 'messenger' ? { embedded: true } : {}"
        />
        <div v-else class="secondary-state"><a-spin size="large" /></div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, provide, reactive, ref, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { CheckOutlined } from '@ant-design/icons-vue'
import CustomerPageHeader from '../components/page-header.vue'
import { useI18n } from '@/hooks/web/useI18n'
import { useCustomerServiceStore } from '@/stores/modules/customer-service'
import { useRobotStore } from '@/stores/modules/robot'
import { getWechatAppList, listRobotApikey } from '@/api/robot/index'
import ApiKeyManage from '@/views/robot/api-key-manage/index.vue'
import WeChatOfficialAccount from '@/views/robot/robot-config/external-service/components/wechat-official-account.vue'
import WeChatMiniProgram from '@/views/robot/robot-config/external-service/components/wechat-mini-program.vue'
import WeChatCustomerService from '@/views/robot/robot-config/external-service/components/wechat-customer-service.vue'
import WeComRobot from '@/views/robot/robot-config/external-service/components/wecom-robot.vue'
import FeishuRobot from '@/views/robot/robot-config/external-service/components/feishu-robot.vue'
import DingTalkRobot from '@/views/robot/robot-config/external-service/components/dingding-robot.vue'
import MessengerRobot from '@/views/robot/robot-config/external-service/components/Messenget.vue'
import TelegramRobot from '@/views/robot/robot-config/external-service/components/telegram-robot.vue'
import LineRobot from '@/views/robot/robot-config/external-service/components/Line.vue'
import WhatsAppRobot from '@/views/robot/robot-config/external-service/components/whats-app.vue'
import wechatIcon from '@/assets/img/channels/wechat.svg'
import miniProgramIcon from '@/assets/img/channels/mini-program.svg'
import wechatCustomerServiceIcon from '@/assets/img/channels/wechat-customer-service.png'
import wecomIcon from '@/assets/img/channels/wecom.svg'
import feishuIcon from '@/assets/img/channels/feishu.png'
import dingtalkIcon from '@/assets/img/channels/dingtalk.svg'
import messengerIcon from '@/assets/img/channels/messenger.svg'
import telegramIcon from '@/assets/img/channels/telegram.svg'
import lineIcon from '@/assets/img/channels/line.svg'
import whatsappIcon from '@/assets/img/channels/whatsapp.svg'
import apiKeyIcon from '@/assets/img/channels/api-key.svg'

const { t } = useI18n('views.customer-service.index')
const customerServiceStore = useCustomerServiceStore()
const robotStore = useRobotStore()
const { robotInfo: robotDetail } = storeToRefs(robotStore)
const { getRobot } = robotStore
const customerRobot = computed(() => customerServiceStore.robot)
const counts = reactive({})
const countErrors = reactive({})
const apikeyCount = ref(null)
const apikeyEnabled = ref(false)
const apikeyError = ref(false)
const activeChannel = ref(null)
const secondaryLoading = ref(false)
const secondaryError = ref('')
const selectedCategory = ref('all')
const categories = [
  { key: 'all', label: 'category_all' },
  { key: 'wechat', label: 'category_wechat' },
  { key: 'social', label: 'category_social' },
  { key: 'other', label: 'category_other' }
]
let countRequest = 0

const channelComponents = {
  official_account: WeChatOfficialAccount,
  mini_program: WeChatMiniProgram,
  wechat_kefu: WeChatCustomerService,
  wecom_robot: WeComRobot,
  feishu_robot: FeishuRobot,
  dingtalk_robot: DingTalkRobot,
  messenger: MessengerRobot,
  telegram_robot: TelegramRobot,
  line: LineRobot,
  whatsapp: WhatsAppRobot
}
const secondaryReady = computed(() => String(robotDetail.value?.id) === String(customerRobot.value?.id))
const activeComponent = computed(() => activeChannel.value ? channelComponents[activeChannel.value.key] : null)
const activeTitleKey = computed(() => activeChannel.value?.key === 'api_key'
  ? 'api_key_manage'
  : activeChannel.value?.nameKey || 'other_channels')

provide('robotInfo', {
  get robotInfo() {
    return robotDetail.value
  },
  getRobot
})

const channelCards = computed(() => [
  {
    key: 'official_account',
    category: 'wechat',
    nameKey: 'channel_official_account',
    descKey: 'channel_official_account_desc',
    background: '#07c160',
    tab: 'WeChatOfficialAccount',
    icon: wechatIcon,
    active: (counts.official_account || 0) > 0
  },
  {
    key: 'mini_program',
    category: 'wechat',
    nameKey: 'channel_mini_program',
    descKey: 'channel_mini_program_desc',
    background: '#07c160',
    tab: 'WeChatMiniProgram',
    icon: miniProgramIcon,
    active: (counts.mini_program || 0) > 0
  },
  {
    key: 'wechat_kefu',
    category: 'wechat',
    nameKey: 'channel_wechat_customer_service',
    descKey: 'channel_wechat_customer_service_desc',
    background: '#07c160',
    tab: 'WeChatCustomerService',
    icon: wechatCustomerServiceIcon,
    smallIcon: true,
    active: (counts.wechat_kefu || 0) > 0
  },
  {
    key: 'wecom_robot',
    category: 'wechat',
    nameKey: 'channel_wecom_robot',
    descKey: 'channel_wecom_robot_desc',
    background: '#0082ef',
    tab: 'WeComRobot',
    icon: wecomIcon,
    smallIcon: true,
    active: (counts.wecom_robot || 0) > 0
  },
  {
    key: 'feishu_robot',
    category: 'other',
    nameKey: 'channel_feishu_robot',
    descKey: 'channel_feishu_robot_desc',
    background: '#ffffff',
    lightIcon: true,
    largeIcon: true,
    tab: 'FeishuRobot',
    icon: feishuIcon,
    active: (counts.feishu_robot || 0) > 0
  },
  {
    key: 'dingtalk_robot',
    category: 'other',
    nameKey: 'channel_dingtalk_robot',
    descKey: 'channel_dingtalk_robot_desc',
    background: '#0089ff',
    tab: 'DingDingRobot',
    icon: dingtalkIcon,
    active: (counts.dingtalk_robot || 0) > 0
  },
  {
    key: 'api_key',
    category: 'other',
    nameKey: 'channel_api_key',
    descKey: 'channel_api_key_desc',
    background: '#e8f1ff',
    tab: '',
    icon: apiKeyIcon,
    active: apikeyEnabled.value
  },
  {
    key: 'messenger',
    category: 'social',
    nameKey: 'channel_messenger',
    descKey: 'channel_messenger_desc',
    background: '#0866ff',
    tab: 'Messenget',
    icon: messengerIcon,
    active: (counts.messenger || 0) > 0
  },
  {
    key: 'telegram_robot',
    category: 'social',
    nameKey: 'channel_telegram',
    descKey: 'channel_telegram_desc',
    background: '#26a5e4',
    tab: 'TelegramRobot',
    icon: telegramIcon,
    active: (counts.telegram_robot || 0) > 0
  },
  {
    key: 'line',
    category: 'social',
    nameKey: 'channel_line',
    descKey: 'channel_line_desc',
    background: '#00c300',
    tab: 'Line',
    icon: lineIcon,
    active: (counts.line || 0) > 0
  },
  {
    key: 'whatsapp',
    category: 'social',
    nameKey: 'channel_whatsapp',
    descKey: 'channel_whatsapp_desc',
    background: '#25d366',
    tab: 'WhatsApp',
    icon: whatsappIcon,
    active: (counts.whatsapp || 0) > 0
  }
])

const filteredChannels = computed(() => channelCards.value.filter((channel) =>
  selectedCategory.value === 'all' || channel.category === selectedCategory.value
))

const badgeOf = (channel) => {
  if (channel.key === 'api_key') {
    if (apikeyError.value) return { text: t('status_unavailable'), active: false }
    if (apikeyCount.value === null) return null
    if (apikeyEnabled.value) return { text: t('status_enabled'), active: true }
    if (apikeyCount.value > 0) return { text: t('status_not_enabled'), active: false }
    return { text: t('status_not_connected'), active: false }
  }
  if (countErrors[channel.key]) return { text: t('status_unavailable'), active: false }
  const count = counts[channel.key]
  if (count === undefined) return null
  if (count > 0) return { text: t('status_connected', { count }), active: true }
  return { text: t('status_not_connected'), active: false }
}

const loadCounts = async () => {
  const request = ++countRequest
  const robot = customerRobot.value
  if (!robot?.id) return
  const channelTypes = {
    official_account: 'official_account',
    mini_program: 'mini_program',
    wechat_kefu: 'wechat_kefu',
    wecom_robot: 'wecom_robot',
    feishu_robot: 'feishu_robot',
    dingtalk_robot: 'dingtalk_robot',
    messenger: 'messenger',
    telegram_robot: 'telegram_robot',
    line: 'line',
    whatsapp: 'whatsapp'
  }
  Object.keys(counts).forEach((key) => delete counts[key])
  Object.keys(countErrors).forEach((key) => delete countErrors[key])
  apikeyCount.value = null
  apikeyEnabled.value = false
  apikeyError.value = false

  await Promise.all(Object.entries(channelTypes).map(async ([key, appType]) => {
    try {
      const res = await getWechatAppList({
        robot_id: robot.id,
        app_type: appType,
        app_name: '',
        is_all: 1
      })
      if (request !== countRequest) return
      counts[key] = Array.isArray(res?.data) ? res.data.length : 0
      delete countErrors[key]
    } catch {
      if (request !== countRequest) return
      delete counts[key]
      countErrors[key] = true
    }
  }))

  if (request !== countRequest || !robot.robot_key) return
  try {
    const res = await listRobotApikey({ robot_key: robot.robot_key })
    if (request !== countRequest) return
    const list = res?.data?.list || []
    apikeyCount.value = list.length
    apikeyEnabled.value = list.some((item) => Number(item.status) === 1)
    apikeyError.value = false
  } catch {
    if (request !== countRequest) return
    apikeyCount.value = 0
    apikeyEnabled.value = false
    apikeyError.value = true
  }
}

const openChannel = async (channel) => {
  if (!channel || !customerRobot.value?.id) return
  activeChannel.value = channel
  secondaryError.value = ''
  if (secondaryReady.value) return

  secondaryLoading.value = true
  try {
    const loaded = await getRobot(customerRobot.value.id)
    if (!loaded || !secondaryReady.value) secondaryError.value = 'error_robot'
  } catch {
    secondaryError.value = 'error_robot'
  } finally {
    secondaryLoading.value = false
  }
}

const backToChannels = () => {
  activeChannel.value = null
  secondaryError.value = ''
  secondaryLoading.value = false
}

watch(
  () => customerRobot.value?.id,
  () => {
    backToChannels()
    loadCounts()
  },
  { immediate: true }
)
onBeforeUnmount(() => { countRequest++ })
</script>

<style lang="less" scoped>
.other-channels-page {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-height: 0;
  margin: 8px 8px 8px 0;
  border-radius: 12px;
  overflow: hidden;
  background: var(--app-bg-container);
}
.channels-body {
  flex: 1;
  min-height: 0;
  padding: 32px 28px;
  overflow-y: auto;
}
.secondary-body {
  flex: 1;
  min-height: 0;
  padding: 32px 28px;
  overflow-y: auto;

  :deep(.ant-alert) {
    border-radius: var(--app-radius-control);
  }
  :deep(.tip-alert.ant-alert) {
    margin-bottom: 16px;
  }
  :deep(.ant-alert-banner) {
    margin-bottom: 16px;
  }
}
.secondary-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 240px;
}
.channels-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}
.channel-filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 24px;

  button {
    height: 32px;
    padding: 0 12px;
    border: 0;
    border-radius: 6px;
    background: var(--app-fill-secondary);
    color: var(--app-text-secondary);
    font: inherit;
    cursor: pointer;

    &.selected,
    &:hover {
      color: var(--app-primary);
      background: var(--app-selected);
    }

    &:focus-visible {
      outline: 2px solid var(--app-primary);
      outline-offset: 2px;
    }
  }
}
.channel-card {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
  min-height: 128px;
  padding: 20px;
  border: 1px solid var(--app-border);
  border-radius: 8px;
  background: var(--app-bg-container);
  font: inherit;
  text-align: left;
  cursor: pointer;
  transition: border-color 200ms, box-shadow 200ms;

  &:hover {
    border-color: var(--app-primary);
    box-shadow: var(--app-shadow-card);
  }

  &:focus-visible {
    outline: 2px solid var(--app-primary);
    outline-offset: 2px;
  }
}
.channel-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  border-radius: 10px;

  &.light {
    border: 1px solid var(--app-border);
  }

  img {
    width: 24px;
    height: 24px;
    object-fit: contain;

    &.large {
      width: 30px;
      height: 30px;
    }

    &.small {
      width: 20px;
      height: 20px;
    }
  }
}
.channel-content {
  flex: 1;
  min-width: 0;
}
.channel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 8px;
}
.channel-badge {
  display: inline-flex;
  align-items: center;
  flex-shrink: 0;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 4px;
  color: var(--app-text-secondary);
  background: var(--app-fill-secondary);
  font-size: 12px;

  &.green {
    // 深绿保证 12px 文本对比度，同时保留“已接入”状态语义。
    color: #067a52;
    background: #e7f8f1;
  }
}
.channel-name {
  min-width: 0;
  color: var(--app-text-primary);
  font-size: 16px;
  font-weight: 500;
}
.channel-desc,
.secondary-body :deep(.ant-alert-message) {
  color: var(--app-text-secondary);
  font-family: var(--app-font-family);
  font-size: 14px;
  font-weight: 400;
  line-height: 22px;
}
.channel-desc {
  display: block;
}
@media (max-width: 1200px) {
  .channels-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
@media (max-width: 760px) {
  .channels-grid {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
