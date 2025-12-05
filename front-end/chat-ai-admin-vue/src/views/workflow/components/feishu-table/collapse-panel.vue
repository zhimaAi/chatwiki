<style lang="less" scoped>
.collapse-panel-warpper {
  position: relative;
  width: 100%;
  overflow: hidden;
}

.collapse-panel-header {
  display: flex;
}

.collapse-panel-title {
  display: flex;
  align-items: center;
  height: 24px;
  padding: 0 4px;
  cursor: pointer;
  color: #1d2129;
  transition: all 0.3s ease;

  &:hover {
    border-radius: 6px;
    background: #E4E6EB;
  }
  
  .title-text {
    display: flex;
    align-items: center;
    font-size: 14px;
    font-weight: 400;
    color: #262626;
  }
}

.collapse-panel-body {
  transition: height 0.3s ease;
  overflow: hidden;
  margin-top: 8px;
}

.arrow-icon {
  transition: transform 0.3s;
}

.arrow-icon.collapsed {
  transform: rotate(-90deg);
}
</style>

<template>
  <div class="collapse-panel-warpper">
    <div class="collapse-panel-header" >
      <div class="collapse-panel-title" @click.stop="toggleCollapse">
        <span class="title-text">{{ title }} <template v-if="count">（{{ count }}）</template></span>
        <svg-icon 
          name="arrow-down" 
          class="arrow-icon" 
          :class="{ 'collapsed': collapsed }" 
          style="font-size: 16px;color: #000;"
        ></svg-icon>
      </div>
    </div>

    <div class="collapse-panel-body" v-show="!collapsed">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

// 定义props
const props = defineProps({
  title: {
    type: String,
  },
  count: {
    type: [Number, String],
    default: 0
  },
  defaultCollapsed: {
    type: Boolean,
    default: false
  }
});

// 折叠状态
const collapsed = ref(props.defaultCollapsed);

// 切换折叠状态
const toggleCollapse = () => {
  collapsed.value = !collapsed.value;
};
</script>
