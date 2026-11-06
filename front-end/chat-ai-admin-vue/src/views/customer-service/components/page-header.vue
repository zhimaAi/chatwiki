<template>
  <header class="customer-page-header">
    <div class="header-heading">
      <button v-if="back" type="button" class="back-button" :aria-label="backLabel" @click="emit('back')"><ArrowLeftOutlined /></button>
      <h1 :title="title">{{ title }}</h1>
      <div v-if="$slots.default" class="header-status"><slot /></div>
    </div>
    <div v-if="tabs.length" class="page-tabs" role="tablist" :aria-label="title">
      <button v-for="tab in tabs" :key="tab.key" :ref="element => setTabElement(tab.key, element)" type="button" role="tab" :class="{ active: activeKey === tab.key }" :aria-selected="activeKey === tab.key" :disabled="tab.disabled" :tabindex="tab.key === focusKey ? 0 : -1" @click="selectTab(tab)" @keydown="onTabKeydown($event, tab)">{{ tab.label }}</button>
    </div>
    <div v-if="$slots.extra" class="header-extra"><slot name="extra" /></div>
  </header>
</template>

<script setup>
import { computed, nextTick } from 'vue'
import { ArrowLeftOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  title: { type: String, required: true },
  tabs: { type: Array, default: () => [] },
  activeKey: { type: [String, Number], default: '' },
  back: { type: Boolean, default: false },
  backLabel: { type: String, default: '返回' }
})
const emit = defineEmits(['update:activeKey', 'change', 'back'])
const enabledTabs = computed(() => props.tabs.filter(tab => !tab.disabled))
const focusKey = computed(() => enabledTabs.value.some(tab => tab.key === props.activeKey) ? props.activeKey : enabledTabs.value[0]?.key)
const tabElements = new Map()
function setTabElement(key, element) {
  if (element) tabElements.set(key, element)
  else tabElements.delete(key)
}
function selectTab(tab) {
  if (tab.disabled || tab.key === props.activeKey) return
  emit('update:activeKey', tab.key)
  emit('change', tab.key)
}
async function onTabKeydown(event, tab) {
  if (!['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) return
  const tabs = enabledTabs.value
  if (!tabs.length) return
  event.preventDefault()
  const index = tabs.findIndex(item => item.key === tab.key)
  const next = event.key === 'Home' ? 0 : event.key === 'End' ? tabs.length - 1 : (index + (event.key === 'ArrowRight' ? 1 : -1) + tabs.length) % tabs.length
  selectTab(tabs[next])
  await nextTick()
  const element = tabElements.get(tabs[next].key)
  element?.focus()
  element?.scrollIntoView?.({ block: 'nearest', inline: 'nearest' })
}
</script>

<style lang="less" scoped>
.customer-page-header { display: flex; align-items: center; flex-shrink: 0; gap: 32px; box-sizing: border-box; height: 55px; min-width: 0; padding: 0 24px; border-bottom: 1px solid #e8e8e8; background: #fff; }
.header-heading { display: flex; align-items: center; gap: 16px; min-width: 0; max-width: 100%; }
h1 { min-width: 0; margin: 0; overflow: hidden; color: #080808; font-size: 18px; font-weight: 500; line-height: 32px; text-overflow: ellipsis; white-space: nowrap; }
.header-status, .header-extra { display: flex; align-items: center; flex-shrink: 0; gap: 8px; }
.header-extra { margin-left: auto; }
.back-button { display: inline-flex; align-items: center; justify-content: center; flex-shrink: 0; width: 32px; height: 32px; padding: 0; border: 0; border-radius: 6px; color: #5a5a5a; background: transparent; font-size: 16px; cursor: pointer; }
.back-button:hover, .back-button:focus-visible { color: #146ef5; background: #f1f6fe; }
.back-button:focus-visible { outline: 2px solid #146ef5; outline-offset: -2px; }
.page-tabs { display: flex; align-self: stretch; flex: 1; gap: 24px; min-width: 0; overflow-x: auto; scrollbar-width: none; &::-webkit-scrollbar { display: none; } }
.page-tabs button { display: inline-flex; align-items: center; flex-shrink: 0; height: 100%; padding: 0 4px; border: 0; border-radius: 0; outline: none; color: #5a5a5a; background: transparent; font: inherit; font-size: 14px; line-height: 22px; white-space: nowrap; cursor: pointer; }
.page-tabs button.active { color: #146ef5; box-shadow: inset 0 -2px #146ef5; }
.page-tabs button:hover:not(:disabled) { color: #146ef5; }
.page-tabs button:focus-visible { color: #146ef5; background: #f1f6fe; box-shadow: inset 0 -2px #146ef5; }
.page-tabs button:disabled { color: #ababab; cursor: not-allowed; }
@media (max-width: 500px) { .customer-page-header { gap: 16px; padding: 0 16px; } .page-tabs { gap: 16px; } }
</style>
