<style lang="less" scoped>
.public-library-layout {
  position: relative;
  display: flex;
  width: 100vw;
  height: 100vh;
  border-radius: 2px;
  background-color: #fff;
  .layout-left-wrapper{
    position: relative;
    height: 100%;

    .sidebar-toggle-btn {
      position: absolute;
      top: 16px;
      display: flex;
      align-items: center;
      justify-content: center;
      height: 32px;
      width: 16px;
      z-index: 100;
      border-radius: 0 8px 8px 0;
      background: #EDEFF2;
      cursor: pointer;
      transition: background 0.2s;

      .toggle-icon {
        font-size: 16px;
        color: #7A8699;
        transform: rotate(180deg);
      }

      &:hover {
        background: #A1A7B3;

        .toggle-icon{
          color: #333; 
        }
      }

      &.is-open{
        .toggle-icon{
          transform: rotate(0);
        }
      }
    }
  }
  .layout-left {
    height: 100%;
    overflow: hidden;

    .layout-left-content{
      display: flex;
      flex-flow: column nowrap;
      height: 100%;
      padding: 12px 8px;
      border-right: 1px solid #D9D9D9;
    }
  }

  .layout-body {
    flex: 1;
    height: 100%;
    overflow-x: hidden;
    overflow-y: auto;
  }

  .side-directory-wrapper {
    flex: 1;
    margin-top: 12px;
    overflow: hidden;
    border-top: 1px solid #f0f0f0;
  }
}
.no-doc {
  width: 100;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.drag-line{
  position: absolute;
  top: 0;
  bottom: 0;
  width: 5px;
  margin-left: -2px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;

  &:hover, &.dragging {
     &::after{
      width: 2px;
      background: #1890ff;
     }
  }
  &::after{
    display: block;
    content: '';
    width: 1px;
    height: 100%;
    background: none;
    cursor: col-resize;
    transition: background-color 0.2s;
  }
}
</style>

<template>
  <div class="public-library-layout">
    <div 
      v-show="!isSidebarCollapsed"
      class="drag-line" 
      :class="{ 'dragging': isDragging }" 
      :style="{ left: `${leftWidth}px` }" 
      @mousedown="handleMouseDown"
    ></div>
    <div class="layout-left-wrapper">
      <span class="sidebar-toggle-btn" :class="{'is-open': !isSidebarCollapsed}" :style="{left: `${isSidebarCollapsed ? collapsedWidth : leftWidth}px`}" @click="toggleSidebar">
        <svg-icon name="arrow-left" class="toggle-icon"></svg-icon>
      </span>
      <div class="layout-left" :style="{ width: `${isSidebarCollapsed ? collapsedWidth : leftWidth}px` }">
        <div class="layout-left-content">
          <SideMenus :menus="menus" :active="activeMenuKey" @menu-click="handleMeunClick" @handleAction="handleMeunAction"  v-if="!isSidebarCollapsed" />

          <div class="side-directory-wrapper" v-if="!isSidebarCollapsed">
            <SideDirectory
              ref="docTreeRef"
              :list="docList"
              :load-data="loadDocList"
              :draggable="docListDraggable"
              :selectedKeys="selectedKeys"
              :iconTemplateConfig="iconTemplateConfig"
              @importDoc="onImportDoc"
              @addDoc="onAddDoc"
              @renameDoc="onRenameDoc"
              @saveName="onSaveName"
              @deleteDoc="onDeleteDoc"
              @copyDoc="onCopyDoc"
              @copyLink="onShareDoc"
              @exportDoc="onExportDoc"
              @dragenter="onDragEnter"
              @drop="onDrop"
              @select="onSelectDoc"
            />
          </div>
        </div>
      </div>
    </div>
    
    <div class="layout-body">
      <router-view @changeDoc="onChangeDoc" />
    </div>
  </div>
</template>

<script setup>
import { OPEN_BOC_BASE_URL } from '@/constants/index'
import { saveAs } from 'file-saver'
import { ref, computed, watch, onMounted, createVNode, nextTick, provide } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { usePublicLibraryStore } from '@/stores/modules/public-library'
import { useCopyShareUrl } from '@/hooks/web/useCopyShareUrl'
import { useStorage } from '@/hooks/web/useStorage'
import { useI18n } from '@/hooks/web/useI18n'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import {
  saveDraftLibDoc,
  getDocList,
  deleteLibDoc,
  getLibDocInfo,
  updateDocSort,
  saveLibDoc,
} from '@/api/public-library'
import SideMenus from './components/side-menus.vue'
import SideDirectory from './components/side-directory.vue'

const { t } = useI18n('views.public-library.index')

const exportMd = (content, filename) => {
  const blob = new Blob([content], {
    type: 'text/markdown;charset=utf-8'
  })
  saveAs(blob, filename)
}

// 树形数据查找方法
const findDocInTree = (tree, targetId) => {
  {
    for (let i = 0; i < tree.length; i++) {
      {
        const node = tree[i]
        if (node.id == targetId) return node
        if (node.children) {
          {
            const found = findDocInTree(node.children, targetId)
            if (found) return found
          }
        }
      }
    }
    return null
  }
}

// 查找父级节点
const findParentDocInTree = (data, targetId, parent = null) => {
  for (let i = 0; i < data.length; i++) {
    const item = data[i]
    if (item.id == targetId) {
      return parent
    }
    if (item.children) {
      const found = findParentDocInTree(item.children, targetId, item)
      if (found) return found
    }
  }
  return null
}

function createDocName(baseName = t('untitled_document'), fileList = []) {
  // 提取所有已存在的文件名（兼容带扩展名的情况）
  const existingNames = new Set(fileList.map((file) => file.title))

  let counter = 1
  let fileName = `${baseName}(${counter})`

  // 动态递增检测（支持无扩展名场景）
  while (existingNames.has(fileName)) {
    counter++
    fileName = `${baseName}(${counter})`
  }

  return fileName
}
const libraryStore = usePublicLibraryStore()
const router = useRouter()
const route = useRoute()

const libraryId = computed(() => libraryStore.library_id)
const libraryKey = computed(() => libraryStore.library_key)
const iconTemplateConfig = computed(() => libraryStore.iconTemplateConfig)

const docId = computed(() => route.query.doc_id)
const docTreeRef = ref(null)
const docList = ref([])
const docListDraggable = ref(true)
const selectedKeys = ref([])
const currentDoc = ref(null)
const activeMenuKey = computed(() => route.meta.subActiveMenu || '')

// 拖拽相关状态
const { setStorage, getStorage } = useStorage('localStorage')
const SIDEBAR_WIDTH_KEY = 'public-library-sidebar-width'
const SIDEBAR_COLLAPSED_KEY = 'public-library-sidebar-collapsed'
const leftWidth = ref(getStorage(SIDEBAR_WIDTH_KEY) || 256)
const isSidebarCollapsed = ref(getStorage(SIDEBAR_COLLAPSED_KEY) || false)
const collapsedWidth = 0 // 收起时的宽度

const startX = ref(0)
const startWidth = ref(0)
const isDragging = ref(false)

// 将isDragging发送给子组件
provide('isDragging', isDragging)

const menus = [
  // {
  //   name: '知识库配置',
  //   key: 'config',
  //   icon: 'jichupeizhi',
  //   path: '/public-library/config',
  //   permissions: ['4']
  // },
  {
    name: t('home'),
    key: 'home',
    icon: 'ai',
    path: '/public-library/home',
    permissions: ['*']
  }
]

const handleMeunClick = (menu) => {
  router.push({
    path: menu.path,
    query: {
      library_id: libraryId.value,
      library_key: libraryKey.value
    }
  })
}

const handleMeunAction = (key, item) => {
  if(key === 'homePreviewStyle'){
    changeHomePreviewStyle(item)
  }
}

const changeHomePreviewStyle = (item) => {
  libraryStore.changeHomePreviewStyle()

  if(route.path !== '/public-library/home'){
    handleMeunClick(item)
  }
}

const getMyDocList = (doc) => {
  let data = {
    library_key: libraryKey.value,
    pid: doc ? doc.id : 0
  }

  return getDocList(data).then((res) => {
    let list = res.data || []
    list.forEach((item) => {
      item.id = item.id * 1
      item.pid = item.pid * 1
      item.key = item.id
      item.level = doc ? doc.level + 1 : 0
      item.sort = item.sort * 1
      item.children_num = item.children_num * 1
      item.children = []
      item.hasLoaded = false
      item.isLeaf = item.children_num == 0
    })

    if (!doc) {
      docList.value = list
    } else {
      doc.hasLoaded = true
      doc.children = list

      docList.value = [...docList.value]
    }
  })
}

const loadDocList = (treeNode) => {
  return new Promise((resolve) => {
    if (treeNode.dataRef.hasLoaded) {
      resolve()
      return
    }

    getMyDocList(treeNode.dataRef).then(() => {
      resolve()
    })
  })
}

const clearSelectedKeys = () => {
  selectedKeys.value = []
}

const { copyShareUrl } = useCopyShareUrl()

const onShareDoc = async (treeNode) => {
  const docUrl = OPEN_BOC_BASE_URL + '/doc/' + treeNode.doc_key

  await copyShareUrl(docUrl)
}

const onImportDoc = (treeNode, docFile) => {
  let isDir = 0;
  let title = docFile.title;
  // let level = treeNode ? treeNode.level + 1 : 0;
  let pid = treeNode ? treeNode.dataRef.id : 0;
  // let iconConfig = iconTemplateConfig.value.levels[level];
  // let doc_icon = isDir == 1 ? iconConfig.folder_icon : iconConfig.doc_icon;

  if (title.endsWith('.md')) {
    title = title.replace(/\.md$/i, '')
  }

  let data = {
    library_key: libraryKey.value,
    doc_id: '',
    doc_type: '4',
    doc_icon: '',
    pid: pid,
    title: title,
    content: docFile.content,
    sort: 0,
    is_dir: isDir,
  }

  saveDraftLibDoc(data).then((res) => {
    message.success(t('import_success'))
    data.id = res.data.doc_id
    data.sort = res.data.sort

    let node = treeNode ? treeNode.dataRef : null

    addDocSuccess(node, data)
  })
}

const onAddDoc = (treeNode, isDir) => {
  let docName = isDir == 1 ? t('new_folder') :  t('untitled_document');
  // let level = treeNode ? treeNode.level + 1 : 0;
  let pid = treeNode ? treeNode.dataRef.id : 0;
  // let iconConfig = iconTemplateConfig.value.levels[level];
  // let doc_icon = isDir == 1 ? iconConfig.folder_icon : iconConfig.doc_icon;

  if(isDir == 0){
    if (treeNode) {
      docName = createDocName(t('untitled_document'), treeNode.dataRef.children)
    } else {
      docName = createDocName(t('untitled_document'), docList.value)
    }
  }

  let data = {
    library_key: libraryKey.value,
    doc_id: '',
    doc_type: '4',
    doc_icon: '',
    pid: pid,
    title: docName,
    content: `# ${t('untitled_document')}`,
    sort: 0,
    is_dir: isDir,
  }

  let api = isDir == 1 ? saveLibDoc : saveDraftLibDoc

  api(data).then((res) => {
    message.success(t('add_success'))
    data.id = res.data.doc_id
    data.sort = res.data.sort

    let node = treeNode ? treeNode.dataRef : null

    addDocSuccess(node, data)
  })
}

const onCopyDoc = (treeNode) => {
  getLibDocInfo({
    library_key: libraryKey.value,
    doc_id: treeNode.dataRef.id,
  }).then((res) => {
    let data = {
      library_key: libraryKey.value,
      doc_id: '',
      doc_type: '4',
      pid: res.data.pid * 1,
      title: res.data.title + '(副本)',
      content: res.data.content,
      doc_icon: res.data.doc_icon || '',
      sort: 0,
      is_dir: 0,
    }

    saveDraftLibDoc(data).then((res) => {
      message.success(t('copy_success'))
      data.id = res.data.doc_id
      data.sort = res.data.sort

      let parentNode = findDocInTree(docList.value, data.pid)

      addDocSuccess(parentNode, data)
    })
  })
}

const addDocSuccess = (parentNode, data) => {
  data.children_num = 0
  data.children = []
  data.hasLoaded = false
  data.key = data.id
  data.isLeaf = true
  // data.sort = 0
  data.level = 0

  if (parentNode) {
    data.level = parentNode.level + 1
    parentNode.children_num++
    parentNode.isLeaf = false

    if (!parentNode.hasLoaded) {
      parentNode.hasLoaded = false
      nextTick(() => {
        docTreeRef.value.toggleExpand(parentNode.id, true)
      })
    } else {
      parentNode.children.unshift(data)
      nextTick(() => {
        docTreeRef.value.toggleExpand(parentNode.id, true)
      })
    }
  } else {
    data.level = 0
    docList.value.unshift(data)
  }

  onSelectDoc(data)
}

const onRenameDoc = (treeNode) => {
  let doc = findDocInTree(docList.value, treeNode.dataRef.id);
  docListDraggable.value = false;

  doc.isEdit = true;
}

const onSaveName = (treeNode, value) => {
  let doc = findDocInTree(docList.value, treeNode.dataRef.id);

  doc.title = value;
  doc.isEdit = false;

  setTimeout(() => {
    docListDraggable.value = true;
  }, 100)

  let data = {
    library_key: libraryKey.value,
    doc_id: doc.id,
    doc_type: '4',
    pid: doc.pid,
    title: value,
    content: doc.content,
    is_dir: doc.is_dir,
  }

  let api = doc.is_dir == 1 ? saveLibDoc : saveDraftLibDoc

  api(data).then(() => {
    message.success(t('modify_success'))
  })
 }

const onDeleteDoc = (treeNode) => {
  Modal.confirm({
    title: t('delete_document'),
    icon: createVNode(ExclamationCircleOutlined),
    content: t('confirm_delete_document', { title: treeNode.dataRef.title }),
    onOk() {
      submitDeleteDoc(treeNode)
    },
    onCancel() {}
  })
}

const submitDeleteDoc = (treeNode) => {
  let handelNodeId = treeNode.dataRef.id;
  let handleNodePid = treeNode.dataRef.pid;

  deleteLibDoc({
    library_key: libraryKey.value,
    doc_id: treeNode.dataRef.id
  }).then(() => {
    message.success(t('delete_success'))
    let parentDoc = findDocInTree(docList.value, handleNodePid)

    if (parentDoc) {
      let index = parentDoc.children.findIndex((item) => item.id == handelNodeId)
      parentDoc.children_num--
      parentDoc.isLeaf = parentDoc.children_num == 0
      parentDoc.children.splice(index, 1)

      if (activeMenuKey.value != 'doc' || currentDoc.value.id != handelNodeId) {
        return
      }

      let currentDocList = parentDoc.children.filter(item => item.is_dir != 1);

      if (currentDocList.length > 0) {
        onSelectDoc(currentDocList[0]);
      } else {
        let currentDocList = docList.value.filter(item => item.is_dir != 1);

        if (currentDocList.length > 0) {
          onSelectDoc(currentDocList[0]);
        } else {
          currentDoc.value = null;

          router.replace({
            path: '/public-library/home',
            query: {
              library_id: libraryId.value,
              library_key: libraryKey.value
            }
          })
        }
      }
    } else {
      let index = docList.value.findIndex((item) => item.id == handelNodeId)

      docList.value.splice(index, 1)

       if (activeMenuKey.value != 'doc' || currentDoc.value.id != handelNodeId) {
        return
      }

      let currentDocList = docList.value.filter(item => item.is_dir != 1);

      if (currentDocList.length > 0) {
        onSelectDoc(currentDocList[0])
      } else {
        currentDoc.value = null
        router.replace({
          path: '/public-library/home',
          query: {
            library_id: libraryId.value,
            library_key: libraryKey.value
          }
        })
      }
    }
  })
}

const onExportDoc = (treeNode) => {
  getLibDocInfo({
    library_key: libraryKey.value,
    doc_id: treeNode.dataRef.id
  }).then((res) => {
    let { content, title } = res.data

    if (!title.endsWith('.md')) {
      title += '.md'
    }

    exportMd(content, title)
  })
}

const onSelectDoc = (data) => {
  if(data.is_dir == 1){
    docTreeRef.value.toggleExpand(data.id)
    return
  }
  
  currentDoc.value = data
  selectedKeys.value = [data.id]
  if (route.path == '/public-library/editor') {
    router.replace({
      path: '/public-library/editor',
      query: {
        library_id: libraryId.value,
        library_key: libraryKey.value,
        doc_id: data.id
      }
    })
  } else {
    router.push({
      path: '/public-library/editor',
      query: {
        library_id: libraryId.value,
        library_key: libraryKey.value,
        doc_id: data.id
      }
    })
  }
}

const onDragEnter = () => {
  // console.log('onDragEnter', info.node)
  // expandedKeys 需要展开时
  // expandedKeys.value = info.expandedKeys;
}

const onDrop = (info) => {
  // console.log(info)
  const data = [...docList.value]
  const dropNode = info.node
  const dragNode = info.dragNode
  const dropKey = dropNode.key
  const dragKey = dragNode.key
  const dropPos = dropNode.pos.split('-')
  const dropPosition = info.dropPosition - Number(dropPos[dropPos.length - 1])
 
   const loop = (data, key, level = 0, callback) => {
     for (let i = 0; i < data.length; i++) {
       if (data[i].key === key) {
         callback(data[i], i, level, data)
         return
       }
 
       if (data[i].children) {
         loop(data[i].children, key, level + 1, callback)
       }
     }
   }

  // Find dragObject
  let dragObj

  loop(data, dragKey, 0, (item, index, level, arr) => {
    dragObj = { ...item }
    arr.splice(index, 1)
  })

  // 父级的children_num先减1
  if (dragNode.parent && dragNode.parent.node.children_num > 0) {
    dragNode.parent.node.children_num--
    dragNode.parent.node.isLeaf = dragNode.parent.node.children_num == 0
  }

  // 判断是不是放置在组上
  if (!info.dropToGap) {
    // Drop on the content
    loop(data, dropKey, 0, (item, index, level) => {
      item.children = item.children || []
      item.children_num++
      item.isLeaf = false
      dragObj.level = level + 1

      // where to insert 示例添加到头部，可以是随意位置
      item.children.unshift(dragObj)
    })
  } else {
    let ar = []
    let i = 0

    loop(data, dropKey, 0, (_item, index, level, arr) => {
      ar = arr
      i = index
      dragObj.level = level
    })

    if (dropPosition === -1) {
      ar.splice(i, 0, dragObj)
    } else {
      ar.splice(i + 1, 0, dragObj)
    }
  }

  docList.value = data

  handleDocSort(dragNode, dropNode)
}

const handleDocSort = (dragNode, dropNode) => {
  let parentDoc = findParentDocInTree(docList.value, dragNode.dataRef.id)

  dragNode.dataRef.pid = parentDoc ? parentDoc.id : 0

  let data = {
    library_key: libraryKey.value,
    doc_id: dragNode.dataRef.id,
    pid: dragNode.dataRef.pid,
    level: parentDoc ? parentDoc.level + 1 : 0,
    sort: 0
  }

  let list = parentDoc ? parentDoc.children : docList.value
  let dragNodeIndex = -1;

  if (list.length > 1) {
    dragNodeIndex = list.findIndex((item) => item.id == dragNode.dataRef.id)
    // 第一个 top=1，取index+1的sort
    if (dragNodeIndex === 0) {
      data.top = 1
      data.sort = list[dragNodeIndex + 1].sort
    } else {
      //其它情况 top 不传，取index-1的sort
      data.sort = list[dragNodeIndex - 1].sort
    }
  }

  updateDocSort(data).then((res) => {
    dragNode.dataRef.sort = res.data.sort
    dragNode.dataRef.level = data.level
 
     if (dropNode.pid == data.pid) {
       return
     }
    // 如果拖到到组内则展开组
    if (dropNode.dataRef.children.length >= dropNode.dataRef.children_num) {
      dropNode.dataRef.hasLoaded = true
      docTreeRef.value.addLoadedKey(dropNode.id)
      docTreeRef.value.toggleExpand(dropNode.id, true)
    } else {
      getMyDocList(dropNode.dataRef).then(() => {
        // 还原dragNode
        list[dragNodeIndex] = dragNode.dataRef

        nextTick(() => {
          dropNode.dataRef.hasLoaded = true
          docTreeRef.value.toggleExpand(dropNode.id, true)
        })
      })
    }
  })
}

const onChangeDoc = (data) => {
  let doc = findDocInTree(docList.value, data.id)

  doc.title = data.title;
  doc.is_draft = data.is_draft;
  doc.doc_icon = data.doc_icon;

  currentDoc.value = doc;
}

// 拖拽功能实现
const handleMouseDown = (e) => {
  isDragging.value = true
  startX.value = e.clientX
  startWidth.value = leftWidth.value
  
  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'
}

const maxWidth = ref(800)
const minWidth = ref(220)

const handleMouseMove = (e) => {
  if (!isDragging.value) return
  
  const deltaX = e.clientX - startX.value
  const newWidth = startWidth.value + deltaX
  
  // 限制最小和最大宽度
  if (newWidth >= minWidth.value && newWidth <= maxWidth.value) {

    leftWidth.value = newWidth
  }
}

const handleMouseUp = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', handleMouseMove)
  document.removeEventListener('mouseup', handleMouseUp)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
  
  // 保存拖拽后的宽度到本地存储
  setStorage(SIDEBAR_WIDTH_KEY, leftWidth.value)
}

// 侧边栏展开收起功能
const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value
  // 保存收起状态到本地存储
  setStorage(SIDEBAR_COLLAPSED_KEY, isSidebarCollapsed.value)
}

watch(
  () => route.path,
  (newPath) => {
    if (newPath !== '/public-library/editor') {
      clearSelectedKeys()
      currentDoc.value = null
    }
  }
)

watch(() => docList.value, (newVal) => { 
  libraryStore.setDocTreeState(JSON.parse(JSON.stringify(newVal)));
}, {
  deep: true
})

watch(docId, (newVal) => {
  if (newVal) {
    selectedKeys.value = [Number(docId.value)]
  }
})

watch(libraryId, () => {
  clearSelectedKeys();
  currentDoc.value = null;
  getMyDocList();
})

onMounted(() => {
  getMyDocList()

  if (docId.value) {
    selectedKeys.value = [Number(docId.value)]
  }
})
</script>
