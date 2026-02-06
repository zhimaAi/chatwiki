<template>
  <div class="main-content-block">
    <a-tabs class="tab-wrapper" v-model:activeKey="activeKey" @change="handleChangeTab">
      <a-tab-pane :key="1" :tab="t('tab_default_library')"></a-tab-pane>
      <a-tab-pane :key="2" :tab="t('tab_related_library')"></a-tab-pane>
       <template #rightExtra>
        <div class="tab-right-extra">
          {{ t('text_hit_rate_label') }}：{{ library_hit_rate }}%
          <a @click="toHitStatics">{{ t('link_view') }}</a>
        </div>
       </template>
    </a-tabs>
    <div class="body-content-box" style="padding-right: 16px" v-if="activeKey == 1">
      <DefaultLibrary />
    </div>
    <div class="body-content-box" v-else>
      <cu-scroll>
        <RelatedLibrary />
      </cu-scroll>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import DefaultLibrary from './default-library.vue'
import RelatedLibrary from './related-library.vue'
import dayjs from 'dayjs'
import { statAiTipAnalyse } from '@/api/manage/index.js'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.library-config.index')
const router = useRouter()
const query = useRoute().query

const activeLocalKey = '/robot/config/library-config/activeKey'
const activeKey = ref(+localStorage.getItem(activeLocalKey) || 1)
const handleChangeTab = () => {
  localStorage.setItem(activeLocalKey, activeKey.value)
}

const library_hit_rate = ref(0)
const getStatistics = async () => {
  statAiTipAnalyse({ 
    robot_id: query.id,
    start_date: dayjs().subtract(6, 'day').format('YYYY-MM-DD'),
    end_date: dayjs().format('YYYY-MM-DD'),
   }).then((res) => {
    library_hit_rate.value = res.data.header.library_hit_rate
  })
}

const toHitStatics = () => {
  localStorage.setItem('/robot/config/statistical_analysis/activeKey', 3)
  router.push({
    path: '/robot/config/statistical_analysis',
    query:{
      id: query.id,
      robot_key: query.robot_key,
    }
  })
}

onMounted(() => {
  getStatistics()
})

</script>

<style lang="less" scoped>
.main-content-block {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  ::v-deep(.ant-tabs-nav-wrap) {
    padding-left: 24px;
  }
  .tab-wrapper ::v-deep(.ant-tabs-nav) {
    margin-bottom: 0;
  }
}
.body-content-box {
  flex: 1;
  overflow: hidden;
}
.tab-right-extra{
  display: flex;
  align-items: center;
  padding-right: 24px;
  gap: 4px;
}
</style>
