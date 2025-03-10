<template>
  <vue-pdf-embed ref="vuePdfRef" @loaded="handleLoad" :source="props.source" />
</template>
<script setup>
import VuePdfEmbed from 'vue-pdf-embed'

const emit = defineEmits(['loaded'])

const props = defineProps({
  source: {
    type: String
  }
})

const handleLoad = (e) => {
  emit('loaded', e)
}
</script>
