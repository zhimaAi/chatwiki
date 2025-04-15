<style lang="less" scoped>
.public-library-layout {
  display: flex;
  width: 100%;
  height: 100%;
  border-radius: 2px;
  background-color: #fff;

  .layout-left {
    display: flex;
    flex-flow: column nowrap;
    width: 256px;
    height: 100%;
    overflow: hidden;
    border-right: 1px solid #f0f0f0;
  }

  .layout-body {
    flex: 1;
    height: 100%;
    overflow-x: hidden;
    overflow-y: auto;
  }

  .library-title {
    display: flex;
    align-items: center;
    overflow: hidden;
    line-height: 22px;
    padding: 24px 24px 24px 24px;

    .library-logo {
      width: 32px;
      height: 32px;
      margin-right: 8px;
      border-radius: 4px;
    }

    .library-name {
      flex: 1;
      font-size: 14px;
      font-weight: 600;
      color: #262626;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
    }

    .publish-status {
      width: 72px;
      height: 22px;
      margin-left: 8px;

      .icon {
        width: 72px;
        height: 22px;
      }
    }
  }

  .side-directory-wrapper {
    flex: 1;
    margin-top: 16px;
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
</style>

<template>
  <div class="public-library-layout">
    <div class="layout-left">
      <div class="library-title">
        <img class="library-logo" :src="state.avatar" alt="" v-if="state.avatar" />
        <span class="library-name">{{ state.library_name }}</span>
        <span class="publish-status" v-if="currentDoc">
          <img
            class="icon"
            src="../../assets/img/editor/unpublished.svg"
            alt=""
            v-if="currentDoc.is_draft == 1"
          />
          <img
            class="icon"
            src="../../assets/img/editor/published .svg"
            alt=""
            v-if="currentDoc.is_draft == 0"
          />
        </span>
      </div>

      <SideMenus :menus="menus" :active="activeMenuKey" @menu-click="handleMeunClick" />

      <div class="side-directory-wrapper">
        <SideDirectory
          ref="docTreeRef"
          :list="docList"
          :load-data="loadDocList"
          :selectedKeys="selectedKeys"
          @importDoc="onImportDoc"
          @addDoc="onAddDoc"
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
    <div class="layout-body">
      <router-view @changeDoc="onChangeDoc" />
    </div>
    <ShareModal ref="shareModalRef" />
  </div>
</template>

<script setup>
import { LIBRARY_OPEN_AVATAR } from '@/constants/index'
import { saveAs } from 'file-saver'
import { ref, reactive, computed, watch, onMounted, createVNode, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { getLibraryInfo } from '@/api/library'
import {
  saveDraftLibDoc,
  getDocList,
  deleteLibDoc,
  getLibDocInfo,
  updateDocSort
} from '@/api/public-library'
import SideMenus from './components/side-menus.vue'
import SideDirectory from './components/side-directory.vue'
import ShareModal from './components/share-modal.vue'

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

function createDocName(baseName = '无标题文档', fileList = []) {
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

const router = useRouter()
const route = useRoute()

const libraryId = computed(() => route.query.library_id)

const state = reactive({
  type: 1,
  library_intro: '',
  library_name: '',
  library_key: '',
  avatar: ''
})

const docId = computed(() => route.query.doc_id)
const docTreeRef = ref(null)
const docList = ref([])
const selectedKeys = ref([])
const currentDoc = ref(null)
const activeMenuKey = computed(() => route.meta.activeMenu || '')

const menus = [
  {
    name: '知识库配置',
    key: 'config',
    icon: 'jichupeizhi',
    path: '/public-library/config',
    permissions: ['4']
  },
  {
    name: '首页',
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
      library_id: route.query.library_id,
      library_key: state.library_key
    }
  })
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

if (docId.value) {
  selectedKeys.value = [Number(docId.value)]
}
watch(docId, (newVal) => {
  if (newVal) {
    selectedKeys.value = [Number(docId.value)]
  }
})

const getMyDocList = (doc) => {
  let data = {
    library_key: state.library_key,
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

const shareModalRef = ref(null)
const onShareDoc = (treeNode) => {
  shareModalRef.value.show({ doc_key: treeNode.doc_key })
}

const onImportDoc = (treeNode, docFile) => {
  let title = docFile.title

  if (title.endsWith('.md')) {
    title = title.replace(/\.md$/i, '')
  }

  let data = {
    library_key: state.library_key,
    doc_id: '',
    doc_type: '4',
    pid: treeNode ? treeNode.dataRef.id : '0',
    title: title,
    content: docFile.content,
    sort: 0
  }

  saveDraftLibDoc(data).then((res) => {
    message.success('导入成功')
    data.id = res.data.doc_id
    data.sort = res.data.sort

    let node = treeNode ? treeNode.dataRef : null

    addDocSuccess(node, data)
  })
}

const onAddDoc = (treeNode) => {
  let docName = '无标题文档'

  if (treeNode) {
    docName = createDocName('无标题文档', treeNode.dataRef.children)
  } else {
    docName = createDocName('无标题文档', docList.value)
  }

  let data = {
    library_key: state.library_key,
    doc_id: '',
    doc_type: '4',
    pid: treeNode ? treeNode.dataRef.id : 0,
    title: docName,
    content: '# 无标题文档',
    sort: 0
  }

  saveDraftLibDoc(data).then((res) => {
    message.success('添加成功')
    data.id = res.data.doc_id
    data.sort = res.data.sort

    let node = treeNode ? treeNode.dataRef : null

    addDocSuccess(node, data)
  })
}

const onCopyDoc = (treeNode) => {
  getLibDocInfo({
    library_key: state.library_key,
    doc_id: treeNode.dataRef.id
  }).then((res) => {
    let data = {
      library_key: state.library_key,
      doc_id: '',
      doc_type: '4',
      pid: res.data.pid * 1,
      title: res.data.title + '(副本)',
      content: res.data.content,
      sort: 0
    }

    saveDraftLibDoc(data).then((res) => {
      message.success('复制成功')
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

const onDeleteDoc = (treeNode) => {
  Modal.confirm({
    title: '删除文档',
    icon: createVNode(ExclamationCircleOutlined),
    content: `确定删除文档 ${treeNode.dataRef.title} 吗？`,
    onOk() {
      submitDeleteDoc(treeNode)
    },
    onCancel() {}
  })
}

const submitDeleteDoc = (treeNode) => {
  deleteLibDoc({
    library_key: state.library_key,
    doc_id: treeNode.dataRef.id
  }).then(() => {
    message.success('删除成功')
    let parentDoc = findDocInTree(docList.value, treeNode.dataRef.pid)

    if (parentDoc) {
      let index = parentDoc.children.findIndex((item) => item.id == treeNode.dataRef.id)
      parentDoc.children_num--
      parentDoc.isLeaf = parentDoc.children_num == 0
      parentDoc.children.splice(index, 1)
      if (activeMenuKey.value != 'doc') {
        return
      }

      if (parentDoc.children.length > 0) {
        onSelectDoc(parentDoc.children[0])
      } else {
        if (docList.value.length > 0) {
          onSelectDoc(docList.value[0])
        } else {
          currentDoc.value = null
          router.replace({
            path: '/public-library/home',
            query: {
              library_id: libraryId.value,
              library_key: state.library_key
            }
          })
        }
      }
    } else {
      let index = docList.value.findIndex((item) => item.id == treeNode.dataRef.id)

      docList.value.splice(index, 1)

      if (activeMenuKey.value != 'doc') {
        return
      }

      if (docList.value.length > 0) {
        onSelectDoc(docList.value[0])
      } else {
        currentDoc.value = null
        router.replace({
          path: '/public-library/home',
          query: {
            library_id: libraryId.value,
            library_key: state.library_key
          }
        })
      }
    }
  })
}

const onExportDoc = (treeNode) => {
  getLibDocInfo({
    library_key: state.library_key,
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
  currentDoc.value = data
  selectedKeys.value = [data.id]
  if (route.path == '/public-library/editor') {
    router.replace({
      path: '/public-library/editor',
      query: {
        library_id: libraryId.value,
        library_key: state.library_key,
        doc_id: data.id
      }
    })
  } else {
    router.push({
      path: '/public-library/editor',
      query: {
        library_id: libraryId.value,
        library_key: state.library_key,
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
        break
      }

      if (data[i].children) {
        loop(data[i].children, key, level + 1, callback)
      }
    }
  }

  // let errorMsg = ''

  // Find dragObject
  let dragObj
  // let dragNodeIndex = 0
  // let dragNodePrarent = null

  loop(data, dragKey, 0, (item, index, level, arr) => {
    dragObj = { ...item }
    // dragNodeIndex = index
    // dragNodePrarent = arr
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
    loop(data, dropKey, 0, (item) => {
      // if (level + 1 > 2) {
      //   errorMsg = '最多支持三级目录1'
      //   return
      // }
      item.children = item.children || []
      item.children_num++
      item.isLeaf = false

      // where to insert 示例添加到头部，可以是随意位置
      item.children.unshift(dragObj)
    })
  } else if (
    (info.node.children || []).length > 0 &&
    // Has children
    info.node.expanded &&
    // Is expanded
    dropPosition === 1 // On the bottom gap
  ) {
    loop(data, dropKey, 0, (item) => {
      // if (level + 1 > 2) {
      //   errorMsg = '最多支持三级目录2'

      //   return
      // }

      item.children = item.children || []
      item.children_num++
      item.isLeaf = false

      // where to insert 示例添加到头部，可以是随意位置
      item.children.unshift(dragObj)
    })
  } else {
    let ar = []
    let i = 0
    // let nextLevel = 0

    loop(data, dropKey, 0, (_item, index, level, arr) => {
      ar = arr
      i = index
      // nextLevel = level + 1
    })

    // if (nextLevel > 2) {
    //   errorMsg = '最多支持三级目录3'
    //   return
    // }

    if (dropPosition === -1) {
      ar.splice(i, 0, dragObj)
    } else {
      ar.splice(i + 1, 0, dragObj)
    }
  }

  // if (errorMsg) {
  // dragNodePrarent.splice(dragNodeIndex, 1, dragObj)
  // message.error(errorMsg)
  // return
  // }

  docList.value = data

  handleDocSort(dragNode, dropNode)
}

const handleDocSort = (dragNode, dropNode) => {
  let parentDoc = findParentDocInTree(docList.value, dragNode.dataRef.id)

  dragNode.dataRef.pid = parentDoc ? parentDoc.id : 0

  let data = {
    library_key: state.library_key,
    doc_id: dragNode.dataRef.id,
    pid: dragNode.dataRef.pid,
    level: parentDoc ? parentDoc.level + 1 : 0,
    sort: 0
  }

  let list = parentDoc ? parentDoc.children : docList.value
  let dragNodeIndex = -1
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

  doc.title = data.title
  doc.is_draft = data.is_draft
  currentDoc.value = doc
}

const getInfo = async () => {
  try {
    const res = await getLibraryInfo({ id: libraryId.value })
    if (!res.data.avatar) {
      res.data.avatar = LIBRARY_OPEN_AVATAR
    }
    Object.assign(state, res.data)
  } catch (error) {
    console.log(error)
  }
}

onMounted(() => {
  getInfo().then(() => {
    getMyDocList()
  })
})
</script>
