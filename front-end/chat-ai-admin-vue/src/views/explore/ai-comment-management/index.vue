<template>
  <div v-if="!batchSendEnabled" class="ability-banner">
    <svg-icon name="tip-icon" style="font-size: 16px; color: white;" />
    <span>文章群发功能暂未开启，您可到探索 > 功能 > 文章群发开启</span>
    <a class="go-enable" href="javascript:void(0);" @click="goEnableBatchSend">去开启</a>
  </div>
  <div class="user-model-page">
  <div class="breadcrumb-wrap">
    <svg-icon name="back" style="font-size: 20px;" @click="goBack" />
    <div v-if="isCopyCustomRule" class="breadcrumb-title" @click="goBack">复制自定义规则</div>
    <div v-else-if="isEditCustomRule" class="breadcrumb-title" @click="goBack">编辑自定义规则</div>
    <div v-else-if="isCreateCustomRule" class="breadcrumb-title" @click="goBack">新建自定义规则</div>
    <template v-else>
      <span class="breadcrumb-title" @click="goBack">AI评论精选</span>
      <a-switch
        :checked="abilitySwitchChecked"
        checked-children="开"
        un-checked-children="关"
        @change="onAbilitySwitchChange"
      />
      <span class="switch-tip">开启后，可设置规则对系统内发布的群发使用AI自动精选功能，自动删评，回复，精选评论。</span>
    </template>
  </div>

    <a-alert class="alert-box" show-icon style="margin: 16px 48px 0;">
      <template #message>
        <div>仅支持通过当前系统群发的文章做Al评论管理，按照设置的规则自动删评 <br /> 自动回复评论需开启文章群发功能并创建文章群发，<a @click="goToArticleGroupSend" href="javascript:void(0);">去创建</a></div>
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
      <span><ExclamationCircleFilled style="color:#faad14;margin-right:8px;" />未开启文章群发</span>
    </template>
    <div><span style="color: red;">文章群发功能暂未开启</span>，您可到探索 > 功能 > 文章群发开启</div>
    <div class="enable-tip-footer">
      <a-checkbox v-model:checked="bannerTipDontRemind">3天内不在提示</a-checkbox>
      <div class="footer-actions">
        <a-button @click="onCancelBannerTip">取消</a-button>
        <a-button type="primary" @click="goEnableBatchSend">去开启</a-button>
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

const route = useRoute()
const router = useRouter()
const isCreateCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule')
const isCopyCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule' && route.query.copy === '1')
const isEditCustomRule = computed(() => route.name === 'exploreAiCommentManagementCreateCustomRule' && route.query.id)

const segmentedOptions = [
  { label: '默认规则', value: 'default-rule' },
  { label: '自定义规则', value: 'custom-rule' },
  { label: '评论处理记录', value: 'comment-processing-record' }
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
      title: '提示',
      content: '关闭后，该功能默认关闭不再支持使用，确认关闭？',
      onOk: () => {
        saveUserAbility({ ability_type: 'official_account_ai_comment', switch_status: newStatus }).then((res) => {
          if (res && res.res == 0) {
            abilitySwitchChecked.value = false
            message.success('操作成功')
          }
        })
      }
    })
    return
  }
  saveUserAbility({ ability_type: 'official_account_ai_comment', switch_status: newStatus }).then((res) => {
    if (res && res.res == 0) {
      abilitySwitchChecked.value = true
      message.success('操作成功')
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
