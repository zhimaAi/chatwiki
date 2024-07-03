<style lang="less" scoped>
.ignore-icons-page {
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
        color: #262626;
        font-size: 20px;
      }

      .icon-name {
        font-size: 14px;
      }
    }
  }
}
</style>

<template>
  <div class="ignore-icons-page">
    <div class="icon-list">
      <div class="icon-item" v-for="item in list" :key="item.name" @click="handleClick(item)">
        <div class="icon-example">
          <svg-icon :name="item.name" />
        </div>
        <div class="icon-name">{{ item.name }}</div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { showDialog } from 'vant';
import { useClipboard } from '@/hooks/web/useClipboard'
import { getAllIcons } from '@/assets/icons/helper'

const { copy } = useClipboard()
let list = getAllIcons()

function handleClick(item) {
  let text = `<svg-icon name="${item.name}" style="font-size: 14px;color: #333;" />`

  copy(text)

  showDialog({
    message: `<svg-icon name="${item.name}" /> copied`,
  }).then(() => {
    // on close
  });
}
</script>
