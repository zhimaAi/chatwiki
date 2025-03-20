<style lang="less" scoped>
.side-meuns {
  padding: 0 8px;

  .side-meun-item {
    display: flex;
    align-items: center;
    padding: 9px 16px;
    line-height: 22px;
    border-radius: 6px;
    cursor: pointer;

    .menu-icon {
      font-size: 14px;
      color: #a1a7b2;
    }

    .menu-name {
      padding-left: 8px;
      font-size: 14px;
      color: rgb(89, 89, 89);
    }

    &.active {
      background: #e5efff;

      .menu-icon {
        color: rgb(36, 117, 252);
      }

      .menu-name {
        color: rgb(36, 117, 252);
      }
    }
  }
}
</style>

<template>
  <div class="side-meuns">
    <template v-for="item in props.menus" :key="item.key">
      <div
        class="side-meun-item"
        :class="{ active: item.key === props.active }"
        @click="handdleMenu(item)"
        v-if="checkPermission(item.permissions)"
      >
        <svg-icon class="menu-icon" :name="item.icon"></svg-icon>
        <span class="menu-name">{{ item.name }}</span>
      </div>
    </template>
  </div>
</template>

<script setup>
import { usePublicLibraryStore } from '@/stores/modules/public-library'

const emit = defineEmits(['menuClick'])

const props = defineProps({
  menus: {
    type: Array,
    default: () => []
  },
  active: {
    type: String,
    default: ''
  }
})

const { checkPermission } = usePublicLibraryStore()

const handdleMenu = (item) => {
  emit('menuClick', item)
}
</script>
