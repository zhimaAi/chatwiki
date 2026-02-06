<template>
  <div>
    <a-modal
      wrapClassName="no-padding-modal"
      :bodyStyle="{ 'padding-right': '24px', 'max-height': '650px', 'overflow-y': 'auto' }"
      v-model:open="open"
      :title="t('modal_title')"
      :width="746"
      :footer="null"
    >
      <div class="list-box">
        <div class="list-item" v-for="item in list" :key="item.id">
          <div class="error-tip">{{ item.split_errmsg }}</div>
          <div class="content-block">{{ item.content }}</div>
          <div class="fragment-img" v-viewer>
            <img v-for="(src, index) in item.images" :key="index" :src="src" alt="" />
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { getFAQFileChunks } from '@/api/library'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.ai-extract-faq.list.components.fail-detail')

const open = ref(false)

const list = ref([])
const show = (data) => {
  open.value = true
  list.value = []
  getFailList(data.id)
}

const getFailList = (id) => {
  getFAQFileChunks({ id }).then((res) => {
    let data = res.data || []
    list.value = data.map((item) => {
      return {
        ...item,
        images: item.images ? JSON.parse(item.images) : []
      }
    })
  })
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.list-box {
  display: flex;
  flex-direction: column;
  gap: 12px;
  .list-item {
    padding: 16px;
    border-radius: 6px;
    background-color: #f2f2f2;
    font-size: 14px;
    line-height: 22px;
    color: #595959;
    white-space: pre-line;
    .error-tip {
      color: #f40;
      margin-bottom: 8px;
    }
  }
}
.fragment-img {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  img {
    width: 80px;
    height: 80px;
    border-radius: 6px;
    cursor: pointer;
  }
}
</style>
