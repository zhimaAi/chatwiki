<template>
  <div class="subsection-box">
    <a-table :data-source="paragraphLists" :pagination="false">
      <a-table-column key="id" data-index="id" :width="1148">
        <template #title>
          <template v-if="isQaLibray"> 问答（共{{ props.total }}个） </template>
          <template v-else>分段 （共{{ props.total }}个分段）</template>
        </template>
        <template #default="{ record, index }">
          <a v-if="record.file_name" class="file-name-text" @click="toFileDetail(record)">{{ record.file_name }}</a>
          <div v-else class="file-name-text" style="color: #F40;">原文档已删除</div>
          <template v-if="isQaLibray">
            <div class="qa-list-box">
              <div class="list-item">
                <div class="list-label">问题</div>
                <div class="list-content">{{ record.question }}</div>
              </div>
              <div
                class="list-item"
                v-if="record.similar_questions && record.similar_questions.length"
              >
                <div class="list-label">相似问法</div>
                <div class="list-content">{{ record.similar_questions.join('/') }}</div>
              </div>
              <div class="list-item">
                <div class="list-label">答案</div>
                <div class="list-content">{{ record.answer }}</div>
              </div>
              <div class="fragment-img" v-viewer>
                <img v-for="(item, index) in record.images" :key="index" :src="item" alt="" />
              </div>
            </div>
          </template>
          <template v-else>
            <div class="common-list-box">
              <div class="top-block">
                <div class="title">分段{{ index + 1 }}</div>
                <div class="title">
                  {{ record.title }}
                </div>
                <span>共{{ record.word_total }}个字符</span>
              </div>
              <div class="content-box" v-html="record.content"></div>
              <div class="fragment-img" v-viewer>
                <img v-for="(item, index) in record.images" :key="index" :src="item" alt="" />
              </div>
            </div>
          </template>
        </template>
      </a-table-column>
      <a-table-column title="嵌入状态" key="status_text" data-index="status_text" :width="128">
        <template #default="{ record }">
          <span class="status-tag" v-if="record.status == 0"><ClockCircleFilled /> 未转换</span>
          <span class="status-tag complete" v-if="record.status == 1"
            ><CheckCircleFilled /> 已转换</span
          >
          <a-tooltip placement="top" v-if="record.status == 2">
            <template #title>
              <span>{{ record.errmsg }}</span>
            </template>
            <span>
              <span class="status-tag status-error"><CloseCircleFilled /> 转换异常</span>
            </span>
          </a-tooltip>
          <span class="status-tag running" v-if="record.status == 3"
            ><a-spin size="small" /> 转换中</span
          >
        </template>
      </a-table-column>
      <a-table-column title="操作" key="action" data-index="action" :width="120">
        <template #default="{ record, index }">
          <div class="right-opration">
            <a-dropdown>
              <div class="hover-btn-box">
                <a-popconfirm
                  v-if="record.category_id > 0"
                  title="是否取消该标记？"
                  ok-text="确定"
                  cancel-text="取消"
                  @confirm="handleSetCategory(record, {})"
                >
                  <StarFilled :style="{ color: getColor(record), 'font-size': '16px' }" />
                </a-popconfirm>
                <StarOutlined
                  @click="handleSetCategory(record, startLists[0])"
                  style="font-size: 16px"
                  v-else
                />
              </div>
              <template #overlay>
                <a-menu>
                  <a-menu-item v-for="star in startLists" :key="star.id">
                    <div class="start-item" @click="handleSetCategory(record, star)">
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
              <div class="hover-btn-box" @click="toReSegmentationPage(record, index)">
                <SyncOutlined />
              </div>
            </a-tooltip>
            <a-dropdown placement="bottomRight">
              <div class="hover-btn-box">
                <MoreOutlined />
              </div>
              <template #overlay>
                <a-menu>
                  <a-menu-item>
                    <div @click.stop="handleOpenEditModal(record)">编辑</div>
                  </a-menu-item>
                  <a-menu-item>
                    <div @click.stop="hanldleDelete(record)">删除</div>
                  </a-menu-item>
                </a-menu>
              </template>
            </a-dropdown>
          </div>
        </template>
      </a-table-column>
    </a-table>
    <ClassificationMarkModal @ok="getCategoryLists" ref="classificationMarkModalRef" />
  </div>
</template>
<script setup>
import { reactive, ref, computed, createVNode } from 'vue'
import { message } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import {
  ExclamationCircleOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  MoreOutlined,
  SyncOutlined,
  ClockCircleFilled,
  LoadingOutlined,
  StarOutlined,
  StarFilled,
  SettingOutlined
} from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import ClassificationMarkModal from '@/views/library/library-preview/components/classification-mark-modal.vue'
import colorLists from '@/utils/starColors.js'
import {
  deleteParagraph,
  editParagraph,
  getCategoryList,
  updateParagraphCategory,
  saveCategoryParagraph
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
  total: {
    type: [Number, String],
    default: 0
  },
  isQaLibray: {
    default: false
  }
})

const route = useRoute()
const query = route.query

const toReSegmentationPage = (item, index) => {
  let { id, title, content, question, answer, images, category_id, library_id } = item
  let similar_questions = item.similar_questions || []
  Modal.confirm({
    title: '重新转换确认',
    icon: null,
    content: `确定要重新转换【分段${index + 1}】吗?`,
    onOk() {
      return new Promise((resolve, reject) => {
        saveCategoryParagraph({
          id,
          library_id,
          title,
          content,
          question,
          answer,
          images,
          category_id,
          similar_questions: JSON.stringify(similar_questions)
        })
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
const hanldleDelete = (record) => {
  Modal.confirm({
    title: '提示',
    icon: createVNode(ExclamationCircleOutlined),
    content: '是否取消该标记?',
    onOk() {
      handleSetCategory(record, {})
    },
    onCancel() {}
  })
}

const startLists = ref([])
const getCategoryLists = () => {
  getCategoryList({library_id: query.id}).then((res) => {
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
    message.success('修改成功')
    if (!star.id) {
      emit('handleDelParagraph', item.id)
    } else {
      item.category_id = star.id
    }
    getCategoryLists()
  })
}

const toFileDetail = (record) => {
  window.open(`/#/library/preview?id=${record.file_id}`)
}

defineExpose({ handleOpenEditModal })
</script>
<style lang="less" scoped>
.subsection-box {
  width: 100%;
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
}

.qa-list-box {
  margin-top: 4px;
  .list-item {
    display: flex;
    flex-wrap: wrap;
    line-height: 22px;
    font-size: 14px;
    color: #595959;
    margin-bottom: 6px;
    .list-label {
      width: 65px;
      font-weight: 600;
      color: #262626;
    }
    .list-content {
      flex: 1;
    }
  }
}
.fragment-img {
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
.common-list-box {
  margin-top: 4px;
  .top-block {
    display: flex;
    align-items: center;
    gap: 4px;
    color: #000000;
    .title {
      font-weight: 600;
    }
    span {
      color: #8c8c8c;
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
.status-tag {
  display: inline-block;
  height: 24px;
  line-height: 24px;
  padding: 0 6px;
  border-radius: 2px;
  font-size: 14px;
  font-weight: 500;
  text-align: center;
  color: #3a4559;
  background-color: #edeff2;

  &.running {
    color: #2475fc;
    background-color: #e8effc;
  }

  &.complete {
    color: #21a665;
    background: #e8fcf3;
  }

  &.error {
    cursor: pointer;
    color: #fb363f;
    background-color: #f5c6c8;
  }

  &.warning {
    cursor: pointer;
    background: #faebe6;
    color: #ed744a;
  }

  &.status-error {
    cursor: pointer;
    color: #fb363f;
  }
}

.right-opration {
  display: flex;
  align-items: center;
  gap: 8px;
}
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
