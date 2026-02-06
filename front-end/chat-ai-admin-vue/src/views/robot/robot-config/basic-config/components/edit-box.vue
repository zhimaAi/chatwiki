<style lang="less" scoped>
.edit-box {
  border-radius: 6px;
  overflow: hidden;
  background-color: #f2f4f7;

  .edit-box-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;

    .left-box {
      display: flex;
      align-items: center;
    }

    .right-box {
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
    .title-tip {
      margin-left: 4px;
    }

    .actions-box {
      display: flex;
      align-items: center;
      line-height: 22px;
      font-size: 14px;
      color: #595959;

      .action-btn {
        cursor: pointer;
      }

      .save-btn {
        color: #2475fc;
      }
    }
  }

  .edit-box-body {
    padding: 0 16px 16px 16px;
  }
}
</style>

<template>
  <div class="edit-box" :class="{ 'is-edit': isEdit }">
    <div class="edit-box-header">
      <div class="left-box">
        <svg-icon class="edit-box-icon" :name="props.iconName"></svg-icon>
        <span class="edit-box-title">{{ props.title }}
          <span v-if="$slots.icon"><slot name="icon"></slot></span>
        </span>
        <span class="title-tip" v-if="$slots.tip">
          <slot name="tip"></slot>
        </span>
      </div>
      <div class="right-box">
        <div class="extra-box">
          <slot name="extra">
            <div class="actions-box">
              <template v-if="props.isEdit">
                <a-flex :gap="8">
                  <a-button @click="handleSave" size="small" type="primary">{{ t('btn_save') }}</a-button>
                  <a-button @click="handleEdit(false)" size="small">{{ t('btn_cancel') }}</a-button>
                </a-flex>
              </template>
              <template v-else>
                <a-button @click="handleEdit(true)" size="small">{{ t('btn_edit') }}</a-button>
              </template>
            </div>
          </slot>
        </div>
      </div>
    </div>

    <div class="edit-box-body" :style="{ ...bodyStyle }">
      <slot></slot>
    </div>
  </div>
</template>
<script setup>
import { useI18n } from '@/hooks/web/useI18n'

const emit = defineEmits(['update:isEdit', 'save', 'edit'])

const props = defineProps({
  iconName: {
    type: String,
    default: ''
  },
  title: {
    type: String,
    default: ''
  },
  isEdit: {
    type: Boolean,
    default: false
  },
  bodyStyle: {
    type: Object,
    default: () => {
      return {}
    }
  }
})

const { t } = useI18n('views.robot.robot-config.basic-config.components.edit-box')

const handleEdit = (val) => {
  emit('edit')
  emit('update:isEdit', val)
}

const handleSave = () => {
  emit('save')
}
</script>
