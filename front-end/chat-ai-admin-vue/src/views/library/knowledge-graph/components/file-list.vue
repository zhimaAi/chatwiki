<style lang="less" scoped>
.file-list-wrapper {
  position: relative;
  display: flex;
  flex-flow: column nowrap;
  height: 100%;
  width: 256px;
  height: 100%;
  padding: 16px 0;
  border-radius: 6px;
  background: #fff;
}

.open-file-box {
  padding: 0 8px;
  user-select: none;
}

.file-list-box {
  flex: 1;
  overflow: hidden;
  user-select: none;

  .file-items {
    padding: 0 8px;
  }

  .loading-box {
    padding-left: 8px !important;
  }
}

.file-item {
  margin-bottom: 4px;

  &.active .file-info,
  .file-info:hover {
    color: #2475FC;
    background: #E6EFFF;
  }

  .file-info {
    display: flex;
    align-items: center;
    padding: 5px 8px;
    border-radius: 6px;
    color: #262626;
    cursor: pointer;
    transition: all 0.2s;

    .caret-icon {
      font-size: 16px;
      color: #595959;
    }

    .file-icon {
      margin-left: 4px;
      font-size: 16px;
    }

    .file-name {
      flex: 1;
      font-size: 14px;
      font-weight: 400;
      margin-left: 4px;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }

  .document-fragments {
    overflow: hidden;

    .loading-box {
      padding-left: 44px !important;
    }
  }

  .fragment-list {
    .fragment-item {
      padding: 5px 8px 5px 44px;
      border-radius: 6px;
      cursor: pointer;

      &:hover,
      &.active {
        background-color: #F2F4F7;
      }

      .doc-content {
        line-height: 22px;
        max-height: 44px;
        font-size: 14px;
        font-weight: 400;
        color: #595959;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        display: -webkit-box;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: normal;
      }
    }
  }
}

.fragment-details {
  display: flex;
  flex-direction: column;
  position: absolute;
  top: 0;
  left: 264px;
  width: 320px;
  height: 60%;
  border-radius: 12px;
  background: #fff;

  .fragment-title {
    line-height: 22px;
    padding: 16px 16px 0 16px;
    font-size: 14px;
    font-weight: 600;
    color: #262626;
  }

  .fragment-details-body {
    flex: 1;
    overflow: hidden;
  }

  .fragment-content {
    line-height: 22px;
    padding: 8px 16px 16px 16px;
    font-size: 14px;
    font-weight: 400;
    color: #595959;
    overflow: hidden;
  }
}

.loading-box {
  display: flex;
  align-items: center;
  margin: 10px 0;

  .loading-text {
    margin-left: 4px;
    font-size: 14px;
    color: #595959;
  }
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>

<template>
  <transition name="slide-fade">
    <div class="file-list-wrapper" v-if="props.show">
      <div class="open-file-box">
        <div class="file-item active" v-if="currentOpenFile.show">
          <div class="file-info">
            <CaretDownFilled class="caret-icon" @click="hideFragmentList()" v-if="fragmentList.show" />
            <CaretRightOutlined class="caret-icon" @click="showFragmentList()" v-else />
            <svg-icon class="file-icon" name="doc-file" style=""></svg-icon>
            <span class="file-name" @click="openFile(currentOpenFile)">{{ currentOpenFile.file_name }}</span>
          </div>
          <div class="document-fragments" :style="{maxHeight: maxHeight + 'px', height: fragmentsHeight > 0 ? fragmentsHeight + 'px' : 'auto' }"
            v-if="fragmentList.show">
            <cu-scroll ref="fragmentsScrollView" group-id="a" @onScrollEnd="onLoadFragmentListt">
              <div class="fragment-list" ref="fragmentListRef">
                <div class="fragment-item" :class="{ active: activeFragmentId == item.id }"
                  @mouseenter="showFragmentDetails(item)" @mouseleave="handleFragmentLeave" @click="openFragment(item)"
                  v-for="item in fragmentList.list" :key="item.id">
                  <div class="doc-content">{{ item.content }}</div>
                </div>
              </div>
              <template v-if="fragmentList.list.length < fragmentList.total">
                <div class="loading-box"><a-spin size="small" /> <span class="loading-text">加载中...</span></div>
              </template>
            </cu-scroll>
          </div>
        </div>
      </div>
      <div class="file-list-box ">
        <cu-scroll group-id="a" @onScrollEnd="onLoadFileList">
          <div class="file-items">
            <template v-for="item in fileList.list" :key="item.id">
              <div class="file-item" v-if="item.id != currentOpenFile.id">
                <div class="file-info" @click="openFile(item)">
                  <CaretRightOutlined class="caret-icon" />
                  <svg-icon class="file-icon" name="doc-file" style=""></svg-icon>
                  <span class="file-name">{{ item.file_name }}</span>
                </div>
              </div>
            </template>

            <template v-if="fileList.showMore">
              <div class="loading-box"><a-spin size="small" /> <span class="loading-text">加载中...</span></div>
            </template>
          </div>
        </cu-scroll>
      </div>

      <div class="fragment-details" v-if="fragmentDetails.show" @mouseenter="handleDetailsEnter"
        @mouseleave="hideFragmentDetails">
        <div class="fragment-title">{{ fragmentDetails.title || '无标题片段' }}</div>
        <div class="fragment-details-body">
          <cu-scroll>
            <div class="fragment-content">
              {{ fragmentDetails.content }}
            </div>
          </cu-scroll>
        </div>
      </div>
    </div>
  </transition>


</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { useRoute } from 'vue-router'
import { getLibraryFileList, getParagraphList, getLibFileInfo } from '@/api/library/index'
import { reactive, ref, nextTick, computed, onMounted } from 'vue'
import { CaretDownFilled, CaretRightOutlined } from '@ant-design/icons-vue'
import CuScroll from './cu-scroll.vue';

const { getStorage, removeStorage } = useStorage('localStorage')

const emit = defineEmits(['openFile', 'openFragment'])

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const route = useRoute()
const autoOpenFileId = getStorage('graph:autoOpenFileId') || null;

const library_id = computed(() => {
  return route.query.id
})

const fileList = reactive({
  list: [],
  library_id: library_id.value,
  status: 2,
  page: 1,
  size: 20,
  total: 0,
  showMore: true,
})

const getFileList = async () => {
  try {
    const res = await getLibraryFileList({
      library_id: fileList.library_id,
      status: fileList.status,
      page: fileList.page,
      size: fileList.size,
    })

    if (res.res != 0) {
      return
    }

    let list = res.data.list || []

    fileList.total = res.data.total

    if (autoOpenFileId) {
      // 过滤掉已打开的文件
      list = list.filter(item => item.id != autoOpenFileId)
    }

    fileList.list = [...fileList.list, ...list]

    if (fileList.list.length >= fileList.total) {
      fileList.showMore = false
    }
    // 自动打开第一个文件
    if (fileList.page == 1 && fileList.list.length > 0) {
      openFile(fileList.list[0])
    }

  }catch(err){
    console.log(err);
    fileList.showMore = false
  }
  
}

const onLoadFileList = () => {
  if (fileList.list.length >= fileList.total) {
    return
  }

  fileList.page++
  getFileList()
}

const fragmentsScrollView = ref(null)
const fragmentListRef = ref(null)
const fragmentsHeight = ref(0)

const currentOpenFile = reactive({
  show: false,
  file_name: '',
  id: '',
})

const fragmentList = reactive({
  show: true,
  page: 1,
  size: 20,
  total: 0,
  list: [],
  showMore: true,
})

const activeFragmentId = ref(null)

const fragmentDetails = reactive({
  show: false,
  title: '',
  content: ''
})

const openFile = async (item) => {
  if(item.id == currentOpenFile.id){
    activeFragmentId.value = null
    emit('openFile', item)

    return
  }
  
  currentOpenFile.show = true
  currentOpenFile.file_name = item.file_name
  currentOpenFile.id = item.id
  activeFragmentId.value = null
  fragmentList.page = 1
  fragmentList.total = 0
  fragmentList.list = []
  fragmentList.show = true
  fragmentList.showMore = true

  await getFragments()

  emit('openFile', item)
}

const hideFragmentList = () => {
  fragmentList.show = false
}

const showFragmentList = () => {
  fragmentList.show = true
}

const maxHeight = ref(300)

const getFragments = async () => {
  const res = await getParagraphList({
    file_id: currentOpenFile.id,
    page: fragmentList.page,
    size: fragmentList.size,
    status: -1,
    graph_status: -1,
    category_id: -1
  })

  if (res.res != 0) {
    return
  }

  fragmentList.total = res.data.total

  if (fragmentList.page == 1) {
    fragmentList.list = res.data.list || []
  } else {
    fragmentList.list = [...fragmentList.list, ...res.data.list]
  }

  if (fragmentList.list.length >= fragmentList.total) {
    fragmentList.showMore = false
  }

  nextTick(() => {
    let height = fragmentListRef.value?.offsetHeight || 0
    if (height > maxHeight.value) {
      fragmentsHeight.value = maxHeight.value
    } else {
      fragmentsHeight.value = height
    }

    nextTick(() => {
      fragmentsScrollView.value?.refresh()
    })
  })
}

const onLoadFragmentListt = () => {
  if (fragmentList.list.length >= fragmentList.total) {
    return
  }

  fragmentList.page++
  getFragments()
}

const openFragment = (item) => {
  activeFragmentId.value = item.id
  emit('openFragment', item)
}

const showFragmentDetails = (item) => {
  clearTimeout(leaveTimer)
  fragmentDetails.title = item.title
  fragmentDetails.content = item.content
  fragmentDetails.show = true
}

let leaveTimer = null

const handleFragmentLeave = () => {
  leaveTimer = setTimeout(() => {
    fragmentDetails.show = false
  }, 200)
}

const handleDetailsEnter = () => {
  clearTimeout(leaveTimer)
}

const hideFragmentDetails = () => {
  fragmentDetails.show = false
}

const getFileInfo = async () => {
  try {
    const res = await getLibFileInfo({
      id: autoOpenFileId
    });

    fileList.list.push(res.data)
  } catch (err) {
    console.log(err);
  }
}

onMounted(async () => {
  if(autoOpenFileId){
    await getFileInfo()
    removeStorage('graph:autoOpenFileId')
  }

  getFileList()
})

defineExpose({
  openFile
})
</script>
