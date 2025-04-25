<template>
    <div class="container">
      <a-modal
        v-model:open="visible"
        :title="modalTitle"
        @ok="handleOk"
        @cancel="handleCancel"
        width="472px"
      >
        <div class="guide-content">
          <img 
            src="@/assets/img/library/guide-image.png" 
            alt="操作引导"
            class="guide-image"
          />
          <div class="content-box">
            <div class="content-info">单次上传多个文档，上传成功后，需要手动点击文档后面 "学习" 进行学习；如果文档学习失败，支持重新学习</div>
            <a-checkbox v-model:checked="dontShowAgain">
              <div class="content-check">不再显示此提示</div>
            </a-checkbox>
          </div>
        </div>
        <template #footer>
          <a-button type="primary" @click="handleOk">我知道了</a-button>
        </template>
      </a-modal>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/stores/modules/user'

const userStore = useUserStore()
const visible = ref(false)
const dontShowAgain = ref(false)
const modalTitle = ref('')

const show = () => {
  // 延迟1秒弹出
  setTimeout(() => {
    visible.value = true
  }, 500);
}

// 处理弹窗确定按钮
const handleOk = () => {
  // 保存用户偏好
  userStore.setGuideLearningTips(dontShowAgain.value)
  visible.value = false
}

// 处理弹窗取消按钮
const handleCancel = () => {
  userStore.setGuideLearningTips(dontShowAgain.value)
  visible.value = false
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.guide-content {
  text-align: left;
  width: 100%;
  height: 333px;
}

.guide-image {
  position: absolute;
  top: 0;
  left: 0;
  max-width: 100%;
  margin-bottom: 24px;
  border-radius: 8px 8px 0 0;
}

.content-box {
  padding-top: 277px;
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;

  .content-info {
    margin-bottom: 20px;
  }
}

.content-check {
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.ant-modal-content {
  position: relative;
}
</style>