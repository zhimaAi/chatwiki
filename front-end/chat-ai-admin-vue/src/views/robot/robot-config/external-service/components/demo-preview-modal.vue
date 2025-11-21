<template>
  <div>
    <a-modal v-model:open="open" :width="746" title="未认证公众号回复效果示例" :footer="null">
      <div class="content-block">
        <div class="img-wrapper-box">
          <div class="title">用户咨询自动回复效果示例</div>
          <img
            src="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/wechat-app/guide1.png"
            alt=""
          />
        </div>
        <div class="img-wrapper-box">
          <div class="title">回复内容超长自动分段效果示例</div>
          <img
            src="https://xkf-upload-oss.xiaokefu.com.cn/static/chat-wiki/wechat-app/guide2.png"
            alt=""
          />
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const open = ref(false)

const show = () => {
  open.value = true
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.content-block {
  display: flex;
  gap: 70px;
  padding: 24px;
  .img-wrapper-box {
    flex: 1;
    text-align: center;
    .title {
      text-align: center;
      font-size: 14px;
      line-height: 22px;
      color: #000000;
      margin-bottom: 8px;
    }
    img {
      width: 240px;
      height: 522px;
      border: 1px solid var(--06, #d9d9d9);
      border-radius: 16px;
    }
  }
}
</style>
