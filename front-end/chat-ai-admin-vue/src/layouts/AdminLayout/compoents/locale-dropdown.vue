<style lang="less" scoped>
.dropdown-link {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 2px 8px;
  .lang-icon {
    font-size: 16px;
  }
  .lang-name {
    line-height: 22px;
    padding: 0 8px;
    font-size: 14px;
    color: #3a4559;
  }
}
</style>

<template>
  <a-dropdown placement="left">
    <div class="dropdown-link" @click.prevent>
      <svg-icon class="lang-icon" name="lang"></svg-icon>
      <span class="lang-name">{{ selectedLocale.name }}</span>
      <RightOutlined />
    </div>
    <template #overlay>
      <a-menu @click="setLang">
        <a-menu-item v-for="item in langMap" :key="item.lang">
          <a href="javascript:;">{{ item.name }}</a>
        </a-menu-item>
      </a-menu>
    </template>
  </a-dropdown>
</template>

<script setup>
import { computed, unref } from 'vue'
import { RightOutlined } from '@ant-design/icons-vue';
import { useLocaleStore } from '@/stores/modules/locale'
import { useLocale } from '@/hooks/web/useLocale'

const localeStore = useLocaleStore()

const langMap = computed(() => localeStore.getLocaleMap)

const currentLang = computed(() => localeStore.getCurrentLocale)
const selectedLocale = computed(() => localeStore.getSelectedLocale)

const setLang = ({ key }) => {
  if (key === unref(currentLang).lang) return

  const { changeLocale } = useLocale()

  changeLocale(key).then(() => {
    // 刷新页面
    window.location.reload()
  })
}
</script>
