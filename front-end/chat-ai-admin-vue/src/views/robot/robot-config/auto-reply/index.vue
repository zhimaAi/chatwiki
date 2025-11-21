<template>
  <div class="auto-reply-home">
    <a-tabs v-model:activeKey="activeKey">
      <a-tab-pane key="keywords" tab="关键词回复">
        <KeywordsReply />
      </a-tab-pane>
      <a-tab-pane key="received" tab="收到消息回复">
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

const route = useRoute()
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