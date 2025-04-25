<template>
  <CuScroll :scrollbar="scrollbar" ref="scrollRef" @scroll="handleScroll" @onScrollEnd="onScrollEnd">
    <VuePdfEmbed
      v-for="page in visiblePages"
      @click="clickPage(page)"
      :key="page"
      :source="source"
      :page="page"
      @rendered="() => pdfRendered(page)"
      :style="getItemStyle"
    />
    <div class="loading"><a-spin v-if="loading"/></div>
  </CuScroll>
</template>

<script setup>
import {ref, onMounted, nextTick, computed} from 'vue'
import CuScroll from "@/components/cu-scroll/cu-scroll.vue";
import VuePdfEmbed from 'vue-pdf-embed'

const props = defineProps({
  source: {
    type: [String],
    required: true
  },
  pageSize: {
    type: Number,
    default: 10
  },
})

const emit = defineEmits(['onScrollEnd', 'select', 'rendered', 'scroll', 'getTotalPage'])
const scrollRef = ref(null)
const visiblePages = ref([])
const totalPages = ref(0)
const loading = ref(false)

const getItemStyle = computed(() => {
  if (!loading.value) {
    return {'border-bottom': '2px solid #d9d9d9'}
  }
  return null
})

const scrollbar = ref({
  fade: false,
  scrollbarTrackClickable: true,
  interactive: true,
  minSize: 40,
})

onMounted(async () => {
  const pdf = await VuePdfEmbed.getDocument(props.source).promise
  // 获取文件总页数
  totalPages.value = pdf.numPages
  emit('getTotalPage', totalPages.value)
  loadMore()
})

const loadMore = (callback=null) => {
  const current = visiblePages.value.length
  if (current >= totalPages.value || loading.value) return
  loading.value = true
  const nextPages = []
  for (let i = 1; i <= props.pageSize && current + i <= totalPages.value; i++) {
    nextPages.push(current + i)
  }
  visiblePages.value.push(...nextPages)
  nextTick(() => {
    typeof callback === 'function' && callback()
  })
}

const onScrollEnd = (e) => {
  loadMore()
  emit('onScrollEnd', e)
}

const pdfRendered = (page) => {
  nextTick(() => {
   if (page === visiblePages.value[visiblePages.value.length - 1]) {
     loading.value = false
     emit('rendered', page, totalPages.value)
   }
  })
}

const clickPage = page => {
  emit('select', page)
}

const getScrollInstance = () => {
  return scrollRef.value
}

const handleScroll = (event)=>{
  emit('scroll', event)
}

defineExpose({
  loadMore,
  getScrollInstance,
})
</script>

<style scoped lang="less">
.loading {
    padding: 12px;
    text-align: center;
}
</style>
