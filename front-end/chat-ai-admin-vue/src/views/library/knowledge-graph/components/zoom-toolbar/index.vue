<style lang="less" scoped>
.custom-control-warpper {
  position: relative;
  user-select: none;

  .custom-control {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;
    padding: 4px 12px;
    border-radius: 8px;
    background-color: #fff;
    box-shadow: 0 4px 16px 0 #0000001a;
  }

  .action-btn {
    border-radius: 6px;
    transition: all 0.2s;
    &:hover {
      background-color: #e4e6eb;
      cursor: pointer;
    }
  }

  .zoom-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    font-size: 16px;
    color: #595959;
  }
}
</style>

<template>
  <div class="custom-control-warpper">
    <div class="custom-control">
      <div class="action-btn zoom-btn" @click="handleReduce">
        <svg-icon name="minus" size="16" />
      </div>
      
      <zoom-select @change="zoomSelectChagne" :title="zoomDisplayValue" />

      <div class="action-btn zoom-btn" @click="handleAmplify">
        <svg-icon name="plus" size="16" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ZoomSelect from './zoom-select.vue'

const emit = defineEmits(['change'])

const props = defineProps({
  value: {
    type: Number,
    default: 1
  }
})


const zoomDisplayValue = computed(() => {
  return `${(props.value * 100).toFixed(0)}%`
})

const zoomSelectChagne = ({ value }) => {
  change(value)
}

const handleReduce = () => {
  // 四舍五入到整数
  let value = props.value - 0.1
  if (value < 0.1) {
    value = 0.1
  }

  change(value)
}

const handleAmplify = () => {
  let value = props.value + 0.1

  if (value > 5) {
    value = 5
  }

  change(value)
}

const change = (value) => {
  value = Number(value.toFixed(1))
  console.log(value)
  emit('change', value)
}
</script>
