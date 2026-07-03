<template>
  <div
    class="mini-card-preview"
    :class="{ selected: selected }"
    @click="$emit('click', card)"
  >
    <div class="card-title">{{ card.title }}</div>
    <div class="card-cover" v-if="card.thumb_url">
      <img :src="card.thumb_url" alt="" />
    </div>
    <div class="card-footer">
      <svg-icon name="mini-card-applet" style="font-size: 16px; color: #BFBFBF"></svg-icon>
      <span class="card-footer-text">小程序</span>
    </div>
    <!-- hover 操作面板 -->
    <div class="hover-actions" v-if="showEdit || showDelete">
      <div class="action-btn" v-if="showEdit" @click.stop="$emit('edit', card)">
        <svg-icon name="edit" style="font-size: 14px"></svg-icon>
      </div>
      <div class="action-btn" v-if="showDelete" @click.stop="$emit('delete', card)">
        <svg-icon name="delete" style="font-size: 14px"></svg-icon>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  card: {
    type: Object,
    default: () => ({})
  },
  selected: {
    type: Boolean,
    default: false
  },
  showEdit: {
    type: Boolean,
    default: false
  },
  showDelete: {
    type: Boolean,
    default: true
  }
})

defineEmits(['click', 'edit', 'delete'])
</script>

<style lang="less" scoped>
.mini-card-preview {
  position: relative;
  width: 188px;
  background: #fff;
  border: 1px solid #F0F0F0;
  border-radius: 8px;
  padding: 8px;
  cursor: pointer;
  flex-shrink: 0;
  transition: border-color 0.2s;

  &:hover {
    border-color: #2475FC;

    .hover-actions {
      opacity: 1;
      visibility: visible;
    }
  }

  &.selected {
    border: 2px solid #2475FC;
    padding: 7px;
  }
}

.card-title {
  font-size: 14px;
  color: #262626;
  line-height: 22px;
  width: 172px;
  margin-bottom: 8px;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
}

.card-cover {
  width: 172px;
  height: 138px;
  border-radius: 6px;
  overflow: hidden;
  margin-bottom: 8px;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.card-footer {
  display: flex;
  align-items: center;
  gap: 4px;
}

.card-footer-text {
  font-size: 12px;
  color: #BFBFBF;
  line-height: 20px;
}

.hover-actions {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: 8px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.85);
  border-radius: 6px;
  box-shadow: 0px 5px 5px -3px rgba(0, 0, 0, 0.1),
              0px 8px 10px 1px rgba(0, 0, 0, 0.06),
              0px 3px 14px 2px rgba(0, 0, 0, 0.05);
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.2s, visibility 0.2s;
}

.action-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #fff;
  border-radius: 2.5px;

  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }
}
</style>
