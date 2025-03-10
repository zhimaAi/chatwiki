<style lang="less" scoped>
.page-tabs {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 0 24px;
  font-size: 14px;
  color: #595959;
  background-color: #fff;

  .tab-item {
    position: relative;
    height: 46px;
    line-height: 46px;
    cursor: pointer;

    &.active {
      font-weight: 600;
      color: #2475fc;
    }

    &.active::before {
      content: '';
      display: block;
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      height: 2px;
      background-color: #2475fc;
    }
  }
}
</style>

<template>
  <div class="page-tabs">
    <div class="tab-item" :class="{ active: props.value == tab.value }" @click="handleClickTab(tab.value)"
      v-for="tab in tabs" :key="tab.value">
      {{ t(tab.langKey) }}
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n()

const emit = defineEmits(['update:value', 'change'])

const props = defineProps({
  value: {
    type: [Number, String]
  }
})

const tabs = ref([
  {
    value: 1,
    title: t('views.user.model.modelsHaveBeenAdded'),
    langKey: 'views.user.model.modelsHaveBeenAdded'
  },
  {
    value: 0,
    title: t('views.user.model.canAddAModel'),
    langKey: 'views.user.model.canAddAModel'
  }
])

const handleClickTab = (val) => {
  emit('update:value', val)
  emit('change', val)
}
</script>
