<template>
  <div class="_container">
    <div class="main">
      <div class="title"><a>{{ t('title_app_charging') }}</a> <span>{{ t('subtitle_value_payment') }}</span></div>
      <div class="content">
        <div class="item-box">
          <div class="tit">{{ t('title_set_pricing_rules') }}</div>
          <div class="desc">{{ t('desc_pricing_modes') }}</div>
          <img src="@/assets/img/robot/app-charging/app-charging-img01.png"/>
        </div>
        <div class="item-box">
          <div class="tit">{{ t('title_purchase_auth_code') }}</div>
          <div class="desc">{{ t('desc_scan_to_purchase') }}</div>
          <img src="@/assets/img/robot/app-charging/app-charging-img02.png"/>
        </div>
        <div class="item-box">
          <div class="tit">{{ t('title_redeem_usage') }}</div>
          <div class="desc">{{ t('desc_redeem_code') }}</div>
          <img src="@/assets/img/robot/app-charging/app-charging-img03.png"/>
        </div>
      </div>
      <a-button type="primary" class="main-btn" @click="linkSetting">{{ t('btn_setup_now') }}</a-button>
    </div>
  </div>
</template>

<script setup>
import {useRoute, useRouter} from 'vue-router';
import { useI18n } from '@/hooks/web/useI18n';

const { t } = useI18n('views.robot.robot-config.payment.guide')

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
