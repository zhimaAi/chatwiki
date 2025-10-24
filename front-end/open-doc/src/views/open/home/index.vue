<template>
  <div class="open-home-page">
    <div class="page-top">
      <div class="home-page-banner" :style="{
        backgroundImage: `url(${bannerImage})`,
      }">
        <div class="upload-banner-btn" v-if="isEdit">
          <div class="upload-btn" @click="uploadBannerInputRef.click()">
            <UploadOutlined class="btn-icon" />
            <span class="btn-text">修改主图</span>
            <input type="file" class="upload-input" ref="uploadBannerInputRef" accept="image/png,image/jpeg"
              @change="handleUploadChange">
          </div>
        </div>
        <action-popover :disabled="!isEdit" name="title" :menus="['edit']" @click="() => handleEditTitle()">
          <h1 class="library-title">
            {{ state.title }}
          </h1>
        </action-popover>

        <action-popover :disabled="!isEdit" name="content" :menus="['edit']" @click="() => handleEditContent()">
          <div class="library-desc" id="library-desc">
            {{ !state.content && isEdit ? '请输入描述' : state.content }}
          </div>
        </action-popover>

        <!-- 搜索框 -->
        <div class="search-box-wrapper">
          <div class="search-box" id="search-box">
            <input class="search-input" id="search-input" ref="searchInputRef" name="search" autocomplete="off"
              placeholder="" type="text" v-model="state.keyword" @keydown.enter="onSidebarSearch()" />
            <span class="search-clear" @click="handleClear" v-if="state.keyword"></span>
            <span class="search-btn" @click="onSidebarSearch()">
              <span class="search-icon"></span>
            </span>
          </div>
        </div>
      </div>

      <!-- 问题引导 -->
      <div class="keyword-list-wrapper">
        <ul class="keyword-list" id="keyword-list">
          <li class="keyword-item" v-for="item in state.question_guide" :key="item.id"
            @click="handleQuestionGuide(item)">
            <action-popover name="question" :disabled="!isEdit" :menus="['edit', 'delete']"
              @click="(name, value) => handleEditQuestion(value, item)">
              <div class="question-guide-item">
                <span class="keyword-text">
                  {{ item.question }}
                </span>
              </div>
            </action-popover>
          </li>
          <li class="keyword-item add-btn-item" @click="handleAddQuestion()"
            v-if="isEdit && state.question_guide.length < 10">
            <img src="@/assets/img/home_add.svg" class="add-btn" />
          </li>
        </ul>
      </div>
    </div>


    <!-- 快捷方式 -->
    <div class="home-page-body">
      <div>
        <VueDraggable class="shortcut-list" ref="el" :disabled="disabledShortcutDraggable" :draggable="'.sortable-drag'"
          :filter="'.add-item'" v-model="shortcutList" @end="onShortcutDragEnd">
          <template v-for="(item, index) in shortcutList" :key="item.key">
            <action-popover :disabled="!isEdit" name="title" :menus="['edit', 'delete']"
              @click="(name, value) => handleEditShortcut(value, index)">
              <div class="shortcut-item sortable-drag" :class="{ 'is-empty': item.children == 0 }"
                v-if="item.is_dir == 1">
                <div class="shortcut-item-header">
                  <div class="shortcut-title">
                    <span class="shortcut-icon">{{ item.doc_icon }}</span>
                    <span class="shortcut-name">{{ item.title }}</span>
                  </div>
                  <div class="shortcut-action">
                    <a class="shortcut-link" @click="handleDetail(item)"
                      :class="{ 'is-empty-link': item.children == 0 }">详情
                      <ArrowIcon direction="right" :color="item.children == 0 ? '#ccc' : 'currentColor'" />
                    </a>
                  </div>
                </div>
                <div class="shortcut-item-body">
                  <div class="post-list">
                    <template v-for="doc in item.children" :key="doc.id">
                      <div class="post-item" @click="handleDetail(doc)">
                        <span class="post-icon">{{ doc.doc_icon }}</span>
                        <span class="post-title">{{ doc.title }}</span>
                      </div>
                    </template>
                  </div>
                </div>
              </div>

              <div class="shortcut-item sortable-drag" v-else>
                <div class="shortcut-item-header">
                  <div class="shortcut-title">
                    <span class="shortcut-icon">{{ item.doc_icon }}</span>
                    <span class="shortcut-name">{{ item.title }}</span>
                  </div>
                  <div class="shortcut-action">
                    <a class="shortcut-link" @click="handleDetail(item)">详情
                      <ArrowIcon direction="right" />
                    </a>
                  </div>
                </div>
                <div class="shortcut-item-body">
                  <div class="post-content">
                    <div v-html="item.content"></div>
                  </div>
                </div>
              </div>
            </action-popover>
          </template>
          <div class="shortcut-item add-item" @click="handleAddShortcut" v-if="isEdit && shortcutList.length < 30">
            <span class="add-tip">
              <PlusCircleOutlined />&nbsp;添加文档快捷方式
            </span>
          </div>
        </VueDraggable>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, reactive, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useOpenDocStore } from '@/stores/open-doc'
import { message } from 'ant-design-vue'
import { PlusCircleOutlined, UploadOutlined } from '@ant-design/icons-vue'
import { VueDraggable } from 'vue-draggable-plus'
import ActionPopover from '@/views/open/home/components/action-popover.vue'
import ArrowIcon from '@/components/icons/arrow-icon.vue'
import { getIconTemplateById } from '@/config/open-doc/icon-template'


const router = useRouter()
const route = useRoute()
const openDocStore = useOpenDocStore()

const uploadBannerInputRef = ref(null)
const searchInputRef = ref(null)
const isEdit = ref(false)

const setEditStatus = (msg) => {
  isEdit.value = msg.type
}

const state = reactive({
  keyword: '',
  title: '',
  content: '',
  question_guide: [],
})
const bannerImage = ref('')
const shortcutList = ref([])
const shortcutKeyList = ref([])

const disabledShortcutDraggable = computed(() => {
  return !isEdit.value
})

const docId = computed(() => {
  return route.params.id
})

const previewKey = computed(() => {
  return openDocStore.previewKey
})

// libraryInfo 暂时未使用，保留以备后续需求
// const libraryInfo = computed(() => {
//   return openDocStore.libraryInfo
// })

const catalogMap = computed(() => {
  return openDocStore.catalogMap
})

const handleClear = () => {
  // 清空输入框
  state.keyword = ''

  // 获得焦点
  searchInputRef.value.focus()
}

const onSidebarSearch = () => {
  if (openDocStore.isEditPage) {
    message.warning('编辑模式不支持搜索功能')
    return
  }

  if (previewKey.value) {
    message.warning('预览模式不支持搜索功能')
    return
  }

  // 失去焦点
  searchInputRef.value.blur()

  router.push({
    name: 'open-search',
    params: {
      id: docId.value,
    },
    query: {
      v: state.keyword,
    },
  })

  state.keyword = ''
}

const handleDetail = (item) => {
  if (item.is_dir == 1 && item.children.length == 0) {
    return
  }

  if (openDocStore.isEditPage) {
    return message.warning('编辑模式不支持查看功能')
  }

  if (item.is_dir == 1) {
    if (!item.children.length) {
      return;
    }

    let docs = item.children.filter(item => item.is_dir == 0);
    if (!docs.length) {
      return;
    }

    router.push({
      path: `/doc/${docs[0].doc_key}`,
      params: {
        id: docId.value,
      }
    })
  } else {
    router.push({
      path: `/doc/${item.doc_key}`,
      params: {
        id: docId.value,
      }
    })
  }
}

const getDocIcon = (item, icon_template_config_id) => {
  let iconTemplate = getIconTemplateById(icon_template_config_id)
  if (item.doc_icon) {
    return item.doc_icon;
  }

  if (catalogMap.value[item.id]) {
    return catalogMap.value[item.id].doc_icon;
  } else {
    let iconConfig = iconTemplate.levels[item.level];

    return item.is_dir == 1 ? iconConfig.folder_icon : iconConfig.doc_icon;
  }
}

const getData = async () => {
  let data = await openDocStore.getHome(docId.value);
  let shortcuts = data.quick_doc_content_value || [];
  let shortcutKeys = data.quick_doc_content || [];

  shortcuts.forEach((item, index) => {
    item.key = shortcutKeys[index].key || item.id;
  })

  state.title = data.title;
  state.content = data.content;
  state.question_guide = data.question_guide;
  bannerImage.value = data.banner_img_url;

  // 提取公共逻辑到函数中 - 处理单个项目
  const processItem = (item, iconTemplateConfigId) => {
    item.level = catalogMap.value[item.id]?.level || 0;
    item.doc_icon = getDocIcon(item, iconTemplateConfigId);
  };

  // 处理快捷方式项目及其子项目
  shortcuts.forEach(item => {
    item.children = item.children.filter((item) => item.is_dir != 1);
    // 处理父项目
    processItem(item, data.icon_template_config_id);

    // 处理子项目
    if (item.children?.length) {
      item.children.forEach(child => {
        processItem(child, data.icon_template_config_id);
      });
    }
  });
  
  shortcutList.value = shortcuts;
  shortcutKeyList.value = shortcutKeys;
}

const handleUploadChange = (e) => {
  const file = e.target.files[0]
  if (!file) {
    return
  }

  // 验证文件大小（2MB = 2 * 1024 * 1024 bytes）
  const maxSize = 2 * 1024 * 1024
  if (file.size > maxSize) {
    message.error('文件大小不能超过2MB')
    // 清空文件输入框
    e.target.value = ''
    return
  }

  let msg = {
    key: 'banner',
    action: 'change',
    data: {
      doc_id: docId.value,
      file: file,
    },
  }

  postMessage(msg)
}

const handleQuestionGuide = (item) => {
  state.keyword = item.question

  onSidebarSearch()
}

const handleEditTitle = () => {
  let msg = {
    key: 'title',
    action: 'edit',
    data: {
      doc_id: docId.value,
      title: state.title,
      content: state.content,
    },
  }

  postMessage(msg)
}

const handleEditContent = () => {
  let msg = {
    key: 'content',
    action: 'edit',
    data: {
      doc_id: docId.value,
      title: state.title,
      content: state.content,
    },
  }

  postMessage(msg)
}

const handleEditQuestion = (actionKey, data) => {
  let msg = {
    key: 'question',
    action: actionKey,
    data: {
      id: data.id,
      question: data.question,
    },
  }

  postMessage(msg)
}

const handleAddQuestion = () => {
  let msg = {
    key: 'question',
    action: 'add',
    data: {
      id: '',
      question: '',
    },
  }

  postMessage(msg)
}

const handleAddShortcut = () => {
  let msg = {
    key: 'shortcut',
    action: 'add',
    data: {},
  }

  postMessage(msg)
}

const handleEditShortcut = (key, index) => {
  let data = shortcutKeyList.value[index];

  let msg = {
    key: 'shortcut',
    action: key,
    data: JSON.parse(JSON.stringify(data)),
  }

  postMessage(msg)
}

const onShortcutDragEnd = () => {
  let list = shortcutList.value.map(item => {
    return { doc_id: item.id * 1, key: item.key }
  })

  let msg = {
    key: 'shortcut',
    action: 'dragEnd',
    data: {
      list: list
    },
  }

  postMessage(msg)
}

const postMessage = (msg) => {
  if (!window.parent) {
    return
  }

  window.parent.postMessage(msg, '*')
}

function onMessage(event) {
  if (!event.data) {
    return
  }

  let { action, data } = event.data

  if (action === 'setEditStatus') {
    setEditStatus(data)
  }
}

onMounted(async () => {
  nextTick(async () => {
    window.addEventListener('message', onMessage)

    await getData()

    postMessage({
      action: 'init',
      key: 'init',
    })
  })
})
onUnmounted(() => {
  window.removeEventListener('message', onMessage)
})
</script>

<style lang="less" scoped>
.open-home-page {
  min-height: 100vh;
  background-color: #f5f6fd;

  .page-top {
    position: relative;
  }

  .home-page-banner {
    position: relative;
    height: 394px;
    background-size: cover;
    background-repeat: no-repeat;

    .upload-banner-btn {
      position: absolute;
      right: 16px;
      top: 16px;


      .upload-btn {
        position: relative;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 8px;
        padding: 5px 16px;
        border-radius: 6px;
        border: 1px solid #D9D9D9;
        background: #FFF;
        cursor: pointer;
        overflow: hidden;

        .upload-input {
          position: absolute;
          left: -1000px;
          top: -1000px;
          opacity: 0;
          cursor: pointer;
          z-index: 0;
        }
      }
    }

    .library-title {
      line-height: 44px;
      padding-top: 38px;
      margin-bottom: 26px;
      font-size: 36px;
      font-weight: 600;
      color: #1a1a1a;
      text-align: center;
    }

    .library-desc {
      line-height: 24px;
      margin-bottom: 34px;
      font-size: 16px;
      font-weight: 400;
      color: #595959;
      text-align: center;
    }

    .search-box-wrapper {
      padding: 0 16px;
    }

    .search-box {
      position: relative;
      max-width: 800px;
      margin: 0 auto;
      margin-bottom: 16px;
      border-radius: 12px;
      overflow: hidden;
    }

    .search-box .search-input {
      width: 100%;
      height: 58px;
      padding: 0 96px 0 24px;
      font-size: 16px;
      font-weight: 600;
      border-radius: 12px;
      color: #1a1a1a;
      border: 2px solid #262626;
      transition: border 0.2s;
    }

    .search-box .search-clear {
      position: absolute;
      display: block;
      right: 66px;
      top: 17px;
      width: 24px;
      height: 24px;
      background: url('@/assets/img/search_clear.svg') no-repeat;
      cursor: pointer;
    }

    .search-box .search-btn {
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      height: 58px;
      width: 58px;
      transition: background 0.2s;
    }

    .search-box .search-icon {
      width: 24px;
      height: 24px;
      background: url('@/assets/img/search_black.svg') no-repeat;
    }

    .search-box .search-btn:hover {
      cursor: pointer;
    }

    .search-box .search-input {
      border: 2px solid #2475fc;
    }

    .search-box .search-input::placeholder {
      color: #ccc;
    }

    .search-box .search-btn {
      border-radius: 0 12px 12px 0;
      background: #2475fc;
    }

    .search-box .search-icon {
      background: url('@/assets/img/search_white.svg') no-repeat;
    }

    .search-box.has-value .search-clear {
      display: block;
    }
  }

  .keyword-list-wrapper {
    margin-top: 24px;
  }

  .keyword-list {

    font-size: 0;
    max-width: 1200px;
    margin: 0 auto;

    .keyword-item {
      position: relative;
      margin-bottom: 12px;
      margin-right: 12px;
      font-size: 0;
      vertical-align: middle;
    }

    .keyword-item:last-child {
      margin-right: 0;
    }

    .keyword-item .keyword-text {
      display: inline-block;
      line-height: 22px;
      padding: 12px 16px;
      font-size: 14px;
      font-weight: 400;
      color: #1a1a1a;
      background: #edeff2;
      border-radius: 4px 12px 12px 12px;
      cursor: pointer;
      transition: all 0.2s;
    }

    .keyword-item .keyword-text:hover {
      background: #e6efff;
    }

    .keyword-item .add-btn {
      width: 32px;
      height: 32px;
      cursor: pointer;
      transition: all 0.2s;
    }

    .keyword-item .add-btn:hover {
      opacity: 0.8;
    }
  }

  .home-page-body {
    max-width: 1200px;
    margin: 0 auto;
    padding: 24px 0;
  }

  .shortcut-list {
    display: flex;
    flex-flow: row wrap;
    gap: 16px;
    padding: 0 16px;

    .shortcut-item {
      width: 100%;
      padding: 16px;
      border-radius: 12px;
      overflow: hidden;
      border: 1px solid #f0f0f0;
      background-color: #fff;

      &.add-item {
        height: 200px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #595959;
        cursor: pointer;
      }

      .shortcut-item-header {
        display: flex;
        justify-content: space-between;
        height: 24px;
        line-height: 24px;
        margin-bottom: 8px;

        .shortcut-title {
          flex: 1;
          display: flex;
          align-items: center;
          overflow: hidden;

          .shortcut-icon {
            margin-right: 8px;
            font-size: 18px;
          }

          .shortcut-name {
            flex: 1;
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            font-size: 16px;
            font-weight: 600;
            color: #000000;
          }
        }

        .shortcut-action {
          display: flex;
          align-items: center;

          .shortcut-link {
            line-height: 22px;
            padding: 1px 8px;
            font-size: 14px;
            font-weight: 400;
            border-radius: 6px;
            color: #595959;
            cursor: pointer;
            transition: all 0.2s;

            &:hover {
              background: #E4E6EB;
            }

            &.is-empty-link {
              color: #ccc;
              cursor: not-allowed;

              &:hover {
                background: transparent;
              }
            }
          }
        }
      }

      .shortcut-item-body {
        .post-content {
          height: 134px;
          line-height: 22px;
          padding: 1px 0;
          font-size: 14px;
          color: #595959;
          display: -webkit-box;
          -webkit-line-clamp: 6;
          line-clamp: 6;
          -webkit-box-orient: vertical;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .post-list {
          height: 134px;
          overflow-y: auto;

          .post-item {
            display: flex;
            align-items: center;
            height: 32px;
            line-height: 32px;
            padding: 0 8px;
            margin-bottom: 2px;
            cursor: pointer;
            border-radius: 6px;
            transition: all 0.2s;
            overflow: hidden;

            &:last-child {
              margin-bottom: 0;
            }

            &:hover {
              background: #f2f4f7;
            }

            .post-icon {
              margin-right: 4px;
              font-size: 16px;
            }

            .post-title {
              flex: 1;
              font-size: 14px;
              color: #595959;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
          }
        }
      }
    }
  }
}

@media (min-width: 992px) {
  .open-home-page {
    .home-page-banner {
      margin: 0 auto;
    }

    .search-box {
      margin: 0 auto;
    }

    .keyword-list-wrapper {
      position: absolute;
      left: 0;
      bottom: 0;
      right: 0;
      height: 148px;
      margin-top: 0;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .keyword-list {
      text-align: center;

      .keyword-item {
        display: inline-block;

        .keyword-text {
          line-height: 22px;
          font-size: 14px;
          font-weight: 400;
          padding: 5px 8px;
          border-radius: 12px;
          color: #595959;
        }
      }

      .keyword-item.add-btn-item {
        height: 32px;
        line-height: 32px;
      }

      .keyword-item .add-btn {
        width: 24px;
        height: 24px;
      }
    }

    .shortcut-list {
      padding: 0;

      .shortcut-item {
        width: calc((100% - 32px) / 3);
      }
    }
  }
}
</style>
