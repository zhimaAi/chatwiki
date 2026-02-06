<template>
  <div class="star-select-box">
    <a-tabs v-model:activeKey="activeKey" size="small" @change="handleChange">
      <a-tab-pane :key="-1">
        <template #tab>{{ t('tab_all') }}</template>
      </a-tab-pane>
      <a-tab-pane :key="0" v-if="!hideId.includes(0)">
        <template #tab>
          <div class="star-item"><StarOutlined />{{ t('tab_unselected') }}</div>
        </template>
      </a-tab-pane>
      <a-tab-pane v-for="item in props.startLists" :key="item.id">
        <template #tab>
          <div class="star-item">
            <StarFilled :style="{ color: colorLists[item.type] }" /> {{ item.name || '-' }} <span v-if="item.data_count > 0">({{ item.data_count }})</span>
          </div>
        </template>
      </a-tab-pane>
    </a-tabs>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import { StarFilled, StarOutlined } from '@ant-design/icons-vue'
import colorLists from '@/utils/starColors.js'

const { t } = useI18n('views.library.library-preview.components.selct-star-box')
const emit = defineEmits(['change'])
const props = defineProps({
  startLists: {
    type: Array,
    default: () => []
  },
  hideId: {
    type: [Number, String, Array],
    default: () => []
  }
})

const activeKey = ref(-1)

const handleChange = () => {
  emit('change', activeKey.value)
}
</script>

<style lang="less" scoped>
.star-select-box {
  width: fit-content;
  background: #edeff2;
  border-radius: 6px;
  padding: 2px;
  .star-item {
    display: flex;
    align-items: center;
    gap: 4px;
    min-width: 32px;
    .anticon {
      margin: 0;
      font-size: 16px;
    }
  }
  ::v-deep(.ant-tabs-nav-operations) {
    display: none !important;
  }
  ::v-deep(.ant-tabs-nav) {
    margin-bottom: 0;
    &::before {
      border: 0;
    }
    .ant-tabs-tab {
      margin: 0;
      padding: 5px 16px;
    }
    .ant-tabs-ink-bar {
      height: 100%;
      background: #fff;
      border-radius: 6px;
    }
    .ant-tabs-tab-active {
      position: relative;
      z-index: 99;
    }
  }
}
</style>
