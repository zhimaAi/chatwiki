<template>
  <div class="subsection-box">
    <!-- <div class="subsection-box-title">
      分段预览
      <span>共{{ props.total }}个分段</span>
    </div> -->
    <a-table :data-source="paragraphLists" :pagination="false">
      <a-table-column key="id" data-index="id" :width="1148">
        <template #title>问答（共{{ props.total }}个）</template>
        <template #default="{ record }">
          <div class="qa-list-box">
            <div class="list-item">
              <div class="list-label">问题</div>
              <div class="list-content">{{ record.question }}</div>
            </div>
            <div class="list-item" v-if="record.similar_questions && record.similar_questions.length">
              <div class="list-label">相似问法</div>
              <div class="list-content">{{record.similar_questions.join('/')}}</div>
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
                <StarFilled
                  @click="handleSetCategory(record, {})"
                  :style="{ color: getColor(record), 'font-size': '16px' }"
                  v-if="record.category_id > 0"
                />
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
                    <div @click.stop="hanldleDelete(record.id)">删除</div>
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
import ClassificationMarkModal from './classification-mark-modal.vue'
import colorLists from '@/utils/starColors.js'
import {
  deleteParagraph,
  editParagraph,
  getCategoryList,
  updateParagraphCategory
} from '@/api/library'
import { useRoute } from 'vue-router'

const route = useRoute()
const query = route.query
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
  }
})

const toReSegmentationPage = (item, index) => {
  let { id, title, content, question, answer, images, category_id } = item
  let similar_questions = item.similar_questions || []
  Modal.confirm({
    title: '重新转换确认',
    icon: null,
    content: `确定要重新转换【分段${index + 1}】吗?`,
    onOk() {
      return new Promise((resolve, reject) => {
        editParagraph({ id, title, content, question, answer, images, category_id, similar_questions: JSON.stringify(similar_questions) })
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

const startLists = ref([])
const getCategoryLists = () => {
  let params = {}
  if (query.id) {
    params.file_id = query.id
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
    getCategoryLists()
  })
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
  .fragment-img {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding-left: 65px;
    img {
      width: 80px;
      height: 80px;
      border-radius: 6px;
      cursor: pointer;
    }
  }
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
