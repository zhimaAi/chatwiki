<template>
  <div class="help-alert" :class="[`help-alert-${type}`, { 'is-closable': closable }]" v-show="visible">
    <span class="help-alert-icon">
      <slot name="icon">
        <InfoCircleFilled v-if="type === 'info'" />
        <CheckCircleFilled v-if="type === 'success'" />
        <ExclamationCircleFilled v-if="type === 'warning'" />
        <CloseCircleFilled v-if="type === 'error'" />
      </slot>
    </span>

    <div class="help-alert-content">
      <div class="help-alert-title">
        <slot name="title">{{ props.title }}</slot>
      </div>
      <div class="help-alert-message" v-if="$slots.default || props.message">
        <slot>{{ props.message }}</slot>
      </div>
    </div>

    <button v-if="closable" class="help-alert-close" @click="handleClose">
      <CloseOutlined />
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { InfoCircleFilled, CheckCircleFilled, ExclamationCircleFilled, CloseCircleFilled, CloseOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  title: {
    type: String,
    default: '使用说明'
  },
  message: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['info', 'success', 'warning', 'error'].includes(value)
  },
  closable: {
    type: Boolean,
    default: false
  }
})

const visible = ref(true)
const emit = defineEmits(['close'])

const handleClose = () => {
  visible.value = false
  emit('close')
}
</script>

<style lang="less" scoped>
.help-alert {
  position: relative;
  display: flex;
  align-items: flex-start;
  padding: 9px 16px;
  border-radius: 6px;

  &.is-closable {
    padding-right: 36px;
  }

  .help-alert-icon {
    height: 22px;
    line-height: 22px;
    margin-right: 8px;
    font-size: 16px;
  }

  .help-alert-content {
    flex: 1;
  }

  .help-alert-title {
    line-height: 22px;
    font-size: 14px;
    font-weight: 600;
    color: #242933;
  }

  .help-alert-message {
    line-height: 22px;
    font-size: 12px;
    font-weight: 400;
    color: #3a4559;
    margin-top: 4px;
  }

  .help-alert-close {
    position: absolute;
    top: 12px;
    right: 16px;
    padding: 0;
    border: none;
    background: transparent;
    font-size: 14px;
    line-height: 1;
    cursor: pointer;
    color: #909399;

    &:hover {
      color: #606266;
    }
  }
}

.help-alert-info {
  background-color: #E9F1FE;
  border: 1px solid #99BFFD;

  .help-alert-icon {
    color: #2475fc;
  }
}

.help-alert-success {
  background-color: #f6ffed;
  border: 1px solid #b7eb8f;

  .help-alert-icon {
    color: #52c41a;
  }
}

.help-alert-warning {
  background-color: #fffbe6;
  border: 1px solid #ffe58f;

  .help-alert-icon {
    color: #faad14;
  }
}

.help-alert-error {
  background-color: #fff2f0;
  border: 1px solid #ffccc7;

  .help-alert-icon {
    color: #ff4d4f;
  }
}
</style>