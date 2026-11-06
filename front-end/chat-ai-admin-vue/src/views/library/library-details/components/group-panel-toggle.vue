<template>
  <button
    type="button"
    class="group-panel-toggle"
    :class="{ 'no-transition': dragging }"
    :style="{ left: (collapsed ? 0 : width - 10) + 'px' }"
    :title="label"
    :aria-label="label"
    :aria-expanded="!collapsed"
    @click="$emit('toggle')"
  >
    <RightOutlined v-if="collapsed" />
    <LeftOutlined v-else />
  </button>
</template>

<script setup>
import { computed } from 'vue'
import { LeftOutlined, RightOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'

const props = defineProps({
  collapsed: { type: Boolean, default: false },
  width: { type: Number, required: true },
  dragging: { type: Boolean, default: false }
})
defineEmits(['toggle'])
const { t } = useI18n('views.library.library-list.library-list')
const label = computed(() => t(props.collapsed ? 'expand_groups' : 'collapse_groups'))
</script>

<style lang="less" scoped>
.group-panel-toggle {
  position: absolute;
  top: 50%;
  z-index: 20;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 44px;
  padding: 0;
  border: 1px solid #e8e8e8;
  border-radius: 6px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, .04);
  color: #8c8c8c;
  font-size: 12px;
  cursor: pointer;
  transform: translateY(-50%);
  transition: left .3s cubic-bezier(.65, 0, .35, 1);
  &.no-transition { transition: none; }
  &:hover { color: #2475fc; background: #f5f9ff; }
  &:focus-visible { outline: 2px solid #2475fc; outline-offset: 2px; }
}
</style>
