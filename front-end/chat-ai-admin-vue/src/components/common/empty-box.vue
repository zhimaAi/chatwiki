<template>
  <div class="empty-box">
    <slot name="image">
      <img src="@/assets/empty.png"/>
    </slot>
    <div class="title">
      <slot name="title">{{ title }}</slot>
    </div>
    <div class="desc">
      <slot name="desc">{{ desc }}</slot>
    </div>
  </div>
</template>

<script setup>
import {ref} from 'vue';

const props = defineProps({
  title: {
    type: String,
    default: '暂无数据'
  },
  desc: {
    type: String,
    default: ''
  }
})

</script>

<style scoped lang="less">
.empty-box {
  height: 100%;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  img {
    width: 200px;
    height: 200px;
  }

  .title {
    color: #262626;
    font-size: 16px;
    font-weight: 600;
    line-height: 24px;
  }

  .desc {
    color: #595959;
    font-size: 14px;
    font-weight: 400;
    margin-top: 12px;
  }
}
</style>
