<template>
  <div class="user-model-page">
    <div class="page-title">{{ t('title_function_center') }}</div>
    <div class="list-wrapper">
      <div class="content-wrapper">
        <a-alert show-icon style="align-items: baseline">
          <template #message>
            <div>
              {{ t('msg_function_center_desc') }} <a href="#/explore/index" target="_blank">{{ t('link_explore_function_center') }}</a>{{ t('msg_function_center_desc_suffix') }}
            </div>
          </template>
        </a-alert>
        <div bg-color="#f5f9ff" class="explore-page">
          <div class="explore-page-body">
            <div class="list-group-box">
              <cu-scroll style="padding-right: 16px; flex: 1">
                <ListEmpty v-if="list.length == 0" size="250" class="empty-box">
                  <div class="empty-content">
                    <p class="empty-title">{{ t('empty_title') }}</p>
                    <p class="empty-desc">{{ t('empty_desc_prefix') }} <a href="#/explore/index" target="_blank">{{ t('link_explore_function_center') }}</a> {{ t('empty_desc_suffix') }}</p>
                    <a-button type="primary" @click="handleCreateClick">{{ t('btn_add') }}</a-button>
                  </div>
                </ListEmpty>
                <ExploreList v-else :list="list" @switchChange="handleSwitchChange" @fixedMenuChange="handleFixedMenuChange" @clickItem="handleClickItem" />
              </cu-scroll>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, h } from 'vue'
import { message, Modal  } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import {
  getRobotAbilityList,
  saveRobotAbilitySwitchStatus,
  saveRobotAbilityFixedMenu
} from '@/api/explore'
import ExploreList from './components/explore-list/index.vue'
import ListEmpty from './components/list-empty.vue'

const { t } = useI18n('views.robot.robot-config.function-center.index')

// 顶部页签暂不使用
const tabs = ref([
  {
    title: 'tab_function',
    value: '1'
  },
  {
    title: 'tab_plugin',
    value: '2'
  }
])
const route = useRoute()
const router = useRouter()

const list = ref([])
const updateTabNumber = () => {
  tabs.value = [
    { title: t('tab_function'), value: '1' },
    { title: t('tab_plugin'), value: '2' }
  ]
}

const robotId = ref('')
let allList = ref([])
const getList = () => {
  getRobotAbilityList({ robot_id: robotId.value }).then((res) => {
    const data = (res?.data || []).map((it) => ({
      ...it,
      explore_name: it.name,
      explore_intro: it.introduction
    }))
    list.value = data
    allList.value = res?.data || []
    updateTabNumber()
  })
}

const handleSwitchChange = (item, checked) => {
  const newStatus = checked ? '1' : '0'
  // 如果是关闭二次弹窗提示：提示：弹框提示：关闭后，触发了关键词不会再回复指定内容
  if (newStatus == '0') {
    let tip = t('msg_close_keyword_no_reply')
    if (item.ability_type == 'robot_smart_menu') {
      tip = t('msg_close_smart_menu_no_reply')
    }
    if (item.ability_type == 'robot_auto_reply') {
      tip = t('msg_close_auto_reply_no_content')
    }
    if (item.ability_type == 'robot_payment') {
      tip = t('msg_close_confirm_payment')
    }
    Modal.confirm({
      title: t('title_confirm_close'),
      icon: h(ExclamationCircleOutlined),
      content: tip,
      onOk: () => {
        saveRobotAbilitySwitchStatus({ robot_id: robotId.value, ability_type: item.ability_type, switch_status: newStatus }).then((res) => {
          if (res && res.res == 0) {
            message.success(t('msg_operation_success'))
            if (item.robot_config) {
              item.robot_config.switch_status = newStatus
            }
            // 通知左侧菜单刷新能力列表（用于动态菜单）
            if (item.robot_only_show != 1) {
              window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: robotId.value } }))
            }
          }
        })
      }
    })
    return
  }
  saveRobotAbilitySwitchStatus({ robot_id: robotId.value, ability_type: item.ability_type, switch_status: newStatus }).then((res) => {
    if (res && res.res == 0) {
      message.success(t('msg_operation_success'))
      if (item.robot_config) {
        item.robot_config.switch_status = newStatus
      }
      window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: robotId.value } }))
    }
  })
}

const handleFixedMenuChange = (item, checked) => {
  const newStatus = checked ? '1' : '0'
  saveRobotAbilityFixedMenu({ robot_id: robotId.value, ability_type: item.ability_type, fixed_menu: newStatus }).then((res) => {
    if (res && res.res == 0) {
      message.success(t('msg_operation_success'))
      if (item.robot_config) {
        item.robot_config.fixed_menu = newStatus
      }
      window.dispatchEvent(new CustomEvent('robotAbilityUpdated', { detail: { robotId: robotId.value } }))
    }
  })
}

const handleClickItem = (item) => {
  let url
  if (item.robot_only_show != 1) {
    url = router.resolve({
      path: item.menu.path,
      query: {
        id: robotId.value,
        robot_key: route.query.robot_key
      }
    })
  } else {
    router.push({
      path: item.menu.path,
      query: {
        id: robotId.value,
        robot_key: route.query.robot_key
      }
    })
    return
  }
  window.open(url.href, '_blank')
}

const handleCreateClick = () => {
  const url = router.resolve({
    path: '/explore/index',
    query: {
      tab: '1'
    }
  })
  window.open(url.href, '_blank')
}


onMounted(() => {
  robotId.value = route.query.id
  getList()
})
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
  margin-top: 4px;
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
// 大于1920px
@media screen and (min-width: 1920px) {
  .library-page {
  }
}

.user-model-page {
  width: 100%;
  height: 100%;
  background-color: #f2f4f7;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 16px;
    background-color: #fff;
    color: #000000;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }

  .list-wrapper {
    background: #fff;
    flex: 1;
    overflow-x: hidden;
    overflow-y: auto;
  }
  .content-wrapper {
    padding: 16px;
    padding-top: 0;
  }
  .actions-box {
    margin: 16px 0;
    display: flex;
  }
  .avatar-box {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    img {
      width: 32px;
      height: 32px;
    }
  }

  .status-block{
    display: flex;
    gap: 4px;
    align-items: center;
    color: #8c8c8c;
    span{
      width: 8px;
      height: 8px;
      border-radius: 8px;
      display: block;
      background: #8c8c8c;
    }
    &.success{
      color: #52C41A;
      span{
        background: #52C41A;
      }
    }
  }
}

.empty-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 0;
  p{
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
  a{
    color: #2475fc;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
  .empty-title{
    color: #262626;
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }
  .empty-desc{
    color: #595959;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
    margin: 12px 0 24px;
  }
}
</style>
