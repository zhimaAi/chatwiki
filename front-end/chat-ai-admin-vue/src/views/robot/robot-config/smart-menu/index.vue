<template>
  <div class="auto-reply-home">
    <!-- <template v-if="hasApps"> -->
      <SmartMenu />
    <!-- </template> -->
    <!-- <template v-else>
      <div class="user-model-page">
        <div class="page-title">智能菜单</div>
        <a-alert show-icon>
          <template #message>
            开启后，该智能菜单可以被关键词回复等功能调用
          </template>
        </a-alert>
        <div class="empty-wrap">
          <ListEmpty size="180">
            <div class="empty-default">暂未绑定公众号</div>
            <div class="empty-sub">请先到对外服务中 机器人>对外服务>微信公众号去绑定公众号</div>
          </ListEmpty>
          <div class="empty-actions">
            <a-button type="primary" @click="goBind">绑定</a-button>
          </div>
        </div>
      </div>
    </template> -->
  </div>
  
</template>

<script setup>
// import { ref, onMounted } from 'vue'
// import { useRoute, useRouter } from 'vue-router'
// import { getWechatAppList } from '@/api/robot'
import SmartMenu from './smart-menu.vue'
// import ListEmpty from '@/views/robot/robot-config/function-center/components/list-empty.vue'

// const route = useRoute()
// const router = useRouter()
// const hasApps = ref(true)

// const fetchApps = async () => {
//   try {
//     const { id } = route.query
//     const res = await getWechatAppList({ robot_id: id, app_type: 'official_account', app_name: '' })
//     const list = Array.isArray(res?.data) ? res.data : []
//     hasApps.value = list.length > 0
//   } catch (_e) {
//     hasApps.value = false
//   }
// }

// const goBind = () => {
//   localStorage.setItem('/robot/config/external-services/activeKey', 'WeChatOfficialAccount')
//   const url = router.resolve({ path: '/robot/config/external-services', query: { id: route.query.id, robot_key: route.query.robot_key } })
//   window.open(url.href, '_blank')
// }

// onMounted(fetchApps)
</script>

<style lang="less" scoped>
.auto-reply-home {
  height: 100%;
  width: 100%;
  padding: 16px;
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
</style>