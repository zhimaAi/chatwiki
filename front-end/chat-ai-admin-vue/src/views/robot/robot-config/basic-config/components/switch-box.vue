<style lang="less" scoped>
.switch-box {
  display: flex;
  align-items: center;
  padding: 16px;
  border-radius: 6px;
  background-color: #f2f4f7;
  overflow: hidden;
  .switch-box-body {
    flex: 1;
  }

  .switch-box-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .left-box {
      display: flex;
      align-items: center;
    }

    .right-box {
      display: flex;
      align-items: center;
    }

    .switch-box-icon {
      margin-right: 8px;
      font-size: 18px;
      color: #262626;
    }

    .switch-box-title {
      line-height: 24px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }
    .title-tip {
      margin-left: 4px;
    }
  }

  .switch-desc {
    line-height: 22px;
    font-size: 14px;
    font-weight: 400;
    color: #595959;
  }
}
</style>

<template>
  <div class="switch-box">
    <div class="switch-box-body">
      <div class="switch-box-header">
        <div class="left-box">
          <svg-icon class="switch-box-icon" :name="props.iconName"></svg-icon>
          <span class="switch-box-title">{{ props.title }}</span>
          <span class="title-tip" v-if="$slots.tip">
            <slot name="tip"></slot>
          </span>
        </div>
        <div class="right-box">
          <div class="extra-box">
            <slot name="extra">
              <div class="extra-box"></div>
            </slot>
          </div>
        </div>
      </div>

      <div class="switch-desc">
        <slot></slot>
      </div>
    </div>

    <div class="switch-box-right">
      <a-switch
        v-model:checked="checked"
        @change="handleChange"
        checked-children="开"
        un-checked-children="关"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const emit = defineEmits(['update:value', 'change'])

const props = defineProps({
  iconName: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  value: {
    type: Boolean,
    default: false
  }
})

const checked = ref(props.value)

const handleChange = (val) => {
  emit('update:value', val)
  emit('change', val)
}
</script>
