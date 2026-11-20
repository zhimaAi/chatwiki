<template>
  <a-tooltip v-if="record.status == 3 || record.status == 8" placement="top" :trigger="['hover', 'focus']">
    <template #title>
      <div>{{ t('failure_reason') }}{{ record.failure_reason || t('failure_unknown') }}</div>
      <div>{{ t('failure_suggestion') }}{{ record.failure_suggestion || t('failure_unknown_suggestion') }}</div>
    </template>
    <slot />
  </a-tooltip>
  <slot v-else />
</template>

<script setup>
import { useI18n } from '@/hooks/web/useI18n'

defineProps({ record: { type: Object, required: true } })
const { t } = useI18n('views.library.library-details.import-record')
</script>
