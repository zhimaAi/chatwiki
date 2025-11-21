<template>
  <div class="auto-reply-home">
    <ReceivedReply />
  </div>
  
</template>

<script setup>
import ReceivedReply from './received-reply.vue'

</script>

<style lang="less" scoped>
.auto-reply-home {
  height: 100%;
  width: 100%;
  padding: 16px;
  overflow-y: auto;
}
</style>