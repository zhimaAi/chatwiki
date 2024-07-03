<template>
  <div class="recall-testing-box">
    <div class="recall-from-box">
      <cu-scroll>
        <RrcallTestingForm @load="handleLoading" @save="handleRecallTest"></RrcallTestingForm>
      </cu-scroll>
    </div>
    <div class="content-list-box">
      <div class="empty-box" v-if="isEmpty || lists.length === 0">
        <img src="@/assets/img/library/detail/empty.png" alt="" />
        <div class="title">暂无结果，请重试</div>
      </div>
      <cu-scroll v-else>
        <div class="list-item" v-for="item in lists" :key="item.id">
          <div class="top-block">
            <div class="title">
              id：{{ item.id }}
              <div v-if="item.title" class="ml4">{{ item.title }}</div>
              <span>共{{ item.word_total }}个字符</span>
            </div>
          </div>
          <div class="info-block">
            <span class="gray-text">来自：</span>
            <div class="link-text" @click="handlePreview(item)">
              <LinkOutlined />
              {{ item.file_name }}
            </div>
            <span class="v-line"></span>
            <span class="gray-text"
              >相似度：
              <svg-icon name="similarity" style="font-size: 16px"></svg-icon>
              {{ item.similarity }}
            </span>
          </div>
          <div class="content-box" v-if="item.question">Q：{{ item.question }}</div>
          <div class="content-box" v-if="item.answer">A：{{ item.answer }}</div>
          <div class="content-box" v-html="item.content"></div>
        </div>
      </cu-scroll>
      <div v-if="loading" class="loading-box"><a-spin /></div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { libraryRecallTest } from '@/api/library'
import { LinkOutlined } from '@ant-design/icons-vue'
import RrcallTestingForm from './components/recall-testing-form.vue'
const route = useRoute()
const router = useRouter()
const isEmpty = ref(false)
const loading = ref(false)
const lists = ref([])
const handleRecallTest = (data) => {
  loading.value = false;
  lists.value = data || []
}
const handleLoading = () => {
  loading.value = true;
}
const handlePreview = (record) => {
  window.open(`/#/library/preview?id=${record.file_id}`)
}
</script>
<style lang="less" scoped>
.recall-testing-box {
  background: #f2f4f7;
  height: 100%;
  display: flex;
  padding: 16px;
  padding-bottom: 0;
  .recall-from-box {
    width: 368px;
    height: 100%;
    border-radius: 2px;
    box-shadow: 0 2px 16px 0 #0e1b3a29;
    background: #fff;
    padding: 16px 0;
  }
  .content-list-box {
    flex: 1;
    height: 100%;
    overflow: hidden;
    background: #f2f4f7;
    margin-left: 16px;
    position: relative;
    .list-item {
      margin-top: 8px;
      width: 100%;
      background: #fff;
      border-radius: 2px;
      padding: 16px;
      .top-block {
        display: flex;
        align-items: center;
        justify-content: space-between;
        .title {
          display: flex;
          align-items: center;
          font-size: 14px;
          line-height: 22px;
          font-weight: 600;
          color: #000000;
          .mr4 {
            margin-right: 4px;
          }
          .ml4 {
            margin-left: 4px;
          }
          span {
            color: #8c8c8c;
            font-weight: 400;
            margin-left: 8px;
          }
        }
      }
      .info-block {
        display: flex;
        align-items: center;
        margin-top: 4px;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        .gray-text {
          color: #8c8c8c;
          display: flex;
          .svg-action {
            margin-right: 4px;
          }
        }
        .link-text {
          display: flex;
          align-items: center;
          color: #164799;
          cursor: pointer;
        }
        .v-line {
          width: 1px;
          height: 12px;
          background: #d9d9d9;
          margin: 0 8px;
        }
      }
      .content-box {
        color: #595959;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        margin-top: 8px;
        white-space: pre-wrap;
        word-wrap: break-word;
      }
    }
  }
  .loading-box {
    position: absolute;
    left: 50%;
    top: 40%;
  }
  .empty-box {
    text-align: center;
    height: 100%;
    padding-top: 148px;
    img {
      width: 200px;
      height: 200px;
    }
    .title {
      font-size: 16px;
      font-style: normal;
      font-weight: 600;
      line-height: 24px;
      color: #262626;
    }
  }
}
</style>
