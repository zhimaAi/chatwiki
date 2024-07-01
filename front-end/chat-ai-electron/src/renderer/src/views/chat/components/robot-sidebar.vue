<style lang="less" scoped>
.robot-sidebar {
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  width: 178px;
  padding-bottom: 20px;
  background-color: #d8e4f3;

  .logo-box {
    width: 52px;
    margin: 0 auto 8px auto;
  }

  .line {
    display: block;
    width: 60px;
    height: 1px;
    margin: 0 auto;
    border-radius: 1px;
    background: #cfd7e5;
    flex-shrink: 0;
  }
  .robot-list-wrapper {
    flex: 1;
    padding-top: 8px;
    overflow: hidden;
  }
  .robot-list {
    height: 100%;
    width: 100%;
    padding: 0 8px;
    overflow-x: hidden;
    overflow-y: auto;
    &::-webkit-scrollbar {
      display: none;
    }

    .robot-item {
      display: flex;
      align-items: center;
      justify-content: start;
      width: 100%;
      height: 40px;
      margin-bottom: 4px;
      padding: 0 8px;
      flex-shrink: 0;
      border-radius: 5px;
      cursor: pointer;
      &:hover {
        background: #ccd7e5;
      }
      &.active {
        background-color: #fff;
      }

      .robot-avatar {
        width: 24px;
        margin-right: 4px;
      }

      .robot-name {
        flex: 1;
        font-size: 14px;
        line-height: 22px;
        color: #242933;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    }
  }
}
</style>

<template>
  <div class="robot-sidebar">
    <div class="logo-box" style="display: none">
      <img class="logo" src="@/assets/logo2.svg" alt="" />
    </div>
    <i class="line" style="display: none"></i>
    <div class="robot-list-wrapper">
      <div class="robot-list">
        <div v-for="robot in props.robotList" :key="robot.id" class="robot-item-warp">
          <a-tooltip placement="right">
            <template #title>
              <span>{{ robot.robot_name }}</span>
            </template>
            <div
              class="robot-item"
              :class="{ active: robot.id === props.value }"
              @click="handleClick(robot)"
            >
              <img class="robot-avatar" :src="robot.avatar" alt="" />
              <div class="robot-name">{{ robot.robot_name }}</div>
            </div>
          </a-tooltip>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const emit = defineEmits(['change'])
const props = defineProps({
  robotList: {
    type: Array,
    default: () => []
  },
  value: {
    type: [String, Number],
    default: ''
  }
})

const handleClick = (item) => {
  emit('change', item)
}
</script>
