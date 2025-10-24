<template>
  <a-popover
    color="#fff"
    trigger="hover"
    placement="bottomRight"
    :overlayStyle="{ 'max-width': '400px', width: '400px' }"
    :overlayInnerStyle="{ padding: '0' }"
    class="share-popup"
    @openChange="openChange"
  >
    <template #content>
      <share-form ref="shareFormRef" layout="vertical" :hide-copy="true" :baseUrl="baseUrl" />
    </template>
    <slot name="default"></slot>
  </a-popover>
</template>

<script setup>
import { OPEN_BOC_BASE_URL } from '@/constants/index'
import { nextTick, ref } from 'vue'
import ShareForm from './share-form.vue'

const props = defineProps({
  docKey: {
    type: String
  },
  libraryId: {
    type: [Number, String],
    default: ''
  },
  libraryKey: {
    type: String,
    default: ''
  },
  baseUrl: {
    type: String,
    default: OPEN_BOC_BASE_URL + '/doc'
  }
})

const shareFormRef = ref(null)

const openChange = async (visible) => {
  if (visible) {
    nextTick(() => {
      shareFormRef.value.init({
        doc_key: props.docKey
      })
    })
  }
}
</script>

<style lang="less" scoped></style>
