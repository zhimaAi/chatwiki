<template>
  <div class="details-library-page">
    <div class="between-content-box">
      <div class="layout-left">
        <div class="library-name-box">
          <img class="avatar" src="@/assets/img/database/base-icon.svg" alt="" />
          <div class="name">
            {{ dataBaseInfo.name }}
          </div>
          <svg-icon @click="toEdit" name="edit"></svg-icon>
        </div>
        <div class="left-menu-box">
          <LeftMenus></LeftMenus>
        </div>
      </div>

      <div class="right-content-box">
        <router-view></router-view>
      </div>
    </div>
    <AddDataSheet @ok="getInfo" ref="addDataSheetRef"></AddDataSheet>
  </div>
</template>

<script>
import { useRoute } from 'vue-router'
import { EditOutlined } from '@ant-design/icons-vue'
import LeftMenus from './components/left-menus.vue'
import { useDatabaseStore } from '@/stores/modules/database'
import { computed, ref, defineComponent } from 'vue'
import AddDataSheet from '../database-list/components/add-data-sheet.vue'
import { getDatabasePermission } from '@/utils/permission'

export default defineComponent({
  components: {
    LeftMenus,
    AddDataSheet,
    EditOutlined
  },

  async beforeRouteEnter(to, from, next) {
    let key = getDatabasePermission(to.query.form_id)
    if (key == 0 || key == 1) {
      next(`/no-permission`)
      return
    }
    next()
  },

  setup() {
    const databaseStore = useDatabaseStore()
    const rotue = useRoute()
    const query = rotue.query

    const dataBaseInfo = computed(() => {
      return databaseStore.databaseInfo
    })
    const getInfo = () => {
      databaseStore.getDatabaseInfo({ id: query.form_id })
    }
    getInfo()

    const addDataSheetRef = ref(null)
    const toEdit = () => {
      addDataSheetRef.value.show({
        ...databaseStore.databaseInfo
      })
    }

    return {
      dataBaseInfo,
      getInfo,
      toEdit,
      addDataSheetRef
    }
  }
})
</script>

<!-- <script setup>
import { useRoute } from 'vue-router'
import { EditOutlined } from '@ant-design/icons-vue'
import LeftMenus from './components/left-menus.vue'
import { useDatabaseStore } from '@/stores/modules/database'
import { computed, ref } from 'vue'
import AddDataSheet from '../database-list/components/add-data-sheet.vue'
const databaseStore = useDatabaseStore()
const rotue = useRoute()
const query = rotue.query

const dataBaseInfo = computed(() => {
  return databaseStore.databaseInfo
})
const getInfo = () => {
  databaseStore.getDatabaseInfo({ id: query.form_id })
}
getInfo()

const addDataSheetRef = ref(null)
const toEdit = () => {
  addDataSheetRef.value.show({
    ...databaseStore.databaseInfo
  })
}
</script> -->

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 2px;
}
.layout-left {
  width: 256px;
  border-right: 1px solid rgba(5, 5, 5, 0.06);
  .library-name-box {
    display: flex;
    align-items: center;
    padding: 24px 24px 16px 24px;
    gap: 8px;
    .avatar {
      width: 20px;
      height: 20px;
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
      max-width: 158px;
    }
    .svg-action {
      cursor: pointer;
    }
  }
}
.between-content-box {
  display: flex;
  flex: 1;
  overflow: hidden;
  .left-menu-box {
    width: 100%;
    ::v-deep(.ant-menu-vertical) {
      border-right: none;
    }
  }
  .right-content-box {
    flex: 1;
    overflow: hidden;
    padding: 24px;
    padding-right: 0;
  }
}
</style>
