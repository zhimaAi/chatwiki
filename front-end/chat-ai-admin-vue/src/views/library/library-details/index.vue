<template>
  <div class="details-library-page">

    <div class="between-content-box">
      <div class="layout-left">
        <div class="library-name-box">
          <img class="avatar" :src="libraryInfo.avatar || DEFAULT_LIBRARY_AVATAR" alt="" />
          <div class="name">
            {{ libraryInfo.library_name }} 
          </div>
        </div>
        <div class="left-menu-box">
          <LeftMenus></LeftMenus>
        </div>
      </div>

      <div class="right-content-box">
        <router-view></router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { EditOutlined } from '@ant-design/icons-vue'
import LeftMenus from './components/left-menus.vue'
import { getLibraryFileList, getLibraryInfo } from '@/api/library'
import { DEFAULT_LIBRARY_AVATAR } from '@/constants/index'

const rotue = useRoute()
const router = useRouter()
const query = rotue.query

const libraryInfo = ref({
  library_name: '',
  avatar: '',
  library_intro: ''
})
const getInfo = () => {
  getLibraryInfo({ id: query.id }).then((res) => {
    libraryInfo.value = res.data
  })
}
getInfo()
</script>

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 2px;
}
.layout-left {
  .library-name-box {
    display: flex;
    align-items: center;
    padding: 24px 24px 16px 24px;
    .avatar {
      width: 20px;
      height: 20px;
      margin-right: 8px;
      border-radius: 2px;
    }
    .name {
      line-height: 22px;
      font-size: 14px;
      font-weight: 600;
      color: #262626;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      max-width: 200px;
    }
  }
}
.between-content-box {
  display: flex;
  flex: 1;
  overflow: hidden;
  .left-menu-box {
    width: 232px;
    margin-right: 24px;
  }
  .right-content-box {
    flex: 1;
    padding: 24px 10px 24px 24px;
  }
}
</style>
