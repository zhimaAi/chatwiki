<style lang="less" scoped>
.wechat-app-item {
  width: 325px;
  height: 74px;
  display: flex;
  align-items: center;
  padding: 16px;
  border-radius: 2px;
  border: 1px solid #d9d9d9;
  overflow: hidden;
  cursor: pointer;

  &:hover {
    box-shadow:
      0 3px 14px 2px #0000000d,
      0 8px 10px 1px #0000000f,
      0 5px 5px -3px #0000001a;
  }

  .item-body {
    flex: 1;
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
    width: 40px;
    height: 40px;
    border-radius: 2px;
    margin-right: 8px;
  }

  .app-name,
  .app-desc {
    width: 100%;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
  }

  .app-name {
    line-height: 22px;
    font-size: 14px;
    font-weight: 600;
    color: #262626;
  }

  .app-desc {
    line-height: 20px;
    font-size: 12px;
    font-weight: 400;
    color: #8c8c8c;
  }

  .item-right {
    display: flex;
    gap: 8px;
  }
  .action-btn {
    color: #8c8c8c;
  }
  .edit-btn:hover {
    color: #2475fc;
  }
  .delete-btn:hover {
    color: #fb363f;
  }
}
</style>

<template>
  <div class="wechat-app-item">
    <div class="item-body">
      <img class="app-avatar" :src="item.app_avatar" alt="" />
      <div class="app-info">
        <div class="app-name">{{ item.app_name }}</div>
        <div class="app-desc">Appid：{{ item.app_id }}</div>
      </div>
    </div>
    <div class="item-right">
      <span class="action-btn edit-btn" @click="handleEdit"><svg-icon name="edit"></svg-icon></span>
      <span class="action-btn delete-btn" @click="handleDelete"
        ><svg-icon name="delete"></svg-icon
      ></span>
    </div>
  </div>
</template>

<script setup>
const emit = defineEmits(['edit', 'delete'])

const props = defineProps({
  item: {
    type: Object,
    default: () => ({})
  }
})

function handleEdit() {
  emit('edit', { ...props.item })
}

function handleDelete() {
  emit('delete', { ...props.item })
}
</script>
