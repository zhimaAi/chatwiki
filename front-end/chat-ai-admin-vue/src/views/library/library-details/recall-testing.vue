<template>
  <div class="recall-testing-box">
    <div class="recall-from-box">
      <cu-scroll>
        <RrcallTestingForm @load="handleLoading" @save="handleRecallTest"></RrcallTestingForm>
      </cu-scroll>
    </div>
    <div class="content-list-box" ref="contRef">
      <div class="empty-box" v-if="isEmpty || lists.length === 0">
        <img src="@/assets/img/library/detail/empty.png" alt="" />
        <div class="title">{{ t('empty_no_results') }}</div>
      </div>
      <cu-scroll v-else>
        <div class="list-item" v-for="item in lists" :key="item.id">
          <div class="top-block">
            <div class="title">
              {{ t('label_id') }}{{ item.id }}
              <div v-if="item.title" class="ml4">{{ item.title }}</div>
              <span>{{ t('text_total_chars', { count: item.word_total }) }}</span>
            </div>
          </div>
          <div class="info-block">
            <span class="gray-text">{{ t('label_from') }}</span>
            <div class="link-text" @click="handlePreview(item)">
              <LinkOutlined />
              {{ item.file_name }}
              <span v-if="!item.file_name">{{ item.library_name }}{{ t('text_selected') }}</span>
            </div>
            <span class="v-line"></span>
            <span class="gray-text"
              >{{ t('label_similarity') }}
              <svg-icon name="similarity" style="font-size: 16px"></svg-icon>
              {{ item.similarity }}
            </span>
          </div>
          <template v-if="libraryType == 2">
            <div class="content-box" v-if="item.question">{{ t('label_question') }}{{ item.question }}</div>
            <div class="content-box similar-questions-box" v-if="item.similar_questions && item.similar_questions.length">
              <span>{{ t('label_similar_questions') }}</span>
              <ul class="similar-questions-list">
                <li v-for="(value, index) in item.similar_questions" :key="index">{{ value }}</li>
              </ul>
            </div>
            <div class="content-box" v-if="item.answer">{{ t('label_answer') }}{{ item.answer }}</div>
          </template>
          <div class="content-box" v-html="item.content" v-else></div>
        </div>
      </cu-scroll>
      <div v-if="loading" class="loading-box"><a-spin /></div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/hooks/web/useI18n'
import { useLibraryStore } from '@/stores/modules/library'
import { LinkOutlined } from '@ant-design/icons-vue'
import RrcallTestingForm from './components/recall-testing-form.vue'
import {useMathJax} from "@/composables/useMathJax.js";

const { t } = useI18n('views.library.library-details.recall-testing')

const route = useRoute()
const router = useRouter()
const libraryStore = useLibraryStore()
const { renderMath } = useMathJax()

const contRef = ref(null)
const isEmpty = ref(false)
const loading = ref(false)
const lists = ref([])


const libraryType = computed(() => libraryStore.type)

const handleRecallTest = (data = []) => {
  loading.value = false;

  data.forEach((item) => {
    if(item.similar_questions){
      item.similar_questions = JSON.parse(item.similar_questions)
    }
  })

  lists.value = data || []
  renderMath(contRef.value)
}
const handleLoading = () => {
  loading.value = true;
}
const handlePreview = (record) => {
  if(!record.file_name){
    window.open(`/#/library/details/categary-manage?id=${record.library_id}`)
    return
  }
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
      .similar-questions-box {
        display: flex;
        .similar-questions-list{
          list-style: none;
          margin: 0;
          padding-left: 0;
          flex: 1;
        }
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
