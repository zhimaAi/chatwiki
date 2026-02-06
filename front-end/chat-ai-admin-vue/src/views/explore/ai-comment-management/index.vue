<template>
  <div v-if="!batchSendEnabled" class="ability-banner">
    <svg-icon name="tip-icon" style="font-size: 16px; color: white;" />
    <span>{{ t('banner_batch_send_not_enabled') }}</span>
    <a class="go-enable" href="javascript:void(0);" @click="goEnableBatchSend">{{ t('common_go_enable') }}</a>
  </div>
  <div class="user-model-page">
  <div class="breadcrumb-wrap">
    <svg-icon name="back" style="font-size: 20px;" @click="goBack" />
    <div v-if="isCopyCustomRule" class="breadcrumb-title" @click="goBack">{{ t('breadcrumb_copy_custom_rule') }}</div>
    <div v-else-if="isEditCustomRule" class="breadcrumb-title" @click="goBack">{{ t('breadcrumb_edit_custom_rule') }}</div>
    <div v-else-if="isCreateCustomRule" class="breadcrumb-title" @click="goBack">{{ t('breadcrumb_create_custom_rule') }}</div>
    <template v-else>
      <span class="breadcrumb-title" @click="goBack">{{ t('breadcrumb_ai_comment_selection') }}</span>
      <a-switch
        :checked="abilitySwitchChecked"
        :checked-children="t('switch_on')"
        :un-checked-children="t('switch_off')"
        @change="onAbilitySwitchChange"
      />
      <span class="switch-tip">{{ t('switch_tip') }}</span>
    </template>
  </div>

    <a-alert class="alert-box" show-icon style="margin: 16px 48px 0;">
      <template #message>
        <div>
          <span>{{ t('alert_support_tip') }}</span>
          <br />
          <span>{{ t('alert_auto_reply_tip') }}</span>
          <a @click="goToArticleGroupSend" href="javascript:void(0);">{{ t('alert_go_create') }}</a>
        </div>
      </template>
    </a-alert>

    <div class="tabs-box" v-if="!isEditCustomRule && !isCopyCustomRule && !isCreateCustomRule">
      <a-segmented
        v-model:value="activeKey"
        :options="segmentedOptions"
        @change="handleChangeTab"
      />
    </div>

    <div class="list-wrapper">
      <div class="content-wrapper">

        <router-view />
      </div>
    </div>
  </div>
  <a-modal v-model:open="bannerTipOpen" :footer="null">
    <template #title>
      <span><ExclamationCircleFilled style="color:#faad14;margin-right:8px;" />{{ t('modal_article_batch_send_not_enabled_title') }}</span>
    </template>
    <div>
      <span style="color: red;">{{ t('modal_article_batch_send_not_enabled') }}</span>{{ t('modal_article_batch_send_not_enabled_tip') }}
    </div>
    <div class="enable-tip-footer">
      <a-checkbox v-model:checked="bannerTipDontRemind">{{ t('checkbox_dont_remind_three_days') }}</a-checkbox>
      <div class="footer-actions">
        <a-button @click="onCancelBannerTip">{{ t('btn_cancel') }}</a-button>
        <a-button type="primary" @click="goEnableBatchSend">{{ t('common_go_enable') }}</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { getSpecifyAbilityConfig, saveUserAbility } from '@/api/explore'
import { ExclamationCircleFilled } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const route = useRoute()
const router = useRouter()
const { t } = useI18n('views.explore.ai-comment-management.index')
const isCreateCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule')
const isCopyCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule' && route.query.copy === '1')
const isEditCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule' && route.query.id)

const segmentedOptions = [
  { label: t('segmented_default_rule'), value: 'default-rule' },
  { label: t('segmented_custom_rule'), value: 'custom-rule' },
  { label: t('segmented_comment_processing_record'), value: 'comment-processing-record' }
]

const getActiveByRoute = () => {
  if (route.name === 'exploreAiCommentManagementCustomRule') return 'custom-rule'
  if (route.name === 'exploreAiCommentManagementCommentProcessingRecord') return 'comment-processing-record'
  return 'default-rule'
}

const goBack = () => {
  // 如果有浏览器记录，就返回上一页
  if (router.options.history.state.back) {
    router.back()
  } else {
    router.push({ path: '/explore/index' })
  }
}

const goToArticleGroupSend = () => {
  const url = router.resolve({ path: '/explore/index/article-group-send' })
  window.open(url.href, '_blank')
}

const activeKey = ref(getActiveByRoute())

// 能力开关：AI评论管理
const abilitySwitchChecked = ref(false)
const initAbilitySwitch = () => {
  getSpecifyAbilityConfig({ ability_type: 'official_account_ai_comment' }).then((res) => {
    const item = res.data
    const st = item?.user_config?.switch_status
    abilitySwitchChecked.value = String(st || '0') == '1'
  })
}

const batchSendEnabled = ref(true)
const initBatchSendSwitch = () => {
  getSpecifyAbilityConfig({ ability_type: 'official_account_batch_send' }).then((res) => {
    const st = res?.data?.user_config?.switch_status
    batchSendEnabled.value = String(st || '0') === '1'
    tryOpenBannerTip()
  })
}
const onAbilitySwitchChange = (checked) => {
  const newStatus = checked ? '1' : '0'
  if (newStatus === '0') {
    Modal.confirm({
      title: t('modal_confirm_title'),
      content: t('modal_confirm_disable_content'),
      onOk: () => {
        saveUserAbility({ ability_type: 'official_account_ai_comment', switch_status: newStatus }).then((res) => {
          if (res && res.res == 0) {
            abilitySwitchChecked.value = false
            message.success(t('message_operation_success'))
          }
        })
      }
    })
    return
  }
  saveUserAbility({ ability_type: 'official_account_ai_comment', switch_status: newStatus }).then((res) => {
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
    'default-rule': '/explore/index/ai-comment-management/default-rule',
    'custom-rule': '/explore/index/ai-comment-management/custom-rule',
    'comment-processing-record': '/explore/index/ai-comment-management/comment-processing-record'
  }
  router.replace(pathMap[val])
}

onMounted(() => {
  if (route.name === 'exploreAiCommentManagement') {
    router.replace('/explore/index/ai-comment-management/default-rule')
  }
  initAbilitySwitch()
  initBatchSendSwitch()
})

const goEnableBatchSend = async () => {
  applyDontRemindIfChecked()
  const url = router.resolve({ path: '/explore/index'})
  window.open(url.href, '_blank')
}

const bannerTipOpen = ref(false)
const bannerTipDontRemind = ref(false)
const DONT_REMIND_KEY = 'ai_comment_management_batch_send_tip_until'
const tryOpenBannerTip = () => {
  const until = Number(localStorage.getItem(DONT_REMIND_KEY) || 0)
  const now = Date.now()
  if (!batchSendEnabled.value && until <= now) {
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
  height: 100%;
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
</style>
