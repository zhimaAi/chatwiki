<template>
  <a-popover v-if="!props.disabled" placement="top" overlay-class-name="custom-popover">
    <template #content>
      <div class="float-action-box">
        <div
          class="action-btn"
          data-type="edit"
          v-if="menus.includes('edit')"
          @click="handleActionClick('edit')"
        >
          <img class="action-icon" title="编辑" src="@/assets/img/edit.svg" alt="" />
        </div>

        <div
          class="action-btn"
          data-type="delete"
          v-if="menus.includes('delete')"
          @click="handleActionClick('delete')"
        >
          <img class="action-icon" title="删除" src="@/assets/img/delete.svg" alt="" />
        </div>
      </div>
    </template>
    <template #title>
      <span></span>
    </template>
    <slot name="default"></slot>
  </a-popover>
  <slot name="default" v-else></slot>
</template>

<script setup>
const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
  name: {
    type: String,
    default: '',
  },
  menus: {
    type: Array,
    default: () => ['edit'],
  },
})

const emit = defineEmits(['click'])

const handleActionClick = (type) => {
  emit('click', props.name, type)
}
</script>

<style lang="less" scoped>
.float-action-box {
  display: flex;
  border-radius: 8px;
  background-color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  padding: 4px;

  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-flow: column nowrap;
    width: 32px;
    height: 32px;
    border-radius: 4px;
    margin-left: 4px;
    transition: all 0.2s ease;
    cursor: pointer;

    &:first-child {
      margin-left: 0;
    }

    &:hover {
      background-color: #f5f5f5;
    }

    .action-icon {
      width: 18px;
      height: 18px;
    }
  }
}
</style>
