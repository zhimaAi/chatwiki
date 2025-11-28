<template>
  <div class="top-menu-box">
    <a-menu
      class="top-menu"
      :items="items"
      v-model:selectedKeys="selectedKeys"
      mode="horizontal"
      @click="handleChangeMenu"
    >
      <template #title> 11 </template>
    </a-menu>
  </div>
</template>

<script setup>
import { ref, h } from 'vue'

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
    label: '嵌入网站',
    title: '嵌入网站'
  },
  {
    key: 'WeChatOfficialAccount',
    id: 'WeChatOfficialAccount',
    label: h('div', { class: 'tag-text-box' }, [
      h('span', {}, '微信公众号'),
      h('span', { class: 'text-xs' }, '支持未认证')
    ])
    // title: '微信公众号'
  },
  {
    key: 'WeChatMiniProgram',
    id: 'WeChatMiniProgram',
    label: '微信小程序',
    title: '微信小程序'
  },
  {
    key: 'WeChatCustomerService',
    id: 'WeChatCustomerService',
    label: '微信客服',
    title: '微信客服'
  },
  {
    key: 'FeishuRobot',
    id: 'FeishuRobot',
    label: '飞书机器人',
    title: '飞书机器人'
  },
  {
    key: 'DingDingRobot',
    id: 'DingDingRobot',
    label: '钉钉机器人',
    title: '钉钉机器人'
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
