<style lang="less" scoped>
.gradient-color-picker {
  .color-picker-box {
    display: flex;
    align-items: center;
    .color-item {
      width: 125px;
      margin: 0 8px;
    }
    .right-arrow {
      font-size: 14px;
      color: #8c8c8c;
    }
  }

  .color-picker-demo {
    width: 100%;
    height: 32px;
    margin-top: 8px;
  }
}
</style>

<template>
  <div class="gradient-color-picker">
    <div class="color-picker-box">
      <div class="gradient-type-select">
        <a-select :value="gradientDirection" style="width: 180px" @change="onChangeType">
          <a-select-option
            :value="opt.value"
            :type="opt.type"
            v-for="opt in options"
            :key="opt.key"
            >{{ opt.label }}</a-select-option
          >
        </a-select>
      </div>
      <div class="color-item">
        <ColorPicker v-model:value="starColor" />
      </div>
      <template v-if="gradientType != 'color'">
        <div class="right-arrow"><RightOutlined /></div>
        <div class="color-item"><ColorPicker v-model:value="endColor" /></div>
      </template>
    </div>
    <!-- <div class="color-picker-demo" :style="{ background: gradientColor }"></div> -->
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { RightOutlined } from '@ant-design/icons-vue'
import ColorPicker from '@/components/color-picker/index.vue'
import { Form } from 'ant-design-vue'

const emit = defineEmits(['update:value'])
const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const formItemContext = Form.useInjectFormItemContext()

const options = ref([
  {
    key: 1,
    label: '单色',
    value: 'color',
    type: 'color'
  },
  {
    key: 2,
    label: '线性渐变（横向）',
    value: 'to right',
    type: 'linear-gradient'
  },
  {
    key: 3,
    label: '线性渐变（竖向）',
    value: 'to bottom',
    type: 'linear-gradient'
  },
  {
    key: 4,
    label: '径向渐变',
    value: 'circle',
    type: 'radial-gradient'
  }
])

const gradientType = ref('linear-gradient')
const gradientDirection = ref('to right')
const starColor = ref('#2435E7')
const endColor = ref('#01A0FB')

const gradientColor = computed(() => {
  if (gradientType.value == 'color') {
    return starColor.value
  }

  return `${gradientType.value}(${gradientDirection.value}, ${starColor.value}, ${endColor.value})`
})

watch(
  () => gradientColor.value,
  () => {
    let val = [gradientType.value, gradientDirection.value, starColor.value, endColor.value].join(
      ','
    )

    emit('update:value', val)

    formItemContext.onFieldChange()
  }
)
watch(
  () => props.value,
  () => {
    const val = props.value.split(',')
    const [type, direction, color1, color2] = val

    gradientType.value = type
    gradientDirection.value = direction
    starColor.value = color1
    endColor.value = color2
  },
  { immediate: true }
)

const onChangeType = (value, option) => {
  gradientType.value = option.type
  gradientDirection.value = value
  formItemContext.onFieldChange()
}
</script>
