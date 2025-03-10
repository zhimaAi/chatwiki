<style lang="less">
.custom-vue-node {
  width: 256px;
  padding: 12px 16px 16px 16px;
  border-radius: 8px;
  background: #fff;

  .cu-node-title {
    line-height: 24px;
    margin-bottom: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .delayed-text {
    line-height: 22px;
    margin-bottom: 4px;
    font-size: 14px;
    font-weight: 400;
    color: #8c8c8c;
  }
  .q-title {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #262626;
  }
  .q-btn-box {
    .q-btn {
      height: 38px;
      line-height: 38px;
      padding: 0 12px;
      margin-top: 12px;
      border-radius: 4px;
      font-size: 14px;
      color: #595959;
      background-color: #f2f4f7;
    }
  }
}
</style>

<template>
  <div class="custom-vue-node">
    <div class="cu-node-title">发送消息</div>
    <div class="delayed-text">延时20s</div>
    <div class="q-title">您好，你是想要了解墙面翻新吗？</div>
    <div>
      <div class="q-btn-box">
        <div class="q-btn">是的没错</div>
        <div class="q-btn">不是不是你误会了</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject } from 'vue'

// const getNode = inject('getNode')
// const nodeModel = getNode()
// const properties = nodeModel.properties
</script>
