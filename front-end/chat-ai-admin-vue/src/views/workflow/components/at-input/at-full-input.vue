<template>
  <div>
    <AtInput
      :type="type"
      :inputStyle="inputStyle"
      :options="options"
      :defaultValue="innerValue"
      :defaultSelectedList="innerTags"
      :placeholder="placeholder"
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
      placeholder="请输入内容，键入“/”可以插入变量"
      :type="type"
      @open="emit('open')"
      @change="handleChange"
      @ok="handleRefreshAtInput"
      ref="fullAtInputRef"
    />
  </div>
</template>

<script setup>
import {ref, watch} from 'vue'
import AtInput from "./at-input.vue";
import FullAtInput from "@/views/workflow/components/at-input/full-at-input.vue";

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
    default: "请输入",
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
  }
})

const atInputRef = ref(null)
const fullAtInputRef = ref(null)
const innerValue = ref(props.defaultValue)
const innerTags = ref(props.defaultSelectedList)

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

<style scoped>

</style>
