<template>
  <div class="auto-reply-home">
    <template v-if="hasApps">
      <ReceivedReply />
    </template>
    <template v-else>
      <div class="user-model-page">
        <!-- <div class="page-title">关注后自动回复</div> -->
        <div class="breadcrumb-wrap" @click="goBack">
          <svg-icon name="back" style="font-size: 20px;" />
          <div class="breadcrumb-title">关注后回复</div>
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
            <a-button type="primary" @click="goBind">绑定</a-button>
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

const route = useRoute()
const router = useRouter()
const hasApps = ref(true)

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

onMounted(fetchApps)
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
  gap: 10px;
  cursor: pointer;
  margin-bottom: 16px;
}
.breadcrumb-title {
  color: #000000;
  font-size: 20px;
  font-style: normal;
  font-weight: 600;
  line-height: 28px;
}
</style>