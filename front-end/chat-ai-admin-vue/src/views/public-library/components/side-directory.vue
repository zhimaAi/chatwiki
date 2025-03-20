<style lang="less" scoped>
.side-directory {
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  width: 254px;
  overflow: hidden;

  .directory-header {
    display: flex;
    align-items: center;
    height: 40px;
    padding: 0 8px 0 24px;

    .label-text {
      flex: 1;
      padding-left: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #595959;
    }

    .action-box {
      display: flex;
      align-items: center;
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

        &:hover {
          background-color: #e4e6eb;
        }
      }
    }
  }
}
.directory-body {
  flex: 1;
  width: 280px;
  overflow: hidden;
  overflow-y: auto;
}
.directory-content {
  width: 254px;
  margin-top: 7px;
  padding: 0 8px;

  .toggle-btn {
    opacity: 1;
    font-size: 12px;
  }

  .action-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    line-height: 24px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 16px;
    color: #595959;

    &:hover {
      background-color: #e4e6eb;
    }
  }

  ::v-deep(.ant-tree-treenode) {
    align-items: center;
    padding: 0 8px 0 10px;
    margin-bottom: 4px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #f2f4f7;
    }

    &.ant-tree-treenode-selected {
      background: #e6efff;
    }

    .ant-tree-node-content-wrapper {
      position: relative;
      padding: 0;
      min-height: 32px;
      height: 32px;
      flex: 1 !important;
      display: flex;
      align-items: center;

      &:hover {
        background: none;
      }

      .ant-tree-title {
        display: block;
        position: relative;
        height: 100%;
        width: 100%;
        overflow: hidden;
      }
    }

    .ant-tree-switcher {
      width: 24px;
      height: 24px;
      line-height: 24px;
      margin-right: 0;
      align-self: unset;
    }
    .ant-tree-node-selected {
      background: none;

      .doc-item {
        .doc-title,
        .doc-icon {
          color: #2475fc;
        }
      }
    }
  }

  .doc-item {
    position: absolute;
    left: 0;
    top: 0;
    right: 0;
    bottom: 0;
    display: flex;
    align-items: center;
    width: 100%;
    overflow: hidden;
    box-sizing: border-box;

    .doc-icon {
      margin-right: 4px;
      font-size: 16px;
      color: #595959;
    }
    .doc-title {
      flex: 1;
      height: 32px;
      line-height: 32px;
      padding-right: 4px;
      font-size: 14px;
      color: #595959;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      .title-text {
        padding-left: 4px;
      }
    }

    .doc-action-box {
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
      opacity: 0;
      display: flex;
      align-items: center;
      height: 32px;
      transition: all 0.2s;

      .menu-name {
        color: #595959;
      }
    }

    .action-btn {
      opacity: 1;
    }

    &:hover {
      padding-right: 48px;

      .doc-action-box {
        opacity: 1 !important;
      }
    }
  }
}
</style>

<template>
  <div class="side-directory">
    <div class="directory-header">
      <svg-icon name="doc-directory" style="font-size: 16px; color: #595959"></svg-icon>
      <div class="label-text">文档目录</div>
      <div class="action-box">
        <a-dropdown>
          <template #overlay>
            <a-menu>
              <a-menu-item key="addDoc" @click="addDoc()"> 添加文档 </a-menu-item>
              <a-menu-item key="importDoc" @click="importDoc()"> 导入文档 </a-menu-item>
            </a-menu>
          </template>
          <span class="action-btn"><PlusOutlined style="color: #595959" /></span>
        </a-dropdown>

        <a-tooltip>
          <template #title>收起全部文档</template>
          <span class="action-btn" @click="toggleExpandAll(false)">
            <svg-icon name="menu-expand" style="font-size: 16px; color: #595959"></svg-icon>
          </span>
        </a-tooltip>
      </div>
    </div>

    <div class="directory-body">
      <div class="directory-content">
        <a-tree
          class="directory-tree"
          v-model:expandedKeys="expandedKeys"
          :selected-keys="selectedKeys"
          :loadedKeys="loadedKeys"
          block-node
          :tree-data="props.list"
          :load-data="loadData"
          :field-names="{ key: 'id', title: 'title', children: 'children' }"
          draggable
          @dragenter="onDragEnter"
          @drop="onDrop"
          @select="onSelect"
        >
          <template #switcherIcon="{ switcherCls }">
            <div style="margin-right: 4px">
              <span class="action-btn toggle-btn">
                <CaretDownOutlined style="font-size: 12px" :class="switcherCls"
              /></span>
            </div>
          </template>
          <template #title="treeNode">
            <div
              class="doc-item"
              :data-level="treeNode.level"
              :class="{ 'is-leaf': treeNode.dataRef.isLeaf }"
            >
              <svg-icon class="doc-icon" name="doc-file"></svg-icon>
              <div class="doc-title">{{ treeNode.dataRef.title }}</div>
              <div class="doc-action-box">
                <!-- action1 -->
                <a-dropdown v-if="treeNode.level < 2">
                  <template #overlay>
                    <a-menu>
                      <a-menu-item key="addDoc" @click.stop.prevent="addDoc(treeNode)">
                        <span class="menu-name">添加文档</span>
                      </a-menu-item>
                      <a-menu-item key="importDoc" @click.stop.prevent="importDoc(treeNode)">
                        <span class="menu-name">导入文档</span>
                      </a-menu-item>
                    </a-menu>
                  </template>
                  <span class="action-btn" @click.stop=""><PlusOutlined /></span>
                </a-dropdown>
                <!-- action2 -->
                <a-dropdown>
                  <template #overlay>
                    <a-menu>
                      <a-menu-item key="copyLink" @click.stop.prevent="handleCopyLink(treeNode)">
                        复制链接<span class="menu-name"></span>
                      </a-menu-item>
                      <a-menu-item key="copyDoc" @click.stop.prevent="handleCopyDoc(treeNode)">
                        <span class="menu-name">复制文档</span>
                      </a-menu-item>
                      <a-menu-item key="exportDoc" @click.stop.prevent="handleExportDoc(treeNode)">
                        <span class="menu-name">导出文档</span>
                      </a-menu-item>
                      <a-menu-item key="deleteDoc" @click.stop.prevent="handleDeleteDoc(treeNode)">
                        <span class="menu-name">删除文档</span>
                      </a-menu-item>
                    </a-menu>
                  </template>
                  <span class="action-btn" @click.stop=""><MoreOutlined /></span>
                </a-dropdown>
              </div>
            </div>
          </template>
        </a-tree>
      </div>
    </div>

    <UploadDoc ref="uploadDocRef" @file-loaded="handleFileLoaded" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { PlusOutlined, CaretDownOutlined, MoreOutlined } from '@ant-design/icons-vue'
import UploadDoc from './upload-doc.vue'

const emit = defineEmits([
  'addDoc',
  'importDoc',
  'copyLink',
  'copyDoc',
  'deleteDoc',
  'exportDoc',
  'dragenter',
  'drop',
  'select'
])

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  },
  loadData: {
    type: Function,
    default: () => {}
  },
  selectedKeys: {
    type: Array,
    default: () => []
  }
})

const expandedKeys = ref([])
const loadedKeys = ref([])

const addLoadedKey = (key) => {
  if (!loadedKeys.value.includes(key)) {
    loadedKeys.value.push(key)
  }
}

const toggleExpand = (key, type) => {
  const index = expandedKeys.value.indexOf(key)

  if (typeof type === 'boolean') {
    if (index > -1 && type === false) {
      expandedKeys.value.splice(index, 1)
    } else {
      expandedKeys.value.push(key)
    }
  } else {
    if (index > -1) {
      expandedKeys.value.splice(index, 1)
    } else {
      expandedKeys.value.push(key)
    }
  }
}

const toggleExpandAll = (type) => {
  if (type === true) {
    props.list.forEach((item) => {
      if (!item.isLeaf && !expandedKeys.value.includes(item.key)) {
        expandedKeys.value.push(item.key)
      }
    })
  } else {
    expandedKeys.value = []
  }
}

const uploadDocRef = ref(null)

const addDoc = (node) => {
  emit('addDoc', node)
}

let currentNode = null

const importDoc = (node) => {
  currentNode = node
  uploadDocRef.value.triggerUpload()
}

const handleFileLoaded = (doc) => {
  emit('importDoc', currentNode, doc)
}

const handleCopyLink = (node) => {
  emit('copyLink', node)
}

const handleCopyDoc = (node) => {
  emit('copyDoc', node)
}

const handleExportDoc = (node) => {
  emit('exportDoc', node)
}

const handleDeleteDoc = (node) => {
  emit('deleteDoc', node)
}

const onDragEnter = (info) => {
  emit('dragenter', info)
}

const onDrop = (info) => {
  emit('drop', info)
}

const onSelect = (selectedKeys, e) => {
  // selectedKeys.value = [e.node.dataRef.key]
  // console.log(selectedKeys)
  // let key = e.node.dataRef.key
  // console.log('key: ', key)
  // selectedKeys.value = []
  // nextTick(() => {
  //   selectedKeys.value = [key]
  // })
  emit('select', e.node.dataRef)
}

defineExpose({
  toggleExpand,
  addLoadedKey
})
</script>
