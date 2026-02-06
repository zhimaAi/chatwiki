<template>
  <div class="_full-input-box">
    <FullscreenOutlined v-if="canFull" class="full-icon" @click="showFull"/>
    <AtInput
      :type="type"
      :inputStyle="inputStyle"
      :options="options"
      :defaultValue="innerValue"
      :defaultSelectedList="innerTags"
      :placeholder="computedPlaceholder"
      :canRepeat="canRepeat"
      :isContain="isContain"
      :checkAnyLevel="checkAnyLevel"
      ref="atInputRef"
      @open="emit('open')"
      @change="handleChange"
    >
      <template #option="{ label, payload }">
        <div class="field-list-item">
          <div class="field-label">{{ label }}</div>
          <div class="field-type">{{ payload.typ }}</div>
        </div>
      </template>
    </AtInput>

    <FullAtInput
      :options="options"
      :defaultValue="innerValue"
      :defaultSelectedList="innerTags"
      :placeholder="t('ph_input_content_with_variable')"
      :type="type"
      @open="emit('open')"
      @change="handleChange"
      @ok="handleRefreshAtInput"
      ref="fullAtInputRef"
    />
  </div>
</template>

<script setup>
import {ref, watch, computed} from 'vue'
import { useI18n } from '@/hooks/web/useI18n'
import {FullscreenOutlined} from '@ant-design/icons-vue'
import AtInput from "./at-input.vue";
import FullAtInput from "@/views/workflow/components/at-input/full-at-input.vue";

const { t } = useI18n('views.workflow.components.at-input.at-full-input')
const emit = defineEmits(['open', 'change'])
const props = defineProps({
  type: {
    type: String,
    default: "input", // input, textarea
    validator(value) {
      return ["input", "textarea"].includes(value);
    },
  },
  placeholder: {
    type: String,
    default: "ph_input",
  },
  defaultValue: {
    type: [String, Number],
    default: "",
  },
  options: {
    type: Array,
    default: () => [],
  },
  defaultSelectedList: {
    type: Array,
    default: () => [],
  },
  canRepeat: {
    type: Boolean,
    default: true,
  },
  isContain: {
    type: Boolean,
    default: false,
  },
  inputStyle: {
    type: String,
    default: "",
  },
  checkAnyLevel:{
    type: Boolean,
    default: false
  },
  canFull:{
    type: Boolean,
    default: false
  }
})

const atInputRef = ref(null)
const fullAtInputRef = ref(null)
const innerValue = ref(props.defaultValue)
const innerTags = ref(props.defaultSelectedList)
const computedPlaceholder = computed(() => {
  return props.placeholder === 'ph_input' ? t('ph_input') : props.placeholder
})

watch(
  () => props.defaultValue,
  (val) => {
    innerValue.value = val
    //atInputRef.value?.refresh()
  }
)

watch(
  () => props.defaultSelectedList,
  (val) => {
    innerTags.value = val
  }
)

function handleChange(val, tags) {
  innerValue.value = val
  innerTags.value = tags
  emit('change', val, tags)
}

function handleRefreshAtInput() {
  atInputRef.value.refresh()
}

function showFull() {
  fullAtInputRef.value.show()
}

defineExpose({
  showFull
})
</script>

<style scoped lang="less">
._full-input-box {
  width: 100%;
  position: relative;

  .full-icon {
    position: absolute;
    bottom: 8px;
    right: 8px;
    cursor: pointer;
    z-index: 999;
  }
}
</style>
