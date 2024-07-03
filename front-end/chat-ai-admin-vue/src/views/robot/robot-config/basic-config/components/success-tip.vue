<style lang="less" scoped>
.success-box {
  width: 478px;
  padding-top: 80px;
  margin: 0 auto;
  text-align: center;
}

.success-icon {
  font-size: 62px;
  color: #fff;
}

.success-label {
  line-height: 22px;
  margin-top: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #000000;
}

.success-desc {
  line-height: 22px;
  margin-top: 4px;
  font-size: 14px;
  color: #8c8c8c;
}

.action-box {
  margin-top: 16px;

  .action-btn {
    margin: 0 4px;
  }
}
</style>

<template>
  <div class="success-box">
    <svg-icon name="success" class="success-icon"></svg-icon>
    <div class="success-label">创建成功</div>
    <div class="success-desc">
      新建/编辑机器人设置之后，建议您对机器人效果进行模拟测试，并根据模拟测试结果适当调整提示词和知识库内容
    </div>

    <div class="action-box">
      <a-button class="action-btn" @click="toLlistPage">返回列表</a-button>
      <a-button class="action-btn" type="primary" @click="toTestPage">去测试</a-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
const emit = defineEmits(['toTest'])
const router = useRouter()

const toLlistPage = () => {
  router.replace('/robot/list')
}

const toTestPage = () => {
  emit('toTest')
}
</script>
