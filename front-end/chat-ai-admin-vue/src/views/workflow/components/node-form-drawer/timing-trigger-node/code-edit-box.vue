<template>
  <div>
    <Codemirror
      @wheel.stop=""
      originalStyle
      class="cm-component"
      v-model:value="code"
      :options="cmOptions"
      ref="cmRef"
      :height="height"
      :width="width"
      @change="onChange"
      @input="onInput"
      @ready="onReady"
    >
    </Codemirror>
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import 'codemirror/mode/shell/shell.js'
import Codemirror from 'codemirror-editor-vue3'
import 'codemirror/theme/dracula.css'

const emit = defineEmits(['update:value', 'change'])

const code = ref('')
const errorMessage = ref('')

const props = defineProps({
  value: {
    type: String,
    default: ''
  },
  width: {
    type: Number,
    default: 512
  },
  height: {
    type: Number,
    default: 150
  }
})

// 验证 cron 表达式的函数
const validateCronExpression = (expression) => {
  if (!expression) {
    return 'Cron表达式不能为空'
  }

  // 去除首尾空格并分割成字段
  const fields = expression.trim().split(/\s+/)
  
  // 标准的 Linux Cron 表达式应该有 5 个字段
  if (fields.length !== 5) {
    return `Cron表达式应包含5个字段（分钟 小时 日 月 周），当前为${fields.length}个字段`
  }

  const [minute, hour, dayOfMonth, month, dayOfWeek] = fields

  // 验证每个字段
  const validations = [
    { value: minute, name: '分钟', range: [0, 59], specialChars: ['*', '/', '-', ','] },
    { value: hour, name: '小时', range: [0, 23], specialChars: ['*', '/', '-', ','] },
    { value: dayOfMonth, name: '日期', range: [1, 31], specialChars: ['*', '/', '-', ',', '?', 'L', 'W'] },
    { value: month, name: '月份', range: [1, 12], specialChars: ['*', '/', '-', ','] },
    { value: dayOfWeek, name: '星期', range: [0, 7], specialChars: ['*', '/', '-', ',', '?', 'L', '#'], aliases: ['SUN', 'MON', 'TUE', 'WED', 'THU', 'FRI', 'SAT'] }
  ]

  for (const validation of validations) {
    const error = validateField(validation.value, validation)
    if (error) {
      return `${validation.name}字段错误: ${error}`
    }
  }

  // 检查 dayOfMonth 和 dayOfWeek 是否同时使用了 ?
  // if (dayOfMonth !== '?' && dayOfWeek !== '?') {
  //   if (dayOfMonth.includes('*') || dayOfWeek.includes('*')) {
  //     return '日期和星期字段不能同时使用通配符，需要其中一个设置为?'
  //   }
  // }

  return null
}

// 验证单个字段
const validateField = (value, config) => {
  if (!value) {
    return '字段不能为空'
  }

  // 处理特殊值和别名
  if (config.aliases) {
    const upperValue = value.toUpperCase()
    if (config.aliases.includes(upperValue)) {
      return null
    }
  }

  // 分割逗号分隔的值
  const parts = value.split(',')
  for (const part of parts) {
    const trimmedPart = part.trim()
    
    // 允许的特殊字符
    if (trimmedPart === '*') continue
    
    // 处理 ? 特殊字符
    if (trimmedPart === '?' && config.specialChars.includes('?')) continue
    
    // 处理范围和步长
    if (trimmedPart.includes('/')) {
      const [range, step] = trimmedPart.split('/')
      if (!isValidNumber(step, 1, 59)) {
        return `步长必须是1-59之间的数字`
      }
      
      if (range === '*') continue
      
      if (range.includes('-')) {
        const [start, end] = range.split('-')
        if (!isValidNumber(start, ...config.range) || !isValidNumber(end, ...config.range)) {
          return `范围值必须在${config.range[0]}-${config.range[1]}之间`
        }
        if (parseInt(start) > parseInt(end)) {
          return `起始值不能大于结束值`
        }
      } else if (!isValidNumber(range, ...config.range)) {
        return `范围值必须在${config.range[0]}-${config.range[1]}之间`
      }
      continue
    }
    
    // 处理范围
    if (trimmedPart.includes('-')) {
      const [start, end] = trimmedPart.split('-')
      if (!isValidNumber(start, ...config.range) || !isValidNumber(end, ...config.range)) {
        return `范围值必须在${config.range[0]}-${config.range[1]}之间`
      }
      if (parseInt(start) > parseInt(end)) {
        return `起始值不能大于结束值`
      }
      continue
    }
    
    // 处理 L, W, # 等特殊字符
    if (trimmedPart.includes('L') || trimmedPart.includes('W') || trimmedPart.includes('#')) {
      if (config.specialChars.some(char => trimmedPart.includes(char))) {
        continue
      } else {
        return `该字段不支持此特殊字符`
      }
    }
    
    // 最后检查是否为有效数字
    if (!isValidNumber(trimmedPart, ...config.range)) {
      return `值必须在${config.range[0]}-${config.range[1]}之间`
    }
  }
  
  return null
}

// 验证数字是否在范围内
const isValidNumber = (value, min, max) => {
  if (isNaN(value)) return false
  const num = parseInt(value, 10)
  return num >= min && num <= max
}

const validateAndSetError = (cronExpression) => {
  const error = validateCronExpression(cronExpression)
  errorMessage.value = error || ''
}

watch(
  () => props.value,
  (val) => {
    code.value = val
    validateAndSetError(val)
  },
  {
    immediate: true
  }
)

// 监听代码变化并验证
watch(code, (newCode) => {
  validateAndSetError(newCode)
})



const cmOptions = {
  mode: 'text/x-sh',
  theme: 'dracula'
}

const cmRef = ref(null)

const onChange = (val) => {
  emit('update:value', val)
  emit('change', val)
}

const handleRefresh = ()=>{
  cmRef.value?.refresh()
}

const onReady = () => {}
const onInput = () => {}

defineExpose({
  handleRefresh
})
</script>

<style lang="less" scoped>
.cm-component {
  border: 1px solid #ddd;
  border-radius: 4px;
}

.error-message {
  color: #ff4d4f;
  font-size: 12px;
  margin-top: 5px;
  padding: 5px;
  background-color: #fff2f0;
  border: 1px solid #ffccc7;
  border-radius: 4px;
}
</style>