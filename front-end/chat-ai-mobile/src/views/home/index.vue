<script setup lang="ts">
import { showDialog } from 'vant'
import { useLocale } from '@/hooks/web/useLocale'

const { changeLocale } = useLocale()

function showAlert() {
  showDialog({
    title: 'Title',
    message: 'The code is written for people to see and can be run on a machine.'
  }).then(() => {
    // on close
  })
}

function setLang(lang: LocaleType) {
  changeLocale(lang)
}

</script>

<template>
  <div class="main">
    <van-button type="primary" @click="showAlert">showAlert</van-button>
    <van-button type="success" @click="setLang('en-US')">en-US</van-button>
    <van-button type="default" @click="setLang('zh-CN')">zh-CN</van-button>
  </div>
</template>

<style>
.main {
  width: 300px;
  margin: 0 auto;
}
</style>
