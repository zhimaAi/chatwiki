<template>
  <div>
    <a-modal v-model:open="open" :title="t('title_field_detail')" :width="746">
      <div class="input-box">
        <at-input
          inputStyle="height: 550px;"
          :options="props.options"
          :defaultSelectedList="props.defaultSelectedList"
          :defaultValue="props.defaultValue"
          :placeholder="computedPlaceholder"
          :type="props.type"
          ref="atInputRef"
          @open="showAtList"
          @change="(text, selectedList) => changeValue(text, selectedList)"
        />
      </div>
      <template #footer>
        <a-button type="primary" @click="handleOk">{{ t('btn_confirm') }}</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, nextTick, watch, computed } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import AtInput from './at-input.vue'

const { t } = useI18n('views.workflow.components.at-input.full-at-input')
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
    default: 'ph_input_message_with_variable'
  },
  type: {
    type: String,
    default: 'textarea'
  }
})

const computedPlaceholder = computed(() => {
  return props.placeholder === 'ph_input_message_with_variable' ? t('ph_input_message_with_variable') : props.placeholder
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