<style lang="less" scoped>
.side-directory {
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  width: 100%;
  overflow: hidden;

  .directory-header {
    display: flex;
    align-items: center;
    flex-wrap: nowrap;
    width: 100%;
    overflow: hidden;
    height: 40px;
    padding: 0 8px 0 16px;

    .label-text {
      flex: 1;
      white-space: nowrap;
      padding-left: 8px;
      font-size: 14px;
      font-weight: 400;
      overflow: hidden;
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
  width: 100%;
  overflow: hidden;
  overflow-y: auto;

  /* 自定义滚动条样式 */
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: #d9d9d9;
    border-radius: 3px;
    transition: background 0.2s;

    &:hover {
      background: #bfbfbf;
    }
  }

  &::-webkit-scrollbar-thumb:active {
    background: #999999;
  }
}

.directory-content {
  width: 100%;
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
    padding: 0 8px 0 8px;
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

    &:hover.no-edit {
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
              <a-menu-item key="addDir" @click="addDoc(null, 0)"> 添加文档 </a-menu-item>
              <a-menu-item key="importDoc" @click="importDoc()"> 导入文档 </a-menu-item>
              <a-menu-item key="addDoc" @click="addDoc(null, 1)"> 添加文件夹 </a-menu-item>
            </a-menu>
          </template>
          <span class="action-btn">
            <PlusOutlined style="color: #595959" />
          </span>
        </a-dropdown>

        <a-tooltip>
          <template #title>收起全部文档</template>
          <span class="action-btn" @click="toggleExpandAll(false)">
            <svg-icon name="menu-expand" style="font-size: 16px; color: #595959"></svg-icon>
          </span>
        </a-tooltip>

        <a-tooltip>
          <template #title>图标模板设置</template>
          <span class="action-btn" @click="openIconTemplate">
            <svg-icon name="system-setting" style="font-size: 16px; color: #595959"></svg-icon>
          </span>
        </a-tooltip>
      </div>
    </div>

    <div class="directory-body">
      <div class="directory-content">
        <a-tree class="directory-tree" block-node v-model:expandedKeys="expandedKeys" :draggable="draggable" :selected-keys="selectedKeys"
          :loadedKeys="loadedKeys" :tree-data="props.list" :load-data="loadData"
          :field-names="{ key: 'id', title: 'title', children: 'children' }" :allowDrop="allowDrop"
          @dragenter="onDragEnter" @drop="onDrop" @select="onSelect">
          <template #switcherIcon="{ switcherCls }">
            <div style="margin-right: 4px">
              <span class="action-btn toggle-btn">
                <CaretDownOutlined style="font-size: 12px" :class="switcherCls" />
              </span>
            </div>
          </template>
          <template #title="treeNode">
            <div class="doc-item" :data-level="treeNode.level"
              :class="[{ 'is-leaf': treeNode.dataRef.isLeaf }, treeNode.dataRef.isEdit ? 'is-edit' : 'no-edit']">
              <span class="doc-icon" v-if="treeNode.data.doc_icon">{{ treeNode.data.doc_icon }}</span>
              <span class="doc-icon" v-else>{{ treeNode.data.is_dir == 1 ?
                iconTemplateConfig.levels[treeNode.level].folder_icon :
                iconTemplateConfig.levels[treeNode.level].doc_icon }}</span>

              <div class="doc-title">
                <span class="doc-title-text" v-if="!treeNode.dataRef.isEdit">{{ treeNode.dataRef.title }}</span>
                <a-input :class="[`rename-input-${treeNode.dataRef.id}`]" name="rename"
                  :defaultValue="treeNode.dataRef.title" placeholder="请输入标题" @blur="saveName(treeNode)"
                  @pressEnter="saveName(treeNode)" v-else />
              </div>

              <div class="doc-action-box" v-if="!treeNode.dataRef.isEdit">
                <!-- action1 -->
                <a-dropdown v-if="treeNode.dataRef.is_dir == 1">
                  <template #overlay>
                    <a-menu>
                      <a-menu-item key="addDir" @click.stop.prevent="addDoc(treeNode, 0)">
                        <span class="menu-name">添加文档</span>
                      </a-menu-item>
                      <a-menu-item key="importDoc" @click.stop.prevent="importDoc(treeNode)">
                        <span class="menu-name">导入文档</span>
                      </a-menu-item>
                      <a-menu-item key="addDoc" @click.stop.prevent="addDoc(treeNode, 1)">
                        <span class="menu-name">添加文件夹</span>
                      </a-menu-item>
                    </a-menu>
                  </template>
                  <span class="action-btn" @click.stop="">
                    <PlusOutlined />
                  </span>
                </a-dropdown>
                <!-- action2 -->
                <a-dropdown>
                  <template #overlay>
                    <a-menu>
                      <template v-if="treeNode.dataRef.is_dir == 1">
                        <a-menu-item key="rename" @click.stop.prevent="handleRenameDoc(treeNode)">
                          <span class="menu-name">重命名</span>
                        </a-menu-item>
                        <a-menu-item key="deleteDoc" @click.stop.prevent="handleDeleteDoc(treeNode)">
                          <span class="menu-name">删除</span>
                        </a-menu-item>
                      </template>
                      <template v-else>
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
                      </template>
                    </a-menu>
                  </template>
                  <span class="action-btn" @click.stop="">
                    <MoreOutlined />
                  </span>
                </a-dropdown>
              </div>
            </div>
          </template>
        </a-tree>
      </div>
    </div>

    <UploadDoc ref="uploadDocRef" @file-loaded="handleFileLoaded" />
    <IconTemplate ref="iconTemplateRef" @select="selectTemplate" />
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { PlusOutlined, CaretDownOutlined, MoreOutlined } from '@ant-design/icons-vue'
import UploadDoc from './upload-doc.vue'
import IconTemplate from './icon-template.vue'

const emit = defineEmits([
  'addDoc',
  'importDoc',
  'copyLink',
  'copyDoc',
  'deleteDoc',
  'exportDoc',
  'dragenter',
  'drop',
  'select',
  'renameDoc',
  'saveName',
])

const props = defineProps({
  list: {
    type: Array,
    default: () => []
  },
  loadData: {
    type: Function,
    default: () => { }
  },
  selectedKeys: {
    type: Array,
    default: () => []
  },
  iconTemplateConfig: {
    type: Object,
    default: () => ({})
  },
  draggable: {
    type: Boolean,
    default: true
  }
})

const debounce = (fn, delay) => {
  let timer = null
  return function (...args) {
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn.apply(this, args)
    }, delay)
  }
}

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

const addDoc = (node, isDir) => {
  emit('addDoc', node, isDir)
}

let currentNode = null

const importDoc = (node) => {
  currentNode = node
  uploadDocRef.value.triggerUpload()
}

const iconTemplateRef = ref(null)
const openIconTemplate = () => {
  iconTemplateRef.value.open()
}

const selectTemplate = (templateId) => {
  emit('selectTemplate', templateId)
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

const handleRenameDoc = (node) => {
  emit('renameDoc', node)

  nextTick(() => {
    let el = document.querySelector('.rename-input-' + node.data.id);

    if (el) {
      el.focus()
      el.select()
    }
  })
}

const saveName = debounce((node) => {
  let el = document.querySelector('.rename-input-' + node.data.id);
  if (!el) return
  let value = el.value.trim()

  emit('saveName', node, value)
}, 200)

const handleDeleteDoc = (node) => {
  emit('deleteDoc', node)
}

const allowDrop = ({ dropNode, dragNode, dropPosition }) => {
  // console.log({ dropNode, dragNode, dropPosition })
  if((dropNode.is_dir == 0 && dragNode.is_dir == 0) && dropPosition == 0){
    return false;
  } else if((dropNode.is_dir == 0 && dragNode.is_dir == 1)  && dropPosition == 0){
    return false;
  }

  return true;
}

const onDragEnter = (info) => {
  emit('dragenter', info)
}

const onDrop = (info) => {
  emit('drop', info)
}

const onSelect = (selectedKeys, e) => {
  emit('select', e.node.dataRef)
}

defineExpose({
  toggleExpand,
  addLoadedKey
})
</script>
