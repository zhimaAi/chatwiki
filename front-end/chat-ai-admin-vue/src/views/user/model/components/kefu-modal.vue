<template>
  <div>
    <a-modal v-model:open="open" title="购买资源包" :footer="null" :width="546">
      <div class="modal-box">
        <a-alert show-icon message="暂不支持线上购买，请扫码联系客服购买"></a-alert>
        <div class="img-box">
          <img src="https://xkf-upload-oss.xiaokefu.com.cn/static/kefu.png" alt="" />
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
.modal-box {
  margin-top: 16px;
  .img-box {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-top: 16px;
    margin-bottom: 48px;
  }
  img {
    width: 200px;
  }
}
</style>
