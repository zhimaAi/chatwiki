<template>
  <a-drawer
    :open="open"
    placement="right"
    :width="400"
    title="基础数据"
    @close="onClose"
  >
    <div class="base-data-drawer" v-if="record">
      <div class="base-cover">
        <img v-if="fileCover" :src="fileCover" alt="" />
      </div>
      <div class="base-title zm-line1">{{ record.file_name }}</div>
      <div class="base-update">
        <span>更新于 {{ record.update_time }}</span>
        <span class="sep">|</span>
        <span>
          <span v-if="record.doc_auto_renew_frequency == 1">不自动更新</span>
          <span v-if="record.doc_auto_renew_frequency == 2">每天</span>
          <span v-if="record.doc_auto_renew_frequency == 3">每3天</span>
          <span v-if="record.doc_auto_renew_frequency == 4">每7天</span>
          <span v-if="record.doc_auto_renew_frequency == 5">每30天</span>
          <span class="ml4" v-if="record.doc_auto_renew_frequency > 1 && record.doc_auto_renew_minute > 0">
            {{ convertTime(record.doc_auto_renew_minute) }}
          </span>
          <a class="ml12 btn-hover-block" @click="emit('editOnline', record)">修改</a>
        </span>
      </div>
      <div class="base-meta-grid">
        <div class="meta-item">
          <div class="value">{{ record.file_ext || '--' }}</div>
          <div class="label">文件格式</div>
        </div>
        <div class="line-vertical"></div>
        <div class="meta-item">
          <div class="value">{{ record.file_size_str || '--' }}</div>
          <div class="label">文档大小</div>
        </div>
        <div class="line-vertical"></div>
        <div class="meta-item">
          <div class="value">{{ record.paragraph_count ?? '-' }}</div>
          <div class="label">文档分段</div>
        </div>
      </div>
      <div class="base-hits-grid">
        <div class="hits-item">
          <div class="value">{{ record.total_hits ?? 0 }}</div>
          <div class="label">合计触发</div>
        </div>
        <div class="line-vertical"></div>
        <div class="hits-item">
          <div class="value">{{ record.today_hits ?? 0 }}</div>
          <div class="label">今日触发</div>
        </div>
        <div class="line-vertical"></div>
        <div class="hits-item">
          <div class="value">{{ record.yesterday_hits ?? 0 }}</div>
          <div class="label">昨日触发</div>
        </div>
      </div>
    </div>
  </a-drawer>
  </template>

<script setup>
import { computed } from 'vue'
import defaultCover from '@/assets/img/default-cover.png'
import emptyDocumentIcon from '@/assets/svg/empty-document.svg'
import { convertTime } from '@/utils/index'

const emit = defineEmits(['close', 'editOnline'])
const props = defineProps({
  open: {
    type: Boolean,
    default: false
  },
  record: {
    type: Object,
    default: null
  }
})

function onClose() {
  emit('close')
}

const getStatusClass = (item) => {
  const s = Number(item?.status ?? -1)
  if ([0, 1, 5, 6, 10].includes(s)) return 'loading'
  if ([3, 7, 8, 9].includes(s)) return 'error'
  return 'success'
}
const getFileCover = (rec) => {
  if (!rec) return ''
  const status = getStatusClass(rec)
  const src =
    rec.thumb_path ||
    rec.cover_url ||
    rec.cover ||
    rec.article_cover_url ||
    ''
  if (src) return src
  if (status === 'loading') return defaultCover
  if (status === 'error') return emptyDocumentIcon
  return ''
}
const fileCover = computed(() => getFileCover(props.record))
</script>

<style scoped lang="less">
.base-data-drawer {
  .base-cover {
    width: 350px;
    height: 200px;
    border-radius: 6px;
    border: 1px solid #F0F0F0;
    overflow: hidden;
    img { width: 100%; display: block; }
  }
  .base-title {
    margin-top: 10px;
    width: 352px;
    color: #262626;
    font-size: 14px;
    font-style: normal;
    font-weight: 600;
    line-height: 22px;
  }
  .base-update {
    margin-top: 4px;
    color: #8c8c8c;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
    .sep {
        display: inline-block;
        height: 15px;
        margin: 0 6px;
        overflow: hidden;
        color: #ddd;
     }
  }
  .base-meta-grid, .base-hits-grid {
    width: 352px;
    height: 80px;
    box-sizing: border-box;
    display: inline-flex;
    padding: 16px 15px 16px 24px;
    justify-content: space-between;
    align-items: center;
    border-radius: 12px;
    border: 1px solid #E4E6EB;
    background: #F2F4F7;
    margin-top: 16px;
    .meta-item, .hits-item {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      gap: 2px;
      background: #f5f5f5;
      border-radius: 8px;
      padding: 8px;
      .value {
        align-self: stretch;
        color: #262626;
        font-size: 16px;
        font-style: normal;
        font-weight: 600;
        line-height: 24px;
      }
      .label {
        align-self: stretch;
        color: #7a8699;
        font-size: 14px;
        font-style: normal;
        font-weight: 400;
        line-height: 22px;
      }
    }
  }
  .base-hits-grid {
    border: 1px solid #ABCAFC;
    background: #F0F7FF;
    .hits-item {
      background: #F0F7FF;
    }
  }
}
.line-vertical {
  width: 1px;
  height: 40px;
  background: #D8DDE5;
  margin: 0 16px;
}
.ml12 {
  margin-left: 12px;
}
</style>
