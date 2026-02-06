<template>
  <div class="privacy-policy">
    <div class="header">
      <div class="back-button" @click="goBack">
        <LeftOutlined />
        {{ t('btn_back') }}
      </div>
    </div>
    <h1>{{ t('title_privacy_policy') }}</h1>

    <p v-html="t('intro_p1')"></p>

    <p v-html="t('intro_p2')"></p>

    <p v-html="t('intro_p3')"></p>

    <p v-html="t('intro_p4')"></p>

    <p v-html="t('intro_p5')"></p>

    <ul>
      <li>{{ t('toc_item_1') }}</li>
      <li>{{ t('toc_item_2') }}</li>
      <li>{{ t('toc_item_3') }}</li>
      <li>{{ t('toc_item_4') }}</li>
      <li>{{ t('toc_item_5') }}</li>
      <li>{{ t('toc_item_6') }}</li>
      <li>{{ t('toc_item_7') }}</li>
      <li>{{ t('toc_item_8') }}</li>
      <li>{{ t('toc_item_9') }}</li>
      <li>{{ t('toc_item_10') }}</li>
      <li>{{ t('toc_item_11') }}</li>
    </ul>

    <h2 v-html="t('section1_title')"></h2>

    <h3 v-html="t('section1_subtitle')"></h3>

    <p v-html="t('section1_p1')"></p>

    <p v-html="t('section1_p2')"></p>

    <h3 v-html="t('section1_1_title')"></h3>

    <p v-html="t('section1_1_p1')"></p>

    <h3 v-html="t('section1_2_title')"></h3>

    <p v-html="t('section1_2_p1')"></p>

    <p v-html="t('section1_2_p2')"></p>

    <p v-html="t('section1_2_p3')"></p>

    <h3 v-html="t('section1_3_title')"></h3>

    <p v-html="t('section1_3_p1')"></p>

    <p v-html="t('section1_3_p2')"></p>

    <p v-html="t('section1_3_p3')"></p>

    <p v-html="t('section1_3_p4')"></p>

    <p v-html="t('section1_3_p5')"></p>

    <p v-html="t('section1_3_p6')"></p>

    <p v-html="t('section1_3_p7')"></p>

    <p v-html="t('section1_3_p8')"></p>

    <p v-html="t('section1_3_p9')"></p>

    <h3 v-html="t('section1_4_title')"></h3>

    <p v-html="t('section1_4_p1')"></p>

    <h3 v-html="t('section1_5_title')"></h3>

    <p v-html="t('section1_5_p1')"></p>

    <h3 v-html="t('section1_6_title')"></h3>

    <p v-html="t('section1_6_p1')"></p>

    <p v-html="t('section1_6_p2')"></p>

    <h3 v-html="t('section1_7_title')"></h3>

    <p v-html="t('section1_7_p1')"></p>

    <h3 v-html="t('section2_subtitle')"></h3>

    <p v-html="t('section2_third_party_p1')"></p>

    <h3 v-html="t('section3_subtitle')"></h3>

    <p v-html="t('section3_p1')"></p>

    <h3 v-html="t('section4_subtitle')"></h3>

    <p v-html="t('section4_p1')"></p>

    <p v-html="t('section4_1')"></p>

    <p v-html="t('section4_2')"></p>

    <p v-html="t('section4_3')"></p>

    <p v-html="t('section4_4')"></p>

    <p v-html="t('section4_4_dev_1')"></p>

    <p v-html="t('section4_4_dev_2')"></p>

    <p v-html="t('section4_4_dev_3')"></p>

    <p v-html="t('section4_4_dev_4')"></p>

    <p v-html="t('section4_5')"></p>

    <h2 v-html="t('section2_title')"></h2>

    <p v-html="t('section2_p1')"></p>

    <p v-html="t('section2_p2')"></p>

    <h2 v-html="t('section3_title')"></h2>

    <h3 v-html="t('section3_1_title')"></h3>

    <p v-html="t('section3_1_p1')"></p>

    <h3 v-html="t('section3_2_title')"></h3>

    <p v-html="t('section3_2_1')"></p>

    <p v-html="t('section3_2_1_1')"></p>

    <p v-html="t('section3_2_1_2')"></p>

    <p v-html="t('section3_2_1_3')"></p>

    <p v-html="t('section3_2_2')"></p>

    <h2 v-html="t('section4_title')"></h2>

    <h3 v-html="t('section4_1_title')"></h3>

    <p v-html="t('section4_1_1')"></p>

    <p v-html="t('section4_1_1_1')"></p>

    <p v-html="t('section4_1_1_2')"></p>

    <p v-html="t('section4_1_1_3')"></p>

    <p v-html="t('section4_1_2')"></p>

    <p v-html="t('section4_1_2_1')"></p>

    <p v-html="t('section4_1_3')"></p>

    <p v-html="t('section4_1_3_1')"></p>

    <p v-html="t('section4_1_4')"></p>

    <p v-html="t('section4_1_4_1')"></p>

    <h3 v-html="t('section4_2_title')"></h3>

    <p v-html="t('section4_2_1')"></p>

    <h3 v-html="t('section4_2_2_title')"></h3>

    <p v-html="t('section4_2_2_1')"></p>

    <p v-html="t('section4_2_2_2')"></p>

    <p v-html="t('section4_2_3')"></p>

    <h3 v-html="t('section4_3_title')"></h3>

    <p v-html="t('section4_3_1')"></p>

    <p v-html="t('section4_3_1_1')"></p>

    <p v-html="t('section4_3_1_2')"></p>

    <p v-html="t('section4_3_1_3')"></p>

    <h3 v-html="t('section4_4_title')"></h3>

    <p v-html="t('section4_4_1')"></p>

    <p v-html="t('section4_4_1_1')"></p>

    <p v-html="t('section4_4_1_2')"></p>

    <h3 v-html="t('section4_4_2_title')"></h3>

    <p v-html="t('section4_4_2_1')"></p>

    <p v-html="t('section4_4_2_2')"></p>

    <p v-html="t('section4_4_2_3')"></p>

    <p v-html="t('section4_4_2_4')"></p>

    <p v-html="t('section4_4_2_5')"></p>

    <p v-html="t('section4_4_2_6')"></p>

    <p v-html="t('section4_4_2_7')"></p>

    <p v-html="t('section4_4_2_8')"></p>

    <p v-html="t('section4_4_2_9')"></p>

    <p v-html="t('section4_4_2_10')"></p>

    <p v-html="t('section4_4_2_11')"></p>

    <p v-html="t('section4_4_2_12')"></p>

    <p v-html="t('section4_4_2_13')"></p>

    <p v-html="t('section4_4_3')"></p>

    <h3 v-html="t('section4_5_title')"></h3>

    <p v-html="t('section4_5_1')"></p>

    <h2 v-html="t('section5_title')"></h2>

    <p v-html="t('section5_1')"></p>

    <p v-html="t('section5_2')"></p>

    <p v-html="t('section5_3')"></p>

    <p v-html="t('section5_4')"></p>

    <p v-html="t('section5_5')"></p>

    <p v-html="t('section5_6')"></p>

    <h2 v-html="t('section6_title')"></h2>

    <p v-html="t('section6_1')"></p>

    <p v-html="t('section6_2')"></p>

    <p v-html="t('section6_3')"></p>

    <h2 v-html="t('section7_title')"></h2>

    <p v-html="t('section7_1')"></p>

    <h3 v-html="t('section7_2_title')"></h3>

    <p v-html="t('section7_2_1')"></p>

    <h3 v-html="t('section7_3_title')"></h3>

    <p v-html="t('section7_3_1')"></p>

    <p v-html="t('section7_3_2')"></p>

    <p v-html="t('section7_3_3')"></p>

    <h3 v-html="t('section7_4_title')"></h3>

    <p v-html="t('section7_4_1')"></p>

    <h3 v-html="t('section7_5_title')"></h3>

    <p v-html="t('section7_5_1')"></p>

    <h3 v-html="t('section7_6_title')"></h3>

    <p v-html="t('section7_6_1')"></p>

    <h3 v-html="t('section7_7_title')"></h3>

    <p v-html="t('section7_7_1')"></p>

    <p v-html="t('section7_7_2')"></p>

    <p v-html="t('section7_7_3')"></p>

    <p v-html="t('section7_7_4')"></p>

    <p v-html="t('section7_7_5')"></p>

    <p v-html="t('section7_7_6')"></p>

    <p v-html="t('section7_7_7')"></p>

    <p v-html="t('section7_7_8')"></p>

    <p v-html="t('section7_7_9')"></p>

    <p v-html="t('section7_7_10')"></p>

    <p v-html="t('section7_7_11')"></p>

    <h3 v-html="t('section7_8_title')"></h3>

    <p v-html="t('section7_8_1')"></p>

    <h2 v-html="t('section8_title')"></h2>

    <p v-html="t('section8_1')"></p>

    <h2 v-html="t('section9_title')"></h2>

    <p v-html="t('section9_1')"></p>

    <p v-html="t('section9_2')"></p>

    <p v-html="t('section9_3')"></p>

    <p v-html="t('section9_4')"></p>

    <p v-html="t('section9_4_1')"></p>

    <p v-html="t('section9_4_2')"></p>

    <p v-html="t('section9_4_3')"></p>

    <p v-html="t('section9_4_4')"></p>

    <p v-html="t('section9_4_5')"></p>

    <h2 v-html="t('section10_title')"></h2>

    <p v-html="t('section10_1')"></p>

    <p v-html="t('section10_2')"></p>

    <p v-html="t('section10_3')"></p>

    <h2 v-html="t('section11_title')"></h2>

    <p v-html="t('section11_1')"></p>
  </div>
</template>

<script setup>
import { LeftOutlined } from '@ant-design/icons-vue';
import { useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.privacy-policy.index')
const router = useRouter()

const goBack = () => {
  router.back()
}
</script>

<style lang="less" scoped>
.privacy-policy {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  line-height: 1.8;
  color: #333;


  .back-button {
    display: flex;
    align-items: center;
    width: fit-content;
    gap: 8px;
    padding: 4px 16px;
    background-color: #fff;
    border: 1px solid #d9d9d9;
    border-radius: 6px;
    color: #333;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      background-color: #f5f5f5;
      border-color: #4096ff;
      color: #4096ff;
    }

    &:active {
      background-color: #e6f7ff;
    }
  }

  h1 {
    font-size: 28px;
    font-weight: 600;
    margin-bottom: 24px;
    color: #1a1a1a;
  }

  h2 {
    font-size: 24px;
    font-weight: 600;
    margin: 32px 0 16px 0;
    color: #1a1a1a;
  }

  h3 {
    font-size: 20px;
    font-weight: 600;
    margin: 24px 0 12px 0;
    color: #333;
  }

  p {
    margin-bottom: 16px;
    text-align: justify;
  }

  ul {
    margin: 16px 0;
    padding-left: 24px;

    li {
      margin-bottom: 8px;
    }
  }

  strong {
    font-weight: 600;
    color: #1a1a1a;
  }
}
</style>
