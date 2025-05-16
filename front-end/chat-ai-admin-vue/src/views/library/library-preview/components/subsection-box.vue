<template>
  <div class="subsection-box content-container" ref="containerRef">
    <div class="subsection-box-title">
      分段预览
      <span>共{{ props.total }}个分段</span>
    </div>
    <div
      class="list-item"
      v-for="(item, index) in props.paragraphLists"
      :key="item.id"
      @click.stop="handleToTargetPage(item, index, $event)"
      :style="{'--status-color': getColor(item)}"
    >
      <div class="top-block">
        <div class="title">
          <!-- id：{{ item.id }} -->
          分段{{ index + 1 }}
          <div class="title-block">
            {{ item.title }}
          </div>
          <span>共{{ item.word_total }}个字符</span>
          <span>
            嵌入状态：{{ item.status_text }}<LoadingOutlined v-if="item.status == 3" />
            <a-tooltip v-if="item.status == 2 && item.errmsg" :title="item.errmsg">
              <strong class="cfb363f"
                >原因<ExclamationCircleOutlined class="err-icon cfb363f"
              /></strong>
            </a-tooltip>
          </span>
          <span v-if="detailsInfo.graph_switch">
            知识图谱状态：{{ item.graph_status_text }}
            <a-tooltip
              v-if="item.graph_status == 3 && item.graph_err_msg"
              :title="item.graph_err_msg"
            >
              <strong class="cfb363f"
                >原因<ExclamationCircleOutlined class="err-icon cfb363f"
              /></strong>
            </a-tooltip>
          </span>
        </div>
        <div class="right-opration">
          <a-dropdown>
            <div class="hover-btn-box">
              <StarFilled
                @click.stop="handleSetCategory(item, {})"
                :style="{ color: getColor(item), 'font-size': '16px' }"
                v-if="item.category_id > 0"
              />
              <StarOutlined
                @click="handleSetCategory(item, startLists[0])"
                style="font-size: 16px;color: #595959;"
                v-else
              />
            </div>

            <template #overlay>
              <a-menu>
                <a-menu-item v-for="star in startLists" :key="star.id">
                  <div class="start-item" @click.stop="handleSetCategory(item, star)">
                    <StarFilled :style="{ color: colorLists[star.type] }" />
                    <div>{{ star.name || '-' }}</div>
                  </div>
                </a-menu-item>
                <a-menu-item>
                  <div class="start-item" @click.stop="handleOpenSetStartModal">
                    <SettingOutlined />
                    <div>精选设置</div>
                  </div>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>

          <template v-if="detailsInfo.graph_switch">
            <a-tooltip>
              <template #title>知识图谱</template>
              <div class="hover-btn-box" @click="openGraphModel(item)">
                <svg-icon name="graph-icon" style="font-size: 16px;color: #595959;"></svg-icon>
              </div>
            </a-tooltip>
          </template>

          <a-tooltip>
            <template #title>重新转换</template>
            <div class="hover-btn-box">
              <SyncOutlined @click.stop="toReSegmentationPage(item, index)" />
            </div>
          </a-tooltip>
          <a-dropdown placement="bottomRight">
            <div class="hover-btn-box">
              <MoreOutlined style="font-size: 16px;color: #595959;" />
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item>
                  <div @click.stop="handleOpenEditModal(item)">编辑</div>
                </a-menu-item>
                <a-menu-item>
                  <div @click.stop="hanldleDelete(item.id)">删除</div>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
      <div class="content-box" v-if="item.question">Q：{{ item.question }}</div>
      <div class="content-box" v-if="item.answer">A：{{ item.answer }}</div>
      <div class="content-box" @mouseup.stop="handleMouseUp($event, index)" v-html="item.content"></div>
      <div class="fragment-img" v-if="item.images.length > 0" v-viewer>
        <img v-for="(item, imageIndex) in item.images" :key="imageIndex" :src="item" alt="" />
      </div>
    </div>
    <ClassificationMarkModal @ok="getCategoryLists" ref="classificationMarkModalRef" />
    <GraphModel ref="graphModelRef" />
    <!-- 选择内容气泡 -->
    <div 
      v-if="showBubble"
      class="bubble-card"
      :style="bubbleStyle"
    >
      <button class="bubble-item" @click.stop="handleAction('separate')">
        <svg-icon
          name="segmentation"
          style="font-size: 16px; color: #262626; margin-right: 4px"
        ></svg-icon>
        单独成段
      </button>
      <button class="bubble-item" @click.stop="handleAction('merge-prev')">
        <svg-icon
          name="segmentation-up"
          style="font-size: 16px; color: #262626; margin-right: 4px"
        ></svg-icon>
        合并到上一分段
      </button>
      <button class="bubble-item" @click.stop="handleAction('merge-next')">
        <svg-icon
          name="segmentation-down"
          style="font-size: 16px; color: #262626; margin-right: 4px"
        ></svg-icon>
        合并到下一分段
      </button>
      <button class="bubble-item" @click.stop="handleAction('delete')">
        <svg-icon
          name="delete"
          style="font-size: 16px; color: #262626; margin-right: 4px"
        ></svg-icon>
        删除
      </button>
    </div>
  </div>
</template>
<script setup>
import { ref, createVNode, computed, onMounted, onUnmounted   } from 'vue'
import { useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import {
  ExclamationCircleOutlined,
  MoreOutlined,
  SyncOutlined,
  LoadingOutlined,
  StarOutlined,
  StarFilled,
  SettingOutlined
} from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import ClassificationMarkModal from './classification-mark-modal.vue'
import GraphModel from './graph-model/index.vue'
import colorLists from '@/utils/starColors.js'
import {
  deleteParagraph,
  editParagraph,
  getCategoryList,
  updateParagraphCategory
} from '@/api/library'

const emit = defineEmits([
  'handleDelParagraph',
  'handleScrollTargetPage',
  'openEditSubscription',
  'handleConvert',
  'getStatrList',
  'handleSplit',
  'handleSplitNext',
  'handleSplitUp',
  'handleSplitDelete'
])
const props = defineProps({
  paragraphLists: {
    type: Array,
    default: () => []
  },
  detailsInfo:{
    type: Object,
    default: () => {}
  },
  total: {
    type: [Number, String],
    default: 0
  }
})

const route = useRoute()

const toReSegmentationPage = (item, index) => {
  let { id, title, content, question, answer, images, category_id } = item
  let similar_questions = item.similar_questions || []
  Modal.confirm({
    title: '重新转换确认',
    icon: null,
    content: `确定要重新转换【分段${index + 1}】吗?`,
    onOk() {
      return new Promise((resolve) => {
        editParagraph({ id, title, content, question, answer, images, category_id, similar_questions: JSON.stringify(similar_questions)  })
          .then(() => {
            emit('handleConvert', item)
            resolve()
          })
          .catch(() => {
            resolve()
          })
      })
    },
    onCancel() {}
  })
}

const handleOpenEditModal = (item) => {
  emit('openEditSubscription', item)
}
const hanldleDelete = (id) => {
  Modal.confirm({
    title: '提示',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认是否删除该分段?',
    onOk() {
      return new Promise((resolve) => {
        deleteParagraph({ id })
          .then(() => {
            message.success('删除成功')
            emit('handleDelParagraph', id)
            resolve()
          })
          .catch(() => {
            resolve()
          })
      })
    },
    onCancel() {}
  })
}

const handleToTargetPage = (item, index, event) => {
  // 如果有选中文本则不执行跳转
  if (window.getSelection().toString().trim()) return;

  emit('handleScrollTargetPage', {
    page_num: item.page_num,
    index
  })
  
  handleClickOutside(event)
}

const startLists = ref([])
const getCategoryLists = () => {
  let params = {
    file_id: route.query.id
  }
  getCategoryList(params).then((res) => {
    startLists.value = res.data || []
    emit('getStatrList', res.data || [])
  })
}
getCategoryLists()

const classificationMarkModalRef = ref(null)
const handleOpenSetStartModal = () => {
  classificationMarkModalRef.value.show()
}

const getColor = (data) => {
  let type = startLists.value.filter((item) => item.id == data.category_id)[0]?.type
  if (type) {
    return colorLists[type]
  } else {
    return ''
  }
}
const handleSetCategory = (item, star = {}) => {
  updateParagraphCategory({
    id: item.id,
    category_id: star.id || 0
  }).then(() => {
    message.success('设置成功')
    item.category_id = star.id
    getCategoryLists()
  })
}

const graphModelRef = ref(null)
const openGraphModel = (item) => {

  graphModelRef.value.open(item)
}

const containerRef = ref(null);
const showBubble = ref(false);
const position = ref({ x: 0, y: 0 });
const selectedText = ref('');
const selectedRange = ref(null);
const selectedData = ref({
  text: '',
  parentIndex: -1,
  originalHtml: ''
});

// 计算气泡卡片样式
const bubbleStyle = computed(() => {
  if (!showBubble.value) return {};
  
  // 获取气泡卡片的预估宽度和高度
  const bubbleWidth = 162; // 根据实际样式调整
  
  // 计算理想位置
  let left = position.value.x;
  let top = position.value.y;
  
  // 获取容器和视口尺寸
  const containerRect = containerRef.value.getBoundingClientRect();
  const viewportWidth = window.innerWidth;
  
  // 防止左侧超出
  if (left < 10) {
    left = 10;
  }
  // 防止右侧超出
  else if (left + bubbleWidth > viewportWidth - 50) {
    left = viewportWidth - bubbleWidth - 50;
  }
  
  // 转换为相对于容器的位置
  const containerLeft = left - containerRect.left;
  const containerTop = top;
  
  return {
    left: `${containerLeft}px`,
    top: `${containerTop}px`,
    // 添加一个箭头指向选中文本
    // '--arrow-offset': `${position.value.x - containerLeft}px`
  };
});

const handleMouseUp = (event, parentIndex) => {
  setTimeout(() => {
    const selection = window.getSelection();
    const selectedTextContent = selection.toString().trim();

    if (selectedTextContent) {
      event.stopPropagation(); // 阻止事件冒泡
      // 保存选中信息
      selectedData.value = {
        text: selectedTextContent,
        parentIndex,
        originalHtml: props.paragraphLists[parentIndex].content
      };

      selectedText.value = selectedTextContent;
      selectedRange.value = selection.getRangeAt(0).cloneRange();
      
      // 获取选中文本的位置
      const range = selection.getRangeAt(0);
      const rect = range.getBoundingClientRect();
      
      // 计算相对于内容区域的位置
      const contentRect = containerRef.value.getBoundingClientRect();
      
      position.value = {
        x: rect.left + rect.width / 2,
        y: rect.top - contentRect.top - 50 // 向上偏移50px显示气泡
      };

      showBubble.value = true;
    } else {
      showBubble.value = false;
    }
  }, 50);
};

const handleAction = (action) => {
  const { parentIndex } = selectedData.value;
  
  switch(action) {
    case 'separate':
      separateParagraph(parentIndex);
      break;
    case 'merge-next':
      mergeWithNext(parentIndex);
      break;
    case 'merge-prev':
      mergeWithPrevious(parentIndex);
      break;
    case 'delete':
      onDeleteParagraph(parentIndex);
      break;
  }

  showBubble.value = false;
  window.getSelection().removeAllRanges();
};

// 新增HTML内容拆分方法
const splitHtmlContent = (container, startOffset, endOffset) => {
  // 增加全选内容判断
  const isFullSelection = 
    startOffset === 0 && 
    endOffset === container.textContent.length &&
    container.textContent.trim().length > 0;

  if (isFullSelection) {
    return ['', container.innerHTML, '']; // 当全选时返回完整内容
  }
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT)
  let node
  let count = 0
  let startNode, endNode
  let startIndex, endIndex

  while ((node = walker.nextNode())) {
    const length = node.nodeValue.length
    if (!startNode && count + length >= startOffset) {
      startNode = node
      startIndex = startOffset - count
    }
    if (!endNode && count + length >= endOffset) {
      endNode = node
      endIndex = endOffset - count
      break
    }
    count += length
  }

  const range = document.createRange()
  range.setStart(startNode, startIndex)
  range.setEnd(endNode, endIndex)

  const beforeRange = document.createRange()
  beforeRange.setStart(container, 0)
  beforeRange.setEnd(startNode, startIndex)

  const afterRange = document.createRange()
  afterRange.setStart(endNode, endIndex)
  afterRange.setEnd(container, container.childNodes.length)

  return [
    getRangeHtml(beforeRange),
    getRangeHtml(range),
    getRangeHtml(afterRange)
  ]
}

const getRangeHtml = (range) => {
  const div = document.createElement('div')
  div.appendChild(range.cloneContents())
  return div.innerHTML
}

// 修改后的分段方法
const separateParagraph = (index) => {
  const content = props.paragraphLists[index].content
  const selection = window.getSelection()
  const range = selection.getRangeAt(0)
  const selectedText = range.toString()

  if (!selectedText) return

  Modal.confirm({
    title: '单独分段',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认将该分段内容单独分段么？分段后，原分段的内容会被删除',
    onOk() {
      // 创建临时容器处理HTML内容
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = content

      // 拆分选中内容
      const [beforeContent, selectedContent, afterContent] = splitHtmlContent(
        tempDiv,
        range.startOffset,
        range.endOffset
      )

      // 新增全选判断逻辑
      const isEmptyBefore = beforeContent.replace(/<[^>]+>/g, '').trim() === '';
      const isEmptyAfter = afterContent.replace(/<[^>]+>/g, '').trim() === '';

      // 触发更新时需要保留原始选中信息
      if (isEmptyBefore && isEmptyAfter) {
        // 当全选内容时直接替换原段落
        emit('handleSplit', {
          index,
          beforeContent: '',       // 清空原内容
          selectedContent,
          afterContent: '',        // 不保留后续内容
          isFullSelection: true    // 添加全选标识
        });
      } else {
        emit('handleSplit', {
          index,
          beforeContent,
          selectedContent,
          afterContent,
          originalSelection: {
            start: range.startOffset,
            end: range.endOffset,
            parentIndex: index
          }
        })
      }
    }
  })
}

// 合并到下一分段
const mergeWithNext = (index) => {
  if (index >= props.paragraphLists.length - 1) {
    return message.error('没有下一个分段')
  };

  const content = props.paragraphLists[index].content
  const selection = window.getSelection()
  const range = selection.getRangeAt(0)
  const selectedText = range.toString()

  if (!selectedText) return

  Modal.confirm({
    title: '合并到下一个分段',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认将该分段内容合并到下一个分段么？分段后，原分段的内容会被删除',
    onOk() {
      // 创建临时容器处理HTML内容
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = content

      // 拆分选中内容
      const [beforeContent, selectedContent, afterContent] = splitHtmlContent(
        tempDiv,
        range.startOffset,
        range.endOffset
      )

      // 新增全选判断逻辑
      const isEmptyBefore = beforeContent.replace(/<[^>]+>/g, '').trim() === '';
      const isEmptyAfter = afterContent.replace(/<[^>]+>/g, '').trim() === '';

      // 触发更新时需要保留原始选中信息
      if (isEmptyBefore && isEmptyAfter) {
        // 当全选内容时直接替换原段落
        emit('handleSplitNext', {
          index,
          beforeContent: '',       // 清空原内容
          selectedContent,
          afterContent: '',        // 不保留后续内容
          isFullSelection: true    // 添加全选标识
        });
      } else {
        // 触发更新时需要保留原始选中信息
        emit('handleSplitNext', {
          index,
          beforeContent,
          selectedContent,
          afterContent,
          originalSelection: {
            start: range.startOffset,
            end: range.endOffset,
            parentIndex: index
          }
        })
      }
    }
  })
};

// 合并到上一分段
const mergeWithPrevious = (index) => {
  if (index <= 0) {
    return message.error('没有上一个分段')
  };

  const content = props.paragraphLists[index].content
  const selection = window.getSelection()
  const range = selection.getRangeAt(0)
  const selectedText = range.toString()

  if (!selectedText) return

  Modal.confirm({
    title: '合并到上一个分段',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认将该分段内容合并到上一个分段么？分段后，原分段的内容会被删除',
    onOk() {
      // 创建临时容器处理HTML内容
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = content

      // 拆分选中内容
      const [beforeContent, selectedContent, afterContent] = splitHtmlContent(
        tempDiv,
        range.startOffset,
        range.endOffset
      )

      // 新增全选判断逻辑
      const isEmptyBefore = beforeContent.replace(/<[^>]+>/g, '').trim() === '';
      const isEmptyAfter = afterContent.replace(/<[^>]+>/g, '').trim() === '';


      // 触发更新时需要保留原始选中信息
      if (isEmptyBefore && isEmptyAfter) {
        // 当全选内容时直接替换原段落
        emit('handleSplitUp', {
          index,
          beforeContent: '',       // 清空原内容
          selectedContent,
          afterContent: '',        // 不保留后续内容
          isFullSelection: true    // 添加全选标识
        });
      } else {
        // 触发更新时需要保留原始选中信息
        emit('handleSplitUp', {
          index,
          beforeContent,
          selectedContent,
          afterContent,
          originalSelection: {
            start: range.startOffset,
            end: range.endOffset,
            parentIndex: index
          }
        })
      }
    }
  })
};

// 删除当前分段
const onDeleteParagraph = (index) => {
  const content = props.paragraphLists[index].content
  const selection = window.getSelection()
  const range = selection.getRangeAt(0)
  const selectedText = range.toString()

  if (!selectedText) return

  Modal.confirm({
    title: '删除',
    icon: createVNode(ExclamationCircleOutlined),
    content: '确认将该分段内容删除么？',
    onOk() {
      // 创建临时容器处理HTML内容
      const tempDiv = document.createElement('div')
      tempDiv.innerHTML = content

      // 拆分选中内容
      const [beforeContent, selectedContent, afterContent] = splitHtmlContent(
        tempDiv,
        range.startOffset,
        range.endOffset
      )

      // 新增全选判断逻辑
      const isEmptyBefore = beforeContent.replace(/<[^>]+>/g, '').trim() === '';
      const isEmptyAfter = afterContent.replace(/<[^>]+>/g, '').trim() === '';

      // 触发更新时需要保留原始选中信息
      if (isEmptyBefore && isEmptyAfter) {
        // 当全选内容时直接替换原段落
        emit('handleSplitDelete', {
          index,
          beforeContent: '',       // 清空原内容
          selectedContent: '',     // 删除-不保留选中内容
          afterContent: '',        // 不保留后续内容
          isFullSelection: true    // 添加全选标识
        });
      } else {
        // 触发更新时需要保留原始选中信息
        emit('handleSplitDelete', {
          index,
          beforeContent,
          selectedContent: '',     // 删除-不保留选中内容
          afterContent,
          originalSelection: {
            start: range.startOffset,
            end: range.endOffset,
            parentIndex: index
          }
        })
      }
    }
  })
};

// 点击其他地方隐藏气泡
const handleClickOutside = (event) => {
  if (showBubble.value && !event.target.closest('.bubble-card')) {
    showBubble.value = false;
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});

defineExpose({ handleOpenEditModal })
</script>
<style lang="less" scoped>
.subsection-box {
  background: #f2f4f7;
  padding: 14px 16px;
  width: 100%;
  .hover-btn-box {
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
    &:hover {
      background: #e4e6eb;
      border-radius: 6px;
    }
  }
  .subsection-box-title {
    display: flex;
    align-items: center;
    font-size: 14px;
    line-height: 22px;
    font-weight: 600;
    color: #242933;
    span {
      color: #7a8699;
      font-weight: 400;
      margin-left: 8px;
    }
  }

  .list-item::before {
    content: "";
    position: absolute;
    left: 0;
    top: 0;
    width: 4px;
    height: 100%;
    background-color: var(--status-color, #FFF);
  }
  .list-item {
    position: relative;
    overflow: hidden;
    margin-top: 8px;
    width: 100%;
    background: #fff;
    border-radius: 6px;
    .top-block {
      padding: 16px 16px 0;
      display: flex;
      align-items: center;
      justify-content: space-between;
      .title {
        display: flex;
        align-items: center;
        font-size: 14px;
        line-height: 22px;
        font-weight: 600;
        color: #000000;
        width: fit-content;
        span {
          color: #8c8c8c;
          font-weight: 400;
          margin-left: 8px;
        }
        .title-block {
          max-width: 320px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          margin-left: 4px;
        }
      }
      .right-opration {
        display: flex;
        align-items: center;
        gap: 8px;
        line-height: 22px;
        .star-btn {
          cursor: pointer;
        }
      }
    }
    .content-box {
      color: #595959;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      white-space: pre-wrap;
      word-wrap: break-word;
      padding: 8px 16px 16px;
    }
    .fragment-img {
      padding: 0px 16px 16px;
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      img {
        width: 80px;
        height: 80px;
        border-radius: 6px;
        cursor: pointer;
      }
    }
  }
}
@keyframes flash-border {
  0%,
  100% {
    background: transparent;
  }
  50% {
    background: #c8d9f4;
  }
}

.flash-border {
  background: #c8d9f4;
  animation: flash-border 1s infinite; /* 持续时间1秒，无限次重复 */
}
.cfb363f {
  color: #fb363f !important;
}
.err-icon {
  margin-left: 4px !important;
}

.start-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #262626;
  font-size: 14px;
  .anticon {
    font-size: 16px;
  }
}

.content-container {
  position: relative;
  padding: 20px;
}

.content-block {
  margin-bottom: 16px;
  line-height: 1.6;
}

.bubble-card {
  position: absolute;
  display: inline-flex;
  padding: 2px;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  border-radius: 6px;
  background: #FFF;
  box-shadow: 0 6px 30px 5px #0000000d, 0 16px 24px 2px #0000000a, 0 8px 10px -5px #00000014;
  
  .bubble-item {
    display: flex;
    padding: 5px 16px;
    align-items: center;
    gap: 8px;
    align-self: stretch;
    border-radius: 6px;
    cursor: pointer;
    background: white;
    border: none;

    &:hover {
      background: #E4E6EB;
    }
  }
  /* 添加箭头 */
  // &::after {
  //   content: '';
  //   position: absolute;
  //   bottom: -10px;
  //   left: var(--arrow-offset, 50%);
  //   transform: translateX(-50%);
  //   border-width: 10px 10px 0;
  //   border-style: solid;
  //   border-color: white transparent transparent;
  //   filter: drop-shadow(0 2px 2px rgba(0, 0, 0, 0.1));
  // }
  
  // /* 箭头边框 */
  // &::before {
  //   content: '';
  //   position: absolute;
  //   bottom: -11px;
  //   left: var(--arrow-offset, 50%);
  //   transform: translateX(-50%);
  //   border-width: 10px 10px 0;
  //   border-style: solid;
  //   border-color: #ddd transparent transparent;
  //   z-index: -1;
  // }
}

/* 响应式调整 */
@media (max-width: 600px) {
  .bubble-card {
    max-width: 250px;
    flex-direction: column;
    gap: 4px;
  }
}

/* 为高亮和评论添加样式 */
.highlight {
  background-color: yellow;
}

.comment {
  background-color: #e6f7ff;
  border-bottom: 1px dashed #1890ff;
}
</style>
