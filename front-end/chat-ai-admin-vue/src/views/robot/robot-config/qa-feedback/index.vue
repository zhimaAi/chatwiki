<style lang="less" scoped>
.user-model-page {
  width: 100%;
  height: 100%;
  border-bottom: 1px solid #fff;
  border-right: 1px solid #fff;
  background-color: #f2f4f7;

  .page-title {
    display: flex;
    align-items: center;
    gap: 24px;
    padding: 24px 24px 16px;
    background-color: #fff;
    color: #000000;
    font-family: "PingFang SC";
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
  }

  .list-wrapper {
    background: #fff;
    height: calc(100% - 178px);
  }

  .overview-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0px 24px;
    gap: 16px;
    background-color: white;

    .item-box {
      display: flex;
      align-items: center;
      justify-content: space-between;
      flex: 1;
      padding: 12px 16px;
      border-radius: 6px;
      background: var(--09, #F2F4F7);
      

      .item-info {
        color: #7a8699;
        font-family: "PingFang SC";
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
      }

      .item-icon-content {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 4px;
        flex: 1 0 0;
        font-weight: bolder;
      }

      .item-like {
        border-right: 1px solid #D8DDE5;
        height: 12px;
      }
    }
  }
}
</style>

<template>
  <div class="user-model-page">
    <div class="page-title">数据概览</div>
    <div class="overview-box">
      <div class="item-box" v-for="item in overviewData" :key="item.titleObj.id">
        <div class="item-info">{{ item.titleObj.title }}</div>
        <div class="item-icon-content item-like"><svg-icon style="font-size: 24px; color: #8C8C8C;" name="like" />{{ item.like_count }}</div>
        <div class="item-icon-content item-dislike"><svg-icon style="font-size: 24px; color: #8C8C8C;" name="dislike" />{{ item.dislike }}</div>
      </div>
    </div>

    <div class="page-title">反馈记录</div>
    <div class="list-wrapper">
      <QaFeedback></QaFeedback>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { ref, onMounted } from 'vue'
import { getFeedbackStats } from '@/api/manage/index.js'
import QaFeedback from './qa-feedback.vue'

const route = useRoute()

const query = route.query

const overviewData = ref([])

const forMatTitle = (key) => {
  let obj = {
    title: '',
    id: 0
  }
  switch (key) {
    case 'today_stats':
        obj.title = '今日'
        obj.id = 1
        break;
    case 'yesterday_stats':
        obj.title = '昨日'
        obj.id = 2
        break;
    case 'week_stats':
        obj.title = '近7日'
        obj.id = 3
        break;
    case 'total_stats':
        obj.title = '总计'
        obj.id = 4
        break;
  }
  return obj
}
const getData = () => {
  // 获取列表
  let parmas = {
    robot_id: query.id, // 机器人ID, // 机器人ID
  }


  getFeedbackStats(parmas).then((res) => {
    const arr = []
    for (const key in res.data) {
      const item = res.data[key]
      arr.push({
        key: key,
        dislike: item.dislike,
        like_count: item.like_count,
        titleObj: forMatTitle(key)
      })
    }
    overviewData.value = arr.sort((a ,b) => a.titleObj.id - b.titleObj.id)
  })
}

onMounted(() => {
  getData()
})

</script>
