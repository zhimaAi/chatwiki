<template>
  <a-modal
    title="版本历史"
    v-model:open="open"
    width="500px"
    :footer="null"
  >
    <a-table :columns="columns" :data-source="list"></a-table>
  </a-modal>
</template>

<script setup>
import {ref} from 'vue';

const open = ref(false)
const list = ref([])
const columns = [
  {
    title: '版本',
    dataIndex: 'version',
    width: 100,
  },
  {
    title: '更新时间',
    dataIndex: 'update_time_text',
    width: 160,
  },
  {
    title: '发布者',
    dataIndex: 'author',
    width: 100,
  }
]

function show(versions= []) {
  list.value = versions
  open.value = true
}

defineExpose({
  show,
})
</script>

<style scoped>

</style>
