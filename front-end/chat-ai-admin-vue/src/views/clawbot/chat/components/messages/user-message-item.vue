<template>
  <div
    class="message-item user-message"
    :id="'msg-' + item.uid"
    :data-msg_type="item.msg_type"
  >
    <div class="itme-left">
      <img class="user-avatar" :src="item.avatar" />
    </div>
    <div class="itme-right">
      <div class="item-body">
        <template v-if="item.msg_type == 3">
          <div class="message-content">
            <img v-viewer class="msg-img" :src="item.media_id_to_oss_url" alt="" />
          </div>
        </template>
        <template v-else-if="item.msg_type == 99">
          <div class="message-content">
            <MultipleMessage :message="item.content" />
          </div>
        </template>
        <template v-else>
          <TextMessage class="message-content" :message="item.content" />
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import TextMessage from './text-message.vue'
import MultipleMessage from './multiple-message.vue'

defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
})
</script>

<style lang="less" scoped>
.user-message {
  .item-body {
    padding: 5px 0;
  }

  .message-content {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #3a4559;
    white-space: pre-wrap;
  }

  .msg-img {
    width: 100%;
    height: 100%;
    max-width: 300px;
    max-height: 300px;
  }
}
</style>
