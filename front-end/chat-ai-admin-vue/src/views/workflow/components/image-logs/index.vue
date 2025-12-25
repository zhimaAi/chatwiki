<template>
  <div class="preview-img-box">
    <a-carousel arrows :dots="false">
      <template #prevArrow>
        <div class="custom-slick-arrow" style="left: 10px; z-index: 1">
          <LeftCircleOutlined />
        </div>
      </template>
      <template #nextArrow>
        <div class="custom-slick-arrow" style="right: 10px">
          <RightCircleOutlined />
        </div>
      </template>
      <div class="img-list-item" v-for="(item, index) in props.currentImageList" :key="item">
        <img @click="handleShowImg(index)" :src="item" alt="" />
      </div>
    </a-carousel>
    <div style="display: none">
      <a-image-preview-group
        :preview="{ visible, current, onVisibleChange: (vis) => (visible = vis) }"
      >
        <a-image v-for="item in props.currentImageList" :key="item" :src="item" />
      </a-image-preview-group>
    </div>
  </div>
</template>

<script setup>
import { LeftCircleOutlined, RightCircleOutlined } from '@ant-design/icons-vue'
import { ref } from 'vue'

const props = defineProps({
  currentImageList: {
    type: Array,
    default: () => []
  }
})

const visible = ref(false)
const current = ref(0)

const handleShowImg = (index) => {
  current.value = index
  visible.value = true
}
</script>

<style lang="less" scoped>
.preview-img-box {
  &::v-deep(.ant-carousel) {
    .slick-slide {
      text-align: center;
      height: 220px;
    }

    .slick-arrow.custom-slick-arrow {
      width: 25px;
      height: 25px;
      font-size: 25px;
      color: #fff;
      background-color: rgba(31, 45, 61, 0.11);
      transition: ease all 0.3s;
      opacity: 0.3;
      z-index: 1;
    }
    .slick-arrow.custom-slick-arrow:before {
      display: none;
    }
    .slick-arrow.custom-slick-arrow:hover {
      color: #fff;
      opacity: 0.5;
    }

    .slick-slide .img-list-item {
      color: #fff;
      width: 100%;
      height: 100%;
      img {
        width: 100%;
      }
    }
  }
}
</style>
