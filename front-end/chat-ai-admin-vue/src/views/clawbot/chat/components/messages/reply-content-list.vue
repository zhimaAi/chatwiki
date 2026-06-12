<template>
  <div class="reply-list" v-if="replyItems.length">
    <template v-for="(rc, idx) in replyItems" :key="idx">
      <div v-if="getReplyType(rc) === 'text'" class="reply-item reply-text">
        <div class="message-content" v-html="rc.description"></div>
      </div>
      <div v-else-if="getReplyType(rc) === 'image'" class="reply-item reply-image">
        <div class="message-content">
          <img v-viewer class="msg-img" :src="rc.pic || rc.thumb_url" />
        </div>
      </div>
      <div v-else-if="getReplyType(rc) === 'url'" class="reply-item reply-url">
        <div class="url-row">
          <a class="url-link" :href="rc.url" target="_blank">{{ rc.url }}</a>
        </div>
      </div>
      <div v-else-if="getReplyType(rc) === 'card'" class="reply-item reply-card">
        <div class="card-row">
          <img v-if="rc.thumb_url" :src="rc.thumb_url" class="card-thumb" />
          <div class="card-title-box">
            <svg-icon class="think-icon" name="applet"></svg-icon>
            <span class="card-title">{{ rc.title }}</span>
          </div>
        </div>
      </div>
      <div v-else-if="getReplyType(rc) === 'imageText'" class="reply-item reply-imageText">
        <a class="imageText-row" :href="rc.url" target="_blank">
          <img v-if="rc.thumb_url" :src="rc.thumb_url" class="imageText-thumb" />
          <div class="imageText-text">
            <div class="imageText-title">{{ rc.title }}</div>
            <div class="imageText-desc">{{ rc.description }}</div>
          </div>
        </a>
      </div>
      <div v-else-if="getReplyType(rc) === 'smartMenu'" class="reply-item reply-smartMenu">
        <div class="smart-menu-box">
          <div class="card-title" v-if="rc.smart_menu && rc.smart_menu.menu_description">
            {{ rc.smart_menu.menu_description }}
          </div>
          <div class="card-text">
            <template v-for="(line, li) in buildMenuLines(rc.smart_menu?.menu_content || [])" :key="li">
              <div class="reply-line" @click="onSmartReplyLineClick($event)">
                <span class="line-text" v-if="line.kind === 'text'">{{ line.text }}</span>
                <div v-else-if="line.kind === 'newline'" class="empty-line"></div>
                <span v-else-if="line.kind === 'html'" v-html="line.html"></span>
                <a
                  v-else-if="line.kind === 'keyword'"
                  href="javascript:;"
                  class="link"
                  @click.prevent="onClickSmartMenuKeyword(line.text)"
                >
                  <div v-if="line.serial_no">{{ line.serial_no }}</div>
                  {{ line.text }}
                </a>
              </div>
            </template>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const emit = defineEmits(['clickMsgMeun'])

const props = defineProps({
  replyContentList: {
    type: [Array, String],
    default: () => []
  }
})

const parseReplyList = (val) => {
  try {
    if (!val) return []
    if (Array.isArray(val)) return val
    if (typeof val === 'string') return JSON.parse(val || '[]')
    return []
  } catch (_e) {
    return []
  }
}

const replyItems = computed(() => parseReplyList(props.replyContentList))

const getReplyType = (item) => {
  return item?.reply_type || item?.type || ''
}

function buildMenuLines(menu_content) {
  const out = []
  ;(Array.isArray(menu_content) ? menu_content : []).forEach((mc) => {
    const menuType = String(mc?.menu_type || '')
    const text = String(mc?.content || '')
    if (menuType === '0') {
      if (text === '') {
        out.push({ kind: 'newline' })
      } else if (/<a[\s\S]*?<\/a>/.test(text)) {
        const sanitized = /href\s*=\s*['"]\s*#\s*['"]/i.test(text)
          ? text.replace(/href\s*=\s*['"]\s*#\s*['"]/ig, 'href="javascript:;"')
          : text.replace(/href=/ig, 'target="_blank" href=')
        out.push({ kind: 'html', html: sanitized })
      } else {
        out.push({ kind: 'text', text })
      }
    } else if (menuType === '1') {
      out.push({ kind: 'keyword', text, serial_no: mc?.serial_no || '' })
    }
  })
  return out.slice(0, 20)
}

function onSmartReplyLineClick(e) {
  const anchor = e.target?.closest?.('a')
  if (!anchor) {
    return
  }
  const href = String(anchor.getAttribute('href') || '')
  if (href === '#' || href === 'javascript:;') {
    e.preventDefault()
    e.stopPropagation()
  }
}

function onClickSmartMenuKeyword(text) {
  const keyword = String(text || '').trim()
  if (!keyword) {
    return
  }
  emit('clickMsgMeun', keyword)
}
</script>

<style lang="less" scoped>
.reply-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 8px;
}

.reply-item {
  border-radius: 4px 16px 16px 16px;
  padding: 16px 12px;
  background: #fff;
}

.message-content {
  line-height: 22px;
  font-size: 14px;
  font-weight: 400;
  color: #3a4559;
  white-space: pre-wrap;
}

.msg-img {
  width: 100%;
  height: 100%;
  max-width: 300px;
  max-height: 300px;
}

.reply-image .msg-img {
  max-width: 500px;
}

.reply-url .url-link {
  color: #2475fc;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  text-decoration: none;
}

.reply-url .url-row {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.reply-card .card-row {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.reply-card .card-title-box {
  display: flex;
  align-items: center;
  gap: 4px;
}

.reply-card .card-thumb {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.reply-card .card-title {
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.reply-imageText .imageText-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.reply-imageText .imageText-thumb {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: cover;
}

.reply-imageText .imageText-text {
  display: flex;
  flex-direction: column;
}

.reply-imageText .imageText-title {
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.reply-imageText .imageText-desc {
  color: #8c8c8c;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: 16px;
}

.reply-smartMenu {
  color: #1a1a1a;

  .smart-menu-box {
    .card-title {
      white-space: pre-wrap;
      align-self: stretch;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 20px;
      margin-bottom: 14px;
    }

    .card-text {
      display: flex;
      flex-direction: column;
      gap: 10px;

      .reply-line {
        line-height: 22px;

        .line-text {
          color: #3a4559;
        }

        .empty-line {
          height: 22px;
        }

        .link {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }
}
</style>
