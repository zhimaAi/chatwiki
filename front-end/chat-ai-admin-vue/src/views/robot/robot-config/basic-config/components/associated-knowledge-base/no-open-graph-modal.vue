<template>
  <div>
    <a-modal v-model:open="open" title="部分知识库未开启生成知识图谱" :footer="null" :width="746">
      <a-alert
        style="margin-top: 16px"
        message="当前机器人关联的部分知识库未开启知识图谱,无法通过知识图谱检索,可能会影响检索结果。请您按需开启对应知识的生成知识图谱功能"
      ></a-alert>
      <cu-scroll :scrollbar="false" style="height: 400px">
        <div class="list-box">
          <div
            class="list-item"
            @click="toEditPage(item)"
            v-for="item in props.list"
            :key="item.id"
          >
            <div class="avater-box">
              <img :src="item.avatar" alt="" />
            </div>
            <div class="content-box">
              <div class="title-block">
                <div class="tag-icon">Graph</div>
                <div class="title-text">{{ item.library_name }}</div>
              </div>
              <div class="desc-box">{{ item.library_intro }}</div>
            </div>
            <div class="right-icon">
              <RightOutlined />
            </div>
          </div>
        </div>
      </cu-scroll>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { RightOutlined } from '@ant-design/icons-vue'

const emit = defineEmits(['refreshList'])
const props = defineProps({
  list: {
    type: Array,
    default: []
  }
})
const open = ref(false)

watch(
  () => props.list,
  () => {
    if (props.list.length == 0) {
      open.value = false
    }
  },
  {}
)

const show = () => {
  open.value = true
}

const handleVisibilityChange = () => {
  let isVisible = document.visibilityState === 'visible'
  if (isVisible) {
    emit('refreshList')
  }
}

onMounted(() => {
  document.addEventListener('visibilitychange', handleVisibilityChange)
})

onUnmounted(() => {
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})

const toEditPage = (item) => {
  window.open(`/#/library/details/knowledge-config?id=${item.id}`)
}

defineExpose({
  show
})
</script>

<style lang="less" scoped>
.list-box {
  margin-top: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  .list-item {
    width: calc(50% - 8px);
    cursor: pointer;
    padding: 16px;
    border: 1px solid var(---, #fb363f);
    border-radius: 6px;
    display: flex;
    align-items: center;
    .avater-box {
      width: 40px;
      height: 40px;
      border-radius: 6px;
      margin-right: 8px;
      img {
        width: 100%;
        height: 100%;
      }
    }
    .content-box {
      flex: 1;
      overflow: hidden;
      .title-block {
        display: flex;
        align-items: center;
        gap: 4px;
      }
      .tag-icon {
        width: 42px;
        height: 18px;
        display: flex;
        align-items: center;
        justify-content: center;
        border: 1px solid #00000026;
        background: #0000000a;
        border-radius: 6px;
        font-size: 12px;
        color: #bfbfbf;
      }
      .title-text {
        max-width: 200px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-size: 14px;
        font-weight: 600;
        line-height: 22px;
        color: #262626;
      }
      .desc-box {
        margin-top: 2px;
        max-width: 250px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        font-size: 12px;
        line-height: 20px;
        color: #8c8c8c;
      }
    }
    .right-icon {
      font-size: 16px;
      color: #8c8c8c;
    }
  }
}
</style>
