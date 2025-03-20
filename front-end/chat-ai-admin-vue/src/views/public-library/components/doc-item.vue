<template>
  <div
    class="doc-item"
    :class="{ active: doc.isActive, 'is-expanded': props.doc.isExpanded }"
    :data-pid="props.doc.pid"
  >
    <a-tooltip v-if="hasChildren">
      <template #title>收起</template>
      <span class="action-btn toggle-btn" @click="toggleExpand">
        <CaretRightOutlined style="" />
      </span>
    </a-tooltip>
    <div class="doc-item-body">
      <svg-icon class="doc-icon" name="doc-file"></svg-icon>
      <div class="doc-name">{{ props.doc.title }}</div>
    </div>
    <div class="doc-item-right">
      <!-- action 1 -->
      <a-dropdown>
        <template #overlay>
          <a-menu @click="handleMenuClick">
            <a-menu-item key="addDoc" @click="handleAddDoc(props.doc.id)"> 添加文档 </a-menu-item>
            <a-menu-item key="importDoc" @click="handleExportDoc(props.doc.id)">
              导入文档
            </a-menu-item>
          </a-menu>
        </template>
        <span class="action-btn"><PlusOutlined /></span>
      </a-dropdown>
      <!-- action2 -->
      <a-dropdown>
        <template #overlay>
          <a-menu @click="handleMenuClick">
            <a-menu-item key="copyLink"> 复制链接 </a-menu-item>
            <a-menu-item key="copyDoc"> 复制文档 </a-menu-item>
            <a-menu-item key="exportDoc"> 导出文档 </a-menu-item>
            <a-menu-item key="deleteDoc"> 删除文档 </a-menu-item>
          </a-menu>
        </template>
        <span class="action-btn"><MoreOutlined /></span>
      </a-dropdown>
    </div>
  </div>
  <!-- 如果有子文档，递归渲染 -->
  <div v-if="props.doc.children && props.doc.children.length > 0" class="directory-items">
    <DocItem
      v-for="child in doc.children"
      :key="child.id"
      :doc="child"
      @add-doc="addDoc"
      @export-doc="exportDoc"
    />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { PlusOutlined, CaretRightOutlined, MoreOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['add-doc', 'export-doc', 'toggle-expand'])

const props = defineProps({
  doc: {
    type: Object,
    required: true
  }
})

// 判断是否有子节点
const hasChildren = computed(() => props.doc.children_num > 0)

// 展开/收起切换
const toggleExpand = () => {
  if (hasChildren.value) {
    emit('toggle-expand', props.doc)
  }
}

const addDoc = (val) => {
  emit('add-doc', val)
}

const exportDoc = (val) => {
  emit('export-doc', val)
}

const handleMenuClick = () => {}
</script>

<style lang="less" scoped>
.directory-items {
  padding-left: 16px;
}
.doc-item {
  position: relative;
  display: flex;
  align-items: center;
  overflow: hidden;
  height: 32px;
  padding: 0 8px;
  padding-left: 28px;
  margin-bottom: 4px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  &.active,
  &:hover {
    background: #f2f4f7;
  }

  .toggle-btn {
    left: 0;
    position: absolute;
  }

  &.is-expanded {
    .toggle-btn {
      transform: rotate(90deg);
    }
  }

  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    margin-left: 4px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 16px;
    color: #595959;

    &:hover {
      background-color: #e4e6eb;
    }
  }

  .doc-item-body {
    flex: 1;
    display: flex;
    align-items: center;
    overflow: hidden;
    padding-left: 4px;
  }

  .doc-item-right {
    display: flex;
    align-items: center;
  }

  .doc-icon,
  .doc-name {
    line-height: 22px;
    font-size: 14px;
    color: #595959;
  }
  .doc-name {
    flex: 1;
    padding-left: 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
