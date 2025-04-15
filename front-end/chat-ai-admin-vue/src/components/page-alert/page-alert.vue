<template>
  <helpAlert type="info" closable  @close="handleCloseTip" v-if="show">
    <template #title>
      <slot name="title">{{ props.title }}</slot>
    </template>
    <slot>{{ message }}</slot>
  </helpAlert>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { useRoute } from 'vue-router';
import { computed } from 'vue';
import helpAlert from '../help-alert/help-alert.vue';

const emit = defineEmits(['close'])
const route = useRoute()
const { setStorage, getStorage } = useStorage('localStorage')

const props = defineProps({
  title: {
    type: String,
    default: '使用说明'
  },
  message: {
    type: String,
    default: ''
  },
  id: {
    type: [String, Number],
    default: ''
  }
})

const closeKeys = getStorage('pageAliertCloseKeys') || []
const show = computed(() => {
  let key = props.id || route.path

  return !closeKeys.includes(key) 
})

const handleCloseTip = () => {
  let key = props.id || route.path

  if (!closeKeys.includes(key)) {
    setStorage('pageAliertCloseKeys', [...closeKeys, key]) 
  }

  emit('close')
}
</script>