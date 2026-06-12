<style lang="less" scoped>
.chat-list-box {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  
  background: #fff;
  box-sizing: border-box;
}

.chat-scroll {
  flex: 1;
  min-height: 0;
}

.chat-list-groups {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0 8px 16px;
}

.chat-list-group {
  display: flex;
  flex-direction: column;
}

.chat-list-group.is-compact {
  gap: 8px;
}

.chat-list-group:not(.is-compact) {
  gap: 8px;
}

.chat-list-title {
  height: 34px;
  padding: 12px 8px 0;
  color: #8c8c8c;
  font-size: 14px;
  line-height: 22px;
  font-weight: 400;
  box-sizing: border-box;
}

.chat-list-item {
  width: 100%;
  min-height: 38px;
  padding: 8px;
  border: 0;
  border-radius: 6px;
  background: #fff;
  display: flex;
  align-items: center;
  gap: 8px;
  text-align: left;
  cursor: pointer;
  box-sizing: border-box;
  transition:
    background-color 0.2s ease,
    color 0.2s ease;

  &:hover {
    background: #f7f7f7;
  }

  &.active {
    border-radius: 8px;
    background: #f2f6ff;

    .chat-item-icon,
    .chat-item-title {
      color: #1668dc;
    }
  }
}

.chat-item-icon {
  flex-shrink: 0;
  color: #262626;
  font-size: 16px;
  line-height: 1;
}

.chat-item-title {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #262626;
  font-size: 14px;
  line-height: 22px;
  font-weight: 400;
}
</style>

<template>
  <div class="chat-list-box">
    <cu-scroll ref="scroller" class="chat-scroll" @onScrollEnd="onScrollEnd">
      <div class="chat-list-groups">
        <section
          v-for="(group, index) in groupedList"
          :key="group.label"
          class="chat-list-group"
          :class="{ 'is-compact': index > 0 }"
        >
          <div class="chat-list-title">{{ group.label }}</div>
          <button
            v-for="item in group.items"
            :key="item.id"
            type="button"
            class="chat-list-item"
            :class="{ active: isActiveItem(item) }"
            @click="handleOpenChat(item)"
          >
            <svg-icon class="chat-item-icon" name="message2" size="16" />
            <span class="chat-item-title">{{ item.subject || item.last_chat_message }}</span>
          </button>
        </section>
      </div>
    </cu-scroll>
  </div>
</template>

<script setup>
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import CuScroll from '@/components/cu-scroll/cu-scroll.vue'

const { t } = useI18n('views.clawbot.chat.components.chat-list')
const emit = defineEmits(['openChat', 'onScrollEnd'])

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  },
  active: {
    type: [String, Number],
    default: ''
  }
})

const scroller = ref(null)

const parseChatDate = (item) => {
  const rawValue = item.update_time || item.updated_at || item.create_time || item.created_at || item.time

  if (!rawValue && rawValue !== 0) {
    return null
  }

  const normalizedValue = String(rawValue).trim()

  if (/^\d{10}$/.test(normalizedValue)) {
    return new Date(Number(normalizedValue) * 1000)
  }

  if (/^\d{13}$/.test(normalizedValue)) {
    return new Date(Number(normalizedValue))
  }

  return new Date(normalizedValue.replace(/-/g, '/'))
}

const isSameDate = (a, b) => {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate()
}

const getGroupLabel = (item) => {
  const date = parseChatDate(item)

  if (!date || Number.isNaN(date.getTime())) {
    return t('group_earlier')
  }

  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const itemDay = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const yesterday = new Date(today)
  yesterday.setDate(today.getDate() - 1)

  if (isSameDate(itemDay, today)) {
    return t('group_today')
  }

  if (isSameDate(itemDay, yesterday)) {
    return t('group_yesterday')
  }

  if (date.getFullYear() === now.getFullYear() && date.getMonth() === now.getMonth()) {
    return t('group_this_month')
  }

  const month = String(date.getMonth() + 1).padStart(2, '0')
  return `${date.getFullYear()}-${month}`
}

const groupedList = computed(() => {
  const groups = []
  const groupMap = new Map()

  props.list.forEach((item) => {
    const label = getGroupLabel(item)

    if (!groupMap.has(label)) {
      const group = { label, items: [] }
      groupMap.set(label, group)
      groups.push(group)
    }

    groupMap.get(label).items.push(item)
  })

  return groups
})

const isActiveItem = (item) => {
  const itemDialogueId = item.dialogue_id ?? item.id
  return String(itemDialogueId) === String(props.active)
}

const handleOpenChat = (item) => {
  emit('openChat', item)
}

const onScrollEnd = () => {
  emit('onScrollEnd')
}

watch(
  () => props.list,
  () => {
    nextTick(() => {
      scroller.value?.refresh()
    })
  }
)
</script>
