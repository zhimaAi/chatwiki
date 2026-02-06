<template>
  <div class="auto-reply-home">
    <a-tabs v-model:activeKey="activeKey">
      <a-tab-pane key="keywords" :tab="t('tab_keywords')">
        <KeywordsReply />
      </a-tab-pane>
      <a-tab-pane key="received" :tab="t('tab_received')">
        <ReceivedReply />
      </a-tab-pane>
    </a-tabs>
  </div>
  
</template>

<script setup>
import { ref } from 'vue'
import KeywordsReply from './keywords-reply.vue'
import ReceivedReply from './received-reply.vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const route = useRoute()
const { t } = useI18n('views.robot.robot-config.auto-reply.index')
const activeKey = ref(route.query.active_key || 'keywords')
</script>

<style lang="less" scoped>
.auto-reply-home {
  height: 100%;
  width: 100%;
  padding: 16px 24px;
  overflow-y: auto;
}
</style>