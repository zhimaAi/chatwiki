<style lang="less" scoped>
.icons-page {
  .icon-list {
    display: flex;
    flex-flow: row wrap;
    gap: 16px;

    .icon-item {
      width: 200px;
      height: 120px;
      padding: 16px;
      text-align: center;
      color: #333;
      background-color: #fff;
      cursor: pointer;

      &:hover {
        color: #fff;
        background-color: #1677ff;
      }

      .icon-example {
        font-size: 20px;
      }

      .icon-name {
        font-size: 14px;
      }
    }
  }
  .search-box{
    padding: 32px;
    text-align: center;
  }
}
</style>

<template>
  <div class="icons-page">
    <!-- 搜索 -->
     <div class="search-box">
      <a-input-search
        v-model:value="searchText"
        placeholder="请输入图标名称搜索"
        style="width: 400px; margin-bottom: 16px;"
        @search="onSearch"
      />
     </div>
    <div class="icon-list">
      <data class="icon-item" v-for="item in list" :key="item.name" @click="handleClick(item)">
        <div class="icon-example">
          <svg-icon :name="item.name"></svg-icon>
        </div>
        <div class="icon-name">{{ item.name }}</div>
      </data>
    </div>
  </div>
</template>
<script setup>
import { message } from 'ant-design-vue'
import { ref } from 'vue'
import { copyText } from '@/utils/index'
import { getAllIcons } from '@/assets/svg/helper'

const originalList = getAllIcons()
const list = ref(originalList)
const searchText = ref('')

function onSearch(value) {
  if (!value) {
    list.value = originalList
    return
  }
  list.value = originalList.filter(item => 
    item.name.toLowerCase().includes(value.toLowerCase())
  )
}

function handleClick(item) {
  let text = `<svg-icon name="${item.name}" style="font-size: 14px;color: #333;"></svg-icon>`
  copyText(text)
  message.success(`<svg-icon name="${item.name}"></svg-icon> copied`)
}
</script>
