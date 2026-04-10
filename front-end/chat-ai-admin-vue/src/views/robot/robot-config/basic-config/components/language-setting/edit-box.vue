<style lang="less" scoped>
.edit-box {
  border-radius: 6px;
  overflow: hidden;
  background-color: #f2f4f7;

  .edit-box-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    .left-box {
      display: flex;
      align-items: center;
    }
    .edit-box-icon {
      margin-right: 8px;
      font-size: 18px;
      color: #262626;
    }

    .edit-box-title {
      display: flex;
      gap: 4px;
      line-height: 24px;
      font-size: 16px;
      font-weight: 600;
      color: #262626;
    }
  }
}
</style>

<template>
  <div class="edit-box">
    <div class="edit-box-header">
      <div class="left-box">
        <svg-icon class="edit-box-icon" :name="props.iconName"></svg-icon>
        <span class="edit-box-title"
          >{{ props.title }}
          <span v-if="$slots.icon"><slot name="icon"></slot></span>
        </span>
        <span class="title-tip" v-if="$slots.tip">
          <slot name="tip"></slot>
        </span>
      </div>
    </div>

    <div class="edit-box-body" :style="{ ...bodyStyle }">
      <slot></slot>
    </div>
  </div>
</template>
<script setup>
import { useI18n } from '@/hooks/web/useI18n'


const props = defineProps({
  iconName: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  bodyStyle: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

const { t } = useI18n('views.robot.robot-config.basic-config.components.edit-box')
</script>
