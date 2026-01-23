<template>
  <a-modal
    v-model:open="open"
    :footer="null"
    wrapClassName="qrcode-wrap"
    width="380px"
  >
    <div class="qrcode-box">
      <div class="title">扫码登录</div>
      <div class="content">
        <div class="tip">请使用微信扫描下方二维码登员</div>
        <img class="qrcode" src="https://xiaokefu.com.cn/static/home/img/zhima.jpg"/>
        <div class="desc">二妞码有效期为5分钟，请尽快扫描</div>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import {ref} from 'vue'

const open = ref(false)

function show() {
  open.value = true
}

defineExpose({
  show,
})
</script>

<style lang="less">
.qrcode-wrap.ant-modal-wrap {
  .ant-modal-content {
    padding: 0;
  }
}
</style>
<style scoped lang="less">
.qrcode-box {
  .title {
    color: #000000;
    font-size: 24px;
    font-weight: 600;
    padding: 32px 24px;
    position: relative;
    &::before {
      content: "";
      position: absolute;
      width: 100%;
      height: 100%;
      left: 0;
      top: -18px;
      background: url("@/assets/img/login-pop-bg.png") no-repeat;
      background-size: 100% 100%;
    }
    &::after {
      content: "扫码登录";
      position: absolute;
      color: #000000;
      font-size: 24px;
      font-weight: 600;
      left: 24px;
      top: 24px;
    }
  }
  .content {
    display: flex;
    align-items: center;
    flex-direction: column;
    gap: 12px;
    color: #262626;
    font-size: 14px;
    padding: 0 24px 24px;

    .qrcode {
      width: 170px;
      height: 170px;
    }

    .desc {
      color: #8c8c8c;
    }
  }
}
</style>
