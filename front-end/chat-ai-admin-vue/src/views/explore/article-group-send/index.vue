<template>
  <div v-if="!aiCommentEnabled" class="ability-banner">
    <svg-icon name="tip-icon" style="font-size: 16px; color: white;" />
    <span>{{ t('banner_ai_comment_not_enabled') }}</span>
    <a class="go-enable" href="javascript:void(0);" @click="goEnableAiComment">{{ t('common_go_enable') }}</a>
  </div>
  <div class="user-model-page">
  <div class="breadcrumb-wrap">
    <svg-icon @click="goBack" name="back" style="font-size: 20px;" />
    <span @click="goBack" class="breadcrumb-title">{{ t('breadcrumb_article_group_send') }}</span>
    <a-switch
      :checked="abilitySwitchChecked"
      :checked-children="t('switch_on')"
      :un-checked-children="t('switch_off')"
      @change="onAbilitySwitchChange"
    />
    <span class="switch-tip">{{ t('switch_tip') }}</span>
  </div>

    <div v-if="loadingApps" class="loading-box"><a-spin /></div>
    <template v-else>
      <div class="empty-wrap" v-if="accountList.length === 0">
        <ListEmpty size="180">
          <div class="empty-default">{{ t('empty_no_mp_bound') }}</div>
          <div class="empty-sub">{{ t('empty_bind_mp_tip') }}</div>
        </ListEmpty>
        <div class="empty-actions">
          <a-button type="primary" @click="toBindMp">{{ t('btn_bind_mp') }}</a-button>
        </div>
      </div>
      <template v-else>
        <div class="mp-list-block">
          <div class="mp-list" :class="{ expanded }" ref="mpListRef">
            <div class="mp-card" v-for="mp in (expanded ? mpAccounts : mpAccounts.slice(0, visibleCount))" :key="mp.id"
              :class="{ selected: mp.appid === selectedAppid }" @click="selectMp(mp)">
              <img :src="mp.logo" class="mp-logo" />
              <span class="mp-name">{{ mp.name }}</span>
            </div>
            <a-button v-if="!expanded && mpAccounts.length > visibleCount" type="dashed" class="more-btn"
              @click="expanded = true">
              {{ t('more_show') }} +{{ mpAccounts.length - visibleCount }}
            </a-button>
          </div>
        </div>

        <div class="tabs-box">
          <a-segmented
            v-model:value="activeKey"
            :options="segmentedOptions"
            @change="handleChangeTab"
          />
        </div>

        <a-alert class="alert-box" style="margin: 16px 48px 0;">
          <template #message>
            <div>{{ t('alert_line_1') }}</div>
            <div>{{ t('alert_line_2') }}</div>
          </template>
        </a-alert>

        <div class="list-wrapper">
          <div class="content-wrapper">
            <router-view v-slot="{ Component }">
              <component :is="Component" :app-id="currentAppId" :access-key="currentAccessKey" />
            </router-view>
          </div>
        </div>
      </template>
    </template>
  </div>
  <a-modal v-model:open="bannerTipOpen" :footer="null">
    <template #title>
      <span><ExclamationCircleFilled style="color:#faad14;margin-right:8px;" />{{ t('modal_ai_comment_not_enabled_title') }}</span>
    </template>
    <div>
      <span style="color: red;">{{ t('modal_ai_comment_not_enabled') }}</span>{{ t('modal_ai_comment_not_enabled_tip') }}
    </div>
    <div class="enable-tip-footer">
      <a-checkbox v-model:checked="bannerTipDontRemind">{{ t('checkbox_dont_remind_three_days') }}</a-checkbox>
      <div class="footer-actions">
        <a-button @click="onCancelBannerTip">{{ t('btn_cancel') }}</a-button>
        <a-button type="primary" @click="goEnableAiComment">{{ t('common_go_enable') }}</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { getSpecifyAbilityConfig, saveUserAbility } from '@/api/explore'
import { getWechatAppList } from '@/api/robot'
import { ExclamationCircleFilled } from '@ant-design/icons-vue'
import ListEmpty from '@/views/robot/robot-config/function-center/components/list-empty.vue'
import { useI18n } from '@/hooks/web/useI18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n('views.explore.article-group-send.index')

const segmentedOptions = [
  { label: t('segmented_group_send'), value: 'group-send' },
  { label: t('segmented_draft_box'), value: 'draft-box' }
]

const getActiveByRoute = () => {
  if (route.name === 'exploreArticleGroupSendDraftBox') return 'draft-box'
  return 'group-send'
}

const goBack = () => {
  if (router.options.history.state.back) {
    router.back()
  } else {
    router.push({ path: '/explore/index' })
  }
}

const activeKey = ref(getActiveByRoute())

const abilitySwitchChecked = ref(false)
const initAbilitySwitch = () => {
  getSpecifyAbilityConfig({ ability_type: 'official_account_batch_send' }).then((res) => {
    const item = res.data
    const st = item?.user_config?.switch_status
    abilitySwitchChecked.value = String(st || '0') === '1'
  })
}

const aiCommentEnabled = ref(true)
const initAiCommentSwitch = () => {
  getSpecifyAbilityConfig({ ability_type: 'official_account_ai_comment' }).then((res) => {
    const st = res?.data?.user_config?.switch_status
    aiCommentEnabled.value = String(st || '0') === '1'
    tryOpenBannerTip()
  })
}
const goEnableAiComment = () => {
  applyDontRemindIfChecked()
  const url = router.resolve({ path: '/explore/index'})
  window.open(url.href, '_blank')
}
const onAbilitySwitchChange = (checked) => {
  const newStatus = checked ? '1' : '0'
  if (newStatus === '0') {
    Modal.confirm({
      title: t('modal_confirm_title'),
      content: t('modal_confirm_disable_content'),
      onOk: () => {
        saveUserAbility({ ability_type: 'official_account_batch_send', switch_status: newStatus }).then((res) => {
          if (res && res.res == 0) {
            abilitySwitchChecked.value = false
            message.success(t('message_operation_success'))
          }
        })
      }
    })
    return
  }
  saveUserAbility({ ability_type: 'official_account_batch_send', switch_status: newStatus }).then((res) => {
    if (res && res.res == 0) {
      abilitySwitchChecked.value = true
      message.success(t('message_operation_success'))
    }
  })
}

watch(
  () => route.name,
  () => {
    activeKey.value = getActiveByRoute()
  }
)

const handleChangeTab = (val) => {
  const pathMap = {
    'group-send': '/explore/index/article-group-send/group-send',
    'draft-box': '/explore/index/article-group-send/draft-box'
  }
  router.replace(pathMap[val])
}

onMounted(() => {
  if (route.name === 'exploreArticleGroupSend') {
    router.replace('/explore/index/article-group-send/group-send')
  }
  loadApps()
  initAbilitySwitch()
  initAiCommentSwitch()
})

onUnmounted(() => { window.removeEventListener('resize', calcVisibleCount) })

const loadingApps = ref(true)
const accountList = ref([])
const currentAppId = ref('')
const currentAccessKey = ref('')
const mpAccounts = ref([])
const selectedAppid = ref('')
const expanded = ref(false)
const mpListRef = ref(null)
const visibleCount = ref(0)
const CARD_WIDTH = 160
const GAP = 8
const MORE_BTN_WIDTH = 96

function calcVisibleCount () {
  const el = mpListRef.value
  if (!el) { visibleCount.value = 0; return }
  const w = el.clientWidth || 0
  const per = CARD_WIDTH + GAP
  const count = Math.floor((w - MORE_BTN_WIDTH) / per)
  visibleCount.value = Math.max(count, 0)
}

const loadApps = () => {
  loadingApps.value = true
  getWechatAppList({ app_type: 'official_account' }).then((res) => {
    const list = res.data || []
    accountList.value = list
    mpAccounts.value = list.map(app => ({ id: app.id, appid: app.app_id, name: app.app_name, logo: app.app_avatar, access_key: app.access_key }))
    const first = mpAccounts.value[0]
    if (first) {
      selectedAppid.value = first.appid
      currentAppId.value = first.appid
      currentAccessKey.value = first.access_key
    }
  }).finally(() => {
    loadingApps.value = false
    nextTick(calcVisibleCount)
    window.addEventListener('resize', calcVisibleCount)
  })
}

function selectMp (mp) {
  selectedAppid.value = mp.appid
  expanded.value = true
  currentAppId.value = mp.appid
  currentAccessKey.value = mp.access_key
}

const toBindMp = () => {
  const querUrl = router.resolve({ path: '/user/official-account' })
  window.open(querUrl.href, '_blank')
}

const bannerTipOpen = ref(false)
const bannerTipDontRemind = ref(false)
const DONT_REMIND_KEY = 'article_group_send_ai_comment_tip_until'
const tryOpenBannerTip = () => {
  const until = Number(localStorage.getItem(DONT_REMIND_KEY) || 0)
  const now = Date.now()
  if (!aiCommentEnabled.value && until <= now) {
    bannerTipOpen.value = true
  }
}
const applyDontRemindIfChecked = () => {
  if (bannerTipDontRemind.value) {
    const until = Date.now() + 3 * 24 * 60 * 60 * 1000
    localStorage.setItem(DONT_REMIND_KEY, String(until))
  }
}
const onCancelBannerTip = () => {
  applyDontRemindIfChecked()
  bannerTipOpen.value = false
}
</script>

<style lang="less" scoped>
.ability-banner {
  height: 40px;
  box-sizing: border-box;
  padding: 9px 48px;
  background: #FFF0EB;
  color: #ed744a;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}
.ability-banner .go-enable {
  color: #2475fc;
}
.enable-tip-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 16px;
}
.footer-actions {
  display: flex;
  gap: 8px;
}
.user-model-page {
  width: 100%;
  height: auto;
  background-color: #fff;
  display: flex;
  flex-direction: column;
  overflow: auto;
  padding: 16px 0px;

  .breadcrumb-wrap {
    width: fit-content;
    display: flex;
    align-items: center;
    cursor: pointer;
    margin: 0 48px 16px;
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

  .tabs-box {
    background: #fff;
    padding-top: 16px;
    margin: 0 48px;
    &::v-deep(.ant-segmented) {
      background: #e4e6eb;
    }
    &::v-deep(.ant-segmented .ant-segmented-item) {
      color: #262626;
    }
    &::v-deep(.ant-segmented .ant-segmented-item-selected) {
      color: #2475fc;
    }
    &::v-deep(.ant-segmented .ant-segmented-item-label) {
      padding: 0 16px;
    }
  }

  .list-wrapper {
    margin: 0 48px;
    background: #fff;
    flex: 1;
    // overflow-x: hidden;
    // overflow-y: auto;
  }
  .content-wrapper {
    padding-top: 16px;
  }
}

.line-box {
  height: 1px;
  background: #F0F0F0;
  margin-top: 16px;
}

.loading-box {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 0;
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

.mp-list-block {
  margin: 0px 0 4px 0;
}

.mp-list {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: nowrap;
  margin: 0 48px;
}

.mp-list.expanded {
  flex-wrap: wrap;
}

.mp-card {
  width: 160px;
  padding: 8px 12px;
  border-radius: 8px;
  background: #fff;
  border: 1px solid #edeff2;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.mp-card.selected {
  border-color: #2475fc;
  box-shadow: 0 0 0 2px rgba(36, 117, 252, 0.1);
}

.mp-logo {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.mp-name {
  color: #262626;
  font-size: 14px;
  line-height: 22px;
}

.more-btn {
  flex: 0 0 auto;
}
</style>
