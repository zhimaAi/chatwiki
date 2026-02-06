<template>
  <div bg-color="#f5f9ff" class="explore-page">
    <div class="explore-page-body">
      <div class="list-toolbar">
        <div class="toolbar-box">
          <MainTab ref="tabRef"/>
        </div>
      </div>

      <div class="list-group-box">
        <cu-scroll style="padding-right: 16px; flex: 1">
          <ExploreList :list="list" @switchChange="handleSwitchChange" @clickItem="handleClick" />
        </cu-scroll>
      </div>
    </div>
    <a-modal v-model:open="enableTipOpen" :title="t('title_tip')" :footer="null">
      <div>{{ t('msg_enable_feature_tip') }}</div>
      <div class="enable-tip-footer">
        <a-checkbox v-model:checked="enableTipDontRemind">{{ t('label_dont_remind') }}</a-checkbox>
        <div class="footer-actions">
          <a-button @click="onCancelTip">{{ t('btn_cancel') }}</a-button>
          <a-button type="primary" @click="goToFunctionCenter">{{ t('btn_go_to_use') }}</a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRouter } from 'vue-router'
import { getRobotList } from '@/api/robot/index.js'
import {
  getAbilityList,
  saveUserAbility
} from '@/api/explore'
import ExploreList from './components/explore-list/index.vue'
import MainTab from "@/views/explore/components/main-tab.vue"
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.explore.explore-index.explore-index')

// 顶部页签暂不使用

const tabRef = ref(null)
const list = ref([])
const enableTipOpen = ref(false)
const enableTipDontRemind = ref(false)
const DONT_REMIND_KEY = 'explore_enable_tip_suppress_until'
const router = useRouter()

onMounted(() => {
  let _active = localStorage.getItem('zm:explore:active')
  if (_active > 1) {
    tabRef.value.change(_active)
  } else {
    getList()
  }
})


let allList = ref([])
const getList = () => {
  // let type = activeKey.value === '1' ? '' : activeKey.value

  getAbilityList().then((res) => {
    const data = (res?.data || []).map((it) => ({
      ...it,
      explore_name: it.name,
      explore_intro: it.introduction
    }))
    list.value = data
    allList.value = res?.data || []
  })
}

const handleSwitchChange = (item, checked) => {
  const newStatus = checked ? '1' : '0'
  if (newStatus === '0') {
    const contMap = {
      'robot_auto_reply': t('msg_close_confirm_default'),
      'library_ability_official_account': t('msg_close_confirm_library'),
      'robot_payment': t('msg_close_confirm_payment'),
    }
    Modal.confirm({
      title: t('title_tip'),
      content: contMap[item.ability_type] || t('msg_close_confirm_default'),
      onOk: () => {
        saveUserAbility({ ability_type: item.ability_type, switch_status: newStatus }).then((res) => {
          if (res && res.res == 0) {
            message.success(t('msg_operation_success'))
            if (item.user_config) {
              item.user_config.switch_status = newStatus
            } else if (item.robot_config) {
              item.robot_config.switch_status = newStatus
            }
          }
        })
      }
    })
    return
  }
  saveUserAbility({ ability_type: item.ability_type, switch_status: newStatus }).then((res) => {
    if (res && res.res == 0) {
      message.success(t('msg_operation_success'))
      if (item.user_config) {
        item.user_config.switch_status = newStatus
      } else if (item.robot_config) {
        item.robot_config.switch_status = newStatus
      }
      if (item.ability_type === 'library_ability_official_account') return
      const until = Number(localStorage.getItem(DONT_REMIND_KEY) || 0)
      const now = Date.now()
      if (until <= now) {
        enableTipDontRemind.value = false
        if (item.robot_only_show != 1) {
          enableTipOpen.value = true
        }
      }
    }
  })
}

const handleClick = (item) => {
  switch (item.ability_type) {
    case 'library_ability_official_account':
      // eslint-disable-next-line no-case-declarations
      const url = router.resolve({ path: '/library/list', query: { active: 3 }})
      window.open(url.href, '_blank')
      break
    default:
      if (item.module_type === 'robot' && item.robot_only_show == 0) {
        goToFunctionCenter(item)
      }
      break
  }
}

const applyDontRemindIfChecked = () => {
  if (enableTipDontRemind.value) {
    const until = Date.now() + 3 * 24 * 60 * 60 * 1000
    localStorage.setItem(DONT_REMIND_KEY, String(until))
  }
}

const onCancelTip = () => {
  applyDontRemindIfChecked()
  enableTipOpen.value = false
}

const goToFunctionCenter = async (item) => {
  applyDontRemindIfChecked()
  enableTipOpen.value = false
  const rid = localStorage.getItem('last_local_robot_id')
  try {
    const { data: lists = [] } = await getRobotList()
    if (!lists.length) {
      const url = router.resolve({ path: '/robot/list' })
      window.open(url.href, '_blank')
      return
    }
    let toDetailRobot
    if (rid) {
      toDetailRobot = lists.find((item) => item.id == rid)
    }
    toDetailRobot = toDetailRobot || lists[0]
    const { id, robot_key } = toDetailRobot
    if (item.menu?.path) {
      const url = router.resolve({ path: item.menu.path, query: { id, robot_key } })
      window.open(url.href, '_blank')
      return
    }
    const url = router.resolve({ path: '/robot/config/function-center', query: { id, robot_key } })
    window.open(url.href, '_blank')
  } catch (e) {
    const url = router.resolve({ path: '/robot/list' })
    window.open(url.href, '_blank')
  }
}

watch(enableTipOpen, (v, ov) => {
  if (ov && !v) {
    applyDontRemindIfChecked()
  }
})
// 能力列表不支持分组删除与拖拽排序，相关逻辑移除
</script>

<style lang="less" scoped>
.explore-page {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  .list-toolbar {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;
  }
}
.create-action {
  display: flex;
  align-items: center;
  .icon {
    width: 20px;
    height: 20px;
    margin-right: 8px;
  }
}

.toolbar-box {
  padding-right: 16px;
}

.explore-page-body {
  padding: 0;
  margin-top: 16px;
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.list-group-box {
  display: flex;
  gap: 16px;
  flex: 1;
  overflow: hidden;
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
// 大于1920px
@media screen and (min-width: 1920px) {
  .library-page {
  }
}
</style>
