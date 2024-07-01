<template>
  <div class="empty-box">
    <img src="../../../../assets/img/library/preview/empty.png" alt="" />
    <div class="title">知识库文档为空，请添加分段</div>
    <div>
      <a-button type="primary" @click="handleOpenEditModal">
        <template #icon>
          <PlusOutlined />
        </template>
        <span>添加分段</span>
      </a-button>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue'
const emit = defineEmits(['openEditSubscription'])
const props = defineProps({
  detailsInfo: {
    type: Object
  }
})

const handleOpenEditModal = () => {
  emit('openEditSubscription', {})
}
</script>
<style lang="less" scoped>
.empty-box {
  padding-top: 148px;
  width: 100%;
  height: calc(100% - 20px);
  background: #f2f4f7;
  text-align: center;
  img {
    width: 200px;
    height: 200px;
  }
  .title {
    font-size: 16px;
    font-style: normal;
    font-weight: 600;
    line-height: 24px;
    color: #262626;
    margin-bottom: 16px;
  }
}
</style>
