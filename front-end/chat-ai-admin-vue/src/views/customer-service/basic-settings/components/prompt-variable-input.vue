<template>
  <AtInput v-bind="$attrs" ref="input" semantic-tags :enable-card-insert="true" :card-options="cardOptions" :get-popup-container="getPopupContainer">
    <template #dropdown="{ options, select, cardOptions: cards, selectCard, close }">
      <SlashPopupPanel workspace :options="options" :card-options="cards" :select="select" :select-card="selectCard" :close="close" @add-tag="$emit('addTag')" @add-card="$emit('addCard')" @edit-card="card => $emit('editCard', card)" />
    </template>
  </AtInput>
</template>
<script setup>
import { ref } from 'vue'
import AtInput from '@/views/robot/robot-config/basic-config/components/at-input/index.vue'
import SlashPopupPanel from '@/views/robot/robot-config/basic-config/components/slash-popup-panel.vue'
defineOptions({ inheritAttrs: false })
defineProps({ cardOptions: { type: Array, default: () => [] }, getPopupContainer: { type: Function, default: null } })
defineEmits(['addTag', 'addCard', 'editCard'])
const input = ref(null)
defineExpose({ refresh: () => input.value?.refresh?.(), updateValue: () => input.value?.updateValue?.() })
</script>
