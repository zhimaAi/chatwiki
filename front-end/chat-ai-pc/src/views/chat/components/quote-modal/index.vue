<template>
  <div>
    <van-popup v-model:show="show" closeable round>
      <div class="modal-box">
        <div class="title-block">{{ baseParmas.file_name }}</div>
        <div class="content-block">
          <div class="list-item" v-for="(item, index) in lists" :key="item.id">
            <div class="item-title">
              {{ t('title_reference', { index: index + 1 }) }}
              <a v-if="item.page_num > 0" @click="viewSourceFile(item)">{{ t('btn_view_source') }}</a>
            </div>
            <div class="content-body">
              <div class="content-text" v-html="item.content"></div>
              <div class="content-text">{{item.question}}</div>
              <div class="content-text">{{item.answer}}</div>
              <div class="img-box">
                <img v-for="img in item.images" :src="img" alt="" v-viewer />
              </div>
            </div>
          </div>
        </div>
      </div>
    </van-popup>
    <van-popup v-model:show="viewOpen" closeable round>
      <div class="modal-box">
        <div class="title-block">{{ baseParmas.file_name }}</div>
        <div class="content-block">
          <VuePdfEmbed :source="sourceUrl"  />
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getAnswerSource } from '@/api/chat/index.js'
import VuePdfEmbed from 'vue-pdf-embed'
import { useChatStore } from '@/stores/modules/chat'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.chat.components.quote-modal.index')
const chatStore = useChatStore()
const { user } = chatStore
const show = ref(false)
const baseParmas = ref({})
const lists = ref([])
const showPopup = (data) => {
  baseParmas.value = {
    ...data
  }
  lists.value = []
  show.value = true
  if(data.answer_source_data && data.answer_source_data.length){
    lists.value = data.answer_source_data
    return
  }
  getLists()

}

const getLists = () => {
  getAnswerSource(baseParmas.value)
    .then((res) => {
      lists.value = res.data || []
    })
    .catch(() => {
      lists.value = []
    })
}

const sourceUrl = ref('')
const viewOpen = ref(false)
const viewSourceFile = (item) => {
  viewOpen.value = true;
  sourceUrl.value = '/manage/getLibRawFileOnePage?id=' + baseParmas.value.file_id + '&page=' + item.page_num + '&admin_user_id=' + user.admin_user_id
}

defineExpose({
  showPopup
})
</script>

<style lang="less" scoped>
.modal-box {
  width: 80vw;
  border-radius: 8px;
  max-width: 520px;
}
.title-block {
  padding: 16px;
  padding-right: 40px;
  font-size: 15px;
  font-weight: 500;
  color: #333;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  border-bottom: 1px solid #e8e8e8;
}
.content-block {
  min-height: 300px;
  padding: 16px;
  max-height: 500px;
  overflow: auto;
  .list-item {
    margin-bottom: 24px;
  }
  .item-title {
    display: flex;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 12px;
    font-size: 14px;
    color: #000;
    span {
      color: #8c8c8c;
      font-weight: 400;
    }
    a{
      color:#2475fc;
      cursor: pointer;
    }
  }
  .content-body {
    background: #f0f5fc;
    padding: 8px;
    line-height: 22px;
    margin-top: 8px;
    font-size: 13px;
    color: #595959;
    white-space: pre-wrap;
  }
  .img-box {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 8px;
    img {
      width: 80px;
      height: 80px;
      border-radius: 6px;
      cursor: pointer;
    }
  }
}
</style>
