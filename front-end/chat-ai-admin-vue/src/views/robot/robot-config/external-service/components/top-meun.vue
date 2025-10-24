<template>
  <div class="top-menu-box">
    <a-menu
      class="top-menu"
      :items="items"
      v-model:selectedKeys="selectedKeys"
      mode="horizontal"
      @click="handleChangeMenu"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'

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
    label: '微信公众号',
    title: '微信公众号'
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
}
</style>
