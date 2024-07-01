<style lang="less" scoped>
.locale-dropdown {
  .dropdown-link {
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    height: 32px;
    padding: 0 16px;
    border-radius: 16px;
    border: 1px solid #c3cbd9;
    .lang-icon {
      width: 16px;
    }
    .lang-name {
      line-height: 22px;
      padding: 0 5px;
      font-size: 14px;
      color: #3a4559;
    }
  }
}
</style>

<template>
  <div class="locale-dropdown">
    <a-dropdown>
      <div class="dropdown-link" @click.prevent>
        <img class="lang-icon" src="../../../assets/img/lang.png" alt="" />
        <span class="lang-name">{{ selectedLocale.name }}</span
        ><svg-icon name="arrow-down" style="font-size: 16px; color: #8c8c8c"></svg-icon>
      </div>
      <template #overlay>
        <a-menu @click="setLang">
          <a-menu-item v-for="item in langMap" :key="item.lang">
            <a href="javascript:;">{{ item.name }}</a>
          </a-menu-item>
        </a-menu>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup>
import { computed, unref } from 'vue'
import { useLocaleStore } from '@/stores/modules/locale'
import { useLocale } from '@/hooks/web/useLocale'

const localeStore = useLocaleStore()

const langMap = computed(() => localeStore.getLocaleMap)

const currentLang = computed(() => localeStore.getCurrentLocale)
const selectedLocale = computed(() => localeStore.getSelectedLocale)

const setLang = ({ key }) => {
  if (key === unref(currentLang).lang) return

  const { changeLocale } = useLocale()

  changeLocale(key)

  // 需要重新加载页面让整个语言多初始化
  // window.location.reload()
}
</script>
