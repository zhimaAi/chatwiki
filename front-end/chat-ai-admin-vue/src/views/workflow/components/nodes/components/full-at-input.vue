<template>
  <div>
    <a-modal v-model:open="open" title="字段详情" :width="746">
      <div class="input-box">
        <at-input
          inputStyle="height: 550px;"
          :options="props.options"
          :defaultSelectedList="props.defaultSelectedList"
          :defaultValue="props.defaultValue"
          :placeholder="props.placeholder"
          :type="props.type"
          ref="atInputRef"
          @open="showAtList"
          @change="(text, selectedList) => changeValue(text, selectedList)"
        />
      </div>
      <template #footer>
        <a-button type="primary" @click="handleOk">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import AtInput from '../at-input/at-input.vue'
const open = ref(false)

const emit = defineEmits(['open', 'change', 'ok'])

watch(
  () => open.value,
  (val) => {
    if (!val) {
      emit('ok')
    }
  }
)

const props = defineProps({
  options: {
    type: Array,
    default: () => []
  },
  defaultSelectedList: {
    type: Array,
    default: () => []
  },
  defaultValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '请输入消息内容，键入“/”可以插入变量'
  },
  type: {
    type: String,
    default: 'textarea'
  }
})

const atInputRef = ref(null)
const show = () => {
  open.value = true
  nextTick(() => {
    atInputRef.value.refresh()
  })
}

const showAtList = (val) => {
  emit('open', val)
}

const changeValue = (text, selectedList) => {
  emit('change', text, selectedList)
}

const handleOk = () => {
  open.value = false
}

defineExpose({ show })
</script>

<style lang="less" scoped>
.input-box {
  margin: 24px 0;
  &::v-deep(.type-textarea) {
    &::-webkit-scrollbar {
      width: 6px;
    }
    &::-webkit-scrollbar-thumb {
      border-radius: 20px;
    }
  }
}
</style>
