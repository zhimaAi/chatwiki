<style lang="less" scoped>
.zoom-select {
  position: relative;
  margin: 0 4px;

  .zoom-select-label {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 80px;
    height: 32px;
    border-radius: 6px;
    transition: all 0.2s;
    .zoom-select-input{
      width: 100%;
      height: 100%;
      border: none;
      outline: none;
      text-align: center;
      font-size: 14px;
      color: #595959;
    }
  }

  .zoom-select-options {
    position: absolute;
    bottom: 48px;
    left: 0;
    padding: 2px;
    font-size: 14px;
    border-radius: 6px;
    color: #595959;
    background-color: #fff;
    box-shadow:
      0 6px 30px 5px #0000000d,
      0 16px 24px 2px #0000000a,
      0 8px 10px -5px #00000014;
  }
  .zoom-select-option {
    width: 77px;
    height: 32px;
    line-height: 32px;
    margin-bottom: 2px;
    border-radius: 6px;
    color: #595959;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .option-line {
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 2px;
  }

  .zoom-select-label:hover,
  .zoom-select-option:hover {
    background-color: #e4e6eb;
    cursor: pointer;
  }
}
</style>

<template>
  <div class="zoom-select">
    <div class="zoom-select-label" @click.stop="">
      <input class="zoom-select-input" type="text" v-model="currentValue" @input="handleInput" @focus="showMenu" @blur="handleBlur">
    </div>
    <div class="zoom-select-options" v-show="isShowMenu">
      <div
        class="zoom-select-option"
        v-for="item in menus"
        :key="item.value"
        @click="handleMenu(item)"
      >
        {{ item.label }}
      </div>
      <div class="option-line"></div>
      <div class="zoom-select-option" @click="handleFitView()">{{ t('btn_fit_view') }}</div>
    </div>
  </div>
</template>

<script>
import { useI18n } from '@/hooks/web/useI18n'

export default {
  name: 'NodeCommon',
  props: {
    modelValue: {
      type: Number,
      default: 100,
    },
  },
  setup() {
    const { t } = useI18n('views.workflow.components.custom-control.zoom-select')
    return { t }
  },
  data() {
    return {
      isShowMenu: false,
      currentValue: `${this.modelValue}%`,
      previousValue: this.modelValue,
      debounceTimer: null,
      menus: [
        { label: '10%', value: 10 },
        { label: '50%', value: 50 },
        { label: '80%', value: 80 },
        { label: '100%', value: 100 },
        { label: '150%', value: 150 },
        { label: '200%', value: 200 },
        { label: '400%', value: 400 },
      ],
    }
  },
  watch: {
    modelValue(newValue) {
      // 当 modelValue 从外部（如父组件）改变时，更新输入框的显示
      // 检查以避免在内部输入时发生不必要的循环更新
      if (newValue !== parseInt(this.currentValue.replace(/[^\d]/g, ''), 10)) {
        this.currentValue = `${newValue}%`
        this.previousValue = newValue
      }
    },
  },
  mounted() {
    document.addEventListener('click', this.documentClick)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.documentClick)
  },
  methods: {
    handleMenu(item) {
      this.$emit('update:modelValue', item.value)
      this.$emit('zoom-change', item.value)
      this.isShowMenu = false
    },
    showMenu() {
      this.isShowMenu = !this.isShowMenu
    },
    handleFitView() {
      this.isShowMenu = false
      this.$emit('fitView')
    },
    handleInput() {
      clearTimeout(this.debounceTimer)
      this.debounceTimer = setTimeout(() => {
        // 提取有效数字
        const inputText = this.currentValue.replace('%', '').trim()
        const numberPart = parseInt(inputText.replace(/[^\d]/g, ''), 10)
        
        if (!inputText) {
          // 空白输入时，显示当前有效值
          this.currentValue = `${this.previousValue}%`
          return
        }
        
        if (isNaN(numberPart)) {
          // 无有效数字，立即恢复显示
          this.currentValue = `${this.previousValue}%`
          return
        }
        
        // 数字超出范围时，自动修正
        let validValue = numberPart
        if (validValue > 800) validValue = 800
        if (validValue < 1) validValue = 1
        
        // 即时显示修正后的值
        this.currentValue = `${validValue}%`
        
        // 发射事件
        this.previousValue = validValue  // 更新记忆值
        this.$emit('update:modelValue', validValue)
        this.$emit('zoom-change', validValue)
      }, 300)
    },
    handleBlur() {
      clearTimeout(this.debounceTimer)
      
      const numberPart = parseInt(this.currentValue.replace(/[^\d]/g, ''), 10)
      
      // 尝试使用最终输入值
      if (!isNaN(numberPart) && numberPart >= 1 && numberPart <= 800) {
        this.$emit('update:modelValue', numberPart)
        this.$emit('zoom-change', numberPart)
        this.currentValue = `${numberPart}%`
        this.previousValue = numberPart
      } else {
        // 恢复显示为记忆的有效值
        this.currentValue = `${this.previousValue}%`
      }
    },
    documentClick(e) {
      if (this.isShowMenu) {
        const menus = this.$el.querySelector('.zoom-select-options')
        if (
          !menus.contains(e.target) &&
          e.target !== this.$el.querySelector('.zoom-select-label')
        ) {
          this.isShowMenu = false
        }
      }
    },
  },
}
</script>
