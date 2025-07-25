<template>
  <div>
    <a-modal v-model:open="open" title="购买记录" :footer="null" :width="896">
      <a-table style="margin-top: 16px" :data-source="list">
        <a-table-column key="create_time" title="购买时间" data-index="create_time">
          <template #default="{ record }">
            {{ formatTime(record.create_time) }}
          </template>
        </a-table-column>
        <a-table-column key="buy_total" title="购买积分数" data-index="buy_total"> </a-table-column>
        <a-table-column key="use_total" title="已消耗" data-index="use_total"> </a-table-column>
        <a-table-column key="expire_total" title="已过期" data-index="expire_total">
        </a-table-column>
        <a-table-column key="surplus_total" title="剩余额度" data-index="surplus_total">
        </a-table-column>
        <a-table-column key="expire_time" title="到期时间" data-index="expire_time">
          <template #default="{ record }">
            {{ formatTime(record.expire_time) }}
          </template>
        </a-table-column>
        <a-table-column key="remark_outer" title="备注" data-index="remark_outer"> </a-table-column>
      </a-table>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import dayjs from 'dayjs'
const props = defineProps({
  list: {
    type: Array,
    default: () => []
  }
})
const open = ref(false)

const formatTime = (time) => {
  if (time <= 0) {
    return '--'
  }
  return dayjs(time * 1000).format('YYYY/MM/DD HH:mm:ss')
}

const show = () => {
  open.value = true
}
defineExpose({
  show
})
</script>

<style lang="less" scoped></style>
