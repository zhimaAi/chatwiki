<style lang="less" scoped>
.wechat-app-item {
  width: 355px;
  height: fit-content;
  padding: 24px;
  border-radius: 12px;
  border: 1px solid #e4e6eb;
  overflow: hidden;
  cursor: pointer;
  position: relative;
  .app-info-block {
    display: flex;
    align-items: center;
  }

  &:hover {
    box-shadow:
      0 3px 14px 2px #0000000d,
      0 8px 10px 1px #0000000f,
      0 5px 5px -3px #0000001a;
  }

  .item-body {
    display: flex;
    align-items: center;
    overflow: hidden;
  }
  .app-info {
    flex: 1;
    overflow: hidden;
  }
  .app-avatar {
    display: block;
    width: 62px;
    height: 62px;
    border-radius: 16px;
    margin-right: 12px;
    border: 1px solid var(--07, #f0f0f0);
  }

  .app-name,
  .app-desc {
    width: 100%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .app-name {
    line-height: 24px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }

  .app-desc {
    line-height: 20px;
    font-size: 12px;
    font-weight: 400;
    color: #8c8c8c;
    margin-top: 4px;
  }

  .ext-info-list {
    margin-top: 12px;
    display: flex;
    flex-wrap: wrap;
    color: #595959;
    line-height: 22px;
    font-size: 14px;
    .status-block {
      display: flex;
      align-items: center;
      border-radius: 6px;
      padding: 0 6px;
      font-size: 14px;
      line-height: 22px;
      gap: 2px;
      &.status-success {
        background: #cafce4;
        color: #21a665;
      }
      &.status-warning {
        background: #fae4dc;
        color: #ed744a;
      }
    }
  }

  .btn-wrapper-box {
    position: absolute;
    width: 24px;
    height: 24px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    right: 24px;
    bottom: 20px;
    transition: all 0.2s ease-in-out;
    &:hover {
      background: #e4e6eb;
    }
  }
}
</style>

<template>
  <div class="wechat-app-item">
    <div class="item-body">
      <img class="app-avatar" :src="item.app_avatar" alt="" />
      <div class="app-info">
        <div class="app-name">{{ item.app_name }}</div>
        <div class="app-desc">{{ t('appid') }}{{ item.app_id }}</div>
      </div>
    </div>
    <template v-if="showExtTypeList.includes(props.app_type)">
      <div class="ext-info-list">
        <div class="ext-content">{{ item.wechat_reply_type }}</div>
      </div>
      <div class="ext-info-list">
        <div class="status-block status-success" v-if="item.account_is_verify == 'true'">
          <CheckCircleFilled />
          {{ t('verified') }}
        </div>
        <a-flex :gap="8" v-else>
          <div class="status-block status-warning">
            <ExclamationCircleFilled />
            {{ t('unverified') }}
          </div>
          <a @click="handleRefresh">{{ t('refresh_status') }}</a>
        </a-flex>
      </div>
    </template>

    <a-dropdown v-if="showMenu && showMenu.length">
      <div class="btn-wrapper-box">
        <EllipsisOutlined />
      </div>
      <template #overlay>
        <a-menu>
          <a-menu-item v-if="showMenu.includes('edit')" @click="handleEdit">
            <a>{{ t('edit') }}</a>
          </a-menu-item>
          <a-menu-item v-if="showMenu.includes('del')" @click="handleDelete">
            <a style="color: #fb363f">{{ t('delete') }}</a>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup>
import { CheckCircleFilled, ExclamationCircleFilled, EllipsisOutlined } from '@ant-design/icons-vue'
import { useI18n } from '@/hooks/web/useI18n'
const emit = defineEmits(['edit', 'delete', 'refresh'])
const showExtTypeList = ['mini_program', 'official_account']

const { t } = useI18n('views.robot.robot-config.external-service.components.wechat-app-item')

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  },
  app_type: {
    type: String,
    default: ''
  },
  showMenu: {
    type: Array,
    default: ['edit', 'del']
  }
})

function handleEdit() {
  emit('edit', { ...props.item })
}

function handleDelete() {
  emit('delete', { ...props.item })
}

function handleRefresh() {
  emit('refresh', { ...props.item })
}
</script>
