<template>
  <div class="top-menu-box">
    <a-menu
      class="top-menu"
      :items="items"
      v-model:selectedKeys="selectedKeys"
      mode="horizontal"
      @click="handleChangeMenu"
    >
    </a-menu>
  </div>
</template>

<script setup>
import { ref, h } from 'vue'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.robot.robot-config.external-service.components.top-meun')

const emit = defineEmits(['change', 'update:value'])
const props = defineProps({
  value: {
    type: String,
    default: ''
  }
})

const selectedKeys = ref([props.value])

const items = ref([
  {
    key: 'WebAPP',
    id: 'WebAPP',
    label: 'WebAPP',
    title: 'WebAPP'
  },
  {
    key: 'EmbedWebsite',
    id: 'EmbedWebsite',
    label: t('menu_embed_website'),
    title: t('menu_embed_website')
  },
  {
    key: 'WeChatOfficialAccount',
    id: 'WeChatOfficialAccount',
    label: h('div', { class: 'tag-text-box' }, [
      h('span', {}, t('menu_wechat_official_account')),
      h('span', { class: 'text-xs' }, t('tip_unverified_supported'))
    ])
  },
  {
    key: 'WeChatMiniProgram',
    id: 'WeChatMiniProgram',
    label: t('menu_wechat_mini_program'),
    title: t('menu_wechat_mini_program')
  },
  {
    key: 'WeChatCustomerService',
    id: 'WeChatCustomerService',
    label: t('menu_wechat_customer_service'),
    title: t('menu_wechat_customer_service')
  },
  {
    key: 'FeishuRobot',
    id: 'FeishuRobot',
    label: t('menu_feishu_robot'),
    title: t('menu_feishu_robot')
  },
  {
    key: 'DingDingRobot',
    id: 'DingDingRobot',
    label: t('menu_dingtalk_robot'),
    title: t('menu_dingtalk_robot')
  }
])

const handleChangeMenu = ({ item }) => {
  if (selectedKeys.value.includes(item.id)) {
    return
  }

  emit('update:value', item.id)
  emit('change', item)
}
</script>

<style lang="less" scoped>
.top-menu-box {
  .top-menu {
    border-right: 0;

    ::v-deep(.menu-icon) {
      color: #a1a7b3;
      font-size: 16px;
      vertical-align: -3px;
    }

    ::v-deep(.ant-menu-item-selected .menu-icon) {
      color: #2475fc;
    }
  }

  ::v-deep(.tag-text-box) {
    display: flex;
    gap: 8px;
    align-items: center;
    .text-xs {
      display: flex;
      align-items: center;
      padding: 0 6px;
      height: 18px;
      width: fit-content;
      border-radius: 6px;
      background: #fb363f;
      color: #ffffff;
      font-size: 12px;
    }
  }
}
</style>
