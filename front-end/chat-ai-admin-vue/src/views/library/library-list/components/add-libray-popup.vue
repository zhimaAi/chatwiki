<template>
  <div class="add-libray-popup" ref="modal">
    <a-modal
      :footer="null"
      :getContainer="() => $refs.modal"
      v-model:open="open"
      title=""
      :width="520"
      @ok="handleOk"
    >
      <div class="add-libray-popup-body">
        <div class="popup-title">新建知识库</div>
        <!-- 知识库类型列表 -->
        <div class="libray-type-list">
          <div class="libray-item" @click="handleSelect(0)">
            <div class="libray-item-header">
              <img class="libray-icon" src="../../../../assets/img/library/a.svg" alt="" />
              <div class="libray-type-name">普通知识库</div>
            </div>
            <div class="libray-desc">
              导入文档或网页创建知识库，知识库可以与聊天机器人或工作流关联。用户和机器人聊天时，机器人可以引用知识库内知识回复用户
            </div>

            <div class="libray-action">
              <RightOutlined />
            </div>
          </div>

          <div class="libray-item" @click="handleSelect(2)">
            <div class="libray-item-header">
              <img class="libray-icon" src="../../../../assets/img/library/q.svg" alt="" />
              <div class="libray-type-name">问答知识库</div>
            </div>
            <div class="libray-desc">
              通过excel或者word导入FAQ问答对构建知识库，适合客服场景或者需要用FAQ问答数据精确回答的场景。用户提问时，如果系统检索到与用户问题匹配度很高的问答对，会直接回复对应的答案。
            </div>
            <div class="libray-action">
              <RightOutlined />
            </div>
          </div>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { RightOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['ok'])

const open = ref(false)
const show = () => {
  open.value = true
}

const close = () => {
  open.value = false
}
const handleOk = (val) => {
  open.value = false
  emit('ok', val)
}

const handleSelect = (val) => {
  handleOk(val)
}

defineExpose({
  show,
  close
})
</script>

<style lang="less" scoped>
.add-libray-popup {
  .add-libray-popup-body {
    padding: 8px;
  }
  .popup-title {
    line-height: 32px;
    margin-bottom: 24px;
    font-size: 24px;
    font-weight: 600;
    color: #242933;
    text-align: center;
  }
  .libray-type-list {
    display: flex;
    justify-content: space-between;
    .libray-item {
      width: 216px;
      height: 362px;
      padding: 16px;
      border-radius: 16px;
      background-color: #fff;
      box-shadow: 0 8px 24px 0 rgba(50, 79, 193, 0.08);
      transition: all 0.2s;
      cursor: pointer;
      &:hover {
        box-shadow:
          0 8px 24px 0 rgba(50, 79, 193, 0.32),
          0 8px 24px 0 rgba(50, 79, 193, 0.08);
      }
    }

    .libray-icon {
      display: block;
      width: 64px;
      height: 64px;
      margin: 16px auto 16px auto;
    }

    .libray-type-name {
      line-height: 24px;
      margin-bottom: 8px;
      font-size: 16px;
      font-weight: 600;
      color: rgb(38, 38, 38);
      text-align: center;
    }

    .libray-desc {
      height: 132px;
      line-height: 22px;
      font-size: 14px;
      font-weight: 400;
      color: rgb(89, 89, 89);
      text-align: center;
    }

    .libray-action {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 32px;
      height: 32px;
      margin: 24px auto 0 auto;
      border-radius: 50%;
      color: rgba(36, 117, 252, 1);
      background: #d6e5ff;
    }
  }
  :deep(.ant-modal-content) {
    background-color: #e5efff;
    box-shadow:
      0 6px 30px 5px rgba(0, 0, 0, 0.05),
      0 16px 24px 2px rgba(0, 0, 0, 0.04),
      0 8px 10px -5px rgba(0, 0, 0, 0.08);
  }
}
</style>
