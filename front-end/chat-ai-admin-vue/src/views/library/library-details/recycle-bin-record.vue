<template>
  <div class="details-library-page">
    <div class="page-body">
      <div class="tip-box">
        <div class="tip-box-left">
          <div class="tip-title">
            <ExclamationCircleFilled style="font-size: 20px; color: #FF9900;" />{{ t('tip_title') }}
          </div>
          <div class="tip-text">{{ t('tip_text') }}</div>
          <div>
            <a class="link-text" href="https://chatwiki.com/price" target="_blank">{{ t('link_text') }} <RightOutlined /></a>
          </div>
        </div>
         <div class="tip-box-right">
          <img src="@/assets/img/library/detail/lock.jpg" alt="">
         </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ExclamationCircleFilled, RightOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.library-details.recycle-bin-record')
</script>

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.page-body{
  width: 584px;
  height: 292px;
  padding: 46px 58px 46px 54px;
  border-radius: 6px;
  box-shadow: 0 4px 24px 0 rgba(0, 0, 0, 0.08);
  background: #FFF;
}
.tip-box{
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  .tip-box-left{
    flex: 1;
  }
  .tip-box-right{
    margin-left: 46px;
  }
  .tip-title{
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;
  }

  .tip-text{
    line-height: 22px;
    margin-bottom: 24px;
    font-size: 14px;
    color: #595959;
  }
  .link-text{
    font-size: 14px;
    line-height: 22px;
  }
}
</style>
