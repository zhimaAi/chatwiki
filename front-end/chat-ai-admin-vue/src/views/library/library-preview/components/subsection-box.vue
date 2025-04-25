<template>
  <div class="subsection-box">
    <div class="subsection-box-title">
      分段预览
      <span>共{{ props.total }}个分段</span>
    </div>
    <div
      class="list-item"
      v-for="(item, index) in props.paragraphLists"
      :key="item.id"
      @click="handleToTargetPage(item, index)"
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
                @click="handleSetCategory(item, {})"
                :style="{ color: getColor(item), 'font-size': '16px' }"
                v-if="item.category_id > 0"
              />
              <StarOutlined
                @click="handleSetCategory(item, startLists[0])"
                style="font-size: 16px"
                v-else
              />
            </div>

            <template #overlay>
              <a-menu>
                <a-menu-item v-for="star in startLists" :key="star.id">
                  <div class="start-item" @click="handleSetCategory(item, star)">
                    <StarFilled :style="{ color: colorLists[star.type] }" />
                    <div>{{ star.name || '-' }}</div>
                  </div>
                </a-menu-item>
                <a-menu-item>
                  <div class="start-item" @click="handleOpenSetStartModal">
                    <SettingOutlined />
                    <div>标记设置</div>
                  </div>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
          <a-tooltip>
            <template #title>重新转换</template>
            <div class="hover-btn-box">
              <SyncOutlined @click="toReSegmentationPage(item, index)" />
            </div>
          </a-tooltip>
          <a-dropdown placement="bottomRight">
            <div class="hover-btn-box">
              <MoreOutlined />
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
      <div class="content-box" v-html="item.content"></div>
      <div class="fragment-img" v-viewer>
        <img v-for="(item, index) in item.images" :key="index" :src="item" alt="" />
      </div>
    </div>
    <ClassificationMarkModal @ok="getCategoryLists" ref="classificationMarkModalRef" />
  </div>
</template>
<script setup>
import { reactive, ref, computed, createVNode } from 'vue'
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
  'getStatrList'
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

const toReSegmentationPage = (item, index) => {
  let { id, title, content, question, answer, images } = item
  let similar_questions = item.similar_questions || []
  Modal.confirm({
    title: '重新转换确认',
    icon: null,
    content: `确定要重新转换【分段${index + 1}】吗?`,
    onOk() {
      return new Promise((resolve, reject) => {
        editParagraph({ id, title, content, question, answer, images, similar_questions: JSON.stringify(similar_questions)  })
          .then((res) => {
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
      return new Promise((resolve, reject) => {
        deleteParagraph({ id })
          .then((res) => {
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

const handleToTargetPage = (item, index) => {
  emit('handleScrollTargetPage', {
    page_num: item.page_num,
    index
  })
}

const startLists = ref([])
const getCategoryLists = () => {
  getCategoryList().then((res) => {
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
    return '#F4EA2A'
  }
}
const handleSetCategory = (item, star = {}) => {
  updateParagraphCategory({
    id: item.id,
    category_id: star.id || 0
  }).then((res) => {
    message.success('设置成功')
    item.category_id = star.id
  })
}

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
  .list-item {
    margin-top: 8px;
    width: 100%;
    background: #fff;
    border-radius: 6px;
    padding: 16px;
    .top-block {
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
      margin-top: 8px;
      white-space: pre-wrap;
      word-wrap: break-word;
    }
    .fragment-img {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-top: 8px;
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
</style>
