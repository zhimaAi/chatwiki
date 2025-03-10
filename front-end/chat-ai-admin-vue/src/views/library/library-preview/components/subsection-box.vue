<template>
  <div class="subsection-box">
    <div class="subsection-box-title">
      分段预览
      <span>共{{ props.total }}个分段</span>
    </div>
    <div
      @click="handleToTargetPage(item.page_num)"
      class="list-item"
      v-for="item in props.paragraphLists"
      :key="item.id"
    >
      <div class="top-block">
        <div class="title">
          id：{{ item.id }}
          {{ item.title }}
          <span>共{{ item.word_total }}个字符</span>
          <span>状态：{{ item.status_text }}</span>
        </div>
        <div class="right-opration">
          <a @click.stop="handleOpenEditModal(item)">编辑</a>
          <span class="v-line"></span>
          <a @click.stop="hanldleDelete(item.id)">删除</a>
        </div>
      </div>
      <div class="content-box" v-if="item.question">Q：{{ item.question }}</div>
      <div class="content-box" v-if="item.answer">A：{{ item.answer }}</div>
      <div class="content-box" v-html="item.content"></div>
      <div class="fragment-img" v-viewer>
        <img v-for="(item, index) in item.images" :key="index" :src="item" alt="" />
      </div>
    </div>
  </div>
</template>
<script setup>
import { reactive, ref, computed, createVNode } from 'vue'
import { message } from 'ant-design-vue'
import { ExclamationCircleOutlined } from '@ant-design/icons-vue'
import { Modal } from 'ant-design-vue'
import { deleteParagraph } from '@/api/library'

const emit = defineEmits(['handleDelParagraph', 'handleScrollTargetPage', 'openEditSubscription'])
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

const handleToTargetPage = (page_num) => {
  emit('handleScrollTargetPage', page_num)
}
defineExpose({ handleOpenEditModal })
</script>
<style lang="less" scoped>
.subsection-box {
  background: #f2f4f7;
  padding: 14px 16px;
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
  .list-item {
    margin-top: 8px;
    width: 100%;
    background: #fff;
    border-radius: 2px;
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
        span {
          color: #8c8c8c;
          font-weight: 400;
          margin-left: 8px;
        }
      }
      .right-opration {
        display: flex;
        align-items: center;
        line-height: 22px;
        .v-line {
          display: block;
          width: 1px;
          height: 12px;
          background: #d8dde6;
          margin: 0 16px;
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
    background: #C8D9F4;
  }
}

.flash-border {
  background: #C8D9F4;
  animation: flash-border 1s infinite; /* 持续时间1秒，无限次重复 */
}
</style>
