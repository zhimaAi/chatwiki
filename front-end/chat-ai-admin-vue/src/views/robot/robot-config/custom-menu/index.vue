<template>
  <div class="auto-reply-home">
    <template v-if="hasApps">
      <CustomMenu />
    </template>
    <template v-else>
      <div class="user-model-page">
        <div class="breadcrumb-wrap">
          <svg-icon @click="goBack" name="back" style="font-size: 20px;" />
          <div @click="goBack" class="breadcrumb-title">自定义菜单</div>
          <a-switch
            :checked="abilitySwitchChecked"
            checked-children="开"
            un-checked-children="关"
            @change="onAbilitySwitchChange"
          />
          <span class="switch-tip">开启后，公众号显示自定义菜单；关闭后，粉丝将无法看到菜单。</span>
        </div>
        <a-alert show-icon>
          <template #message>
            开启后，用户关注公众号后，回复指定的内容，<span style="color: #FF4D4F;">该功能仅支持公众号内回复</span>
          </template>
        </a-alert>
        <div class="empty-wrap">
          <ListEmpty size="180">
            <div class="empty-default">暂未绑定公众号</div>
            <div class="empty-sub">请选到系统设置>公众号管理绑定公众号</div>
          </ListEmpty>
          <div class="empty-actions">
            <a-button type="primary" @click="goBind">去绑定公众号</a-button>
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
import CustomMenu from './custom-menu.vue'
import ListEmpty from '@/views/robot/robot-config/function-center/components/list-empty.vue'
import { getSpecifyAbilityConfig, saveUserAbility } from '@/api/explore'
import { message, Modal } from 'ant-design-vue'

const route = useRoute()
const router = useRouter()
const hasApps = ref(true)
const abilitySwitchChecked = ref(false)

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

const initAbilitySwitch = () => {
  getSpecifyAbilityConfig({ ability_type: 'official_custom_menu' }).then((res) => {
    const st = res?.data?.user_config?.switch_status
    abilitySwitchChecked.value = String(st || '0') === '1'
  })
}

const onAbilitySwitchChange = (checked) => {
  const newStatus = checked ? '1' : '0'
  if (newStatus === '0') {
    Modal.confirm({
      title: '提示',
      content: '关闭后，该功能默认关闭不再支持使用，所有的公众号菜单都会停用，确认关闭？',
      onOk: () => {
        saveUserAbility({ ability_type: 'official_custom_menu', switch_status: newStatus }).then((res) => {
          if (res && res.res == 0) {
            abilitySwitchChecked.value = false
            message.success('操作成功')
          }
        })
      }
    })
    return
  }
  saveUserAbility({ ability_type: 'official_custom_menu', switch_status: newStatus }).then((res) => {
    if (res && res.res == 0) {
      abilitySwitchChecked.value = true
      message.success('操作成功')
    }
  })
}

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

onMounted(() => {
  initAbilitySwitch()
  fetchApps()
})
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
