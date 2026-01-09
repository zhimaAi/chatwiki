<template>
  <div class="_container">
    <div class="main">
      <div class="title"><a>应用收费</a> <span>让用户为价值买单！</span></div>
      <div class="content">
        <div class="item-box">
          <div class="tit">设置收费规则</div>
          <div class="desc">支持设置按次收费/按天收费两种收费模式</div>
          <img src="@/assets/img/robot/app-charging/app-charging-img01.png"/>
        </div>
        <div class="item-box">
          <div class="tit">用户购买授权码</div>
          <div class="desc">用户扫码添加管理员购买授权码</div>
          <img src="@/assets/img/robot/app-charging/app-charging-img02.png"/>
        </div>
        <div class="item-box">
          <div class="tit">用户兑换使用</div>
          <div class="desc">用户回复授权码，获取使用次数</div>
          <img src="@/assets/img/robot/app-charging/app-charging-img03.png"/>
        </div>
      </div>
      <a-button type="primary" class="main-btn" @click="linkSetting">立即设置</a-button>
    </div>
  </div>
</template>

<script setup>
import {useRoute, useRouter} from 'vue-router';

const route = useRoute()
const router = useRouter()

function linkSetting() {
  router.push({
    path: '/robot/ability/payment',
    query: {
      ...route.query,
      form_guide: 1
    }
  })
}
</script>

<style scoped lang="less">
._container {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 86vh;

  .main {
    display: flex;
    align-items: center;
    flex-direction: column;

    .title {
      display: flex;
      align-items: center;
      gap: 8px;
      color: #262626;
      font-size: 36px;
      font-weight: 600;
      margin-bottom: 56px;
    }

    .content {
      display: flex;
      gap: 16px;
      text-align: center;

      .item-box {
        position: relative;
        width: 242px;
        &:before {
          content: "";
          display: block;
          width: 12px;
          height: 12px;
          border-radius: 24px;
          background: #2475FC;
          position: absolute;
          left: 0;
          right: 0;
          margin: auto;
          top: -24px;
          z-index: 99;
        }
        &:nth-child(2):after {
          content: "";
          display: block;
          width: 513px;
          height: 2px;
          background: #ADCCFF;
          position: absolute;
          left: -136px;
          top: -20px;
          z-index: 88;
        }
      }

      .tit {
        color: #000000;
        font-size: 14px;
        font-weight: 600;
        margin-bottom: 8px;
      }

      .desc {
        height: 43px;
        color: #595959;
        font-size: 14px;
      }

      img {
        width: 242px;
        height: 232px;
      }
    }

    .main-btn {
      width: 240px;
      margin-top: 48px;
    }
  }
}
</style>
